package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

// Version information set by ldflags during build.
var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

var rootCmd = &cobra.Command{
	Use:   "stasks",
	Short: "Structured Tasks - machine-readable project task lists",
	Long: `stasks is a CLI tool for managing structured task lists.

It provides commands to validate, generate, and analyze TASKS.json files,
producing deterministic TASKS.md output.

Example usage:
  stasks validate TASKS.json
  stasks generate -i TASKS.json -o TASKS.md
  stasks stats TASKS.json`,
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print version information",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("stasks %s\n", version)
		fmt.Printf("  commit: %s\n", commit)
		fmt.Printf("  built:  %s\n", date)
	},
}

func init() {
	rootCmd.AddCommand(validateCmd)
	rootCmd.AddCommand(generateCmd)
	rootCmd.AddCommand(statsCmd)
	rootCmd.AddCommand(depsCmd)
	rootCmd.AddCommand(versionCmd)
}
