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
	Use:   "sroadmap",
	Short: "Structured Roadmap - machine-readable project roadmaps",
	Long: `sroadmap is a CLI tool for managing structured roadmaps.

It provides commands to validate, generate, and analyze ROADMAP.json files,
producing deterministic ROADMAP.md output.

Example usage:
  sroadmap validate ROADMAP.json
  sroadmap generate -i ROADMAP.json -o ROADMAP.md
  sroadmap stats ROADMAP.json`,
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print version information",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("sroadmap %s\n", version)
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
