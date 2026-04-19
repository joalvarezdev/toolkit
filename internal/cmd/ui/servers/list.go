// Package servers renders server-related console views.
package servers

import (
	"context"
	"fmt"
	"io"
	"slices"
	"strings"

	domainservers "github.com/joalvarez/ui-test/internal/servers"
)

var listServers = domainservers.List

func ListView(ctx context.Context) (string, error) {
	items, err := listServers(ctx)
	if err != nil {
		return "", err
	}

	return formatList(items), nil
}

func RunList(ctx context.Context, stdout io.Writer) error {
	view, err := ListView(ctx)
	if err != nil {
		return err
	}

	_, err = fmt.Fprintln(stdout, view)
	return err
}

func formatList(items []domainservers.Server) string {
	if len(items) == 0 {
		return "No configured servers found"
	}

	lines := make([]string, 0, len(items))
	for _, item := range items {
		lines = append(lines, formatListItem(item))
	}

	if len(lines) == 1 {
		return lines[0]
	}

	leftCount := (len(lines) + 1) / 2
	left := slices.Clone(lines[:leftCount])
	right := slices.Clone(lines[leftCount:])
	leftWidth := maxLineWidth(left)

	rows := make([]string, 0, leftCount)

	for i := 0; i < leftCount; i++ {
		if i >= len(right) {
			rows = append(rows, left[i])
			continue
		}

		rows = append(rows, fmt.Sprintf("%-*s  %s", leftWidth, left[i], right[i]))
	}

	return strings.Join(rows, "\n")
}

func formatListItem(item domainservers.Server) string {
	return fmt.Sprintf("%s - %s (%s)", item.Name, item.IPAddress, item.Username)
}

func maxLineWidth(lines []string) int {
	width := 0
	for _, line := range lines {
		if len(line) > width {
			width = len(line)
		}
	}

	return width
}
