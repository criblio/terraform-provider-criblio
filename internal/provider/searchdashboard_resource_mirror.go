// Re-applies API response fields after refreshPlan merges the Terraform plan onto the refreshed
// model. refreshPlan overwrites Optional+Computed attributes omitted in HCL with null, which
// diverges from what Read returns (pure API) and causes perpetual plan diffs.

package provider

import (
	"encoding/json"
	"strings"

	tfTypes "github.com/criblio/terraform-provider-criblio/internal/provider/types"
	"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func restoreSearchDashboardAfterPlanMerge(dst, api *SearchDashboardResourceModel) {
	if dst == nil || api == nil {
		return
	}
	// Replace unknown/null category with API (including explicit null); refreshPlan often leaves unknown.
	if stringAttrEmpty(dst.Category) {
		dst.Category = api.Category
	}
	if len(dst.Groups) == 0 && len(api.Groups) > 0 {
		dst.Groups = api.Groups
	} else if len(api.Groups) == 0 {
		dst.Groups = nil
	}
	for i := range dst.Elements {
		if i >= len(api.Elements) {
			break
		}
		restoreDashboardElementUnion(&dst.Elements[i], &api.Elements[i])
	}
}

func stringAttrEmpty(s types.String) bool {
	return s.IsNull() || s.IsUnknown()
}

func restoreDashboardElementUnion(dst, api *tfTypes.DashboardElementUnion) {
	if dst == nil || api == nil {
		return
	}
	if dst.DashboardElementVisualization != nil && api.DashboardElementVisualization != nil {
		restoreVisualization(dst.DashboardElementVisualization, api.DashboardElementVisualization)
	}
	if dst.DashboardElementInput != nil && api.DashboardElementInput != nil {
		restoreInput(dst.DashboardElementInput, api.DashboardElementInput)
	}
	if dst.DashboardElement != nil && api.DashboardElement != nil {
		restoreMarkdownElement(dst.DashboardElement, api.DashboardElement)
	}
}

func restoreVisualization(dst, api *tfTypes.DashboardElementVisualization) {
	mergeConfigMaps(&dst.Config, api.Config)
	if dst.Search == nil && api.Search != nil {
		dst.Search = api.Search
	} else if dst.Search != nil && api.Search != nil {
		mergeSearchQuery(dst.Search, api.Search)
	}
	restoreTitleAction(&dst.TitleAction, api.TitleAction)
}

func restoreInput(dst, api *tfTypes.DashboardElementInput) {
	mergeConfigMaps(&dst.Config, api.Config)
	if dst.Search == nil && api.Search != nil {
		dst.Search = api.Search
	} else if dst.Search != nil && api.Search != nil {
		mergeSearchQuery(dst.Search, api.Search)
	}
	restoreTitleAction(&dst.TitleAction, api.TitleAction)
}

func restoreMarkdownElement(dst, api *tfTypes.DashboardElement) {
	mergeConfigMaps(&dst.Config, api.Config)
	if dst.Search == nil && api.Search != nil {
		dst.Search = api.Search
	}
	restoreTitleAction(&dst.TitleAction, api.TitleAction)
}

// elementConfigPresence records whether each element's config map had any keys in the prior
// Terraform state or plan snapshot. We only strip JSON-null-only keys when the user did not
// configure that map (no keys); if they set keys such as color = null, stripping would remove
// them from state and cause perpetual drift against configuration.
type elementConfigPresence struct {
	visHadKeys, inputHadKeys, mdHadKeys bool
}

func snapshotSearchDashboardConfigMapPresence(m *SearchDashboardResourceModel) []elementConfigPresence {
	if m == nil {
		return nil
	}
	out := make([]elementConfigPresence, len(m.Elements))
	for i := range m.Elements {
		el := &m.Elements[i]
		if el.DashboardElementVisualization != nil {
			out[i].visHadKeys = len(el.DashboardElementVisualization.Config) > 0
		}
		if el.DashboardElementInput != nil {
			out[i].inputHadKeys = len(el.DashboardElementInput.Config) > 0
		}
		if el.DashboardElement != nil {
			out[i].mdHadKeys = len(el.DashboardElement.Config) > 0
		}
	}
	return out
}

// normalizeSearchDashboardConfigMaps removes map entries whose JSON value is null when the map
// was not configured in Terraform (no keys in the given snapshot). The API often returns keys
// such as "color": null while RefreshFrom encodes them as Normalized("null"); stripping those
// avoids drift when config omits the map. When the snapshot has keys (e.g. color = null), we
// keep the API shape so state matches configuration.
func normalizeSearchDashboardConfigMaps(m *SearchDashboardResourceModel, presence []elementConfigPresence) {
	if m == nil {
		return
	}
	for i := range m.Elements {
		var p elementConfigPresence
		if i < len(presence) {
			p = presence[i]
		}
		el := &m.Elements[i]
		if el.DashboardElementVisualization != nil {
			if !p.visHadKeys {
				stripJSONNullKeysFromConfigMap(&el.DashboardElementVisualization.Config)
			}
		}
		if el.DashboardElementInput != nil {
			if !p.inputHadKeys {
				stripJSONNullKeysFromConfigMap(&el.DashboardElementInput.Config)
			}
		}
		if el.DashboardElement != nil {
			if !p.mdHadKeys {
				stripJSONNullKeysFromConfigMap(&el.DashboardElement.Config)
			}
		}
	}
	ensureSearchDashboardElementConfigMaps(m)
}

// ensureSearchDashboardElementConfigMaps sets each element branch's config map to a non-nil empty
// map when the API omitted it or normalization stripped all keys. Terraform otherwise stores null
// and the next plan shows a spurious "+ config = {}" drift against a known empty map.
func ensureSearchDashboardElementConfigMaps(m *SearchDashboardResourceModel) {
	if m == nil {
		return
	}
	for i := range m.Elements {
		el := &m.Elements[i]
		if el.DashboardElementVisualization != nil && el.DashboardElementVisualization.Config == nil {
			el.DashboardElementVisualization.Config = make(map[string]jsontypes.Normalized)
		}
		if el.DashboardElementInput != nil && el.DashboardElementInput.Config == nil {
			el.DashboardElementInput.Config = make(map[string]jsontypes.Normalized)
		}
		if el.DashboardElement != nil && el.DashboardElement.Config == nil {
			el.DashboardElement.Config = make(map[string]jsontypes.Normalized)
		}
	}
}

func stripJSONNullKeysFromConfigMap(m *map[string]jsontypes.Normalized) {
	if m == nil {
		return
	}
	if *m == nil {
		*m = make(map[string]jsontypes.Normalized)
		return
	}
	if len(*m) == 0 {
		*m = make(map[string]jsontypes.Normalized)
		return
	}
	for k, v := range *m {
		if isJSONNullNormalized(v) {
			delete(*m, k)
		}
	}
	if len(*m) == 0 {
		*m = make(map[string]jsontypes.Normalized)
	}
}

func isJSONNullNormalized(v jsontypes.Normalized) bool {
	if v.IsNull() || v.IsUnknown() {
		return true
	}
	s := strings.TrimSpace(v.ValueString())
	if s == "null" {
		return true
	}
	var any interface{}
	if err := json.Unmarshal([]byte(s), &any); err != nil {
		return false
	}
	return any == nil
}

func mergeConfigMaps(dst *map[string]jsontypes.Normalized, api map[string]jsontypes.Normalized) {
	if len(api) == 0 {
		return
	}
	if len(*dst) == 0 {
		m := make(map[string]jsontypes.Normalized, len(api))
		for k, v := range api {
			m[k] = v
		}
		*dst = m
		return
	}
	d := *dst
	for k, v := range api {
		if _, ok := d[k]; !ok {
			d[k] = v
		}
	}
}

func restoreTitleAction(dst **tfTypes.TitleAction, api *tfTypes.TitleAction) {
	if dst == nil {
		return
	}
	if api == nil {
		if *dst != nil && titleActionEmptyOrUnknown(*dst) {
			*dst = nil
		}
		return
	}
	if *dst == nil {
		t := *api
		*dst = &t
		return
	}
	mergeTitleAction(*dst, api)
}

func titleActionEmptyOrUnknown(t *tfTypes.TitleAction) bool {
	if t == nil {
		return true
	}
	return stringAttrEmpty(t.Label) && stringAttrEmpty(t.URL) &&
		(t.OpenInNewTab.IsNull() || t.OpenInNewTab.IsUnknown())
}

func mergeTitleAction(dst, api *tfTypes.TitleAction) {
	if dst == nil || api == nil {
		return
	}
	dst.Label = mergeTFString(dst.Label, api.Label)
	dst.URL = mergeTFString(dst.URL, api.URL)
	dst.OpenInNewTab = mergeTFBool(dst.OpenInNewTab, api.OpenInNewTab)
}

func mergeTFString(dst, api types.String) types.String {
	if dst.IsNull() || dst.IsUnknown() {
		return api
	}
	return dst
}

func mergeTFFloat64(dst, api types.Float64) types.Float64 {
	if dst.IsNull() || dst.IsUnknown() {
		return api
	}
	return dst
}

func mergeTFBool(dst, api types.Bool) types.Bool {
	if dst.IsNull() || dst.IsUnknown() {
		return api
	}
	return dst
}

func mergeSearchQuery(dst, api *tfTypes.SearchQuery) {
	if dst == nil || api == nil {
		return
	}
	if dst.SearchQueryInline == nil && api.SearchQueryInline != nil {
		inline := *api.SearchQueryInline
		dst.SearchQueryInline = &inline
	} else if dst.SearchQueryInline != nil && api.SearchQueryInline != nil {
		mergeSearchQueryInline(dst.SearchQueryInline, api.SearchQueryInline)
	}
	if dst.SearchQuerySaved == nil && api.SearchQuerySaved != nil {
		saved := *api.SearchQuerySaved
		dst.SearchQuerySaved = &saved
	} else if dst.SearchQuerySaved != nil && api.SearchQuerySaved != nil {
		mergeSearchQuerySaved(dst.SearchQuerySaved, api.SearchQuerySaved)
	}
	if dst.SearchQueryValues == nil && api.SearchQueryValues != nil {
		vals := *api.SearchQueryValues
		dst.SearchQueryValues = &vals
	} else if dst.SearchQueryValues != nil && api.SearchQueryValues != nil {
		mergeSearchQueryValues(dst.SearchQueryValues, api.SearchQueryValues)
	}
}

func mergeSearchQueryInline(dst, api *tfTypes.SearchQueryInline) {
	if dst == nil || api == nil {
		return
	}
	if dst.Earliest == nil && api.Earliest != nil {
		e := *api.Earliest
		dst.Earliest = &e
	} else if dst.Earliest != nil && api.Earliest != nil {
		dst.Earliest.Str = mergeTFString(dst.Earliest.Str, api.Earliest.Str)
		dst.Earliest.Number = mergeTFFloat64(dst.Earliest.Number, api.Earliest.Number)
	}
	if dst.Latest == nil && api.Latest != nil {
		l := *api.Latest
		dst.Latest = &l
	} else if dst.Latest != nil && api.Latest != nil {
		dst.Latest.Str = mergeTFString(dst.Latest.Str, api.Latest.Str)
		dst.Latest.Number = mergeTFFloat64(dst.Latest.Number, api.Latest.Number)
	}
	dst.ParentSearchID = mergeTFString(dst.ParentSearchID, api.ParentSearchID)
	dst.Query = mergeTFString(dst.Query, api.Query)
	dst.SampleRate = mergeTFFloat64(dst.SampleRate, api.SampleRate)
	dst.Timezone = mergeTFString(dst.Timezone, api.Timezone)
	dst.Type = mergeTFString(dst.Type, api.Type)
}

func mergeSearchQuerySaved(dst, api *tfTypes.SearchQuerySaved) {
	if dst == nil || api == nil {
		return
	}
	dst.Query = mergeTFString(dst.Query, api.Query)
	dst.QueryID = mergeTFString(dst.QueryID, api.QueryID)
	dst.RunMode = mergeTFString(dst.RunMode, api.RunMode)
	dst.Type = mergeTFString(dst.Type, api.Type)
}

func mergeSearchQueryValues(dst, api *tfTypes.SearchQueryValues) {
	if dst == nil || api == nil {
		return
	}
	dst.Type = mergeTFString(dst.Type, api.Type)
	if len(dst.Values) == 0 && len(api.Values) > 0 {
		dst.Values = api.Values
	}
}
