package sandbox

import (
	"fmt"
	"os"
	"path/filepath"
)

// ProfileConfig represents the configuration for a sandbox profile
type ProfileConfig struct {
	// WritablePaths are paths where write access is allowed
	WritablePaths []string
	// NetworkAccess specifies if network access is allowed
	NetworkAccess bool
	// AllowedPorts are the ports that can be accessed (if NetworkAccess is true)
	AllowedPorts []int
}

// GenerateProfile creates a sandbox profile for macOS sandbox-exec
func GenerateProfile(config ProfileConfig) string {
	profile := `(version 1)
(deny default)
(allow sysctl-read)
(allow process-fork)
(allow process-exec)
(allow file-read*)
(allow file-read-metadata)
`

	// Add network access if enabled
	if config.NetworkAccess {
		if len(config.AllowedPorts) > 0 {
			for _, port := range config.AllowedPorts {
				profile += fmt.Sprintf("(allow network* (remote ip \"localhost:%d\"))\n", port)
			}
		} else {
			profile += "(allow network*)\n"
		}
	}

	// Add writable paths
	for _, path := range config.WritablePaths {
		absPath, _ := filepath.Abs(path)
		profile += fmt.Sprintf("(allow file-write* (subpath \"%s\"))\n", absPath)
	}

	// Explicitly deny all other writes
	profile += "(deny file-write*)\n"

	return profile
}

// WriteProfileToFile writes a sandbox profile to a file
func WriteProfileToFile(profile string, filename string) error {
	return os.WriteFile(filename, []byte(profile), 0644)
}