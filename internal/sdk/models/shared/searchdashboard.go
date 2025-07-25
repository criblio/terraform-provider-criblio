// Code generated by Speakeasy (https://speakeasy.com). DO NOT EDIT.

package shared

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/criblio/terraform-provider-criblio/internal/sdk/internal/utils"
)

type SearchDashboardType string

const (
	SearchDashboardTypeMarkdownDefault SearchDashboardType = "markdown.default"
)

func (e SearchDashboardType) ToPointer() *SearchDashboardType {
	return &e
}
func (e *SearchDashboardType) UnmarshalJSON(data []byte) error {
	var v string
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	switch v {
	case "markdown.default":
		*e = SearchDashboardType(v)
		return nil
	default:
		return fmt.Errorf("invalid value for SearchDashboardType: %v", v)
	}
}

type Variant string

const (
	VariantMarkdown Variant = "markdown"
)

func (e Variant) ToPointer() *Variant {
	return &e
}
func (e *Variant) UnmarshalJSON(data []byte) error {
	var v string
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	switch v {
	case "markdown":
		*e = Variant(v)
		return nil
	default:
		return fmt.Errorf("invalid value for Variant: %v", v)
	}
}

type ElementMarkdown struct {
	Description *string             `json:"description,omitempty"`
	Empty       *bool               `json:"empty,omitempty"`
	HidePanel   *bool               `json:"hidePanel,omitempty"`
	ID          string              `json:"id"`
	Index       *float64            `json:"index,omitempty"`
	Layout      DashboardLayout     `json:"layout"`
	Title       *string             `json:"title,omitempty"`
	Type        SearchDashboardType `json:"type"`
	Value       *string             `json:"value,omitempty"`
	Variant     Variant             `json:"variant"`
}

func (o *ElementMarkdown) GetDescription() *string {
	if o == nil {
		return nil
	}
	return o.Description
}

func (o *ElementMarkdown) GetEmpty() *bool {
	if o == nil {
		return nil
	}
	return o.Empty
}

func (o *ElementMarkdown) GetHidePanel() *bool {
	if o == nil {
		return nil
	}
	return o.HidePanel
}

func (o *ElementMarkdown) GetID() string {
	if o == nil {
		return ""
	}
	return o.ID
}

func (o *ElementMarkdown) GetIndex() *float64 {
	if o == nil {
		return nil
	}
	return o.Index
}

func (o *ElementMarkdown) GetLayout() DashboardLayout {
	if o == nil {
		return DashboardLayout{}
	}
	return o.Layout
}

func (o *ElementMarkdown) GetTitle() *string {
	if o == nil {
		return nil
	}
	return o.Title
}

func (o *ElementMarkdown) GetType() SearchDashboardType {
	if o == nil {
		return SearchDashboardType("")
	}
	return o.Type
}

func (o *ElementMarkdown) GetValue() *string {
	if o == nil {
		return nil
	}
	return o.Value
}

func (o *ElementMarkdown) GetVariant() Variant {
	if o == nil {
		return Variant("")
	}
	return o.Variant
}

type Element struct {
	Description     *string                  `json:"description,omitempty"`
	Empty           *bool                    `json:"empty,omitempty"`
	HidePanel       *bool                    `json:"hidePanel,omitempty"`
	HorizontalChart *bool                    `json:"horizontalChart,omitempty"`
	ID              string                   `json:"id"`
	Index           *float64                 `json:"index,omitempty"`
	InputID         *string                  `json:"inputId,omitempty"`
	Layout          DashboardLayout          `json:"layout"`
	Search          SearchQuery              `json:"search"`
	Title           *string                  `json:"title,omitempty"`
	Type            DashboardElementType     `json:"type"`
	Value           map[string]any           `json:"value,omitempty"`
	Variant         *DashboardElementVariant `json:"variant,omitempty"`
}

func (o *Element) GetDescription() *string {
	if o == nil {
		return nil
	}
	return o.Description
}

func (o *Element) GetEmpty() *bool {
	if o == nil {
		return nil
	}
	return o.Empty
}

func (o *Element) GetHidePanel() *bool {
	if o == nil {
		return nil
	}
	return o.HidePanel
}

func (o *Element) GetHorizontalChart() *bool {
	if o == nil {
		return nil
	}
	return o.HorizontalChart
}

func (o *Element) GetID() string {
	if o == nil {
		return ""
	}
	return o.ID
}

func (o *Element) GetIndex() *float64 {
	if o == nil {
		return nil
	}
	return o.Index
}

func (o *Element) GetInputID() *string {
	if o == nil {
		return nil
	}
	return o.InputID
}

func (o *Element) GetLayout() DashboardLayout {
	if o == nil {
		return DashboardLayout{}
	}
	return o.Layout
}

func (o *Element) GetSearch() SearchQuery {
	if o == nil {
		return SearchQuery{}
	}
	return o.Search
}

func (o *Element) GetTitle() *string {
	if o == nil {
		return nil
	}
	return o.Title
}

func (o *Element) GetType() DashboardElementType {
	if o == nil {
		return DashboardElementType("")
	}
	return o.Type
}

func (o *Element) GetValue() map[string]any {
	if o == nil {
		return nil
	}
	return o.Value
}

func (o *Element) GetVariant() *DashboardElementVariant {
	if o == nil {
		return nil
	}
	return o.Variant
}

type ElementUnionType string

const (
	ElementUnionTypeElement         ElementUnionType = "element"
	ElementUnionTypeElementMarkdown ElementUnionType = "element_Markdown"
)

type ElementUnion struct {
	Element         *Element         `queryParam:"inline"`
	ElementMarkdown *ElementMarkdown `queryParam:"inline"`

	Type ElementUnionType
}

func CreateElementUnionElement(element Element) ElementUnion {
	typ := ElementUnionTypeElement

	return ElementUnion{
		Element: &element,
		Type:    typ,
	}
}

func CreateElementUnionElementMarkdown(elementMarkdown ElementMarkdown) ElementUnion {
	typ := ElementUnionTypeElementMarkdown

	return ElementUnion{
		ElementMarkdown: &elementMarkdown,
		Type:            typ,
	}
}

func (u *ElementUnion) UnmarshalJSON(data []byte) error {

	var elementMarkdown ElementMarkdown = ElementMarkdown{}
	if err := utils.UnmarshalJSON(data, &elementMarkdown, "", true, true); err == nil {
		u.ElementMarkdown = &elementMarkdown
		u.Type = ElementUnionTypeElementMarkdown
		return nil
	}

	var element Element = Element{}
	if err := utils.UnmarshalJSON(data, &element, "", true, true); err == nil {
		u.Element = &element
		u.Type = ElementUnionTypeElement
		return nil
	}

	return fmt.Errorf("could not unmarshal `%s` into any supported union types for ElementUnion", string(data))
}

func (u ElementUnion) MarshalJSON() ([]byte, error) {
	if u.Element != nil {
		return utils.MarshalJSON(u.Element, "", true)
	}

	if u.ElementMarkdown != nil {
		return utils.MarshalJSON(u.ElementMarkdown, "", true)
	}

	return nil, errors.New("could not marshal union type ElementUnion: all fields are null")
}

type SearchDashboard struct {
	CacheTTLSeconds    *float64            `json:"cacheTTLSeconds,omitempty"`
	Category           *string             `json:"category,omitempty"`
	Created            float64             `json:"created"`
	CreatedBy          string              `json:"createdBy"`
	Description        *string             `json:"description,omitempty"`
	DisplayCreatedBy   *string             `json:"displayCreatedBy,omitempty"`
	DisplayModifiedBy  *string             `json:"displayModifiedBy,omitempty"`
	Elements           []ElementUnion      `json:"elements"`
	ID                 string              `json:"id"`
	Modified           float64             `json:"modified"`
	ModifiedBy         *string             `json:"modifiedBy,omitempty"`
	Name               string              `json:"name"`
	PackID             *string             `json:"packId,omitempty"`
	RefreshRate        *float64            `json:"refreshRate,omitempty"`
	ResolvedDatasetIds []string            `json:"resolvedDatasetIds,omitempty"`
	Schedule           *SavedQuerySchedule `json:"schedule,omitempty"`
}

func (o *SearchDashboard) GetCacheTTLSeconds() *float64 {
	if o == nil {
		return nil
	}
	return o.CacheTTLSeconds
}

func (o *SearchDashboard) GetCategory() *string {
	if o == nil {
		return nil
	}
	return o.Category
}

func (o *SearchDashboard) GetCreated() float64 {
	if o == nil {
		return 0.0
	}
	return o.Created
}

func (o *SearchDashboard) GetCreatedBy() string {
	if o == nil {
		return ""
	}
	return o.CreatedBy
}

func (o *SearchDashboard) GetDescription() *string {
	if o == nil {
		return nil
	}
	return o.Description
}

func (o *SearchDashboard) GetDisplayCreatedBy() *string {
	if o == nil {
		return nil
	}
	return o.DisplayCreatedBy
}

func (o *SearchDashboard) GetDisplayModifiedBy() *string {
	if o == nil {
		return nil
	}
	return o.DisplayModifiedBy
}

func (o *SearchDashboard) GetElements() []ElementUnion {
	if o == nil {
		return []ElementUnion{}
	}
	return o.Elements
}

func (o *SearchDashboard) GetID() string {
	if o == nil {
		return ""
	}
	return o.ID
}

func (o *SearchDashboard) GetModified() float64 {
	if o == nil {
		return 0.0
	}
	return o.Modified
}

func (o *SearchDashboard) GetModifiedBy() *string {
	if o == nil {
		return nil
	}
	return o.ModifiedBy
}

func (o *SearchDashboard) GetName() string {
	if o == nil {
		return ""
	}
	return o.Name
}

func (o *SearchDashboard) GetPackID() *string {
	if o == nil {
		return nil
	}
	return o.PackID
}

func (o *SearchDashboard) GetRefreshRate() *float64 {
	if o == nil {
		return nil
	}
	return o.RefreshRate
}

func (o *SearchDashboard) GetResolvedDatasetIds() []string {
	if o == nil {
		return nil
	}
	return o.ResolvedDatasetIds
}

func (o *SearchDashboard) GetSchedule() *SavedQuerySchedule {
	if o == nil {
		return nil
	}
	return o.Schedule
}
