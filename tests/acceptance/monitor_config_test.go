package tests

import (
	"context"
	"os"
	"strings"
	"testing"

	"github.com/criblio/terraform-provider-criblio/internal/provider"
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hcldec"
	"github.com/hashicorp/hcl/v2/hclparse"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/zclconf/go-cty/cty"
	"github.com/zclconf/go-cty/cty/function"
	"github.com/zclconf/go-cty/cty/function/stdlib"
	ctyjson "github.com/zclconf/go-cty/cty/json"
)

// Validate the actual HCL used by live tests and examples without API credentials.
func TestMonitorConfigurations(t *testing.T) {
	configs := map[string]string{
		"threshold-create": monitorThreshold(80),
		"threshold-update": monitorThreshold(90),
		"forecast-create":  monitorForecast(2),
		"forecast-update":  monitorForecast(3),
		"outlier":          monitorOutlier(),
		"inherited-query": strings.Replace(monitorThreshold(80), `mode       = "promql"
      dataset_id = "metrics"
      promql     = "cpu_usage_percent"`, `config_refs = [{ id = "shared", fingerprint = "hash" }]`, 1),
		"metadata": strings.Replace(monitorThreshold(80), "aetos_metadata = {}", `aetos_metadata = { owner = "example" }`, 1),
	}
	for _, path := range []string{"../../examples/monitor/main.tf", "../../examples/resources/criblio_monitor/resource.tf"} {
		content, err := os.ReadFile(path)
		if err != nil {
			t.Fatal(err)
		}
		configs[path] = string(content)
	}
	for name, config := range configs {
		t.Run(name, func(t *testing.T) {
			validateMonitorHCL(t, name, config)
		})
	}
}

func validateMonitorHCL(t *testing.T, name, config string) {
	t.Helper()
	ctx := context.Background()
	var response resource.SchemaResponse
	provider.NewMonitorResource().Schema(ctx, resource.SchemaRequest{}, &response)
	typeJSON, err := response.Schema.Type().TerraformType(ctx).MarshalJSON()
	if err != nil {
		t.Fatal(err)
	}
	typ, err := ctyjson.UnmarshalType(typeJSON)
	if err != nil {
		t.Fatal(err)
	}
	spec := hcldec.ObjectSpec{}
	for name, attribute := range response.Schema.Attributes {
		spec[name] = &hcldec.AttrSpec{Name: name, Type: monitorOptionalObjectType(typ.AttributeType(name)), Required: attribute.IsRequired()}
	}
	file, diagnostics := hclparse.NewParser().ParseHCL([]byte(config), name)
	if diagnostics.HasErrors() {
		t.Fatal(diagnostics.Error())
	}
	content, _, diagnostics := file.Body.PartialContent(&hcl.BodySchema{Blocks: []hcl.BlockHeaderSchema{
		{Type: "resource", LabelNames: []string{"type", "name"}},
		{Type: "variable", LabelNames: []string{"name"}},
	}})
	if diagnostics.HasErrors() {
		t.Fatal(diagnostics.Error())
	}
	variables := map[string]cty.Value{}
	for _, block := range content.Blocks {
		if block.Type != "variable" {
			continue
		}
		attributes, diagnostics := block.Body.JustAttributes()
		if diagnostics.HasErrors() {
			t.Fatal(diagnostics.Error())
		}
		value, diagnostics := attributes["default"].Expr.Value(nil)
		if diagnostics.HasErrors() {
			t.Fatal(diagnostics.Error())
		}
		variables[block.Labels[0]] = value
	}
	eval := &hcl.EvalContext{Variables: map[string]cty.Value{"var": cty.ObjectVal(variables)}, Functions: map[string]function.Function{"jsonencode": stdlib.JSONEncodeFunc}}
	server := providerserver.NewProtocol6(provider.New("test")())()
	count := 0
	for _, block := range content.Blocks {
		if block.Type != "resource" || block.Labels[0] != "criblio_monitor" {
			continue
		}
		count++
		attributes, diagnostics := block.Body.JustAttributes()
		if diagnostics.HasErrors() {
			t.Fatal(diagnostics.Error())
		}
		for name, attribute := range attributes {
			if !typ.HasAttribute(name) {
				t.Fatalf("unsupported monitor attribute %s", name)
			}
			value, diagnostics := attribute.Expr.Value(eval)
			if diagnostics.HasErrors() {
				t.Fatal(diagnostics.Error())
			}
			monitorCheckObjectFields(t, value, typ.AttributeType(name), name)
		}
		value, diagnostics := hcldec.Decode(block.Body, spec, eval)
		if diagnostics.HasErrors() {
			t.Fatal(diagnostics.Error())
		}
		encoded, err := ctyjson.Marshal(value, typ)
		if err != nil {
			t.Fatal(err)
		}
		validation, err := server.ValidateResourceConfig(ctx, &tfprotov6.ValidateResourceConfigRequest{
			TypeName: "criblio_monitor", Config: &tfprotov6.DynamicValue{JSON: encoded},
		})
		if err != nil {
			t.Fatal(err)
		}
		for _, diagnostic := range validation.Diagnostics {
			if diagnostic.Severity == tfprotov6.DiagnosticSeverityError {
				t.Errorf("%s: %s: %s", block.Labels[1], diagnostic.Summary, diagnostic.Detail)
			}
		}
	}
	if count == 0 {
		t.Fatal("configuration contains no monitor resource")
	}
}

func monitorOptionalObjectType(typ cty.Type) cty.Type {
	switch {
	case typ.IsObjectType():
		attributes := typ.AttributeTypes()
		var names []string
		for name, attribute := range attributes {
			attributes[name] = monitorOptionalObjectType(attribute)
			names = append(names, name)
		}
		return cty.ObjectWithOptionalAttrs(attributes, names)
	case typ.IsMapType():
		return cty.Map(monitorOptionalObjectType(typ.ElementType()))
	case typ.IsListType():
		return cty.List(monitorOptionalObjectType(typ.ElementType()))
	case typ.IsSetType():
		return cty.Set(monitorOptionalObjectType(typ.ElementType()))
	default:
		return typ
	}
}

func monitorCheckObjectFields(t *testing.T, value cty.Value, typ cty.Type, name string) {
	t.Helper()
	if value.IsNull() || !value.IsKnown() {
		return
	}
	if typ.IsObjectType() && value.Type().IsObjectType() {
		for key, nested := range value.AsValueMap() {
			if !typ.HasAttribute(key) {
				t.Fatalf("unsupported monitor attribute %s.%s", name, key)
			}
			monitorCheckObjectFields(t, nested, typ.AttributeType(key), name+"."+key)
		}
	} else if (typ.IsMapType() || typ.IsListType() || typ.IsSetType()) && value.CanIterateElements() {
		for iterator := value.ElementIterator(); iterator.Next(); {
			_, nested := iterator.Element()
			monitorCheckObjectFields(t, nested, typ.ElementType(), name+"[]")
		}
	}
}
