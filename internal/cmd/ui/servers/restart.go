// Package servers connect to server
package servers

import (
	"context"

	domainservers "github.com/joalvarez/ui-test/internal/servers"
)

func Restart(ctx context.Context, name string) error {
	err := domainservers.Restart(ctx, name)
	if err != nil {
		return err
	}

	return nil
}
