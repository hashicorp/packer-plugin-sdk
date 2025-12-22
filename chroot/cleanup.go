// Copyright IBM Corp. 2013, 2025
// SPDX-License-Identifier: MPL-2.0

package chroot

import (
	"github.com/hashicorp/packer-plugin-sdk/multistep"
)

// Cleanup is an interface that some steps implement for early cleanup.
type Cleanup interface {
	CleanupFunc(multistep.StateBag) error
}
