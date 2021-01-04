// +build aix darwin dragonfly freebsd js,wasm linux netbsd openbsd solaris

package pkrplugin

import (
	sdk "github.com/hashicorp/packer/packer-plugin-sdk/plugin"
)

var (
	FileExtension = "_x" + sdk.APIVersion // OS-Specific plugin file extention
)
