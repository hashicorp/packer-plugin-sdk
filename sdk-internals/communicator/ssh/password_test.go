// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package ssh

import (
	"reflect"
	"testing"

	"golang.org/x/crypto/ssh"
)

func TestPasswordKeyboardInteractive_Impl(t *testing.T) {
	var raw interface{} = PasswordKeyboardInteractive("foo")
	if _, ok := raw.(ssh.KeyboardInteractiveChallenge); !ok {
		t.Fatal("PasswordKeyboardInteractive must implement KeyboardInteractiveChallenge")
	}
}

func TestPasswordKeyboardInteractive_Challenge(t *testing.T) {
	p := PasswordKeyboardInteractive("foo")
	result, err := p("foo", "bar", []string{"one", "two"}, nil)
	if err != nil {
		t.Fatalf("err not nil: %s", err)
	}

	if !reflect.DeepEqual(result, []string{"foo", "foo"}) {
		t.Fatalf("invalid password: %#v", result)
	}
}
