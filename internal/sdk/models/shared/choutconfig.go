// Code generated by Speakeasy (https://speakeasy.com). DO NOT EDIT.

package shared

import (
	"encoding/json"
	"fmt"
)

type CHOutConfigAuthType string

const (
	CHOutConfigAuthTypeToken              CHOutConfigAuthType = "token"
	CHOutConfigAuthTypeNone               CHOutConfigAuthType = "none"
	CHOutConfigAuthTypeTextSecret         CHOutConfigAuthType = "textSecret"
	CHOutConfigAuthTypeBasic              CHOutConfigAuthType = "basic"
	CHOutConfigAuthTypeCredentialsSecret  CHOutConfigAuthType = "credentialsSecret"
	CHOutConfigAuthTypeSecret             CHOutConfigAuthType = "secret"
	CHOutConfigAuthTypeManual             CHOutConfigAuthType = "manual"
	CHOutConfigAuthTypeManualAPIKey       CHOutConfigAuthType = "manualAPIKey"
	CHOutConfigAuthTypeSslUserCertificate CHOutConfigAuthType = "sslUserCertificate"
)

func (e CHOutConfigAuthType) ToPointer() *CHOutConfigAuthType {
	return &e
}
func (e *CHOutConfigAuthType) UnmarshalJSON(data []byte) error {
	var v string
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	switch v {
	case "token":
		fallthrough
	case "none":
		fallthrough
	case "textSecret":
		fallthrough
	case "basic":
		fallthrough
	case "credentialsSecret":
		fallthrough
	case "secret":
		fallthrough
	case "manual":
		fallthrough
	case "manualAPIKey":
		fallthrough
	case "sslUserCertificate":
		*e = CHOutConfigAuthType(v)
		return nil
	default:
		return fmt.Errorf("invalid value for CHOutConfigAuthType: %v", v)
	}
}

type CHOutConfigColumnMapping struct {
	ColumnName            string `json:"columnName"`
	ColumnType            string `json:"columnType"`
	ColumnValueExpression string `json:"columnValueExpression"`
}

func (o *CHOutConfigColumnMapping) GetColumnName() string {
	if o == nil {
		return ""
	}
	return o.ColumnName
}

func (o *CHOutConfigColumnMapping) GetColumnType() string {
	if o == nil {
		return ""
	}
	return o.ColumnType
}

func (o *CHOutConfigColumnMapping) GetColumnValueExpression() string {
	if o == nil {
		return ""
	}
	return o.ColumnValueExpression
}

type CHOutConfig struct {
	AsyncInserts                  bool                         `json:"asyncInserts"`
	Auth                          *HTTPOutAuthConfig           `json:"auth,omitempty"`
	AuthType                      *CHOutConfigAuthType         `json:"authType,omitempty"`
	ColumnMappings                []CHOutConfigColumnMapping   `json:"columnMappings,omitempty"`
	Compress                      *bool                        `json:"compress,omitempty"`
	Concurrency                   *float64                     `json:"concurrency,omitempty"`
	Database                      string                       `json:"database"`
	DumpFormatErrorsToDisk        *bool                        `json:"dumpFormatErrorsToDisk,omitempty"`
	ExcludeMappingFields          []string                     `json:"excludeMappingFields,omitempty"`
	ExtraHTTPHeaders              []NameValue                  `json:"extraHttpHeaders,omitempty"`
	ExtraParams                   []HTTPOutExtraParamConfig    `json:"extraParams,omitempty"`
	FailedRequestLoggingMode      *string                      `json:"failedRequestLoggingMode,omitempty"`
	FlushPeriodSec                float64                      `json:"flushPeriodSec"`
	Format                        FormatEnum                   `json:"format"`
	KeepAlive                     *bool                        `json:"keepAlive,omitempty"`
	LoadBalanced                  bool                         `json:"loadBalanced"`
	MappingType                   MappingType                  `json:"mappingType"`
	MaxConnectionReuseSec         *float64                     `json:"maxConnectionReuseSec,omitempty"`
	MaxPayloadEvents              *float64                     `json:"maxPayloadEvents,omitempty"`
	MaxPayloadSizeKB              *float64                     `json:"maxPayloadSizeKB,omitempty"`
	MaxSockets                    *float64                     `json:"maxSockets,omitempty"`
	Method                        *string                      `json:"method,omitempty"`
	Password                      *string                      `json:"password,omitempty"`
	RejectUnauthorized            *bool                        `json:"rejectUnauthorized,omitempty"`
	ResponseHonorRetryAfterHeader *bool                        `json:"responseHonorRetryAfterHeader,omitempty"`
	ResponseRetrySettings         []HTTPOutResponseRetryConfig `json:"responseRetrySettings,omitempty"`
	SafeHeaders                   []string                     `json:"safeHeaders,omitempty"`
	SQLUsername                   *string                      `json:"sqlUsername,omitempty"`
	TableName                     string                       `json:"tableName"`
	TableNameExpression           *string                      `json:"tableNameExpression,omitempty"`
	TimeoutRetrySettings          *RetryBackoffOptions         `json:"timeoutRetrySettings,omitempty"`
	TimeoutSec                    *float64                     `json:"timeoutSec,omitempty"`
	TLS                           *TLSClientParams             `json:"tls,omitempty"`
	Token                         *string                      `json:"token,omitempty"`
	TotalMemoryLimitKB            *float64                     `json:"totalMemoryLimitKB,omitempty"`
	URL                           string                       `json:"url"`
	Urls                          []string                     `json:"urls,omitempty"`
	UseRoundRobinDNS              *bool                        `json:"useRoundRobinDns,omitempty"`
	Username                      *string                      `json:"username,omitempty"`
	WaitForAsyncInserts           *bool                        `json:"waitForAsyncInserts,omitempty"`
}

func (o *CHOutConfig) GetAsyncInserts() bool {
	if o == nil {
		return false
	}
	return o.AsyncInserts
}

func (o *CHOutConfig) GetAuth() *HTTPOutAuthConfig {
	if o == nil {
		return nil
	}
	return o.Auth
}

func (o *CHOutConfig) GetAuthType() *CHOutConfigAuthType {
	if o == nil {
		return nil
	}
	return o.AuthType
}

func (o *CHOutConfig) GetColumnMappings() []CHOutConfigColumnMapping {
	if o == nil {
		return nil
	}
	return o.ColumnMappings
}

func (o *CHOutConfig) GetCompress() *bool {
	if o == nil {
		return nil
	}
	return o.Compress
}

func (o *CHOutConfig) GetConcurrency() *float64 {
	if o == nil {
		return nil
	}
	return o.Concurrency
}

func (o *CHOutConfig) GetDatabase() string {
	if o == nil {
		return ""
	}
	return o.Database
}

func (o *CHOutConfig) GetDumpFormatErrorsToDisk() *bool {
	if o == nil {
		return nil
	}
	return o.DumpFormatErrorsToDisk
}

func (o *CHOutConfig) GetExcludeMappingFields() []string {
	if o == nil {
		return nil
	}
	return o.ExcludeMappingFields
}

func (o *CHOutConfig) GetExtraHTTPHeaders() []NameValue {
	if o == nil {
		return nil
	}
	return o.ExtraHTTPHeaders
}

func (o *CHOutConfig) GetExtraParams() []HTTPOutExtraParamConfig {
	if o == nil {
		return nil
	}
	return o.ExtraParams
}

func (o *CHOutConfig) GetFailedRequestLoggingMode() *string {
	if o == nil {
		return nil
	}
	return o.FailedRequestLoggingMode
}

func (o *CHOutConfig) GetFlushPeriodSec() float64 {
	if o == nil {
		return 0.0
	}
	return o.FlushPeriodSec
}

func (o *CHOutConfig) GetFormat() FormatEnum {
	if o == nil {
		return FormatEnum("")
	}
	return o.Format
}

func (o *CHOutConfig) GetKeepAlive() *bool {
	if o == nil {
		return nil
	}
	return o.KeepAlive
}

func (o *CHOutConfig) GetLoadBalanced() bool {
	if o == nil {
		return false
	}
	return o.LoadBalanced
}

func (o *CHOutConfig) GetMappingType() MappingType {
	if o == nil {
		return MappingType("")
	}
	return o.MappingType
}

func (o *CHOutConfig) GetMaxConnectionReuseSec() *float64 {
	if o == nil {
		return nil
	}
	return o.MaxConnectionReuseSec
}

func (o *CHOutConfig) GetMaxPayloadEvents() *float64 {
	if o == nil {
		return nil
	}
	return o.MaxPayloadEvents
}

func (o *CHOutConfig) GetMaxPayloadSizeKB() *float64 {
	if o == nil {
		return nil
	}
	return o.MaxPayloadSizeKB
}

func (o *CHOutConfig) GetMaxSockets() *float64 {
	if o == nil {
		return nil
	}
	return o.MaxSockets
}

func (o *CHOutConfig) GetMethod() *string {
	if o == nil {
		return nil
	}
	return o.Method
}

func (o *CHOutConfig) GetPassword() *string {
	if o == nil {
		return nil
	}
	return o.Password
}

func (o *CHOutConfig) GetRejectUnauthorized() *bool {
	if o == nil {
		return nil
	}
	return o.RejectUnauthorized
}

func (o *CHOutConfig) GetResponseHonorRetryAfterHeader() *bool {
	if o == nil {
		return nil
	}
	return o.ResponseHonorRetryAfterHeader
}

func (o *CHOutConfig) GetResponseRetrySettings() []HTTPOutResponseRetryConfig {
	if o == nil {
		return nil
	}
	return o.ResponseRetrySettings
}

func (o *CHOutConfig) GetSafeHeaders() []string {
	if o == nil {
		return nil
	}
	return o.SafeHeaders
}

func (o *CHOutConfig) GetSQLUsername() *string {
	if o == nil {
		return nil
	}
	return o.SQLUsername
}

func (o *CHOutConfig) GetTableName() string {
	if o == nil {
		return ""
	}
	return o.TableName
}

func (o *CHOutConfig) GetTableNameExpression() *string {
	if o == nil {
		return nil
	}
	return o.TableNameExpression
}

func (o *CHOutConfig) GetTimeoutRetrySettings() *RetryBackoffOptions {
	if o == nil {
		return nil
	}
	return o.TimeoutRetrySettings
}

func (o *CHOutConfig) GetTimeoutSec() *float64 {
	if o == nil {
		return nil
	}
	return o.TimeoutSec
}

func (o *CHOutConfig) GetTLS() *TLSClientParams {
	if o == nil {
		return nil
	}
	return o.TLS
}

func (o *CHOutConfig) GetToken() *string {
	if o == nil {
		return nil
	}
	return o.Token
}

func (o *CHOutConfig) GetTotalMemoryLimitKB() *float64 {
	if o == nil {
		return nil
	}
	return o.TotalMemoryLimitKB
}

func (o *CHOutConfig) GetURL() string {
	if o == nil {
		return ""
	}
	return o.URL
}

func (o *CHOutConfig) GetUrls() []string {
	if o == nil {
		return nil
	}
	return o.Urls
}

func (o *CHOutConfig) GetUseRoundRobinDNS() *bool {
	if o == nil {
		return nil
	}
	return o.UseRoundRobinDNS
}

func (o *CHOutConfig) GetUsername() *string {
	if o == nil {
		return nil
	}
	return o.Username
}

func (o *CHOutConfig) GetWaitForAsyncInserts() *bool {
	if o == nil {
		return nil
	}
	return o.WaitForAsyncInserts
}
