package main

import (
	"os"
	"path/filepath"
)

// FileExists reports whether the named file or directory exists.
func FileExists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	return !os.IsNotExist(err)
}

// WriteStringToFile writes content to a file named filename in dir.
// It creates the directory (and parents) if needed and sets file mode to 0644.
func WriteStringToFile(dir, filename, content string) error {
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}
	fullPath := filepath.Join(dir, filename)
	return os.WriteFile(fullPath, []byte(content), 0o644)
}
