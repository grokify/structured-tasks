package main

import (
	"github.com/spf13/cobra"
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

func init() {
	rootCmd.AddCommand(validateCmd)
	rootCmd.AddCommand(generateCmd)
	rootCmd.AddCommand(statsCmd)
	rootCmd.AddCommand(depsCmd)
}
