package cmd

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/criblio/terraform-provider-criblio/internal/restclient"
	"github.com/stretchr/testify/assert"
)

func TestIsContextualDiscoveryError(t *testing.T) {
	tests := []struct {
		name      string
		err       error
		migration bool
		expected  bool
	}{
		{
			name:     "forbidden resource",
			err:      fmt.Errorf("mapping ruleset: %w", &restclient.HTTPError{StatusCode: http.StatusForbidden}),
			expected: true,
		},
		{
			name:      "license prohibited during migration",
			err:       &restclient.HTTPError{StatusCode: http.StatusInternalServerError, Body: `{"message":"The current license prohibits notifications"}`},
			migration: true,
			expected:  true,
		},
		{
			name:      "missing endpoint during migration",
			err:       &restclient.NotFoundError{Path: "/newer/endpoint"},
			migration: true,
			expected:  true,
		},
		{
			name:      "unexpected server failure",
			err:       &restclient.HTTPError{StatusCode: http.StatusInternalServerError, Body: `{"message":"internal error"}`},
			migration: true,
			expected:  false,
		},
		{
			name:      "unauthorized remains fatal",
			err:       &restclient.HTTPError{StatusCode: http.StatusUnauthorized},
			migration: true,
			expected:  false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert.Equal(t, test.expected, isContextualDiscoveryError(test.err, test.migration))
		})
	}
}
