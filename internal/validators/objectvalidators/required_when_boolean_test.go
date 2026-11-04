package objectvalidators

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestRequiredWhenBoolean(t *testing.T) {
	for _, test := range []struct {
		name      string
		enabled   types.Bool
		earliest  types.String
		condition bool
		wantError bool
	}{
		{"enabled-missing", types.BoolValue(true), types.StringNull(), true, true},
		{"enabled-configured", types.BoolValue(true), types.StringValue("-30d"), true, false},
		{"disabled-missing", types.BoolValue(false), types.StringNull(), true, false},
		{"omitted-trigger", types.BoolNull(), types.StringNull(), true, false},
		{"unknown-trigger", types.BoolUnknown(), types.StringNull(), true, false},
		{"unknown-required", types.BoolValue(true), types.StringUnknown(), true, false},
		{"false-condition", types.BoolValue(false), types.StringNull(), false, true},
	} {
		t.Run(test.name, func(t *testing.T) {
			object, diagnostics := types.ObjectValue(map[string]attr.Type{"enabled": types.BoolType, "earliest": types.StringType},
				map[string]attr.Value{"enabled": test.enabled, "earliest": test.earliest})
			if diagnostics.HasError() {
				t.Fatal(diagnostics)
			}
			var response validator.ObjectResponse
			RequiredWhenBoolean("enabled", test.condition, "earliest").ValidateObject(context.Background(), validator.ObjectRequest{
				Path: path.Root("metadata"), ConfigValue: object,
			}, &response)
			if response.Diagnostics.HasError() != test.wantError {
				t.Fatalf("diagnostics = %v, want error %t", response.Diagnostics, test.wantError)
			}
			if test.wantError {
				want := "metadata.earliest is required when metadata.enabled is "
				if test.condition {
					want += "true."
				} else {
					want += "false."
				}
				if response.Diagnostics[0].Detail() != want {
					t.Fatalf("detail = %q, want %q", response.Diagnostics[0].Detail(), want)
				}
			}
		})
	}
}
