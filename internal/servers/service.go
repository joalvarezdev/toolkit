// Package servers
package servers

import (
	"context"
)

type Server struct {
	ID        string
	Name      string
	Username  string
	IPAddress string
	Password  string
	Vpn       string
	Logs      string
}

type Source interface {
	List(ctx context.Context) ([]Server, error)
	FindByName(ctx context.Context, name string) (Server, error)
	Connect(ctx context.Context, name string) error
	Restart(ctx context.Context, name string) error
	Logs(ctx context.Context, name string, base string) error
}
