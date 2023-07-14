// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package fix

// Fixer applies all defined fixes on a plugin dir; a Fixer should be idempotent.
// The caller of any fix is responsible for checking if the file context have changed.
type fixer interface {
	fix(filename string, data []byte) ([]byte, error)
}

type fix struct {
	name, description string
	scan              func(dir string) ([]string, error)
	fixer
}

var (
	// availableFixes to apply to a plugin - refer to init func
	availableFixes []fix
)

func init() {
	availableFixes = []fix{
		goctyFix,
	}
}
