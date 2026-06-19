package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestBumpProviderVersion(t *testing.T) {
	path := filepath.Join(t.TempDir(), "index.md")
	content := []byte(`terraform {
  required_providers {
    criblio = {
      source  = "criblio/criblio"
      version = "1.23.46"
    }
  }
}
`)
	if err := os.WriteFile(path, content, 0644); err != nil {
		t.Fatalf("write fixture: %v", err)
	}

	if err := bumpProviderVersion(path); err != nil {
		t.Fatalf("bumpProviderVersion returned error: %v", err)
	}

	got, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read fixture: %v", err)
	}
	if !strings.Contains(string(got), `version = "1.23.47"`) {
		t.Fatalf("version was not bumped: %s", got)
	}
}

func TestBumpProviderVersionMissingVersion(t *testing.T) {
	path := filepath.Join(t.TempDir(), "index.md")
	if err := os.WriteFile(path, []byte("no provider version here\n"), 0644); err != nil {
		t.Fatalf("write fixture: %v", err)
	}

	if err := bumpProviderVersion(path); err == nil {
		t.Fatalf("bumpProviderVersion returned nil error")
	}
}
