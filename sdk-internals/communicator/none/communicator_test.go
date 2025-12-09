// Copyright IBM Corp. 2013, 2025
// SPDX-License-Identifier: MPL-2.0

package none

import (
	"testing"

	packersdk "github.com/hashicorp/packer-plugin-sdk/packer"
)

func TestCommIsCommunicator(t *testing.T) {
	var raw interface{} = &comm{}
	if _, ok := raw.(packersdk.Communicator); !ok {
		t.Fatalf("comm must be a communicator")
	}
}
