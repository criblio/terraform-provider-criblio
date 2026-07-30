package provider

import (
	"encoding/json"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestDestinationMicrosoftFabricSASLUsesClientIDAPIKey(t *testing.T) {
	values := make(map[string]attr.Value, len(OutputMicrosoftFabricSaslAttrTypes()))
	for name, typ := range OutputMicrosoftFabricSaslAttrTypes() {
		value, err := DestinationTerraformNullValue(typ)
		if err != nil {
			t.Fatalf("null value for %s: %v", name, err)
		}
		values[name] = value
	}
	values["disabled"] = types.BoolValue(false)
	values["mechanism"] = types.StringValue("oauthbearer")
	values["client_secret_auth_type"] = types.StringValue("secret")
	values["client_id"] = types.StringValue("fabric-client-id")

	payload, err := DestinationTerraformValueToJSON(types.ObjectValueMust(OutputMicrosoftFabricSaslAttrTypes(), values))
	if err != nil {
		t.Fatalf("convert SASL payload: %v", err)
	}
	sasl := payload.(map[string]any)
	if got := sasl["clientId"]; got != "fabric-client-id" {
		t.Fatalf("clientId = %#v, want fabric-client-id; payload=%#v", got, sasl)
	}
	if _, exists := sasl["client_id"]; exists {
		t.Fatalf("payload contains invalid client_id key: %#v", sasl)
	}
}

func TestDestinationMicrosoftFabricSASLReadsClientIDAPIKey(t *testing.T) {
	value, err := DestinationAPIValueToTerraformValue(
		map[string]any{
			"disabled":             false,
			"mechanism":            "oauthbearer",
			"clientSecretAuthType": "secret",
			"clientId":             "fabric-client-id",
		},
		types.ObjectType{AttrTypes: OutputMicrosoftFabricSaslAttrTypes()},
	)
	if err != nil {
		t.Fatalf("convert SASL response: %v", err)
	}
	clientID := value.(types.Object).Attributes()["client_id"].(types.String)
	if clientID.IsNull() || clientID.ValueString() != "fabric-client-id" {
		t.Fatalf("client_id = %#v, want fabric-client-id", clientID)
	}
}

func TestDestinationSplunkResponseResolvesPlannedUnknowns(t *testing.T) {
	var api DestinationModel
	if err := json.Unmarshal([]byte(`{
		"id":"splunk-test",
		"type":"splunk",
		"host":"localhost",
		"port":9997
	}`), &api); err != nil {
		t.Fatalf("unmarshal Splunk response: %v", err)
	}
	if api.OutputSplunk == nil {
		t.Fatal("Splunk response did not select OutputSplunk variant")
	}

	state := DestinationModel{OutputSplunk: &OutputSplunkModel{
		AuthToken: types.StringUnknown(),
		AuthType:  types.StringUnknown(),
		Compress:  types.StringUnknown(),
	}}
	applyDestinationAPIToState(&api, &state, true, false)

	for name, value := range map[string]types.String{
		"auth_token": state.OutputSplunk.AuthToken,
		"auth_type":  state.OutputSplunk.AuthType,
		"compress":   state.OutputSplunk.Compress,
	} {
		if value.IsUnknown() {
			t.Errorf("output_splunk.%s remained unknown after apply", name)
		}
	}
}

func TestDestinationRouterResponseResolvesUnknownRuleDescriptions(t *testing.T) {
	ruleType := types.ObjectType{AttrTypes: OutputRouterRulesAttrTypes()}
	plannedRule := types.ObjectValueMust(OutputRouterRulesAttrTypes(), map[string]attr.Value{
		"filter":      types.StringValue("true"),
		"output":      types.StringValue("default"),
		"description": types.StringUnknown(),
		"final":       types.BoolValue(false),
	})
	apiRule := types.ObjectValueMust(OutputRouterRulesAttrTypes(), map[string]attr.Value{
		"filter":      types.StringValue("true"),
		"output":      types.StringValue("default"),
		"description": types.StringNull(),
		"final":       types.BoolValue(false),
	})
	state := DestinationModel{OutputRouter: &OutputRouterModel{
		Rules: types.ListValueMust(ruleType, []attr.Value{plannedRule}),
	}}
	api := DestinationModel{OutputRouter: &OutputRouterModel{
		Rules: types.ListValueMust(ruleType, []attr.Value{apiRule}),
	}}

	applyDestinationAPIToState(&api, &state, true, false)

	description := state.OutputRouter.Rules.Elements()[0].(types.Object).Attributes()["description"].(types.String)
	if description.IsUnknown() || !description.IsNull() {
		t.Fatalf("expected resolved null description, got %#v", description)
	}
}
