package utils

import (
	"github.com/funvibe/funxy/internal/config"
	"path/filepath"
	"strings"
)

// ResolveImportPath resolves an import path relative to a base directory if it starts with a dot.
// Otherwise returns the import path as is.
func ResolveImportPath(baseDir, importPath string) string {
	if len(importPath) > 0 && importPath[0] == '.' {
		if baseDir != "." && baseDir != "" {
			return filepath.Join(baseDir, importPath)
		}
	}
	return importPath
}

// ExtractModuleName derives a module name from a file path.
// It takes the base filename and removes the source extension.
func ExtractModuleName(path string) string {
	// Get the base filename
	name := filepath.Base(path)

	// Remove extension if present
	name = strings.TrimSuffix(name, config.SourceFileExt)

	return name
}

// GetModuleDir returns the directory context for a module path.
// If the path points to a file (ends with .lang), returns the file's directory.
// If the path points to a directory (no extension), returns the path itself.
func GetModuleDir(path string) string {
	if strings.HasSuffix(path, config.SourceFileExt) {
		return filepath.Dir(path)
	}
	return path
}
