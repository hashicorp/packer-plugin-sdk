// Copyright IBM Corp. 2013, 2025
// SPDX-License-Identifier: MPL-2.0

//go:build !solaris

package filelock

import "github.com/gofrs/flock"

type Flock = flock.Flock

func New(path string) *Flock {
	return flock.New(path)
}
