// Package converter: model type registry for the import CLI.
// ResourceModelTypes returns reflect.Type for each provider resource model so the
// converter can instantiate models and call RefreshFrom* via reflection.
// Uses existing provider types only; no provider code changes required.
// When adding a new resource to the provider, add a corresponding entry here.
package converter

import (
	"reflect"

	"github.com/criblio/terraform-provider-criblio/internal/provider"
)

// OneOfBlockNamesFromModel returns the tfsdk attribute names for oneOf-style nested blocks
// (pointer-to-struct fields with tfsdk tag), excluding "id" and "items". Used to derive
// supported block names from the provider model so unsupported API types can be skipped dynamically.
func OneOfBlockNamesFromModel(modelTypeName string) ([]string, error) {
	types := ResourceModelTypes()
	typ, ok := types[modelTypeName]
	if !ok {
		return nil, nil
	}
	var names []string
	for i := 0; i < typ.NumField(); i++ {
		f := typ.Field(i)
		tfsdk := f.Tag.Get("tfsdk")
		if tfsdk == "" || tfsdk == "id" || tfsdk == "items" {
			continue
		}
		if f.Type.Kind() == reflect.Ptr {
			elm := f.Type.Elem()
			if elm.Kind() == reflect.Struct {
				names = append(names, tfsdk)
			}
		}
	}
	return names, nil
}

// AllAttributeNamesFromModel returns all tfsdk attribute names for the given resource model type.
// Used to filter generated HCL so only schema-supported attributes are emitted (unsupported fields ignored).
func AllAttributeNamesFromModel(modelTypeName string) []string {
	types := ResourceModelTypes()
	typ, ok := types[modelTypeName]
	if !ok {
		return nil
	}
	var names []string
	for i := 0; i < typ.NumField(); i++ {
		tfsdk := typ.Field(i).Tag.Get("tfsdk")
		if tfsdk != "" && tfsdk != "-" {
			names = append(names, tfsdk)
		}
	}
	return names
}

// ResourceModelTypes returns the reflect.Type (struct type) for each resource model
// by its ModelTypeName (e.g. "SourceResourceModel"). Used by the converter to
// instantiate provider models and invoke RefreshFrom* methods.
func ResourceModelTypes() map[string]reflect.Type {
	return map[string]reflect.Type{
		"AppscopeConfigResourceModel":              reflect.TypeOf((*provider.AppscopeConfigResourceModel)(nil)).Elem(),
		"CertificateResourceModel":                reflect.TypeOf((*provider.CertificateResourceModel)(nil)).Elem(),
		"CollectorResourceModel":                  reflect.TypeOf((*provider.CollectorResourceModel)(nil)).Elem(),
		"CommitResourceModel":                      reflect.TypeOf((*provider.CommitResourceModel)(nil)).Elem(),
		"CriblLakeDatasetResourceModel":           reflect.TypeOf((*provider.CriblLakeDatasetResourceModel)(nil)).Elem(),
		"CriblLakeHouseResourceModel":             reflect.TypeOf((*provider.CriblLakeHouseResourceModel)(nil)).Elem(),
		"DatabaseConnectionResourceModel":         reflect.TypeOf((*provider.DatabaseConnectionResourceModel)(nil)).Elem(),
		"DeployResourceModel":                      reflect.TypeOf((*provider.DeployResourceModel)(nil)).Elem(),
		"DestinationResourceModel":                reflect.TypeOf((*provider.DestinationResourceModel)(nil)).Elem(),
		"EventBreakerRulesetResourceModel":        reflect.TypeOf((*provider.EventBreakerRulesetResourceModel)(nil)).Elem(),
		"GlobalVarResourceModel":                  reflect.TypeOf((*provider.GlobalVarResourceModel)(nil)).Elem(),
		"GrokResourceModel":                       reflect.TypeOf((*provider.GrokResourceModel)(nil)).Elem(),
		"GroupResourceModel":                       reflect.TypeOf((*provider.GroupResourceModel)(nil)).Elem(),
		"GroupSystemSettingsResourceModel":         reflect.TypeOf((*provider.GroupSystemSettingsResourceModel)(nil)).Elem(),
		"HmacFunctionResourceModel":               reflect.TypeOf((*provider.HmacFunctionResourceModel)(nil)).Elem(),
		"KeyResourceModel":                         reflect.TypeOf((*provider.KeyResourceModel)(nil)).Elem(),
		"LakehouseDatasetConnectionResourceModel": reflect.TypeOf((*provider.LakehouseDatasetConnectionResourceModel)(nil)).Elem(),
		"LookupFileResourceModel":                 reflect.TypeOf((*provider.LookupFileResourceModel)(nil)).Elem(),
		"MappingRulesetResourceModel":             reflect.TypeOf((*provider.MappingRulesetResourceModel)(nil)).Elem(),
		"NotificationResourceModel":               reflect.TypeOf((*provider.NotificationResourceModel)(nil)).Elem(),
		"NotificationTargetResourceModel":         reflect.TypeOf((*provider.NotificationTargetResourceModel)(nil)).Elem(),
		"PackResourceModel":                        reflect.TypeOf((*provider.PackResourceModel)(nil)).Elem(),
		"PackBreakersResourceModel":               reflect.TypeOf((*provider.PackBreakersResourceModel)(nil)).Elem(),
		"PackDestinationResourceModel":            reflect.TypeOf((*provider.PackDestinationResourceModel)(nil)).Elem(),
		"PackLookupsResourceModel":                reflect.TypeOf((*provider.PackLookupsResourceModel)(nil)).Elem(),
		"PackPipelineResourceModel":               reflect.TypeOf((*provider.PackPipelineResourceModel)(nil)).Elem(),
		"PackRoutesResourceModel":                 reflect.TypeOf((*provider.PackRoutesResourceModel)(nil)).Elem(),
		"PackSourceResourceModel":                 reflect.TypeOf((*provider.PackSourceResourceModel)(nil)).Elem(),
		"PackVarsResourceModel":                   reflect.TypeOf((*provider.PackVarsResourceModel)(nil)).Elem(),
		"ParquetSchemaResourceModel":              reflect.TypeOf((*provider.ParquetSchemaResourceModel)(nil)).Elem(),
		"ParserLibEntryResourceModel":             reflect.TypeOf((*provider.ParserLibEntryResourceModel)(nil)).Elem(),
		"PipelineResourceModel":                   reflect.TypeOf((*provider.PipelineResourceModel)(nil)).Elem(),
		"ProjectResourceModel":                     reflect.TypeOf((*provider.ProjectResourceModel)(nil)).Elem(),
		"RegexResourceModel":                      reflect.TypeOf((*provider.RegexResourceModel)(nil)).Elem(),
		"RoutesResourceModel":                      reflect.TypeOf((*provider.RoutesResourceModel)(nil)).Elem(),
		"SchemaResourceModel":                      reflect.TypeOf((*provider.SchemaResourceModel)(nil)).Elem(),
		"SearchDashboardResourceModel":            reflect.TypeOf((*provider.SearchDashboardResourceModel)(nil)).Elem(),
		"SearchDashboardCategoryResourceModel":    reflect.TypeOf((*provider.SearchDashboardCategoryResourceModel)(nil)).Elem(),
		"SearchDatasetResourceModel":              reflect.TypeOf((*provider.SearchDatasetResourceModel)(nil)).Elem(),
		"SearchDatasetProviderResourceModel":      reflect.TypeOf((*provider.SearchDatasetProviderResourceModel)(nil)).Elem(),
		"SearchMacroResourceModel":                reflect.TypeOf((*provider.SearchMacroResourceModel)(nil)).Elem(),
		"SearchSavedQueryResourceModel":           reflect.TypeOf((*provider.SearchSavedQueryResourceModel)(nil)).Elem(),
		"SearchUsageGroupResourceModel":           reflect.TypeOf((*provider.SearchUsageGroupResourceModel)(nil)).Elem(),
		"SecretResourceModel":                     reflect.TypeOf((*provider.SecretResourceModel)(nil)).Elem(),
		"SourceResourceModel":                     reflect.TypeOf((*provider.SourceResourceModel)(nil)).Elem(),
		"SubscriptionResourceModel":               reflect.TypeOf((*provider.SubscriptionResourceModel)(nil)).Elem(),
		"WorkspaceResourceModel":                  reflect.TypeOf((*provider.WorkspaceResourceModel)(nil)).Elem(),
	}
}
