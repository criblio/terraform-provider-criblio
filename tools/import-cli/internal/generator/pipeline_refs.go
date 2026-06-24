package generator

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/criblio/terraform-provider-criblio/tools/import-cli/internal/hcl"
)

// ApplyPipelineFunctionProcessorReferences rewrites pipeline chain function
// conf.processor values that point at another exported pipeline into Terraform
// references. This lets Terraform build a graph edge and create linked pipelines
// in the right order.
func ApplyPipelineFunctionProcessorReferences(items []ResourceItem) {
	pipelineNamesByID := make(map[string]string)
	for _, it := range items {
		if it.TypeName != "criblio_pipeline" {
			continue
		}
		id, ok := stringAttr(it.Attrs, "id")
		if !ok || id == "" {
			continue
		}
		pipelineNamesByID[id] = it.Name
	}
	if len(pipelineNamesByID) == 0 {
		return
	}

	for i := range items {
		if items[i].TypeName != "criblio_pipeline" {
			continue
		}
		rewritePipelineFunctionProcessorRefs(items[i].Attrs, pipelineNamesByID)
	}
}

func rewritePipelineFunctionProcessorRefs(attrs map[string]hcl.Value, pipelineNamesByID map[string]string) {
	conf, ok := attrs["conf"]
	if !ok || conf.Kind != hcl.KindMap {
		return
	}
	functions, ok := conf.Map["functions"]
	if !ok || functions.Kind != hcl.KindList {
		return
	}
	for i := range functions.List {
		fn := functions.List[i]
		if fn.Kind != hcl.KindMap {
			continue
		}
		confValue, ok := fn.Map["conf"]
		if !ok || confValue.Kind != hcl.KindString {
			continue
		}
		expr, rewritten := pipelineFunctionConfExpr(confValue.String, pipelineNamesByID)
		if !rewritten {
			continue
		}
		fn.Map["conf"] = hcl.Value{Kind: hcl.KindExpression, Expr: expr}
		functions.List[i] = fn
	}
	conf.Map["functions"] = functions
	attrs["conf"] = conf
}

func pipelineFunctionConfExpr(confJSON string, pipelineNamesByID map[string]string) (string, bool) {
	dec := json.NewDecoder(bytes.NewBufferString(confJSON))
	dec.UseNumber()
	var data map[string]interface{}
	if err := dec.Decode(&data); err != nil {
		return "", false
	}
	processor, ok := data["processor"].(string)
	if !ok || processor == "" {
		return "", false
	}
	name, ok := pipelineNamesByID[processor]
	if !ok {
		return "", false
	}
	confValue := interfaceToHCLValue(data)
	confValue.Map["processor"] = hcl.Value{
		Kind: hcl.KindExpression,
		Expr: fmt.Sprintf("criblio_pipeline.%s.id", name),
	}
	return "jsonencode(" + confValue.ToHCLExpr() + ")", true
}

func interfaceToHCLValue(v interface{}) hcl.Value {
	switch t := v.(type) {
	case nil:
		return hcl.Value{Kind: hcl.KindNull}
	case string:
		return hcl.Value{Kind: hcl.KindString, String: t}
	case bool:
		return hcl.Value{Kind: hcl.KindBool, Bool: t}
	case json.Number:
		f, err := t.Float64()
		if err != nil {
			return hcl.Value{Kind: hcl.KindString, String: t.String()}
		}
		return hcl.Value{Kind: hcl.KindNumber, Number: f}
	case float64:
		return hcl.Value{Kind: hcl.KindNumber, Number: t}
	case []interface{}:
		list := make([]hcl.Value, len(t))
		for i, item := range t {
			list[i] = interfaceToHCLValue(item)
		}
		return hcl.Value{Kind: hcl.KindList, List: list}
	case map[string]interface{}:
		m := make(map[string]hcl.Value, len(t))
		for k, item := range t {
			m[k] = interfaceToHCLValue(item)
		}
		return hcl.Value{Kind: hcl.KindMap, Map: m}
	default:
		return hcl.Value{Kind: hcl.KindString, String: fmt.Sprint(t)}
	}
}

func stringAttr(attrs map[string]hcl.Value, name string) (string, bool) {
	v, ok := attrs[name]
	if !ok || v.Kind != hcl.KindString {
		return "", false
	}
	return v.String, true
}
