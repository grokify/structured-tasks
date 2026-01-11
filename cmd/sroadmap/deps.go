package main

import (
	"fmt"
	"io"
	"strings"

	"github.com/grokify/structured-roadmap/roadmap"
	"github.com/spf13/cobra"
)

type edge struct {
	from, to string
}

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

	// Collect items with dependencies
	var edges []edge
	itemMap := make(map[string]roadmap.Item)

	for _, item := range r.Items {
		itemMap[item.ID] = item
		for _, dep := range item.DependsOn {
			edges = append(edges, edge{from: dep, to: item.ID})
		}
	}

	out := cmd.OutOrStdout()

	if len(edges) == 0 {
		fmt.Fprintln(cmd.ErrOrStderr(), "No dependencies found in roadmap")
		return nil
	}

	switch depsFormat {
	case "mermaid":
		renderMermaid(out, r, edges, itemMap)
	case "dot":
		renderDot(out, r, edges, itemMap)
	default:
		return fmt.Errorf("unknown format: %s", depsFormat)
	}
	return nil
}

func renderMermaid(out io.Writer, _ *roadmap.Roadmap, edges []edge, itemMap map[string]roadmap.Item) {
	fmt.Fprintln(out, "```mermaid")
	fmt.Fprintln(out, "graph TD")

	// Define nodes with labels
	seen := make(map[string]bool)
	for _, e := range edges {
		if !seen[e.from] {
			item := itemMap[e.from]
			label := sanitizeMermaid(item.Title)
			status := statusShape(item.Status)
			fmt.Fprintf(out, "    %s%s%s\n", e.from, status[0], label)
			fmt.Fprintf(out, "    %s%s\n", e.from, status[1])
			seen[e.from] = true
		}
		if !seen[e.to] {
			item := itemMap[e.to]
			label := sanitizeMermaid(item.Title)
			status := statusShape(item.Status)
			fmt.Fprintf(out, "    %s%s%s\n", e.to, status[0], label)
			fmt.Fprintf(out, "    %s%s\n", e.to, status[1])
			seen[e.to] = true
		}
	}

	fmt.Fprintln(out)

	// Define edges
	for _, e := range edges {
		fmt.Fprintf(out, "    %s --> %s\n", e.from, e.to)
	}

	fmt.Fprintln(out, "```")
}

func renderDot(out io.Writer, r *roadmap.Roadmap, edges []edge, itemMap map[string]roadmap.Item) {
	fmt.Fprintf(out, "digraph \"%s\" {\n", r.Project)
	fmt.Fprintln(out, "    rankdir=LR;")
	fmt.Fprintln(out, "    node [shape=box];")
	fmt.Fprintln(out)

	// Define nodes
	seen := make(map[string]bool)
	for _, e := range edges {
		for _, id := range []string{e.from, e.to} {
			if !seen[id] {
				item := itemMap[id]
				color := statusColor(item.Status)
				fmt.Fprintf(out, "    %s [label=\"%s\" color=\"%s\"];\n", id, sanitizeDot(item.Title), color)
				seen[id] = true
			}
		}
	}

	fmt.Fprintln(out)

	// Define edges
	for _, e := range edges {
		fmt.Fprintf(out, "    %s -> %s;\n", e.from, e.to)
	}

	fmt.Fprintln(out, "}")
}

func sanitizeMermaid(s string) string {
	s = strings.ReplaceAll(s, "\"", "'")
	s = strings.ReplaceAll(s, "[", "(")
	s = strings.ReplaceAll(s, "]", ")")
	return "[\"" + s + "\"]"
}

func sanitizeDot(s string) string {
	s = strings.ReplaceAll(s, "\"", "\\\"")
	return s
}

func statusShape(status roadmap.Status) [2]string {
	switch status {
	case roadmap.StatusCompleted:
		return [2]string{"([", "])"}
	case roadmap.StatusInProgress:
		return [2]string{"{{", "}}"}
	case roadmap.StatusPlanned:
		return [2]string{"[", "]"}
	default:
		return [2]string{"((", "))"}
	}
}

func statusColor(status roadmap.Status) string {
	switch status {
	case roadmap.StatusCompleted:
		return "green"
	case roadmap.StatusInProgress:
		return "orange"
	case roadmap.StatusPlanned:
		return "blue"
	default:
		return "gray"
	}
}
