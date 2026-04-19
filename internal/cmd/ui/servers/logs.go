// Package servers connect to server
package servers

import (
	"context"

	domainservers "github.com/joalvarez/ui-test/internal/servers"
)

func Logs(ctx context.Context, name string, base string) error {
	err := domainservers.Logs(ctx, name, base)
	if err != nil {
		return err
	}

	return nil
}
