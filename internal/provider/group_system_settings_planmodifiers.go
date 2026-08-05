package provider

import (
	custom_objectplanmodifier "github.com/criblio/terraform-provider-criblio/internal/tfplanmodifiers/objectplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
)

func groupSystemSettingsEmptyObjectPlanModifiers() []planmodifier.Object {
	return []planmodifier.Object{
		custom_objectplanmodifier.EmptyObjectAsNull(),
	}
}
