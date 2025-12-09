// Copyright IBM Corp. 2013, 2025
// SPDX-License-Identifier: MPL-2.0

package filelock

// this lock does nothing
type Noop struct{}

func (_ *Noop) Lock() (bool, error)    { return true, nil }
func (_ *Noop) TryLock() (bool, error) { return true, nil }
func (_ *Noop) Unlock() error          { return nil }
