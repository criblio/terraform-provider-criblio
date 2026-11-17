package provider

import (
	"context"
	"fmt"
	"regexp"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var searchSourceEncryptedToken = regexp.MustCompile(`^#[0-9]+:[^:]*:[^#]+#$`)

// searchSourceAuthTokensToAPI adapts the shared Terraform token objects to the
// underlying input's protocol-specific authentication schema.
func searchSourceAuthTokensToAPI(body map[string]any) error {
	tokens, ok := body["authTokens"].([]any)
	if !ok {
		return nil
	}
	for i, item := range tokens {
		token, ok := item.(map[string]any)
		if !ok {
			return fmt.Errorf("auth_tokens[%d]: expected a non-null token object", i)
		}
		if _, ok := token["token"].(string); !ok {
			return fmt.Errorf("auth_tokens[%d].token: expected a known token string", i)
		}
	}
	switch body["type"] {
	case "cribl_http":
		for _, item := range tokens {
			token := item.(map[string]any)
			token["tokenSecret"] = token["token"]
			delete(token, "token")
		}
	case "http_raw":
		for i, item := range tokens {
			token := item.(map[string]any)
			if enabled, ok := token["enabled"].(bool); ok && !enabled {
				return fmt.Errorf("auth_tokens[%d].enabled: Raw HTTP does not support disabling individual tokens; remove the token entry instead", i)
			}
			delete(token, "enabled")
		}
		// Clear both representations when removing tokens. The input runtime
		// may backport extended tokens into authTokens during initialization.
		body["authTokens"] = []any{}
		body["authTokensExt"] = tokens
	}
	return nil
}

// searchSourceAuthTokensFromAPI exposes legacy Raw HTTP strings and extended
// objects through the existing Terraform list-of-objects attribute.
func searchSourceAuthTokensFromAPI(body map[string]any) error {
	if body["type"] == "http_raw" {
		legacy, hasLegacy := body["authTokens"].([]any)
		extended, hasExtended := body["authTokensExt"].([]any)
		if !hasLegacy && !hasExtended {
			return nil
		}
		tokens := make([]any, 0, len(extended))
		// BaseHttpIn treats authTokensExt as authoritative whenever present,
		// including an empty list, and mirrors its values into authTokens for
		// older workers. Reading both lists would duplicate every credential.
		if !hasExtended {
			for i, item := range legacy {
				switch token := item.(type) {
				case string:
					tokens = append(tokens, map[string]any{"token": token})
				case map[string]any:
					// Accept object responses from older Search API implementations.
					tokens = append(tokens, token)
				default:
					return fmt.Errorf("authTokens[%d]: expected token string or object", i)
				}
			}
		}
		for i, item := range extended {
			token, ok := item.(map[string]any)
			if !ok {
				return fmt.Errorf("authTokensExt[%d]: expected token object", i)
			}
			if token["authType"] == "secret" || token["tokenSecret"] != nil {
				return fmt.Errorf("authTokensExt[%d]: stored-secret references on Raw HTTP are not supported by auth_tokens", i)
			}
			if value, ok := token["token"].(string); ok && value != "" {
				tokens = append(tokens, token)
			}
		}
		body["authTokens"] = tokens
	}
	if tokens, ok := body["authTokens"].([]any); ok {
		for _, item := range tokens {
			if token, ok := item.(map[string]any); ok {
				if body["type"] == "cribl_http" && token["tokenSecret"] != nil {
					token["token"] = token["tokenSecret"]
				} else if token["token"] == nil {
					token["token"] = token["tokenSecret"]
				}
			}
		}
	}
	return nil
}

func mergeSearchSourceAuthTokenState(api, prior types.List, sourceType string) types.List {
	if api.IsNull() || api.IsUnknown() {
		if prior.IsUnknown() || prior.IsNull() {
			return types.ListNull(types.ObjectType{AttrTypes: SearchSourceAuthTokensAttrTypes()})
		}
		return prior
	}
	// The API represents cleared tokens as [], while omitted Terraform
	// configuration plans null. Preserve that distinction for consistent state
	// after apply and subsequent refreshes; populated API lists still expose drift.
	if len(api.Elements()) == 0 && prior.IsNull() {
		return prior
	}
	if prior.IsNull() || prior.IsUnknown() || len(api.Elements()) != len(prior.Elements()) {
		return api
	}
	values := make([]attr.Value, len(api.Elements()))
	for i, value := range api.Elements() {
		apiObject := value.(types.Object)
		priorObject := prior.Elements()[i].(types.Object)
		if apiObject.IsNull() || apiObject.IsUnknown() || priorObject.IsNull() || priorObject.IsUnknown() {
			values[i] = value
			continue
		}
		attributes := apiObject.Attributes()
		previous := priorObject.Attributes()
		apiToken := attributes["token"].(types.String)
		priorToken := previous["token"].(types.String)
		// The plaintext is not recoverable from a masked or encrypted response.
		// Retain the matching plan/state value, including newly rotated tokens.
		if !priorToken.IsNull() && !priorToken.IsUnknown() && !apiToken.IsNull() && !apiToken.IsUnknown() {
			secret := apiToken.ValueString()
			if isLikelyMaskedSecret(secret) || searchSourceEncryptedToken.MatchString(secret) {
				attributes["token"] = priorToken
			}
		}
		for _, name := range []string{"description", "enabled"} {
			// These optional fields are not computed: API defaults must not turn
			// an omitted nested field into a configured Terraform value.
			if previous[name].IsNull() || attributes[name].IsNull() || (sourceType == "http_raw" && name == "enabled") {
				attributes[name] = previous[name]
			}
		}
		values[i] = types.ObjectValueMust(SearchSourceAuthTokensAttrTypes(), attributes)
	}
	return types.ListValueMust(types.ObjectType{AttrTypes: SearchSourceAuthTokensAttrTypes()}, values)
}

const searchSourceAuthManagedKey = "search_source_auth_tokens_managed"

// ModifyPlan preserves unconfigured legacy/imported credentials. Private state
// records ownership only once configuration or a new resource takes control.
func (r *SearchSourceResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	if req.Plan.Raw.IsNull() {
		return
	}
	var configured types.List
	resp.Diagnostics.Append(req.Config.GetAttribute(ctx, path.Root("auth_tokens"), &configured)...)
	if resp.Diagnostics.HasError() {
		return
	}
	managed := req.State.Raw.IsNull() || !configured.IsNull()
	if req.Private != nil {
		data, diags := req.Private.GetKey(ctx, searchSourceAuthManagedKey)
		resp.Diagnostics.Append(diags...)
		managed = managed || string(data) == "true"
	}
	if resp.Diagnostics.HasError() {
		return
	}
	if managed && resp.Private != nil {
		resp.Diagnostics.Append(resp.Private.SetKey(ctx, searchSourceAuthManagedKey, []byte("true"))...)
	}
	if !configured.IsNull() {
		return
	}
	tokens := types.ListValueMust(types.ObjectType{AttrTypes: SearchSourceAuthTokensAttrTypes()}, nil)
	if !managed {
		resp.Diagnostics.Append(req.State.GetAttribute(ctx, path.Root("auth_tokens"), &tokens)...)
	}
	if !resp.Diagnostics.HasError() {
		resp.Diagnostics.Append(resp.Plan.SetAttribute(ctx, path.Root("auth_tokens"), tokens)...)
	}
}

// searchSourceAuthRequestModel omits API-computed credentials on unrelated
// updates. Resending encrypted or masked legacy values could corrupt tokens.
func searchSourceAuthRequestModel(ctx context.Context, model SearchSourceModel, req resource.UpdateRequest) (SearchSourceModel, diag.Diagnostics) {
	var configured types.List
	diags := req.Config.GetAttribute(ctx, path.Root("auth_tokens"), &configured)
	if diags.HasError() || !configured.IsNull() {
		return model, diags
	}
	managed := false
	if req.Private != nil {
		data, privateDiags := req.Private.GetKey(ctx, searchSourceAuthManagedKey)
		diags.Append(privateDiags...)
		managed = string(data) == "true"
	}
	if !managed {
		model.AuthTokens = types.ListNull(types.ObjectType{AttrTypes: SearchSourceAuthTokensAttrTypes()})
	}
	return model, diags
}
