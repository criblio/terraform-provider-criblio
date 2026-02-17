// Package discovery implements the discovery engine that iterates through the
// registry and calls SDK List* endpoints to enumerate existing resources.
package discovery

import (
	"context"
	"fmt"
	"reflect"

	"github.com/criblio/terraform-provider-criblio/internal/sdk"
	"github.com/criblio/terraform-provider-criblio/internal/sdk/models/operations"
	"github.com/criblio/terraform-provider-criblio/tools/import-cli/internal/registry"
)

// Result holds the discovery result for one resource type: count and any error
// with resource context. Items are left for future HCL generation.
// Details is optional (e.g. group names for criblio_group) for dry-run display.
// PerGroupCounts is set for group-scoped resources (key = group label e.g. "default (stream)").
type Result struct {
	TypeName       string
	Count          int
	Err            error
	Details        []string         // optional; e.g. group names for criblio_group
	PerGroupCounts map[string]int   // optional; per-group count for preview/export
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
	groupIDs := make([]string, 0, len(streamIDs)+len(edgeIDs))
	groupIDs = append(groupIDs, streamIDs...)
	groupIDs = append(groupIDs, edgeIDs...)
	if len(groupIDs) == 0 && len(groupFilter) == 0 {
		groupIDs = []string{"default"}
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
		if e.SDKService == "" || e.ListMethod == "" {
			count := 0
			var details []string
			if e.TypeName == "criblio_group" {
				count = len(streamNames) + len(edgeNames)
				details = groupNames
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
	// Skip API call when request requires a pack ID we don't have (avoids 404 from /p// in path).
	if len(args) >= 2 && requestRequiresPack(args[1]) {
		return 0, nil, fmt.Errorf("skipped: list request requires pack ID (pack-scoped resource; use --import-ids-file to add import IDs)")
	}
	// If the method has no request or request has no GroupID, call once.
	if len(args) < 2 || !requestHasGroupID(args[1]) {
		n, err := callListAndCount(ctx, method, args)
		return n, nil, err
	}
	// Request has GroupID: call once per group and sum counts; record per-group.
	perGroup = make(map[string]int)
	for _, gid := range groupIDs {
		args := buildListArgs(ctx, method, e, gid)
		if len(args) >= 2 && requestRequiresPack(args[1]) {
			return total, perGroup, nil
		}
		n, err := callListAndCount(ctx, method, args)
		total += n
		perGroup[gid] = n
		if err != nil {
			return total, perGroup, err
		}
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
	f := req.FieldByName("GroupID")
	if f.IsValid() && f.CanSet() && f.Kind() == reflect.String {
		f.SetString(defaultGroupID)
	}
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
