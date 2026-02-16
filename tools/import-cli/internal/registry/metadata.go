// Package registry: SDK method and import ID metadata for the import CLI.
//
// ResourceMetadata is filled from ImportMetadata() in import_metadata.go.
package registry

// ResourceMetadata holds SDK service/method names and import ID format for a resource type.
type ResourceMetadata struct {
	SDKService     string // e.g. "Inputs", "Pipelines"
	ListMethod     string // e.g. "ListInput"; empty if no working list API
	GetMethod      string // e.g. "GetInputByID"
	ImportIDFormat string // e.g. "json:group_id,id", "id"
}
