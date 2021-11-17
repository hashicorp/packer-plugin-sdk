// Package filepath helps deal with configuration paths that come from user
// config files. It supports a 'slash-only' way of interpretting paths, and
// allows transforming those paths into a system agnostic path. See
// https://github.com/hashicorp/packer/issues/6188 for concrete examples.

package filepath

import (
	"path/filepath"
	"strings"
)

// A List represents a list of filepaths.
type List []string

func ListFromString(path string) List {
	return strings.Split(path, "/")
}

// String returns the user path separated by the systems path separator; on
// windows \, on other systems: /.
func (l *List) String() string {
	return filepath.Join(*l...)
}
