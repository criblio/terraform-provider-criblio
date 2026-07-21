// Package discovery implements REST-backed discovery for existing resources.
package discovery

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"os"
	"slices"
	"strings"

	"github.com/criblio/terraform-provider-criblio/internal/restclient"
	importclient "github.com/criblio/terraform-provider-criblio/tools/import-cli/internal/client"
	"github.com/criblio/terraform-provider-criblio/tools/import-cli/internal/custom"
	"github.com/criblio/terraform-provider-criblio/tools/import-cli/internal/registry"
)

// Result holds the discovery result for one resource type.
type Result struct {
	TypeName       string
	Count          int
	Err            error
	Details        []string
	PerGroupCounts map[string]int
}

// IsRecoverableListDecodeError reports whether err should not abort export.
func IsRecoverableListDecodeError(err error) bool {
	for e := err; e != nil; e = errors.Unwrap(e) {
		if strings.Contains(e.Error(), "decode response") ||
			strings.Contains(e.Error(), "unmarshal") {
			return true
		}
	}
	return false
}

// Discover enumerates resources through REST list endpoints.
func Discover(ctx context.Context, client *importclient.Client, reg *registry.Registry, include, exclude, groupFilter []string, onPrem bool) ([]Result, error) {
	if client == nil || client.REST == nil {
		return nil, fmt.Errorf("REST client is nil")
	}

	includeSet := sliceToSet(include)
	excludeSet := sliceToSet(exclude)

	streamIDs, streamNames, streamErr := fetchGroupsByProduct(ctx, client, "stream")
	if streamErr != nil {
		return nil, fmt.Errorf("fetch stream groups: %w", streamErr)
	}
	edgeIDs, edgeNames, edgeErr := fetchGroupsByProduct(ctx, client, "edge")
	if edgeErr != nil {
		return nil, fmt.Errorf("fetch edge groups: %w", edgeErr)
	}

	streamIDs, streamNames = filterGroups(streamIDs, streamNames, " (stream)", groupFilter)
	edgeIDs, edgeNames = filterGroups(edgeIDs, edgeNames, " (edge)", groupFilter)

	groupIDs := make([]string, 0, len(streamIDs)+len(edgeIDs)+1)
	groupIDs = append(groupIDs, streamIDs...)
	groupIDs = append(groupIDs, edgeIDs...)
	if len(groupIDs) == 0 && len(groupFilter) == 0 {
		groupIDs = fallbackGroupIDs(onPrem)
	} else if len(groupFilter) == 0 {
		groupIDs = ensureDefaultGroups(groupIDs)
	}
	if onPrem {
		groupIDs = filterOutDefaultSearch(groupIDs)
	}

	idToLabel := make(map[string]string)
	for i := range streamIDs {
		idToLabel[streamIDs[i]] = streamNames[i] + " (stream)"
	}
	for i := range edgeIDs {
		idToLabel[edgeIDs[i]] = edgeNames[i] + " (edge)"
	}

	groupNames := make([]string, 0, len(streamNames)+len(edgeNames))
	for _, name := range streamNames {
		groupNames = append(groupNames, name+" (stream)")
	}
	for _, name := range edgeNames {
		groupNames = append(groupNames, name+" (edge)")
	}

	var results []Result
	for _, e := range reg.Entries() {
		if !matchesFilter(e.TypeName, includeSet, excludeSet) {
			continue
		}
		if len(groupFilter) > 0 && skipDiscoveryForGroupFilter(e.TypeName, groupIDs) {
			results = append(results, Result{TypeName: e.TypeName})
			continue
		}

		res := Result{TypeName: e.TypeName}
		switch {
		case e.TypeName == "criblio_group":
			res.Count = len(streamNames) + len(edgeNames)
			res.Details = groupNames
		case e.TypeName == "criblio_custom_banner":
			_, err := restclient.Get[map[string]json.RawMessage](ctx, client.REST, "/system/banners/custom-banner")
			if restclient.IsNotFound(err) {
				res.Count = 0
			} else if err != nil {
				res.Err = err
			} else {
				res.Count = 1
			}
		case e.TypeName == "criblio_lakehouse_dataset_connection":
			ids, err := listLakehouseDatasetConnectionIdentifiers(ctx, client)
			res.Count = len(ids)
			res.Err = err
		case e.TypeName == "criblio_pack_routes":
			ids, err := listPackRoutesIdentifiers(ctx, client, groupIDs)
			res.Count = len(ids)
			res.Err = err
		case e.TypeName == "criblio_search_dataset_ruleset" || e.TypeName == "criblio_search_datatype_ruleset":
			if slices.Contains(groupIDs, "default_search") {
				res.Count = 1
			}
		case e.RESTListPath != "":
			count, perGroup, err := listOneREST(ctx, client, e, groupIDs)
			res.Count = count
			res.Err = err
			if len(perGroup) > 0 {
				res.PerGroupCounts = make(map[string]int)
				for gid, count := range perGroup {
					if label, ok := idToLabel[gid]; ok {
						res.PerGroupCounts[label] = count
					} else {
						res.PerGroupCounts[gid] = count
					}
				}
			}
		}
		if res.Err != nil {
			res.Err = fmt.Errorf("%s: %w", e.TypeName, res.Err)
		}
		results = append(results, res)
	}
	return results, nil
}

// GetGroupIDs returns group IDs used for list/export.
func GetGroupIDs(ctx context.Context, client *importclient.Client, groupFilter []string, onPrem bool) ([]string, error) {
	streamIDs, streamNames, err := fetchGroupsByProduct(ctx, client, "stream")
	if err != nil {
		return nil, fmt.Errorf("fetch stream groups: %w", err)
	}
	edgeIDs, edgeNames, err := fetchGroupsByProduct(ctx, client, "edge")
	if err != nil {
		return nil, fmt.Errorf("fetch edge groups: %w", err)
	}
	streamIDs, _ = filterGroups(streamIDs, streamNames, " (stream)", groupFilter)
	edgeIDs, _ = filterGroups(edgeIDs, edgeNames, " (edge)", groupFilter)

	groupIDs := make([]string, 0, len(streamIDs)+len(edgeIDs)+1)
	groupIDs = append(groupIDs, streamIDs...)
	groupIDs = append(groupIDs, edgeIDs...)
	if len(groupIDs) == 0 && len(groupFilter) == 0 {
		groupIDs = fallbackGroupIDs(onPrem)
	} else if len(groupFilter) == 0 {
		groupIDs = ensureDefaultGroups(groupIDs)
	}
	if onPrem {
		groupIDs = filterOutDefaultSearch(groupIDs)
	}
	return groupIDs, nil
}

// ListItemIdentifiers returns one identifier map per item.
func ListItemIdentifiers(ctx context.Context, client *importclient.Client, e registry.Entry, groupIDs []string) ([]map[string]string, error) {
	if len(groupIDs) == 0 {
		return nil, nil
	}
	switch e.TypeName {
	case "criblio_custom_banner":
		return []map[string]string{{"id": "custom-banner"}}, nil
	case "criblio_group":
		idMaps, _, err := ListGroupIdentifiersAndItems(ctx, client, groupIDs)
		return idMaps, err
	case "criblio_lakehouse_dataset_connection":
		return listLakehouseDatasetConnectionIdentifiers(ctx, client)
	case "criblio_search_dataset_ruleset", "criblio_search_datatype_ruleset":
		if slices.Contains(groupIDs, "default_search") {
			return []map[string]string{{"id": "default", "group_id": "default_search"}}, nil
		}
		return nil, nil
	case "criblio_pack_routes":
		return listPackRoutesIdentifiers(ctx, client, groupIDs)
	}
	if e.RESTListPath == "" {
		return nil, nil
	}
	ids, _, err := listRESTIdentifiers(ctx, client, e, groupIDs)
	return ids, err
}

// ListGroupIdentifiersAndItems returns identifier maps and raw group API items.
func ListGroupIdentifiersAndItems(ctx context.Context, client *importclient.Client, groupIDs []string) (idMaps []map[string]string, items []json.RawMessage, err error) {
	groupIDSet := sliceToSet(groupIDs)
	for _, product := range []string{"stream", "edge"} {
		rawItems, err := getRESTItems(ctx, client, fmt.Sprintf("/products/%s/groups", product))
		if err != nil {
			return nil, nil, err
		}
		for _, raw := range rawItems {
			item, err := rawMap(raw)
			if err != nil {
				return nil, nil, err
			}
			id := rawString(item, "id", "ID", "name")
			if id == "" {
				continue
			}
			if len(groupIDSet) > 0 && !inSet(id, groupIDSet) {
				continue
			}
			idMaps = append(idMaps, map[string]string{"group_id": id, "id": id, "product": product})
			items = append(items, raw)
		}
	}
	return idMaps, items, nil
}

func fetchGroupsByProduct(ctx context.Context, client *importclient.Client, product string) (ids []string, names []string, err error) {
	items, err := getRESTItems(ctx, client, fmt.Sprintf("/products/%s/groups", product))
	if err != nil {
		return nil, nil, err
	}
	for _, raw := range items {
		item, err := rawMap(raw)
		if err != nil {
			return nil, nil, err
		}
		id := rawString(item, "id", "ID", "name")
		if id == "" {
			continue
		}
		ids = append(ids, id)
		if name := rawString(item, "name"); name != "" {
			names = append(names, name)
		} else {
			names = append(names, id)
		}
	}
	return ids, names, nil
}

func listOneREST(ctx context.Context, client *importclient.Client, e registry.Entry, groupIDs []string) (int, map[string]int, error) {
	ids, perGroup, err := listRESTIdentifiers(ctx, client, e, groupIDs)
	if err != nil {
		return 0, nil, err
	}
	return len(ids), perGroup, nil
}

func listRESTIdentifiers(ctx context.Context, client *importclient.Client, e registry.Entry, groupIDs []string) ([]map[string]string, map[string]int, error) {
	if client == nil || client.REST == nil {
		return nil, nil, fmt.Errorf("REST client is nil")
	}

	var out []map[string]string
	perGroup := map[string]int{}
	for _, scope := range restScopes(e.RESTListPath, groupIDs) {
		gid := scope["group_id"]
		if gid != "" && skipGroupScopedSingleton(e.TypeName, gid) {
			continue
		}
		if pathUsesRESTParam(e.RESTListPath, "pack") {
			packIDs, err := getPackIDsByGroup(ctx, client, gid)
			if err != nil {
				return nil, nil, err
			}
			groupCount := 0
			for _, packID := range packIDs {
				scope["pack"] = packID
				ids, err := listRESTPathIdentifiers(ctx, client, e, scope)
				if err != nil {
					return nil, nil, err
				}
				for _, id := range ids {
					id["pack"] = packID
				}
				groupCount += len(ids)
				out = append(out, ids...)
			}
			if gid != "" {
				perGroup[gid] = groupCount
			}
			continue
		}
		ids, err := listRESTPathIdentifiers(ctx, client, e, scope)
		if err != nil {
			return nil, nil, err
		}
		if gid != "" {
			perGroup[gid] = len(ids)
		}
		out = append(out, ids...)
	}
	if len(perGroup) == 0 {
		perGroup = nil
	}
	return out, perGroup, nil
}

func listRESTPathIdentifiers(ctx context.Context, client *importclient.Client, e registry.Entry, scope map[string]string) ([]map[string]string, error) {
	path := renderRESTPath(e.RESTListPath, scope)
	items, err := getRESTItems(ctx, client, path)
	if restclient.IsNotFound(err) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return identifiersFromRawItems(items, scope, e)
}

func restScopes(path string, groupIDs []string) []map[string]string {
	if pathUsesRESTParam(path, "group_id") {
		scopes := make([]map[string]string, 0, len(groupIDs))
		for _, gid := range groupIDs {
			scopes = append(scopes, map[string]string{"group_id": gid})
		}
		return scopes
	}
	if pathUsesRESTParam(path, "lake_id") {
		return []map[string]string{{"lake_id": "default"}}
	}
	if pathUsesRESTParam(path, "product") {
		return []map[string]string{{"product": "stream"}, {"product": "edge"}}
	}
	return []map[string]string{{}}
}

func identifiersFromRawItems(items []json.RawMessage, scope map[string]string, e registry.Entry) ([]map[string]string, error) {
	out := make([]map[string]string, 0, len(items))
	for _, raw := range items {
		item, err := rawMap(raw)
		if err != nil {
			return nil, err
		}
		if rawString(item, "lib", "library") == custom.EventBreakerLibCribl {
			continue
		}
		id := rawString(item, "id", "ID", "Id", "keyID", "keyId", "key_id", "name", "Name")
		if id == "" && e.ListUseGroupIDAsItemID && scope["group_id"] != "" {
			id = scope["group_id"]
		}
		if id == "" {
			continue
		}
		if shouldSkipRawID(e.TypeName, id, item) {
			continue
		}
		if e.TypeName == "criblio_pack_pipeline" && strings.HasPrefix(id, "pack:") {
			id = strings.TrimPrefix(id, "pack:")
		}
		m := map[string]string{"id": id}
		for _, key := range []string{"group_id", "lake_id", "pack", "product"} {
			if scope[key] != "" {
				m[key] = scope[key]
			}
		}
		out = append(out, m)
	}
	return out, nil
}

func shouldSkipRawID(typeName, id string, item map[string]any) bool {
	switch {
	case isLookupFileType(typeName) && (strings.HasPrefix(id, "cribl.") || rawHasCriblDefaultTag(item)):
		return true
	case typeName == "criblio_search_dataset" && rawHasCriblDefaultTag(item):
		return true
	case typeName == "criblio_search_dataset_provider" && custom.DefaultSearchDatasetProviderIDs[id]:
		return true
	case typeName == "criblio_search_dataset" && rawString(item, "type") == custom.SearchDatasetTypeCriblLake:
		return true
	case typeName == "criblio_cribl_lake_dataset" && custom.DefaultCriblLakeDatasetIDs[id]:
		return true
	case (typeName == "criblio_event_breaker_ruleset" || typeName == "criblio_search_macro") && strings.Contains(id, "."):
		return true
	case typeName == "criblio_pipeline" && strings.HasPrefix(id, "pack:"):
		return true
	case (typeName == "criblio_destination" || typeName == "criblio_pack_destination") && custom.DefaultDestinationIDs[id]:
		return true
	}
	return false
}

func isLookupFileType(typeName string) bool {
	return typeName == "criblio_lookup_file" || typeName == "criblio_pack_lookups"
}

func listLakehouseDatasetConnectionIdentifiers(ctx context.Context, client *importclient.Client) ([]map[string]string, error) {
	lakehouses, err := getRESTItems(ctx, client, "/products/lake/lakes/default/lakehouses")
	if restclient.IsNotFound(err) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	datasets, err := getRESTItems(ctx, client, "/products/lake/lakes/default/datasets")
	if restclient.IsNotFound(err) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	out := make([]map[string]string, 0, len(lakehouses)*len(datasets))
	for _, lhRaw := range lakehouses {
		lh, err := rawMap(lhRaw)
		if err != nil {
			return nil, err
		}
		lakehouseID := rawString(lh, "id", "ID", "name")
		if lakehouseID == "" {
			continue
		}
		for _, dsRaw := range datasets {
			ds, err := rawMap(dsRaw)
			if err != nil {
				return nil, err
			}
			datasetID := rawString(ds, "id", "ID", "name")
			if datasetID == "" {
				continue
			}
			out = append(out, map[string]string{"lakehouse_id": lakehouseID, "lake_dataset_id": datasetID})
		}
	}
	return out, nil
}

func listPackRoutesIdentifiers(ctx context.Context, client *importclient.Client, groupIDs []string) ([]map[string]string, error) {
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

func getPackIDsByGroup(ctx context.Context, client *importclient.Client, groupID string) ([]string, error) {
	items, err := getRESTItems(ctx, client, fmt.Sprintf("/m/%s/packs", url.PathEscape(groupID)))
	if restclient.IsNotFound(err) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	ids := make([]string, 0, len(items))
	for _, raw := range items {
		item, err := rawMap(raw)
		if err != nil {
			return nil, err
		}
		id := rawString(item, "id", "ID", "name")
		if id != "" && !custom.SkipPacks[id] {
			ids = append(ids, id)
		}
	}
	return ids, nil
}

func getRESTItems(ctx context.Context, client *importclient.Client, path string) ([]json.RawMessage, error) {
	items, err := restclient.Get[[]json.RawMessage](ctx, client.REST, path)
	if err == nil {
		if items == nil {
			return nil, nil
		}
		return *items, nil
	}
	var notFound *restclient.NotFoundError
	if errors.As(err, &notFound) {
		return nil, err
	}
	item, itemErr := restclient.Get[map[string]json.RawMessage](ctx, client.REST, path)
	if itemErr != nil {
		return nil, err
	}
	raw, err := json.Marshal(item)
	if err != nil {
		return nil, err
	}
	return []json.RawMessage{raw}, nil
}

func rawMap(raw json.RawMessage) (map[string]any, error) {
	var item map[string]any
	if err := json.Unmarshal(raw, &item); err != nil {
		return nil, err
	}
	return item, nil
}

func rawString(item map[string]any, keys ...string) string {
	for _, key := range keys {
		value, ok := item[key]
		if !ok || value == nil {
			continue
		}
		switch typed := value.(type) {
		case string:
			return strings.TrimSpace(typed)
		case fmt.Stringer:
			return strings.TrimSpace(typed.String())
		default:
			if s := fmt.Sprint(typed); s != "" && s != "<nil>" {
				return strings.TrimSpace(s)
			}
		}
	}
	return ""
}

func rawHasCriblDefaultTag(item map[string]any) bool {
	tags, ok := item["tags"]
	if !ok {
		return false
	}
	switch typed := tags.(type) {
	case string:
		return strings.EqualFold(strings.TrimSpace(typed), custom.CriblDefaultTag)
	case []any:
		for _, tag := range typed {
			if s, ok := tag.(string); ok && strings.EqualFold(strings.TrimSpace(s), custom.CriblDefaultTag) {
				return true
			}
		}
	}
	return false
}

func pathUsesRESTParam(path, name string) bool {
	return strings.Contains(path, "{"+name+"}")
}

func renderRESTPath(path string, values map[string]string) string {
	rendered := path
	for key, value := range values {
		rendered = strings.ReplaceAll(rendered, "{"+key+"}", url.PathEscape(value))
	}
	return rendered
}

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

func sliceToSet(values []string) map[string]struct{} {
	set := make(map[string]struct{}, len(values))
	for _, value := range values {
		set[value] = struct{}{}
	}
	return set
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

func inSet(value string, set map[string]struct{}) bool {
	_, ok := set[value]
	return ok
}

func skipDiscoveryForGroupFilter(typeName string, groupIDs []string) bool {
	hasDefaultSearch := slices.Contains(groupIDs, "default_search")
	if strings.HasPrefix(typeName, "criblio_search_") {
		return !hasDefaultSearch
	}
	switch typeName {
	case "criblio_cribl_lake_dataset", "criblio_cribl_lake_house", "criblio_lakehouse_dataset_connection", "criblio_notification_target":
		return true
	}
	return false
}

func skipGroupScopedSingleton(typeName, groupID string) bool {
	return typeName == "criblio_routes" && (groupID == "default_search" || groupID == "search")
}

func filterOutDefaultSearch(groupIDs []string) []string {
	out := make([]string, 0, len(groupIDs))
	for _, gid := range groupIDs {
		if gid != "default_search" {
			out = append(out, gid)
		}
	}
	return out
}

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

func fallbackGroupIDs(onPrem bool) []string {
	ids := []string{"default"}
	if !onPrem {
		ids = append(ids, "default_search")
	}
	if workspace := os.Getenv("CRIBL_WORKSPACE_ID"); workspace != "" && workspace != "default" && workspace != "default_search" {
		ids = append(ids, workspace)
	}
	return ids
}
