// Package export converts discovery results into generator ResourceItems by
// listing identifiers, fetching each resource via the converter, and building
// HCL attributes and import IDs.
package export

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/criblio/terraform-provider-criblio/internal/sdk"
	"github.com/criblio/terraform-provider-criblio/internal/sdk/models/operations"
	"github.com/criblio/terraform-provider-criblio/internal/sdk/models/shared"
	"github.com/criblio/terraform-provider-criblio/tools/import-cli/internal/converter"
	"github.com/criblio/terraform-provider-criblio/tools/import-cli/internal/discovery"
	"github.com/criblio/terraform-provider-criblio/tools/import-cli/internal/generator"
	"github.com/criblio/terraform-provider-criblio/tools/import-cli/internal/hcl"
	"github.com/criblio/terraform-provider-criblio/tools/import-cli/internal/registry"
)

// ErrUnsupportedOneOfType is returned when a oneOf resource has a discriminator value the provider does not support; the exporter skips the resource.
var ErrUnsupportedOneOfType = errors.New("unsupported oneOf type")

// ErrSkipResourceLibCribl is returned when a resource has lib = "cribl" (built-in/system); the exporter skips it.
var ErrSkipResourceLibCribl = errors.New("lib is cribl (built-in, skip export)")

// ListSkipReason describes why a resource type produced no items at list stage.
type ListSkipReason struct {
	TypeName string
	Reason   string
	Count    int // number of resources skipped (from discovery) for this type
}

// ExportResult holds the result of ToResourceItems for reporting.
type ExportResult struct {
	Items           []generator.ResourceItem
	ListSkipped     []ListSkipReason // types skipped at list (no metadata, list failed, or list returned 0 ids)
	ConvertSkipped  []string         // one message per resource that failed convert/hcl/import
	DiscoveredTotal int              // sum of discovery counts for types we attempted (set by caller)
}

// ProgressFunc reports progress to the user; nil means no progress output.
type ProgressFunc func(format string, args ...interface{})

// ToResourceItems turns discovery results into generator ResourceItems for types
// that have GetMethod and ImportIDFormat. Uses groupIDs for list/get requests.
// parallel limits concurrent API calls (default 5); use 1 for sequential.
// progress, when non-nil, is called to report progress per resource type.
// Continues on list-level and per-item errors so as many resources as possible are exported;
// failed types or items are recorded in result.ListSkipped and result.ConvertSkipped.
// Caller should set result.DiscoveredTotal to the sum of discovery counts for reporting.
func ToResourceItems(ctx context.Context, client *sdk.CriblIo, reg *registry.Registry, results []discovery.Result, groupIDs []string, parallel int, progress ProgressFunc) (result *ExportResult, err error) {
	if parallel < 1 {
		parallel = 1
	}
	out := &ExportResult{}
	for _, r := range results {
		if r.Err != nil {
			out.ListSkipped = append(out.ListSkipped, ListSkipReason{TypeName: r.TypeName, Reason: r.Err.Error(), Count: r.Count})
			continue
		}
		if r.Count == 0 {
			out.ListSkipped = append(out.ListSkipped, ListSkipReason{TypeName: r.TypeName, Reason: "list returned 0 identifiers", Count: 0})
			continue
		}
		e, ok := reg.ByTypeName(r.TypeName)
		if !ok || e.ImportIDFormat == "" {
			out.ListSkipped = append(out.ListSkipped, ListSkipReason{TypeName: r.TypeName, Reason: "no GetMethod or ImportIDFormat", Count: r.Count})
			continue
		}
		if e.GetMethod == "" && r.TypeName != "criblio_lakehouse_dataset_connection" {
			out.ListSkipped = append(out.ListSkipped, ListSkipReason{TypeName: r.TypeName, Reason: "no GetMethod or ImportIDFormat", Count: r.Count})
			continue
		}
		// criblio_group: SDK GetGroupsByID response body is empty; use list response (GetProductsGroupsByProduct) and refresh from CreateProductsGroupsByProductResponseBody.
		if r.TypeName == "criblio_group" {
			if progress != nil {
				progress("criblio_group: %d items", r.Count)
			}
			idMaps, groupItems, listErr := discovery.ListGroupIdentifiersAndItems(ctx, client, groupIDs)
			if listErr != nil {
				out.ListSkipped = append(out.ListSkipped, ListSkipReason{TypeName: r.TypeName, Reason: listErr.Error(), Count: r.Count})
				continue
			}
			if len(idMaps) == 0 {
				out.ListSkipped = append(out.ListSkipped, ListSkipReason{TypeName: r.TypeName, Reason: "list returned 0 identifiers", Count: 0})
				continue
			}
			for i, idMap := range idMaps {
				body := &operations.CreateProductsGroupsByProductResponseBody{Items: []shared.ConfigGroup{groupItems[i]}}
				requestParams := toRequestParams(idMap)
				model, convErr := converter.ConvertFromResponseBodyWithIdentifiers(ctx, e, body, requestParams)
				if convErr != nil {
					out.ConvertSkipped = append(out.ConvertSkipped, fmt.Sprintf("%s %v: %s", r.TypeName, idMap, sanitizeConvertError(convErr)))
					continue
				}
				if appendErr := appendResourceItemFromModel(out, r.TypeName, e, idMap, model); appendErr != nil {
					if errors.Is(appendErr, ErrSkipResourceLibCribl) {
						out.ConvertSkipped = append(out.ConvertSkipped, fmt.Sprintf("%s %v: lib is cribl (built-in, skip export)", r.TypeName, idMap))
					} else {
						out.ConvertSkipped = append(out.ConvertSkipped, fmt.Sprintf("%s %v: %s", r.TypeName, idMap, sanitizeConvertError(appendErr)))
					}
				}
			}
			continue
		}
		// criblio_lakehouse_dataset_connection: no Get API; build minimal HCL from identifiers only.
		if r.TypeName == "criblio_lakehouse_dataset_connection" && e.GetMethod == "" {
			if progress != nil {
				progress("criblio_lakehouse_dataset_connection: %d items", r.Count)
			}
			idMaps, listErr := discovery.ListItemIdentifiers(ctx, client, e, groupIDs)
			if listErr != nil {
				out.ListSkipped = append(out.ListSkipped, ListSkipReason{TypeName: r.TypeName, Reason: listErr.Error(), Count: r.Count})
				continue
			}
			if len(idMaps) == 0 {
				out.ListSkipped = append(out.ListSkipped, ListSkipReason{TypeName: r.TypeName, Reason: "list returned 0 identifiers", Count: 0})
				continue
			}
			for _, idMap := range idMaps {
				importID, idErr := generator.BuildImportID(e.ImportIDFormat, idMap)
				if idErr != nil {
					out.ConvertSkipped = append(out.ConvertSkipped, fmt.Sprintf("%s %v: import ID: %s", r.TypeName, idMap, sanitizeConvertError(idErr)))
					continue
				}
				attrs := map[string]hcl.Value{
					"lakehouse_id":    {Kind: hcl.KindString, String: idMap["lakehouse_id"]},
					"lake_dataset_id": {Kind: hcl.KindString, String: idMap["lake_dataset_id"]},
				}
				name := generator.StableResourceNameFromMap(e.TypeName, idMap)
				out.Items = append(out.Items, generator.ResourceItem{
					TypeName:  e.TypeName,
					Name:      name,
					Attrs:     attrs,
					ImportID:  importID,
					GroupID:   "global", // lakehouse_dataset_connection has no group_id
				})
			}
			continue
		}
		idMaps, listErr := discovery.ListItemIdentifiers(ctx, client, e, groupIDs)
		if listErr != nil {
			out.ListSkipped = append(out.ListSkipped, ListSkipReason{TypeName: r.TypeName, Reason: listErr.Error(), Count: r.Count})
			continue
		}
		if len(idMaps) == 0 {
			out.ListSkipped = append(out.ListSkipped, ListSkipReason{TypeName: r.TypeName, Reason: "list returned 0 identifiers", Count: 0})
			continue
		}
		if progress != nil {
			progress("%s: %d items", r.TypeName, len(idMaps))
		}
		if parallel <= 1 {
			for _, idMap := range idMaps {
				item, skipMsg := convertOneResource(ctx, client, r, e, idMap)
				if skipMsg != "" {
					out.ConvertSkipped = append(out.ConvertSkipped, skipMsg)
				} else if item != nil {
					out.Items = append(out.Items, *item)
				}
			}
		} else {
			sem := make(chan struct{}, parallel)
			var mu sync.Mutex
			var wg sync.WaitGroup
			for _, idMap := range idMaps {
				idMap := idMap
				wg.Add(1)
				go func() {
					defer wg.Done()
					sem <- struct{}{}
					defer func() { <-sem }()
					item, skipMsg := convertOneResource(ctx, client, r, e, idMap)
					mu.Lock()
					if skipMsg != "" {
						out.ConvertSkipped = append(out.ConvertSkipped, skipMsg)
					} else if item != nil {
						out.Items = append(out.Items, *item)
					}
					mu.Unlock()
				}()
			}
			wg.Wait()
		}
	}
	// ConvertSkipped is informational (oneOf unsupported, skip by config, etc.); do not treat as fatal error.
	// printExportSummary reports skipped resources; export succeeds with best-effort output.
	return out, nil
}
