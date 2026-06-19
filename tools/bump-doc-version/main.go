package main

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
)

const docsIndexPath = "docs/index.md"

var providerVersionRE = regexp.MustCompile(`version = "([0-9]+)\.([0-9]+)\.([0-9]+)"`)

func main() {
	if err := bumpProviderVersion(docsIndexPath); err != nil {
		fmt.Fprintf(os.Stderr, "bump docs index version: %v\n", err)
		os.Exit(1)
	}
}

func bumpProviderVersion(path string) error {
	content, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("read %s: %w", path, err)
	}

	matches := providerVersionRE.FindSubmatchIndex(content)
	if matches == nil {
		return fmt.Errorf("provider version not found in %s", path)
	}

	patch, err := strconv.Atoi(string(content[matches[6]:matches[7]]))
	if err != nil {
		return fmt.Errorf("parse patch version: %w", err)
	}
	next := fmt.Sprintf(`version = "%s.%s.%d"`,
		content[matches[2]:matches[3]],
		content[matches[4]:matches[5]],
		patch+1,
	)

	output := append([]byte{}, content[:matches[0]]...)
	output = append(output, next...)
	output = append(output, content[matches[1]:]...)
	if err := os.WriteFile(path, output, 0644); err != nil {
		return fmt.Errorf("write %s: %w", path, err)
	}
	return nil
}
