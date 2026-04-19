// Package main for commands
package main

import (
	"context"
	"fmt"
	"io"
	"os"

	uiservers "github.com/joalvarez/ui-test/internal/cmd/ui/servers"
)

func main() {
	if err := run(context.Background(), os.Args[1:], os.Stdin, os.Stdout, os.Stderr); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(ctx context.Context, args []string, stdin io.Reader, stdout, stderr io.Writer) error {
	if len(args) == 0 {
		writeUsage(stderr)
		return fmt.Errorf("subcommand is required")
	}

	switch args[0] {
	case "list":
		return uiservers.RunList(ctx, stdout)
	case "show":
		if len(args) < 2 {
			return fmt.Errorf("usage: servers show <name>")
		}
		return uiservers.ShowServerInfo(ctx, args[1], stdout)
	case "connect":
		if len(args) < 2 {
			return fmt.Errorf("usage: servers connect <name>")
		}
		return uiservers.Connect(ctx, args[1])
	case "restart":
		if len(args) < 2 {
			return fmt.Errorf("usage: servers restart <name>")
		}
		return uiservers.Restart(ctx, args[1])
	case "logs":
		if len(args) < 2 {
			return fmt.Errorf("usage: servers logs <name>")
		}
		if len(args) < 3 {
			return uiservers.Logs(ctx, args[1], "")
		}
		return uiservers.Logs(ctx, args[1], args[2])
	case "help", "-h", "--help":
		writeUsage(stdout)
		return nil
	default:
		writeUsage(stderr)
		return fmt.Errorf("unknown subcommand %q", args[0])
	}
}

func writeUsage(w io.Writer) {
	fmt.Fprintln(w, "Usage: servers <subcommand>")
	fmt.Fprintln(w)
	fmt.Fprintln(w, "Available subcommands:")
	fmt.Fprintln(w, "  list    List configured servers")
	fmt.Fprintln(w, "  show    Show a configured server")
	fmt.Fprintln(w, "  connect    Connect a configured server")
	fmt.Fprintln(w, "  restart    Restart a configured server")
}
