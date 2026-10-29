package provider

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"testing"

	"github.com/criblio/terraform-provider-criblio/internal/restclient"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

func TestSearchDashboardMetricRoundTrip(t *testing.T) {
	for _, strategy := range []string{`{"type":"auto","targetPoints":240}`, `{"type":"fixed","step":60}`} {
		t.Run(strategy, func(t *testing.T) {
			payload := searchDashboardMetricFixture(t)
			dashboard := searchDashboardJSON(t, payload)
			search := dashboard["elements"].([]any)[0].(map[string]any)["search"].(map[string]any)
			search["stepStrategy"] = searchDashboardJSON(t, []byte(strategy))
			payload, err := json.Marshal(dashboard)
			if err != nil {
				t.Fatal(err)
			}
			var model SearchDashboardModel
			if err := json.Unmarshal(payload, &model); err != nil {
				t.Fatal(err)
			}
			assertSearchDashboardJSON(t, dashboard, &model)
			state := SearchDashboardModel{ID: types.StringValue("pai7wm0s3")}
			if !isSearchDashboardImportState(&state) {
				t.Fatal("ID-only state must be recognized as an import")
			}
			applySearchDashboardAPIToState(&model, &state, true, true)
			assertSearchDashboardJSON(t, dashboard, &state)
			applySearchDashboardAPIToState(&model, &state, true, false)
			assertSearchDashboardJSON(t, dashboard, &state)
		})
	}
}

func TestSearchDashboardQueryRoundTripBackwardCompatibility(t *testing.T) {
	for _, query := range []string{
		`{"type":"saved","queryId":"saved-query","query":"dataset=logs","runMode":"adhoc"}`,
		`{"type":"inline","query":"dataset=logs","earliest":-300,"latest":"now","timezone":"UTC"}`,
		`{"type":"values","values":["first","second"]}`,
		`{"type":"metric","earliest":"-1h","latest":"now","queries":[{"localId":"A","query":"sum(events_total)"}]}`,
		`{"type":"metric","source":"systemInsightsSql","expectedOutputType":"instant","earliest":0,"latest":100,"queries":[{"localId":"A","query":null}]}`,
	} {
		t.Run(query, func(t *testing.T) {
			expected := searchDashboardJSON(t, []byte(query))
			value, err := SearchDashboardAPIValueToTerraformValue(expected, types.ObjectType{
				AttrTypes: SearchDashboardElementsDashboardElementVisualizationSearchAttrTypes(),
			})
			if err != nil {
				t.Fatal(err)
			}
			actual, err := SearchDashboardTerraformValueToJSON(value)
			if err != nil {
				t.Fatal(err)
			}
			if queries, ok := expected["queries"].([]any); ok {
				for _, item := range queries {
					row := item.(map[string]any)
					if row["query"] == nil {
						delete(row, "query")
					}
				}
			}
			if !reflect.DeepEqual(expected, actual) {
				t.Fatalf("query round trip differs: got %v, want %v", actual, expected)
			}
		})
	}
}

func TestSearchDashboardMetricAPIReadUpdate(t *testing.T) {
	payload := searchDashboardMetricFixture(t)
	expected := searchDashboardJSON(t, payload)
	var updates int
	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		if request.URL.Path != "/api/v1/m/default_search/search/dashboards/pai7wm0s3" {
			t.Errorf("unexpected API path: %s", request.URL.Path)
			writer.WriteHeader(http.StatusNotFound)
			return
		}
		if request.Method == http.MethodPatch {
			var body map[string]any
			if err := json.NewDecoder(request.Body).Decode(&body); err != nil {
				t.Errorf("decode update: %v", err)
				writer.WriteHeader(http.StatusBadRequest)
				return
			}
			if !reflect.DeepEqual(expected, body) {
				t.Errorf("update lost dashboard fields: got %v, want %v", body, expected)
			}
			updates++
		} else if request.Method != http.MethodGet {
			t.Errorf("unexpected API method: %s", request.Method)
		}
		writer.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(writer).Encode(map[string]any{"items": []any{expected}}); err != nil {
			t.Errorf("encode API response: %v", err)
		}
	}))
	defer server.Close()
	api := newSearchDashboardAPI(restclient.New(restclient.Config{BaseURL: server.URL, BearerToken: "test-token"}))
	ctx := context.Background()
	res := &SearchDashboardResource{api: api}
	var schemaResponse resource.SchemaResponse
	res.Schema(ctx, resource.SchemaRequest{}, &schemaResponse)
	imported := resource.ImportStateResponse{State: tfsdk.State{
		Schema: schemaResponse.Schema,
		Raw:    tftypes.NewValue(schemaResponse.Schema.Type().TerraformType(ctx), nil),
	}}
	res.ImportState(ctx, resource.ImportStateRequest{ID: "pai7wm0s3"}, &imported)
	if imported.Diagnostics.HasError() {
		t.Fatal(imported.Diagnostics)
	}
	read := resource.ReadResponse{State: imported.State}
	res.Read(ctx, resource.ReadRequest{State: imported.State}, &read)
	if read.Diagnostics.HasError() {
		t.Fatal(read.Diagnostics)
	}
	var state SearchDashboardModel
	if diags := read.State.Get(ctx, &state); diags.HasError() {
		t.Fatal(diags)
	}
	assertSearchDashboardJSON(t, expected, &state)
	updated, err := api.Update(context.Background(), state)
	if err != nil {
		t.Fatal(err)
	}
	assertSearchDashboardJSON(t, expected, updated)
	if updates != 1 {
		t.Fatalf("update requests = %d, want 1", updates)
	}
}

func TestSearchDashboardMetricFieldsRemainOptionalComputed(t *testing.T) {
	var response resource.SchemaResponse
	(&SearchDashboardResource{}).Schema(context.Background(), resource.SchemaRequest{}, &response)
	if response.Diagnostics.HasError() {
		t.Fatal(response.Diagnostics)
	}
	elements := response.Schema.Attributes["elements"].(schema.ListNestedAttribute)
	for _, branch := range []string{"dashboard_element_visualization", "dashboard_element_input"} {
		t.Run(branch, func(t *testing.T) {
			element := elements.NestedObject.Attributes[branch].(schema.SingleNestedAttribute)
			search := element.Attributes["search"].(schema.SingleNestedAttribute)
			metric := search.Attributes["search_query_metric"].(schema.SingleNestedAttribute)
			assertSearchDashboardOptionalComputed(t, metric.Attributes)
		})
	}
}

func assertSearchDashboardOptionalComputed(t *testing.T, attributes map[string]schema.Attribute) {
	t.Helper()
	for name, attribute := range attributes {
		if attribute.IsRequired() || !attribute.IsOptional() || !attribute.IsComputed() {
			t.Errorf("%s must remain optional and computed", name)
		}
		switch nested := attribute.(type) {
		case schema.SingleNestedAttribute:
			assertSearchDashboardOptionalComputed(t, nested.Attributes)
		case schema.ListNestedAttribute:
			assertSearchDashboardOptionalComputed(t, nested.NestedObject.Attributes)
		}
	}
}

func searchDashboardMetricFixture(t *testing.T) []byte {
	t.Helper()
	payload, err := os.ReadFile("testdata/search_dashboard_metric.json")
	if err != nil {
		t.Fatal(err)
	}
	return payload
}

func searchDashboardJSON(t *testing.T, payload []byte) map[string]any {
	t.Helper()
	var value map[string]any
	if err := json.Unmarshal(payload, &value); err != nil {
		t.Fatal(err)
	}
	return value
}

func assertSearchDashboardJSON(t *testing.T, expected map[string]any, model *SearchDashboardModel) {
	t.Helper()
	payload, err := json.Marshal(model)
	if err != nil {
		t.Fatal(err)
	}
	if actual := searchDashboardJSON(t, payload); !reflect.DeepEqual(expected, actual) {
		t.Fatalf("dashboard round trip differs: got %s, want %v", payload, expected)
	}
}
