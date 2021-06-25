// +build darwin freebsd linux netbsd openbsd solaris

package packer

import (
	"os"
	"path/filepath"
)

func getDefaultCacheDir() string {
	var defaultConfigFileDir string

	if xdgCacheHome := os.Getenv("XDG_CACHE_HOME"); xdgCacheHome != "" {
		defaultConfigFileDir = xdgCacheHome
	} else {
		defaultConfigFileDir = filepath.Join(os.Getenv("HOME"), "cache")
	}

	return filepath.Join(defaultConfigFileDir, "packer")
}
