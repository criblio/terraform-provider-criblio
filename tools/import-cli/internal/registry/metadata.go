// Package registry: SDK method and import ID metadata for the import CLI.
//
// ResourceMetadata is filled from ImportMetadata() in import_metadata.go.
package registry

// OneOfConfig describes a resource whose config uses a oneOf pattern: the API returns
// a single payload in a read-only list (e.g. Items), and the Terraform schema uses
// type-specific blocks keyed by a discriminator (e.g. output_prometheus, output_webhook).
// The import CLI uses this to emit the correct block from the first list item.
type OneOfConfig struct {
	// ReadOnlyAttr is the model attribute that holds the API payload list (e.g. "Items").
	// This attribute is skipped when writing HCL (read-only/computed).
	ReadOnlyAttr string
	// DiscriminatorField is the key in each list item that identifies the variant (e.g. "type").
	DiscriminatorField string
	// BlockNamePrefix is prepended to the normalized discriminator value to form the block name
	// (e.g. "output_" -> output_prometheus, "input_collector_" -> input_collector_rest).
	BlockNamePrefix string
	// BlockNameSuffix is appended to the normalized discriminator value when prefix is empty.
	// Use when provider blocks are {type}_target (e.g. smtp_target, slack_target for notification_target).
	BlockNameSuffix string
	// KeysToSkip are item keys to omit from the emitted block (e.g. "status").
	KeysToSkip []string
	// DiscriminatorAlias maps API discriminator values to provider block names (without prefix).
	// Use when the API returns a type that is not a valid provider block (e.g. "collection" -> "rest" for collector).
	DiscriminatorAlias map[string]string
	// UnsupportedDiscriminatorValues lists API discriminator values that the provider cannot manage; resources with these types are skipped during export.
	UnsupportedDiscriminatorValues []string
	// NestedDiscriminatorField specifies a dot-separated path to a nested discriminator
	// (e.g. "collector.type"). When set, after checking UnsupportedDiscriminatorValues
	// the exporter parses the parent object from the item map and reads the inner field
	// to resolve the actual block type. Used for SavedJob/collector where the top-level
	// "type" is the job kind ("collection") and the real collector type is at collector.type.
	NestedDiscriminatorField string
	// SupportedBlockNames is the list of provider block names (e.g. smtp_target, slack_target) for this oneOf. When set, API discriminator is resolved to one of these (with heuristics) and resources that don't match are skipped. Populated dynamically from the provider model when nil in metadata.
	SupportedBlockNames []string
}

// ResourceMetadata holds SDK service/method names and import ID format for a resource type.
type ResourceMetadata struct {
	SDKService     string // e.g. "Inputs", "Pipelines"
	ListMethod     string // e.g. "ListInput"; empty if no working list API
	GetMethod      string // e.g. "GetInputByID"
	ImportIDFormat string // e.g. "json:group_id,id", "id"
	// OneOf configures oneOf-style resources (destination, collector, pack_destination, etc.).
	OneOf *OneOfConfig
	// RefreshFromMethod overrides the RefreshFrom* method name used when converting Get response to model.
	// When set, the converter uses this instead of RefreshFromOperations+GetMethod+"ResponseBody".
	// Use when the provider has no RefreshFrom* for the Get response but has one for a compatible type (e.g. ListParserResponseBody for parser lib entry).
	RefreshFromMethod string
	// ListItemIDMethod is the method name on list items to get the ID (e.g. "GetKeyID"). Empty = use GetID or map "id"/"Id".
	ListItemIDMethod string
	// ListUseGroupIDAsItemID when true: when list is called per-group and the item has no id, use the group ID as the identifier (e.g. criblio_group_system_settings).
	ListUseGroupIDAsItemID bool
}
