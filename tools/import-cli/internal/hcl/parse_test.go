package hcl

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseHCL_valid_passes(t *testing.T) {
	valid := []byte(`resource "criblio_source" "example" {
  id       = "input-1"
  group_id = "default"
}`)
	err := ParseHCL(valid, "test.tf")
	assert.NoError(t, err)
}

func TestParseHCL_invalid_fails_with_clear_error(t *testing.T) {
	tests := []struct {
		name    string
		invalid string
		wantErr string
	}{
		{
			name:    "unclosed block",
			invalid: `resource "criblio_source" "x" {`,
			wantErr: "parse",
		},
		{
			name:    "unexpected token",
			invalid: `resource "criblio_source" "x" { id = }`,
			wantErr: "parse",
		},
		{
			name:    "invalid syntax",
			invalid: `resource "a" "b" { broken =`,
			wantErr: "parse",
		},
		{
			name:    "unclosed string",
			invalid: `resource "criblio_source" "unclosed`,
			wantErr: "parse",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ParseHCL([]byte(tt.invalid), "test.tf")
			require.Error(t, err)
			assert.Contains(t, err.Error(), tt.wantErr)
			assert.Contains(t, err.Error(), "test.tf")
		})
	}
}

func TestParseHCL_empty_file_passes(t *testing.T) {
	// Empty or comment-only files are valid HCL
	err := ParseHCL([]byte(""), "empty.tf")
	assert.NoError(t, err)

	err = ParseHCL([]byte("# comment only\n"), "comments.tf")
	assert.NoError(t, err)
}

func TestParseHCL_import_block_passes(t *testing.T) {
	valid := []byte(`import {
  to = criblio_source.example
  id = "default/input-1"
}`)
	err := ParseHCL(valid, "import.tf")
	assert.NoError(t, err)
}

func TestParseHCL_invalid_returns_diagnostics(t *testing.T) {
	invalid := []byte(`resource "x" "y" { id = `)
	err := ParseHCL(invalid, "bad.tf")
	require.Error(t, err)
	// Error should include filename and diagnostic details
	assert.True(t, strings.Contains(err.Error(), "bad.tf") || strings.Contains(err.Error(), "parse"))
}
