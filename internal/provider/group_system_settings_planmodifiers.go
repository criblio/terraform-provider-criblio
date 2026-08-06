package provider

import (
	"context"

	custom_objectplanmodifier "github.com/criblio/terraform-provider-criblio/internal/tfplanmodifiers/objectplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func groupSystemSettingsEmptyObjectPlanModifiers() []planmodifier.Object {
	return []planmodifier.Object{
		custom_objectplanmodifier.PreferConfigOrState(),
	}
}

func groupSystemSettingsObjectFromAPIOrPrior(api, prior types.Object) types.Object {
	if api.IsNull() || api.IsUnknown() {
		return api
	}
	if prior.IsNull() || prior.IsUnknown() {
		return api
	}

	apiAttributes := api.Attributes()
	priorAttributes := prior.Attributes()
	merged := make(map[string]attr.Value, len(apiAttributes))
	for name, apiValue := range apiAttributes {
		merged[name] = apiValue
	}
	for name, priorValue := range priorAttributes {
		apiValue, ok := merged[name]
		if !ok || apiValue == nil {
			continue
		}

		apiObject, apiIsObject := apiValue.(types.Object)
		priorObject, priorIsObject := priorValue.(types.Object)
		if apiIsObject && priorIsObject {
			merged[name] = groupSystemSettingsObjectFromAPIOrPrior(apiObject, priorObject)
			continue
		}

		if name == "passphrase" {
			apiString, apiIsString := apiValue.(types.String)
			priorString, priorIsString := priorValue.(types.String)
			if apiIsString && priorIsString {
				if apiString.IsNull() || apiString.IsUnknown() {
					merged[name] = priorString
				} else {
					merged[name] = stringFromAPIOrPrior(apiString.ValueString(), priorString)
				}
			}
		}
	}

	value, diagnostics := types.ObjectValue(api.AttributeTypes(context.Background()), merged)
	if diagnostics.HasError() {
		return api
	}
	return value
}
