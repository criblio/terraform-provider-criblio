package provider

import (
	"context"
	"fmt"
	"regexp"
	"sort"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Query labels are provider-managed: the author gives a label-less `queries` list (+ optional
// `alert_expression`); the hooks key it A, B, C..., build the one "Alert Query" expr, and reject direct query/expr edits.

const alertQueryLabel = "Alert Query"

// queryLabelSequence is the positional label alphabet: queries[0]->A, queries[1]->B, ... (max 26).
var queryLabelSequence = func() []string {
	s := make([]string, 26)
	for i := range s {
		s[i] = string(rune('A' + i))
	}
	return s
}()

// monitorQuerySynthesize builds the `query` map by re-keying the author's `queries` list onto A/B/C.
func monitorQuerySynthesize() []planmodifier.Map {
	return []planmodifier.Map{monitorQueryMapSynthesizer{}}
}

type monitorQueryMapSynthesizer struct{}

func (monitorQueryMapSynthesizer) Description(_ context.Context) string {
	return "Builds the provider-managed query map by assigning labels A, B, C... to the queries list."
}

func (m monitorQueryMapSynthesizer) MarkdownDescription(ctx context.Context) string {
	return m.Description(ctx)
}

func (m monitorQueryMapSynthesizer) PlanModifyMap(ctx context.Context, req planmodifier.MapRequest, resp *planmodifier.MapResponse) {
	if !req.ConfigValue.IsNull() {
		resp.Diagnostics.AddAttributeError(req.Path, "query is provider-managed",
			"Set `queries` (and optionally `alert_expression`) instead; the provider assigns the A/B/C labels.")
		return
	}

	var queries types.List
	resp.Diagnostics.Append(req.Plan.GetAttribute(ctx, path.Root("queries"), &queries)...)
	if resp.Diagnostics.HasError() {
		return
	}

	elemType := types.ObjectType{AttrTypes: MonitorQueryAttrTypes()}

	// No operands in config: keep the prior value on import/refresh, else require queries.
	if queries.IsNull() || queries.IsUnknown() {
		if !req.StateValue.IsNull() && !req.StateValue.IsUnknown() {
			resp.PlanValue = req.StateValue
		} else {
			resp.Diagnostics.AddAttributeError(path.Root("queries"), "queries is required",
				"Provide at least one query in `queries`; the provider builds the query map from it.")
		}
		return
	}

	elems := queries.Elements()
	if len(elems) > len(queryLabelSequence) {
		resp.Diagnostics.AddAttributeError(path.Root("queries"), "too many queries",
			"A monitor supports at most 26 queries (labels A-Z).")
		return
	}

	kv := make(map[string]attr.Value, len(elems))
	for i, e := range elems {
		kv[queryLabelSequence[i]] = e
	}
	mapVal, d := types.MapValue(elemType, kv)
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.PlanValue = mapVal
}

// monitorExprSynthesize builds the single "Alert Query" expr over the queries' labels.
func monitorExprSynthesize() []planmodifier.List {
	return []planmodifier.List{monitorExprListSynthesizer{}}
}

type monitorExprListSynthesizer struct{}

func (monitorExprListSynthesizer) Description(_ context.Context) string {
	return "Synthesizes the single auto-managed Alert Query expression over the queries' labels."
}

func (m monitorExprListSynthesizer) MarkdownDescription(ctx context.Context) string {
	return m.Description(ctx)
}

func (m monitorExprListSynthesizer) PlanModifyList(ctx context.Context, req planmodifier.ListRequest, resp *planmodifier.ListResponse) {
	if !req.ConfigValue.IsNull() {
		resp.Diagnostics.AddAttributeError(req.Path, "expr is provider-managed",
			"The Alert Query expression is built from `queries`/`alert_expression`; do not set expr.")
		return
	}

	exprElemType := types.ObjectType{AttrTypes: MonitorExprAttrTypes()}

	// Every monitor type fires on the one "Alert Query" expression (resolve.ts resolveDetector reads
	// expr for threshold/change/anomaly/forecast alike), so it is synthesized whenever operands exist.
	var queries types.List
	resp.Diagnostics.Append(req.Plan.GetAttribute(ctx, path.Root("queries"), &queries)...)
	if resp.Diagnostics.HasError() {
		return
	}
	if queries.IsNull() || queries.IsUnknown() {
		// import/refresh: keep the prior expr (query hook raises the fresh-create "required" error).
		if !req.StateValue.IsNull() && !req.StateValue.IsUnknown() {
			resp.PlanValue = req.StateValue
		}
		return
	}
	n := len(queries.Elements())
	if n == 0 || n > len(queryLabelSequence) {
		return
	}

	var alertExpr types.String
	resp.Diagnostics.Append(req.Plan.GetAttribute(ctx, path.Root("alert_expression"), &alertExpr)...)
	if resp.Diagnostics.HasError() {
		return
	}

	exprModel, dangling := buildAlertExpression(queryLabelSequence[:n], alertExpr.ValueString())
	if len(dangling) > 0 {
		resp.Diagnostics.AddAttributeError(path.Root("alert_expression"),
			"alert_expression references unknown query labels",
			fmt.Sprintf("Labels [%s] have no matching query; the monitor has %d queries (%s). Reference only those labels.",
				strings.Join(dangling, ", "), n, strings.Join(queryLabelSequence[:n], ", ")))
		return
	}

	exprObj, d := types.ObjectValueFrom(ctx, MonitorExprAttrTypes(), exprModel)
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}
	listVal, d := types.ListValue(exprElemType, []attr.Value{exprObj})
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.PlanValue = listVal
}

// buildAlertExpression builds the "Alert Query" expr; exprText wins when set, else a lone query is its
// label and multiples sum. Returns labels exprText names with no query so the caller can reject them.
func buildAlertExpression(labels []string, exprText string) (MonitorExprModel, []string) {
	known := make(map[string]struct{}, len(labels))
	for _, l := range labels {
		known[l] = struct{}{}
	}

	var dangling []string
	text := strings.TrimSpace(exprText)
	switch {
	case text != "":
		for _, ref := range extractQueryLabels(text) {
			if _, ok := known[ref]; !ok {
				dangling = append(dangling, ref)
			}
		}
	case len(labels) == 1:
		text = labels[0]
	default:
		text = strings.Join(labels, " + ")
	}

	refs := extractQueryLabels(text)
	labelVals := make([]attr.Value, len(refs))
	for i, l := range refs {
		labelVals[i] = types.StringValue(l)
	}
	queryLabels, _ := types.ListValue(types.StringType, labelVals)

	return MonitorExprModel{
		Label:       types.StringValue(alertQueryLabel),
		Text:        types.StringValue(text),
		QueryLabels: queryLabels,
		DatasetID:   types.StringNull(),
		Left:        types.StringNull(),
		Operation:   types.StringNull(),
		Right:       types.StringNull(),
		Scalar:      types.Float64Null(),
	}, dangling
}

var queryLabelRef = regexp.MustCompile(`\b([A-Z])\b`)

// extractQueryLabels returns the distinct single-capital labels in text, sorted — a port of queryCompiler.ts.
func extractQueryLabels(text string) []string {
	seen := map[string]struct{}{}
	for _, m := range queryLabelRef.FindAllStringSubmatch(text, -1) {
		seen[m[1]] = struct{}{}
	}
	out := make([]string, 0, len(seen))
	for l := range seen {
		out = append(out, l)
	}
	sort.Strings(out)
	return out
}
