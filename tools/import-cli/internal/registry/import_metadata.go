// Package registry provides import metadata for the import CLI: SDK service and
// method names and import ID format per resource type.
//
// When you add a new resource to the provider (internal/provider/provider.go
// Resources()), add a corresponding entry to importMetadataBase below (and
// importMetadataOverrides or importMetadataClearList if needed). A test
// ensures every provider resource type has an entry here.
package registry

// oneOf configs for resources whose API returns a single payload in a list and schema uses type-specific blocks.
// These align with OpenAPI oneOf schemas: Output -> destination/pack_destination; Input -> source; InputCollector/SavedJob -> collector; NotificationTarget -> notification_target.
// ReadOnlyAttr is the tfsdk attribute name (e.g. "items") so it can be skipped in HCL; Go field name is derived for reflection.
var (
	oneOfOutput = &OneOfConfig{
		ReadOnlyAttr:       "items",
		DiscriminatorField: "type",
		BlockNamePrefix:    "output_",
		KeysToSkip:         []string{"status"},
	}
	oneOfInput = &OneOfConfig{
		ReadOnlyAttr:       "items",
		DiscriminatorField: "type",
		BlockNamePrefix:    "input_",
		KeysToSkip:         []string{"status"},
	}
	oneOfNotificationTarget = &OneOfConfig{
		ReadOnlyAttr:       "items",
		DiscriminatorField: "type",
		BlockNamePrefix:    "",
		BlockNameSuffix:    "_target",
		KeysToSkip:         []string{"status", "on_backpressure"},
		// SupportedBlockNames populated dynamically from provider model (NotificationTargetResourceModel) in registry.
	}
	// oneOfSearchDataset and oneOfSearchDatasetProvider: GenericDataset/GenericProvider oneOf in OpenAPI; block names from provider model (e.g. api_elastic_search_dataset).
	oneOfSearchDataset = &OneOfConfig{
		ReadOnlyAttr:       "items",
		DiscriminatorField: "type",
		BlockNamePrefix:    "",
		BlockNameSuffix:    "",
		KeysToSkip:         []string{},
		// SupportedBlockNames populated dynamically from SearchDatasetResourceModel.
	}
	oneOfSearchDatasetProvider = &OneOfConfig{
		ReadOnlyAttr:       "items",
		DiscriminatorField: "type",
		BlockNamePrefix:    "",
		BlockNameSuffix:    "",
		KeysToSkip:         []string{},
		// SupportedBlockNames populated dynamically from SearchDatasetProviderResourceModel.
	}
)

// importMetadataBase is the default metadata per type (from resource SDK usage).
// Overrides and clearListMethodTypes are applied in ImportMetadata().
var importMetadataBase = map[string]ResourceMetadata{
	"criblio_appscope_config":              {SDKService: "AppscopeConfigs", ListMethod: "ListAppscopeLibEntry", GetMethod: "GetAppscopeLibEntryByID", ImportIDFormat: ""},
	"criblio_certificate":                  {SDKService: "Certificates", ListMethod: "ListCertificate", GetMethod: "GetCertificateByID", ImportIDFormat: ""},
	"criblio_collector":                    {SDKService: "SavedJobs", ListMethod: "ListCollectors", GetMethod: "GetSavedJobByID", ImportIDFormat: "", OneOf: &OneOfConfig{ReadOnlyAttr: "items", DiscriminatorField: "type", BlockNamePrefix: "input_collector_", KeysToSkip: []string{"status"}, UnsupportedDiscriminatorValues: []string{"scheduledSearch", "executor"}, NestedDiscriminatorField: "collector.type"}},
	"criblio_commit":                       {SDKService: "", ListMethod: "", GetMethod: "", ImportIDFormat: ""},
	"criblio_cribl_lake_dataset":           {SDKService: "Lake", ListMethod: "GetCriblLakeDatasetByLakeID", GetMethod: "GetCriblLakeDatasetByLakeIDAndID", ImportIDFormat: ""},
	"criblio_cribl_lake_house":             {SDKService: "LakeHouse", ListMethod: "ListDefaultLakeLakehouse", GetMethod: "GetDefaultLakeLakehouseByID", ImportIDFormat: ""},
	"criblio_database_connection":          {SDKService: "DatabaseConnections", ListMethod: "GetDatabaseConnectionConfig", GetMethod: "GetDatabaseConnectionConfigByID", ImportIDFormat: ""},
	"criblio_deploy":                       {SDKService: "", ListMethod: "", GetMethod: "", ImportIDFormat: ""},
	"criblio_destination":                  {SDKService: "Outputs", ListMethod: "ListOutput", GetMethod: "GetOutputByID", ImportIDFormat: "", OneOf: oneOfOutput},
	"criblio_event_breaker_ruleset":        {SDKService: "EventBreakerRules", ListMethod: "ListEventBreakerRuleset", GetMethod: "GetEventBreakerRulesetByID", ImportIDFormat: ""},
	"criblio_global_var":                   {SDKService: "GlobalVariables", ListMethod: "GetGlobalVariable", GetMethod: "GetGlobalVariableByID", ImportIDFormat: ""},
	"criblio_grok":                         {SDKService: "Grokfiles", ListMethod: "ListGrokFile", GetMethod: "GetGrokFileByID", ImportIDFormat: ""},
	"criblio_group":                        {SDKService: "Groups", ListMethod: "ListGroups", GetMethod: "GetGroupsByID", ImportIDFormat: ""},
	"criblio_group_system_settings":        {SDKService: "System", ListMethod: "", GetMethod: "GetSystemSettingsConf", ImportIDFormat: ""},
	"criblio_hmac_function":                {SDKService: "HmacFunctions", ListMethod: "ListHmacFunction", GetMethod: "GetHmacFunctionByID", ImportIDFormat: ""},
	"criblio_key":                          {SDKService: "", ListMethod: "", GetMethod: "", ImportIDFormat: ""},
	"criblio_lakehouse_dataset_connection": {SDKService: "", ListMethod: "", GetMethod: "", ImportIDFormat: ""},
	"criblio_lookup_file":                  {SDKService: "", ListMethod: "", GetMethod: "", ImportIDFormat: ""},
	"criblio_mapping_ruleset":              {SDKService: "", ListMethod: "", GetMethod: "", ImportIDFormat: ""},
	"criblio_notification":                 {SDKService: "Notifications", ListMethod: "ListNotification", GetMethod: "GetNotificationByID", ImportIDFormat: ""},
	"criblio_notification_target":          {SDKService: "NotificationTargets", ListMethod: "ListNotificationTarget", GetMethod: "GetNotificationTargetByID", ImportIDFormat: "", OneOf: oneOfNotificationTarget},
	"criblio_pack":                         {SDKService: "Packs", ListMethod: "GetPacksByGroup", GetMethod: "GetPacksByID", ImportIDFormat: ""},
	"criblio_pack_breakers":                {SDKService: "Routes", ListMethod: "GetBreakersByPack", GetMethod: "GetBreakersByPackAndID", ImportIDFormat: ""},
	"criblio_pack_destination":             {SDKService: "Outputs", ListMethod: "ListPackOutput", GetMethod: "GetPackOutputByID", ImportIDFormat: "", OneOf: oneOfOutput},
	"criblio_pack_lookups":                 {SDKService: "Routes", ListMethod: "GetSystemLookupsByPack", GetMethod: "GetSystemLookupsByPackAndID", ImportIDFormat: ""},
	"criblio_pack_pipeline":                {SDKService: "Routes", ListMethod: "", GetMethod: "GetPipelinesByPackWithID", ImportIDFormat: ""},
	"criblio_pack_routes":                  {SDKService: "Routes", ListMethod: "", GetMethod: "GetRoutesByPack", ImportIDFormat: ""},
	"criblio_pack_source":                  {SDKService: "Routes", ListMethod: "", GetMethod: "GetSystemInputsByPackAndID", ImportIDFormat: "", OneOf: oneOfInput},
	"criblio_pack_vars":                    {SDKService: "GlobalVariables", ListMethod: "", GetMethod: "GetGlobalVariableLibVarsByPackAndID", ImportIDFormat: ""},
	"criblio_parquet_schema":               {SDKService: "Parquetschemas", ListMethod: "ListSchema", GetMethod: "GetSchemaByID", ImportIDFormat: ""},
	"criblio_parser_lib_entry":             {SDKService: "", ListMethod: "", GetMethod: "", ImportIDFormat: ""},
	"criblio_pipeline":                     {SDKService: "Pipelines", ListMethod: "ListPipeline", GetMethod: "GetPipelineByID", ImportIDFormat: ""},
	"criblio_project":                      {SDKService: "Projects", ListMethod: "ListProject", GetMethod: "GetProjectByID", ImportIDFormat: ""},
	"criblio_regex":                        {SDKService: "Regexes", ListMethod: "ListRegexLibEntry", GetMethod: "GetRegexLibEntryByID", ImportIDFormat: ""},
	"criblio_routes":                       {SDKService: "Routes", ListMethod: "", GetMethod: "GetRoutesByGroupID", ImportIDFormat: ""},
	"criblio_schema":                       {SDKService: "Schemas", ListMethod: "ListLibSchemas", GetMethod: "GetLibSchemasByID", ImportIDFormat: ""},
	"criblio_search_dashboard":             {SDKService: "Dashboards", ListMethod: "ListSearchDashboard", GetMethod: "GetSearchDashboardByID", ImportIDFormat: ""},
	"criblio_search_dashboard_category":    {SDKService: "DashboardCategories", ListMethod: "ListDashboardCategory", GetMethod: "GetDashboardCategoryByID", ImportIDFormat: ""},
	"criblio_search_dataset":               {SDKService: "Datasets", ListMethod: "ListDataset", GetMethod: "GetDatasetByID", ImportIDFormat: "", OneOf: oneOfSearchDataset},
	"criblio_search_dataset_provider":      {SDKService: "Datasets", ListMethod: "ListDatasetProvider", GetMethod: "GetDatasetProviderByID", ImportIDFormat: "", OneOf: oneOfSearchDatasetProvider},
	"criblio_search_macro":                 {SDKService: "Macros", ListMethod: "ListSearchMacro", GetMethod: "GetSearchMacroByID", ImportIDFormat: ""},
	"criblio_search_saved_query":           {SDKService: "SavedQueries", ListMethod: "ListSavedQuery", GetMethod: "GetSavedQueryByID", ImportIDFormat: ""},
	"criblio_search_usage_group":           {SDKService: "", ListMethod: "", GetMethod: "", ImportIDFormat: ""},
	"criblio_secret":                       {SDKService: "Secrets", ListMethod: "ListRestSecret", GetMethod: "GetRestSecretByID", ImportIDFormat: ""},
	"criblio_source":                       {SDKService: "Inputs", ListMethod: "ListInput", GetMethod: "GetInputByID", ImportIDFormat: "", OneOf: oneOfInput},
	"criblio_subscription":                 {SDKService: "Subscriptions", ListMethod: "ListSubscription", GetMethod: "GetSubscriptionByID", ImportIDFormat: ""},
	"criblio_workspace":                    {SDKService: "", ListMethod: "", GetMethod: "", ImportIDFormat: ""},
}

var importMetadataOverrides = map[string]ResourceMetadata{
	"criblio_appscope_config":              {ImportIDFormat: "json:group_id,id"},
	"criblio_collector":                    {ListMethod: "ListCollectors", ImportIDFormat: "json:group_id,id"},
	"criblio_cribl_lake_house":             {ListMethod: "ListDefaultLakeLakehouses", ImportIDFormat: "id"},
	"criblio_global_var":                   {ImportIDFormat: "json:group_id,id"},
	"criblio_grok":                         {ImportIDFormat: "json:group_id,id"},
	"criblio_group":                        {ImportIDFormat: "group_id", RefreshFromMethod: "RefreshFromOperationsCreateProductsGroupsByProductResponseBody"},
	"criblio_group_system_settings":        {ListMethod: "GetSystemSettingsConf", ImportIDFormat: "group_id", ListUseGroupIDAsItemID: true},
	"criblio_key":                          {SDKService: "Keys", ListMethod: "ListKeyMetadataEntity", GetMethod: "GetKeyMetadataEntityByID", ImportIDFormat: "json:group_id,id,key_id", ListItemIDMethod: "GetKeyID", RefreshFromMethod: "RefreshFromOperationsCreateKeyMetadataEntityResponseBody"},
	"criblio_pack":                         {ImportIDFormat: "json:group_id,id,pack"},
	"criblio_pack_source":                  {ListMethod: "GetSystemInputsByPack", GetMethod: "GetSystemInputsByPackAndID", ImportIDFormat: "json:group_id,id,pack"},
	"criblio_pack_pipeline":                {ListMethod: "GetPipelinesByPack", GetMethod: "GetPipelinesByPackWithID", ImportIDFormat: "json:group_id,id,pack"},
	"criblio_pack_breakers":                {ImportIDFormat: "json:group_id,id,pack"},
	"criblio_pack_lookups":                 {ImportIDFormat: "json:group_id,id,pack"},
	"criblio_pack_routes":                  {ImportIDFormat: "json:group_id,pack"},
	"criblio_routes":                       {ListMethod: "GetRoutesByGroupID", GetMethod: "GetRoutesByGroupID", ImportIDFormat: "group_id", ListUseGroupIDAsItemID: true},
	"criblio_parser_lib_entry":             {SDKService: "Parsers", ListMethod: "ListParser", GetMethod: "GetParserByID", ImportIDFormat: "json:group_id,id", RefreshFromMethod: "RefreshFromOperationsGetParserByIDResponseBody"},
	"criblio_search_dashboard":             {GetMethod: "GetSearchDashboardByID", ImportIDFormat: "id"},
	"criblio_notification":                 {ImportIDFormat: "id"},
	"criblio_notification_target":          {ImportIDFormat: "id"},
	"criblio_lookup_file":                  {SDKService: "Lookups", ListMethod: "ListLookupFile", GetMethod: "GetLookupFileByID", ImportIDFormat: "json:group_id,id", ListUseGroupIDAsItemID: false},
	"criblio_search_dashboard_category":    {ListMethod: "ListDashboardCategory", ImportIDFormat: "id"},
	"criblio_parquet_schema":               {ImportIDFormat: "json:group_id,id"},
	"criblio_regex":                        {ImportIDFormat: "json:group_id,id"},
	"criblio_schema":                       {ImportIDFormat: "json:group_id,id"},
	"criblio_search_macro":                 {ImportIDFormat: "id"},
	"criblio_search_saved_query":           {ListMethod: "ListSavedQuery", ImportIDFormat: "id"},
	"criblio_pack_destination":             {ListMethod: "ListPackOutput", ImportIDFormat: "json:group_id,id,pack", RefreshFromMethod: "RefreshFromOperationsGetOutputByIDResponseBody"},
	"criblio_pack_vars":                    {ListMethod: "GetGlobalVariableLibVarsByPack", GetMethod: "GetGlobalVariableLibVarsByPackAndID", ImportIDFormat: "json:group_id,id,pack"},
	"criblio_source":                       {ImportIDFormat: "json:group_id,id"},
	"criblio_destination":                  {ImportIDFormat: "json:group_id,id"},
	"criblio_pipeline":                     {ImportIDFormat: "json:group_id,id"},
	"criblio_certificate":                  {ImportIDFormat: "json:group_id,id"},
	"criblio_cribl_lake_dataset":           {ImportIDFormat: "json:lake_id,id"},
	"criblio_database_connection":          {ImportIDFormat: "json:group_id,id"},
	"criblio_event_breaker_ruleset":        {ImportIDFormat: "json:group_id,id"},
	"criblio_hmac_function":                {ImportIDFormat: "json:group_id,id"},
	"criblio_secret":                       {ImportIDFormat: "json:group_id,id"},
	"criblio_project":                      {ImportIDFormat: "json:group_id,id"},
	"criblio_subscription":                 {ImportIDFormat: "json:group_id,id"},
	"criblio_search_dataset":               {ImportIDFormat: "id"},
	"criblio_search_dataset_provider":      {ImportIDFormat: "id"},
	"criblio_search_usage_group":           {SDKService: "UsageGroups", ListMethod: "ListUsageGroup", GetMethod: "GetUsageGroupByID", ImportIDFormat: "id", RefreshFromMethod: "RefreshFromOperationsGetUsageGroupByIDResponseBody"},
	"criblio_lakehouse_dataset_connection": {ImportIDFormat: "json:lakehouse_id,lake_dataset_id"},
}

// importMetadataClearList: types with no working list API; ListMethod is cleared so discovery skips them.
// criblio_group: list via GetProductsGroupsByProduct (stream/edge); special-case in export.
var importMetadataClearList = []string{
	"criblio_group",
}

// ImportMetadata returns merged import metadata for all resource types.
func ImportMetadata() map[string]ResourceMetadata {
	out := make(map[string]ResourceMetadata, len(importMetadataBase))
	for k, v := range importMetadataBase {
		out[k] = v
	}
	for k, o := range importMetadataOverrides {
		m := out[k]
		if o.SDKService != "" {
			m.SDKService = o.SDKService
		}
		if o.ListMethod != "" {
			m.ListMethod = o.ListMethod
		}
		if o.GetMethod != "" {
			m.GetMethod = o.GetMethod
		}
		if o.ImportIDFormat != "" {
			m.ImportIDFormat = o.ImportIDFormat
		}
		if o.RefreshFromMethod != "" {
			m.RefreshFromMethod = o.RefreshFromMethod
		}
		if o.ListItemIDMethod != "" {
			m.ListItemIDMethod = o.ListItemIDMethod
		}
		if o.ListUseGroupIDAsItemID {
			m.ListUseGroupIDAsItemID = o.ListUseGroupIDAsItemID
		}
		out[k] = m
	}
	for _, name := range importMetadataClearList {
		m := out[name]
		m.ListMethod = ""
		out[name] = m
	}
	return out
}
