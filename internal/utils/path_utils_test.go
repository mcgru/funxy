package utils

import (
	"testing"
)

func TestExtractModuleName(t *testing.T) {
	tests := []struct {
		path     string
		expected string
	}{
		{"simple.lang", "simple"},
		{"path/to/module.lang", "module"},
		{"module", "module"},
		{"/absolute/path/to/mod.lang", "mod"},
		{".lang", ""}, // Edge case: just extension
		{"name.with.dots.lang", "name.with.dots"},
	}

	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			got := ExtractModuleName(tt.path)
			if got != tt.expected {
				t.Errorf("ExtractModuleName(%q) = %q; want %q", tt.path, got, tt.expected)
			}
		})
	}
}

func TestGetModuleDir(t *testing.T) {
	tests := []struct {
		path     string
		expected string
	}{
		{"path/to/file.lang", "path/to"},
		{"file.lang", "."},
		{"/abs/file.lang", "/abs"},
		// Add directory cases since behavior changed
		{"path/to/dir", "path/to/dir"},
		{"/abs/dir", "/abs/dir"},
	}

	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			got := GetModuleDir(tt.path)
			if got != tt.expected {
				t.Errorf("GetModuleDir(%q) = %q; want %q", tt.path, got, tt.expected)
			}
		})
	}
}
