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
