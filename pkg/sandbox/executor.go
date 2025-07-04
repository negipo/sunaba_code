package sandbox

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

// Executor represents a sandboxed command executor
type Executor struct {
	profilePath string
	config      ProfileConfig
}

// NewExecutor creates a new sandbox executor
func NewExecutor(config ProfileConfig) (*Executor, error) {
	// Create temporary profile file
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get home directory: %w", err)
	}

	profilePath := filepath.Join(homeDir, ".sunaba_code.sandbox.pf")
	profile := GenerateProfile(config)

	if err := WriteProfileToFile(profile, profilePath); err != nil {
		return nil, fmt.Errorf("failed to write sandbox profile: %w", err)
	}

	return &Executor{
		profilePath: profilePath,
		config:      config,
	}, nil
}

// Execute runs a command in the sandbox
func (e *Executor) Execute(command string, args []string, env []string) error {
	// Build sandbox-exec command
	sandboxArgs := []string{
		"-f", e.profilePath,
	}

	// Add sandbox-specific environment variables to prevent Claude from creating lock files
	env = append(env, "CLAUDE_CONFIG_READONLY=1")
	env = append(env, "NODE_ENV=production")

	// Add environment variables
	if len(env) > 0 {
		sandboxArgs = append(sandboxArgs, "env")
		sandboxArgs = append(sandboxArgs, env...)
	}

	// Add the actual command
	sandboxArgs = append(sandboxArgs, command)
	sandboxArgs = append(sandboxArgs, args...)

	cmd := exec.Command("sandbox-exec", sandboxArgs...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

// ExecuteShell runs a shell command in the sandbox
func (e *Executor) ExecuteShell(shellCommand string, env []string) error {
	return e.Execute("/bin/bash", []string{"-c", shellCommand}, env)
}

// Cleanup removes the temporary profile file
func (e *Executor) Cleanup() error {
	return os.Remove(e.profilePath)
}