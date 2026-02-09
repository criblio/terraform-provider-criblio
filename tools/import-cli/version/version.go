// Package version holds build-time version and app name, set via ldflags.
package version

// Build-time variables (set via -ldflags).
var (
	Version string
	Commit  string
	Date    string
	AppName string
)

// AppNameOrDefault returns AppName if set, otherwise "goatify".
func AppNameOrDefault() string {
	if AppName != "" {
		return AppName
	}
	return "goatify"
}
