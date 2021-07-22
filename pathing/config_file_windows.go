// +build windows

package pathing

import (
	"log"
	"os"
	"path/filepath"
)

const (
	defaultConfigFile = "packer.config"
	defaultConfigDir  = "packer.d"
)

func configDir() (path string, err error) {
	var dir string
	if cd := os.Getenv("PACKER_CONFIG_DIR"); cd != "" {
		log.Printf("Detected config directory from env var: %s", cd)
		dir = filepath.Join(cd, defaultConfigDir)
		return dir, nil
	}
		homedir, err := homeDir()
		if err != nil {
			return "", err
		}
		dir = filepath.Join(homedir, defaultConfigDir)
	}

	return dir, nil
}
