// Package servers datasource
package servers

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"strings"
)

type dataSource struct {
	Path string
}

type ServerJSON struct {
	Name      string `json:"name"`
	Username  string `json:"username"`
	IPAddress string `json:"ip_address"`
	Password  string `json:"password"`
	Vpn       string `json:"vpn"`
	Logs      string `json:"logs"`
}

func newDataSource() dataSource {
	return dataSource{
		Path: "~/.config/ui-test/servers/servers.json",
	}
}

func (s dataSource) Servers() ([]ServerJSON, error) {
	path, err := expandHomePath(s.Path)
	if err != nil {
		return nil, err
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var raw []ServerJSON
	if err := json.Unmarshal(data, &raw); err != nil {
		return nil, err
	}

	return raw, nil
}

func expandHomePath(path string) (string, error) {
	if path == "" {
		return "", errors.New("config source path is empty")
	}
	if path == "~" {
		return os.UserHomeDir()
	}
	if strings.HasPrefix(path, "~/") {
		home, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		return filepath.Join(home, strings.TrimPrefix(path, "~/")), nil
	}
	return path, nil
}
