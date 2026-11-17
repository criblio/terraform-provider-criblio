package provider

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"

	"github.com/criblio/terraform-provider-criblio/internal/restclient"
	"github.com/criblio/terraform-provider-criblio/tools/codegen/parser"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

// These are actions or per-scope settings with no collection GET endpoint.
var listDataSourceExclusions = map[string]string{
	"criblio_commit":                       "action",
	"criblio_worker_group_upgrade":         "action",
	"criblio_deploy":                       "action",
	"criblio_group_system_settings":        "singleton settings",
	"criblio_system_limits":                "singleton settings",
	"criblio_workspace_system_settings":    "singleton settings",
	"criblio_lakehouse_dataset_connection": "connection action has no listing API",
	"criblio_search_dataset_acl":           "per-dataset ACL has no collection API",
	"criblio_search_dataset_intelligence":  "per-dataset settings have no collection API",
}

func TestListDataSourceCoverage(t *testing.T) {
	ctx := context.Background()
	definitions, err := parser.ParseFile(filepath.Join("..", "..", "merged-spec.yml"))
	if err != nil {
		t.Fatal(err)
	}
	lists := map[string]parser.ResourceDef{}
	for _, definition := range definitions {
		if definition.List.Path != "" {
			lists[definition.TypeName] = definition
		}
	}
	for managed, discovery := range map[string]string{"criblio_pack": "criblio_pack_discovery", "criblio_group": "criblio_group_discovery", "criblio_custom_banner": "criblio_custom_banner_discovery"} {
		lists[managed] = lists[discovery]
	}
	registered := map[string]bool{}
	for _, factory := range (&CriblioProvider{}).DataSources(ctx) {
		var metadata datasource.MetadataResponse
		factory().Metadata(ctx, datasource.MetadataRequest{ProviderTypeName: "criblio"}, &metadata)
		if registered[metadata.TypeName] {
			t.Errorf("duplicate data source %s", metadata.TypeName)
		}
		registered[metadata.TypeName] = true
	}
	for _, factory := range (&CriblioProvider{}).Resources(ctx) {
		var metadata resource.MetadataResponse
		factory().Metadata(ctx, resource.MetadataRequest{ProviderTypeName: "criblio"}, &metadata)
		if _, excluded := listDataSourceExclusions[metadata.TypeName]; excluded {
			continue
		}
		definition, ok := lists[metadata.TypeName]
		if !ok || definition.ListTypeName == "" {
			t.Errorf("%s lacks list data source or documented exclusion", metadata.TypeName)
			continue
		}
		if !registered[definition.ListTypeName] {
			t.Errorf("%s is not registered", definition.ListTypeName)
		}
		for _, file := range []string{"docs/data-sources/" + definition.ListFileStem + ".md", "examples/data-sources/" + definition.ListTypeName + "/data-source.tf", "tests/acceptance/" + definition.FileStem + "_data_source_test.go"} {
			if _, err := os.Stat(filepath.Join("..", "..", file)); err != nil {
				t.Errorf("%s: %v", metadata.TypeName, err)
			}
		}
	}
}

func TestAllListDataSourcesRead(t *testing.T) {
	ctx := context.Background()
	definitions, err := parser.ParseFile(filepath.Join("..", "..", "merged-spec.yml"))
	if err != nil {
		t.Fatal(err)
	}
	lists := map[string]parser.ResourceDef{}
	for _, definition := range definitions {
		if definition.List.Path != "" && definition.StructName != "GroupSystemSettings" {
			lists[definition.ListTypeName] = definition
		}
	}
	for _, factory := range (&CriblioProvider{}).DataSources(ctx) {
		d := factory()
		var metadata datasource.MetadataResponse
		d.Metadata(ctx, datasource.MetadataRequest{ProviderTypeName: "criblio"}, &metadata)
		definition, ok := lists[metadata.TypeName]
		if !ok {
			continue
		}
		for _, empty := range []bool{false, true} {
			suffix := "/populated"
			if empty {
				suffix = "/empty"
			}
			t.Run(metadata.TypeName+suffix, func(t *testing.T) {
				var schema datasource.SchemaResponse
				d.Schema(ctx, datasource.SchemaRequest{}, &schema)
				if schema.Diagnostics.HasError() {
					t.Fatal(schema.Diagnostics)
				}
				config := map[string]any{}
				expectedPath := definition.List.Path
				for _, field := range definition.List.PathParams {
					value := "scope"
					if definition.StructName == "InstanceSettings" || definition.StructName == "SystemInfo" {
						value = "default"
					}
					if field.TerraformName == "group_id" && strings.HasPrefix(definition.TypeName, "criblio_search_") {
						value = "default_search"
					}
					config[field.TerraformName] = value
					expectedPath = strings.ReplaceAll(expectedPath, "{"+field.APIName+"}", value)
				}
				if strings.ContainsAny(expectedPath, "{}") {
					t.Fatalf("list route has an unresolved placeholder: %s", expectedPath)
				}
				requests := 0
				server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					if strings.Contains(definition.List.Path, "/p/{pack}/") && strings.HasSuffix(r.URL.Path, "/packs") {
						if err := json.NewEncoder(w).Encode(map[string]any{"items": []map[string]string{{"id": "scope"}}, "offset": 0, "totalCount": 1}); err != nil {
							t.Error(err)
						}
						return
					}
					requests++
					if r.Method != http.MethodGet || !strings.HasSuffix(r.URL.Path, expectedPath) {
						t.Errorf("unexpected request: %s %s, want %s", r.Method, r.URL.Path, expectedPath)
					}
					if definition.StructName == "NotificationTemplate" && r.URL.Query().Get("engine") != "handlebars" {
						t.Errorf("engine query: %s", r.URL.RawQuery)
					}
					offset := 0
					if definition.List.Paginated {
						var err error
						offset, err = strconv.Atoi(r.URL.Query().Get("offset"))
						if err != nil {
							t.Error(err)
						}
						if r.URL.Query().Get("limit") != "100" {
							t.Errorf("limit query: %s", r.URL.RawQuery)
						}
					}
					total := 2
					if empty {
						total = 0
					}
					items := []map[string]any{}
					if offset < total {
						items = append(items, map[string]any{"id": "first", "workspaceId": "first"})
						if !definition.List.Paginated {
							items = append(items, map[string]any{"id": "second", "workspaceId": "second"})
						} else if offset == 1 {
							items[0]["id"] = "second"
							items[0]["workspaceId"] = "second"
						}
					}
					response := map[string]any{"items": items, "count": len(items)}
					if definition.List.Paginated {
						response["totalCount"] = total
						response["offset"] = offset
						response["limit"] = 1
					}
					if err := json.NewEncoder(w).Encode(response); err != nil {
						t.Error(err)
					}
				}))
				defer server.Close()
				targetURL, err := url.Parse(server.URL)
				if err != nil {
					t.Fatal(err)
				}
				client := &http.Client{Transport: roundTripFunc(func(request *http.Request) (*http.Response, error) {
					cloned := request.Clone(request.Context())
					cloned.URL.Scheme = targetURL.Scheme
					cloned.URL.Host = targetURL.Host
					return http.DefaultTransport.RoundTrip(cloned)
				})}
				var configured datasource.ConfigureResponse
				d.(datasource.DataSourceWithConfigure).Configure(ctx, datasource.ConfigureRequest{ProviderData: &ProviderClients{RC: restclient.New(restclient.Config{BaseURL: server.URL, BearerToken: "test", HTTPClient: client})}}, &configured)
				raw := collectorPlanTestObject(schema.Schema.Type().TerraformType(ctx).(tftypes.Object), config)
				response := datasource.ReadResponse{State: tfsdk.State{Schema: schema.Schema}}
				d.Read(ctx, datasource.ReadRequest{Config: tfsdk.Config{Schema: schema.Schema, Raw: raw}}, &response)
				if response.Diagnostics.HasError() {
					t.Fatal(response.Diagnostics)
				}
				var values []attr.Value
				if schema.Schema.Attributes["items"].GetType() == types.DynamicType {
					var items types.Dynamic
					if diagnostics := response.State.GetAttribute(ctx, path.Root("items"), &items); diagnostics.HasError() {
						t.Fatal(diagnostics)
					}
					values = items.UnderlyingValue().(types.Tuple).Elements()
				} else {
					var items types.List
					if diagnostics := response.State.GetAttribute(ctx, path.Root("items"), &items); diagnostics.HasError() {
						t.Fatal(diagnostics)
					}
					if items.IsNull() || items.IsUnknown() {
						t.Fatal("collection must be known and non-null")
					}
					values = items.Elements()
				}
				want := 2
				if empty {
					want = 0
				}
				if len(values) != want {
					t.Fatalf("got %d items, want %d", len(values), want)
				}
				for i, value := range values {
					attrs := value.(types.Object).Attributes()
					if id, ok := attrs["id"].(types.String); ok && !id.IsNull() {
						expected := "first"
						if i == 1 {
							expected = "second"
						}
						if id.ValueString() != expected {
							t.Errorf("item %d ID = %s, want %s", i, id.ValueString(), expected)
						}
					}
				}
				wantRequests := 1
				if !empty && definition.List.Paginated {
					wantRequests = 2
				}
				if requests != wantRequests {
					t.Fatalf("got %d requests, want %d", requests, wantRequests)
				}
			})
		}
	}
}

func TestNotificationsListUsesConcreteGroupRoute(t *testing.T) {
	ctx := context.Background()
	requests := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requests++
		if r.URL.Path != "/api/v1/m/stream-group/notifications" || r.URL.Query().Get("groupId") != "" {
			t.Errorf("incorrect notification route: %s", r.URL)
		}
		if _, err := w.Write([]byte(`{"items":[],"totalCount":0}`)); err != nil {
			t.Error(err)
		}
	}))
	defer server.Close()
	d := NewNotificationsDataSource()
	var schema datasource.SchemaResponse
	d.Schema(ctx, datasource.SchemaRequest{}, &schema)
	var configured datasource.ConfigureResponse
	d.(datasource.DataSourceWithConfigure).Configure(ctx, datasource.ConfigureRequest{ProviderData: &ProviderClients{RC: restclient.New(restclient.Config{BaseURL: server.URL, BearerToken: "test"})}}, &configured)
	raw := collectorPlanTestObject(schema.Schema.Type().TerraformType(ctx).(tftypes.Object), map[string]any{"group_id": "stream-group"})
	response := datasource.ReadResponse{State: tfsdk.State{Schema: schema.Schema}}
	d.Read(ctx, datasource.ReadRequest{Config: tfsdk.Config{Schema: schema.Schema, Raw: raw}}, &response)
	if response.Diagnostics.HasError() || requests != 1 {
		t.Fatalf("diagnostics=%v requests=%d", response.Diagnostics, requests)
	}
}
