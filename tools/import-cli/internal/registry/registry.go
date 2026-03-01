// Package registry discovers Terraform resource types by reflecting over the
// provider's Resources() definitions, so the CLI stays in sync with the provider.
// SDK/list/get and import ID metadata come from import_metadata.go in this package.
package registry

import (
	"context"
	"reflect"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/resource"
)

const providerTypeName = "criblio"

// Entry holds a single Terraform resource type and its model type information,
// plus SDK List/Get method names and Terraform import ID format (from ImportState logic).
type Entry struct {
	// TypeName is the full Terraform resource type (e.g. criblio_source).
	// It matches the provider definition exactly.
	TypeName string
	// ModelTypeName is the Go type name of the resource's state model (e.g. SourceResourceModel).
	// Used for reflection-based conversion and HCL generation.
	ModelTypeName string
	// SDKService is the name of the service field on sdk.CriblIo (e.g. "Inputs", "Pipelines").
	// Empty if not known; used by discovery to call List*.
	SDKService string
	// ListMethod is the SDK method name to list resources (e.g. "ListInput", "ListPipeline").
	// Empty if not known or not applicable.
	ListMethod string
	// GetMethod is the SDK method name to get a single resource by ID (e.g. "GetInputByID").
	// Empty if not known or not applicable.
	GetMethod string
	// ImportIDFormat describes the Terraform import ID format, matching provider ImportState.
	// Examples: "json:group_id,id", "id", "json:group_id,id,pack". Empty if not known.
	ImportIDFormat string
	// OneOf configures oneOf-style resources; when set, export emits a type-specific block from Items.
	OneOf *OneOfConfig
	// RefreshFromMethod overrides the RefreshFrom* method name used when converting Get response to model. Empty = derive from GetMethod.
	RefreshFromMethod string
	// ListItemIDMethod is the method name on list items to get the ID (e.g. "GetKeyID"). Empty = use GetID or map "id"/"Id".
	ListItemIDMethod string
	// ListUseGroupIDAsItemID when true: when list is called per-group and the item has no id, use the group ID as the identifier.
	ListUseGroupIDAsItemID bool
}

// EntryOverride provides static overrides for derived SDK/list/get/ImportID metadata.
// Only non-empty fields replace the derived values for that type.
type EntryOverride struct {
	SDKService     string
	ListMethod     string
	GetMethod      string
	ImportIDFormat string
}

// Registry holds all discovered Terraform resource types from the provider.
type Registry struct {
	entries    []Entry
	byTypeName map[string]Entry
}

// MetadataFromProvider returns import metadata for use with NewFromResources.
func MetadataFromProvider() map[string]ResourceMetadata {
	return ImportMetadata()
}

// NewFromResources discovers resource types by calling each resource constructor
// and reading Metadata. Terraform type names and model type names come from the
// provider. metadata is from ImportMetadata(); overrides replace those values when set.
// oneOfBlockNamesFromModel is optional; when set it is used to populate SupportedBlockNames for oneOf resources from the provider model (e.g. converter.OneOfBlockNamesFromModel).
func NewFromResources(ctx context.Context, constructors []func() resource.Resource, metadata map[string]ResourceMetadata, overrides map[string]EntryOverride, oneOfBlockNamesFromModel func(string) ([]string, error)) (*Registry, error) {
	byTypeName := make(map[string]Entry)
	var entries []Entry

	req := resource.MetadataRequest{ProviderTypeName: providerTypeName}

	for _, newRes := range constructors {
		res := newRes()
		var resp resource.MetadataResponse
		res.Metadata(ctx, req, &resp)

		typeName := resp.TypeName
		if typeName == "" {
			continue
		}

		modelTypeName := modelTypeNameFromResource(res)
		e := Entry{TypeName: typeName, ModelTypeName: modelTypeName}

		if meta, ok := metadata[typeName]; ok {
			e.SDKService = meta.SDKService
			e.ListMethod = meta.ListMethod
			e.GetMethod = meta.GetMethod
			e.ImportIDFormat = meta.ImportIDFormat
			e.OneOf = meta.OneOf
			e.RefreshFromMethod = meta.RefreshFromMethod
			e.ListItemIDMethod = meta.ListItemIDMethod
			e.ListUseGroupIDAsItemID = meta.ListUseGroupIDAsItemID
			// Populate OneOf.SupportedBlockNames from provider model when callback provided.
			if oneOfBlockNamesFromModel != nil && meta.OneOf != nil && modelTypeName != "" {
				if names, err := oneOfBlockNamesFromModel(modelTypeName); err == nil && len(names) > 0 {
					oneOfCopy := *meta.OneOf
					oneOfCopy.SupportedBlockNames = names
					e.OneOf = &oneOfCopy
				}
			}
		}
		if o, ok := overrides[typeName]; ok {
			if o.SDKService != "" {
				e.SDKService = o.SDKService
			}
			if o.ListMethod != "" {
				e.ListMethod = o.ListMethod
			}
			if o.GetMethod != "" {
				e.GetMethod = o.GetMethod
			}
			if o.ImportIDFormat != "" {
				e.ImportIDFormat = o.ImportIDFormat
			}
		}

		entries = append(entries, e)
		byTypeName[typeName] = e
	}

	return &Registry{entries: entries, byTypeName: byTypeName}, nil
}

// modelTypeNameFromResource derives the model type name from the resource's
// concrete type (e.g. *SourceResource -> "SourceResourceModel").
func modelTypeNameFromResource(res resource.Resource) string {
	typ := reflect.TypeOf(res)
	if typ.Kind() == reflect.Pointer {
		typ = typ.Elem()
	}
	name := typ.Name()
	// Provider convention: SourceResource -> SourceResourceModel
	if strings.HasSuffix(name, "Resource") {
		return strings.TrimSuffix(name, "Resource") + "ResourceModel"
	}
	return name + "Model"
}

// Entries returns all registry entries in discovery order.
func (r *Registry) Entries() []Entry {
	return append([]Entry(nil), r.entries...)
}

// ByTypeName returns the entry for the given Terraform type name, or false if not found.
func (r *Registry) ByTypeName(typeName string) (Entry, bool) {
	e, ok := r.byTypeName[typeName]
	return e, ok
}

// TypeNames returns all Terraform resource type names.
func (r *Registry) TypeNames() []string {
	out := make([]string, len(r.entries))
	for i, e := range r.entries {
		out[i] = e.TypeName
	}
	return out
}

// Len returns the number of registered resource types.
func (r *Registry) Len() int {
	return len(r.entries)
}
