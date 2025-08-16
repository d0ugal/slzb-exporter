package version

import "runtime"

var (
	// Version is the version of the application
	Version = "dev"
	// Commit is the git commit hash
	Commit = "unknown"
	// BuildDate is the build date
	BuildDate = "unknown"
)

// Info contains version information
type Info struct {
	Version   string `json:"version"`
	Commit    string `json:"commit"`
	BuildDate string `json:"build_date"`
	GoVersion string `json:"go_version"`
}

// Get returns the version information
func Get() Info {
	return Info{
		Version:   Version,
		Commit:    Commit,
		BuildDate: BuildDate,
		GoVersion: runtime.Version(),
	}
}
