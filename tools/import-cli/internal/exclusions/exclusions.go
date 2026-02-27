// Package exclusions centralizes type-level and ID-level exclusions for the import CLI.
// Single source of truth for resources never exported and IDs to skip per type.
package exclusions

// NoExportTypes are resource types never exported by the import CLI.
// Merged with --exclude so these types are skipped in discovery and export.
var NoExportTypes = []string{
	"criblio_commit",                       // No list/get API in SDK; apply-only resource.
	"criblio_deploy",                       // No list/get API in SDK; apply-only resource.
	"criblio_group_system_settings",        // Cloud rejects api host/port updates on default group; config/state drift causes apply failures.
	"criblio_key",                          // Skipped from import; keys are sensitive.
	"criblio_lakehouse_dataset_connection", // Provider has no import state operation; do not generate.
	"criblio_lookup_file",                  // Control Plane API may not return content; UI uses different endpoint (knowledge/lookups).
	"criblio_mapping_ruleset",              // List API is on root CriblIo (GetAdminProductsMappingsByProduct), not a service; no standard discovery.
	"criblio_pack_lookups",                 // Same as lookup_file; content often missing from API response.
	"criblio_workspace",                    // No list/get API in SDK; workspace is implicit from config.
}

// SkipExportIDs lists resource IDs to never export, by type.
// Use for resources that fail apply (e.g. missing required attrs, API restrictions).
var SkipExportIDs = map[string]map[string]bool{
	"criblio_notification_target": {
		"system_email":          true, // smtp_target requires host/port; system_email is built-in placeholder
		"system_notifications":  true, // oneOf type unsupported by provider
	},
	"criblio_pack_destination": {
		"default": true, // read-only in Pack context
		"devnull": true, // read-only in Pack context
	},
	"criblio_pack_source": {
		"test_pack_source": true, // provider marshal fails: union type Input all fields null
	},
	"criblio_source": {
		"in_syslog":         true, // provider marshal fails: union type Input all fields null
		"in_syslog_default": true,
		"in_syslog_tls":     true,
	},
}
