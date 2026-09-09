// Package export converts discovery results into generator ResourceItems by
// listing identifiers, fetching each resource via the converter, and building
// HCL attributes and import IDs.
package export

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"sort"
	"sync"

	"github.com/criblio/terraform-provider-criblio/internal/restclient"
	importclient "github.com/criblio/terraform-provider-criblio/tools/import-cli/internal/client"
	"github.com/criblio/terraform-provider-criblio/tools/import-cli/internal/converter"
	"github.com/criblio/terraform-provider-criblio/tools/import-cli/internal/custom"
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
	DefaultsSkipped int              // count of resources skipped due to --exclude-defaults
	DiscoveredTotal int              // sum of discovery counts for types we attempted (set by caller)
}

type conversionTask struct {
	result discovery.Result
	entry  registry.Entry
	idMap  map[string]string
}

type conversionOutcome struct {
	item             *generator.ResourceItem
	skipMessage      string
	defaultsSkipped  int
	requestAttempted bool
	requestErr       error
}

type conversionRequest struct {
	attempted bool
	err       error
}

// ProgressFunc reports progress to the user; nil means no progress output.
type ProgressFunc func(format string, args ...interface{})

// ToResourceItems turns discovery results into generator ResourceItems for types
// that have GetMethod and ImportIDFormat. Uses groupIDs for list/get requests.
// parallel limits concurrent API calls (default 5); use 1 for sequential.
// excludeDefaults, when true, skips built-in Cribl resources (lib=cribl, tags=cribl:default, known default IDs).
// progress, when non-nil, is called to report progress for each group.
// Continues on list-level and per-item errors so as many resources as possible are exported;
// failed types or items are recorded in result.ListSkipped and result.ConvertSkipped.
// Caller should set result.DiscoveredTotal to the sum of discovery counts for reporting.
// groupFilter is the CLI --group slice; when non-empty, export only resources whose output folder
// matches a resolved worker/search group (see skipExportForGroupFilter).
func ToResourceItems(ctx context.Context, client *importclient.Client, reg *registry.Registry, results []discovery.Result, groupIDs []string, groupFilter []string, parallel int, excludeDefaults bool, includeOverride IncludeOverride, progress ProgressFunc) (result *ExportResult, err error) {
	if parallel < 1 {
		parallel = 1
	}
	out := &ExportResult{}
	tasksByGroup := make(map[string][]conversionTask)
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
		if e.RESTGetPath == "" && r.TypeName != "criblio_lakehouse_dataset_connection" {
			out.ListSkipped = append(out.ListSkipped, ListSkipReason{TypeName: r.TypeName, Reason: "no GetMethod or ImportIDFormat", Count: r.Count})
			continue
		}
		// criblio_group: use the list response because the legacy master GET omits fields on some deployments.
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
				requestParams := toRequestParams(idMap)
				model, convErr := converter.ConvertRawItemWithIdentifiers(e, groupItems[i], requestParams)
				if convErr != nil {
					out.ConvertSkipped = append(out.ConvertSkipped, fmt.Sprintf("%s %v: %s", r.TypeName, idMap, sanitizeConvertError(convErr)))
					continue
				}
				if appendErr := appendResourceItemFromModel(out, r.TypeName, e, idMap, model, groupFilter, groupIDs, excludeDefaults, includeOverride); appendErr != nil {
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
			idMaps := r.Identifiers
			var listErr error
			if !r.InventoryComplete {
				idMaps, listErr = discovery.ListItemIdentifiers(ctx, client, e, groupIDs)
			}
			if listErr != nil {
				out.ListSkipped = append(out.ListSkipped, ListSkipReason{TypeName: r.TypeName, Reason: listErr.Error(), Count: r.Count})
				continue
			}
			if len(idMaps) == 0 {
				out.ListSkipped = append(out.ListSkipped, ListSkipReason{TypeName: r.TypeName, Reason: "list returned 0 identifiers", Count: 0})
				continue
			}
			for _, idMap := range idMaps {
				if skipExportForGroupFilter(r.TypeName, idMap, groupFilter, groupIDs) {
					continue
				}
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
					TypeName: e.TypeName,
					Name:     name,
					Attrs:    attrs,
					ImportID: importID,
					GroupID:  "global", // lakehouse_dataset_connection has no group_id
				})
			}
			continue
		}
		idMaps := r.Identifiers
		var listErr error
		if !r.InventoryComplete {
			idMaps, listErr = discovery.ListItemIdentifiers(ctx, client, e, groupIDs)
		}
		if listErr != nil {
			out.ListSkipped = append(out.ListSkipped, ListSkipReason{TypeName: r.TypeName, Reason: listErr.Error(), Count: r.Count})
			continue
		}
		if len(idMaps) == 0 {
			out.ListSkipped = append(out.ListSkipped, ListSkipReason{TypeName: r.TypeName, Reason: "list returned 0 identifiers", Count: 0})
			continue
		}
		for _, idMap := range idMaps {
			groupID := idMap["group_id"]
			tasksByGroup[groupID] = append(tasksByGroup[groupID], conversionTask{
				result: r,
				entry:  e,
				idMap:  idMap,
			})
		}
	}
	convertTasksByGroup(ctx, client, out, tasksByGroup, groupIDs, groupFilter, parallel, excludeDefaults, includeOverride, progress)
	// ConvertSkipped is informational (oneOf unsupported, skip by config, etc.); do not treat as fatal error.
	// printExportSummary reports skipped resources; export succeeds with best-effort output.

	// Post-process: if excludeDefaults, remove default group resources that have no user-created resources.
	if excludeDefaults {
		filterEmptyDefaultGroups(out)
	}

	return out, nil
}

func convertTasksByGroup(ctx context.Context, client *importclient.Client, out *ExportResult, tasksByGroup map[string][]conversionTask, groupIDs, groupFilter []string, parallel int, excludeDefaults bool, includeOverride IncludeOverride, progress ProgressFunc) {
	groupOrder := orderedTaskGroups(tasksByGroup, groupIDs)
	for _, groupID := range groupOrder {
		tasks := tasksByGroup[groupID]
		if len(tasks) == 0 {
			continue
		}
		if progress != nil {
			label := groupID
			if label == "" {
				label = "global"
			}
			progress("group %s: exporting %d items", label, len(tasks))
		}

		remaining := tasks
		for len(remaining) > 0 {
			first := runConversionTask(ctx, client, remaining[0], groupIDs, groupFilter, excludeDefaults, includeOverride)
			mergeConversionOutcome(out, first)
			remaining = remaining[1:]
			if !first.requestAttempted {
				continue
			}
			if first.requestErr == nil {
				break
			}
			if isTooManyRequests(first.requestErr) {
				for _, task := range remaining {
					out.ConvertSkipped = append(out.ConvertSkipped, fmt.Sprintf("%s %v: group bootstrap failed: %s", task.result.TypeName, task.idMap, sanitizeConvertError(first.requestErr)))
				}
				remaining = nil
				break
			}
		}
		if len(remaining) == 0 {
			continue
		}

		outcomes := make(chan conversionOutcome, len(remaining))
		sem := make(chan struct{}, parallel)
		var wg sync.WaitGroup
		for _, task := range remaining {
			task := task
			wg.Add(1)
			go func() {
				defer wg.Done()
				sem <- struct{}{}
				defer func() { <-sem }()
				outcomes <- runConversionTask(ctx, client, task, groupIDs, groupFilter, excludeDefaults, includeOverride)
			}()
		}
		wg.Wait()
		close(outcomes)
		for outcome := range outcomes {
			mergeConversionOutcome(out, outcome)
		}
	}
}

func runConversionTask(ctx context.Context, client *importclient.Client, task conversionTask, groupIDs, groupFilter []string, excludeDefaults bool, includeOverride IncludeOverride) conversionOutcome {
	local := &ExportResult{}
	request := &conversionRequest{}
	item, skipMessage := convertOneResource(ctx, client, task.result, task.entry, task.idMap, groupFilter, groupIDs, excludeDefaults, includeOverride, local, request)
	return conversionOutcome{
		item:             item,
		skipMessage:      skipMessage,
		defaultsSkipped:  local.DefaultsSkipped,
		requestAttempted: request.attempted,
		requestErr:       request.err,
	}
}

func isTooManyRequests(err error) bool {
	var httpErr *restclient.HTTPError
	return errors.As(err, &httpErr) && httpErr.StatusCode == http.StatusTooManyRequests
}

func mergeConversionOutcome(out *ExportResult, outcome conversionOutcome) {
	out.DefaultsSkipped += outcome.defaultsSkipped
	if outcome.skipMessage != "" {
		out.ConvertSkipped = append(out.ConvertSkipped, outcome.skipMessage)
	} else if outcome.item != nil {
		out.Items = append(out.Items, *outcome.item)
	}
}

func orderedTaskGroups(tasksByGroup map[string][]conversionTask, groupIDs []string) []string {
	order := make([]string, 0, len(tasksByGroup))
	seen := make(map[string]bool, len(tasksByGroup))
	for i := len(groupIDs) - 1; i >= 0; i-- {
		groupID := groupIDs[i]
		if len(tasksByGroup[groupID]) > 0 && !seen[groupID] {
			order = append(order, groupID)
			seen[groupID] = true
		}
	}
	if len(tasksByGroup[""]) > 0 {
		order = append(order, "")
		seen[""] = true
	}
	var remaining []string
	for groupID, tasks := range tasksByGroup {
		if len(tasks) > 0 && !seen[groupID] {
			remaining = append(remaining, groupID)
		}
	}
	sort.Strings(remaining)
	return append(order, remaining...)
}

// filterEmptyDefaultGroups removes criblio_group resources for default groups (default, default_fleet, defaultHybrid)
// if no other user-created resources exist in that group. Modifies out.Items in place.
func filterEmptyDefaultGroups(out *ExportResult) {
	// Count non-group resources per group_id
	resourcesPerGroup := make(map[string]int)
	for _, item := range out.Items {
		if item.TypeName == "criblio_group" {
			continue
		}
		if item.GroupID != "" {
			resourcesPerGroup[item.GroupID]++
		}
	}

	// Filter out default group resources that have no other resources
	filtered := make([]generator.ResourceItem, 0, len(out.Items))
	for _, item := range out.Items {
		if item.TypeName == "criblio_group" && custom.DefaultGroupIDs[item.GroupID] {
			if resourcesPerGroup[item.GroupID] == 0 {
				out.DefaultsSkipped++
				continue
			}
		}
		filtered = append(filtered, item)
	}
	out.Items = filtered
}
