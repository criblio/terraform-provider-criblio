package stringvalidators

import (
	"context"
	"fmt"
	"strings"
	"sync"

	"github.com/criblio/terraform-provider-criblio/internal/restclient"
	"github.com/criblio/terraform-provider-criblio/internal/sdk"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

// IsCriblPipelineFunctionID returns a validator that fetches valid function IDs
// from the Cribl API at plan time. client is a pointer to the resource's SDK
// client field; it will be nil until Configure runs, at which point
// ValidateString dereferences it to make the API call.
//
// Results are cached in the validator instance so the API is only called once
// per resource lifetime regardless of how many function entries are validated.
func IsCriblPipelineFunctionID(client **sdk.CriblIo) validator.String {
	return &criblPipelineFunctionIDValidator{client: client}
}

// IsCriblPipelineFunctionIDWithRestClient returns a validator that fetches
// valid function IDs with the migrated REST client.
func IsCriblPipelineFunctionIDWithRestClient(client **restclient.Client) validator.String {
	return &criblPipelineFunctionIDRestValidator{client: client}
}

type criblPipelineFunctionIDValidator struct {
	client **sdk.CriblIo

	mu        sync.Mutex
	cachedIDs map[string]struct{}
}

type criblPipelineFunctionIDRestValidator struct {
	client **restclient.Client

	mu        sync.Mutex
	cachedIDs map[string]struct{}
}

func (v *criblPipelineFunctionIDValidator) Description(_ context.Context) string {
	return "value must be a valid Cribl pipeline function ID"
}

func (v *criblPipelineFunctionIDValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

func (v *criblPipelineFunctionIDRestValidator) Description(_ context.Context) string {
	return "value must be a valid Cribl pipeline function ID"
}

func (v *criblPipelineFunctionIDRestValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

func (v *criblPipelineFunctionIDValidator) ValidateString(ctx context.Context, req validator.StringRequest, resp *validator.StringResponse) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}
	value := req.ConfigValue.ValueString()

	ids, err := v.functionIDs(ctx)
	if err != nil {
		// Warn rather than error so a transient API failure doesn't block plan.
		resp.Diagnostics.AddAttributeWarning(
			req.Path,
			"Could not validate pipeline function ID",
			fmt.Sprintf("Failed to fetch function list from Cribl API: %s. Skipping validation.", err),
		)
		return
	}

	if _, ok := ids[value]; !ok {
		resp.Diagnostics.AddAttributeError(
			req.Path,
			"Invalid pipeline function ID",
			fmt.Sprintf("%s: %q is not a known Cribl pipeline function ID. Known IDs: %s", req.Path, value, sortedKeys(ids)),
		)
	}
}

func (v *criblPipelineFunctionIDRestValidator) ValidateString(ctx context.Context, req validator.StringRequest, resp *validator.StringResponse) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}
	value := req.ConfigValue.ValueString()

	ids, err := v.functionIDs(ctx)
	if err != nil {
		resp.Diagnostics.AddAttributeWarning(
			req.Path,
			"Could not validate pipeline function ID",
			fmt.Sprintf("Failed to fetch function list from Cribl API: %s. Skipping validation.", err),
		)
		return
	}

	if _, ok := ids[value]; !ok {
		resp.Diagnostics.AddAttributeError(
			req.Path,
			"Invalid pipeline function ID",
			fmt.Sprintf("%s: %q is not a known Cribl pipeline function ID. Known IDs: %s", req.Path, value, sortedKeys(ids)),
		)
	}
}

func (v *criblPipelineFunctionIDValidator) functionIDs(ctx context.Context) (map[string]struct{}, error) {
	v.mu.Lock()
	defer v.mu.Unlock()

	if v.cachedIDs != nil {
		return v.cachedIDs, nil
	}

	client := v.resolveClient()
	if client == nil {
		return nil, fmt.Errorf("SDK client not configured and no CRIBL_* environment variables found")
	}

	resp, err := client.Functions.ListFunction(ctx)
	if err != nil {
		return nil, err
	}
	if resp.Object == nil {
		return nil, fmt.Errorf("empty response (status %d)", resp.StatusCode)
	}

	ids := make(map[string]struct{}, len(resp.Object.Items))
	for _, fn := range resp.Object.Items {
		if id := strings.TrimSpace(fn.GetID()); id != "" {
			ids[id] = struct{}{}
		}
	}

	v.cachedIDs = ids
	return ids, nil
}

func (v *criblPipelineFunctionIDRestValidator) functionIDs(ctx context.Context) (map[string]struct{}, error) {
	v.mu.Lock()
	defer v.mu.Unlock()

	if v.cachedIDs != nil {
		return v.cachedIDs, nil
	}

	client := v.resolveClient()
	if client == nil {
		return nil, fmt.Errorf("REST client not configured")
	}

	functions, err := restclient.Get[[]criblFunction](ctx, client, "/functions")
	if err != nil {
		return nil, err
	}
	if functions == nil {
		return nil, fmt.Errorf("empty response")
	}

	ids := make(map[string]struct{}, len(*functions))
	for _, fn := range *functions {
		if id := strings.TrimSpace(fn.ID); id != "" {
			ids[id] = struct{}{}
		}
	}

	v.cachedIDs = ids
	return ids, nil
}

// resolveClient returns the configured resource client if available, otherwise
// falls back to a bare sdk.New() whose CriblTerraformHook will pick up auth
// from CRIBL_* environment variables — the same path used during terraform validate.
func (v *criblPipelineFunctionIDValidator) resolveClient() *sdk.CriblIo {
	if v.client != nil && *v.client != nil {
		return *v.client
	}
	return sdk.New()
}

func (v *criblPipelineFunctionIDRestValidator) resolveClient() *restclient.Client {
	if v.client != nil && *v.client != nil {
		return *v.client
	}
	return nil
}

type criblFunction struct {
	ID string `json:"id"`
}

func sortedKeys(m map[string]struct{}) string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	// sort inline — avoid importing "sort" just for error messages
	for i := 1; i < len(keys); i++ {
		for j := i; j > 0 && keys[j] < keys[j-1]; j-- {
			keys[j], keys[j-1] = keys[j-1], keys[j]
		}
	}
	return strings.Join(keys, ", ")
}
