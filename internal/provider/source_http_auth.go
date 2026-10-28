package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// normalizeHTTPRawExtendedAuthTokens drops empty API placeholders before source
// and pack-source models decode them. Keep legacy authTokens and every extended
// entry with a value, including token-only rows without an authType.
func normalizeHTTPRawExtendedAuthTokens(raw map[string]any) {
	if raw["type"] != "http_raw" {
		return
	}
	tokens, ok := raw["authTokensExt"].([]any)
	if !ok || len(tokens) == 0 {
		return
	}
	filtered := make([]any, 0, len(tokens))
	for _, token := range tokens {
		if hasHTTPAuthTokenValue(token) {
			filtered = append(filtered, token)
		}
	}
	if len(filtered) == len(tokens) {
		return
	}
	if len(filtered) == 0 {
		delete(raw, "authTokensExt")
		return
	}
	raw["authTokensExt"] = filtered
}

func hasHTTPAuthTokenValue(value any) bool {
	switch v := value.(type) {
	case nil:
		return false
	case string:
		return v != ""
	case map[string]any:
		for _, child := range v {
			if hasHTTPAuthTokenValue(child) {
				return true
			}
		}
		return false
	case []any:
		for _, child := range v {
			if hasHTTPAuthTokenValue(child) {
				return true
			}
		}
		return false
	default:
		return true
	}
}

// Clear placeholders retained in state by older provider versions, without
// discarding configured values or unknown values during a refresh.
func normalizeHTTPRawExtendedAuthTokenState(tokens types.List) types.List {
	if tokens.IsNull() || tokens.IsUnknown() || len(tokens.Elements()) == 0 {
		return tokens
	}
	filtered := make([]attr.Value, 0, len(tokens.Elements()))
	for _, token := range tokens.Elements() {
		if hasHTTPAuthTokenStateValue(token) {
			filtered = append(filtered, token)
		}
	}
	if len(filtered) == len(tokens.Elements()) {
		return tokens
	}
	elementType := tokens.ElementType(context.Background())
	if len(filtered) == 0 {
		return types.ListNull(elementType)
	}
	return types.ListValueMust(elementType, filtered)
}

func hasHTTPAuthTokenStateValue(value attr.Value) bool {
	if value.IsUnknown() {
		return true
	}
	if value.IsNull() {
		return false
	}
	switch v := value.(type) {
	case types.String:
		return v.ValueString() != ""
	case types.Object:
		for _, child := range v.Attributes() {
			if hasHTTPAuthTokenStateValue(child) {
				return true
			}
		}
		return false
	case types.List:
		for _, child := range v.Elements() {
			if hasHTTPAuthTokenStateValue(child) {
				return true
			}
		}
		return false
	default:
		return true
	}
}
