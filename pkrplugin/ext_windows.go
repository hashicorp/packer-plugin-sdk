package pkrplugin

import sdk "github.com/hashicorp/packer/packer-plugin-sdk/plugin"

var (
	FileExtension = "_x" + sdk.APIVersion + ".exe" // OS-Specific plugin file extention
)
