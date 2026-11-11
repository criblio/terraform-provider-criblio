package export

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/criblio/terraform-provider-criblio/internal/provider"
	"github.com/criblio/terraform-provider-criblio/internal/restclient"
	"github.com/criblio/terraform-provider-criblio/tools/import-cli/internal/converter"
	"github.com/criblio/terraform-provider-criblio/tools/import-cli/internal/hcl"
	frameworkprovider "github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	frameworkresource "github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/stretchr/testify/require"
)

type savedQueryExportTestProvider struct {
	systemLimitsExportTestProvider
	baseURL string
}

func (p *savedQueryExportTestProvider) Configure(_ context.Context, _ frameworkprovider.ConfigureRequest, resp *frameworkprovider.ConfigureResponse) {
	resp.ResourceData = &provider.ProviderClients{RC: restclient.New(restclient.Config{BaseURL: p.baseURL, BearerToken: "test"})}
}

func (p *savedQueryExportTestProvider) Resources(context.Context) []func() frameworkresource.Resource {
	return []func() frameworkresource.Resource{provider.NewSearchSavedQueryResource}
}

func TestSavedQueryNotificationExportImportCleanPlan(t *testing.T) {
	entry, ok := buildTestRegistry(t).ByTypeName("criblio_search_saved_query")
	require.True(t, ok)

	for _, tc := range []struct {
		name          string
		notifications string
		disabled      *bool
	}{
		{name: "empty", notifications: `{}`},
		{name: "disabled", notifications: `{"disabled":true}`, disabled: boolPointer(true)},
		{name: "enabled", notifications: `{"disabled":false}`, disabled: boolPointer(false)},
		{name: "null disabled", notifications: `{"disabled":null}`},
		{name: "null notifications", notifications: `null`},
		{name: "omitted"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			notificationField := ""
			if tc.notifications != "" {
				notificationField = `,"notifications":` + tc.notifications
			}
			fixture := json.RawMessage(fmt.Sprintf(`{
  "id":"export-notifications", "name":"Scheduled logs", "query":"dataset=\"default_logs\" | limit 1",
  "schedule":{"cronSchedule":"0 * * * *","enabled":false,"keepLastN":2,"tz":"UTC","jitterPercent":10%s}
}`, notificationField))
			ids := map[string]string{"id": "export-notifications"}
			model, err := converter.ConvertRawItemWithIdentifiers(entry, fixture, ids)
			require.NoError(t, err)
			var out ExportResult
			require.NoError(t, appendResourceItemFromModel(&out, entry.TypeName, entry, ids, model, nil, nil, false, IncludeOverride{}))
			require.Len(t, out.Items, 1)
			item := out.Items[0]
			schedule := item.Attrs["schedule"].Map
			require.Equal(t, "0 * * * *", schedule["cron_schedule"].String)
			require.Equal(t, hcl.Value{Kind: hcl.KindBool, Bool: false}, schedule["enabled"])
			require.Equal(t, float64(2), schedule["keep_last_n"].Number)
			require.Equal(t, "UTC", schedule["tz"].String)
			require.Equal(t, float64(10), schedule["jitter_percent"].Number)
			if tc.disabled != nil {
				require.Equal(t, hcl.Value{Kind: hcl.KindBool, Bool: *tc.disabled}, schedule["notifications"].Map["disabled"])
			}
			file, err := hcl.FileWithResources([]hcl.ResourceInput{{TypeName: item.TypeName, Name: item.Name, Attrs: item.Attrs}}, hcl.DefaultResourceBlockOptions())
			require.NoError(t, err)
			config := string(file.Bytes())
			if tc.notifications == `{}` || tc.notifications == `{"disabled":null}` {
				require.Regexp(t, `notifications\s*=\s*\{\}`, config)
			}

			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
				if req.URL.Path != "/api/v1/m/default_search/search/saved/export-notifications" {
					t.Errorf("unexpected saved-query path: %s", req.URL.Path)
					w.WriteHeader(http.StatusBadRequest)
					return
				}
				w.Header().Set("Content-Type", "application/json")
				switch req.Method {
				case http.MethodGet:
					if err := json.NewEncoder(w).Encode(map[string]any{"count": 1, "items": []json.RawMessage{fixture}}); err != nil {
						t.Errorf("encode saved query: %v", err)
					}
				case http.MethodDelete:
					w.WriteHeader(http.StatusNoContent)
				default:
					t.Errorf("export/import/unchanged apply caused a write: %s", req.Method)
					w.WriteHeader(http.StatusBadRequest)
				}
			}))
			defer server.Close()

			name := item.TypeName + "." + item.Name
			check := resource.TestCheckNoResourceAttr(name, "schedule.notifications.disabled")
			if tc.disabled != nil {
				check = resource.TestCheckResourceAttr(name, "schedule.notifications.disabled", fmt.Sprint(*tc.disabled))
			}
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
					"criblio": providerserver.NewProtocol6WithError(&savedQueryExportTestProvider{baseURL: server.URL}),
				},
				Steps: []resource.TestStep{
					{Config: config, PlanOnly: true, ExpectNonEmptyPlan: true},
					{Config: config, ResourceName: name, ImportState: true, ImportStateId: item.ImportID, ImportStatePersist: true},
					{Config: config, PlanOnly: true},
					{Config: config, Check: check},
				},
			})
		})
	}
}

func boolPointer(value bool) *bool {
	return &value
}
