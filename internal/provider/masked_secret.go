package provider

import "github.com/hashicorp/terraform-plugin-framework/types"

// isLikelyMaskedSecret reports whether the API returned a redacted placeholder instead of the
// real secret (common pattern: all asterisks).
func isLikelyMaskedSecret(s string) bool {
	if len(s) < 4 {
		return false
	}
	for _, r := range s {
		if r != '*' {
			return false
		}
	}
	return true
}

// stringFromAPIOrPrior maps an API string into Terraform state. When the API returns a masked
// value and we already have a known non-null value from state/plan, keep the prior value so
// config and state stay aligned (avoids perpetual drift on refresh).
func stringFromAPIOrPrior(api string, prior types.String) types.String {
	if isLikelyMaskedSecret(api) && !prior.IsUnknown() && !prior.IsNull() {
		return prior
	}
	return types.StringValue(api)
}

func stringPointerFromAPIOrPrior(api *string, prior types.String) types.String {
	if api == nil {
		return types.StringNull()
	}
	return stringFromAPIOrPrior(*api, prior)
}
