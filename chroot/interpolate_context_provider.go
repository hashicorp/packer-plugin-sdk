// Copyright IBM Corp. 2013, 2025
// SPDX-License-Identifier: MPL-2.0

package chroot

import "github.com/hashicorp/packer-plugin-sdk/template/interpolate"

type interpolateContextProvider interface {
	GetContext() interpolate.Context
}
