package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type ignoreSet map[string]bool

func readIgnoreFile(filename string) (ignoreSet, error) {
	ignored := ignoreSet{}
	file, err := os.Open(filename)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return ignored, nil
		}
		return nil, fmt.Errorf("read ignore file: %v", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		ignored[filepath.ToSlash(filepath.Clean(line))] = true
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("scan ignore file: %v", err)
	}
	return ignored, nil
}

func (s ignoreSet) ignored(path string) bool {
	return s[filepath.ToSlash(filepath.Clean(path))]
}
