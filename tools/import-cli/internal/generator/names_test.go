package generator

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStableResourceName_terraform_safe(t *testing.T) {
	tests := []struct {
		parts []string
		want  string
	}{
		{[]string{"input-hec-1"}, "source_input-hec-1"},
		{[]string{"default", "input-1"}, "source_default_input-1"},
		{[]string{"group/id", "x"}, "source_group_id_x"},
		{[]string{"a:b", "c"}, "source_a_b_c"},
		{[]string{"UPPER", "lower"}, "source_UPPER_lower"},
		{nil, "source"},
		{[]string{""}, "source"},
	}
	for _, tt := range tests {
		got := StableResourceName("criblio_source", tt.parts)
		assert.Equal(t, tt.want, got, "parts=%v", tt.parts)
		// Must be Terraform-safe: letters, digits, underscore, hyphen only
		for _, r := range got {
			ok := (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') || r == '_' || r == '-'
			assert.True(t, ok, "StableResourceName must be Terraform-safe, got %q", got)
		}
	}
}

func TestStableResourceName_deterministic(t *testing.T) {
	parts := []string{"default", "input-hec-1"}
	a := StableResourceName("criblio_source", parts)
	b := StableResourceName("criblio_source", parts)
	assert.Equal(t, a, b)
}

func TestStableResourceNameFromMap_deterministic(t *testing.T) {
	m := map[string]string{"group_id": "g1", "id": "x", "pack": "p1"}
	a := StableResourceNameFromMap("criblio_pack_source", m)
	b := StableResourceNameFromMap("criblio_pack_source", m)
	assert.Equal(t, a, b)
	// Key order in map should not matter
	m2 := map[string]string{"id": "x", "pack": "p1", "group_id": "g1"}
	c := StableResourceNameFromMap("criblio_pack_source", m2)
	assert.Equal(t, a, c)
}

func TestStableResourceName_truncates_long(t *testing.T) {
	long := strings.Repeat("a", 100)
	got := StableResourceName("criblio_source", []string{long})
	assert.LessOrEqual(t, len(got), MaxResourceNameLength)
}
