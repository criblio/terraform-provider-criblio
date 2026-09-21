package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestChangelogEntryClassifiesCommitSubjects(t *testing.T) {
	tests := []struct {
		subject  string
		category string
		entry    string
	}{
		{subject: "fixing worker group replacement issue.", category: "Fixed", entry: "Fixed worker group replacement issue."},
		{subject: "feat(import): support new resources", category: "Added", entry: "Added support new resources."},
		{subject: "remove legacy SDK", category: "Removed", entry: "Removed legacy SDK."},
		{subject: "update generated docs", category: "Changed", entry: "Update generated docs."},
	}

	for _, test := range tests {
		t.Run(test.subject, func(t *testing.T) {
			category, entry := changelogEntry(test.subject)
			if category != test.category || entry != test.entry {
				t.Fatalf("changelogEntry(%q) = (%q, %q), want (%q, %q)", test.subject, category, entry, test.category, test.entry)
			}
		})
	}
}

func TestUpdateChangelogAddsAndDeduplicatesUnreleasedEntries(t *testing.T) {
	path := filepath.Join(t.TempDir(), "CHANGELOG.md")
	content := `# Changelog

## [Unreleased]

### Changed
- Existing change.

## [1.0.0] - 2026-01-01
`
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatalf("write changelog fixture: %v", err)
	}

	subjects := []string{"fixing imports", "add retry support", "fixing imports"}
	if err := updateChangelog(path, subjects); err != nil {
		t.Fatalf("updateChangelog returned error: %v", err)
	}
	if err := updateChangelog(path, subjects); err != nil {
		t.Fatalf("second updateChangelog returned error: %v", err)
	}

	got, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read updated changelog: %v", err)
	}
	text := string(got)
	if strings.Count(text, "- Fixed imports.\n") != 1 {
		t.Fatalf("fixed entry count = %d, want 1\n%s", strings.Count(text, "- Fixed imports.\n"), text)
	}
	if strings.Count(text, "- Added retry support.\n") != 1 {
		t.Fatalf("added entry count = %d, want 1\n%s", strings.Count(text, "- Added retry support.\n"), text)
	}
	if !strings.Contains(text, "### Fixed\n- Fixed imports.\n") {
		t.Fatalf("missing Fixed section:\n%s", text)
	}
	if !strings.Contains(text, "### Added\n- Added retry support.\n") {
		t.Fatalf("missing Added section:\n%s", text)
	}
}
