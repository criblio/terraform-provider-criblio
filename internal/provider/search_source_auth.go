package provider

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"regexp"

	"github.com/criblio/terraform-provider-criblio/internal/restclient"
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
		protocol, _ := body["type"].(string)
		if err := validateSearchSourceAuthToken(token, protocol); err != nil {
			return fmt.Errorf("auth_tokens[%d]: %w", i, err)
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
		for _, item := range tokens {
			token := item.(map[string]any)
			delete(token, "enabled")
			if token["tokenSecret"] != nil {
				token["authType"] = "secret"
			}
		}
		// Clear both representations when removing tokens. The input runtime
		// may backport extended tokens into authTokens during initialization.
		body["authTokens"] = []any{}
		body["authTokensExt"] = tokens
	case "splunk_hec":
		for _, item := range tokens {
			token := item.(map[string]any)
			if token["tokenSecret"] != nil {
				token["authType"] = "secret"
			}
		}
	}
	return nil
}

func validateSearchSourceAuthToken(token map[string]any, protocol string) error {
	value, hasToken := token["token"].(string)
	secret, hasSecret := token["tokenSecret"].(string)
	if hasToken == hasSecret {
		return errors.New("configure exactly one of token or token_secret")
	}
	if (hasToken && value == "") || (hasSecret && secret == "") {
		return errors.New("token or token_secret must not be empty")
	}
	if hasSecret && protocol != "http_raw" && protocol != "splunk_hec" {
		return errors.New("token_secret is supported only for Raw HTTP and Splunk HEC; use token for this protocol")
	}
	if secret == "****/******" || (protocol == "cribl_http" && value == "****/******") {
		return errors.New("configure an accessible secret ID instead of the masked global-secret placeholder")
	}
	if enabled, ok := token["enabled"].(bool); protocol == "http_raw" && ok && !enabled {
		return errors.New("raw HTTP does not support disabling individual tokens; remove the token entry instead")
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
					if err := normalizeSearchSourceRawAuthToken(token); err != nil {
						return fmt.Errorf("authTokens[%d]: %w", i, err)
					}
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
			if err := normalizeSearchSourceRawAuthToken(token); err != nil {
				return fmt.Errorf("authTokensExt[%d]: %w", i, err)
			}
			if token["tokenSecret"] != nil {
				tokens = append(tokens, token)
				continue
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
				if body["type"] == "splunk_hec" {
					if err := normalizeSearchSourceRawAuthToken(token); err != nil {
						return err
					}
					continue
				}
				if body["type"] == "cribl_http" && token["tokenSecret"] != nil {
					token["token"] = token["tokenSecret"]
				} else if body["type"] != "http_raw" && token["token"] == nil {
					token["token"] = token["tokenSecret"]
				}
				if body["type"] != "http_raw" {
					// Keep the established Cribl HTTP token attribute; token_secret
					// is the additional Raw HTTP representation only.
					delete(token, "tokenSecret")
				}
			}
		}
	}
	return nil
}

func normalizeSearchSourceRawAuthToken(token map[string]any) error {
	secret, _ := token["tokenSecret"].(string)
	if token["authType"] == "secret" || (token["authType"] != "manual" && secret != "") {
		if secret == "" {
			return errors.New("stored-secret reference is missing tokenSecret")
		}
		// A resolved runtime token is not the configured credential. Keep only
		// the reference so refresh/import never converts it to an inline value.
		delete(token, "token")
	} else {
		// Manual rows may include an empty secret selector from the UI.
		delete(token, "tokenSecret")
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
			if isLikelyMaskedSecret(secret) || (sourceType == "cribl_http" && secret == "****/******") || (sourceType != "cribl_http" && searchSourceEncryptedToken.MatchString(secret)) {
				attributes["token"] = priorToken
			}
		}
		apiReference := attributes["token_secret"].(types.String)
		priorReference := previous["token_secret"].(types.String)
		if apiReference.ValueString() == "****/******" && !priorReference.IsNull() && !priorReference.IsUnknown() {
			attributes["token_secret"] = priorReference
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

var _ resource.ResourceWithValidateConfig = &SearchSourceResource{}

// ValidateConfig checks protocol-specific token choices before making API requests.
func (r *SearchSourceResource) ValidateConfig(ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse) {
	var sourceType types.String
	var tokens types.List
	resp.Diagnostics.Append(req.Config.GetAttribute(ctx, path.Root("type"), &sourceType)...)
	resp.Diagnostics.Append(req.Config.GetAttribute(ctx, path.Root("auth_tokens"), &tokens)...)
	if resp.Diagnostics.HasError() || sourceType.IsUnknown() || tokens.IsNull() || tokens.IsUnknown() {
		return
	}
	for i, item := range tokens.Elements() {
		object := item.(types.Object)
		if object.IsNull() {
			resp.Diagnostics.AddAttributeError(path.Root("auth_tokens").AtListIndex(i), "Invalid authentication token", "Authentication token entries must not be null.")
			continue
		}
		if object.IsUnknown() {
			continue
		}
		unknown := false
		for _, attribute := range object.Attributes() {
			unknown = unknown || attribute.IsUnknown()
		}
		if unknown {
			continue
		}
		value, err := SearchSourceTerraformValueToJSON(object)
		if err == nil {
			err = validateSearchSourceAuthToken(value.(map[string]any), sourceType.ValueString())
		}
		if err != nil {
			resp.Diagnostics.AddAttributeError(path.Root("auth_tokens").AtListIndex(i), "Invalid authentication token", err.Error())
		}
	}
}

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

// searchSourceAuthRequestModel marks unowned credentials for preservation from
// the raw API representation rather than re-encoding their Terraform values.
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

// updateWithPreservedAuth respects the v1 API's full-replacement semantics.
// Keep unowned credentials and protocol-specific properties from the backend,
// removing inactive runtime values to respect its exclusive token-row schema.
func (a SearchSourceAPI) updateWithPreservedAuth(ctx context.Context, model SearchSourceModel) (*SearchSourceModel, error) {
	if !model.AuthTokens.IsNull() {
		return a.Update(ctx, model)
	}
	path := fmt.Sprintf("/m/default_search/search/local_search/sources/%s", model.ID.ValueString())
	existing, err := restclient.Get[map[string]any](ctx, a.client, path)
	if err != nil {
		return nil, err
	}
	if existing == nil {
		return nil, errors.New("source response is empty; cannot preserve authentication")
	}
	encoded, err := json.Marshal(model)
	if err != nil {
		return nil, err
	}
	var body map[string]any
	if err := json.Unmarshal(encoded, &body); err != nil {
		return nil, err
	}
	for _, field := range []string{"authTokens", "authTokensExt"} {
		if tokens, ok := (*existing)[field]; ok {
			if field == "authTokens" && model.Type.ValueString() == "http_raw" && (*existing)["authTokensExt"] != nil {
				body[field] = []any{}
				continue
			}
			if entries, ok := tokens.([]any); ok {
				for _, item := range entries {
					if token, ok := item.(map[string]any); ok && (model.Type.ValueString() == "http_raw" || model.Type.ValueString() == "splunk_hec") {
						if err := normalizeSearchSourceRawAuthToken(token); err != nil {
							return nil, err
						}
						if token["tokenSecret"] != nil {
							token["authType"] = "secret"
						}
					}
				}
			}
			if err := validatePreservedSearchSourceAuth(tokens); err != nil {
				return nil, err
			}
			body[field] = tokens
		}
	}
	return restclient.Patch[map[string]any, SearchSourceModel](ctx, a.client, path, body)
}

func validatePreservedSearchSourceAuth(value any) error {
	switch value := value.(type) {
	case string:
		if isLikelyMaskedSecret(value) || value == "****/******" {
			return errors.New("cannot safely preserve masked authentication during a source update; grant access to the secret folder or configure auth_tokens with known credentials")
		}
	case []any:
		for _, item := range value {
			if err := validatePreservedSearchSourceAuth(item); err != nil {
				return err
			}
		}
	case map[string]any:
		for _, field := range []string{"token", "tokenSecret"} {
			if err := validatePreservedSearchSourceAuth(value[field]); err != nil {
				return err
			}
		}
	}
	return nil
}
