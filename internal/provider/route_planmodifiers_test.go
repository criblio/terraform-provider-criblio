package provider

import (
	"encoding/json"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/require"
)

func TestRouteIDsFollowRoutes(t *testing.T) {
	list := func(raw string) types.List {
		t.Helper()
		var table RoutesModel
		require.NoError(t, json.Unmarshal([]byte(`{"routes":`+raw+`}`), &table))
		return table.Routes
	}
	prior := list(`[{"name":"duplicate","pipeline":"one","id":"a"},{"name":"duplicate","pipeline":"two","id":"b"},{"name":"last","pipeline":"three","id":"c"}]`)
	t.Run("reserve specific filters before matching omitted filters", func(t *testing.T) {
		state := list(`[{"name":"same","pipeline":"main","filter":"specific","id":"filtered"},{"name":"same","pipeline":"main","filter":"true","id":"all"}]`)
		config := list(`[{"name":"same","pipeline":"main"},{"name":"same","pipeline":"main","filter":"specific"}]`)
		resp := planmodifier.ListResponse{PlanValue: config}
		routeIDsByName{}.PlanModifyList(t.Context(), planmodifier.ListRequest{ConfigValue: config, PlanValue: config, StateValue: state}, &resp)
		require.Empty(t, resp.Diagnostics)
		ids := resp.PlanValue.Elements()
		require.Equal(t, "all", objectString(ids[0], "id"))
		require.Equal(t, "filtered", objectString(ids[1], "id"))
	})
	for _, tc := range []struct {
		name, routes string
		ids          []types.String
	}{
		{"reorder and rename", `[{"name":"duplicate","pipeline":"two"},{"name":"duplicate","pipeline":"one"},{"name":"renamed","pipeline":"three"}]`,
			[]types.String{types.StringValue("b"), types.StringValue("a"), types.StringValue("c")}},
		{"insert", `[{"name":"new","pipeline":"one"},{"name":"last","pipeline":"three"},{"name":"duplicate","pipeline":"two"},{"name":"duplicate","pipeline":"one"}]`,
			[]types.String{types.StringUnknown(), types.StringValue("c"), types.StringValue("b"), types.StringValue("a")}},
		{"remove and edit", `[{"name":"last","pipeline":"four"},{"name":"duplicate","pipeline":"two"}]`,
			[]types.String{types.StringValue("c"), types.StringValue("b")}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			plan := list(tc.routes)
			resp := planmodifier.ListResponse{PlanValue: plan}
			routeIDsByName{}.PlanModifyList(t.Context(), planmodifier.ListRequest{ConfigValue: plan, PlanValue: plan, StateValue: prior}, &resp)
			require.Empty(t, resp.Diagnostics)
			for i, element := range resp.PlanValue.Elements() {
				require.Equal(t, tc.ids[i], element.(types.Object).Attributes()["id"])
			}
		})
	}
}
