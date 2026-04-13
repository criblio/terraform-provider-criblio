// Manual extension (not Speakeasy-generated). Override default User-Agent from New().
package sdk

// WithUserAgent sets the User-Agent header for API requests. The default in generated
// New() targets the Terraform provider; callers such as import-cli pass a product-specific value.
func WithUserAgent(userAgent string) SDKOption {
	return func(s *CriblIo) {
		s.sdkConfiguration.UserAgent = userAgent
	}
}
