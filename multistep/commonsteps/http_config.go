// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

//go:generate packer-sdc struct-markdown

package commonsteps

import (
	"errors"
	"fmt"

	"github.com/hashicorp/packer-plugin-sdk/template/interpolate"
)

// These are the different valid network procotol values for "http_network_protocol"
const (
	NetworkProtocolTCP       string = "tcp"
	NetworkProcotolTCP4             = "tcp4"
	NetworkProtocolTCP6             = "tcp6"
	NetworkProtocolUnix             = "unix"
	NetworkProcotlUnixPacket        = "unixpacket"
)

// Packer will create an http server serving `http_directory` when it is set, a
// random free port will be selected and the architecture of the directory
// referenced will be available in your builder.
//
// Example usage from a builder:
//
// ```
// wget http://{{ .HTTPIP }}:{{ .HTTPPort }}/foo/bar/preseed.cfg
// ```
type HTTPConfig struct {
	// Path to a directory to serve using an HTTP server. The files in this
	// directory will be available over HTTP that will be requestable from the
	// virtual machine. This is useful for hosting kickstart files and so on.
	// By default this is an empty string, which means no HTTP server will be
	// started. The address and port of the HTTP server will be available as
	// variables in `boot_command`. This is covered in more detail below.
	HTTPDir string `mapstructure:"http_directory"`
	// Key/Values to serve using an HTTP server. `http_content` works like and
	// conflicts with `http_directory`. The keys represent the paths and the
	// values contents, the keys must start with a slash, ex: `/path/to/file`.
	// `http_content` is useful for hosting kickstart files and so on. By
	// default this is empty, which means no HTTP server will be started. The
	// address and port of the HTTP server will be available as variables in
	// `boot_command`. This is covered in more detail below.
	// Example:
	// ```hcl
	//   http_content = {
	//     "/a/b"     = file("http/b")
	//     "/foo/bar" = templatefile("${path.root}/preseed.cfg", { packages = ["nginx"] })
	//   }
	// ```
	HTTPContent map[string]string `mapstructure:"http_content"`
	// These are the minimum and maximum port to use for the HTTP server
	// started to serve the `http_directory`. Because Packer often runs in
	// parallel, Packer will choose a randomly available port in this range to
	// run the HTTP server. If you want to force the HTTP server to be on one
	// port, make this minimum and maximum port the same. By default the values
	// are `8000` and `9000`, respectively.
	HTTPPortMin int `mapstructure:"http_port_min"`
	HTTPPortMax int `mapstructure:"http_port_max"`
	// This is the bind address for the HTTP server. Defaults to 0.0.0.0 so that
	// it will work with any network interface.
	HTTPAddress string `mapstructure:"http_bind_address"`
	// Use to specify a specific ip/fqdn a vm should use to reach the callback http server upon completion.
	// This is required when running via workflows/pipelines which are running within a kubernetes cluster.
	HTTPCallbackAddress string `mapstructure:"http_callback_address"`
	// This is the bind interface for the HTTP server. Defaults to the first
	// interface with a non-loopback address. Either `http_bind_address` or
	// `http_interface` can be specified.
	HTTPInterface string `mapstructure:"http_interface" undocumented:"true"`
	// Defines the HTTP Network protocol. Valid options are `tcp`, `tcp4`, `tcp6`,
	// `unix`, and `unixpacket`. This value defaults to `tcp`.
	HTTPNetworkProtocol string `mapstructure:"http_network_protocol"`
}

func (c *HTTPConfig) Prepare(ctx *interpolate.Context) []error {
	// Validation
	var errs []error

	if c.HTTPPortMin == 0 {
		c.HTTPPortMin = 8000
	}

	if c.HTTPPortMax == 0 {
		c.HTTPPortMax = 9000
	}

	if c.HTTPInterface != "" && c.HTTPAddress != "" {
		errs = append(errs,
			errors.New("either http_interface or http_bind_address can be specified"))
	}

	if c.HTTPAddress == "" {
		c.HTTPAddress = "0.0.0.0"
	}

	if c.HTTPPortMin > c.HTTPPortMax {
		errs = append(errs,
			errors.New("http_port_min must be less than http_port_max"))
	}

	if len(c.HTTPContent) > 0 && len(c.HTTPDir) > 0 {
		errs = append(errs,
			errors.New("http_content cannot be used in conjunction with http_dir. Consider using the file function to load file in memory and serve them with http_content: https://www.packer.io/docs/templates/hcl_templates/functions/file/file"))
	}

	if c.HTTPNetworkProtocol == "" {
		c.HTTPNetworkProtocol = "tcp"
	}

	validProtocol := false
	validProtocols := []string{
		NetworkProtocolTCP,
		NetworkProcotolTCP4,
		NetworkProtocolTCP6,
		NetworkProtocolUnix,
		NetworkProcotlUnixPacket,
	}

	for _, protocol := range validProtocols {
		if c.HTTPNetworkProtocol == protocol {
			validProtocol = true
			break
		}
	}

	if !validProtocol {
		errs = append(errs,
			fmt.Errorf("http_network_protocol is invalid. Must be one of: %v", validProtocols))
	}

	return errs
}
