package main

import (
	"fmt"

	"github.com/grokify/structured-roadmap/renderer"
	"github.com/grokify/structured-roadmap/roadmap"
	"github.com/spf13/cobra"
)

var depsFormat string

var depsCmd = &cobra.Command{
	Use:   "deps <file>",
	Short: "Generate dependency graph",
	Long:  `Generate a dependency graph from item dependencies in Mermaid or DOT format.`,
	Args:  cobra.ExactArgs(1),
	RunE:  runDeps,
}

func init() {
	depsCmd.Flags().StringVar(&depsFormat, "format", "mermaid", "Output format: mermaid, dot")
}

func runDeps(cmd *cobra.Command, args []string) error {
	path := args[0]

	r, err := roadmap.ParseFile(path)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	deps := renderer.BuildDependencyGraph(r)
	out := cmd.OutOrStdout()

	if len(deps.Edges) == 0 {
		fmt.Fprintln(cmd.ErrOrStderr(), "No dependencies found in roadmap")
		return nil
	}

	switch depsFormat {
	case "mermaid":
		renderer.RenderMermaid(out, r, deps)
	case "dot":
		renderer.RenderDOT(out, r, deps)
	default:
		return fmt.Errorf("unknown format: %s", depsFormat)
	}
	return nil
}
