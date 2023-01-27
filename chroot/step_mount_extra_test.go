// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package chroot

import "testing"

func TestMountExtraCleanupFunc_ImplementsCleanupFunc(t *testing.T) {
	var raw interface{} = new(StepMountExtra)
	if _, ok := raw.(Cleanup); !ok {
		t.Fatalf("cleanup func should be a CleanupFunc")
	}
}
