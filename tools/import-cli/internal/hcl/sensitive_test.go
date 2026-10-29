package hcl

import (
	"strings"
	"testing"
)

func TestReplaceSchemaSensitiveValues(t *testing.T) {
	attrs := map[string]Value{
		"input_http_raw": {Kind: KindMap, Map: map[string]Value{
			"auth_tokens":  {Kind: KindList, List: []Value{{Kind: KindString, String: "dummy-token-one"}, {Kind: KindString, String: "dummy-token-two"}}},
			"token_secret": {Kind: KindString, String: "stored-secret-name"},
		}},
		"connections": {Kind: KindList, List: []Value{{Kind: KindMap, Map: map[string]Value{
			"opaque": {Kind: KindString, String: "dummy-opaque-credential"},
			"label":  {Kind: KindString, String: "public-label"},
		}}}},
		"accounts": {Kind: KindMap, Map: map[string]Value{
			"first": {Kind: KindMap, Map: map[string]Value{"opaque": {Kind: KindString, String: "dummy-map-credential"}}},
		}},
		"unset": {Kind: KindNull},
		"empty": {Kind: KindString},
	}
	paths := map[string]bool{"input_http_raw.auth_tokens": true, "connections.opaque": true, "accounts.*.opaque": true, "unset": true, "empty": true}
	ReplaceSensitiveValuesWithVariableRefs(attrs, "source_example", paths)
	ReplaceSecretValuesWithVariableRefs(attrs, "source_example")
	file, err := FileWithResources([]ResourceInput{{TypeName: "criblio_source", Name: "example", Attrs: attrs}}, DefaultResourceBlockOptions())
	if err != nil {
		t.Fatal(err)
	}
	output := string(file.Bytes())
	if strings.Contains(output, "dummy-") {
		t.Fatalf("credential leaked: %s", output)
	}
	if !strings.Contains(output, "stored-secret-name") || !strings.Contains(output, "public-label") {
		t.Fatalf("non-sensitive settings must be preserved: %s", output)
	}
	if len(CollectSecretVariableNames(attrs)) != 4 {
		t.Fatal("each sensitive string must have a distinct variable")
	}
	if !attrs["unset"].IsNull() || attrs["empty"].Kind != KindString {
		t.Fatal("unset and empty credentials must retain their values")
	}
	if strings.Contains(output, "jsonencode(var.") {
		t.Fatal("plain credential strings must use plain variable references")
	}
}

func TestSecretReferenceNamesAreNotCredentials(t *testing.T) {
	for _, name := range []string{"token_secret", "text_secret", "credentials_secret", "aws_secret", "api_key_secret", "service_account_credentials_secret"} {
		if pathRefersToSensitiveAttribute("config."+name) || isSensitiveKey(name) {
			t.Errorf("%s is a stored secret name", name)
		}
	}
	if !pathRefersToSensitiveAttribute("config.application_secret") || !isSensitiveKey("client_secret") {
		t.Fatal("literal secret fields must retain fallback masking")
	}
}
