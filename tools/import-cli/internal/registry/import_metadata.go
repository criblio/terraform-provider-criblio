// Package registry provides import metadata for the import CLI: SDK service and
// method names and import ID format per resource type.
//
// When you add a new resource to the provider (internal/provider/provider.go
// Resources()), add a corresponding entry to importMetadataBase below (and
// importMetadataOverrides or importMetadataClearList if needed). A test
// ensures every provider resource type has an entry here.
package registry

// importMetadataBase is the default metadata per type (from resource SDK usage).
// Overrides and clearListMethodTypes are applied in ImportMetadata().
var importMetadataBase = map[string]ResourceMetadata{
	"criblio_appscope_config":              {SDKService: "AppscopeConfigs", ListMethod: "ListAppscopeLibEntry", GetMethod: "GetAppscopeLibEntryByID", ImportIDFormat: ""},
	"criblio_certificate":                  {SDKService: "Certificates", ListMethod: "ListCertificate", GetMethod: "GetCertificateByID", ImportIDFormat: ""},
	"criblio_collector":                    {SDKService: "SavedJobs", ListMethod: "ListSavedJob", GetMethod: "GetSavedJobByID", ImportIDFormat: ""},
	"criblio_commit":                       {SDKService: "", ListMethod: "", GetMethod: "", ImportIDFormat: ""},
	"criblio_cribl_lake_dataset":           {SDKService: "Lake", ListMethod: "", GetMethod: "GetCriblLakeDatasetByLakeIDAndID", ImportIDFormat: ""},
	"criblio_cribl_lake_house":             {SDKService: "LakeHouse", ListMethod: "ListDefaultLakeLakehouse", GetMethod: "GetDefaultLakeLakehouseByID", ImportIDFormat: ""},
	"criblio_database_connection":          {SDKService: "DatabaseConnections", ListMethod: "ListDatabaseConnectionConfig", GetMethod: "GetDatabaseConnectionConfigByID", ImportIDFormat: ""},
	"criblio_deploy":                       {SDKService: "", ListMethod: "", GetMethod: "", ImportIDFormat: ""},
	"criblio_destination":                  {SDKService: "Outputs", ListMethod: "ListOutput", GetMethod: "GetOutputByID", ImportIDFormat: ""},
	"criblio_event_breaker_ruleset":        {SDKService: "EventBreakerRules", ListMethod: "ListEventBreakerRuleset", GetMethod: "GetEventBreakerRulesetByID", ImportIDFormat: ""},
	"criblio_global_var":                   {SDKService: "GlobalVariables", ListMethod: "ListGlobalVariable", GetMethod: "GetGlobalVariableByID", ImportIDFormat: ""},
	"criblio_grok":                         {SDKService: "Grokfiles", ListMethod: "ListGrokFile", GetMethod: "GetGrokFileByID", ImportIDFormat: ""},
	"criblio_group":                        {SDKService: "Groups", ListMethod: "ListGroups", GetMethod: "GetGroupsByID", ImportIDFormat: ""},
	"criblio_group_system_settings":        {SDKService: "System", ListMethod: "", GetMethod: "GetSystemSettingsConf", ImportIDFormat: ""},
	"criblio_hmac_function":                {SDKService: "HmacFunctions", ListMethod: "ListHmacFunction", GetMethod: "GetHmacFunctionByID", ImportIDFormat: ""},
	"criblio_key":                          {SDKService: "", ListMethod: "", GetMethod: "", ImportIDFormat: ""},
	"criblio_lakehouse_dataset_connection": {SDKService: "", ListMethod: "", GetMethod: "", ImportIDFormat: ""},
	"criblio_lookup_file":                  {SDKService: "", ListMethod: "", GetMethod: "", ImportIDFormat: ""},
	"criblio_mapping_ruleset":              {SDKService: "", ListMethod: "", GetMethod: "", ImportIDFormat: ""},
	"criblio_notification":                 {SDKService: "Notifications", ListMethod: "ListNotification", GetMethod: "GetNotificationByID", ImportIDFormat: ""},
	"criblio_notification_target":          {SDKService: "NotificationTargets", ListMethod: "ListNotificationTarget", GetMethod: "GetNotificationTargetByID", ImportIDFormat: ""},
	"criblio_pack":                         {SDKService: "Packs", ListMethod: "ListPacks", GetMethod: "GetPacksByID", ImportIDFormat: ""},
	"criblio_pack_breakers":                {SDKService: "Routes", ListMethod: "", GetMethod: "GetBreakersByPackAndID", ImportIDFormat: ""},
	"criblio_pack_destination":             {SDKService: "Outputs", ListMethod: "ListPackOutput", GetMethod: "GetPackOutputByID", ImportIDFormat: ""},
	"criblio_pack_lookups":                 {SDKService: "Routes", ListMethod: "", GetMethod: "GetSystemLookupsByPackAndID", ImportIDFormat: ""},
	"criblio_pack_pipeline":                {SDKService: "Routes", ListMethod: "", GetMethod: "GetPipelinesByPackWithID", ImportIDFormat: ""},
	"criblio_pack_routes":                  {SDKService: "Routes", ListMethod: "", GetMethod: "GetRoutesByPack", ImportIDFormat: ""},
	"criblio_pack_source":                  {SDKService: "Routes", ListMethod: "", GetMethod: "GetSystemInputsByPackAndID", ImportIDFormat: ""},
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
	"criblio_search_dataset":               {SDKService: "Datasets", ListMethod: "ListDataset", GetMethod: "GetDatasetByID", ImportIDFormat: ""},
	"criblio_search_dataset_provider":      {SDKService: "Datasets", ListMethod: "ListDatasetProvider", GetMethod: "GetDatasetProviderByID", ImportIDFormat: ""},
	"criblio_search_macro":                 {SDKService: "Macros", ListMethod: "ListSearchMacro", GetMethod: "GetSearchMacroByID", ImportIDFormat: ""},
	"criblio_search_saved_query":           {SDKService: "SavedQueries", ListMethod: "ListSavedQuery", GetMethod: "GetSavedQueryByID", ImportIDFormat: ""},
	"criblio_search_usage_group":           {SDKService: "", ListMethod: "", GetMethod: "", ImportIDFormat: ""},
	"criblio_secret":                       {SDKService: "Secrets", ListMethod: "ListRestSecret", GetMethod: "GetRestSecretByID", ImportIDFormat: ""},
	"criblio_source":                       {SDKService: "Inputs", ListMethod: "ListInput", GetMethod: "GetInputByID", ImportIDFormat: ""},
	"criblio_subscription":                 {SDKService: "Subscriptions", ListMethod: "ListSubscription", GetMethod: "GetSubscriptionByID", ImportIDFormat: ""},
	"criblio_workspace":                    {SDKService: "", ListMethod: "", GetMethod: "", ImportIDFormat: ""},
}

var importMetadataOverrides = map[string]ResourceMetadata{
	"criblio_collector":                 {ListMethod: "ListCollectors"},
	"criblio_cribl_lake_house":          {ListMethod: "ListDefaultLakeLakehouses"},
	"criblio_global_var":                {ImportIDFormat: "json:group_id,id"},
	"criblio_group_system_settings":     {ListMethod: "GetSystemSettingsConf", ImportIDFormat: "group_id"},
	"criblio_pack":                      {ImportIDFormat: "json:group_id,id,pack"},
	"criblio_pack_source":               {ListMethod: "GetSystemInputsByPack", GetMethod: "GetSystemInputsByPackAndID", ImportIDFormat: "json:group_id,id,pack"},
	"criblio_pack_pipeline":             {ListMethod: "GetPipelinesByPack", GetMethod: "GetPipelinesByPackWithID", ImportIDFormat: "json:group_id,id,pack"},
	"criblio_routes":                    {GetMethod: "GetRoutesByID", ImportIDFormat: "json:group_id,id"},
	"criblio_parser_lib_entry":          {SDKService: "Parsers", ListMethod: "ListParser", GetMethod: "GetParserByID", ImportIDFormat: "json:group_id,id"},
	"criblio_search_dashboard":          {GetMethod: "GetSearchDashboardByID", ImportIDFormat: "id"},
	"criblio_notification":              {ImportIDFormat: "id"},
	"criblio_notification_target":       {ImportIDFormat: "id"},
	"criblio_lookup_file":               {SDKService: "Lookups", ListMethod: "ListLookupFile", GetMethod: "GetLookupFileByID", ImportIDFormat: "group_id"},
	"criblio_search_dashboard_category": {ListMethod: "ListDashboardCategory", ImportIDFormat: "id"},
	"criblio_search_macro":              {ImportIDFormat: "id"},
	"criblio_search_saved_query":        {ListMethod: "ListSavedQuery", ImportIDFormat: "id"},
	"criblio_key":                       {SDKService: "Keys", ListMethod: "ListKeyMetadataEntity", GetMethod: "GetKeyMetadataEntityByID", ImportIDFormat: "json:group_id,id"},
	"criblio_pack_destination":          {ListMethod: "ListPackOutput", ImportIDFormat: "json:group_id,id,pack"},
	"criblio_source":                    {ImportIDFormat: "json:group_id,id"},
	"criblio_destination":               {ImportIDFormat: "json:group_id,id"},
	"criblio_pipeline":                  {ImportIDFormat: "json:group_id,id"},
	"criblio_certificate":               {ImportIDFormat: "json:group_id,id"},
}

// importMetadataClearList: types with no working list API; ListMethod is cleared so discovery skips them.
var importMetadataClearList = []string{
	"criblio_database_connection",
	"criblio_group",
	"criblio_routes",
	"criblio_search_dataset",
	"criblio_search_dataset_provider",
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
		out[k] = m
	}
	for _, name := range importMetadataClearList {
		m := out[name]
		m.ListMethod = ""
		out[name] = m
	}
	return out
}
