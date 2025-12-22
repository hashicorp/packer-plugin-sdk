// Copyright IBM Corp. 2013, 2025
// SPDX-License-Identifier: MPL-2.0

/*
Package rpc contains the implementation of the remote procedure call code that
the Packer core uses to communicate with packer plugins. As a plugin maintainer,
you are unlikely to need to directly import or use this package, but it
underpins the packer server that all plugins must implement.
*/
package rpc

import (
	"encoding/gob"

	"github.com/zclconf/go-cty/cty"
)

// Test that cty types implement the gob.GobEncoder interface.
// Support for encoding/gob was removed in github.com/zclconf/go-cty@v1.11.0.
// Refer to issue https://github.com/hashicorp/packer-plugin-sdk/issues/187
var _ gob.GobEncoder = cty.Value{}

func init() {
	gob.Register(new(map[string]string))
	gob.Register(make([]interface{}, 0))
	gob.Register(new(BasicError))
}
