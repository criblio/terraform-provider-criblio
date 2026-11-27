package provider

import (
	"context"
	"encoding/json"
	"strings"
	"unicode"

	custom_listplanmodifier "github.com/criblio/terraform-provider-criblio/internal/tfplanmodifiers/listplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type pipelineFunctionSettings struct{}

func pipelineFunctionsPlanModifiers() []planmodifier.List {
	return []planmodifier.List{
		custom_listplanmodifier.SuppressDiff(custom_listplanmodifier.ExplicitSuppress),
		pipelineFunctionSettings{},
	}
}

func (pipelineFunctionSettings) Description(context.Context) string {
	return "Preserves API settings from the matching pipeline function."
}

func (m pipelineFunctionSettings) MarkdownDescription(ctx context.Context) string {
	return m.Description(ctx)
}

func (pipelineFunctionSettings) PlanModifyList(ctx context.Context, req planmodifier.ListRequest, resp *planmodifier.ListResponse) {
	config, prior, planned := req.ConfigValue.Elements(), req.StateValue.Elements(), req.PlanValue.Elements()
	if len(prior) == 0 || len(config) == 0 || len(config) != len(planned) || req.PlanValue.Equal(req.StateValue) {
		return
	}
	matches := matchOrderedConfig(config, prior, "id", pipelineFunctionMatchesConfig)
	for i, j := range matches {
		if j < 0 || planned[i].IsUnknown() || planned[i].IsNull() {
			continue
		}
		function := planned[i].(types.Object)
		attrs := function.Attributes()
		old := prior[j].(types.Object).Attributes()
		for key, value := range config[i].(types.Object).Attributes() {
			if value.IsNull() || (key == "conf" && equivalentFunctionConf(objectString(config[i], "id"), value, old[key])) {
				attrs[key] = old[key]
			}
		}
		updated, diags := types.ObjectValue(function.AttributeTypes(ctx), attrs)
		resp.Diagnostics.Append(diags...)
		planned[i] = updated
	}
	value, diags := types.ListValue(req.PlanValue.ElementType(ctx), planned)
	resp.Diagnostics.Append(diags...)
	resp.PlanValue = value
}

func pipelineFunctionMatchesConfig(config, prior types.Object) bool {
	old := prior.Attributes()
	for key, value := range config.Attributes() {
		if value.IsNull() || value.IsUnknown() || value.Equal(old[key]) {
			continue
		}
		if key != "conf" || !equivalentFunctionConf(objectString(config, "id"), value, old[key]) {
			return false
		}
	}
	return true
}

func equivalentFunctionConf(id string, a, b attr.Value) bool {
	if a.IsUnknown() || a.IsNull() || b.IsUnknown() || b.IsNull() {
		return false
	}
	return normalizePipelineConf(id, a.(jsontypes.Normalized).ValueString()) == normalizePipelineConf(id, b.(jsontypes.Normalized).ValueString())
}

func normalizePipelineConf(functionID, value string) string {
	var decoded any
	if err := json.Unmarshal([]byte(value), &decoded); err != nil {
		return normalizePipelineCode(value)
	}
	// The eval and aggregation functions default printUndefineds to false.
	// Other function types and nested keys belong to the user's configuration.
	if obj, ok := decoded.(map[string]any); ok && (functionID == "eval" || functionID == "aggregation") && obj["printUndefineds"] == false {
		delete(obj, "printUndefineds")
	}
	normalized := normalizePipelineConfValue(decoded)
	encoded, err := json.Marshal(normalized)
	if err != nil {
		return normalizePipelineCode(value)
	}
	return string(encoded)
}

func normalizePipelineConfValue(value any) any {
	switch typed := value.(type) {
	case map[string]any:
		output := make(map[string]any, len(typed))
		for key, item := range typed {
			if key == "code" {
				if code, ok := item.(string); ok {
					output[key] = normalizePipelineCode(code)
					continue
				}
			}
			output[key] = normalizePipelineConfValue(item)
		}
		return output
	case []any:
		output := make([]any, len(typed))
		for index, item := range typed {
			output[index] = normalizePipelineConfValue(item)
		}
		return output
	default:
		return value
	}
}

func normalizePipelineCode(value string) string {
	value = strings.ReplaceAll(value, "\r\n", "\n")
	value = strings.ReplaceAll(value, "\r", "\n")
	lines := strings.Split(strings.TrimSpace(value), "\n")
	minIndent := -1
	for _, line := range lines {
		if strings.TrimSpace(line) == "" {
			continue
		}
		indent := leadingWhitespace(line)
		if minIndent == -1 || indent < minIndent {
			minIndent = indent
		}
	}
	if minIndent <= 0 {
		return strings.Join(lines, "\n")
	}
	for index, line := range lines {
		lines[index] = trimLeadingWhitespace(line, minIndent)
	}
	return strings.Join(lines, "\n")
}

func leadingWhitespace(value string) int {
	count := 0
	for _, char := range value {
		if !unicode.IsSpace(char) || char == '\n' || char == '\r' {
			break
		}
		count++
	}
	return count
}

func trimLeadingWhitespace(value string, count int) string {
	removed := 0
	for index, char := range value {
		if removed == count || !unicode.IsSpace(char) || char == '\n' || char == '\r' {
			return value[index:]
		}
		removed++
	}
	return ""
}
