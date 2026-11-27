package tests

import (
	"fmt"
	"testing"

	"github.com/criblio/terraform-provider-criblio/internal/auth"
	"github.com/criblio/terraform-provider-criblio/internal/restclient"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/stretchr/testify/require"
)

// Exercise drift on a disposable pack so the test never edits the group's routing table.
func TestPackPipelineAndRoutesDrift(t *testing.T) {
	pack := acctest.RandomWithPrefix("tf-drift")
	const pipeline = "criblio_pack_pipeline.test"
	const routes = "criblio_pack_routes.test"
	config := packMetadataConfig(pack, "drift", "1.0.0") + `
resource "criblio_pack_pipeline" "test" {
  group_id = "default"
  pack = criblio_pack.metadata.id
  id = "drift"
  conf = {
    functions = [{
      id = "eval"
      filter = "true"
      conf = jsonencode({add = [{name = "answer", value = "42"}]})
    }, { id = "drop", disabled = true, conf = jsonencode({}) }]
  }
}
resource "criblio_pack_routes" "test" {
  group_id = "default"
  pack = criblio_pack.metadata.id
  routes = [
    { name = "first", pipeline = criblio_pack_pipeline.test.id },
    { name = "first", pipeline = criblio_pack_pipeline.test.id, filter = "false" },
    { name = "second", pipeline = "main" }
  ]
}
`
	packPath := "/m/default/p/" + pack
	pipelinePath := packPath + "/pipelines/drift"
	routesPath := packPath + "/routes/default"
	var client *restclient.Client
	type route struct{ Name, ID, Filter string }
	readRoutes := func() []route {
		t.Helper()
		table, err := restclient.Get[struct{ Routes []route }](t.Context(), client, routesPath)
		require.NoError(t, err)
		require.Len(t, table.Routes, 3)
		return table.Routes
	}
	var original []route
	patch := func(path string, body map[string]any) {
		t.Helper()
		_, err := restclient.Patch[map[string]any, any](t.Context(), client, path, body)
		require.NoError(t, err)
	}
	setFunctions := func(value string, printUndefineds, reverse bool) {
		eval := map[string]any{"id": "eval", "filter": "true", "conf": map[string]any{
			"printUndefineds": printUndefineds,
			"add":             []any{map[string]string{"name": "answer", "value": value}},
		}}
		drop := map[string]any{"id": "drop", "disabled": true, "conf": map[string]any{}}
		functions := []any{eval, drop}
		if reverse {
			functions[0], functions[1] = functions[1], functions[0]
		}
		patch(pipelinePath, map[string]any{"id": "drift", "conf": map[string]any{"functions": functions}})
	}
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: providerFactory,
		PreCheck: func() {
			credentials, err := auth.GetCredentials()
			require.NoError(t, err)
			client = restclient.New(restclient.Config{BaseURL: auth.ConstructBaseURL(auth.ConstructBaseURLInput{}, credentials), Credentials: credentials})
		},
		Steps: []resource.TestStep{
			{Config: config},
			{PreConfig: func() { setFunctions("42", false, false) }, Config: config, PlanOnly: true},
			{PreConfig: func() {
				original = readRoutes()
				require.Equal(t, "first", original[0].Name)
				require.Equal(t, "first", original[1].Name)
				require.Equal(t, "second", original[2].Name)
				require.NotEmpty(t, original[0].ID)
				require.NotEmpty(t, original[1].ID)
				require.NotEqual(t, original[0].ID, original[1].ID)
				setFunctions("17", true, true)
				patch(routesPath, map[string]any{"id": "default", "routes": []any{
					map[string]string{"name": "second", "id": original[2].ID, "pipeline": "main"},
					map[string]string{"name": "first", "id": original[1].ID, "pipeline": "drift", "filter": "false"},
					map[string]string{"name": "first", "id": original[0].ID, "pipeline": "drift"},
				}})
			}, RefreshState: true, ExpectNonEmptyPlan: true, Check: resource.TestCheckResourceAttr(routes, "routes.0.name", "second")},
			// Import must reproduce the refreshed state while the API still contains the edits.
			{ResourceName: pipeline, ImportState: true, ImportStateVerify: true,
				ImportStateId: fmt.Sprintf(`{"group_id":"default","id":"drift","pack":%q}`, pack)},
			{Config: config, ConfigPlanChecks: resource.ConfigPlanChecks{PreApply: []plancheck.PlanCheck{
				plancheck.ExpectResourceAction(pipeline, plancheck.ResourceActionUpdate),
				plancheck.ExpectResourceAction(routes, plancheck.ResourceActionUpdate),
			}}},
			{PreConfig: func() {
				type function struct {
					ID       string
					Disabled bool
					Conf     struct {
						Add             []struct{ Value string }
						PrintUndefineds bool
					}
				}
				deployed, err := restclient.Get[struct {
					Conf struct{ Functions []function }
				}](t.Context(), client, pipelinePath)
				require.NoError(t, err)
				require.Len(t, deployed.Conf.Functions, 2)
				require.Equal(t, "eval", deployed.Conf.Functions[0].ID)
				require.False(t, deployed.Conf.Functions[0].Disabled)
				require.Equal(t, "drop", deployed.Conf.Functions[1].ID)
				require.True(t, deployed.Conf.Functions[1].Disabled)
				conf := deployed.Conf.Functions[0].Conf
				require.Len(t, conf.Add, 1)
				require.Equal(t, "42", conf.Add[0].Value)
				require.False(t, conf.PrintUndefineds)
				require.Equal(t, original, readRoutes())
			}, Config: config, PlanOnly: true},
		},
	})
}
