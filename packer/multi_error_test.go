// Copyright IBM Corp. 2013, 2025
// SPDX-License-Identifier: MPL-2.0

package packer

import (
	"errors"
	"fmt"
	"testing"
)

func TestMultiError_Impl(t *testing.T) {
	var raw interface{} = &MultiError{}
	if _, ok := raw.(error); !ok {
		t.Fatal("MultiError must implement error")
	}
}

func TestMultiErrorError(t *testing.T) {
	expected := `2 error(s) occurred:

* foo
* bar`

	errors := []error{
		errors.New("foo"),
		errors.New("bar"),
	}

	multi := &MultiError{errors}
	if multi.Error() != expected {
		t.Fatalf("bad: %s", multi.Error())
	}
}

func TestMultiErrorAppend_MultiError(t *testing.T) {
	original := &MultiError{
		Errors: []error{errors.New("foo")},
	}

	result := MultiErrorAppend(original, errors.New("bar"))
	if len(result.Errors) != 2 {
		t.Fatalf("wrong len: %d", len(result.Errors))
	}

	original = &MultiError{}
	result = MultiErrorAppend(original, errors.New("bar"))
	if len(result.Errors) != 1 {
		t.Fatalf("wrong len: %d", len(result.Errors))
	}
}

func TestMultiErrorAppend_NonMultiError(t *testing.T) {
	original := errors.New("foo")
	result := MultiErrorAppend(original, errors.New("bar"))
	if len(result.Errors) != 2 {
		t.Fatalf("wrong len: %d", len(result.Errors))
	}
}

func TestMultiErrorAppend_TypedNil(t *testing.T) {
	var errs *MultiError
	result := MultiErrorAppend(errs, errors.New("fail"))
	if len(result.Errors) != 1 {
		t.Fatalf("wrong len: %d", len(result.Errors))
	}
	if result.Errors[0].Error() != "fail" {
		t.Fatalf("wrong error: %v", result.Errors[0])
	}

	result = MultiErrorAppend(&MultiError{}, (*MultiError)(nil))
	if len(result.Errors) != 0 {
		t.Fatalf("nil multi-error was kept: %#v", result.Errors)
	}
}

func TestMultiErrorAppend_Wrapped(t *testing.T) {
	original := &MultiError{Errors: []error{errors.New("foo")}}
	result := MultiErrorAppend(fmt.Errorf("wrap: %w", original), errors.New("bar"))
	if len(result.Errors) != 2 {
		t.Fatalf("wrong len: %d", len(result.Errors))
	}

	nested := &MultiError{Errors: []error{errors.New("baz")}}
	result = MultiErrorAppend(original, fmt.Errorf("wrap: %w", nested))
	if len(result.Errors) != 3 {
		t.Fatalf("wrong len: %d", len(result.Errors))
	}
}
