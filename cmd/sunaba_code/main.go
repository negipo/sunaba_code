package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/negipo/sunaba_code/internal/claude"
	"github.com/negipo/sunaba_code/pkg/sandbox"
	"github.com/spf13/cobra"
)

var (
	writablePaths []string
	claudeConfig  bool
)

var rootCmd = &cobra.Command{
	Use:   "sunaba_code [flags] -- [claude command and args]",
	Short: "Run Claude Code in a sandboxed environment",
	Long: `sunaba_code runs Claude Code in a sandboxed environment using macOS sandbox-exec.
By default, only the current directory is writable, and no network access is allowed.

Examples:
  # Run Claude Code with default settings (current dir writable)
  sunaba_code

  # Allow writing to specific directories
  sunaba_code -w ~/projects -w /tmp

  # Enable network access
  sunaba_code --network

  # Run a different command in sandbox
  sunaba_code -- /usr/bin/python script.py`,
	Args: cobra.ArbitraryArgs,
	RunE: runCommand,
}

func init() {
	// Get current working directory for default writable path
	cwd, err := os.Getwd()
	if err != nil {
		cwd = "."
	}

	rootCmd.Flags().StringSliceVarP(&writablePaths, "writable", "w", []string{cwd}, "Paths where write access is allowed")
	rootCmd.Flags().BoolVar(&claudeConfig, "claude-config", true, "Allow access to Claude configuration files")
}

func runCommand(cmd *cobra.Command, args []string) error {

	// Expand paths to absolute paths
	expandedPaths := make([]string, len(writablePaths))
	for i, path := range writablePaths {
		absPath, err := filepath.Abs(os.ExpandEnv(path))
		if err != nil {
			return fmt.Errorf("invalid path %s: %w", path, err)
		}
		expandedPaths[i] = absPath
	}

	// Add /dev/null for output redirection
	expandedPaths = append(expandedPaths, "/dev/null")

	// Add Claude configuration directories if enabled
	if claudeConfig {
		homeDir, err := os.UserHomeDir()
		if err == nil {
			// Add specific Claude-related files and directories for minimal security risk
			claudeDirs := []string{
				filepath.Join(homeDir, ".claude"),
				filepath.Join(homeDir, ".claude.json"),
				filepath.Join(homeDir, ".config", "claude"),
			}
			for _, dir := range claudeDirs {
				expandedPaths = append(expandedPaths, dir)
			}
			
			// Add home directory for Claude's temporary files (required for OAuth and locking)
			// This is a security trade-off needed for Claude Code to function properly
			expandedPaths = append(expandedPaths, homeDir)
		}
	}

	// Create sandbox configuration (network access enabled by default)
	config := sandbox.ProfileConfig{
		WritablePaths: expandedPaths,
		NetworkAccess: true,
		AllowedPorts:  []int{}, // No specific port restrictions
	}

	// Create executor
	executor, err := sandbox.NewExecutor(config)
	if err != nil {
		return fmt.Errorf("failed to create sandbox executor: %w", err)
	}
	defer executor.Cleanup()

	// Handle command - if no command specified, try to find claude
	var command string
	var commandArgs []string

	if len(args) == 0 || args[0] == "claude" {
		// Find claude executable
		claudePath, err := claude.FindClaudeExecutable()
		if err != nil {
			return err
		}
		command = claudePath
		if len(args) > 1 {
			commandArgs = args[1:]
		}
	} else {
		command = args[0]
		commandArgs = args[1:]
	}

	fmt.Fprintf(os.Stderr, "Running '%s' in sandbox with write access to: %v\n", command, expandedPaths)

	return executor.Execute(command, commandArgs, os.Environ())
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}