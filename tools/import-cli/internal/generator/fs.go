package generator

import (
	"fmt"
	"os"
	"path/filepath"
)

// FileSystem abstracts directory and file writes for safe, testable behavior.
// Implementations can use the real OS or a mock for unit tests.
type FileSystem interface {
	// MkdirAll creates the directory and any parents. Idempotent; safe to call when dir exists.
	MkdirAll(path string, perm os.FileMode) error
	// WriteFileAtomic writes data to name in dir by writing to a temp file then renaming.
	// Callers get either the full content on disk or no file (atomic).
	WriteFileAtomic(dir, name string, data []byte, perm os.FileMode) error
}

// DefaultFS is the filesystem used by WriteModuleDirectory when no FS is provided.
// Tests can replace this with a mock; production uses osFS.
var DefaultFS FileSystem = &osFS{}

// osFS uses the real os package and writes files atomically (temp + rename).
type osFS struct{}

func (*osFS) MkdirAll(path string, perm os.FileMode) error {
	return os.MkdirAll(path, perm)
}

func (*osFS) WriteFileAtomic(dir, name string, data []byte, perm os.FileMode) error {
	tmp, err := os.CreateTemp(dir, ".terraform-write-*")
	if err != nil {
		return fmt.Errorf("create temp file: %w", err)
	}
	tmpName := tmp.Name()
	defer os.Remove(tmpName) // best-effort cleanup on failure
	if _, err := tmp.Write(data); err != nil {
		_ = tmp.Close()
		return fmt.Errorf("write: %w", err)
	}
	if err := tmp.Sync(); err != nil {
		_ = tmp.Close()
		return fmt.Errorf("sync: %w", err)
	}
	if err := tmp.Close(); err != nil {
		return fmt.Errorf("close: %w", err)
	}
	target := filepath.Join(dir, name)
	if err := os.Rename(tmpName, target); err != nil {
		return fmt.Errorf("rename: %w", err)
	}
	return nil
}
