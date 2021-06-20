// +build darwin freebsd linux netbsd openbsd solaris

package packer

import (
	"os"
	"path/filepath"
)

func getDefaultCacheDir() string {
	var defaultConfigFileDir string

	if xdgConfigHome := os.Getenv("XDG_CONFIG_HOME"); xdgConfigHome != "" {
		defaultConfigFileDir = xdgConfigHome
	} else {
		defaultConfigFileDir = filepath.Join(os.Getenv("HOME"), "cache")
	}

	return filepath.Join(defaultConfigFileDir, "packer")
}
