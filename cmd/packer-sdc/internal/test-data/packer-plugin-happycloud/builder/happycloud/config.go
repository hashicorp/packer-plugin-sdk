// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

//go:generate packer-sdc struct-markdown
//go:generate packer-sdc mapstructure-to-hcl2 -type Config,CustomerEncryptionKey

package happycloud

// Config is the configuration structure for the happycloud builder. It stores
// both the publicly settable state as well as the privately generated state of
// the config object.
type Config struct {

	// The JSON file containing your account credentials. Not required if you
	// run Packer on a HappyCloud instance with a service account. Instructions for
	// creating the file or using service accounts are above.
	AccountFile string `mapstructure:"account_file" required:"false"`
	// The project ID that will be used to launch instances and store images.
	ProjectId string `mapstructure:"project_id" required:"true"`
	// Number of guest accelerator cards to add to the launched instance.
	AcceleratorCount int64 `mapstructure:"accelerator_count" required:"false"`
	// The name of a pre-allocated static external IP address. Note, must be the
	// name and not the actual IP address.
	Address string `mapstructure:"address" required:"false"`
	// If true, the default service account will not be used if
	// service_account_email is not specified. Set this value to true and omit
	// service_account_email to provision a VM with no service account.
	DisableDefaultServiceAccount bool `mapstructure:"disable_default_service_account" required:"false"`
}

// CustomerEncryptionKey helps configure a customer encryption key
type CustomerEncryptionKey struct {
	// KeyName: The name of the encryption key that is stored in happycloud
	KeyName string `mapstructure:"key_name" json:"key_name,omitempty"`

	// RawKey: Specifies a 256-bit customer-supplied encryption key, encoded in
	// RFC 4648 base64 to either encrypt or decrypt this resource.
	RawKey string `mapstructure:"raw_key" json:"raw_key,omitempty"`
}
