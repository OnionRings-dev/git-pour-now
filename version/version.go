// Package version holds build information injected from a Git release tag.
package version

var (
	Version   = "dev"
	BuildTime = ""
	GitCommit = ""
)

func GetVersionInfo() string {
	if BuildTime != "" && GitCommit != "" {
		return Version + " (built " + BuildTime + " from " + GitCommit + ")"
	}
	return Version
}
