package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/negipo/sunaba_code/pkg/sandbox"
	"github.com/spf13/cobra"
)

var (
	writablePaths []string
	networkAccess bool
	allowedPorts  []int
)

var rootCmd = &cobra.Command{
	Use:   "sunaba_code [flags] -- [claude command and args]",
	Short: "Run Claude Code in a sandboxed environment",
	Long: `sunaba_code runs Claude Code in a sandboxed environment using macOS sandbox-exec.
By default, only the current directory is writable, and no network access is allowed.

Examples:
  # Run Claude Code with default settings (current dir writable)
  sunaba_code -- claude

  # Allow writing to specific directories
  sunaba_code -w ~/projects -w /tmp -- claude

  # Enable network access
  sunaba_code --network -- claude`,
	Args: cobra.MinimumNArgs(1),
	RunE: runCommand,
}

func init() {
	// Get current working directory for default writable path
	cwd, err := os.Getwd()
	if err != nil {
		cwd = "."
	}

	rootCmd.Flags().StringSliceVarP(&writablePaths, "writable", "w", []string{cwd}, "Paths where write access is allowed")
	rootCmd.Flags().BoolVar(&networkAccess, "network", false, "Allow network access")
	rootCmd.Flags().IntSliceVar(&allowedPorts, "ports", []int{}, "Specific ports to allow (requires --network)")
}

func runCommand(cmd *cobra.Command, args []string) error {
	// Validate flags
	if len(allowedPorts) > 0 && !networkAccess {
		return fmt.Errorf("--ports requires --network flag")
	}

	// Expand paths to absolute paths
	expandedPaths := make([]string, len(writablePaths))
	for i, path := range writablePaths {
		absPath, err := filepath.Abs(os.ExpandEnv(path))
		if err != nil {
			return fmt.Errorf("invalid path %s: %w", path, err)
		}
		expandedPaths[i] = absPath
	}

	// Create sandbox configuration
	config := sandbox.ProfileConfig{
		WritablePaths: expandedPaths,
		NetworkAccess: networkAccess,
		AllowedPorts:  allowedPorts,
	}

	// Create executor
	executor, err := sandbox.NewExecutor(config)
	if err != nil {
		return fmt.Errorf("failed to create sandbox executor: %w", err)
	}
	defer executor.Cleanup()

	// Execute the command
	command := args[0]
	commandArgs := args[1:]

	fmt.Fprintf(os.Stderr, "Running '%s' in sandbox with write access to: %v\n", command, expandedPaths)
	if networkAccess {
		if len(allowedPorts) > 0 {
			fmt.Fprintf(os.Stderr, "Network access enabled for ports: %v\n", allowedPorts)
		} else {
			fmt.Fprintf(os.Stderr, "Full network access enabled\n")
		}
	}

	return executor.Execute(command, commandArgs, os.Environ())
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}