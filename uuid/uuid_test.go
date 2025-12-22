// Copyright IBM Corp. 2013, 2025
// SPDX-License-Identifier: MPL-2.0

package uuid

import (
	"testing"
)

func TestTimeOrderedUuid(t *testing.T) {
	uuid := TimeOrderedUUID()
	if len(uuid) != 36 {
		t.Fatalf("bad: %s", uuid)
	}
}
