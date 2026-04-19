// Package servers connect to server
package servers

import (
	"context"

	domainservers "github.com/joalvarez/ui-test/internal/servers"
)

func Connect(ctx context.Context, name string) error {
	err := domainservers.Connect(ctx, name)
	if err != nil {
		return err
	}

	return nil
}
