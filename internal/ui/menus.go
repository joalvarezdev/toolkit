// Package ui all menus
package ui

type MenuNode struct {
	ID       string
	Label    string
	Hint     string
	Children []*MenuNode
	Action   string
}

func (n *MenuNode) IsAction() bool {
	return n.Action != ""
}

func (n *MenuNode) HasChildren() bool {
	return len(n.Children) > 0
}

var RootMenu = &MenuNode{
	ID:    "root",
	Label: "Main Menu",
	Hint:  "↑/↓ or j/k move • enter open • q quit",
	Children: []*MenuNode{
		{
			ID:    "servers",
			Label: "Servers",
			Hint:  "Manage configured servers",
			Children: []*MenuNode{
				{
					ID:     "servers.list",
					Label:  "List Servers",
					Hint:   "Show all registered servers",
					Action: "servers.list",
				},
				{
					ID:     "servers.show",
					Label:  "Show Server Information",
					Hint:   "Show basic information for selected server",
					Action: "servers.show",
				},
				{
					ID:     "servers.connect",
					Label:  "Connect Server",
					Hint:   "Connect for selected server",
					Action: "servers.connect",
				},
				{
					ID:     "servers.restart",
					Label:  "Restart Server",
					Hint:   "Restart selected server",
					Action: "servers.restart",
				},
				{
					ID:     "servers.logs",
					Label:  "Show Server Logs",
					Hint:   "Show Server logs",
					Action: "servers.logs",
				},
			},
		},
	},
}
