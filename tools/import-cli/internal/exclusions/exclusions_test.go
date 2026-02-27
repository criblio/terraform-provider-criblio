package exclusions

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNoExportTypes_ContainsExpectedTypes(t *testing.T) {
	expected := map[string]bool{
		"criblio_commit":                       true,
		"criblio_deploy":                       true,
		"criblio_group_system_settings":        true,
		"criblio_key":                          true,
		"criblio_lakehouse_dataset_connection": true,
		"criblio_lookup_file":                  true,
		"criblio_mapping_ruleset":               true,
		"criblio_pack_lookups":                 true,
		"criblio_workspace":                    true,
	}
	for _, typ := range NoExportTypes {
		assert.True(t, expected[typ], "NoExportTypes should contain %q", typ)
	}
	assert.Len(t, NoExportTypes, len(expected), "NoExportTypes count should match expected")
}

func TestSkipExportIDs_ContainsExpectedTypesAndIDs(t *testing.T) {
	tests := []struct {
		typeName string
		id       string
		wantSkip bool
	}{
		{"criblio_notification_target", "system_email", true},
		{"criblio_notification_target", "system_notifications", true},
		{"criblio_notification_target", "other", false},
		{"criblio_pack_destination", "default", true},
		{"criblio_pack_destination", "devnull", true},
		{"criblio_pack_destination", "my_output", false},
		{"criblio_pack_source", "test_pack_source", true},
		{"criblio_source", "in_syslog", true},
		{"criblio_source", "in_syslog_default", true},
		{"criblio_source", "in_syslog_tls", true},
		{"criblio_source", "in_http", false},
		{"unknown_type", "any_id", false},
	}
	for _, tt := range tests {
		ids, ok := SkipExportIDs[tt.typeName]
		gotSkip := ok && ids[tt.id]
		assert.Equal(t, tt.wantSkip, gotSkip, "SkipExportIDs[%q][%q] = %v, want %v", tt.typeName, tt.id, gotSkip, tt.wantSkip)
	}
}
