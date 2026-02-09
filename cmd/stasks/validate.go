package main

import (
	"fmt"

	"github.com/grokify/structured-tasks/tasks"
	"github.com/spf13/cobra"
)

var validateCmd = &cobra.Command{
	Use:   "validate <file>",
	Short: "Validate a TASKS.json file",
	Long:  `Validate a TASKS.json file against the schema and check for errors.`,
	Args:  cobra.ExactArgs(1),
	RunE:  runValidate,
}

func runValidate(cmd *cobra.Command, args []string) error {
	path := args[0]

	tl, err := tasks.ParseFile(path)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	result := tasks.Validate(tl)

	if result.Valid {
		fmt.Fprintf(cmd.ErrOrStderr(), "✅ %s is valid\n", path)
		fmt.Fprintf(cmd.ErrOrStderr(), "   Project: %s\n", tl.Project)
		fmt.Fprintf(cmd.ErrOrStderr(), "   Tasks: %d\n", len(tl.Tasks))
		fmt.Fprintf(cmd.ErrOrStderr(), "   Areas: %d\n", len(tl.Areas))
		return nil
	}

	fmt.Fprintf(cmd.ErrOrStderr(), "❌ %s has %d error(s)\n\n", path, len(result.Errors))
	for _, e := range result.Errors {
		fmt.Fprintf(cmd.ErrOrStderr(), "  • %s: %s\n", e.Field, e.Message)
	}
	return fmt.Errorf("validation failed with %d error(s)", len(result.Errors))
}
