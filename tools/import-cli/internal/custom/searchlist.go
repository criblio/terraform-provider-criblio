// Package custom: search list capture and parsing for criblio_search_dataset and
// criblio_search_dataset_provider when the SDK cannot unmarshal (e.g. cribl_lake).
// Provides HTTP capture transport, identifier parsing, and filtering so only user-created
// resources are imported: skip items with tag cribl:default,
// and for dataset providers skip known default IDs (cribl_leader, cribl_edge, S3, etc.).
package custom

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"sync"
)

// eventBreakerRulesetListPathRe matches m/{groupId}/lib/breakers at end of path (list, not get-by-id or pack).
// Extracts groupID from paths like /api/v1/m/default/lib/breakers or m/default/lib/breakers.
var eventBreakerRulesetListPathRe = regexp.MustCompile(`m/([^/]+)/lib/breakers$`)

const (
	PathSearchDatasets         = "/m/default_search/search/datasets"
	PathSearchDatasetProviders = "/m/default_search/search/dataset-providers"
)

// CriblDefaultTag is the tag used for built-in Cribl datasets; we skip these so only
// user-created datasets (tags != cribl:default) are imported.
const CriblDefaultTag = "cribl:default"

// EventBreakerLibCribl is the lib value for built-in event breaker rulesets. We ignore these;
// only user-created rulesets (lib=custom or cribl-custom) are imported.
const EventBreakerLibCribl = "cribl"

// SkipPacks lists pack names to never discover or export (e.g. HelloPacks used in examples/tests).
var SkipPacks = map[string]bool{
	"HelloPacks": true,
}

// DefaultSearchDatasetProviderIDs are built-in dataset provider IDs that are never imported.
// We skip these when parsing the dataset-providers list so only user-created providers are imported.
var DefaultSearchDatasetProviderIDs = map[string]bool{
	"cribl_leader":            true,
	"cribl_edge":              true,
	"S3":                      true,
	"cribl_s3sample_provider": true,
	"cribl_meta":              true,
	"cribl_lake":              true,
}

// DefaultSearchDatasetIDs are built-in search dataset IDs that are never imported.
// We skip these when parsing the datasets list so only user-created datasets are imported.
var DefaultSearchDatasetIDs = map[string]bool{
	"cribl_logs":       true,
	"cribl_metrics":    true,
	"default_events":   true,
	"default_logs":     true,
	"default_metrics":  true,
	"default_spans":    true,
	"S3":               true, // default S3 dataset
}

// DefaultDestinationIDs are built-in destination IDs (default, devnull) that are never imported.
// We skip these when listing destinations so only user-created destinations are exported.
var DefaultDestinationIDs = map[string]bool{
	"default": true,
	"devnull": true,
}

// DefaultCriblLakeDatasetIDs are built-in Cribl Lake dataset IDs that are never imported.
// We skip these when listing lake datasets so only user-created lake datasets are exported.
var DefaultCriblLakeDatasetIDs = map[string]bool{
	"cribl_logs":       true,
	"cribl_metrics":    true,
	"default_events":   true,
	"default_logs":     true,
	"default_metrics":  true,
	"default_spans":    true,
	"S3":               true,
}

// SearchDatasetTypeCriblLake is the API type for Cribl Lake datasets. They appear in the search
// dataset list but are managed via /products/lake/lakes/{lakeId}/datasets and should be imported
// as criblio_cribl_lake_dataset, not criblio_search_dataset.
const SearchDatasetTypeCriblLake = "cribl_lake"

var (
	mu       sync.Mutex
	captured = make(map[string][]byte)
)

// SearchListTransport wraps a RoundTripper and copies response bodies for search list
// URLs so they can be parsed when the SDK returns a union unmarshal error.
type SearchListTransport struct {
	Base http.RoundTripper
}

// eventBreakerRulesetListKey builds the capture key for GET /m/{groupId}/lib/breakers list response.
// Returns "" if path is not the event breaker ruleset list URL. SDK ListEventBreakerRulesetResponseBody has no Items.
func eventBreakerRulesetListKey(path string) string {
	path = strings.Trim(path, "/")
	path = strings.TrimSuffix(path, "/")
	// Exclude pack paths (/m/{groupId}/p/{pack}/lib/breakers)
	if strings.Contains(path, "/p/") {
		return ""
	}
	// Use regex to match m/{groupId}/lib/breakers at end (list has no trailing id; get-by-id would be lib/breakers/{id})
	if m := eventBreakerRulesetListPathRe.FindStringSubmatch(path); len(m) >= 2 {
		return "event_breaker_ruleset_list:" + m[1]
	}
	return ""
}

// packBreakersListPathRe matches m/{groupId}/p/{pack}/lib/breakers at end of path.
var packBreakersListPathRe = regexp.MustCompile(`m/([^/]+)/p/([^/]+)/lib/breakers$`)

// packBreakersListKey builds the capture key for GET /m/{groupId}/p/{pack}/lib/breakers list response.
// Used when SDK fails to unmarshal (lib="cribl" not in EventBreakerRuleset Library enum).
func packBreakersListKey(path string) string {
	path = strings.Trim(path, "/")
	path = strings.TrimSuffix(path, "/")
	if m := packBreakersListPathRe.FindStringSubmatch(path); len(m) >= 3 {
		return "pack_breakers:" + m[1] + ":" + m[2]
	}
	return ""
}

// packOutputsListKey builds the capture key for GET /m/{groupId}/p/{pack}/system/outputs list response.
// Returns "" if path is not a pack outputs list URL.
func packOutputsListKey(path string) string {
	if !strings.HasSuffix(path, "/system/outputs") || strings.HasPrefix(path, "/system/outputs") {
		return ""
	}
	// Path like /api/v1/m/default/p/my-pack/system/outputs (list has no trailing id).
	parts := strings.Split(strings.Trim(path, "/"), "/")
	var groupID, pack string
	for i := 0; i < len(parts)-1; i++ {
		if parts[i] == "m" && i+1 < len(parts) {
			groupID = parts[i+1]
		}
		if parts[i] == "p" && i+1 < len(parts) {
			pack = parts[i+1]
			break
		}
	}
	if groupID == "" || pack == "" {
		return ""
	}
	return "pack_outputs:" + groupID + ":" + pack
}

// packInputsListKey builds the capture key for GET /m/{groupId}/p/{pack}/system/inputs list response.
// Returns "" if path is not a pack inputs list URL. Used when SDK fails to unmarshal Input union (all fields null).
func packInputsListKey(path string) string {
	if !strings.HasSuffix(path, "/system/inputs") || strings.HasPrefix(path, "/system/inputs") {
		return ""
	}
	if !strings.Contains(path, "/p/") {
		return ""
	}
	// Path like /api/v1/m/default/p/my-pack/system/inputs (list has no trailing id).
	parts := strings.Split(strings.Trim(path, "/"), "/")
	var groupID, pack string
	for i := 0; i < len(parts)-1; i++ {
		if parts[i] == "m" && i+1 < len(parts) {
			groupID = parts[i+1]
		}
		if parts[i] == "p" && i+1 < len(parts) {
			pack = parts[i+1]
			break
		}
	}
	if groupID == "" || pack == "" {
		return ""
	}
	return "pack_inputs:" + groupID + ":" + pack
}

// savedJobsListKey builds the capture key for GET /m/{groupId}/lib/jobs (saved jobs / collectors list).
// Returns "" if path is not the saved jobs list URL (excludes get-by-id paths with trailing /{id}).
func savedJobsListKey(path string) string {
	path = strings.Trim(path, "/")
	if !strings.HasSuffix(path, "/lib/jobs") {
		return ""
	}
	parts := strings.Split(path, "/")
	for i, p := range parts {
		if p == "m" && i+1 < len(parts) {
			return "saved_jobs_list:" + parts[i+1]
		}
	}
	return ""
}

// notificationTargetListKey builds the capture key for GET /notification-targets list.
// Returns "" if path is not the notification target list URL (excludes get-by-id paths).
func notificationTargetListKey(path string) string {
	path = strings.Trim(path, "/")
	if strings.HasSuffix(path, "/notification-targets") {
		return "notification_targets_list"
	}
	return ""
}

// packOutputGetKey builds the capture key for GET /m/{groupId}/p/{pack}/system/outputs/{id} (single pack output).
// Returns "" if path is not a pack output get URL.
func packOutputGetKey(path string) string {
	if !strings.Contains(path, "/system/outputs/") || !strings.Contains(path, "/p/") {
		return ""
	}
	parts := strings.Split(strings.Trim(path, "/"), "/")
	var groupID, pack, id string
	for i := 0; i < len(parts); i++ {
		if parts[i] == "m" && i+1 < len(parts) {
			groupID = parts[i+1]
		}
		if parts[i] == "p" && i+1 < len(parts) {
			pack = parts[i+1]
		}
		if parts[i] == "outputs" && i+1 < len(parts) {
			id = parts[i+1]
			break
		}
	}
	if groupID == "" || pack == "" || id == "" {
		return ""
	}
	return "pack_output_get:" + groupID + ":" + pack + ":" + id
}

// RoundTrip implements http.RoundTripper.
func (t *SearchListTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	base := t.Base
	if base == nil {
		base = http.DefaultTransport
	}
	res, err := base.RoundTrip(req)
	if err != nil || res == nil || res.Body == nil {
		return res, err
	}
	// Use Path (decoded); fallback to EscapedPath if Path is empty (some URL formats)
	path := req.URL.Path
	if path == "" && req.URL.RawPath != "" {
		if p, err := url.PathUnescape(req.URL.RawPath); err == nil {
			path = p
		}
	}
	if req.Method != http.MethodGet {
		return res, nil
	}
	var key string
	switch {
	case strings.HasSuffix(path, PathSearchDatasets) && path != "":
		key = PathSearchDatasets
	case strings.HasSuffix(path, PathSearchDatasetProviders) && path != "":
		key = PathSearchDatasetProviders
	default:
		if k := packOutputGetKey(path); k != "" {
			key = k
		} else if k := packOutputsListKey(path); k != "" {
			key = k
		} else if k := packInputsListKey(path); k != "" {
			key = k
		} else if k := eventBreakerRulesetListKey(path); k != "" {
			key = k
		} else if k := packBreakersListKey(path); k != "" {
			key = k
		} else if k := savedJobsListKey(path); k != "" {
			key = k
		} else if k := notificationTargetListKey(path); k != "" {
			key = k
		} else if strings.Contains(path, "/search/datasets/") {
			// Capture GET .../search/datasets/{id} response so we can use full body when SDK unmarshal fails.
			parts := strings.Split(strings.Trim(path, "/"), "/")
			for i, p := range parts {
				if p == "datasets" && i+1 < len(parts) {
					id := parts[i+1]
					if id != "" {
						key = "dataset_get:" + id
					}
					break
				}
			}
		}
		if key == "" {
			return res, nil
		}
	}
	body, err := io.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		return res, nil
	}
	mu.Lock()
	captured[key] = body
	mu.Unlock()
	res.Body = io.NopCloser(strings.NewReader(string(body)))
	return res, nil
}

// GetAndClearSearchListBody returns the captured body for the given path key and removes it.
func GetAndClearSearchListBody(key string) []byte {
	mu.Lock()
	defer mu.Unlock()
	b := captured[key]
	delete(captured, key)
	return b
}

// GetAndClearSearchDatasetGetBody returns the captured response body for GET .../search/datasets/{id} and removes it.
// The body is the raw API response (e.g. {"items":[{...}]}). Use when SDK Get succeeded but unmarshal failed.
func GetAndClearSearchDatasetGetBody(id string) []byte {
	return GetAndClearSearchListBody("dataset_get:" + id)
}

// GetAndClearPackOutputsListBody returns the captured response body for GET /m/{groupID}/p/{pack}/system/outputs
// and removes it. The SDK ListPackOutput response body has no Items in the generated type; we parse the raw JSON here.
func GetAndClearPackOutputsListBody(groupID, pack string) []byte {
	return GetAndClearSearchListBody("pack_outputs:" + groupID + ":" + pack)
}

// GetAndClearPackInputsListBody returns the captured response body for GET /m/{groupID}/p/{pack}/system/inputs
// and removes it. Used when SDK fails to unmarshal Input union (GetSystemInputsByPack returns empty items or union error).
func GetAndClearPackInputsListBody(groupID, pack string) []byte {
	return GetAndClearSearchListBody("pack_inputs:" + groupID + ":" + pack)
}

// GetAndClearPackOutputGetBody returns the captured response body for GET /m/{groupID}/p/{pack}/system/outputs/{id}
// and removes it. The SDK GetPackOutputByIDResponseBody is empty; we capture the raw response and use GetOutputByIDResponseBody shape for conversion.
func GetAndClearPackOutputGetBody(groupID, pack, id string) []byte {
	return GetAndClearSearchListBody("pack_output_get:" + groupID + ":" + pack + ":" + id)
}

// packOutputFirstItemCache holds raw items[0] JSON for pack_destination so export can emit the oneOf block
// when the converted model (DestinationResourceModel) has no Items field and addOneOfBlockFromFirstItem does nothing.
var (
	packOutputFirstItemMu    sync.Mutex
	packOutputFirstItemCache = make(map[string][]byte)
)

func packOutputFirstItemKey(groupID, pack, id string) string {
	return "pack_output_first_item:" + groupID + ":" + pack + ":" + id
}

// StorePackOutputFirstItem stores the raw JSON of items[0] from the pack output GET response.
// Used so export can emit the correct oneOf block (e.g. output_cribl_lake) when the model does not expose Items.
func StorePackOutputFirstItem(groupID, pack, id string, itemJSON []byte) {
	if len(itemJSON) == 0 {
		return
	}
	packOutputFirstItemMu.Lock()
	packOutputFirstItemCache[packOutputFirstItemKey(groupID, pack, id)] = itemJSON
	packOutputFirstItemMu.Unlock()
}

// GetAndClearPackOutputFirstItem returns and removes the stored first item JSON for the given pack output.
func GetAndClearPackOutputFirstItem(groupID, pack, id string) ([]byte, bool) {
	packOutputFirstItemMu.Lock()
	defer packOutputFirstItemMu.Unlock()
	key := packOutputFirstItemKey(groupID, pack, id)
	b, ok := packOutputFirstItemCache[key]
	if ok {
		delete(packOutputFirstItemCache, key)
	}
	return b, ok
}

// GetAndClearSavedJobsListBody returns the captured response body for GET /m/{groupID}/lib/jobs
// and removes it. Used when the SDK fails to unmarshal scheduledSearch/executor items into InputCollector.
func GetAndClearSavedJobsListBody(groupID string) []byte {
	return GetAndClearSearchListBody("saved_jobs_list:" + groupID)
}

// UnsupportedSavedJobTypes are SavedJob top-level type values that do not map to any InputCollector variant.
var UnsupportedSavedJobTypes = map[string]bool{
	"scheduledSearch": true,
	"executor":        true,
}

// ParseSavedJobsListBody parses the raw saved jobs list response and returns identifiers
// for items whose type is supported by the provider (i.e. "collection"), skipping
// scheduledSearch and executor types that cannot be unmarshaled into InputCollector.
func ParseSavedJobsListBody(body []byte, groupID string) ([]map[string]string, error) {
	var resp struct {
		Items []json.RawMessage `json:"items"`
	}
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, err
	}
	out := make([]map[string]string, 0, len(resp.Items))
	for _, raw := range resp.Items {
		var item struct {
			ID   string `json:"id"`
			Type string `json:"type"`
		}
		if err := json.Unmarshal(raw, &item); err != nil || item.ID == "" {
			continue
		}
		if UnsupportedSavedJobTypes[item.Type] {
			continue
		}
		out = append(out, map[string]string{"id": item.ID, "group_id": groupID})
	}
	return out, nil
}

// GetAndClearNotificationTargetListBody returns the captured response body for GET /notification-targets
// and removes it. Used when the SDK fails to unmarshal bulletin_message items into NotificationTarget.
func GetAndClearNotificationTargetListBody() []byte {
	return GetAndClearSearchListBody("notification_targets_list")
}

// UnsupportedNotificationTargetTypes are notification target type values not in the SDK NotificationTarget union.
var UnsupportedNotificationTargetTypes = map[string]bool{
	"bulletin_message": true,
}

// ParseNotificationTargetListBody parses the raw notification target list response and returns
// identifiers for items whose type is supported by the provider, skipping bulletin_message.
// Notification targets are not group-scoped so identifiers only contain "id".
func ParseNotificationTargetListBody(body []byte) ([]map[string]string, error) {
	var resp struct {
		Items []json.RawMessage `json:"items"`
	}
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, err
	}
	out := make([]map[string]string, 0, len(resp.Items))
	for _, raw := range resp.Items {
		var item struct {
			ID   string `json:"id"`
			Type string `json:"type"`
		}
		if err := json.Unmarshal(raw, &item); err != nil || item.ID == "" {
			continue
		}
		if UnsupportedNotificationTargetTypes[item.Type] {
			continue
		}
		out = append(out, map[string]string{"id": item.ID})
	}
	return out, nil
}

// GetAndClearEventBreakerRulesetListBody returns the captured response body for GET /m/{groupID}/lib/breakers
// and removes it. SDK ListEventBreakerRulesetResponseBody has no Items; we parse the raw JSON.
func GetAndClearEventBreakerRulesetListBody(groupID string) []byte {
	return GetAndClearSearchListBody("event_breaker_ruleset_list:" + groupID)
}

// GetAndClearPackBreakersListBody returns the captured response body for GET /m/{groupID}/p/{pack}/lib/breakers
// and removes it. Used when SDK fails to unmarshal (lib="cribl" not in EventBreakerRuleset Library enum).
func GetAndClearPackBreakersListBody(groupID, packID string) []byte {
	return GetAndClearSearchListBody("pack_breakers:" + groupID + ":" + packID)
}

// ParsePackBreakersListBody parses the pack breakers list response ({"items":[{"id":"...", "lib":"..."},...]})
// and returns identifier maps with group_id, id, pack. Skips items with lib="cribl" (built-in).
func ParsePackBreakersListBody(body []byte, groupID, packID string) ([]map[string]string, error) {
	var resp struct {
		Items []json.RawMessage `json:"items"`
	}
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, err
	}
	out := make([]map[string]string, 0, len(resp.Items))
	for _, raw := range resp.Items {
		var item struct {
			ID  string  `json:"id"`
			Lib *string `json:"lib"`
		}
		if err := json.Unmarshal(raw, &item); err != nil || item.ID == "" {
			continue
		}
		if item.Lib != nil && *item.Lib == EventBreakerLibCribl {
			continue
		}
		out = append(out, map[string]string{"group_id": groupID, "id": item.ID, "pack": packID})
	}
	return out, nil
}

// ParseEventBreakerRulesetListBody parses the list event breaker ruleset API response ({"items":[{"id":"...", ...},...]})
// and returns one identifier map per item with "id" and "group_id" set.
// Uses raw JSON parsing for items so we don't fail when lib has values like "cribl" (built-in) that aren't in the SDK enum.
// Skips built-in rulesets (lib=cribl); only exports user-created.
func ParseEventBreakerRulesetListBody(body []byte, groupID string) ([]map[string]string, error) {
	var resp struct {
		Items []json.RawMessage `json:"items"`
	}
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, err
	}
	out := make([]map[string]string, 0, len(resp.Items))
	for _, raw := range resp.Items {
		var item struct {
			ID  string  `json:"id"`
			Lib *string `json:"lib"`
		}
		if err := json.Unmarshal(raw, &item); err != nil || item.ID == "" {
			continue
		}
		if item.Lib != nil && *item.Lib == EventBreakerLibCribl {
			continue
		}
		out = append(out, map[string]string{"id": item.ID, "group_id": groupID})
	}
	return out, nil
}

// ParsePackOutputsListBody parses the list pack outputs API response ({"items":[{"id":"..."},...]}) and returns
// one identifier map per item with "id" set. Used when SDK returns empty items (ListPackOutputResponseBody has no Items).
func ParsePackOutputsListBody(body []byte) ([]map[string]string, error) {
	var resp struct {
		Items []struct {
			ID string `json:"id"`
		} `json:"items"`
	}
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, err
	}
	out := make([]map[string]string, 0, len(resp.Items))
	for _, it := range resp.Items {
		if it.ID != "" {
			out = append(out, map[string]string{"id": it.ID})
		}
	}
	return out, nil
}

// ParsePackInputsListBody parses the list pack inputs API response ({"items":[{"id":"..."},...]}) and returns
// one identifier map per item with "id" set. Used when SDK fails to unmarshal Input union (empty items or union error).
func ParsePackInputsListBody(body []byte) ([]map[string]string, error) {
	var resp struct {
		Items []struct {
			ID string `json:"id"`
		} `json:"items"`
	}
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, err
	}
	out := make([]map[string]string, 0, len(resp.Items))
	for _, it := range resp.Items {
		if it.ID != "" {
			out = append(out, map[string]string{"id": it.ID})
		}
	}
	return out, nil
}

func itemHasCriblDefaultTag(raw json.RawMessage) bool {
	var it struct {
		Tags interface{} `json:"tags,omitempty"`
	}
	if err := json.Unmarshal(raw, &it); err != nil || it.Tags == nil {
		return false
	}
	switch v := it.Tags.(type) {
	case string:
		return v == CriblDefaultTag
	case []interface{}:
		for _, e := range v {
			if s, ok := e.(string); ok && s == CriblDefaultTag {
				return true
			}
		}
		return false
	}
	return false
}

// itemIsCriblLakeDataset reports whether the item has type "cribl_lake". Such datasets
// are Cribl Lake datasets (products/lake API) and should be exported as criblio_cribl_lake_dataset.
func itemIsCriblLakeDataset(raw json.RawMessage) bool {
	var it struct {
		Type string `json:"type,omitempty"`
	}
	if err := json.Unmarshal(raw, &it); err != nil {
		return false
	}
	return it.Type == SearchDatasetTypeCriblLake
}

var itemCache map[string]map[string]json.RawMessage // path key -> id -> item JSON

// IdentifiersFromSearchListBody parses body as { "items": [ { "id": "..." }, ... ] }
// and returns one map per item with key "id". pathKey is PathSearchDatasets or
// PathSearchDatasetProviders; each full item is cached by id so Convert can build
// the model from the list item (reuse list response, no Get call).
// Search datasets: skip items whose tags include "cribl:default"; skip type "cribl_lake" (those
// are exported as criblio_cribl_lake_dataset only). All other items are criblio_search_dataset.
func IdentifiersFromSearchListBody(body []byte, pathKey string) ([]map[string]string, int, error) {
	if len(body) == 0 {
		return nil, 0, nil
	}
	var resp struct {
		Items []json.RawMessage `json:"items,omitempty"`
	}
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, 0, err
	}
	mu.Lock()
	if itemCache == nil {
		itemCache = make(map[string]map[string]json.RawMessage)
	}
	c := itemCache[pathKey]
	if c == nil {
		c = make(map[string]json.RawMessage)
		itemCache[pathKey] = c
	}
	mu.Unlock()

	out := make([]map[string]string, 0, len(resp.Items))
	for _, raw := range resp.Items {
		var it struct {
			ID string `json:"id"`
		}
		if err := json.Unmarshal(raw, &it); err != nil || it.ID == "" {
			continue
		}
		// Skip search datasets/providers tagged cribl:default (built-in); only user-created are imported.
		if itemHasCriblDefaultTag(raw) {
			continue
		}
		if pathKey == PathSearchDatasetProviders && DefaultSearchDatasetProviderIDs[it.ID] {
			continue
		}
		// Skip cribl_lake datasets: they are managed via Lake API and exported as criblio_cribl_lake_dataset only.
		if pathKey == PathSearchDatasets && itemIsCriblLakeDataset(raw) {
			continue
		}
		out = append(out, map[string]string{"id": it.ID})
		c[it.ID] = raw
	}
	mu.Lock()
	itemCache[pathKey] = c
	mu.Unlock()
	return out, len(out), nil
}

// GetCachedSearchListItem returns the full item JSON for the given path key and id,
// if cached from a prior list parse. Used by the converter to build the model without calling Get.
func GetCachedSearchListItem(pathKey, id string) []byte {
	mu.Lock()
	defer mu.Unlock()
	if itemCache == nil {
		return nil
	}
	c := itemCache[pathKey]
	if c == nil {
		return nil
	}
	raw, ok := c[id]
	if !ok {
		return nil
	}
	return raw
}
