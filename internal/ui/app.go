// Package ui view initial
package ui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/joalvarez/toolkit/internal/servers"
)

type navigationState struct {
	node   *MenuNode
	cursor int
}

type Model struct {
	width   int
	height  int
	root    *MenuNode
	current *MenuNode
	history []navigationState
	cursor  int

	listedServersOutput string
	listErr             error

	serverLookupInput string
	serverLookupErr   error
	selectedServer    *servers.Server
}

type serversListedMsg struct {
	output string
	err    error
}

type serverInfoMsg struct {
	query  string
	server *servers.Server
	err    error
}

func New() Model {
	return Model{
		root:    RootMenu,
		current: RootMenu,
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

	case serversListedMsg:
		m.listedServersOutput = msg.output
		m.listErr = msg.err

	case serverInfoMsg:
		m.serverLookupInput = msg.query
		m.selectedServer = msg.server
		m.serverLookupErr = msg.err

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "up", "k":
			if cmd, handled := m.handleActionKey(msg); handled {
				return m, cmd
			}
			m.moveSelection(-1)
		case "down", "j":
			if cmd, handled := m.handleActionKey(msg); handled {
				return m, cmd
			}
			m.moveSelection(1)
		case "enter":
			if cmd, handled := m.handleActionKey(msg); handled {
				return m, cmd
			}
			return m, m.enter()
		case "esc", "backspace":
			if cmd, handled := m.handleActionKey(msg); handled {
				return m, cmd
			}
			m.goBack()
		default:
			if cmd, handled := m.handleActionKey(msg); handled {
				return m, cmd
			}
		}

	}

	return m, nil
}

func (m Model) View() string {
	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("205"))

	bodyWidth := m.bodyWidth()
	bodyInnerWidth := bodyWidth - 2

	headingStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("86")).
		Width(bodyInnerWidth).
		Align(lipgloss.Center)

	itemStyle := lipgloss.NewStyle().
		PaddingLeft(1)

	selectedStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("230")).
		Background(lipgloss.Color("62")).
		Padding(0, 1)

	hintStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("241")).
		Width(bodyInnerWidth).
		Align(lipgloss.Center)

	body := lipgloss.NewStyle().
		Width(bodyWidth).
		Padding(1, 1).
		Border(lipgloss.RoundedBorder()).
		Render(strings.Join([]string{
			headingStyle.Render(m.currentTitle()),
			headingStyle.Render(m.currentSubtitle()),
			"",
			m.renderContent(itemStyle, selectedStyle),
			lipgloss.NewStyle().
				Align(lipgloss.Center).
				Render(),
			"",
			hintStyle.Render(m.currentHint()),
		}, "\n"))

	content := lipgloss.JoinVertical(
		lipgloss.Center,
		titleStyle.Render(title()),
		"",
		body,
	)

	return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, content)
}

func (m Model) bodyWidth() int {
	width := 55
	if m.current != nil && m.current.Action == "servers.list" {
		width = 84
	}

	if m.width <= 0 {
		return width
	}

	available := m.width - 4
	if available < 40 {
		return 40
	}
	if available < width {
		return available
	}

	return width
}

func (m *Model) moveSelection(delta int) {
	items := m.currentItems()
	if len(items) == 0 {
		return
	}

	m.cursor = (m.cursor + delta + len(items)) % len(items)
}

func (m Model) currentTitle() string {
	if m.current == nil {
		return "Menu"
	}

	return m.current.Label
}

func (m Model) currentNodes() []*MenuNode {
	if m.current == nil {
		return nil
	}

	return m.current.Children
}

func (m Model) currentItems() []string {
	nodes := m.currentNodes()
	items := make([]string, 0, len(nodes))
	for _, node := range nodes {
		items = append(items, node.Label)
	}

	return items
}

func (m Model) currentHint() string {
	if m.current != nil && m.current.IsAction() {
		return "esc/backspace go back • q quit"
	}

	nav := "↑/↓ or j/k move • enter open • q quit"

	if len(m.history) > 0 {
		nav = "↑/↓ or j/k move • enter open • esc/backspace go back • q quit"
	}

	return nav
}

func (m Model) renderItems(itemStyle, selectedStyle lipgloss.Style) string {
	items := m.currentItems()
	if len(items) == 0 {
		return "No items available"
	}

	rendered := make([]string, 0, len(items))
	for i, item := range items {
		line := "  " + item
		if i == m.cursor {
			line = selectedStyle.Render("› " + item)
		} else {
			line = itemStyle.Render(line)
		}
		rendered = append(rendered, line)
	}

	return strings.Join(rendered, "\n")
}

func (m *Model) enter() tea.Cmd {
	node := m.selectedNode()
	if node == nil {
		return nil
	}

	if node.HasChildren() {
		m.history = append(m.history, navigationState{node: m.current, cursor: m.cursor})
		m.current = node
		m.cursor = 0
		return nil
	}

	if node.IsAction() {
		m.history = append(m.history, navigationState{node: m.current, cursor: m.cursor})
		m.current = node
		m.cursor = 0
		m.prepareAction(node.Action)
		return m.runAction(node.Action)
	}

	return nil
}

func (m *Model) goBack() {
	if len(m.history) == 0 {
		return
	}

	previous := m.history[len(m.history)-1]
	m.history = m.history[:len(m.history)-1]
	m.current = previous.node
	m.cursor = previous.cursor
}

func (m Model) selectedNode() *MenuNode {
	nodes := m.currentNodes()
	if len(nodes) == 0 || m.cursor < 0 || m.cursor >= len(nodes) {
		return nil
	}

	return nodes[m.cursor]
}

func (m Model) currentSubtitle() string {
	if node := m.selectedNode(); node != nil {
		return node.Hint
	}

	if m.current != nil {
		return m.current.Hint
	}

	return ""
}

func (m Model) renderContent(itemStyle, selectedStyle lipgloss.Style) string {
	if m.current != nil && m.current.IsAction() {
		switch m.current.Action {
		case "servers.list":
			if m.listErr != nil {
				return itemStyle.Render("Error listing servers: " + m.listErr.Error())
			}

			if strings.TrimSpace(m.listedServersOutput) == "" {
				return itemStyle.Render("Loading servers...")
			}

			lines := strings.Split(m.listedServersOutput, "\n")
			for i, line := range lines {
				lines[i] = itemStyle.Render(line)
			}

			return strings.Join(lines, "\n")
		case "servers.show":
			return m.renderServerInfo(itemStyle)
		}
	}
	return m.renderItems(itemStyle, selectedStyle)
}

func (m *Model) prepareAction(action string) {
	switch action {
	case "servers.show":
		m.serverLookupInput = ""
		m.serverLookupErr = nil
		m.selectedServer = nil
	}
}

func (m *Model) handleActionKey(msg tea.KeyMsg) (tea.Cmd, bool) {
	if m.current == nil || !m.current.IsAction() {
		return nil, false
	}

	switch m.current.Action {
	case "servers.show":
		switch msg.Type {
		case tea.KeyRunes:
			m.serverLookupInput += string(msg.Runes)
			m.serverLookupErr = nil
			m.selectedServer = nil
			return nil, true
		case tea.KeySpace:
			m.serverLookupInput += " "
			m.serverLookupErr = nil
			m.selectedServer = nil
			return nil, true
		case tea.KeyBackspace:
			if m.serverLookupInput == "" {
				return nil, false
			}

			runes := []rune(m.serverLookupInput)
			m.serverLookupInput = string(runes[:len(runes)-1])
			m.serverLookupErr = nil
			m.selectedServer = nil
			return nil, true
		case tea.KeyEnter:
			query := strings.TrimSpace(m.serverLookupInput)
			if query == "" {
				m.serverLookupErr = fmt.Errorf("server name is required")
				m.selectedServer = nil
				return nil, true
			}

			m.serverLookupErr = nil
			m.selectedServer = nil

			switch m.current.Action {
			case "servers.show":
				return m.showServerInfoCmd(query), true
			default:
				return nil, false
			}
		default:
			return nil, false
		}
	default:
		return nil, false
	}
}

func (m Model) renderServerInfo(itemStyle lipgloss.Style) string {
	lines := []string{
		itemStyle.Render("Enter server name and press enter:"),
		itemStyle.Render("> " + m.serverLookupInput),
	}

	if m.serverLookupErr != nil {
		lines = append(lines, "", itemStyle.Render("Error: "+m.serverLookupErr.Error()))
	}

	if m.selectedServer != nil {
		lines = append(lines,
			"",
			itemStyle.Render("Name: "+m.selectedServer.Name),
			itemStyle.Render("Username: "+m.selectedServer.Username),
			itemStyle.Render("IP Address: "+displayValue(m.selectedServer.IPAddress)),
			itemStyle.Render("Password: "+m.selectedServer.Password),
		)
	}

	return strings.Join(lines, "\n")
}

func displayValue(value string) string {
	if strings.TrimSpace(value) == "" {
		return "(not configured)"
	}

	return value
}
