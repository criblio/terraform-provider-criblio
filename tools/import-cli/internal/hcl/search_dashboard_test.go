package hcl

import (
	"encoding/json"
	"os"
	"strings"
	"testing"

	"github.com/criblio/terraform-provider-criblio/internal/provider"
	hclv2 "github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/hashicorp/hcl/v2/hclwrite"
)

func TestSearchDashboardMetricExportPreservesNestedFields(t *testing.T) {
	payload, err := os.ReadFile("../../../../internal/provider/testdata/search_dashboard_metric.json")
	if err != nil {
		t.Fatal(err)
	}
	var model provider.SearchDashboardModel
	if err := json.Unmarshal(payload, &model); err != nil {
		t.Fatal(err)
	}
	attributes, err := ModelToValue(&model, nil)
	if err != nil {
		t.Fatal(err)
	}
	block, err := ResourceBlock("criblio_search_dashboard", "metrics", attributes, DefaultResourceBlockOptions(), nil)
	if err != nil {
		t.Fatal(err)
	}
	file := hclwrite.NewEmptyFile()
	file.Body().AppendBlock(block)
	output := string(file.Bytes())
	for _, expected := range []string{
		`source`, `"customerMetricsPromql"`, `dataset_id`, `expected_output_type`,
		`step_strategy`, `target_points`, `timezone`, `builder`, `enabled`,
		`mode`, `show_promql_preview`, `group_by`, `filters`, `input_id`,
		`variant`, `title_action`, `"metric"`,
	} {
		if !strings.Contains(output, expected) {
			t.Errorf("export dropped %s: %s", expected, output)
		}
	}
	elements := block.Body().GetAttribute("elements").Expr().BuildTokens(nil)
	if strings.Contains(string(elements.Bytes()), `search_query_inline`) || strings.Contains(string(elements.Bytes()), `search_query_saved`) {
		t.Fatalf("export emitted inactive query wrappers: %s", output)
	}
	parsed, diags := hclsyntax.ParseConfig(file.Bytes(), "dashboard.tf", hclv2.InitialPos)
	if diags.HasErrors() {
		t.Fatal(diags)
	}
	body := parsed.Body.(*hclsyntax.Body).Blocks[0].Body
	groups, diags := body.Attributes["groups"].Expr.Value(nil)
	if diags.HasErrors() || groups.IsNull() || groups.LengthInt() != 0 {
		t.Fatalf("export did not preserve empty groups: %v, %v", groups, diags)
	}
	value, diags := body.Attributes["elements"].Expr.Value(nil)
	if diags.HasErrors() {
		t.Fatal(diags)
	}
	visualization := value.AsValueSlice()[0].GetAttr("dashboard_element_visualization")
	if config := visualization.GetAttr("config"); config.IsNull() || config.LengthInt() != 0 {
		t.Fatalf("export did not preserve empty config: %v", config)
	}
	queries := visualization.GetAttr("search").GetAttr("search_query_metric").GetAttr("queries").AsValueSlice()
	for index, field := range []string{"filters", "group_by"} {
		collection := queries[index].GetAttr("builder").GetAttr(field)
		if collection.IsNull() || collection.LengthInt() != 0 {
			t.Errorf("export did not preserve empty %s: %v", field, collection)
		}
	}
}
