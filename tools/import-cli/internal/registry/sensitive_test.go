package registry

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

func TestSensitiveSchemaPaths(t *testing.T) {
	resourceSchema := schema.Schema{Attributes: map[string]schema.Attribute{
		"opaque": schema.StringAttribute{Optional: true, Sensitive: true},
		"list": schema.ListNestedAttribute{Optional: true, NestedObject: schema.NestedAttributeObject{Attributes: map[string]schema.Attribute{
			"opaque": schema.StringAttribute{Optional: true, Sensitive: true},
			"label":  schema.StringAttribute{Optional: true},
		}}},
		"map": schema.MapNestedAttribute{Optional: true, NestedObject: schema.NestedAttributeObject{Attributes: map[string]schema.Attribute{
			"opaque": schema.StringAttribute{Optional: true, Sensitive: true},
		}}},
	}, Blocks: map[string]schema.Block{
		"connection": schema.SingleNestedBlock{Attributes: map[string]schema.Attribute{
			"opaque": schema.StringAttribute{Optional: true, Sensitive: true},
		}},
	}}
	paths := sensitiveSchemaPaths(resourceSchema)
	for _, path := range []string{"opaque", "list.opaque", "map.*.opaque", "connection.opaque"} {
		if !paths[path] {
			t.Errorf("missing sensitive path %s", path)
		}
	}
	if paths["list.label"] || len(paths) != 4 {
		t.Fatalf("unexpected sensitive paths: %v", paths)
	}
}
