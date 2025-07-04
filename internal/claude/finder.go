package claude

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
)

// FindClaudeExecutable finds the Claude Code executable in the system
func FindClaudeExecutable() (string, error) {
	// Try common names
	candidates := []string{"claude", "claude-code"}

	// First check if it's in PATH
	for _, name := range candidates {
		if path, err := exec.LookPath(name); err == nil {
			return path, nil
		}
	}

	// Check common installation locations based on OS
	var searchPaths []string
	homeDir, _ := os.UserHomeDir()

	switch runtime.GOOS {
	case "darwin":
		searchPaths = []string{
			"/usr/local/bin/claude",
			"/opt/homebrew/bin/claude",
			filepath.Join(homeDir, ".local/bin/claude"),
			filepath.Join(homeDir, "bin/claude"),
		}
	case "linux":
		searchPaths = []string{
			"/usr/local/bin/claude",
			"/usr/bin/claude",
			filepath.Join(homeDir, ".local/bin/claude"),
			filepath.Join(homeDir, "bin/claude"),
		}
	}

	// Try each path
	for _, path := range searchPaths {
		if _, err := os.Stat(path); err == nil {
			return path, nil
		}
	}

	return "", fmt.Errorf("claude executable not found. Please ensure Claude Code is installed and in your PATH")
}