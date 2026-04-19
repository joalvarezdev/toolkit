// Package servers show all info server
package servers

import (
	"context"
	"fmt"
	"io"

	domainservers "github.com/joalvarez/toolkit/internal/servers"
)

var showServer = domainservers.FindByName

func ShowServerInfo(ctx context.Context, name string, stdout io.Writer) error {
	server, err := showServer(ctx, name)
	if err != nil {
		return err
	}

	_, err = fmt.Fprintln(stdout, "Name: "+server.Name+"\nAddress: "+server.IPAddress+"\nPassword: "+server.Password+"\nUsername: "+server.Username)

	return err
}
