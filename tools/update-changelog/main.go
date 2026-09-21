package main

import (
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"
)

const changelogPath = "CHANGELOG.md"

var subjectPrefixRE = regexp.MustCompile(`(?i)^(add|added|adding|feat|fix|fixed|fixing|remove|removed|removing|security)(?:\([^)]*\))?:?\s+`)
var providerTagRE = regexp.MustCompile(`^v[0-9]+\.[0-9]+\.[0-9]+$`)

func main() {
	subjects, err := unreleasedCommitSubjects(".")
	if err != nil {
		fmt.Fprintf(os.Stderr, "find unreleased commits: %v\n", err)
		os.Exit(1)
	}
	if err := updateChangelog(changelogPath, subjects); err != nil {
		fmt.Fprintf(os.Stderr, "update changelog: %v\n", err)
		os.Exit(1)
	}
}

func unreleasedCommitSubjects(repo string) ([]string, error) {
	tagOutput, err := gitOutput(repo, "tag", "--merged", "HEAD", "--sort=-version:refname")
	if err != nil {
		return nil, fmt.Errorf("list provider tags: %w", err)
	}
	tag := latestProviderTag(string(tagOutput))
	if tag == "" {
		return nil, fmt.Errorf("no reachable provider release tag matching vX.Y.Z; fetch tags and full history before generating the changelog")
	}

	logOutput, err := gitOutput(repo, "log", "--no-merges", "--reverse", "--format=%s", tag+"..HEAD")
	if err != nil {
		return nil, fmt.Errorf("read commits after %s: %w", tag, err)
	}
	lines := strings.Split(strings.TrimSpace(string(logOutput)), "\n")
	if len(lines) == 1 && lines[0] == "" {
		return nil, nil
	}
	return lines, nil
}

func gitOutput(repo string, args ...string) ([]byte, error) {
	cmd := exec.Command("git", args...)
	cmd.Dir = repo
	return cmd.Output()
}

func latestProviderTag(output string) string {
	for _, tag := range strings.Fields(output) {
		if providerTagRE.MatchString(tag) {
			return tag
		}
	}
	return ""
}

func updateChangelog(path string, subjects []string) error {
	content, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("read %s: %w", path, err)
	}

	updated := string(content)
	for _, subject := range subjects {
		category, entry := changelogEntry(subject)
		if entry == "" {
			continue
		}
		unreleased, err := unreleasedSection(updated)
		if err != nil {
			return err
		}
		if strings.Contains(unreleased, "- "+entry+"\n") {
			continue
		}
		updated, err = insertUnreleasedEntry(updated, category, entry)
		if err != nil {
			return err
		}
	}
	if updated == string(content) {
		return nil
	}
	if err := os.WriteFile(path, []byte(updated), 0644); err != nil {
		return fmt.Errorf("write %s: %w", path, err)
	}
	return nil
}

func changelogEntry(subject string) (string, string) {
	subject = strings.TrimSpace(subject)
	if subject == "" {
		return "", ""
	}

	category := "Changed"
	prefix := strings.ToLower(strings.TrimSpace(subjectPrefixRE.FindString(subject)))
	switch {
	case strings.HasPrefix(prefix, "add"), strings.HasPrefix(prefix, "feat"):
		category = "Added"
	case strings.HasPrefix(prefix, "fix"):
		category = "Fixed"
	case strings.HasPrefix(prefix, "remove"):
		category = "Removed"
	case strings.HasPrefix(prefix, "security"):
		category = "Security"
	}
	if prefix != "" {
		subject = subjectPrefixRE.ReplaceAllString(subject, "")
	}
	subject = strings.TrimSpace(strings.TrimSuffix(subject, "."))
	if subject == "" {
		return "", ""
	}
	switch category {
	case "Added", "Fixed", "Removed":
		subject = category + " " + subject
	case "Security":
		subject = "Security: " + subject
	default:
		subject = strings.ToUpper(subject[:1]) + subject[1:]
	}
	subject += "."
	return category, subject
}

func insertUnreleasedEntry(content, category, entry string) (string, error) {
	sectionStart, sectionEnd, err := unreleasedSectionBounds(content)
	if err != nil {
		return "", err
	}
	section := content[sectionStart:sectionEnd]
	heading := "### " + category

	if headingOffset := strings.Index(section, heading); headingOffset >= 0 {
		entryStart := headingOffset + len(heading)
		nextHeadingOffset := strings.Index(section[entryStart:], "\n### ")
		insertAt := sectionEnd
		if nextHeadingOffset >= 0 {
			insertAt = sectionStart + entryStart + nextHeadingOffset
		}
		return content[:insertAt] + "- " + entry + "\n" + content[insertAt:], nil
	}

	block := "\n### " + category + "\n- " + entry + "\n"
	return content[:sectionEnd] + block + content[sectionEnd:], nil
}

func unreleasedSection(content string) (string, error) {
	start, end, err := unreleasedSectionBounds(content)
	if err != nil {
		return "", err
	}
	return content[start:end], nil
}

func unreleasedSectionBounds(content string) (int, int, error) {
	const unreleasedHeader = "## [Unreleased]"
	unreleasedStart := strings.Index(content, unreleasedHeader)
	if unreleasedStart < 0 {
		return 0, 0, fmt.Errorf("%s heading not found", unreleasedHeader)
	}
	sectionStart := unreleasedStart + len(unreleasedHeader)
	nextReleaseOffset := strings.Index(content[sectionStart:], "\n## [")
	if nextReleaseOffset < 0 {
		nextReleaseOffset = len(content) - sectionStart
	}
	sectionEnd := sectionStart + nextReleaseOffset
	return sectionStart, sectionEnd, nil
}
