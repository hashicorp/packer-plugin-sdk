// Copyright IBM Corp. 2013, 2025
// SPDX-License-Identifier: MPL-2.0

package config

import (
	"testing"
)

func TestTrilianParsing(t *testing.T) {
	type testCase struct {
		Input       string
		Output      Trilean
		ErrExpected bool
	}
	testCases := []testCase{
		{"true", TriTrue, false}, {"True", TriTrue, false},
		{"false", TriFalse, false}, {"False", TriFalse, false},
		{"", TriUnset, false}, {"badvalue", TriUnset, true},
		{"FAlse", TriUnset, true}, {"TrUe", TriUnset, true},
	}
	for _, tc := range testCases {
		tril, err := TrileanFromString(tc.Input)
		if err != nil {
			if tc.ErrExpected == false {
				t.Fatalf("Didn't expect error: %v", tc)
			}
		}
		if tc.Output != tril {
			t.Fatalf("Didn't return proper trilean. %v", tc)
		}
	}
}
