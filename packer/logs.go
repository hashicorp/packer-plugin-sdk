// Copyright IBM Corp. 2013, 2025
// SPDX-License-Identifier: MPL-2.0

package packer

import (
	"io"
	"sort"
	"strings"
	"sync"
)

type secretFilter struct {
	s      map[string]struct{}
	sorted []string // non-empty secrets, longest first; rebuilt on Set
	m      sync.Mutex
	w      io.Writer
}

func (l *secretFilter) Set(secrets ...string) {
	l.m.Lock()
	defer l.m.Unlock()
	for _, s := range secrets {
		l.s[s] = struct{}{}
	}
	l.sorted = sortedSecrets(l.s)
}

func (l *secretFilter) SetOutput(output io.Writer) {
	l.m.Lock()
	defer l.m.Unlock()
	l.w = output
}

func (l *secretFilter) Write(p []byte) (n int, err error) {
	secrets, w := l.snapshot()
	return w.Write([]byte(redact(string(p), secrets)))
}

// FilterString will overwrite any sensitive variables in a string, returning
// the filtered string.
func (l *secretFilter) FilterString(message string) string {
	secrets, _ := l.snapshot()
	return redact(message, secrets)
}

// snapshot reads the longest-first secret list and the output writer together
// under the mutex, so a concurrent Set or SetOutput can't race with a redaction
// in progress. Set replaces the slice wholesale rather than mutating it, so the
// returned slice stays safe to range over after the mutex is released.
func (l *secretFilter) snapshot() ([]string, io.Writer) {
	l.m.Lock()
	defer l.m.Unlock()
	return l.sorted, l.w
}

// redact replaces every occurrence of a secret with "<sensitive>". It scans the
// original input once, so a "<sensitive>" marker it writes is never itself
// searched for secrets. Feeding each replacement's output into the next lets a
// secret like "sitive" match text a previous replacement introduced, turning
// "my_token" into "<sen<sensitive>>"; scanning the original avoids that.
// secrets must be ordered longest first so the longest match wins where secrets
// overlap at a position.
func redact(s string, secrets []string) string {
	if len(secrets) == 0 {
		return s
	}
	var b strings.Builder
	b.Grow(len(s))
	for i := 0; i < len(s); {
		matched := false
		for _, secret := range secrets {
			if strings.HasPrefix(s[i:], secret) {
				b.WriteString("<sensitive>")
				i += len(secret)
				matched = true
				break
			}
		}
		if !matched {
			b.WriteByte(s[i])
			i++
		}
	}
	return b.String()
}

// sortedSecrets returns the non-empty secrets ordered longest first, and
// lexicographically for equal lengths so the order is deterministic. Redacting
// the longest match first keeps a shorter secret that is a substring of a
// longer one (e.g. "ubuntu" and "ubuntu-22.04") from leaving the remainder of
// the longer secret in the output.
func sortedSecrets(set map[string]struct{}) []string {
	secrets := make([]string, 0, len(set))
	for s := range set {
		if s != "" {
			secrets = append(secrets, s)
		}
	}
	sort.Slice(secrets, func(i, j int) bool {
		if len(secrets[i]) != len(secrets[j]) {
			return len(secrets[i]) > len(secrets[j])
		}
		return secrets[i] < secrets[j]
	})
	return secrets
}

var LogSecretFilter secretFilter

func init() {
	LogSecretFilter.s = make(map[string]struct{})
}
