package main

import (
	"fmt"

	"github.com/grokify/structured-roadmap/roadmap"
	"github.com/spf13/cobra"
)

var validateCmd = &cobra.Command{
	Use:   "validate <file>",
	Short: "Validate a ROADMAP.json file",
	Long:  `Validate a ROADMAP.json file against the schema and check for errors.`,
	Args:  cobra.ExactArgs(1),
	RunE:  runValidate,
}

func runValidate(cmd *cobra.Command, args []string) error {
	path := args[0]

	r, err := roadmap.ParseFile(path)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	result := roadmap.Validate(r)

	if result.Valid {
		fmt.Fprintf(cmd.ErrOrStderr(), "✅ %s is valid\n", path)
		fmt.Fprintf(cmd.ErrOrStderr(), "   Project: %s\n", r.Project)
		fmt.Fprintf(cmd.ErrOrStderr(), "   Items: %d\n", len(r.Items))
		fmt.Fprintf(cmd.ErrOrStderr(), "   Phases: %d\n", len(r.Phases))
		fmt.Fprintf(cmd.ErrOrStderr(), "   Areas: %d\n", len(r.Areas))
		return nil
	}

	fmt.Fprintf(cmd.ErrOrStderr(), "❌ %s has %d error(s)\n\n", path, len(result.Errors))
	for _, e := range result.Errors {
		fmt.Fprintf(cmd.ErrOrStderr(), "  • %s: %s\n", e.Field, e.Message)
	}
	return fmt.Errorf("validation failed with %d error(s)", len(result.Errors))
}
