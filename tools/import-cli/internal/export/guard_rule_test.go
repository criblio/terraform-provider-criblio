package export

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/criblio/terraform-provider-criblio/internal/restclient"
	importclient "github.com/criblio/terraform-provider-criblio/tools/import-cli/internal/client"
	"github.com/criblio/terraform-provider-criblio/tools/import-cli/internal/discovery"
	"github.com/criblio/terraform-provider-criblio/tools/import-cli/internal/hcl"
	hclv2 "github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGuardRuleExport(t *testing.T) {
	reg := buildTestRegistry(t)
	for _, tc := range []struct {
		name, fields string
		wantRulesets []string
		skip         bool
	}{
		{name: "empty", fields: `,"lib":"custom","rulesets":[]`},
		{name: "omitted", fields: `,"lib":"custom"`},
		{name: "null", fields: `,"lib":"custom","rulesets":null`},
		{name: "attached", fields: `,"lib":"custom","rulesets":["custom-ruleset"]`, wantRulesets: []string{"custom-ruleset"}},
		{name: "built-in", fields: `,"lib":"cribl","rulesets":[]`, skip: true},
		{name: "customized-built-in", fields: `,"lib":"cribl-custom","rulesets":[]`, skip: true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
				if req.Method != http.MethodGet || req.URL.Path != "/api/v1/m/default/lib/sds-rules/employee_badge_id" {
					t.Errorf("unexpected request: %s %s", req.Method, req.URL.Path)
					http.NotFound(w, req)
					return
				}
				w.Header().Set("Content-Type", "application/json")
				_, err := w.Write([]byte(`{"items":[{"id":"employee_badge_id","regex":"[0-9]+"` + tc.fields + `}]}`))
				if err != nil {
					t.Error(err)
				}
			}))
			defer server.Close()
			client := &importclient.Client{REST: restclient.New(restclient.Config{BaseURL: server.URL, BearerToken: "test"})}
			for _, excludeDefaults := range []bool{false, true} {
				result, err := ToResourceItems(context.Background(), client, reg, []discovery.Result{{
					TypeName: "criblio_guard_rule", Count: 1, InventoryComplete: true,
					Identifiers: []map[string]string{{"group_id": "default", "id": "employee_badge_id"}},
				}}, []string{"default"}, nil, 1, excludeDefaults, IncludeOverride{}, nil)
				require.NoError(t, err)
				if tc.skip {
					require.Empty(t, result.Items)
					require.Len(t, result.ConvertSkipped, 1)
					continue
				}
				require.Empty(t, result.ConvertSkipped)
				require.Len(t, result.Items, 1)
				item := result.Items[0]
				assert.JSONEq(t, `{"group_id":"default","id":"employee_badge_id"}`, item.ImportID)
				file, err := hcl.FileWithResources([]hcl.ResourceInput{{TypeName: item.TypeName, Name: item.Name, Attrs: item.Attrs}}, hcl.DefaultResourceBlockOptions())
				require.NoError(t, err)
				parsed, diags := hclsyntax.ParseConfig(file.Bytes(), "guard.tf", hclv2.InitialPos)
				require.False(t, diags.HasErrors(), "%s", diags.Error())
				body := parsed.Body.(*hclsyntax.Body)
				rulesets, ok := body.Blocks[0].Body.Attributes["rulesets"]
				require.True(t, ok, "required rulesets missing from HCL: %s", file.Bytes())
				value, diags := rulesets.Expr.Value(nil)
				require.False(t, diags.HasErrors(), "%s", diags.Error())
				require.Equal(t, len(tc.wantRulesets), value.LengthInt())
				for i, entry := range value.AsValueSlice() {
					assert.Equal(t, tc.wantRulesets[i], entry.AsString())
				}
			}
		})
	}
}
