package main

import (
	"os"
	"os/exec"
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

func TestUpdateChangelogAllowsEntryFromHistoricalRelease(t *testing.T) {
	path := filepath.Join(t.TempDir(), "CHANGELOG.md")
	content := `# Changelog

## [Unreleased]

## [1.0.0] - 2026-01-01

### Fixed
- Fixed imports.
`
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatalf("write changelog fixture: %v", err)
	}

	if err := updateChangelog(path, []string{"fixing imports"}); err != nil {
		t.Fatalf("updateChangelog returned error: %v", err)
	}

	got, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read updated changelog: %v", err)
	}
	if count := strings.Count(string(got), "- Fixed imports.\n"); count != 2 {
		t.Fatalf("fixed entry count = %d, want 2\n%s", count, got)
	}
}

func TestLatestProviderTagIgnoresOtherTagFamilies(t *testing.T) {
	output := "v0.0.34-goatify\ngoatify-v0.0.34\nv1.25.92\nv1.25.91\n"
	if got := latestProviderTag(output); got != "v1.25.92" {
		t.Fatalf("latestProviderTag() = %q, want %q", got, "v1.25.92")
	}
}

func TestUnreleasedCommitSubjectsWithoutTagsReturnsActionableError(t *testing.T) {
	repo := t.TempDir()
	runGit(t, repo, "init")
	runGit(t, repo, "config", "user.email", "test@example.com")
	runGit(t, repo, "config", "user.name", "Test User")
	if err := os.WriteFile(filepath.Join(repo, "file.txt"), []byte("content\n"), 0644); err != nil {
		t.Fatalf("write repository fixture: %v", err)
	}
	runGit(t, repo, "add", "file.txt")
	runGit(t, repo, "commit", "-m", "fix without a release tag")

	_, err := unreleasedCommitSubjects(repo)
	if err == nil {
		t.Fatal("unreleasedCommitSubjects returned nil error")
	}
	if !strings.Contains(err.Error(), "fetch tags and full history") {
		t.Fatalf("unreleasedCommitSubjects error = %q, want fetch guidance", err)
	}
}

func runGit(t *testing.T, repo string, args ...string) {
	t.Helper()
	cmd := exec.Command("git", args...)
	cmd.Dir = repo
	if output, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("git %s: %v\n%s", strings.Join(args, " "), err, output)
	}
}
