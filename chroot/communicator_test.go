// Copyright IBM Corp. 2013, 2025
// SPDX-License-Identifier: MPL-2.0

package chroot

import (
	"testing"

	packersdk "github.com/hashicorp/packer-plugin-sdk/packer"
)

func TestCommunicator_ImplementsCommunicator(t *testing.T) {
	var raw interface{} = &Communicator{}

	if _, ok := raw.(packersdk.Communicator); !ok {
		t.Fatalf("Communicator should be a communicator")
	}
}
