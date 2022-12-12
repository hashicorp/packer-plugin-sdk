// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package chroot

import "github.com/hashicorp/packer-plugin-sdk/template/interpolate"

type interpolateContextProvider interface {
	GetContext() interpolate.Context
}
