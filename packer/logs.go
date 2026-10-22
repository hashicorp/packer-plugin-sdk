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
	sorted []string // non-empty secrets, longest first
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
	_, err = w.Write([]byte(redact(string(p), secrets)))
	if err == nil {
		n = len(p)
	}
	return
}

// FilterString returns message with registered secrets replaced by "<sensitive>".
func (l *secretFilter) FilterString(message string) string {
	secrets, _ := l.snapshot()
	return redact(message, secrets)
}

// snapshot returns the current secrets and output writer.
func (l *secretFilter) snapshot() ([]string, io.Writer) {
	l.m.Lock()
	defer l.m.Unlock()
	return l.sorted, l.w
}

// redact replaces each matching secret with "<sensitive>". Secrets are
// expected to be ordered longest first; matches are found in the original input.
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

// sortedSecrets returns non-empty secrets longest first, breaking ties
// lexicographically for deterministic ordering.
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
