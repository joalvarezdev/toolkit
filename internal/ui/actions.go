// Package ui
package ui

import (
	"context"

	tea "github.com/charmbracelet/bubbletea"
	uiservers "github.com/joalvarez/toolkit/internal/cmd/ui/servers"
	"github.com/joalvarez/toolkit/internal/servers"
)

func (m *Model) runAction(action string) tea.Cmd {
	switch action {
	case "servers.list":
		return m.listServersCmd()
	case "servers.show":
		return nil
	default:
		return nil
	}
}

func (m *Model) listServersCmd() tea.Cmd {
	return func() tea.Msg {
		output, err := uiservers.ListView(context.Background())
		return serversListedMsg{
			output: output,
			err:    err,
		}
	}
}

func (m *Model) showServerInfoCmd(query string) tea.Cmd {
	return func() tea.Msg {
		server, err := servers.FindByName(context.Background(), query)
		if err != nil {
			return serverInfoMsg{
				query:  query,
				server: nil,
				err:    err,
			}
		}

		return serverInfoMsg{
			query:  query,
			server: &server,
			err:    nil,
		}
	}
}
