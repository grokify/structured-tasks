package main

import (
	"fmt"
	"os"

	"github.com/grokify/structured-roadmap/renderer"
	"github.com/grokify/structured-roadmap/roadmap"
	"github.com/spf13/cobra"
)

var (
	genInput           string
	genOutput          string
	genGroupBy         string
	genCheckbox        bool
	genEmoji           bool
	genLegend          bool
	genNoIntro         bool
	genTOC             bool
	genTOCDepth        int
	genOverview        bool
	genAreaSubheadings bool
	genNumbered        bool
	genNoRules         bool
)

var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate ROADMAP.md from ROADMAP.json",
	Long:  `Generate a Markdown roadmap file from a JSON intermediate representation.`,
	RunE:  runGenerate,
}

func init() {
	generateCmd.Flags().StringVarP(&genInput, "input", "i", "ROADMAP.json", "Input JSON file")
	generateCmd.Flags().StringVarP(&genOutput, "output", "o", "", "Output Markdown file (default: stdout)")
	generateCmd.Flags().StringVar(&genGroupBy, "group-by", "area", "Grouping: area, type, phase, status, quarter, priority")
	generateCmd.Flags().BoolVar(&genCheckbox, "checkboxes", true, "Use [x]/[ ] checkbox syntax")
	generateCmd.Flags().BoolVar(&genEmoji, "emoji", true, "Include emoji status indicators")
	generateCmd.Flags().BoolVar(&genLegend, "legend", false, "Show legend table")
	generateCmd.Flags().BoolVar(&genNoIntro, "no-intro", false, "Omit introductory paragraph")
	generateCmd.Flags().BoolVar(&genTOC, "toc", false, "Show table of contents")
	generateCmd.Flags().IntVar(&genTOCDepth, "toc-depth", 1, "TOC depth: 1 = sections only, 2 = sections + items")
	generateCmd.Flags().BoolVar(&genOverview, "overview", false, "Show overview table")
	generateCmd.Flags().BoolVar(&genAreaSubheadings, "area-subheadings", false, "Show area sub-sections within phases (use with --group-by phase)")
	generateCmd.Flags().BoolVar(&genNumbered, "numbered", false, "Number items")
	generateCmd.Flags().BoolVar(&genNoRules, "no-rules", false, "Omit horizontal rules between sections")
}

func runGenerate(cmd *cobra.Command, args []string) error {
	r, err := roadmap.ParseFile(genInput)
	if err != nil {
		return fmt.Errorf("failed to read input: %w", err)
	}

	// Validate first
	result := roadmap.Validate(r)
	if !result.Valid {
		fmt.Fprintf(cmd.ErrOrStderr(), "Validation errors in %s:\n", genInput)
		for _, e := range result.Errors {
			fmt.Fprintf(cmd.ErrOrStderr(), "  • %s: %s\n", e.Field, e.Message)
		}
		return fmt.Errorf("validation failed with %d error(s)", len(result.Errors))
	}

	// Build options
	opts := renderer.DefaultOptions()
	opts.UseCheckboxes = genCheckbox
	opts.UseEmoji = genEmoji
	opts.ShowLegend = genLegend
	opts.ShowIntro = !genNoIntro
	opts.ShowTOC = genTOC
	opts.TOCDepth = genTOCDepth
	opts.ShowOverviewTable = genOverview
	opts.ShowAreaSubheadings = genAreaSubheadings
	opts.NumberItems = genNumbered
	opts.HorizontalRules = !genNoRules

	switch genGroupBy {
	case "area":
		opts.GroupBy = renderer.GroupByArea
	case "type":
		opts.GroupBy = renderer.GroupByType
	case "phase":
		opts.GroupBy = renderer.GroupByPhase
	case "status":
		opts.GroupBy = renderer.GroupByStatus
	case "quarter":
		opts.GroupBy = renderer.GroupByQuarter
	case "priority":
		opts.GroupBy = renderer.GroupByPriority
	default:
		return fmt.Errorf("unknown group-by value: %s", genGroupBy)
	}

	// Render
	output := renderer.Render(r, opts)

	// Write output
	if genOutput == "" {
		fmt.Fprint(cmd.OutOrStdout(), output)
	} else {
		if err := os.WriteFile(genOutput, []byte(output), 0600); err != nil {
			return fmt.Errorf("failed to write output: %w", err)
		}
		fmt.Fprintf(cmd.ErrOrStderr(), "Generated %s\n", genOutput)
	}
	return nil
}
