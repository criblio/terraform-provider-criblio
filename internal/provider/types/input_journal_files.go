// Code generated by Speakeasy (https://speakeasy.com). DO NOT EDIT.

package types

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type InputJournalFiles struct {
	Connections  []InputJournalFilesConnection `tfsdk:"connections"`
	CurrentBoot  types.Bool                    `tfsdk:"current_boot"`
	Description  types.String                  `tfsdk:"description"`
	Disabled     types.Bool                    `tfsdk:"disabled"`
	Environment  types.String                  `tfsdk:"environment"`
	ID           types.String                  `tfsdk:"id"`
	Interval     types.Float64                 `tfsdk:"interval"`
	Journals     []types.String                `tfsdk:"journals"`
	MaxAgeDur    types.String                  `tfsdk:"max_age_dur"`
	Metadata     []InputJournalFilesMetadatum  `tfsdk:"metadata"`
	Path         types.String                  `tfsdk:"path"`
	Pipeline     types.String                  `tfsdk:"pipeline"`
	Pq           *InputJournalFilesPq          `tfsdk:"pq"`
	PqEnabled    types.Bool                    `tfsdk:"pq_enabled"`
	Rules        []InputJournalFilesRule       `tfsdk:"rules"`
	SendToRoutes types.Bool                    `tfsdk:"send_to_routes"`
	Streamtags   []types.String                `tfsdk:"streamtags"`
	Type         types.String                  `tfsdk:"type"`
}
