package config

import (
	"fmt"
	"os"
	"os/user"
	"path/filepath"
)

var userConfigsPath string

func PrepareUserConfigs() error {
	u, err := user.Current()
	if err != nil {
		return fmt.Errorf("cannot get current user: %s", err)
	}

	configPath := filepath.Join(u.HomeDir, ".cb")
	info, err := os.Stat(configPath)

	// Present, check if it's a folder
	if err == nil {
		if !info.IsDir() {
			return fmt.Errorf("user configs should be a folder: ~/.cb")
		}
		return nil
	}

	// Another kind of error
	if !os.IsNotExist(err) {
		return fmt.Errorf("cannot stat user configs: %s", err)
	}

	// Create it if not present
	if err := os.MkdirAll(configPath, 0755); err != nil {
		return fmt.Errorf("cannot create user configs: %s", err)
	}

	return nil
}

func GetUserConfigsPath() string {
	return userConfigsPath
}
