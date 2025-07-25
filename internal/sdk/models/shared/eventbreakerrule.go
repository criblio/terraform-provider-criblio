// Code generated by Speakeasy (https://speakeasy.com). DO NOT EDIT.

package shared

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/criblio/terraform-provider-criblio/internal/sdk/internal/utils"
)

type TypeKvp string

const (
	TypeKvpKvp TypeKvp = "kvp"
)

func (e TypeKvp) ToPointer() *TypeKvp {
	return &e
}
func (e *TypeKvp) UnmarshalJSON(data []byte) error {
	var v string
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	switch v {
	case "kvp":
		*e = TypeKvp(v)
		return nil
	default:
		return fmt.Errorf("invalid value for TypeKvp: %v", v)
	}
}

type ParserKvp struct {
	AllowedKeyChars   *bool      `json:"allowedKeyChars,omitempty"`
	AllowedValueChars []string   `json:"allowedValueChars,omitempty"`
	CleanFields       []string   `json:"cleanFields,omitempty"`
	DelimChar         *string    `json:"delimChar,omitempty"`
	DstField          *string    `json:"dstField,omitempty"`
	EscapeChar        *string    `json:"escapeChar,omitempty"`
	FieldFilterExpr   *string    `json:"fieldFilterExpr,omitempty"`
	Fields            []string   `json:"fields,omitempty"`
	Keep              []string   `json:"keep,omitempty"`
	Mode              ParserMode `json:"mode"`
	NullValue         *string    `json:"nullValue,omitempty"`
	QuoteChar         *string    `json:"quoteChar,omitempty"`
	Remove            []string   `json:"remove,omitempty"`
	SrcField          string     `json:"srcField"`
	Type              TypeKvp    `json:"type"`
}

func (o *ParserKvp) GetAllowedKeyChars() *bool {
	if o == nil {
		return nil
	}
	return o.AllowedKeyChars
}

func (o *ParserKvp) GetAllowedValueChars() []string {
	if o == nil {
		return nil
	}
	return o.AllowedValueChars
}

func (o *ParserKvp) GetCleanFields() []string {
	if o == nil {
		return nil
	}
	return o.CleanFields
}

func (o *ParserKvp) GetDelimChar() *string {
	if o == nil {
		return nil
	}
	return o.DelimChar
}

func (o *ParserKvp) GetDstField() *string {
	if o == nil {
		return nil
	}
	return o.DstField
}

func (o *ParserKvp) GetEscapeChar() *string {
	if o == nil {
		return nil
	}
	return o.EscapeChar
}

func (o *ParserKvp) GetFieldFilterExpr() *string {
	if o == nil {
		return nil
	}
	return o.FieldFilterExpr
}

func (o *ParserKvp) GetFields() []string {
	if o == nil {
		return nil
	}
	return o.Fields
}

func (o *ParserKvp) GetKeep() []string {
	if o == nil {
		return nil
	}
	return o.Keep
}

func (o *ParserKvp) GetMode() ParserMode {
	if o == nil {
		return ParserMode("")
	}
	return o.Mode
}

func (o *ParserKvp) GetNullValue() *string {
	if o == nil {
		return nil
	}
	return o.NullValue
}

func (o *ParserKvp) GetQuoteChar() *string {
	if o == nil {
		return nil
	}
	return o.QuoteChar
}

func (o *ParserKvp) GetRemove() []string {
	if o == nil {
		return nil
	}
	return o.Remove
}

func (o *ParserKvp) GetSrcField() string {
	if o == nil {
		return ""
	}
	return o.SrcField
}

func (o *ParserKvp) GetType() TypeKvp {
	if o == nil {
		return TypeKvp("")
	}
	return o.Type
}

type TypeJSON string

const (
	TypeJSONJSON TypeJSON = "json"
)

func (e TypeJSON) ToPointer() *TypeJSON {
	return &e
}
func (e *TypeJSON) UnmarshalJSON(data []byte) error {
	var v string
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	switch v {
	case "json":
		*e = TypeJSON(v)
		return nil
	default:
		return fmt.Errorf("invalid value for TypeJSON: %v", v)
	}
}

type ParserJSON struct {
	CleanFields     []string   `json:"cleanFields,omitempty"`
	DstField        *string    `json:"dstField,omitempty"`
	FieldFilterExpr *string    `json:"fieldFilterExpr,omitempty"`
	Fields          []string   `json:"fields,omitempty"`
	Keep            []string   `json:"keep,omitempty"`
	Mode            ParserMode `json:"mode"`
	Remove          []string   `json:"remove,omitempty"`
	SrcField        string     `json:"srcField"`
	Type            TypeJSON   `json:"type"`
}

func (o *ParserJSON) GetCleanFields() []string {
	if o == nil {
		return nil
	}
	return o.CleanFields
}

func (o *ParserJSON) GetDstField() *string {
	if o == nil {
		return nil
	}
	return o.DstField
}

func (o *ParserJSON) GetFieldFilterExpr() *string {
	if o == nil {
		return nil
	}
	return o.FieldFilterExpr
}

func (o *ParserJSON) GetFields() []string {
	if o == nil {
		return nil
	}
	return o.Fields
}

func (o *ParserJSON) GetKeep() []string {
	if o == nil {
		return nil
	}
	return o.Keep
}

func (o *ParserJSON) GetMode() ParserMode {
	if o == nil {
		return ParserMode("")
	}
	return o.Mode
}

func (o *ParserJSON) GetRemove() []string {
	if o == nil {
		return nil
	}
	return o.Remove
}

func (o *ParserJSON) GetSrcField() string {
	if o == nil {
		return ""
	}
	return o.SrcField
}

func (o *ParserJSON) GetType() TypeJSON {
	if o == nil {
		return TypeJSON("")
	}
	return o.Type
}

type PatternList struct {
	Pattern string `json:"pattern"`
}

func (o *PatternList) GetPattern() string {
	if o == nil {
		return ""
	}
	return o.Pattern
}

type TypeGrok string

const (
	TypeGrokGrok TypeGrok = "grok"
)

func (e TypeGrok) ToPointer() *TypeGrok {
	return &e
}
func (e *TypeGrok) UnmarshalJSON(data []byte) error {
	var v string
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	switch v {
	case "grok":
		*e = TypeGrok(v)
		return nil
	default:
		return fmt.Errorf("invalid value for TypeGrok: %v", v)
	}
}

type ParserGrok struct {
	DstField    *string       `json:"dstField,omitempty"`
	Mode        ParserMode    `json:"mode"`
	Pattern     *string       `json:"pattern,omitempty"`
	PatternList []PatternList `json:"patternList,omitempty"`
	Source      *string       `json:"source,omitempty"`
	SrcField    string        `json:"srcField"`
	Type        TypeGrok      `json:"type"`
}

func (o *ParserGrok) GetDstField() *string {
	if o == nil {
		return nil
	}
	return o.DstField
}

func (o *ParserGrok) GetMode() ParserMode {
	if o == nil {
		return ParserMode("")
	}
	return o.Mode
}

func (o *ParserGrok) GetPattern() *string {
	if o == nil {
		return nil
	}
	return o.Pattern
}

func (o *ParserGrok) GetPatternList() []PatternList {
	if o == nil {
		return nil
	}
	return o.PatternList
}

func (o *ParserGrok) GetSource() *string {
	if o == nil {
		return nil
	}
	return o.Source
}

func (o *ParserGrok) GetSrcField() string {
	if o == nil {
		return ""
	}
	return o.SrcField
}

func (o *ParserGrok) GetType() TypeGrok {
	if o == nil {
		return TypeGrok("")
	}
	return o.Type
}

type ParserType string

const (
	ParserTypeClf   ParserType = "clf"
	ParserTypeCsv   ParserType = "csv"
	ParserTypeDelim ParserType = "delim"
	ParserTypeElff  ParserType = "elff"
)

func (e ParserType) ToPointer() *ParserType {
	return &e
}
func (e *ParserType) UnmarshalJSON(data []byte) error {
	var v string
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	switch v {
	case "clf":
		fallthrough
	case "csv":
		fallthrough
	case "delim":
		fallthrough
	case "elff":
		*e = ParserType(v)
		return nil
	default:
		return fmt.Errorf("invalid value for ParserType: %v", v)
	}
}

type Parser struct {
	AllowedKeyChars   *bool      `json:"allowedKeyChars,omitempty"`
	AllowedValueChars []string   `json:"allowedValueChars,omitempty"`
	CleanFields       []string   `json:"cleanFields,omitempty"`
	DelimChar         *string    `json:"delimChar,omitempty"`
	DstField          *string    `json:"dstField,omitempty"`
	EscapeChar        *string    `json:"escapeChar,omitempty"`
	FieldFilterExpr   *string    `json:"fieldFilterExpr,omitempty"`
	Fields            []string   `json:"fields,omitempty"`
	Keep              []string   `json:"keep,omitempty"`
	Mode              ParserMode `json:"mode"`
	NullValue         *string    `json:"nullValue,omitempty"`
	QuoteChar         *string    `json:"quoteChar,omitempty"`
	Remove            []string   `json:"remove,omitempty"`
	SrcField          string     `json:"srcField"`
	Type              ParserType `json:"type"`
}

func (o *Parser) GetAllowedKeyChars() *bool {
	if o == nil {
		return nil
	}
	return o.AllowedKeyChars
}

func (o *Parser) GetAllowedValueChars() []string {
	if o == nil {
		return nil
	}
	return o.AllowedValueChars
}

func (o *Parser) GetCleanFields() []string {
	if o == nil {
		return nil
	}
	return o.CleanFields
}

func (o *Parser) GetDelimChar() *string {
	if o == nil {
		return nil
	}
	return o.DelimChar
}

func (o *Parser) GetDstField() *string {
	if o == nil {
		return nil
	}
	return o.DstField
}

func (o *Parser) GetEscapeChar() *string {
	if o == nil {
		return nil
	}
	return o.EscapeChar
}

func (o *Parser) GetFieldFilterExpr() *string {
	if o == nil {
		return nil
	}
	return o.FieldFilterExpr
}

func (o *Parser) GetFields() []string {
	if o == nil {
		return nil
	}
	return o.Fields
}

func (o *Parser) GetKeep() []string {
	if o == nil {
		return nil
	}
	return o.Keep
}

func (o *Parser) GetMode() ParserMode {
	if o == nil {
		return ParserMode("")
	}
	return o.Mode
}

func (o *Parser) GetNullValue() *string {
	if o == nil {
		return nil
	}
	return o.NullValue
}

func (o *Parser) GetQuoteChar() *string {
	if o == nil {
		return nil
	}
	return o.QuoteChar
}

func (o *Parser) GetRemove() []string {
	if o == nil {
		return nil
	}
	return o.Remove
}

func (o *Parser) GetSrcField() string {
	if o == nil {
		return ""
	}
	return o.SrcField
}

func (o *Parser) GetType() ParserType {
	if o == nil {
		return ParserType("")
	}
	return o.Type
}

type ParserUnionType string

const (
	ParserUnionTypeParser     ParserUnionType = "parser"
	ParserUnionTypeParserGrok ParserUnionType = "parser_Grok"
	ParserUnionTypeParserJSON ParserUnionType = "parser_JSON"
	ParserUnionTypeParserKvp  ParserUnionType = "parser_Kvp"
)

type ParserUnion struct {
	Parser     *Parser     `queryParam:"inline"`
	ParserGrok *ParserGrok `queryParam:"inline"`
	ParserJSON *ParserJSON `queryParam:"inline"`
	ParserKvp  *ParserKvp  `queryParam:"inline"`

	Type ParserUnionType
}

func CreateParserUnionParser(parser Parser) ParserUnion {
	typ := ParserUnionTypeParser

	return ParserUnion{
		Parser: &parser,
		Type:   typ,
	}
}

func CreateParserUnionParserGrok(parserGrok ParserGrok) ParserUnion {
	typ := ParserUnionTypeParserGrok

	return ParserUnion{
		ParserGrok: &parserGrok,
		Type:       typ,
	}
}

func CreateParserUnionParserJSON(parserJSON ParserJSON) ParserUnion {
	typ := ParserUnionTypeParserJSON

	return ParserUnion{
		ParserJSON: &parserJSON,
		Type:       typ,
	}
}

func CreateParserUnionParserKvp(parserKvp ParserKvp) ParserUnion {
	typ := ParserUnionTypeParserKvp

	return ParserUnion{
		ParserKvp: &parserKvp,
		Type:      typ,
	}
}

func (u *ParserUnion) UnmarshalJSON(data []byte) error {

	var parserGrok ParserGrok = ParserGrok{}
	if err := utils.UnmarshalJSON(data, &parserGrok, "", true, true); err == nil {
		u.ParserGrok = &parserGrok
		u.Type = ParserUnionTypeParserGrok
		return nil
	}

	var parserJSON ParserJSON = ParserJSON{}
	if err := utils.UnmarshalJSON(data, &parserJSON, "", true, true); err == nil {
		u.ParserJSON = &parserJSON
		u.Type = ParserUnionTypeParserJSON
		return nil
	}

	var parser Parser = Parser{}
	if err := utils.UnmarshalJSON(data, &parser, "", true, true); err == nil {
		u.Parser = &parser
		u.Type = ParserUnionTypeParser
		return nil
	}

	var parserKvp ParserKvp = ParserKvp{}
	if err := utils.UnmarshalJSON(data, &parserKvp, "", true, true); err == nil {
		u.ParserKvp = &parserKvp
		u.Type = ParserUnionTypeParserKvp
		return nil
	}

	return fmt.Errorf("could not unmarshal `%s` into any supported union types for ParserUnion", string(data))
}

func (u ParserUnion) MarshalJSON() ([]byte, error) {
	if u.Parser != nil {
		return utils.MarshalJSON(u.Parser, "", true)
	}

	if u.ParserGrok != nil {
		return utils.MarshalJSON(u.ParserGrok, "", true)
	}

	if u.ParserJSON != nil {
		return utils.MarshalJSON(u.ParserJSON, "", true)
	}

	if u.ParserKvp != nil {
		return utils.MarshalJSON(u.ParserKvp, "", true)
	}

	return nil, errors.New("could not marshal union type ParserUnion: all fields are null")
}

type EventBreakerRuleTimestampType string

const (
	EventBreakerRuleTimestampTypeAuto    EventBreakerRuleTimestampType = "auto"
	EventBreakerRuleTimestampTypeFormat  EventBreakerRuleTimestampType = "format"
	EventBreakerRuleTimestampTypeCurrent EventBreakerRuleTimestampType = "current"
)

func (e EventBreakerRuleTimestampType) ToPointer() *EventBreakerRuleTimestampType {
	return &e
}
func (e *EventBreakerRuleTimestampType) UnmarshalJSON(data []byte) error {
	var v string
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	switch v {
	case "auto":
		fallthrough
	case "format":
		fallthrough
	case "current":
		*e = EventBreakerRuleTimestampType(v)
		return nil
	default:
		return fmt.Errorf("invalid value for EventBreakerRuleTimestampType: %v", v)
	}
}

type Timestamp struct {
	Format *string                       `json:"format,omitempty"`
	Length *float64                      `json:"length,omitempty"`
	Type   EventBreakerRuleTimestampType `json:"type"`
}

func (o *Timestamp) GetFormat() *string {
	if o == nil {
		return nil
	}
	return o.Format
}

func (o *Timestamp) GetLength() *float64 {
	if o == nil {
		return nil
	}
	return o.Length
}

func (o *Timestamp) GetType() EventBreakerRuleTimestampType {
	if o == nil {
		return EventBreakerRuleTimestampType("")
	}
	return o.Type
}

type TimestampTimezone struct {
	Name    string    `json:"name"`
	Offsets []float64 `json:"offsets"`
	Untils  []float64 `json:"untils"`
}

func (o *TimestampTimezone) GetName() string {
	if o == nil {
		return ""
	}
	return o.Name
}

func (o *TimestampTimezone) GetOffsets() []float64 {
	if o == nil {
		return []float64{}
	}
	return o.Offsets
}

func (o *TimestampTimezone) GetUntils() []float64 {
	if o == nil {
		return []float64{}
	}
	return o.Untils
}

type TimestampTimezoneUnionType string

const (
	TimestampTimezoneUnionTypeStr               TimestampTimezoneUnionType = "str"
	TimestampTimezoneUnionTypeTimestampTimezone TimestampTimezoneUnionType = "timestampTimezone"
)

type TimestampTimezoneUnion struct {
	Str               *string            `queryParam:"inline"`
	TimestampTimezone *TimestampTimezone `queryParam:"inline"`

	Type TimestampTimezoneUnionType
}

func CreateTimestampTimezoneUnionStr(str string) TimestampTimezoneUnion {
	typ := TimestampTimezoneUnionTypeStr

	return TimestampTimezoneUnion{
		Str:  &str,
		Type: typ,
	}
}

func CreateTimestampTimezoneUnionTimestampTimezone(timestampTimezone TimestampTimezone) TimestampTimezoneUnion {
	typ := TimestampTimezoneUnionTypeTimestampTimezone

	return TimestampTimezoneUnion{
		TimestampTimezone: &timestampTimezone,
		Type:              typ,
	}
}

func (u *TimestampTimezoneUnion) UnmarshalJSON(data []byte) error {

	var timestampTimezone TimestampTimezone = TimestampTimezone{}
	if err := utils.UnmarshalJSON(data, &timestampTimezone, "", true, true); err == nil {
		u.TimestampTimezone = &timestampTimezone
		u.Type = TimestampTimezoneUnionTypeTimestampTimezone
		return nil
	}

	var str string = ""
	if err := utils.UnmarshalJSON(data, &str, "", true, true); err == nil {
		u.Str = &str
		u.Type = TimestampTimezoneUnionTypeStr
		return nil
	}

	return fmt.Errorf("could not unmarshal `%s` into any supported union types for TimestampTimezoneUnion", string(data))
}

func (u TimestampTimezoneUnion) MarshalJSON() ([]byte, error) {
	if u.Str != nil {
		return utils.MarshalJSON(u.Str, "", true)
	}

	if u.TimestampTimezone != nil {
		return utils.MarshalJSON(u.TimestampTimezone, "", true)
	}

	return nil, errors.New("could not marshal union type TimestampTimezoneUnion: all fields are null")
}

type EventBreakerRuleType string

const (
	EventBreakerRuleTypeRegex         EventBreakerRuleType = "regex"
	EventBreakerRuleTypeTimestamp     EventBreakerRuleType = "timestamp"
	EventBreakerRuleTypeJSON          EventBreakerRuleType = "json"
	EventBreakerRuleTypeCsv           EventBreakerRuleType = "csv"
	EventBreakerRuleTypeJSONArray     EventBreakerRuleType = "json_array"
	EventBreakerRuleTypeHeader        EventBreakerRuleType = "header"
	EventBreakerRuleTypeAwsCloudtrail EventBreakerRuleType = "aws_cloudtrail"
	EventBreakerRuleTypeAwsVpcflow    EventBreakerRuleType = "aws_vpcflow"
)

func (e EventBreakerRuleType) ToPointer() *EventBreakerRuleType {
	return &e
}
func (e *EventBreakerRuleType) UnmarshalJSON(data []byte) error {
	var v string
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	switch v {
	case "regex":
		fallthrough
	case "timestamp":
		fallthrough
	case "json":
		fallthrough
	case "csv":
		fallthrough
	case "json_array":
		fallthrough
	case "header":
		fallthrough
	case "aws_cloudtrail":
		fallthrough
	case "aws_vpcflow":
		*e = EventBreakerRuleType(v)
		return nil
	default:
		return fmt.Errorf("invalid value for EventBreakerRuleType: %v", v)
	}
}

type EventBreakerRule struct {
	CleanFields          *bool                    `json:"cleanFields,omitempty"`
	Condition            string                   `json:"condition"`
	Delimiter            *string                  `json:"delimiter,omitempty"`
	DelimiterRegex       *string                  `json:"delimiterRegex,omitempty"`
	Disabled             *bool                    `json:"disabled,omitempty"`
	EscapeChar           *string                  `json:"escapeChar,omitempty"`
	EventBreakerRegex    *string                  `default:"/^/" json:"eventBreakerRegex"`
	Fields               []EventBreakerRuleFields `json:"fields,omitempty"`
	FieldsLineRegex      *string                  `json:"fieldsLineRegex,omitempty"`
	HeaderLineRegex      *string                  `json:"headerLineRegex,omitempty"`
	Index                *float64                 `json:"index,omitempty"`
	JSONArrayField       *string                  `json:"jsonArrayField,omitempty"`
	JSONExtractAll       *bool                    `json:"jsonExtractAll,omitempty"`
	JSONTimeField        *string                  `json:"jsonTimeField,omitempty"`
	MaxEventBytes        float64                  `json:"maxEventBytes"`
	Name                 string                   `json:"name"`
	NullFieldVal         *string                  `json:"nullFieldVal,omitempty"`
	Parser               *ParserUnion             `json:"parser,omitempty"`
	ParserEnabled        *bool                    `json:"parserEnabled,omitempty"`
	QuoteChar            *string                  `json:"quoteChar,omitempty"`
	ShouldUseDataRaw     *bool                    `json:"shouldUseDataRaw,omitempty"`
	TimeField            *string                  `json:"timeField,omitempty"`
	Timestamp            Timestamp                `json:"timestamp"`
	TimestampAnchorRegex string                   `json:"timestampAnchorRegex"`
	TimestampEarliest    *string                  `json:"timestampEarliest,omitempty"`
	TimestampLatest      *string                  `json:"timestampLatest,omitempty"`
	TimestampTimezone    TimestampTimezoneUnion   `json:"timestampTimezone"`
	Type                 *EventBreakerRuleType    `json:"type,omitempty"`
}

func (e EventBreakerRule) MarshalJSON() ([]byte, error) {
	return utils.MarshalJSON(e, "", false)
}

func (e *EventBreakerRule) UnmarshalJSON(data []byte) error {
	if err := utils.UnmarshalJSON(data, &e, "", false, false); err != nil {
		return err
	}
	return nil
}

func (o *EventBreakerRule) GetCleanFields() *bool {
	if o == nil {
		return nil
	}
	return o.CleanFields
}

func (o *EventBreakerRule) GetCondition() string {
	if o == nil {
		return ""
	}
	return o.Condition
}

func (o *EventBreakerRule) GetDelimiter() *string {
	if o == nil {
		return nil
	}
	return o.Delimiter
}

func (o *EventBreakerRule) GetDelimiterRegex() *string {
	if o == nil {
		return nil
	}
	return o.DelimiterRegex
}

func (o *EventBreakerRule) GetDisabled() *bool {
	if o == nil {
		return nil
	}
	return o.Disabled
}

func (o *EventBreakerRule) GetEscapeChar() *string {
	if o == nil {
		return nil
	}
	return o.EscapeChar
}

func (o *EventBreakerRule) GetEventBreakerRegex() *string {
	if o == nil {
		return nil
	}
	return o.EventBreakerRegex
}

func (o *EventBreakerRule) GetFields() []EventBreakerRuleFields {
	if o == nil {
		return nil
	}
	return o.Fields
}

func (o *EventBreakerRule) GetFieldsLineRegex() *string {
	if o == nil {
		return nil
	}
	return o.FieldsLineRegex
}

func (o *EventBreakerRule) GetHeaderLineRegex() *string {
	if o == nil {
		return nil
	}
	return o.HeaderLineRegex
}

func (o *EventBreakerRule) GetIndex() *float64 {
	if o == nil {
		return nil
	}
	return o.Index
}

func (o *EventBreakerRule) GetJSONArrayField() *string {
	if o == nil {
		return nil
	}
	return o.JSONArrayField
}

func (o *EventBreakerRule) GetJSONExtractAll() *bool {
	if o == nil {
		return nil
	}
	return o.JSONExtractAll
}

func (o *EventBreakerRule) GetJSONTimeField() *string {
	if o == nil {
		return nil
	}
	return o.JSONTimeField
}

func (o *EventBreakerRule) GetMaxEventBytes() float64 {
	if o == nil {
		return 0.0
	}
	return o.MaxEventBytes
}

func (o *EventBreakerRule) GetName() string {
	if o == nil {
		return ""
	}
	return o.Name
}

func (o *EventBreakerRule) GetNullFieldVal() *string {
	if o == nil {
		return nil
	}
	return o.NullFieldVal
}

func (o *EventBreakerRule) GetParser() *ParserUnion {
	if o == nil {
		return nil
	}
	return o.Parser
}

func (o *EventBreakerRule) GetParserEnabled() *bool {
	if o == nil {
		return nil
	}
	return o.ParserEnabled
}

func (o *EventBreakerRule) GetQuoteChar() *string {
	if o == nil {
		return nil
	}
	return o.QuoteChar
}

func (o *EventBreakerRule) GetShouldUseDataRaw() *bool {
	if o == nil {
		return nil
	}
	return o.ShouldUseDataRaw
}

func (o *EventBreakerRule) GetTimeField() *string {
	if o == nil {
		return nil
	}
	return o.TimeField
}

func (o *EventBreakerRule) GetTimestamp() Timestamp {
	if o == nil {
		return Timestamp{}
	}
	return o.Timestamp
}

func (o *EventBreakerRule) GetTimestampAnchorRegex() string {
	if o == nil {
		return ""
	}
	return o.TimestampAnchorRegex
}

func (o *EventBreakerRule) GetTimestampEarliest() *string {
	if o == nil {
		return nil
	}
	return o.TimestampEarliest
}

func (o *EventBreakerRule) GetTimestampLatest() *string {
	if o == nil {
		return nil
	}
	return o.TimestampLatest
}

func (o *EventBreakerRule) GetTimestampTimezone() TimestampTimezoneUnion {
	if o == nil {
		return TimestampTimezoneUnion{}
	}
	return o.TimestampTimezone
}

func (o *EventBreakerRule) GetType() *EventBreakerRuleType {
	if o == nil {
		return nil
	}
	return o.Type
}
