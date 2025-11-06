package converter

// FieldMapper handles conversion between Cribl API field names and Terraform attribute names

import (
	"strings"
)

// FieldMapper maps Cribl API field names (camelCase) to Terraform attributes (snake_case)
type FieldMapper struct {
	mappings map[string]string
}

// NewFieldMapper creates a new field mapper with common mappings
func NewFieldMapper() *FieldMapper {
	return &FieldMapper{
		mappings: getDefaultMappings(),
	}
}

// MapField converts a field name from Cribl format to Terraform format
func (fm *FieldMapper) MapField(fieldName string) string {
	// Check explicit mappings first
	if mapped, ok := fm.mappings[fieldName]; ok {
		return mapped
	}

	// Auto-convert camelCase to snake_case
	return toSnakeCase(fieldName)
}

// getDefaultMappings returns common field name mappings
// These are based on existing modules in modules/json-config-*/locals.tf
func getDefaultMappings() map[string]string {
	return map[string]string{
		// Common fields
		"maxBufferSize":   "max_buffer_size",
		"maxFileSize":     "max_file_size",
		"maxSize":         "max_size",
		"minVersion":      "min_version",
		"maxVersion":      "max_version",
		"clientId":        "client_id",
		"clientSecret":    "client_secret",
		"accessKeyId":     "access_key_id",
		"secretAccessKey": "secret_access_key",
		
		// TLS fields
		"certPath":           "cert_path",
		"certificateName":    "certificate_name",
		"commonNameRegex":    "common_name_regex",
		"privKeyPath":        "priv_key_path",
		"rejectUnauthorized": "reject_unauthorized",
		"requestCert":        "request_cert",
		
		// PQ (persistent queue) fields
		"commitFrequency": "commit_frequency",
		
		// Network fields
		"bindAddress":  "bind_address",
		"sourceIP":     "source_ip",
		"keepAlive":    "keep_alive",
		"soLinger":     "so_linger",
		"tcpNoDelay":   "tcp_no_delay",
		"writeTimeout": "write_timeout",
		"readTimeout":  "read_timeout",
		
		// Authentication fields
		"apiKey":        "api_key",
		"authToken":     "auth_token",
		"bearerToken":   "bearer_token",
		"oauthClientId": "oauth_client_id",
		
		// AWS fields
		"awsAccessKeyId":     "aws_access_key_id",
		"awsSecretAccessKey": "aws_secret_access_key",
		"awsRegion":          "aws_region",
		"awsRoleArn":         "aws_role_arn",
		
		// Kafka fields
		"brokerTimeout":    "broker_timeout",
		"topicName":        "topic_name",
		"compressionType":  "compression_type",
		"acks":             "acks",
		"batchSize":        "batch_size",
		"lingerMs":         "linger_ms",
		"maxInFlight":      "max_in_flight",
		"requestTimeout":   "request_timeout",
		"enableIdempotence": "enable_idempotence",
		
		// S3 fields
		"bucketName":      "bucket_name",
		"objectKeyPrefix": "object_key_prefix",
		
		// HTTP fields
		"httpTimeout":      "http_timeout",
		"maxRetries":       "max_retries",
		"retryDelay":      "retry_delay",
		"retryOnTimeout":  "retry_on_timeout",
		"contentType":     "content_type",
		"userAgent":       "user_agent",
		"connectTimeout": "connect_timeout",
		
		// JSON/Data fields
		"jsonKey":      "json_key",
		"jsonValue":    "json_value",
		"dataFormat":   "data_format",
		"messageKey":   "message_key",
		"timestampKey": "timestamp_key",
	}
}

// toSnakeCase converts camelCase or PascalCase to snake_case
func toSnakeCase(s string) string {
	if s == "" {
		return s
	}

	var result strings.Builder
	runes := []rune(s)

	for i, r := range runes {
		if i > 0 {
			prev := runes[i-1]
			// Insert underscore before uppercase if previous was lowercase or digit
			if r >= 'A' && r <= 'Z' {
				if (prev >= 'a' && prev <= 'z') || (prev >= '0' && prev <= '9') {
					result.WriteRune('_')
				}
			}
		}
		result.WriteRune(r)
	}

	return strings.ToLower(result.String())
}

