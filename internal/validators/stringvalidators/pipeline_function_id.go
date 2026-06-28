package stringvalidators

import (
	"context"
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/criblio/terraform-provider-criblio/internal/auth"
	"github.com/criblio/terraform-provider-criblio/internal/restclient"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

// IsCriblPipelineFunctionIDWithRestClient returns a validator that fetches
// valid function IDs with the migrated REST client.
func IsCriblPipelineFunctionIDWithRestClient(client **restclient.Client) validator.String {
	return &criblPipelineFunctionIDRestValidator{client: client}
}

type criblPipelineFunctionIDRestValidator struct {
	client **restclient.Client

	mu        sync.Mutex
	cachedIDs map[string]struct{}
}

func (v *criblPipelineFunctionIDRestValidator) Description(_ context.Context) string {
	return "value must be a valid Cribl pipeline function ID"
}

func (v *criblPipelineFunctionIDRestValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
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

func (v *criblPipelineFunctionIDRestValidator) resolveClient() *restclient.Client {
	if v.client != nil && *v.client != nil {
		return *v.client
	}
	credentials, err := auth.GetCredentials()
	if err != nil {
		return nil
	}
	return restclient.New(restclient.Config{
		Credentials: credentials,
		BearerToken: os.Getenv("CRIBL_BEARER_TOKEN"),
	})
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
