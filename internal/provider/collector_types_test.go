package provider

import (
	"encoding/json"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestCollectorMarshalJSONIncludesRestOAuthClientSecret(t *testing.T) {
	confValues := make(map[string]attr.Value, len(InputCollectorRestCollectorConfAttrTypes()))
	for name, typ := range InputCollectorRestCollectorConfAttrTypes() {
		value, err := CollectorTerraformNullValue(typ)
		if err != nil {
			t.Fatalf("null value for %s: %v", name, err)
		}
		confValues[name] = value
	}
	confValues["authentication"] = types.StringValue("oauth")
	confValues["client_secret_param_name"] = types.StringValue("client_secret")
	confValues["client_secret_param_value"] = types.StringValue("vault-sourced-secret")
	collector := types.ObjectValueMust(InputCollectorRestCollectorAttrTypes(), map[string]attr.Value{
		"type": types.StringValue("rest"),
		"conf": types.ObjectValueMust(InputCollectorRestCollectorConfAttrTypes(), confValues),
	})
	model := CollectorModel{InputCollectorRest: &InputCollectorRestModel{
		ID:        types.StringValue("salesforce"),
		Collector: collector,
	}}

	data, err := json.Marshal(model)
	if err != nil {
		t.Fatalf("json.Marshal returned error: %v", err)
	}
	var body map[string]any
	if err := json.Unmarshal(data, &body); err != nil {
		t.Fatalf("json.Unmarshal returned error: %v", err)
	}
	collectorBody := body["collector"].(map[string]any)
	confBody := collectorBody["conf"].(map[string]any)
	if got := confBody["clientSecretParamValue"]; got != "vault-sourced-secret" {
		t.Fatalf("clientSecretParamValue = %#v, want Vault-sourced value; body=%s", got, data)
	}
}

func TestCollectorMarshalJSONIncludesSavedJobType(t *testing.T) {
	model := CollectorModel{
		GroupID: types.StringValue("default"),
		ID:      types.StringValue("rest-api-demo-collector"),
		InputCollectorRest: &InputCollectorRestModel{
			ID: types.StringValue("rest-api-demo-collector"),
		},
	}

	data, err := json.Marshal(model)
	if err != nil {
		t.Fatalf("json.Marshal returned error: %v", err)
	}

	var got map[string]any
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatalf("json.Unmarshal returned error: %v", err)
	}

	if got["type"] != "collection" {
		t.Fatalf("expected saved job type collection, got %#v in %s", got["type"], data)
	}
	if got["id"] != "rest-api-demo-collector" {
		t.Fatalf("expected collector id in body, got %#v in %s", got["id"], data)
	}
}

func TestCollectorModelUnmarshalGoogleCloudStorageAlias(t *testing.T) {
	var model CollectorModel
	err := json.Unmarshal([]byte(`{
		"id": "my-gcs-collector",
		"type": "collection",
		"collector": {
			"type": "google_cloud_storage",
			"conf": {
				"bucket": "my-bucket",
				"path": "logs/"
			}
		}
	}`), &model)
	if err != nil {
		t.Fatalf("UnmarshalJSON returned error: %v", err)
	}
	if model.InputCollectorGCS == nil {
		t.Fatalf("InputCollectorGCS was not selected")
	}
	if model.InputCollectorGCS.Collector.IsNull() || model.InputCollectorGCS.Collector.IsUnknown() {
		t.Fatalf("collector = %#v", model.InputCollectorGCS.Collector)
	}
}

func TestCollectorModelUnmarshalAzureBlobAlias(t *testing.T) {
	var model CollectorModel
	err := json.Unmarshal([]byte(`{
		"id": "my-azure-collector",
		"type": "collection",
		"collector": {
			"type": "azure_blob",
			"conf": {
				"containerName": "logs"
			}
		}
	}`), &model)
	if err != nil {
		t.Fatalf("UnmarshalJSON returned error: %v", err)
	}
	if model.InputCollectorAzureBlob == nil {
		t.Fatalf("InputCollectorAzureBlob was not selected")
	}
}

func TestCollectorModelUnmarshalHealthCheckAlias(t *testing.T) {
	var model CollectorModel
	err := json.Unmarshal([]byte(`{
		"id": "my-health-collector",
		"type": "collection",
		"collector": {
			"type": "health_check",
			"conf": {
				"collectUrl": "https://api.example.com/health"
			}
		}
	}`), &model)
	if err != nil {
		t.Fatalf("UnmarshalJSON returned error: %v", err)
	}
	if model.InputCollectorHealthCheck == nil {
		t.Fatalf("InputCollectorHealthCheck was not selected")
	}
}
