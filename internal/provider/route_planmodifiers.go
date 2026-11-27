package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type routeIDsByName struct{}

func routeListPlanModifiers() []planmodifier.List {
	return []planmodifier.List{routeIDsByName{}}
}

func (routeIDsByName) Description(context.Context) string {
	return "Keeps the server-assigned ID with its route when routes change order."
}

func (m routeIDsByName) MarkdownDescription(ctx context.Context) string {
	return m.Description(ctx)
}

func (routeIDsByName) PlanModifyList(ctx context.Context, req planmodifier.ListRequest, resp *planmodifier.ListResponse) {
	if req.StateValue.IsNull() || req.StateValue.IsUnknown() || req.PlanValue.IsNull() || req.PlanValue.IsUnknown() || req.PlanValue.Equal(req.StateValue) {
		return
	}
	prior := req.StateValue.Elements()
	planned := req.PlanValue.Elements()
	configured := req.ConfigValue.Elements()
	if len(configured) != len(planned) {
		return
	}
	matches := matchOrderedConfig(configured, prior, "name", routeMatchesSettings)
	oldNames, newNames := make(map[string]bool), make(map[string]int)
	for _, row := range prior {
		oldNames[objectString(row, "name")] = true
	}
	for _, row := range configured {
		newNames[objectString(row, "name")]++
	}
	for i, element := range planned {
		route := element.(types.Object)
		if route.IsNull() || route.IsUnknown() {
			continue
		}
		index := matches[i]
		name := objectString(configured[i], "name")
		if index < 0 && len(prior) == len(planned) && name != "" && !oldNames[name] && newNames[name] == 1 {
			// Preserve an in-place rename if the old route disappeared and its settings match.
			oldName := objectString(prior[i], "name")
			if oldName != "" && newNames[oldName] == 0 && routeMatchesSettings(configured[i].(types.Object), prior[i].(types.Object)) {
				index = i
			}
		}
		attrs := route.Attributes()
		attrs["id"] = types.StringUnknown()
		if index >= 0 {
			if id := prior[index].(types.Object).Attributes()["id"]; !id.IsNull() {
				attrs["id"] = id
			}
		}
		updated, diags := types.ObjectValue(route.AttributeTypes(ctx), attrs)
		resp.Diagnostics.Append(diags...)
		planned[i] = updated
	}
	value, diags := types.ListValue(req.PlanValue.ElementType(ctx), planned)
	resp.Diagnostics.Append(diags...)
	resp.PlanValue = value
}

// Compare only configured values; the API can populate omitted defaults.
func routeMatchesSettings(config, prior types.Object) bool {
	if config.IsNull() || config.IsUnknown() {
		return false
	}
	old := prior.Attributes()
	for key, value := range config.Attributes() {
		if key != "name" && key != "id" && !value.IsNull() && !value.IsUnknown() && !value.Equal(old[key]) {
			return false
		}
	}
	return true
}
