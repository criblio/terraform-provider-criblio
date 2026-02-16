package stringplanmodifier

import (
	"context"
	"encoding/json"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
)

// PipelineConfSuppressWhitespaceDiff returns a plan modifier that uses the state value
// when the planned and state JSON conf are semantically equal after normalizing
// whitespace in the "code" field (e.g. from API vs heredoc formatting).
// This prevents perpetual diff when pipeline functions use a "code" block.
func PipelineConfSuppressWhitespaceDiff() planmodifier.String {
	return pipelineConfSuppressWhitespaceDiff{}
}

type pipelineConfSuppressWhitespaceDiff struct{}

func (pipelineConfSuppressWhitespaceDiff) Description(_ context.Context) string {
	return "Suppresses diff when the only difference is whitespace in the code field."
}

func (pipelineConfSuppressWhitespaceDiff) MarkdownDescription(_ context.Context) string {
	return "Suppresses diff when the only difference is whitespace in the code field."
}

func (pipelineConfSuppressWhitespaceDiff) PlanModifyString(ctx context.Context, req planmodifier.StringRequest, resp *planmodifier.StringResponse) {
	if req.PlanValue.IsUnknown() || req.StateValue.IsNull() || req.StateValue.IsUnknown() {
		return
	}
	planStr := req.PlanValue.ValueString()
	stateStr := req.StateValue.ValueString()
	if normalizePipelineFunctionConfJSON(planStr) == normalizePipelineFunctionConfJSON(stateStr) {
		resp.PlanValue = req.StateValue
	}
}

// normalizePipelineFunctionConfJSON parses the JSON, normalizes the "code" field
// (trim and normalize line endings), and re-encodes so comparison is stable.
func normalizePipelineFunctionConfJSON(s string) string {
	var m map[string]interface{}
	if err := json.Unmarshal([]byte(s), &m); err != nil {
		return s
	}
	if code, ok := m["code"].(string); ok {
		m["code"] = normalizeCodeString(code)
	}
	out, err := json.Marshal(m)
	if err != nil {
		return s
	}
	return string(out)
}

func normalizeCodeString(code string) string {
	code = strings.ReplaceAll(code, "\r\n", "\n")
	code = strings.TrimSpace(code)
	// Remove common leading indent so heredoc formatting matches API formatting
	lines := strings.Split(code, "\n")
	if len(lines) == 0 {
		return code
	}
	minIndent := -1
	for _, line := range lines {
		if strings.TrimSpace(line) == "" {
			continue
		}
		trimmed := strings.TrimLeft(line, " \t")
		indent := len(line) - len(trimmed)
		if minIndent == -1 || indent < minIndent {
			minIndent = indent
		}
	}
	if minIndent <= 0 {
		return code
	}
	for i, line := range lines {
		if len(line) >= minIndent {
			lines[i] = line[minIndent:]
		}
	}
	return strings.TrimSpace(strings.Join(lines, "\n"))
}
