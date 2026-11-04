package provider

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"net/http"
	"os"
	"strconv"
)

func (data CriblioProviderModel) httpTransport() (*http.Transport, error) {
	if data.InsecureSkipVerify.IsUnknown() || data.CACertFile.IsUnknown() {
		return nil, fmt.Errorf("insecure_skip_verify and ca_cert_file must be known before configuring the provider")
	}

	skipVerify := data.InsecureSkipVerify.ValueBool()
	if data.InsecureSkipVerify.IsNull() {
		if value := os.Getenv("CRIBL_INSECURE_SKIP_VERIFY"); value != "" {
			var err error
			skipVerify, err = strconv.ParseBool(value)
			if err != nil {
				return nil, fmt.Errorf("CRIBL_INSECURE_SKIP_VERIFY must be a boolean: %w", err)
			}
		}
	}

	caFile := data.CACertFile.ValueString()
	if data.CACertFile.IsNull() {
		caFile = os.Getenv("CRIBL_CA_CERT_FILE")
	}

	base, ok := http.DefaultTransport.(*http.Transport)
	if !ok {
		return nil, fmt.Errorf("default HTTP transport must be an *http.Transport")
	}
	transport := base.Clone()
	if transport.TLSClientConfig == nil {
		transport.TLSClientConfig = &tls.Config{}
	}
	transport.TLSClientConfig.InsecureSkipVerify = skipVerify

	if caFile != "" {
		pem, err := os.ReadFile(caFile)
		if err != nil {
			return nil, fmt.Errorf("read ca_cert_file %q: %w", caFile, err)
		}
		roots, err := x509.SystemCertPool()
		if err != nil {
			return nil, fmt.Errorf("load system CA certificates: %w", err)
		}
		if !roots.AppendCertsFromPEM(pem) {
			return nil, fmt.Errorf("ca_cert_file %q contains no valid PEM certificates", caFile)
		}
		transport.TLSClientConfig.RootCAs = roots
	}
	return transport, nil
}
