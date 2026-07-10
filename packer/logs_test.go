// Copyright IBM Corp. 2013, 2025
// SPDX-License-Identifier: MPL-2.0

package packer

import (
	"bytes"
	"testing"
)

// When one registered secret is a substring of another, the shorter value must
// not be replaced first. Ranging over the map is randomly ordered, so without
// longest-first ordering the short secret ("ubuntu") can be redacted before the
// long one ("ubuntu-22.04"), leaving the tail "-22.04" (part of a real secret)
// in the output. The loop makes the random map order show up reliably.
func TestSecretFilterFilterStringOverlapping(t *testing.T) {
	const in = "connecting to ubuntu-22.04 now"
	const want = "connecting to <sensitive> now"

	for i := 0; i < 100; i++ {
		l := &secretFilter{s: map[string]struct{}{
			"ubuntu-22.04": {},
			"ubuntu":       {},
		}}
		if got := l.FilterString(in); got != want {
			t.Fatalf("secret partially leaked: got %q, want %q", got, want)
		}
	}
}

func TestSecretFilterWriteOverlapping(t *testing.T) {
	const in = "connecting to ubuntu-22.04 now"
	const want = "connecting to <sensitive> now"

	for i := 0; i < 100; i++ {
		var buf bytes.Buffer
		l := &secretFilter{
			s: map[string]struct{}{
				"ubuntu-22.04": {},
				"ubuntu":       {},
			},
			w: &buf,
		}
		if _, err := l.Write([]byte(in)); err != nil {
			t.Fatal(err)
		}
		if got := buf.String(); got != want {
			t.Fatalf("secret partially leaked: got %q, want %q", got, want)
		}
	}
}
