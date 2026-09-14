package updater

import (
	"testing"
)

func TestCompareVersions(t *testing.T) {
	tests := []struct {
		v1       string
		v2       string
		expected int
	}{
		{"v1.2.3", "v1.2.0", 1},
		{"1.0.0", "v1.0.1", -1},
		{"v2.0.0", "v1.9.9", 1},
		{"v1.2.3", "1.2.3", 0},
		{"v1.2.3-beta", "v1.2.3", 0},
		{"v1.10.0", "v1.9.0", 1},
		{"v1.2.4", "v1.2.3+4", 1},
	}

	for _, tt := range tests {
		got := CompareVersions(tt.v1, tt.v2)
		if got != tt.expected {
			t.Errorf("CompareVersions(%q, %q) = %d; want %d", tt.v1, tt.v2, got, tt.expected)
		}
	}
}

func TestRenderUpdateNotification(t *testing.T) {
	out := RenderUpdateNotification("1.0.0", "v1.2.0")
	if out == "" {
		t.Errorf("RenderUpdateNotification returned empty string")
	}
}
