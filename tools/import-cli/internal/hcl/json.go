package hcl

import (
	"bytes"
	"encoding/json"
	"strings"
)

// MarshalJSONNoEscape returns compact JSON without escaping HTML characters.
func MarshalJSONNoEscape(v any) ([]byte, error) {
	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	enc.SetEscapeHTML(false)
	if err := enc.Encode(v); err != nil {
		return nil, err
	}
	return []byte(strings.TrimSuffix(buf.String(), "\n")), nil
}

func normalizeJSONNoEscape(raw string) (string, error) {
	var v any
	if err := json.Unmarshal([]byte(raw), &v); err != nil {
		return "", err
	}
	out, err := MarshalJSONNoEscape(v)
	if err != nil {
		return "", err
	}
	return string(out), nil
}
