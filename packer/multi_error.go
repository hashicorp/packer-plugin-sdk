// Copyright IBM Corp. 2013, 2025
// SPDX-License-Identifier: MPL-2.0

package packer

import (
	"errors"
	"fmt"
	"strings"
)

// MultiError is an error type to track multiple errors. This is used to
// accumulate errors in cases such as configuration parsing, and returning
// them as a single error.
type MultiError struct {
	Errors []error
}

func (e *MultiError) Error() string {
	points := make([]string, len(e.Errors))
	for i, err := range e.Errors {
		points[i] = fmt.Sprintf("* %s", err)
	}

	return fmt.Sprintf(
		"%d error(s) occurred:\n\n%s",
		len(e.Errors), strings.Join(points, "\n"))
}

// MultiErrorAppend is a helper function that will append more errors
// onto a MultiError in order to create a larger multi-error. If the
// original error is not a MultiError, it will be turned into one.
func MultiErrorAppend(err error, errs ...error) *MultiError {
	var me *MultiError
	if err == nil {
		me = new(MultiError)
	} else if !errors.As(err, &me) || me == nil {
		newErrs := make([]error, len(errs)+1)
		newErrs[0] = err
		copy(newErrs[1:], errs)
		return &MultiError{
			Errors: newErrs,
		}
	}

	for _, verr := range errs {
		var rhs *MultiError
		if errors.As(verr, &rhs) && rhs != nil {
			me.Errors = append(me.Errors, rhs.Errors...)
		} else {
			me.Errors = append(me.Errors, verr)
		}
	}
	return me
}
