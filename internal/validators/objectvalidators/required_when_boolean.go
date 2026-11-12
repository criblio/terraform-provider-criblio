package objectvalidators

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type requiredWhenBoolean struct {
	property string
	value    bool
	required []string
}

// RequiredWhenBoolean requires sibling attributes when property equals value.
// Unknown values are deferred until Terraform can validate their resolved values.
func RequiredWhenBoolean(property string, value bool, required ...string) validator.Object {
	return requiredWhenBoolean{property: property, value: value, required: append([]string(nil), required...)}
}

func (v requiredWhenBoolean) Description(context.Context) string {
	return fmt.Sprintf("%s must be configured when %s is %t", strings.Join(v.required, ", "), v.property, v.value)
}

func (v requiredWhenBoolean) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

func (v requiredWhenBoolean) ValidateObject(ctx context.Context, req validator.ObjectRequest, resp *validator.ObjectResponse) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}
	attributes := req.ConfigValue.Attributes()
	trigger, ok := attributes[v.property].(types.Bool)
	if !ok {
		resp.Diagnostics.AddAttributeError(req.Path, "Invalid conditional validation", "The condition must refer to a boolean attribute: "+v.property)
		return
	}
	if trigger.IsNull() || trigger.IsUnknown() || trigger.ValueBool() != v.value {
		return
	}
	for _, name := range v.required {
		value, ok := attributes[name]
		if !ok {
			resp.Diagnostics.AddAttributeError(req.Path, "Invalid conditional validation", "The required attribute is not in the object schema: "+name)
			continue
		}
		if value.IsNull() {
			resp.Diagnostics.AddAttributeError(req.Path.AtName(name), "Missing required attribute",
				fmt.Sprintf("%s is required when %s is %t.", req.Path.AtName(name), req.Path.AtName(v.property), v.value))
		}
	}
}
