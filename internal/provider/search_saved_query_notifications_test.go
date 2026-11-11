package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func savedQueryScheduleObject(t *testing.T, disabled types.Bool) types.Object {
	t.Helper()
	notifications, diags := types.ObjectValue(SearchSavedQueryScheduleNotificationsAttrTypes(), map[string]attr.Value{
		"disabled": disabled,
	})
	if diags.HasError() {
		t.Fatalf("build notifications object: %v", diags)
	}
	schedule, diags := types.ObjectValue(SearchSavedQueryScheduleAttrTypes(), map[string]attr.Value{
		"cron_schedule":  types.StringValue("0 * * * *"),
		"enabled":        types.BoolValue(false),
		"jitter_percent": types.Int64Value(10),
		"keep_last_n":    types.Int64Value(2),
		"notifications":  notifications,
		"tz":             types.StringValue("UTC"),
	})
	if diags.HasError() {
		t.Fatalf("build schedule object: %v", diags)
	}
	return schedule
}

func savedQueryScheduleDisabled(t *testing.T, schedule types.Object) types.Bool {
	t.Helper()
	notifications, ok := schedule.Attributes()["notifications"].(types.Object)
	if !ok {
		t.Fatal("schedule has no notifications object")
	}
	disabled, ok := notifications.Attributes()["disabled"].(types.Bool)
	if !ok {
		t.Fatal("notifications has no disabled attribute")
	}
	return disabled
}

// A plan with "notifications = {}" leaves disabled unknown; Create and Update
// preserve the planned schedule, so the apply helper must resolve the unknown
// before the model is persisted as state.
func TestSavedQueryNotificationsUnknownResolved(t *testing.T) {
	for _, tc := range []struct {
		name string
		api  types.Bool
		want types.Bool
	}{
		{name: "api known", api: types.BoolValue(true), want: types.BoolValue(true)},
		{name: "api null", api: types.BoolNull(), want: types.BoolNull()},
	} {
		t.Run(tc.name, func(t *testing.T) {
			apiModel := &SearchSavedQueryModel{Schedule: savedQueryScheduleObject(t, tc.api)}
			state := &SearchSavedQueryModel{Schedule: savedQueryScheduleObject(t, types.BoolUnknown())}
			applySearchSavedQueryAPIToState(apiModel, state, true, false)
			got := savedQueryScheduleDisabled(t, state.Schedule)
			if !got.Equal(tc.want) {
				t.Fatalf("disabled = %v, want %v", got, tc.want)
			}
		})
	}
}

// Known values from prior state or configuration must never be overwritten by
// the unknown-resolution pass during refresh.
func TestSavedQueryNotificationsKnownPreserved(t *testing.T) {
	apiModel := &SearchSavedQueryModel{Schedule: savedQueryScheduleObject(t, types.BoolValue(false))}
	state := &SearchSavedQueryModel{Schedule: savedQueryScheduleObject(t, types.BoolValue(true))}
	applySearchSavedQueryAPIToState(apiModel, state, true, false)
	got := savedQueryScheduleDisabled(t, state.Schedule)
	if !got.Equal(types.BoolValue(true)) {
		t.Fatalf("disabled = %v, want true (preserved input)", got)
	}
}
