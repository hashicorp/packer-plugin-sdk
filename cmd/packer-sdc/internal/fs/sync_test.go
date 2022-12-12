// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package fs

import (
	"os"
	"path/filepath"
	"testing"

	"golang.org/x/mod/sumdb/dirhash"
)

func TestSyncDir(t *testing.T) {
	tests := []struct {
		folder  string
		srcSha  string
		dstSha  string
		wantErr bool
	}{
		{
			"sync_files",
			"h1:FFP6oHed1iAdwdO7q9n1Z3KCEHZxUv/TTj00nJpqAp8=",
			"h1:mKX1iSsFSBpp3SJY8AN97L5oB557oxOMBl2iJOGIVTQ=",
			false},
	}
	for _, tt := range tests {
		t.Run(tt.folder, func(t *testing.T) {
			root := filepath.Join("test-data", tt.folder)
			src := filepath.Join(root, "src")
			toCopyDst := filepath.Join(root, "dst")
			tempDst := filepath.Join(root, "temp_dst")

			hash, err := dirhash.HashDir(src, "", dirhash.Hash1)
			if err != nil {
				t.Fatal(err)
			}
			if hash != tt.srcSha {
				t.Fatalf("Incorrect src hash, found %s", hash)
			}
			hash, err = dirhash.HashDir(toCopyDst, "", dirhash.Hash1)
			if err != nil {
				t.Fatal(err)
			}
			if hash != tt.dstSha {
				t.Fatalf("Incorrect dst hash, found %s", hash)
			}

			// copying to temp_dst
			err = SyncDir(toCopyDst, tempDst)
			defer os.RemoveAll(tempDst)
			if err != nil {
				t.Fatalf("failed syncing to tempdir: %v", err)
			}
			hash, err = dirhash.HashDir(tempDst, "", dirhash.Hash1)
			if err != nil {
				t.Fatal(err)
			}
			if hash != tt.dstSha {
				t.Fatalf("Incorrect dst hash for temp dir, found %s", hash)
			}

			if err := SyncDir(src, tempDst); (err != nil) != tt.wantErr {
				t.Errorf("SyncDir() error = %v, wantErr %v", err, tt.wantErr)
			}

			hash, err = dirhash.HashDir(tempDst, "", dirhash.Hash1)
			if err != nil {
				t.Fatal(err)
			}
			if hash != tt.srcSha {
				t.Fatalf("Incorrect dst hash for end temp dir, found %s", hash)
			}
		})
	}
}
