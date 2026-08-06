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

func groupSystemSettingsCustomLogoPlanModifiers() []planmodifier.Object {
	return []planmodifier.Object{
		custom_objectplanmodifier.PreferConfigOrState(),
		custom_objectplanmodifier.EmptyStringsAsNull(),
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

func groupSystemSettingsObjectAfterWrite(api, planned types.Object) types.Object {
	if api.IsNull() || api.IsUnknown() {
		return planned
	}
	if planned.IsNull() || planned.IsUnknown() {
		return api
	}

	apiAttributes := api.Attributes()
	plannedAttributes := planned.Attributes()
	merged := make(map[string]attr.Value, len(apiAttributes))
	for name, apiValue := range apiAttributes {
		merged[name] = apiValue
	}
	for name, plannedValue := range plannedAttributes {
		apiValue, ok := merged[name]
		if !ok || apiValue == nil || apiValue.IsNull() || apiValue.IsUnknown() {
			merged[name] = plannedValue
			continue
		}
		apiObject, apiIsObject := apiValue.(types.Object)
		plannedObject, plannedIsObject := plannedValue.(types.Object)
		if apiIsObject && plannedIsObject {
			merged[name] = groupSystemSettingsObjectAfterWrite(apiObject, plannedObject)
		}
	}

	value, diagnostics := types.ObjectValue(api.AttributeTypes(context.Background()), merged)
	if diagnostics.HasError() {
		return planned
	}
	return value
}

func groupSystemSettingsAPIAfterWrite(api, planned *GroupSystemSettingsModel) *GroupSystemSettingsModel {
	if api == nil || planned == nil {
		return api
	}
	api.API = groupSystemSettingsObjectAfterWrite(api.API, planned.API)
	api.Apps = groupSystemSettingsObjectAfterWrite(api.Apps, planned.Apps)
	api.Backups = groupSystemSettingsObjectAfterWrite(api.Backups, planned.Backups)
	api.CustomLogo = groupSystemSettingsObjectAfterWrite(api.CustomLogo, planned.CustomLogo)
	api.Pii = groupSystemSettingsObjectAfterWrite(api.Pii, planned.Pii)
	api.Proxy = groupSystemSettingsObjectAfterWrite(api.Proxy, planned.Proxy)
	api.Rollback = groupSystemSettingsObjectAfterWrite(api.Rollback, planned.Rollback)
	api.Shutdown = groupSystemSettingsObjectAfterWrite(api.Shutdown, planned.Shutdown)
	api.Sni = groupSystemSettingsObjectAfterWrite(api.Sni, planned.Sni)
	api.Sockets = groupSystemSettingsObjectAfterWrite(api.Sockets, planned.Sockets)
	api.Support = groupSystemSettingsObjectAfterWrite(api.Support, planned.Support)
	api.System = groupSystemSettingsObjectAfterWrite(api.System, planned.System)
	api.TLS = groupSystemSettingsObjectAfterWrite(api.TLS, planned.TLS)
	api.UpgradeGroupSettings = groupSystemSettingsObjectAfterWrite(api.UpgradeGroupSettings, planned.UpgradeGroupSettings)
	api.UpgradeSettings = groupSystemSettingsObjectAfterWrite(api.UpgradeSettings, planned.UpgradeSettings)
	api.Workers = groupSystemSettingsObjectAfterWrite(api.Workers, planned.Workers)
	return api
}
