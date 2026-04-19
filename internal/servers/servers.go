// Package servers implementation methods
package servers

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func List(ctx context.Context) ([]Server, error) {
	_ = ctx

	ds := newDataSource()

	dsData, err := ds.Servers()
	if err != nil {
		return nil, err
	}

	servers := make([]Server, 0, len(dsData))
	for _, item := range dsData {
		servers = append(servers, serverFromJSON(item))
	}

	return servers, nil
}

func FindByName(ctx context.Context, name string) (Server, error) {
	_ = ctx

	ds := newDataSource()
	items, err := ds.Servers()
	if err != nil {
		return Server{}, err
	}

	for _, server := range items {
		if strings.EqualFold(server.Name, strings.TrimSpace(name)) {
			return serverFromJSON(server), nil
		}
	}

	return Server{}, fmt.Errorf("server %q not found", name)
}

func serverFromJSON(item ServerJSON) Server {
	return Server{
		Name:      item.Name,
		Username:  item.Username,
		IPAddress: item.IPAddress,
		Password:  item.Password,
		Vpn:       item.Vpn,
		Logs:      item.Logs,
	}
}

// Connect SSh
func Connect(ctx context.Context, name string) error {
	_ = ctx

	server, err := validations(ctx, name)
	if err != nil {
		return err
	}

	cmd := exec.CommandContext(
		ctx,
		"sshpass",
		"-p", server.Password,
		"ssh",
		"-o", "StrictHostKeyChecking=no",
		fmt.Sprintf("%s@%s", server.Username, server.IPAddress),
	)

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

func Restart(ctx context.Context, name string) error {
	_ = ctx

	server, err := validations(ctx, name)
	if err != nil {
		return err
	}

	cmd := exec.Command(
		"sshpass",
		"-p", server.Password,
		"ssh",
		"-o", "StrictHostKeyChecking=no",
		fmt.Sprintf("%s@%s", server.Username, server.IPAddress),
		"echo "+server.Password+" | sudo -S reboot",
	)

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

func Logs(ctx context.Context, name string, base string) error {
	_ = ctx

	server, err := validations(ctx, name)
	if err != nil {
		return err
	}

	targetBase := "logzeus7"
	if base != "" {
		targetBase = base
	}

	cmd := exec.CommandContext(
		ctx,
		"sshpass",
		"-p", server.Password,
		"ssh",
		"-o", "StrictHostKeyChecking=no",
		fmt.Sprintf("%s@%s", server.Username, server.IPAddress),
		"tail", "-F", fmt.Sprintf(server.Logs, targetBase),
	)

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

func checkSSHPass() error {
	cmd := exec.Command("which", "sshpass")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("sshpass is not installed. Please install it with: sudo apt install sshpass (or your distro's equivalent)")
	}
	return nil
}

func validateVPNConnected(ctx context.Context, vpn string) error {
	vpn = strings.TrimSpace(vpn)
	if vpn == "" {
		return nil
	}

	ok, err := vpnNMCLIActive(ctx, vpn)
	if err == nil && ok {
		return nil
	}

	return fmt.Errorf("vpn %q is not connected", vpn)
}

func vpnNMCLIActive(ctx context.Context, name string) (bool, error) {
	if _, err := exec.LookPath("nmcli"); err != nil {
		return false, err
	}

	cmd := exec.CommandContext(ctx, "nmcli", "-t", "-f", "NAME", "connection", "show", "--active")

	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out

	if err := cmd.Run(); err != nil {
		return false, err
	}

	for _, line := range strings.Split(out.String(), "\n") {
		if strings.TrimSpace(line) == name {
			return true, nil
		}
	}

	return false, nil
}

func validations(ctx context.Context, name string) (Server, error) {
	server, err := FindByName(ctx, name)
	if err != nil {
		return Server{}, err
	}

	if server.Username == "" || server.IPAddress == "" || server.Password == "" {
		return Server{}, fmt.Errorf("server %q is missing username, ip_address or password", name)
	}

	if err := checkSSHPass(); err != nil {
		return Server{}, err
	}

	if err := validateVPNConnected(ctx, server.Vpn); err != nil {
		return Server{}, err
	}

	return server, nil
}
