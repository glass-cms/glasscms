package version

import (
	"fmt"
	"log"
	"runtime"

	"github.com/hashicorp/go-version"
)

var (
	// This will be set during build time via -ldflags.
	Version = "dev"

	// This will be set during build time via -ldflags.
	Commit = "unknown"

	// This will be set during build time via -ldflags.
	Date = "unknown"
)

// Info contains version information.
type Info struct {
	Version   string `json:"version"`
	Commit    string `json:"commit"`
	Date      string `json:"date"`
	GoVersion string `json:"go_version"`
	Platform  string `json:"platform"`
}

// Get returns the current version information.
func Get() Info {
	return Info{
		Version:   Version,
		Commit:    Commit,
		Date:      Date,
		GoVersion: runtime.Version(),
		Platform:  fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH),
	}
}

// String returns a human-readable version string with version on first line.
func (i Info) String() string {
	return fmt.Sprintf("%s\ncommit: %s, built: %s, go: %s, platform: %s",
		i.Version, i.Commit, i.Date, i.GoVersion, i.Platform)
}

// If version is not "dev", it must be a valid semantic version or the program will crash.
func ValidateVersion() {
	if Version == "dev" {
		return // Skip validation for dev builds
	}

	// Parse the version using go-version (handles both v1.0.0 and 1.0.0 formats)
	v, err := version.NewVersion(Version)
	if err != nil {
		log.Fatalf("Invalid semantic version '%s': %v", Version, err)
	}

	// go-version normalizes the version, so we need to check if the input format is acceptable
	// It should accept both "v1.0.0" and "1.0.0" formats
	normalizedVersion := v.String()
	versionWithV := "v" + normalizedVersion

	if Version != normalizedVersion && Version != versionWithV {
		log.Fatalf("Version '%s' is not in proper semantic version format. Expected: %s or %s",
			Version, normalizedVersion, versionWithV)
	}
}
