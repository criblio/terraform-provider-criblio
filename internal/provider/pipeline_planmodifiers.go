package provider

import (
	"context"
	"encoding/json"
	"strings"
	"unicode"

	custom_stringplanmodifier "github.com/criblio/terraform-provider-criblio/internal/planmodifiers/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
)

type pipelineConfSuppressWhitespaceDiff struct{}

func pipelineConfPlanModifiers() []planmodifier.String {
	return []planmodifier.String{
		custom_stringplanmodifier.SuppressDiff(custom_stringplanmodifier.ExplicitSuppress),
		pipelineConfSuppressWhitespaceDiff{},
	}
}

func (pipelineConfSuppressWhitespaceDiff) Description(context.Context) string {
	return "Suppresses pipeline function conf diffs caused only by API code whitespace normalization."
}

func (pipelineConfSuppressWhitespaceDiff) MarkdownDescription(ctx context.Context) string {
	return pipelineConfSuppressWhitespaceDiff{}.Description(ctx)
}

func (pipelineConfSuppressWhitespaceDiff) PlanModifyString(ctx context.Context, req planmodifier.StringRequest, resp *planmodifier.StringResponse) {
	if req.PlanValue.IsUnknown() || req.PlanValue.IsNull() || req.StateValue.IsUnknown() || req.StateValue.IsNull() {
		return
	}
	if normalizePipelineConf(req.PlanValue.ValueString()) == normalizePipelineConf(req.StateValue.ValueString()) {
		resp.PlanValue = req.StateValue
	}
}

func normalizePipelineConf(value string) string {
	var decoded any
	if err := json.Unmarshal([]byte(value), &decoded); err != nil {
		return normalizePipelineCode(value)
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
