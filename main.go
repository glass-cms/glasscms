package main

import (
	"github.com/glass-cms/glasscms/cmd"
	"github.com/glass-cms/glasscms/internal/version"
)

// Build-time variables that will be injected during build.
var (
	Version = "dev"
	Commit  = "unknown"
	Date    = "unknown"
)

func main() {
	// Set version information from build-time variables
	version.Version = Version
	version.Commit = Commit
	version.Date = Date

	// Validate semantic version (will crash if invalid and not "dev")
	version.ValidateVersion()

	cmd.Execute()
}
