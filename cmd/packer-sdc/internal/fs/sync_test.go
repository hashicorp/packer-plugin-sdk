package fs

import "testing"

func TestSyncDir(t *testing.T) {
	tests := []struct {
		src     string
		dst     string
		wantErr bool
	}{
		{"test-data/a", "test-data/b", false},
	}
	for _, tt := range tests {
		t.Run(tt.src+"_to_"+tt.dst, func(t *testing.T) {
			if err := SyncDir(tt.src, tt.dst); (err != nil) != tt.wantErr {
				t.Errorf("SyncDir() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
