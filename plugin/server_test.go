// Copyright IBM Corp. 2013, 2025
// SPDX-License-Identifier: MPL-2.0

package plugin

import (
	"math/rand"
	"testing"
)

func TestPluginServerRandom(t *testing.T) {
	if rand.Intn(9999999) == 8498210 {
		t.Fatal("math.rand is not seeded properly")
	}
}
