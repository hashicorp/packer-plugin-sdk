// +build darwin freebsd linux netbsd openbsd solaris

package pathing

import (
	"log"
	"os"
	"path/filepath"
)

const (
	defaultConfigFile = ".packerconfig"
	defaultConfigDir  = ".packer.d"
)

func configDir() (path string, err error) {
	var dir string
	homedir := os.Getenv("HOME")
	if homedir == "" {
		return "", err
	}
	if cd := os.Getenv("PACKER_CONFIG_DIR"); cd != "" {
		log.Printf("Detected config directory from env var: %s", cd)
		dir = filepath.Join(cd, defaultConfigDir)
	} else if hasDefaultConfigFileLocation(homedir) {
		dir = filepath.Join(homedir, defaultConfigDir)
		log.Printf("Old default config directory found: %s", dir)
	} else if xdgConfigHome := os.Getenv("XDG_CONFIG_HOME"); xdgConfigHome != "" {
		log.Printf("Detected xdg config directory from env var: %s", xdgConfigHome)
		dir = xdgConfigHome
	} else {
		dir = filepath.Join(homedir, ".config", "packer")
	}

	return dir, nil
}

func hasDefaultConfigFileLocation(homedir string) bool {
	if _, err := os.Stat(filepath.Join(homedir, defaultConfigDir)); err != nil {
		return false
	}
	return true
}
