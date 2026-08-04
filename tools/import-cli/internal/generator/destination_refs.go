package generator

import (
	"fmt"

	"github.com/criblio/terraform-provider-criblio/tools/import-cli/internal/hcl"
)

// ApplyDestinationRouterReferences rewrites router rule outputs that point to
// another exported destination into Terraform references. The references give
// Terraform the dependency edges required for safe create and destroy ordering.
func ApplyDestinationRouterReferences(items []ResourceItem) {
	if len(items) == 0 {
		return
	}
	typeName := items[0].TypeName
	if typeName != "criblio_destination" && typeName != "criblio_pack_destination" {
		return
	}

	namesByKey := make(map[string]string)
	for _, item := range items {
		id, ok := stringAttr(item.Attrs, "id")
		if !ok || id == "" {
			continue
		}
		namesByKey[destinationReferenceKey(item, id)] = item.Name
	}

	for i := range items {
		rewriteDestinationRouterReferences(&items[i], namesByKey)
	}
}

func rewriteDestinationRouterReferences(item *ResourceItem, namesByKey map[string]string) {
	router, ok := item.Attrs["output_router"]
	if !ok || router.Kind != hcl.KindMap {
		return
	}
	rules, ok := router.Map["rules"]
	if !ok || rules.Kind != hcl.KindList {
		return
	}

	for i := range rules.List {
		rule := rules.List[i]
		if rule.Kind != hcl.KindMap {
			continue
		}
		output, ok := rule.Map["output"]
		if !ok || output.Kind != hcl.KindString || output.String == "" {
			continue
		}
		name, ok := namesByKey[destinationReferenceKey(*item, output.String)]
		if !ok || name == item.Name {
			continue
		}
		rule.Map["output"] = hcl.Value{
			Kind: hcl.KindExpression,
			Expr: fmt.Sprintf("%s.%s.id", item.TypeName, name),
		}
		rules.List[i] = rule
	}

	router.Map["rules"] = rules
	item.Attrs["output_router"] = router
}

func destinationReferenceKey(item ResourceItem, id string) string {
	if item.TypeName != "criblio_pack_destination" {
		return item.GroupID + "\x00" + id
	}
	pack, _ := stringAttr(item.Attrs, "pack")
	return item.GroupID + "\x00" + pack + "\x00" + id
}
