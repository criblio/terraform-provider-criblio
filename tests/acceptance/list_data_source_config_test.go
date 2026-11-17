package tests

import (
	"os"
	"regexp"
	"strconv"
	"strings"
	"testing"
)

var listScopePattern = regexp.MustCompile(`(?m)^  (group_id|pack|project_id|organization_id) = "[^"]*"$`)

// listDataSourceAcceptanceConfig selects existing scopes without mutating them.
func listDataSourceAcceptanceConfig(t *testing.T, config string) string {
	t.Helper()
	return listScopePattern.ReplaceAllStringFunc(config, func(line string) string {
		name := listScopePattern.FindStringSubmatch(line)[1]
		var value, env string
		switch name {
		case "group_id":
			env = "CRIBL_TEST_GROUP_ID"
			value = "default"
			if strings.Contains(config, `data "criblio_search_`) {
				value = "default_search"
			}
		case "pack":
			env = "CRIBL_TEST_PACK_ID"
		case "project_id":
			env = "CRIBL_TEST_PROJECT_ID"
		case "organization_id":
			env = "CRIBL_ORGANIZATION_ID"
		}
		if configured := os.Getenv(env); configured != "" {
			value = configured
		}
		if value == "" {
			t.Skipf("%s must identify an existing scope for this list acceptance test", env)
		}
		return "  " + name + " = " + strconv.Quote(value)
	})
}
