package hooks

import (
	"errors"
	"os"
	"strings"
)

func trimPath(path string) string {
	// Validate that the requested endpoint is supported for on-prem deployments
	path = strings.TrimLeft(path, "/")

	// Remove /api/v1 if already present in path
	path = strings.TrimPrefix(path, "api/v1/")
	path = strings.TrimPrefix(path, "api/v1")
	return path
}
