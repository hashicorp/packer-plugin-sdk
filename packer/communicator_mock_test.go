// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package packer

import (
	"testing"
)

func TestMockCommunicator_impl(t *testing.T) {
	var raw interface{}
	raw = new(MockCommunicator)
	if _, ok := raw.(Communicator); !ok {
		t.Fatal("should be a communicator")
	}
}
