package packer

import (
	"os"
	"path/filepath"
)

// CachePath returns an absolute path to a cache file or directory
//
// When the directory is not absolute, CachePath will try to make a
// a cache depending on the operating system.
//
// NOTE: cache directory will change depending on operating system dependent
// ex:
//   PACKER_CACHE_DIR=""            CacheDir() => "./packer_cache/
//   PACKER_CACHE_DIR=""            CacheDir("foo") => "./packer_cache/foo
//   PACKER_CACHE_DIR="bar"         CacheDir("foo") => "./bar/foo
//   PACKER_CACHE_DIR="/home/there" CacheDir("foo", "bar") => "/home/there/foo/bar
func CachePath(paths ...string) (path string, err error) {
	defer func() {
		// create the dir based on return path if it doesn't exist
		os.MkdirAll(filepath.Dir(path), os.ModePerm)
	}()
	cacheDir := getDefaultCacheDir()
	if cd := os.Getenv("PACKER_CACHE_DIR"); cd != "" {
		cacheDir = cd
	}

	paths = append([]string{cacheDir}, paths...)
	result, err := filepath.Abs(filepath.Join(paths...))
	if err != nil {
		return "", err
	}
	return result, err
}
