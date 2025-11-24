package hooks

import (
	"strings"
)

func trimPath(path string) string {
	path = strings.TrimLeft(path, "/")
	path = strings.TrimPrefix(path, "api/v1/")
	path = strings.TrimPrefix(path, "api/v1")

	return path
}
