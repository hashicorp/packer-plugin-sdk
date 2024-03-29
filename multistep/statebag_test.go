// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package multistep

import (
	"testing"
)

func TestBasicStateBag_ImplRunner(t *testing.T) {
	var raw interface{} = &BasicStateBag{}
	if _, ok := raw.(StateBag); !ok {
		t.Fatalf("must be a StateBag")
	}
}

func TestBasicStateBag(t *testing.T) {
	b := new(BasicStateBag)
	if b.Get("foo") != nil {
		t.Fatalf("bad: %#v", b.Get("foo"))
	}

	if _, ok := b.GetOk("foo"); ok {
		t.Fatal("should not have foo")
	}

	b.Put("foo", "bar")

	if b.Get("foo").(string) != "bar" {
		t.Fatalf("bad")
	}
}
