package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"testing"

	"github.com/criblio/terraform-provider-criblio/internal/auth"
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

var listDataSourceSpec struct {
	sync.Once
	definitions []parser.ResourceDef
	err         error
}

// Build from checked-in inputs rather than requiring a prior make generate.
// Both merge outputs stay in a temporary directory so tests never rewrite source files.
func listDataSourceDefinitions(t *testing.T) []parser.ResourceDef {
	t.Helper()
	listDataSourceSpec.Do(func() {
		dir := t.TempDir()
		specPath := filepath.Join(dir, "merged-spec.yml")
		command := exec.CommandContext(t.Context(), "go", "run", "./tools/merge-spec",
			"--output", specPath,
			"--cloud-only-paths-output", filepath.Join(dir, "cloud_only_paths.go"))
		command.Dir = filepath.Join("..", "..")
		if output, err := command.CombinedOutput(); err != nil {
			listDataSourceSpec.err = fmt.Errorf("merge test spec: %w\n%s", err, output)
			return
		}
		listDataSourceSpec.definitions, listDataSourceSpec.err = parser.ParseFile(specPath)
	})
	if listDataSourceSpec.err != nil {
		t.Fatal(listDataSourceSpec.err)
	}
	return listDataSourceSpec.definitions
}

func TestListDataSourceCoverage(t *testing.T) {
	ctx := context.Background()
	definitions := listDataSourceDefinitions(t)
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
	for _, deployment := range []string{"cloud", "onprem"} {
		t.Run(deployment, func(t *testing.T) {
			testAllListDataSourcesRead(t, deployment == "onprem")
		})
	}
}

func testAllListDataSourcesRead(t *testing.T, onPrem bool) {
	t.Helper()
	ctx := context.Background()
	definitions := listDataSourceDefinitions(t)
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
					pathMatches := strings.HasSuffix(r.URL.Path, expectedPath)
					if onPrem {
						pathMatches = r.URL.Path == "/api/v1"+expectedPath
					}
					if r.Method != http.MethodGet || !pathMatches {
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
				clientConfig := restclient.Config{BaseURL: server.URL, BearerToken: "test", HTTPClient: client}
				if onPrem {
					clientConfig.Credentials = &auth.CriblConfig{OnpremServerURL: server.URL}
				}
				var configured datasource.ConfigureResponse
				d.(datasource.DataSourceWithConfigure).Configure(ctx, datasource.ConfigureRequest{ProviderData: &ProviderClients{RC: restclient.New(clientConfig)}}, &configured)
				raw := collectorPlanTestObject(schema.Schema.Type().TerraformType(ctx).(tftypes.Object), config)
				response := datasource.ReadResponse{State: tfsdk.State{Schema: schema.Schema}}
				d.Read(ctx, datasource.ReadRequest{Config: tfsdk.Config{Schema: schema.Schema, Raw: raw}}, &response)
				if onPrem && auth.IsRestrictedOnPremEndpoint(expectedPath) {
					if !response.Diagnostics.HasError() || !strings.Contains(response.Diagnostics.Errors()[0].Detail(), "not supported for on-prem deployments") {
						t.Fatalf("cloud-only list must report on-prem restriction: %v", response.Diagnostics)
					}
					if requests != 0 {
						t.Fatalf("cloud-only list sent %d requests on-prem", requests)
					}
					return
				}
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

func TestOnPremMappingsListReportsLicenseFailure(t *testing.T) {
	ctx := context.Background()
	requests := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requests++
		if r.Method != http.MethodGet || r.URL.Path != "/api/v1/admin/products/stream/mappings" {
			t.Errorf("incorrect mappings request: %s %s", r.Method, r.URL)
		}
		w.WriteHeader(http.StatusForbidden)
		if _, err := w.Write([]byte(`{"status":"error","message":"Forbidden"}`)); err != nil {
			t.Error(err)
		}
	}))
	defer server.Close()
	d := NewMappingsDataSource()
	var schema datasource.SchemaResponse
	d.Schema(ctx, datasource.SchemaRequest{}, &schema)
	var configured datasource.ConfigureResponse
	d.(datasource.DataSourceWithConfigure).Configure(ctx, datasource.ConfigureRequest{ProviderData: &ProviderClients{RC: restclient.New(restclient.Config{
		Credentials: &auth.CriblConfig{OnpremServerURL: server.URL}, BearerToken: "test",
	})}}, &configured)
	raw := collectorPlanTestObject(schema.Schema.Type().TerraformType(ctx).(tftypes.Object), map[string]any{"product": "stream"})
	response := datasource.ReadResponse{State: tfsdk.State{Schema: schema.Schema}}
	d.Read(ctx, datasource.ReadRequest{Config: tfsdk.Config{Schema: schema.Schema, Raw: raw}}, &response)
	if !response.Diagnostics.HasError() || !strings.Contains(response.Diagnostics.Errors()[0].Detail(), "HTTP 403") {
		t.Fatalf("license failure must remain visible: %v", response.Diagnostics)
	}
	if requests != 1 {
		t.Fatalf("got %d requests, want 1", requests)
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

func TestKeysListAvoidsBrokenPagination(t *testing.T) {
	ctx := context.Background()
	requests := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requests++
		if r.Method != http.MethodGet || r.URL.Path != "/api/v1/m/key-group/system/keys" {
			t.Errorf("incorrect keys request: %s %s", r.Method, r.URL)
		}
		if r.URL.Query().Has("offset") || r.URL.Query().Has("limit") {
			// Cribl's pagination comparator calls id.localeCompare for keyId-only items.
			http.Error(w, `{"status":"error","message":"Cannot read properties of undefined (reading 'localeCompare')"}`, http.StatusInternalServerError)
			return
		}
		if _, err := w.Write([]byte(`{"items":[{"keyId":"first","algorithm":"aes-256-gcm"},{"keyId":"second","algorithm":"aes-256-cbc"}],"count":2}`)); err != nil {
			t.Error(err)
		}
	}))
	defer server.Close()
	d := NewKeysDataSource()
	var schema datasource.SchemaResponse
	d.Schema(ctx, datasource.SchemaRequest{}, &schema)
	var configured datasource.ConfigureResponse
	d.(datasource.DataSourceWithConfigure).Configure(ctx, datasource.ConfigureRequest{ProviderData: &ProviderClients{RC: restclient.New(restclient.Config{BaseURL: server.URL, BearerToken: "test"})}}, &configured)
	raw := collectorPlanTestObject(schema.Schema.Type().TerraformType(ctx).(tftypes.Object), map[string]any{"group_id": "key-group"})
	response := datasource.ReadResponse{State: tfsdk.State{Schema: schema.Schema}}
	d.Read(ctx, datasource.ReadRequest{Config: tfsdk.Config{Schema: schema.Schema, Raw: raw}}, &response)
	if response.Diagnostics.HasError() {
		t.Fatal(response.Diagnostics)
	}
	var model KeysListDataSourceModel
	if diagnostics := response.State.Get(ctx, &model); diagnostics.HasError() {
		t.Fatal(diagnostics)
	}
	if requests != 1 || len(model.Items.Elements()) != 2 {
		t.Fatalf("requests=%d items=%v; want one request and two keys", requests, model.Items)
	}
	for i, expected := range []string{"first", "second"} {
		item := model.Items.Elements()[i].(types.Object)
		if got := item.Attributes()["key_id"].(types.String).ValueString(); got != expected {
			t.Errorf("key %d ID = %q, want %q", i, got, expected)
		}
	}
}

func TestProjectPipelinesListAvoidsBrokenRouteRewrite(t *testing.T) {
	for _, empty := range []bool{false, true} {
		t.Run(fmt.Sprintf("empty=%t", empty), func(t *testing.T) {
			ctx := context.Background()
			requests := 0
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				requests++
				if r.Method != http.MethodGet || r.URL.Path != "/api/v1/m/project-group/system/projects/project-id/pipelines" {
					t.Errorf("incorrect project pipeline request: %s %s", r.Method, r.URL)
				}
				// Cribl matches project middleware against req.url, including its query.
				if r.URL.RawQuery != "" {
					http.Error(w, "Cannot GET /api/v1/system/projects/project-id/pipelines", http.StatusNotFound)
					return
				}
				items := []map[string]any{}
				if !empty {
					for _, id := range []string{"first", "second"} {
						items = append(items, map[string]any{"id": id, "conf": map[string]any{"description": id, "functions": []any{}}})
					}
				}
				if err := json.NewEncoder(w).Encode(map[string]any{"items": items, "count": len(items)}); err != nil {
					t.Error(err)
				}
			}))
			defer server.Close()
			d := NewProjectPipelinesDataSource()
			var schema datasource.SchemaResponse
			d.Schema(ctx, datasource.SchemaRequest{}, &schema)
			var configured datasource.ConfigureResponse
			d.(datasource.DataSourceWithConfigure).Configure(ctx, datasource.ConfigureRequest{ProviderData: &ProviderClients{RC: restclient.New(restclient.Config{BaseURL: server.URL, BearerToken: "test"})}}, &configured)
			raw := collectorPlanTestObject(schema.Schema.Type().TerraformType(ctx).(tftypes.Object), map[string]any{"group_id": "project-group", "project_id": "project-id"})
			response := datasource.ReadResponse{State: tfsdk.State{Schema: schema.Schema}}
			d.Read(ctx, datasource.ReadRequest{Config: tfsdk.Config{Schema: schema.Schema, Raw: raw}}, &response)
			if response.Diagnostics.HasError() {
				t.Fatal(response.Diagnostics)
			}
			var model ProjectPipelinesListDataSourceModel
			if diagnostics := response.State.Get(ctx, &model); diagnostics.HasError() {
				t.Fatal(diagnostics)
			}
			want := 2
			if empty {
				want = 0
			}
			if requests != 1 || model.Items.IsNull() || model.Items.IsUnknown() || len(model.Items.Elements()) != want {
				t.Fatalf("requests=%d items=%v; want one request and %d pipelines", requests, model.Items, want)
			}
			if !empty {
				for i, expected := range []string{"first", "second"} {
					attributes := model.Items.Elements()[i].(types.Object).Attributes()
					if got := attributes["id"].(types.String).ValueString(); got != expected {
						t.Errorf("pipeline %d ID = %q, want %q", i, got, expected)
					}
					if got := attributes["conf"].(types.Object).Attributes()["description"].(types.String).ValueString(); got != expected {
						t.Errorf("pipeline %d description = %q, want %q", i, got, expected)
					}
				}
			}
		})
	}
}

func TestSearchSourcesListReadsFullCollectionOnce(t *testing.T) {
	for _, count := range []int{0, 105} {
		t.Run(fmt.Sprintf("count=%d", count), func(t *testing.T) {
			ctx := context.Background()
			requests := 0
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				requests++
				if r.Method != http.MethodGet || r.URL.Path != "/api/v1/m/default_search/search/local_search/sources" {
					t.Errorf("incorrect Search sources request: %s %s", r.Method, r.URL)
				}
				if r.URL.Query().Has("offset") || r.URL.Query().Has("limit") {
					t.Errorf("Search source discovery must omit pagination: %s", r.URL)
				}
				// Bound the regression itself if pagination is accidentally reenabled.
				if requests > 2 {
					http.Error(w, "Search source list did not stop after the full collection", http.StatusInternalServerError)
					return
				}
				// Older handlers ignore pagination and return the same sources with changing status.
				items := make([]map[string]any, 0, count)
				for i := 0; i < count; i++ {
					items = append(items, map[string]any{
						"id": fmt.Sprintf("source-%d", i), "type": "http_raw",
						"description": fmt.Sprintf("source %d", i),
						"status":      map[string]any{"receivedEvents": requests},
					})
				}
				if err := json.NewEncoder(w).Encode(map[string]any{"items": items, "count": len(items)}); err != nil {
					t.Error(err)
				}
			}))
			defer server.Close()
			d := NewSearchSourcesDataSource()
			var schema datasource.SchemaResponse
			d.Schema(ctx, datasource.SchemaRequest{}, &schema)
			var configured datasource.ConfigureResponse
			d.(datasource.DataSourceWithConfigure).Configure(ctx, datasource.ConfigureRequest{ProviderData: &ProviderClients{RC: restclient.New(restclient.Config{BaseURL: server.URL, BearerToken: "test"})}}, &configured)
			raw := collectorPlanTestObject(schema.Schema.Type().TerraformType(ctx).(tftypes.Object), map[string]any{})
			response := datasource.ReadResponse{State: tfsdk.State{Schema: schema.Schema}}
			d.Read(ctx, datasource.ReadRequest{Config: tfsdk.Config{Schema: schema.Schema, Raw: raw}}, &response)
			if response.Diagnostics.HasError() {
				t.Fatal(response.Diagnostics)
			}
			var model SearchSourcesListDataSourceModel
			if diagnostics := response.State.Get(ctx, &model); diagnostics.HasError() {
				t.Fatal(diagnostics)
			}
			if requests != 1 || model.Items.IsNull() || model.Items.IsUnknown() || len(model.Items.Elements()) != count {
				t.Fatalf("requests=%d items=%v; want one request and %d sources", requests, model.Items, count)
			}
			for i, value := range model.Items.Elements() {
				attributes := value.(types.Object).Attributes()
				if got, want := attributes["id"].(types.String).ValueString(), fmt.Sprintf("source-%d", i); got != want {
					t.Errorf("source %d ID = %q, want %q", i, got, want)
				}
			}
		})
	}
}
