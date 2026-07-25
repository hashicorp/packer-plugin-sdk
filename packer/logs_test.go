// Copyright IBM Corp. 2013, 2025
// SPDX-License-Identifier: MPL-2.0

package packer

import (
	"bytes"
	"io"
	"reflect"
	"sync"
	"testing"
)

func newSecretFilter(secrets ...string) *secretFilter {
	l := &secretFilter{s: map[string]struct{}{}}
	l.Set(secrets...)
	return l
}

// Set must record secrets longest first (lexicographically for equal lengths)
// and drop empty values, so redaction is deterministic regardless of the map's
// iteration order.
func TestSecretFilterSetSortsLongestFirst(t *testing.T) {
	l := newSecretFilter("ubuntu", "", "ubuntu-22.04", "bb", "aa")
	want := []string{"ubuntu-22.04", "ubuntu", "aa", "bb"}
	if !reflect.DeepEqual(l.sorted, want) {
		t.Fatalf("sorted secrets: got %q, want %q", l.sorted, want)
	}
}

// When one secret is a substring of another, the longest match must be redacted
// so the tail of the longer secret ("-22.04") can't leak.
func TestSecretFilterFilterStringOverlapping(t *testing.T) {
	l := newSecretFilter("ubuntu-22.04", "ubuntu")
	const want = "connecting to <sensitive> now"
	if got := l.FilterString("connecting to ubuntu-22.04 now"); got != want {
		t.Fatalf("secret partially leaked: got %q, want %q", got, want)
	}
}

// Redaction runs against the original input, so a secret that happens to appear
// inside the "<sensitive>" marker ("sitive") must not match text the marker
// introduced and turn "my_token" into "<sen<sensitive>>".
func TestSecretFilterFilterStringDoesNotRefilterMarker(t *testing.T) {
	l := newSecretFilter("my_token", "sitive")
	if got := l.FilterString("my_token"); got != "<sensitive>" {
		t.Fatalf("marker was re-filtered: got %q, want %q", got, "<sensitive>")
	}
}

func TestSecretFilterWriteOverlapping(t *testing.T) {
	var buf bytes.Buffer
	l := newSecretFilter("ubuntu-22.04", "ubuntu")
	l.SetOutput(&buf)

	if _, err := l.Write([]byte("connecting to ubuntu-22.04 now")); err != nil {
		t.Fatal(err)
	}
	const want = "connecting to <sensitive> now"
	if got := buf.String(); got != want {
		t.Fatalf("secret partially leaked: got %q, want %q", got, want)
	}
}

// Registering secrets while a redaction reads them must not race. Run under
// `go test -race`: reading the map without the mutex trips the detector.
func TestSecretFilterConcurrentSetAndRead(t *testing.T) {
	l := newSecretFilter()
	l.SetOutput(io.Discard)

	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(3)
		go func() { defer wg.Done(); l.Set("secret") }()
		go func() { defer wg.Done(); l.FilterString("a secret in a message") }()
		go func() {
			defer wg.Done()
			_, _ = l.Write([]byte("a secret in a message"))
		}()
	}
	wg.Wait()
}
