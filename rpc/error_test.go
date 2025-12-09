// Copyright IBM Corp. 2013, 2025
// SPDX-License-Identifier: MPL-2.0

package rpc

import (
	"errors"
	"testing"
)

func TestBasicError_ImplementsError(t *testing.T) {
	var _ error = new(BasicError)
}

func TestBasicError_MatchesMessage(t *testing.T) {
	err := errors.New("foo")
	wrapped := NewBasicError(err)

	if wrapped.Error() != err.Error() {
		t.Fatalf("bad: %#v", wrapped.Error())
	}
}
