// Code generated by Speakeasy (https://speakeasy.com). DO NOT EDIT.

package shared

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/criblio/terraform-provider-criblio/internal/sdk/internal/utils"
)

type CompatibilityChecks struct {
	Datatypes *bool    `json:"datatypes,omitempty"`
	StageIds  []string `json:"stageIds,omitempty"`
}

func (o *CompatibilityChecks) GetDatatypes() *bool {
	if o == nil {
		return nil
	}
	return o.Datatypes
}

func (o *CompatibilityChecks) GetStageIds() []string {
	if o == nil {
		return nil
	}
	return o.StageIds
}

type SearchJobEarliestType string

const (
	SearchJobEarliestTypeStr    SearchJobEarliestType = "str"
	SearchJobEarliestTypeNumber SearchJobEarliestType = "number"
)

type SearchJobEarliest struct {
	Str    *string  `queryParam:"inline"`
	Number *float64 `queryParam:"inline"`

	Type SearchJobEarliestType
}

func CreateSearchJobEarliestStr(str string) SearchJobEarliest {
	typ := SearchJobEarliestTypeStr

	return SearchJobEarliest{
		Str:  &str,
		Type: typ,
	}
}

func CreateSearchJobEarliestNumber(number float64) SearchJobEarliest {
	typ := SearchJobEarliestTypeNumber

	return SearchJobEarliest{
		Number: &number,
		Type:   typ,
	}
}

func (u *SearchJobEarliest) UnmarshalJSON(data []byte) error {

	var str string = ""
	if err := utils.UnmarshalJSON(data, &str, "", true, true); err == nil {
		u.Str = &str
		u.Type = SearchJobEarliestTypeStr
		return nil
	}

	var number float64 = float64(0)
	if err := utils.UnmarshalJSON(data, &number, "", true, true); err == nil {
		u.Number = &number
		u.Type = SearchJobEarliestTypeNumber
		return nil
	}

	return fmt.Errorf("could not unmarshal `%s` into any supported union types for SearchJobEarliest", string(data))
}

func (u SearchJobEarliest) MarshalJSON() ([]byte, error) {
	if u.Str != nil {
		return utils.MarshalJSON(u.Str, "", true)
	}

	if u.Number != nil {
		return utils.MarshalJSON(u.Number, "", true)
	}

	return nil, errors.New("could not marshal union type SearchJobEarliest: all fields are null")
}

type SearchJobLatestType string

const (
	SearchJobLatestTypeStr    SearchJobLatestType = "str"
	SearchJobLatestTypeNumber SearchJobLatestType = "number"
)

type SearchJobLatest struct {
	Str    *string  `queryParam:"inline"`
	Number *float64 `queryParam:"inline"`

	Type SearchJobLatestType
}

func CreateSearchJobLatestStr(str string) SearchJobLatest {
	typ := SearchJobLatestTypeStr

	return SearchJobLatest{
		Str:  &str,
		Type: typ,
	}
}

func CreateSearchJobLatestNumber(number float64) SearchJobLatest {
	typ := SearchJobLatestTypeNumber

	return SearchJobLatest{
		Number: &number,
		Type:   typ,
	}
}

func (u *SearchJobLatest) UnmarshalJSON(data []byte) error {

	var str string = ""
	if err := utils.UnmarshalJSON(data, &str, "", true, true); err == nil {
		u.Str = &str
		u.Type = SearchJobLatestTypeStr
		return nil
	}

	var number float64 = float64(0)
	if err := utils.UnmarshalJSON(data, &number, "", true, true); err == nil {
		u.Number = &number
		u.Type = SearchJobLatestTypeNumber
		return nil
	}

	return fmt.Errorf("could not unmarshal `%s` into any supported union types for SearchJobLatest", string(data))
}

func (u SearchJobLatest) MarshalJSON() ([]byte, error) {
	if u.Str != nil {
		return utils.MarshalJSON(u.Str, "", true)
	}

	if u.Number != nil {
		return utils.MarshalJSON(u.Number, "", true)
	}

	return nil, errors.New("could not marshal union type SearchJobLatest: all fields are null")
}

type SearchJobStatus string

const (
	SearchJobStatusNew       SearchJobStatus = "new"
	SearchJobStatusFailed    SearchJobStatus = "failed"
	SearchJobStatusRunning   SearchJobStatus = "running"
	SearchJobStatusCompleted SearchJobStatus = "completed"
	SearchJobStatusCanceled  SearchJobStatus = "canceled"
	SearchJobStatusQueued    SearchJobStatus = "queued"
)

func (e SearchJobStatus) ToPointer() *SearchJobStatus {
	return &e
}
func (e *SearchJobStatus) UnmarshalJSON(data []byte) error {
	var v string
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	switch v {
	case "new":
		fallthrough
	case "failed":
		fallthrough
	case "running":
		fallthrough
	case "completed":
		fallthrough
	case "canceled":
		fallthrough
	case "queued":
		*e = SearchJobStatus(v)
		return nil
	default:
		return fmt.Errorf("invalid value for SearchJobStatus: %v", v)
	}
}

type SearchJobType string

const (
	SearchJobTypeCommand         SearchJobType = "command"
	SearchJobTypeStandard        SearchJobType = "standard"
	SearchJobTypeDatatypePreview SearchJobType = "datatypePreview"
	SearchJobTypeScheduled       SearchJobType = "scheduled"
	SearchJobTypeDashboard       SearchJobType = "dashboard"
)

func (e SearchJobType) ToPointer() *SearchJobType {
	return &e
}
func (e *SearchJobType) UnmarshalJSON(data []byte) error {
	var v string
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	switch v {
	case "command":
		fallthrough
	case "standard":
		fallthrough
	case "datatypePreview":
		fallthrough
	case "scheduled":
		fallthrough
	case "dashboard":
		*e = SearchJobType(v)
		return nil
	default:
		return fmt.Errorf("invalid value for SearchJobType: %v", v)
	}
}

type SearchJob struct {
	Accelerated                 *bool                      `json:"accelerated,omitempty"`
	AliasOfOriginalJobID        *string                    `json:"aliasOfOriginalJobId,omitempty"`
	CompatibilityChecks         *CompatibilityChecks       `json:"compatibilityChecks,omitempty"`
	CompletionInfo              *string                    `json:"completionInfo,omitempty"`
	Context                     *string                    `json:"context,omitempty"`
	CorrelationID               *string                    `json:"correlationId,omitempty"`
	CPUMetrics                  *CPUTimeMetric             `json:"cpuMetrics,omitempty"`
	DatatypeOverrides           *DatatypeOverrides         `json:"datatypeOverrides,omitempty"`
	DisableNotifications        *bool                      `json:"disableNotifications,omitempty"`
	DisplayUsername             string                     `json:"displayUsername"`
	Earliest                    *SearchJobEarliest         `json:"earliest,omitempty"`
	EarliestEpoch               *float64                   `json:"earliestEpoch,omitempty"`
	ErrorStateConfig            *SearchJobErrorStateConfig `json:"errorStateConfig,omitempty"`
	Group                       string                     `json:"group"`
	ID                          string                     `json:"id"`
	IsPrivate                   *bool                      `json:"isPrivate,omitempty"`
	Latest                      *SearchJobLatest           `json:"latest,omitempty"`
	LatestEpoch                 *float64                   `json:"latestEpoch,omitempty"`
	Metadata                    *SearchJobMetadata         `json:"metadata,omitempty"`
	NumEventsAfter              *float64                   `json:"numEventsAfter,omitempty"`
	NumEventsBefore             *float64                   `json:"numEventsBefore,omitempty"`
	Query                       string                     `json:"query"`
	QueryWithMacrosResolved     *string                    `json:"queryWithMacrosResolved,omitempty"`
	SampleRate                  *float64                   `json:"sampleRate,omitempty"`
	SavedQueryName              *string                    `json:"savedQueryName,omitempty"`
	SearchParameterDeclarations []SearchParameter          `json:"searchParameterDeclarations,omitempty"`
	SearchParameterValues       any                        `json:"searchParameterValues,omitempty"`
	Stages                      []SearchJobStageConfig     `json:"stages,omitempty"`
	Status                      SearchJobStatus            `json:"status"`
	TableConfig                 *TableViewSettings         `json:"tableConfig,omitempty"`
	TargetEventTime             *float64                   `json:"targetEventTime,omitempty"`
	TimeCompleted               *float64                   `json:"timeCompleted,omitempty"`
	TimeCreated                 float64                    `json:"timeCreated"`
	TimeStarted                 float64                    `json:"timeStarted"`
	TimeToFirstByte             *float64                   `json:"timeToFirstByte,omitempty"`
	TotalBytesScanned           *float64                   `json:"totalBytesScanned,omitempty"`
	TotalEventCount             *float64                   `json:"totalEventCount,omitempty"`
	Type                        *SearchJobType             `json:"type,omitempty"`
	UsageGroupID                *string                    `json:"usageGroupId,omitempty"`
	UsageMetrics                *SearchAuditMetrics        `json:"usageMetrics,omitempty"`
	User                        string                     `json:"user"`
}

func (o *SearchJob) GetAccelerated() *bool {
	if o == nil {
		return nil
	}
	return o.Accelerated
}

func (o *SearchJob) GetAliasOfOriginalJobID() *string {
	if o == nil {
		return nil
	}
	return o.AliasOfOriginalJobID
}

func (o *SearchJob) GetCompatibilityChecks() *CompatibilityChecks {
	if o == nil {
		return nil
	}
	return o.CompatibilityChecks
}

func (o *SearchJob) GetCompletionInfo() *string {
	if o == nil {
		return nil
	}
	return o.CompletionInfo
}

func (o *SearchJob) GetContext() *string {
	if o == nil {
		return nil
	}
	return o.Context
}

func (o *SearchJob) GetCorrelationID() *string {
	if o == nil {
		return nil
	}
	return o.CorrelationID
}

func (o *SearchJob) GetCPUMetrics() *CPUTimeMetric {
	if o == nil {
		return nil
	}
	return o.CPUMetrics
}

func (o *SearchJob) GetDatatypeOverrides() *DatatypeOverrides {
	if o == nil {
		return nil
	}
	return o.DatatypeOverrides
}

func (o *SearchJob) GetDisableNotifications() *bool {
	if o == nil {
		return nil
	}
	return o.DisableNotifications
}

func (o *SearchJob) GetDisplayUsername() string {
	if o == nil {
		return ""
	}
	return o.DisplayUsername
}

func (o *SearchJob) GetEarliest() *SearchJobEarliest {
	if o == nil {
		return nil
	}
	return o.Earliest
}

func (o *SearchJob) GetEarliestEpoch() *float64 {
	if o == nil {
		return nil
	}
	return o.EarliestEpoch
}

func (o *SearchJob) GetErrorStateConfig() *SearchJobErrorStateConfig {
	if o == nil {
		return nil
	}
	return o.ErrorStateConfig
}

func (o *SearchJob) GetGroup() string {
	if o == nil {
		return ""
	}
	return o.Group
}

func (o *SearchJob) GetID() string {
	if o == nil {
		return ""
	}
	return o.ID
}

func (o *SearchJob) GetIsPrivate() *bool {
	if o == nil {
		return nil
	}
	return o.IsPrivate
}

func (o *SearchJob) GetLatest() *SearchJobLatest {
	if o == nil {
		return nil
	}
	return o.Latest
}

func (o *SearchJob) GetLatestEpoch() *float64 {
	if o == nil {
		return nil
	}
	return o.LatestEpoch
}

func (o *SearchJob) GetMetadata() *SearchJobMetadata {
	if o == nil {
		return nil
	}
	return o.Metadata
}

func (o *SearchJob) GetNumEventsAfter() *float64 {
	if o == nil {
		return nil
	}
	return o.NumEventsAfter
}

func (o *SearchJob) GetNumEventsBefore() *float64 {
	if o == nil {
		return nil
	}
	return o.NumEventsBefore
}

func (o *SearchJob) GetQuery() string {
	if o == nil {
		return ""
	}
	return o.Query
}

func (o *SearchJob) GetQueryWithMacrosResolved() *string {
	if o == nil {
		return nil
	}
	return o.QueryWithMacrosResolved
}

func (o *SearchJob) GetSampleRate() *float64 {
	if o == nil {
		return nil
	}
	return o.SampleRate
}

func (o *SearchJob) GetSavedQueryName() *string {
	if o == nil {
		return nil
	}
	return o.SavedQueryName
}

func (o *SearchJob) GetSearchParameterDeclarations() []SearchParameter {
	if o == nil {
		return nil
	}
	return o.SearchParameterDeclarations
}

func (o *SearchJob) GetSearchParameterValues() any {
	if o == nil {
		return nil
	}
	return o.SearchParameterValues
}

func (o *SearchJob) GetStages() []SearchJobStageConfig {
	if o == nil {
		return nil
	}
	return o.Stages
}

func (o *SearchJob) GetStatus() SearchJobStatus {
	if o == nil {
		return SearchJobStatus("")
	}
	return o.Status
}

func (o *SearchJob) GetTableConfig() *TableViewSettings {
	if o == nil {
		return nil
	}
	return o.TableConfig
}

func (o *SearchJob) GetTargetEventTime() *float64 {
	if o == nil {
		return nil
	}
	return o.TargetEventTime
}

func (o *SearchJob) GetTimeCompleted() *float64 {
	if o == nil {
		return nil
	}
	return o.TimeCompleted
}

func (o *SearchJob) GetTimeCreated() float64 {
	if o == nil {
		return 0.0
	}
	return o.TimeCreated
}

func (o *SearchJob) GetTimeStarted() float64 {
	if o == nil {
		return 0.0
	}
	return o.TimeStarted
}

func (o *SearchJob) GetTimeToFirstByte() *float64 {
	if o == nil {
		return nil
	}
	return o.TimeToFirstByte
}

func (o *SearchJob) GetTotalBytesScanned() *float64 {
	if o == nil {
		return nil
	}
	return o.TotalBytesScanned
}

func (o *SearchJob) GetTotalEventCount() *float64 {
	if o == nil {
		return nil
	}
	return o.TotalEventCount
}

func (o *SearchJob) GetType() *SearchJobType {
	if o == nil {
		return nil
	}
	return o.Type
}

func (o *SearchJob) GetUsageGroupID() *string {
	if o == nil {
		return nil
	}
	return o.UsageGroupID
}

func (o *SearchJob) GetUsageMetrics() *SearchAuditMetrics {
	if o == nil {
		return nil
	}
	return o.UsageMetrics
}

func (o *SearchJob) GetUser() string {
	if o == nil {
		return ""
	}
	return o.User
}
