// Package discovery implements the discovery engine that iterates through the
// registry and calls SDK List* endpoints to enumerate existing resources.
package discovery

import (
	"context"
	"fmt"
	"os"
	"reflect"
	"strings"

	"github.com/criblio/terraform-provider-criblio/internal/sdk"
	"github.com/criblio/terraform-provider-criblio/internal/sdk/models/operations"
	"github.com/criblio/terraform-provider-criblio/internal/sdk/models/shared"
	"github.com/criblio/terraform-provider-criblio/tools/import-cli/internal/registry"
	"github.com/criblio/terraform-provider-criblio/tools/import-cli/internal/custom"
)

// eventBreakerRulesetIdentifiersFromCapture returns identifiers from the captured
// event breaker list response when the SDK failed to unmarshal (lib="cribl" not in enum).
func eventBreakerRulesetIdentifiersFromCapture(gid string) ([]map[string]string, error) {
	body := custom.GetAndClearEventBreakerRulesetListBody(gid)
	if len(body) == 0 {
		return nil, nil
	}
	return custom.ParseEventBreakerRulesetListBody(body, gid)
}

// searchListIdentifiersFromCapture returns identifiers (and count) from the captured
// list response body when the SDK failed to unmarshal (e.g. cribl_lake). Only for
// criblio_search_dataset and criblio_search_dataset_provider.
func searchListIdentifiersFromCapture(e registry.Entry) ([]map[string]string, int, error) {
	var key string
	switch e.TypeName {
	case "criblio_search_dataset":
		key = custom.PathSearchDatasets
	case "criblio_search_dataset_provider":
		key = custom.PathSearchDatasetProviders
	default:
		return nil, 0, nil
	}
	body := custom.GetAndClearSearchListBody(key)
	return custom.IdentifiersFromSearchListBody(body, key)
}

// unionUnmarshalIdentifiersFromCapture returns identifiers from the captured raw list response
// for resource types where the SDK union unmarshal fails for some items (e.g. scheduledSearch in
// InputCollector, bulletin_message in NotificationTarget). Unsupported item types are filtered out.
func unionUnmarshalIdentifiersFromCapture(e registry.Entry, groupID string) ([]map[string]string, int, error) {
	switch e.TypeName {
	case "criblio_collector":
		body := custom.GetAndClearSavedJobsListBody(groupID)
		if len(body) == 0 {
			return nil, 0, nil
		}
		ids, err := custom.ParseSavedJobsListBody(body, groupID)
		return ids, len(ids), err
	case "criblio_notification_target":
		body := custom.GetAndClearNotificationTargetListBody()
		if len(body) == 0 {
			return nil, 0, nil
		}
		ids, err := custom.ParseNotificationTargetListBody(body)
		return ids, len(ids), err
	default:
		return nil, 0, nil
	}
}

// isSDKUnionUnmarshalError reports whether err is the SDK failing to unmarshal a response
// into a oneOf union type (GenericDataset, GenericProvider, InputCollector, NotificationTarget).
// When true, discovery falls back to parsing the captured raw response body instead of failing.
func isSDKUnionUnmarshalError(err error) bool {
	if err == nil {
		return false
	}
	s := err.Error()
	if !strings.Contains(s, "could not unmarshal") {
		return false
	}
	return strings.Contains(s, "GenericDataset") ||
		strings.Contains(s, "GenericProvider") ||
		strings.Contains(s, "InputCollector") ||
		strings.Contains(s, "NotificationTarget")
}

// isSDKLibraryUnmarshalError reports whether err is the SDK failing to unmarshal lib="cribl"
// (EventBreakerRuleset Library enum only has custom/cribl-custom; API returns cribl for built-ins).
// When true, discovery falls back to parsing the captured raw response body instead of failing.
func isSDKLibraryUnmarshalError(err error) bool {
	if err == nil {
		return false
	}
	return strings.Contains(err.Error(), "invalid value for Library")
}

// Result holds the discovery result for one resource type: count and any error
// with resource context. Items are left for future HCL generation.
// Details is optional (e.g. group names for criblio_group) for dry-run display.
// PerGroupCounts is set for group-scoped resources (key = group label e.g. "default (stream)").
type Result struct {
	TypeName       string
	Count          int
	Err            error
	Details        []string       // optional; e.g. group names for criblio_group
	PerGroupCounts map[string]int // optional; per-group count for preview/export
}

// Discover runs the discovery engine: for each registry entry that has
// SDKService and ListMethod (and passes include/exclude), calls the SDK List*
// endpoint and records count and errors. Uses group IDs from the API when
// available so list requests (e.g. /m/{groupId}/...) succeed; falls back to
// ["default"] if the groups API is unavailable.
// groupFilter restricts to specific groups (by ID or label e.g. "default (stream)"); empty = all groups.
func Discover(ctx context.Context, client *sdk.CriblIo, reg *registry.Registry, include, exclude, groupFilter []string) ([]Result, error) {
	includeSet := sliceToSet(include)
	excludeSet := sliceToSet(exclude)

	streamIDs, streamNames, streamErr := fetchGroupsByProduct(ctx, client, operations.GetProductsGroupsByProductProductStream)
	if streamErr != nil {
		return nil, fmt.Errorf("fetch stream groups: %w", streamErr)
	}
	edgeIDs, edgeNames, edgeErr := fetchGroupsByProduct(ctx, client, operations.GetProductsGroupsByProductProductEdge)
	if edgeErr != nil {
		return nil, fmt.Errorf("fetch edge groups: %w", edgeErr)
	}

	// Apply group filter: keep only groups whose ID or label matches.
	streamIDs, streamNames = filterGroups(streamIDs, streamNames, " (stream)", groupFilter)
	edgeIDs, edgeNames = filterGroups(edgeIDs, edgeNames, " (edge)", groupFilter)

	// Use both stream and edge group IDs for list API (edge fleets have their own sources, pipelines, etc.).
	groupIDs := make([]string, 0, len(streamIDs)+len(edgeIDs)+1)
	groupIDs = append(groupIDs, streamIDs...)
	groupIDs = append(groupIDs, edgeIDs...)
	if len(groupIDs) == 0 && len(groupFilter) == 0 {
		groupIDs = fallbackGroupIDs()
	} else {
		// Always include "default" and "default_search" so resources under those groups (e.g. certificates, lookups, parser_lib_entry) are discovered.
		groupIDs = ensureDefaultGroups(groupIDs)
	}
	// Map group ID -> label for PerGroupCounts (stream and edge).
	idToLabel := make(map[string]string)
	for i := range streamIDs {
		idToLabel[streamIDs[i]] = streamNames[i] + " (stream)"
	}
	for i := range edgeIDs {
		idToLabel[edgeIDs[i]] = edgeNames[i] + " (edge)"
	}
	// All group names for criblio_group display (stream and edge, labeled).
	var groupNames []string
	for _, n := range streamNames {
		groupNames = append(groupNames, n+" (stream)")
	}
	for _, n := range edgeNames {
		groupNames = append(groupNames, n+" (edge)")
	}
	if len(streamNames) == 0 && len(edgeNames) == 0 {
		groupNames = nil
	}

	var results []Result
	for _, e := range reg.Entries() {
		if !matchesFilter(e.TypeName, includeSet, excludeSet) {
			continue
		}
		// No list method: show in preview. criblio_group uses group count and names from API.
		// criblio_lakehouse_dataset_connection has no list API; discover via lakehouses × lake datasets.
		if e.SDKService == "" || e.ListMethod == "" {
			count := 0
			var details []string
			if e.TypeName == "criblio_group" {
				count = len(streamNames) + len(edgeNames)
				details = groupNames
			} else if e.TypeName == "criblio_lakehouse_dataset_connection" {
				ids, lhErr := listLakehouseDatasetConnectionIdentifiers(ctx, client)
				if lhErr == nil {
					count = len(ids)
				}
			} else if e.TypeName == "criblio_pack_routes" {
				ids, prErr := listPackRoutesIdentifiers(ctx, client, groupIDs)
				if prErr == nil {
					count = len(ids)
				}
			}
			results = append(results, Result{TypeName: e.TypeName, Count: count, Details: details})
			continue
		}
		count, perGroup, err := listOne(ctx, client, e, groupIDs)
		res := Result{TypeName: e.TypeName, Count: count}
		if err != nil {
			res.Err = fmt.Errorf("%s: %w", e.TypeName, err)
		}
		if len(perGroup) > 0 {
			res.PerGroupCounts = make(map[string]int)
			for gid, n := range perGroup {
				if label, ok := idToLabel[gid]; ok {
					res.PerGroupCounts[label] = n
				} else {
					res.PerGroupCounts[gid] = n
				}
			}
		}
		results = append(results, res)
	}
	return results, nil
}

// filterGroups keeps only (ids[i], names[i]) where ids[i] or (names[i]+suffix) is in filter. If filter is empty, returns all.
func filterGroups(ids, names []string, suffix string, filter []string) (filteredIDs, filteredNames []string) {
	if len(filter) == 0 {
		return ids, names
	}
	filterSet := sliceToSet(filter)
	for i := range ids {
		id := ids[i]
		label := names[i] + suffix
		if inSet(id, filterSet) || inSet(label, filterSet) {
			filteredIDs = append(filteredIDs, id)
			filteredNames = append(filteredNames, names[i])
		}
	}
	return filteredIDs, filteredNames
}

// fetchGroupsByProduct returns group IDs and display names for a given product (stream or edge).
// Names use GetName() with ID as fallback. Used for both stream (listOne) and stream+edge (criblio_group display).
func fetchGroupsByProduct(ctx context.Context, client *sdk.CriblIo, product operations.GetProductsGroupsByProductProduct) (ids []string, names []string, err error) {
	if client == nil || client.Groups == nil {
		return nil, nil, fmt.Errorf("client or Groups service nil")
	}
	req := operations.GetProductsGroupsByProductRequest{Product: product}
	resp, err := client.Groups.GetProductsGroupsByProduct(ctx, req)
	if err != nil {
		return nil, nil, err
	}
	if resp == nil || resp.Object == nil {
		return nil, nil, nil
	}
	items := resp.Object.GetItems()
	if len(items) == 0 {
		return nil, nil, nil
	}
	ids = make([]string, 0, len(items))
	names = make([]string, 0, len(items))
	for _, g := range items {
		id := g.GetID()
		if id == "" {
			continue
		}
		ids = append(ids, id)
		if n := g.GetName(); n != nil && *n != "" {
			names = append(names, *n)
		} else {
			names = append(names, id)
		}
	}
	return ids, names, nil
}

// listGroupIdentifiers returns one identifier map per group by calling GetProductsGroupsByProduct for stream and edge.
// If groupIDs is non-empty, only groups whose id is in groupIDs are returned.
func listGroupIdentifiers(ctx context.Context, client *sdk.CriblIo, groupIDs []string) ([]map[string]string, error) {
	idMaps, _, err := ListGroupIdentifiersAndItems(ctx, client, groupIDs)
	return idMaps, err
}

// listLakehouseDatasetConnectionIdentifiers returns (lakehouse_id, lake_dataset_id) pairs by listing
// lakehouses and lake datasets (default lake). The API has no list endpoint for connections, so we
// derive potential connection resources from the cartesian product.
func listLakehouseDatasetConnectionIdentifiers(ctx context.Context, client *sdk.CriblIo) ([]map[string]string, error) {
	lakehousesResp, err := client.LakeHouse.ListDefaultLakeLakehouses(ctx)
	if err != nil || lakehousesResp == nil || lakehousesResp.Object == nil {
		return nil, err
	}
	datasetsResp, err := client.Lake.GetCriblLakeDatasetByLakeID(ctx, operations.GetCriblLakeDatasetByLakeIDRequest{
		LakeID: operations.GetCriblLakeDatasetByLakeIDLakeIDDefault,
	})
	if err != nil || datasetsResp == nil || datasetsResp.Object == nil {
		return nil, err
	}
	lakehouses := lakehousesResp.Object.GetItems()
	datasets := datasetsResp.Object.GetItems()
	out := make([]map[string]string, 0, len(lakehouses)*len(datasets))
	for _, lh := range lakehouses {
		lakehouseID := lh.GetID()
		if lakehouseID == "" {
			continue
		}
		for _, ds := range datasets {
			datasetID := ds.GetID()
			if datasetID == "" {
				continue
			}
			out = append(out, map[string]string{"lakehouse_id": lakehouseID, "lake_dataset_id": datasetID})
		}
	}
	return out, nil
}

// ListGroupIdentifiersAndItems returns identifier maps and full ConfigGroup items for each group (same order).
// Used by export for criblio_group so we can refresh from list response instead of GetGroupsByID (whose response body is empty in SDK).
func ListGroupIdentifiersAndItems(ctx context.Context, client *sdk.CriblIo, groupIDs []string) (idMaps []map[string]string, items []shared.ConfigGroup, err error) {
	streamIDs, streamItems, err := fetchGroupsByProductWithItems(ctx, client, operations.GetProductsGroupsByProductProductStream)
	if err != nil {
		return nil, nil, err
	}
	edgeIDs, edgeItems, err := fetchGroupsByProductWithItems(ctx, client, operations.GetProductsGroupsByProductProductEdge)
	if err != nil {
		return nil, nil, err
	}
	groupIDSet := sliceToSet(groupIDs)
	for i, id := range streamIDs {
		if len(groupIDSet) > 0 && !inSet(id, groupIDSet) {
			continue
		}
		idMaps = append(idMaps, map[string]string{"group_id": id, "product": "stream"})
		items = append(items, streamItems[i])
	}
	for i, id := range edgeIDs {
		if len(groupIDSet) > 0 && !inSet(id, groupIDSet) {
			continue
		}
		idMaps = append(idMaps, map[string]string{"group_id": id, "product": "edge"})
		items = append(items, edgeItems[i])
	}
	return idMaps, items, nil
}

// fetchGroupsByProductWithItems returns group IDs and full ConfigGroup items for a product (stream or edge).
func fetchGroupsByProductWithItems(ctx context.Context, client *sdk.CriblIo, product operations.GetProductsGroupsByProductProduct) (ids []string, items []shared.ConfigGroup, err error) {
	if client == nil || client.Groups == nil {
		return nil, nil, fmt.Errorf("client or Groups service nil")
	}
	req := operations.GetProductsGroupsByProductRequest{Product: product}
	resp, err := client.Groups.GetProductsGroupsByProduct(ctx, req)
	if err != nil {
		return nil, nil, err
	}
	if resp == nil || resp.Object == nil {
		return nil, nil, nil
	}
	groupItems := resp.Object.GetItems()
	if len(groupItems) == 0 {
		return nil, nil, nil
	}
	ids = make([]string, 0, len(groupItems))
	items = make([]shared.ConfigGroup, 0, len(groupItems))
	for _, g := range groupItems {
		id := g.GetID()
		if id == "" {
			continue
		}
		ids = append(ids, id)
		items = append(items, g)
	}
	return ids, items, nil
}

func sliceToSet(s []string) map[string]struct{} {
	m := make(map[string]struct{}, len(s))
	for _, v := range s {
		m[v] = struct{}{}
	}
	return m
}

func matchesFilter(typeName string, include, exclude map[string]struct{}) bool {
	if len(exclude) > 0 && inSet(typeName, exclude) {
		return false
	}
	if len(include) > 0 && !inSet(typeName, include) {
		return false
	}
	return true
}

func inSet(s string, m map[string]struct{}) bool {
	_, ok := m[s]
	return ok
}

// GetGroupIDs returns group IDs used for list/export (stream + edge, filtered by groupFilter).
// Same logic as inside Discover. Use when exporting so list and get calls use the same groups.
func GetGroupIDs(ctx context.Context, client *sdk.CriblIo, groupFilter []string) ([]string, error) {
	streamIDs, streamNames, err := fetchGroupsByProduct(ctx, client, operations.GetProductsGroupsByProductProductStream)
	if err != nil {
		return nil, fmt.Errorf("fetch stream groups: %w", err)
	}
	edgeIDs, edgeNames, err := fetchGroupsByProduct(ctx, client, operations.GetProductsGroupsByProductProductEdge)
	if err != nil {
		return nil, fmt.Errorf("fetch edge groups: %w", err)
	}
	streamIDs, _ = filterGroups(streamIDs, streamNames, " (stream)", groupFilter)
	edgeIDs, _ = filterGroups(edgeIDs, edgeNames, " (edge)", groupFilter)
	groupIDs := make([]string, 0, len(streamIDs)+len(edgeIDs)+1)
	groupIDs = append(groupIDs, streamIDs...)
	groupIDs = append(groupIDs, edgeIDs...)
	if len(groupIDs) == 0 && len(groupFilter) == 0 {
		groupIDs = fallbackGroupIDs()
	} else {
		// Always include "default" and "default_search" so resources under those groups are listed.
		groupIDs = ensureDefaultGroups(groupIDs)
	}
	return groupIDs, nil
}

// ensureDefaultGroups prepends "default" and "default_search" to groupIDs if not already present,
// so resources under those groups (e.g. certificates, parser_lib_entry under default_search) are discovered.
func ensureDefaultGroups(groupIDs []string) []string {
	needDefault := true
	needDefaultSearch := true
	for _, gid := range groupIDs {
		if gid == "default" {
			needDefault = false
		}
		if gid == "default_search" {
			needDefaultSearch = false
		}
		if !needDefault && !needDefaultSearch {
			break
		}
	}
	if !needDefault && !needDefaultSearch {
		return groupIDs
	}
	prepend := make([]string, 0, 2)
	if needDefault {
		prepend = append(prepend, "default")
	}
	if needDefaultSearch {
		prepend = append(prepend, "default_search")
	}
	return append(prepend, groupIDs...)
}

// fallbackGroupIDs returns group IDs to try when the groups API returned none.
// Uses "default", "default_search" (search/parser resources), and, if set, CRIBL_WORKSPACE_ID
// (e.g. cloud workspace "main") so list calls like /m/{groupId}/lib/grok succeed.
func fallbackGroupIDs() []string {
	ids := []string{"default", "default_search"}
	if w := os.Getenv("CRIBL_WORKSPACE_ID"); w != "" && w != "default" && w != "default_search" {
		ids = append(ids, w)
	}
	return ids
}

// ListItemIdentifiers calls the SDK List* method for the entry and returns one identifier map per item.
// Each map has lowercase keys expected by BuildImportID: "id", "group_id", "pack" (as applicable).
// Used by export to fetch each resource and generate import blocks.
func ListItemIdentifiers(ctx context.Context, client *sdk.CriblIo, e registry.Entry, groupIDs []string) ([]map[string]string, error) {
	if client == nil {
		return nil, fmt.Errorf("client is nil")
	}
	if len(groupIDs) == 0 {
		return nil, nil
	}
	// criblio_group has no ListMethod; list via GetProductsGroupsByProduct (stream + edge) and filter by groupIDs.
	if e.TypeName == "criblio_group" && e.ListMethod == "" {
		return listGroupIdentifiers(ctx, client, groupIDs)
	}
	if e.TypeName == "criblio_lakehouse_dataset_connection" {
		return listLakehouseDatasetConnectionIdentifiers(ctx, client)
	}
	if e.TypeName == "criblio_pack_routes" {
		return listPackRoutesIdentifiers(ctx, client, groupIDs)
	}
	clientVal := reflect.ValueOf(client)
	if clientVal.Kind() == reflect.Ptr {
		clientVal = clientVal.Elem()
	}
	svcField := clientVal.FieldByName(e.SDKService)
	if !svcField.IsValid() {
		return nil, fmt.Errorf("SDK service %q not found on client", e.SDKService)
	}
	if svcField.Kind() == reflect.Ptr && svcField.IsNil() {
		return nil, fmt.Errorf("SDK service %q is nil", e.SDKService)
	}
	svc := svcField
	if svc.Kind() == reflect.Ptr {
		svc = svc.Elem()
	}
	method := svc.Addr().MethodByName(e.ListMethod)
	if !method.IsValid() {
		return nil, fmt.Errorf("method %q not found on service %s", e.ListMethod, e.SDKService)
	}
	args := buildListArgs(ctx, method, e, groupIDs[0])
	if len(args) >= 2 && requestRequiresPack(args[1]) {
		return listPackScopedIdentifiers(ctx, client, e, groupIDs, method)
	}
	if len(args) < 2 || !requestHasGroupID(args[1]) {
		items, err := callListAndGetItems(ctx, method, args)
		if err != nil {
			if isSDKUnionUnmarshalError(err) {
				if ids, _, parseErr := searchListIdentifiersFromCapture(e); parseErr == nil && len(ids) > 0 {
					return ids, nil
				}
				if ids, _, parseErr := unionUnmarshalIdentifiersFromCapture(e, groupIDs[0]); parseErr == nil && len(ids) > 0 {
					return ids, nil
				}
				return nil, nil
			}
			return nil, err
		}
		// Pass list scope: groupIDs[0] for group-scoped; "default" for criblio_cribl_lake_dataset (Lake API only has lake "default").
		scope := groupIDs[0]
		if e.TypeName == "criblio_cribl_lake_dataset" {
			scope = "default"
		}
		return identifiersFromItems(items, scope, e)
	}
	var out []map[string]string
	for _, gid := range groupIDs {
		args := buildListArgs(ctx, method, e, gid)
		if len(args) >= 2 && requestRequiresPack(args[1]) {
			// Pack-scoped: list via listPackScopedIdentifiers (called once above with first group).
			// Here we're in the per-group loop but pack-scoped already returned; skip.
			continue
		}
		// One resource per group (e.g. criblio_routes: route table per group); do not iterate items.
		if e.ListUseGroupIDAsItemID {
			out = append(out, map[string]string{"group_id": gid, "id": gid})
			continue
		}
		var ids []map[string]string
		var err error
		switch e.TypeName {
		case "criblio_event_breaker_ruleset":
			ids, err = listIdentifiersForEventBreakerRuleset(ctx, method, e, gid, args)
		case "criblio_search_dataset", "criblio_search_dataset_provider":
			ids, err = listIdentifiersForSearchTypes(ctx, method, e, gid, args)
		default:
			items, listErr := callListAndGetItems(ctx, method, args)
			if listErr != nil {
				if isSDKUnionUnmarshalError(listErr) {
					if parsed, _, parseErr := searchListIdentifiersFromCapture(e); parseErr == nil && len(parsed) > 0 {
						out = append(out, parsed...)
					} else if parsed, _, parseErr := unionUnmarshalIdentifiersFromCapture(e, gid); parseErr == nil && len(parsed) > 0 {
						out = append(out, parsed...)
					}
					continue
				}
				return nil, listErr
			}
			ids, err = identifiersFromItems(items, gid, e)
		}
		if err != nil {
			return nil, err
		}
		if e.TypeName == "criblio_pack" {
			ids = filterPackIdentifiers(ids)
		}
		out = append(out, ids...)
	}
	return out, nil
}

// listIdentifiersForEventBreakerRuleset lists event breaker ruleset identifiers for one group.
// SDK ListEventBreakerRulesetResponseBody has no Items; falls back to captured list response when SDK fails.
func listIdentifiersForEventBreakerRuleset(ctx context.Context, method reflect.Value, e registry.Entry, gid string, args []reflect.Value) ([]map[string]string, error) {
	items, err := callListAndGetItems(ctx, method, args)
	if err != nil {
		if isSDKLibraryUnmarshalError(err) {
			return eventBreakerRulesetIdentifiersFromCapture(gid)
		}
		return nil, err
	}
	ids, err := identifiersFromItems(items, gid, e)
	if err != nil {
		return nil, err
	}
	if len(ids) == 0 {
		if body := custom.GetAndClearEventBreakerRulesetListBody(gid); len(body) > 0 {
			if parsed, parseErr := custom.ParseEventBreakerRulesetListBody(body, gid); parseErr == nil {
				return parsed, nil
			}
		}
	}
	return ids, nil
}

// listIdentifiersForSearchTypes lists search dataset or provider identifiers for one group.
// Falls back to captured list when SDK union unmarshal fails (e.g. cribl_lake).
// When SDK fails and capture parse fails, returns (nil, nil) so caller can continue (no ids for this group).
func listIdentifiersForSearchTypes(ctx context.Context, method reflect.Value, e registry.Entry, gid string, args []reflect.Value) ([]map[string]string, error) {
	items, err := callListAndGetItems(ctx, method, args)
	if err != nil {
		if isSDKUnionUnmarshalError(err) {
			ids, _, parseErr := searchListIdentifiersFromCapture(e)
			if parseErr != nil {
				return nil, nil // continue without ids
			}
			return ids, nil
		}
		return nil, err
	}
	return identifiersFromItems(items, gid, e)
}

// filterPackIdentifiers removes pack IDs in SkipPacks (e.g. HelloPacks).
func filterPackIdentifiers(ids []map[string]string) []map[string]string {
	filtered := ids[:0]
	for _, m := range ids {
		if !custom.SkipPacks[m["id"]] {
			filtered = append(filtered, m)
		}
	}
	return filtered
}

// callListAndGetItems invokes the list method and returns the slice of items (reflect.Value of slice).
func callListAndGetItems(ctx context.Context, method reflect.Value, args []reflect.Value) (reflect.Value, error) {
	var items reflect.Value
	err := func() (err error) {
		defer func() {
			if r := recover(); r != nil {
				err = fmt.Errorf("list method call failed: %v", r)
			}
		}()
		outs := method.Call(args)
		if len(outs) != 2 {
			return fmt.Errorf("unexpected method signature")
		}
		// Safe error extraction: second return may be nil (success); type-asserting nil interface to error panics.
		if v := outs[1].Interface(); v != nil {
			if e, ok := v.(error); ok {
				return e
			}
		}
		respVal := outs[0]
		if respVal.IsNil() {
			return nil
		}
		if respVal.Kind() == reflect.Ptr {
			respVal = respVal.Elem()
		}
		objectField := respVal.FieldByName("Object")
		if !objectField.IsValid() || (objectField.Kind() == reflect.Ptr && objectField.IsNil()) {
			return nil
		}
		// Object is already *ResponseBody; call GetItems on it (not Addr() which would be **ResponseBody).
		getItems := objectField.MethodByName("GetItems")
		if !getItems.IsValid() {
			return nil
		}
		itemOuts := getItems.Call(nil)
		if len(itemOuts) == 0 {
			return nil
		}
		items = itemOuts[0]
		return nil
	}()
	if err != nil {
		return reflect.Value{}, err
	}
	return items, nil
}

// identifiersFromItems builds one identifier map per item. groupID is set for group-scoped resources.
// e is used for ListItemIDMethod (e.g. GetKeyID) and ListUseGroupIDAsItemID (use groupID as id when item has none).
// Items with lib="cribl" (built-in) are skipped so they are not exported.
func identifiersFromItems(items reflect.Value, groupID string, e registry.Entry) ([]map[string]string, error) {
	if items.Kind() != reflect.Slice {
		return nil, nil
	}
	n := items.Len()
	out := make([]map[string]string, 0, n)
	for i := 0; i < n; i++ {
		item := items.Index(i)
		if getLibFromItem(item) == custom.EventBreakerLibCribl {
			continue
		}
		id := getIDFromItem(item, e.ListItemIDMethod)
		if id == "" && e.ListUseGroupIDAsItemID && groupID != "" {
			id = groupID
		}
		if id == "" {
			continue
		}
		// Skip search datasets tagged cribl:default (built-in); only user-created are imported.
		if e.TypeName == "criblio_search_dataset" && itemHasCriblDefaultTagFromItem(item) {
			continue
		}
		// Skip default search dataset providers (cribl_leader, S3, cribl_lake, etc.) so only user-created are imported.
		if e.TypeName == "criblio_search_dataset_provider" && custom.DefaultSearchDatasetProviderIDs[id] {
			continue
		}
		// Skip cribl_lake datasets: they are managed via Lake API and exported as criblio_cribl_lake_dataset only.
		if e.TypeName == "criblio_search_dataset" && getTypeFromItem(item) == custom.SearchDatasetTypeCriblLake {
			continue
		}
		// Skip default Cribl Lake datasets (cribl_logs, default_*, etc.) so only user-created are imported.
		if e.TypeName == "criblio_cribl_lake_dataset" && custom.DefaultCriblLakeDatasetIDs[id] {
			continue
		}
		// Skip built-in destinations (default, devnull) so only user-created are imported.
		if (e.TypeName == "criblio_destination" || e.TypeName == "criblio_pack_destination") && custom.DefaultDestinationIDs[id] {
			continue
		}
		m := map[string]string{"id": id}
		if e.TypeName == "criblio_cribl_lake_dataset" {
			// Lake API uses lake_id (default "default"); always set so HCL has required lake_id.
			if groupID != "" {
				m["lake_id"] = groupID
			} else {
				m["lake_id"] = "default"
			}
		} else if e.TypeName == "criblio_cribl_lake_house" {
			// Lakehouse API is not group-scoped; omit group_id so export puts under global.
		} else if groupID != "" {
			m["group_id"] = groupID
		}
		out = append(out, m)
	}
	return out, nil
}

// getLibFromItem returns the lib value from a list item (struct or map). Used to skip built-in items (lib="cribl").
func getLibFromItem(item reflect.Value) string {
	if item.Kind() == reflect.Ptr && !item.IsNil() {
		item = item.Elem()
	}
	if item.Kind() == reflect.Map {
		for _, key := range []string{"lib", "Lib"} {
			k := item.MapIndex(reflect.ValueOf(key))
			if !k.IsValid() {
				continue
			}
			if k.Kind() == reflect.Interface {
				k = reflect.ValueOf(k.Interface())
			}
			for k.Kind() == reflect.Ptr && k.IsValid() && !k.IsNil() {
				k = k.Elem()
			}
			if k.Kind() == reflect.String {
				return k.String()
			}
		}
		return ""
	}
	getLib := item.Addr().MethodByName("GetLib")
	if !getLib.IsValid() {
		return ""
	}
	outs := getLib.Call(nil)
	if len(outs) == 0 {
		return ""
	}
	v := outs[0]
	if v.Kind() == reflect.Ptr && v.IsNil() {
		return ""
	}
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	if v.Kind() == reflect.String {
		return v.String()
	}
	// Named types with underlying string (e.g. shared.CriblLib "custom"/"cribl") for Knowledge resources.
	if v.IsValid() && v.CanConvert(reflect.TypeOf("")) {
		return v.Convert(reflect.TypeOf("")).String()
	}
	return ""
}

// itemHasCriblDefaultTagFromItem returns true if the item has tags containing "cribl:default" (string or slice).
// Used to skip built-in search datasets when the SDK list unmarshals successfully.
func itemHasCriblDefaultTagFromItem(item reflect.Value) bool {
	if item.Kind() == reflect.Ptr && !item.IsNil() {
		item = item.Elem()
	}
	tag := custom.CriblDefaultTag
	if item.Kind() == reflect.Map {
		k := item.MapIndex(reflect.ValueOf("tags"))
		if !k.IsValid() {
			return false
		}
		if k.Kind() == reflect.Interface {
			k = reflect.ValueOf(k.Interface())
		}
		for k.Kind() == reflect.Ptr && k.IsValid() && !k.IsNil() {
			k = k.Elem()
		}
		if k.Kind() == reflect.String {
			return k.String() == tag
		}
		if k.Kind() == reflect.Slice {
			for j := 0; j < k.Len(); j++ {
				el := k.Index(j)
				if el.Kind() == reflect.Interface {
					el = reflect.ValueOf(el.Interface())
				}
				if el.Kind() == reflect.String && el.String() == tag {
					return true
				}
			}
		}
		return false
	}
	// Struct: field Tags or method GetTags
	f := item.FieldByName("Tags")
	if f.IsValid() {
		for f.Kind() == reflect.Ptr && f.IsValid() && !f.IsNil() {
			f = f.Elem()
		}
		if f.Kind() == reflect.String {
			return f.String() == tag
		}
		if f.Kind() == reflect.Slice {
			for j := 0; j < f.Len(); j++ {
				el := f.Index(j)
				if el.Kind() == reflect.Interface {
					el = reflect.ValueOf(el.Interface())
				}
				if el.Kind() == reflect.String && el.String() == tag {
					return true
				}
			}
		}
	}
	getTags := item.Addr().MethodByName("GetTags")
	if getTags.IsValid() && getTags.Type().NumOut() >= 1 {
		outs := getTags.Call(nil)
		v := outs[0]
		for v.Kind() == reflect.Ptr && v.IsValid() && !v.IsNil() {
			v = v.Elem()
		}
		if v.Kind() == reflect.String {
			return v.String() == tag
		}
		if v.Kind() == reflect.Slice {
			for j := 0; j < v.Len(); j++ {
				el := v.Index(j)
				if el.Kind() == reflect.Interface {
					el = reflect.ValueOf(el.Interface())
				}
				if el.Kind() == reflect.String && el.String() == tag {
					return true
				}
			}
		}
	}
	return false
}

// getTypeFromItem returns the type string from a list item (e.g. "cribl_lake" for search datasets).
// Used to skip cribl_lake items so they are only exported as criblio_cribl_lake_dataset.
func getTypeFromItem(item reflect.Value) string {
	if item.Kind() == reflect.Ptr && !item.IsNil() {
		item = item.Elem()
	}
	if item.Kind() == reflect.Map {
		for _, key := range []string{"type", "Type"} {
			k := item.MapIndex(reflect.ValueOf(key))
			if !k.IsValid() {
				continue
			}
			if k.Kind() == reflect.Interface {
				k = reflect.ValueOf(k.Interface())
			}
			for k.Kind() == reflect.Ptr && k.IsValid() && !k.IsNil() {
				k = k.Elem()
			}
			if k.Kind() == reflect.String {
				return k.String()
			}
		}
		return ""
	}
	getType := item.Addr().MethodByName("GetType")
	if !getType.IsValid() {
		return ""
	}
	outs := getType.Call(nil)
	if len(outs) == 0 {
		return ""
	}
	v := outs[0]
	if v.Kind() == reflect.Ptr && v.IsNil() {
		return ""
	}
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	if v.Kind() == reflect.String {
		return v.String()
	}
	if v.IsValid() && v.CanConvert(reflect.TypeOf("")) {
		return v.Convert(reflect.TypeOf("")).String()
	}
	return ""
}

// getIDFromItem returns the ID from a list item (struct or map). itemIDMethod if non-empty is tried first (e.g. "GetKeyID").
func getIDFromItem(item reflect.Value, itemIDMethod string) string {
	if item.Kind() == reflect.Ptr && !item.IsNil() {
		item = item.Elem()
	}
	if item.Kind() == reflect.Map {
		for _, key := range []string{"id", "Id"} {
			k := item.MapIndex(reflect.ValueOf(key))
			if !k.IsValid() {
				continue
			}
			if k.Kind() == reflect.Interface {
				k = reflect.ValueOf(k.Interface())
			}
			for k.Kind() == reflect.Ptr && k.IsValid() && !k.IsNil() {
				k = k.Elem()
			}
			if k.Kind() == reflect.String {
				return k.String()
			}
		}
		return ""
	}
	// Struct: try custom method first, then GetID
	if itemIDMethod != "" {
		m := item.Addr().MethodByName(itemIDMethod)
		if m.IsValid() {
			outs := m.Call(nil)
			if len(outs) > 0 {
				v := outs[0]
				if v.Kind() == reflect.Ptr && !v.IsNil() {
					v = v.Elem()
				}
				if v.Kind() == reflect.String {
					return v.String()
				}
			}
		}
	}
	getID := item.Addr().MethodByName("GetID")
	if getID.IsValid() {
		outs := getID.Call(nil)
		if len(outs) > 0 {
			v := outs[0]
			if v.Kind() == reflect.Ptr && !v.IsNil() {
				v = v.Elem()
			}
			if v.Kind() == reflect.String {
				return v.String()
			}
		}
	}
	// Union-type structs (e.g. GenericDataset, GenericProvider) have no GetID on the wrapper; get ID from the non-nil inner field.
	return getIDFromUnionStruct(item)
}

// getIDFromUnionStruct returns the ID from a union-style struct whose fields are pointers to concrete types with GetID.
// Used for ListDataset/ListDatasetProvider items (GenericDataset, GenericProvider).
func getIDFromUnionStruct(item reflect.Value) string {
	if item.Kind() == reflect.Ptr && !item.IsNil() {
		item = item.Elem()
	}
	if item.Kind() != reflect.Struct {
		return ""
	}
	for i := 0; i < item.NumField(); i++ {
		f := item.Field(i)
		if f.Kind() != reflect.Ptr || f.IsNil() {
			continue
		}
		inner := f.Elem()
		if inner.Kind() != reflect.Struct {
			continue
		}
		getID := inner.Addr().MethodByName("GetID")
		if !getID.IsValid() {
			continue
		}
		outs := getID.Call(nil)
		if len(outs) == 0 {
			continue
		}
		v := outs[0]
		if v.Kind() == reflect.Ptr && !v.IsNil() {
			v = v.Elem()
		}
		if v.Kind() == reflect.String && v.String() != "" {
			return v.String()
		}
	}
	return ""
}

// listOne calls the SDK list method for the given entry and returns total count and per-group counts.
// When the list method takes a request with GroupID, it is called once per groupID; perGroup maps groupID -> count.
// When the method has no GroupID, perGroup is nil.
func listOne(ctx context.Context, client *sdk.CriblIo, e registry.Entry, groupIDs []string) (total int, perGroup map[string]int, err error) {
	// No groups (e.g. user filtered to edge-only): return zero, no API calls.
	if len(groupIDs) == 0 {
		return 0, nil, nil
	}
	clientVal := reflect.ValueOf(client)
	if clientVal.Kind() == reflect.Ptr {
		clientVal = clientVal.Elem()
	}
	svcField := clientVal.FieldByName(e.SDKService)
	if !svcField.IsValid() {
		return 0, nil, fmt.Errorf("SDK service %q not found on client", e.SDKService)
	}
	if svcField.Kind() == reflect.Ptr && svcField.IsNil() {
		return 0, nil, fmt.Errorf("SDK service %q is nil", e.SDKService)
	}
	svc := svcField
	if svc.Kind() == reflect.Ptr {
		svc = svc.Elem()
	}
	method := svc.Addr().MethodByName(e.ListMethod)
	if !method.IsValid() {
		return 0, nil, fmt.Errorf("method %q not found on service %s", e.ListMethod, e.SDKService)
	}

	args := buildListArgs(ctx, method, e, groupIDs[0])
	if len(args) >= 2 && requestRequiresPack(args[1]) {
		return countPackScoped(ctx, client, e, groupIDs, method)
	}
	// If the method has no request or request has no GroupID, call once.
	if len(args) < 2 || !requestHasGroupID(args[1]) {
		// For search dataset/provider types, count after filtering out defaults so preview matches export.
		if e.TypeName == "criblio_search_dataset_provider" || e.TypeName == "criblio_search_dataset" {
			items, listErr := callListAndGetItems(ctx, method, args)
			if listErr != nil && isSDKUnionUnmarshalError(listErr) {
				if _, count, parseErr := searchListIdentifiersFromCapture(e); parseErr == nil && count > 0 {
					return count, nil, nil
				}
				return 0, nil, nil
			}
			if listErr != nil {
				return 0, nil, listErr
			}
			ids, idErr := identifiersFromItems(items, groupIDs[0], e)
			if idErr != nil {
				return 0, nil, idErr
			}
			return len(ids), nil, nil
		}
		// criblio_cribl_lake_dataset: count after filtering DefaultCriblLakeDatasetIDs so discovery count matches export.
		if e.TypeName == "criblio_cribl_lake_dataset" {
			items, listErr := callListAndGetItems(ctx, method, args)
			if listErr != nil {
				return 0, nil, listErr
			}
			ids, idErr := identifiersFromItems(items, groupIDs[0], e)
			if idErr != nil {
				return 0, nil, idErr
			}
			return len(ids), nil, nil
		}
		n, err := callListAndCount(ctx, method, args)
		if err != nil && isSDKUnionUnmarshalError(err) {
			if _, count, parseErr := searchListIdentifiersFromCapture(e); parseErr == nil && count > 0 {
				return count, nil, nil
			}
			if _, count, parseErr := unionUnmarshalIdentifiersFromCapture(e, groupIDs[0]); parseErr == nil && count > 0 {
				return count, nil, nil
			}
			return 0, nil, nil
		}
		return n, nil, err
	}
	// Request has GroupID: call once per group and sum counts; record per-group.
	perGroup = make(map[string]int)
	for _, gid := range groupIDs {
		args := buildListArgs(ctx, method, e, gid)
		if len(args) >= 2 && requestRequiresPack(args[1]) {
			return total, perGroup, nil
		}
		var n int
		switch e.TypeName {
		case "criblio_search_dataset", "criblio_search_dataset_provider":
			ids, listErr := listIdentifiersForSearchTypes(ctx, method, e, gid, args)
			if listErr != nil {
				return total, perGroup, listErr
			}
			n = len(ids)
		case "criblio_event_breaker_ruleset":
			ids, listErr := listIdentifiersForEventBreakerRuleset(ctx, method, e, gid, args)
			if listErr != nil {
				return total, perGroup, listErr
			}
			n = len(ids)
		default:
			var err error
			n, err = callListAndCount(ctx, method, args)
			if err != nil {
				if isSDKUnionUnmarshalError(err) {
					if _, count, parseErr := searchListIdentifiersFromCapture(e); parseErr == nil && count > 0 {
						total += count
						perGroup[gid] = count
					} else if _, count, parseErr := unionUnmarshalIdentifiersFromCapture(e, gid); parseErr == nil && count > 0 {
						total += count
						perGroup[gid] = count
					}
					continue
				}
				return total, perGroup, err
			}
			// criblio_routes: API may return 0 items; count as 1 resource per group when ListUseGroupIDAsItemID.
			if n == 0 && e.ListUseGroupIDAsItemID {
				n = 1
			}
		}
		total += n
		perGroup[gid] = n
	}
	return total, perGroup, nil
}

// callListAndCount invokes the list method with the given args and returns the item count.
// args are built from method.Type() in buildListArgs so they match the method signature.
// We use reflection because the registry drives which SDK List* method to call by name;
// a panic from type mismatch is recovered and returned as an error.
func callListAndCount(ctx context.Context, method reflect.Value, args []reflect.Value) (n int, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("list method call failed: %v", r)
		}
	}()
	outs := method.Call(args)
	if len(outs) != 2 {
		return 0, fmt.Errorf("unexpected method signature")
	}
	respVal := outs[0]
	errVal := outs[1]
	if !errVal.IsNil() {
		return 0, errVal.Interface().(error)
	}
	if respVal.IsNil() {
		return 0, nil
	}
	return countItemsFromResponse(respVal)
}

func buildListArgs(ctx context.Context, method reflect.Value, e registry.Entry, groupID string) []reflect.Value {
	if groupID == "" {
		groupID = "default"
	}
	// Lake API list (criblio_cribl_lake_dataset) only accepts lakeId "default"; do not use groupIDs[0] (e.g. "default_search").
	if e.TypeName == "criblio_cribl_lake_dataset" {
		groupID = "default"
	}
	mt := method.Type()
	args := []reflect.Value{reflect.ValueOf(ctx)}
	if mt.NumIn() >= 2 {
		param1 := mt.In(1)
		if param1.Kind() != reflect.Slice {
			reqVal := reflect.New(param1)
			setGroupIDDefault(reqVal, groupID)
			args = append(args, reqVal.Elem())
		}
	}
	return args
}

// requestHasGroupID reports whether the request struct has a GroupID field (used for path /m/{groupId}/...).
func requestHasGroupID(reqVal reflect.Value) bool {
	if reqVal.Kind() == reflect.Ptr {
		reqVal = reqVal.Elem()
	}
	if !reqVal.IsValid() || reqVal.Kind() != reflect.Struct {
		return false
	}
	f := reqVal.FieldByName("GroupID")
	return f.IsValid() && f.Kind() == reflect.String
}

func setGroupIDDefault(req reflect.Value, defaultGroupID string) {
	if req.Kind() == reflect.Ptr {
		req = req.Elem()
	}
	if !req.IsValid() || req.Kind() != reflect.Struct {
		return
	}
	// Set GroupID (path /m/{groupId}/...) or LakeID (path /lakes/{lakeId}/...) so list request is valid.
	for _, name := range []string{"GroupID", "LakeID"} {
		f := req.FieldByName(name)
		if f.IsValid() && f.CanSet() && f.Kind() == reflect.String {
			f.SetString(defaultGroupID)
		}
	}
}

func setPackDefault(req reflect.Value, packID string) {
	if req.Kind() == reflect.Ptr {
		req = req.Elem()
	}
	if !req.IsValid() || req.Kind() != reflect.Struct {
		return
	}
	f := req.FieldByName("Pack")
	if f.IsValid() && f.CanSet() && f.Kind() == reflect.String {
		f.SetString(packID)
	}
}

// getPackIDsByGroup returns pack IDs for the given group via Packs.GetPacksByGroup.
// Used to discover pack-scoped resources (pack_destination, pack_pipeline, pack_source).
func getPackIDsByGroup(ctx context.Context, client *sdk.CriblIo, groupID string) ([]string, error) {
	if client == nil || client.Packs == nil {
		return nil, fmt.Errorf("packs service nil")
	}
	resp, err := client.Packs.GetPacksByGroup(ctx, operations.GetPacksByGroupRequest{GroupID: groupID})
	if err != nil {
		return nil, err
	}
	if resp == nil || resp.Object == nil {
		return nil, nil
	}
	items := resp.Object.GetItems()
	ids := make([]string, 0, len(items))
	for _, item := range items {
		if id := item.GetID(); id != "" && !custom.SkipPacks[id] {
			ids = append(ids, id)
		}
	}
	return ids, nil
}

// buildListArgsWithPack builds [ctx, request] with GroupID and Pack set for pack-scoped list methods.
func buildListArgsWithPack(ctx context.Context, method reflect.Value, e registry.Entry, groupID, packID string) []reflect.Value {
	if groupID == "" {
		groupID = "default"
	}
	mt := method.Type()
	args := []reflect.Value{reflect.ValueOf(ctx)}
	if mt.NumIn() >= 2 {
		param1 := mt.In(1)
		if param1.Kind() != reflect.Slice {
			reqVal := reflect.New(param1)
			setGroupIDDefault(reqVal, groupID)
			setPackDefault(reqVal, packID)
			args = append(args, reqVal.Elem())
		}
	}
	return args
}

// listPackRoutesIdentifiers discovers criblio_pack_routes: one resource per (group_id, pack).
// There is no list API that returns items; we enumerate groups × packs and return one identifier map per pair.
func listPackRoutesIdentifiers(ctx context.Context, client *sdk.CriblIo, groupIDs []string) ([]map[string]string, error) {
	var out []map[string]string
	for _, gid := range groupIDs {
		packIDs, err := getPackIDsByGroup(ctx, client, gid)
		if err != nil {
			return nil, fmt.Errorf("criblio_pack_routes: get packs for group %s: %w", gid, err)
		}
		for _, packID := range packIDs {
			out = append(out, map[string]string{"group_id": gid, "pack": packID})
		}
	}
	return out, nil
}

// listPackScopedIdentifiers discovers identifiers for pack-scoped resources (pack_destination, pack_pipeline, pack_source)
// by listing packs per group then calling the list method for each (groupID, packID).
func listPackScopedIdentifiers(ctx context.Context, client *sdk.CriblIo, e registry.Entry, groupIDs []string, method reflect.Value) ([]map[string]string, error) {
	var out []map[string]string
	for _, gid := range groupIDs {
		packIDs, err := getPackIDsByGroup(ctx, client, gid)
		if err != nil {
			return nil, fmt.Errorf("%s: get packs for group %s: %w", e.TypeName, gid, err)
		}
		for _, packID := range packIDs {
			args := buildListArgsWithPack(ctx, method, e, gid, packID)
			items, err := callListAndGetItems(ctx, method, args)
			var ids []map[string]string
			if err != nil {
				// criblio_pack_source: SDK may fail to unmarshal Input union; use captured list response.
				if e.TypeName == "criblio_pack_source" {
					body := custom.GetAndClearPackInputsListBody(gid, packID)
					if len(body) > 0 {
						capturedIDs, parseErr := custom.ParsePackInputsListBody(body)
						if parseErr == nil {
							for _, m := range capturedIDs {
								m["group_id"] = gid
								m["pack"] = packID
								out = append(out, m)
							}
						}
					}
					continue
				}
				// criblio_pack_breakers: SDK may fail when lib="cribl" (not in EventBreakerRuleset enum); use captured list response.
				if e.TypeName == "criblio_pack_breakers" && isSDKLibraryUnmarshalError(err) {
					body := custom.GetAndClearPackBreakersListBody(gid, packID)
					if len(body) > 0 {
						if capturedIDs, parseErr := custom.ParsePackBreakersListBody(body, gid, packID); parseErr == nil {
							out = append(out, capturedIDs...)
						}
					}
					continue
				}
				return nil, fmt.Errorf("%s: list pack %s/%s: %w", e.TypeName, gid, packID, err)
			}
			ids, err = identifiersFromItems(items, gid, e)
			if err != nil {
				return nil, err
			}
			// criblio_pack_destination: SDK ListPackOutputResponseBody has no Items; use captured list response.
			if e.TypeName == "criblio_pack_destination" && len(ids) == 0 {
				body := custom.GetAndClearPackOutputsListBody(gid, packID)
				if len(body) > 0 {
					capturedIDs, parseErr := custom.ParsePackOutputsListBody(body)
					if parseErr == nil {
						for _, m := range capturedIDs {
							m["group_id"] = gid
						}
						ids = capturedIDs
					}
				}
			}
			// criblio_pack_source: SDK may return empty items when Input union unmarshals with all fields null; use captured list response.
			if e.TypeName == "criblio_pack_source" && len(ids) == 0 {
				body := custom.GetAndClearPackInputsListBody(gid, packID)
				if len(body) > 0 {
					capturedIDs, parseErr := custom.ParsePackInputsListBody(body)
					if parseErr == nil {
						for _, m := range capturedIDs {
							m["group_id"] = gid
						}
						ids = capturedIDs
					}
				}
			}
			for _, m := range ids {
				if e.TypeName == "criblio_pack_destination" && custom.DefaultDestinationIDs[m["id"]] {
					continue
				}
				m["pack"] = packID
				out = append(out, m)
			}
		}
	}
	return out, nil
}

// countPackScoped returns total count and per-group counts for pack-scoped list methods.
func countPackScoped(ctx context.Context, client *sdk.CriblIo, e registry.Entry, groupIDs []string, method reflect.Value) (total int, perGroup map[string]int, err error) {
	perGroup = make(map[string]int)
	for _, gid := range groupIDs {
		packIDs, err := getPackIDsByGroup(ctx, client, gid)
		if err != nil {
			return 0, nil, fmt.Errorf("%s: get packs for group %s: %w", e.TypeName, gid, err)
		}
		var groupCount int
		for _, packID := range packIDs {
			args := buildListArgsWithPack(ctx, method, e, gid, packID)
			n, err := callListAndCount(ctx, method, args)
			if err != nil {
				// criblio_pack_source: SDK may fail to unmarshal Input union; use captured list body for count.
				if e.TypeName == "criblio_pack_source" {
					body := custom.GetAndClearPackInputsListBody(gid, packID)
					if len(body) > 0 {
						if ids, parseErr := custom.ParsePackInputsListBody(body); parseErr == nil {
							n = len(ids)
						}
					} else {
						n = 0
					}
				} else if e.TypeName == "criblio_pack_breakers" && isSDKLibraryUnmarshalError(err) {
					body := custom.GetAndClearPackBreakersListBody(gid, packID)
					if len(body) > 0 {
						if ids, parseErr := custom.ParsePackBreakersListBody(body, gid, packID); parseErr == nil {
							n = len(ids)
						}
					} else {
						n = 0
					}
				} else {
					return 0, nil, fmt.Errorf("%s: list pack %s/%s: %w", e.TypeName, gid, packID, err)
				}
			}
			// criblio_pack_destination: SDK response has no Items; use captured list body for count.
			if e.TypeName == "criblio_pack_destination" && n == 0 {
				body := custom.GetAndClearPackOutputsListBody(gid, packID)
				if len(body) > 0 {
					if ids, parseErr := custom.ParsePackOutputsListBody(body); parseErr == nil {
						n = len(ids)
					}
				}
			}
			// criblio_pack_source: SDK may return 0 when Input union unmarshals with all fields null; use captured list body.
			if e.TypeName == "criblio_pack_source" && n == 0 {
				body := custom.GetAndClearPackInputsListBody(gid, packID)
				if len(body) > 0 {
					if ids, parseErr := custom.ParsePackInputsListBody(body); parseErr == nil {
						n = len(ids)
					}
				}
			}
			groupCount += n
		}
		total += groupCount
		if groupCount > 0 {
			perGroup[gid] = groupCount
		}
	}
	return total, perGroup, nil
}

// requestRequiresPack reports whether the request struct has a Pack field that is
// required for the API path (e.g. /m/{groupId}/p/{pack}/...) and is currently empty.
func requestRequiresPack(reqVal reflect.Value) bool {
	if reqVal.Kind() == reflect.Ptr {
		reqVal = reqVal.Elem()
	}
	if !reqVal.IsValid() || reqVal.Kind() != reflect.Struct {
		return false
	}
	f := reqVal.FieldByName("Pack")
	if !f.IsValid() || f.Kind() != reflect.String {
		return false
	}
	return f.String() == ""
}

func countItemsFromResponse(respVal reflect.Value) (int, error) {
	if respVal.Kind() == reflect.Ptr {
		respVal = respVal.Elem()
	}
	objectField := respVal.FieldByName("Object")
	if !objectField.IsValid() || (objectField.Kind() == reflect.Ptr && objectField.IsNil()) {
		return 0, nil
	}
	// Call GetItems on the Object (pointer receiver)
	getItems := objectField.MethodByName("GetItems")
	if !getItems.IsValid() {
		return 0, nil
	}
	outs := getItems.Call(nil)
	if len(outs) == 0 {
		return 0, nil
	}
	items := outs[0]
	if items.Kind() == reflect.Slice {
		return items.Len(), nil
	}
	return 0, nil
}
