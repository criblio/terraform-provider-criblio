package tests

import (
	"fmt"
	"os"
	"slices"
	"strconv"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

// TestWorkerGroupUpgrade submits one real upgrade to an explicitly selected group.
// The group must have running workers eligible for a job-based upgrade. Terraform
// cleanup only forgets the action; it does not cancel or roll back that upgrade.
func TestWorkerGroupUpgrade(t *testing.T) {
	if os.Getenv("TF_ACC") == "" {
		t.Skip("Acceptance tests require TF_ACC")
	}
	if os.Getenv("CRIBL_TEST_WORKER_GROUP_UPGRADE") != "1" {
		t.Skip("Set CRIBL_TEST_WORKER_GROUP_UPGRADE=1 to submit a live Worker Group upgrade")
	}
	groupID := strings.TrimSpace(os.Getenv("CRIBL_TEST_WORKER_GROUP_UPGRADE_GROUP_ID"))
	if groupID == "" {
		t.Fatal("CRIBL_TEST_WORKER_GROUP_UPGRADE_GROUP_ID must identify a disposable self-managed Worker Group with eligible running workers")
	}

	const name = "criblio_worker_group_upgrade.test"
	groupConfig := fmt.Sprintf(`data "criblio_group" "target" {
  id = %q
}
`, groupID)
	config := func(trigger, settings string) string {
		return groupConfig + fmt.Sprintf(`resource "criblio_worker_group_upgrade" "test" {
  group_id = data.criblio_group.target.id
  trigger  = %q
  %s
}
`, trigger, settings)
	}
	initial := config("acceptance-1", "")
	var submittedJobs []string
	checkJobs := func(state *terraform.State) error {
		entry := state.RootModule().Resources[name]
		if entry == nil || entry.Primary == nil {
			return fmt.Errorf("upgrade action is missing from state")
		}
		attributes := entry.Primary.Attributes
		count, err := strconv.Atoi(attributes["job_ids.#"])
		if err != nil || count < 1 {
			return fmt.Errorf("upgrade returned no job IDs; the selected group must have eligible running workers (job_ids.#=%q)", attributes["job_ids.#"])
		}
		jobs := make([]string, count)
		for i := range jobs {
			jobs[i] = attributes[fmt.Sprintf("job_ids.%d", i)]
			if jobs[i] == "" {
				return fmt.Errorf("upgrade job ID %d is empty", i)
			}
		}
		if submittedJobs == nil {
			submittedJobs = jobs
			t.Logf("Submitted Worker Group upgrade jobs: %v", jobs)
			return nil
		}
		if !slices.Equal(jobs, submittedJobs) {
			return fmt.Errorf("upgrade job IDs changed without a new submission: got %v, want %v", jobs, submittedJobs)
		}
		return nil
	}
	check := resource.ComposeAggregateTestCheckFunc(
		resource.TestCheckResourceAttr(name, "group_id", groupID),
		resource.TestCheckResourceAttr(name, "trigger", "acceptance-1"),
		resource.TestCheckResourceAttr(name, "upgrade_mode", "rolling"),
		checkJobs,
	)
	noop := resource.ConfigPlanChecks{PreApply: []plancheck.PlanCheck{plancheck.ExpectResourceAction(name, plancheck.ResourceActionNoop)}}
	replacement := resource.ConfigPlanChecks{PostApplyPreRefresh: []plancheck.PlanCheck{plancheck.ExpectResourceAction(name, plancheck.ResourceActionReplace)}}
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories:  providerFactory,
		PreventPostDestroyRefresh: true,
		Steps: []resource.TestStep{
			// Confirm the group exists before submitting the action.
			{Config: groupConfig, Check: resource.TestCheckResourceAttr("data.criblio_group.target", "id", groupID)},
			{Config: initial, Check: check, ConfigPlanChecks: resource.ConfigPlanChecks{
				PreApply: []plancheck.PlanCheck{plancheck.ExpectResourceAction(name, plancheck.ResourceActionCreate)},
			}},
			{Config: initial, Check: check, ConfigPlanChecks: noop},
			{RefreshState: true, Check: check},
			// Plan replacements without submitting overlapping live upgrades.
			{Config: config("acceptance-2", ""), PlanOnly: true, ExpectNonEmptyPlan: true, ConfigPlanChecks: replacement},
			{Config: config("acceptance-1", "upgrade_percentage = 0.5"), PlanOnly: true, ExpectNonEmptyPlan: true, ConfigPlanChecks: replacement},
			{Config: initial, Check: check, ConfigPlanChecks: noop},
			{Config: groupConfig, ConfigPlanChecks: resource.ConfigPlanChecks{
				PreApply: []plancheck.PlanCheck{plancheck.ExpectResourceAction(name, plancheck.ResourceActionDestroy)},
			}, Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr("data.criblio_group.target", "id", groupID),
				func(state *terraform.State) error {
					if state.RootModule().Resources[name] != nil {
						return fmt.Errorf("upgrade action remains in Terraform state after destroy")
					}
					return nil
				},
			)},
		},
	})
}
