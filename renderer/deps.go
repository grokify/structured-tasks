package renderer

import (
	"fmt"
	"io"
	"strings"

	"github.com/grokify/structured-roadmap/roadmap"
)

// Edge represents a dependency relationship between two items.
type Edge struct {
	From, To string
}

// DepsResult contains the dependency graph data.
type DepsResult struct {
	Edges   []Edge
	ItemMap map[string]roadmap.Item
}

// BuildDependencyGraph extracts dependency edges from a roadmap.
func BuildDependencyGraph(r *roadmap.Roadmap) DepsResult {
	var edges []Edge
	itemMap := make(map[string]roadmap.Item)

	for _, item := range r.Items {
		itemMap[item.ID] = item
		for _, dep := range item.DependsOn {
			edges = append(edges, Edge{From: dep, To: item.ID})
		}
	}

	return DepsResult{
		Edges:   edges,
		ItemMap: itemMap,
	}
}

// RenderMermaid renders a dependency graph in Mermaid format.
func RenderMermaid(w io.Writer, r *roadmap.Roadmap, deps DepsResult) {
	fmt.Fprintln(w, "```mermaid")
	fmt.Fprintln(w, "graph TD")

	// Define nodes with labels
	seen := make(map[string]bool)
	for _, e := range deps.Edges {
		if !seen[e.From] {
			item := deps.ItemMap[e.From]
			label := sanitizeMermaid(item.Title)
			shape := StatusShape(item.Status)
			fmt.Fprintf(w, "    %s%s%s\n", e.From, shape[0], label)
			fmt.Fprintf(w, "    %s%s\n", e.From, shape[1])
			seen[e.From] = true
		}
		if !seen[e.To] {
			item := deps.ItemMap[e.To]
			label := sanitizeMermaid(item.Title)
			shape := StatusShape(item.Status)
			fmt.Fprintf(w, "    %s%s%s\n", e.To, shape[0], label)
			fmt.Fprintf(w, "    %s%s\n", e.To, shape[1])
			seen[e.To] = true
		}
	}

	fmt.Fprintln(w)

	// Define edges
	for _, e := range deps.Edges {
		fmt.Fprintf(w, "    %s --> %s\n", e.From, e.To)
	}

	fmt.Fprintln(w, "```")
}

// RenderDOT renders a dependency graph in Graphviz DOT format.
func RenderDOT(w io.Writer, r *roadmap.Roadmap, deps DepsResult) {
	fmt.Fprintf(w, "digraph \"%s\" {\n", r.Project)
	fmt.Fprintln(w, "    rankdir=LR;")
	fmt.Fprintln(w, "    node [shape=box];")
	fmt.Fprintln(w)

	// Define nodes
	seen := make(map[string]bool)
	for _, e := range deps.Edges {
		for _, id := range []string{e.From, e.To} {
			if !seen[id] {
				item := deps.ItemMap[id]
				color := StatusColor(item.Status)
				fmt.Fprintf(w, "    %s [label=\"%s\" color=\"%s\"];\n", id, sanitizeDOT(item.Title), color)
				seen[id] = true
			}
		}
	}

	fmt.Fprintln(w)

	// Define edges
	for _, e := range deps.Edges {
		fmt.Fprintf(w, "    %s -> %s;\n", e.From, e.To)
	}

	fmt.Fprintln(w, "}")
}

// StatusShape returns the Mermaid node shape for a status.
// Returns [opening, closing] brackets.
func StatusShape(status roadmap.Status) [2]string {
	switch status {
	case roadmap.StatusCompleted:
		return [2]string{"([", "])"} // Stadium/rounded
	case roadmap.StatusInProgress:
		return [2]string{"{{", "}}"} // Hexagon
	case roadmap.StatusPlanned:
		return [2]string{"[", "]"} // Rectangle
	default:
		return [2]string{"((", "))"} // Circle
	}
}

// StatusColor returns the DOT node color for a status.
func StatusColor(status roadmap.Status) string {
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

// sanitizeMermaid escapes special characters for Mermaid labels.
func sanitizeMermaid(s string) string {
	s = strings.ReplaceAll(s, "\"", "'")
	s = strings.ReplaceAll(s, "[", "(")
	s = strings.ReplaceAll(s, "]", ")")
	return "[\"" + s + "\"]"
}

// sanitizeDOT escapes special characters for DOT labels.
func sanitizeDOT(s string) string {
	s = strings.ReplaceAll(s, "\"", "\\\"")
	return s
}
