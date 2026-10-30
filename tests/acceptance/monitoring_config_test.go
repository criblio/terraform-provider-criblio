package tests

import (
	"context"
	"os"
	"testing"

	"github.com/criblio/terraform-provider-criblio/internal/provider"
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hcldec"
	"github.com/hashicorp/hcl/v2/hclparse"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/zclconf/go-cty/cty"
	"github.com/zclconf/go-cty/cty/function"
	"github.com/zclconf/go-cty/cty/function/stdlib"
	ctyjson "github.com/zclconf/go-cty/cty/json"
)

func TestMonitoringConfigurations(t *testing.T) {
	configs := map[string]string{
		"shared":         monitorSharedConfig("fixture-shared", "P3"),
		"shared-update":  monitorSharedConfig("fixture-shared", "P2"),
		"profile":        monitorProfile("fixture-profile", "Platform"),
		"profile-update": monitorProfile("fixture-profile", "Updated"),
		"policy":         notificationPolicy("fixture-policy", 10),
		"policy-update":  notificationPolicy("fixture-policy", 20),
		"silence":        monitorSilence("fixture-silence", "Maintenance", 1893456000000),
		"silence-update": monitorSilence("fixture-silence", "Updated", 1893456000000),
	}
	for _, name := range []string{"monitor_shared_config", "monitor_profile", "monitor_silence", "notification_policy"} {
		for _, file := range []string{"../../examples/" + name + "/main.tf", "../../examples/resources/criblio_" + name + "/resource.tf"} {
			content, err := os.ReadFile(file)
			if err != nil {
				t.Fatal(err)
			}
			configs[file] = string(content)
		}
	}
	for name, config := range configs {
		t.Run(name, func(t *testing.T) { validateMonitoringHCL(t, name, config) })
	}
}

func validateMonitoringHCL(t *testing.T, name, config string) {
	t.Helper()
	constructors := map[string]func() resource.Resource{
		"criblio_monitor":               provider.NewMonitorResource,
		"criblio_monitor_shared_config": provider.NewMonitorSharedConfigResource,
		"criblio_monitor_profile":       provider.NewMonitorProfileResource,
		"criblio_monitor_silence":       provider.NewMonitorSilenceResource,
		"criblio_notification_policy":   provider.NewNotificationPolicyResource,
		"criblio_notification_target":   provider.NewNotificationTargetResource,
		"criblio_notification_template": provider.NewNotificationTemplateResource,
	}
	file, diagnostics := hclparse.NewParser().ParseHCL([]byte(config), name)
	if diagnostics.HasErrors() {
		t.Fatal(diagnostics.Error())
	}
	content, _, diagnostics := file.Body.PartialContent(&hcl.BodySchema{Blocks: []hcl.BlockHeaderSchema{{Type: "resource", LabelNames: []string{"type", "name"}}}})
	if diagnostics.HasErrors() {
		t.Fatal(diagnostics.Error())
	}
	refs := map[string]map[string]cty.Value{}
	for _, block := range content.Blocks {
		if refs[block.Labels[0]] == nil {
			refs[block.Labels[0]] = map[string]cty.Value{}
		}
		refs[block.Labels[0]][block.Labels[1]] = cty.ObjectVal(map[string]cty.Value{"id": cty.StringVal("fixture-id"), "fingerprint": cty.StringVal("fixture-hash")})
	}
	eval := &hcl.EvalContext{Variables: map[string]cty.Value{}, Functions: map[string]function.Function{"jsonencode": stdlib.JSONEncodeFunc}}
	for name, values := range refs {
		eval.Variables[name] = cty.ObjectVal(values)
	}
	ctx := context.Background()
	server := providerserver.NewProtocol6(provider.New("test")())()
	for _, block := range content.Blocks {
		constructor, ok := constructors[block.Labels[0]]
		if !ok {
			t.Fatalf("unknown resource %s", block.Labels[0])
		}
		var response resource.SchemaResponse
		constructor().Schema(ctx, resource.SchemaRequest{}, &response)
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
		attrs, diagnostics := block.Body.JustAttributes()
		if diagnostics.HasErrors() {
			t.Fatal(diagnostics.Error())
		}
		for name, attr := range attrs {
			if !typ.HasAttribute(name) {
				t.Fatalf("unsupported %s.%s", block.Labels[0], name)
			}
			value, diagnostics := attr.Expr.Value(eval)
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
		validation, err := server.ValidateResourceConfig(ctx, &tfprotov6.ValidateResourceConfigRequest{TypeName: block.Labels[0], Config: &tfprotov6.DynamicValue{JSON: encoded}})
		if err != nil {
			t.Fatal(err)
		}
		for _, diagnostic := range validation.Diagnostics {
			if diagnostic.Severity == tfprotov6.DiagnosticSeverityError {
				t.Errorf("%s: %s: %s", block.Labels[0], diagnostic.Summary, diagnostic.Detail)
			}
		}
	}
}

func TestMonitoringDataSourceConfigurations(t *testing.T) {
	constructors := map[string]func() datasource.DataSource{
		"monitor_shared_config":  provider.NewMonitorSharedConfigDataSource,
		"monitor_shared_configs": provider.NewMonitorSharedConfigsDataSource,
		"monitor_profile":        provider.NewMonitorProfileDataSource,
		"monitor_profiles":       provider.NewMonitorProfilesDataSource,
		"monitor_silence":        provider.NewMonitorSilenceDataSource,
		"monitor_silences":       provider.NewMonitorSilencesDataSource,
		"notification_policy":    provider.NewNotificationPolicyDataSource,
		"notification_policies":  provider.NewNotificationPoliciesDataSource,
		"monitor_active_alerts":  provider.NewMonitorActiveAlertsDataSource,
	}
	for name, constructor := range constructors {
		t.Run(name, func(t *testing.T) {
			ctx := context.Background()
			filename := "../../examples/data-sources/criblio_" + name + "/data-source.tf"
			config, err := os.ReadFile(filename)
			if err != nil {
				t.Fatal(err)
			}
			file, diagnostics := hclparse.NewParser().ParseHCL(config, filename)
			if diagnostics.HasErrors() {
				t.Fatal(diagnostics.Error())
			}
			content, diagnostics := file.Body.Content(&hcl.BodySchema{Blocks: []hcl.BlockHeaderSchema{{Type: "data", LabelNames: []string{"type", "name"}}}})
			if diagnostics.HasErrors() {
				t.Fatal(diagnostics.Error())
			}
			var schema datasource.SchemaResponse
			constructor().Schema(ctx, datasource.SchemaRequest{}, &schema)
			typeJSON, err := schema.Schema.Type().TerraformType(ctx).MarshalJSON()
			if err != nil {
				t.Fatal(err)
			}
			typ, err := ctyjson.UnmarshalType(typeJSON)
			if err != nil {
				t.Fatal(err)
			}
			spec := hcldec.ObjectSpec{}
			for name, attribute := range schema.Schema.Attributes {
				spec[name] = &hcldec.AttrSpec{Name: name, Type: monitorOptionalObjectType(typ.AttributeType(name)), Required: attribute.IsRequired()}
			}
			for _, block := range content.Blocks {
				value, diagnostics := hcldec.Decode(block.Body, spec, nil)
				if diagnostics.HasErrors() {
					t.Fatal(diagnostics.Error())
				}
				encoded, err := ctyjson.Marshal(value, typ)
				if err != nil {
					t.Fatal(err)
				}
				response, err := providerserver.NewProtocol6(provider.New("test")())().ValidateDataResourceConfig(ctx, &tfprotov6.ValidateDataResourceConfigRequest{TypeName: "criblio_" + name, Config: &tfprotov6.DynamicValue{JSON: encoded}})
				if err != nil {
					t.Fatal(err)
				}
				for _, diagnostic := range response.Diagnostics {
					if diagnostic.Severity == tfprotov6.DiagnosticSeverityError {
						t.Errorf("%s: %s", diagnostic.Summary, diagnostic.Detail)
					}
				}
			}
		})
	}
}
