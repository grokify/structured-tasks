package renderer

import (
	"fmt"
	"io"
	"strings"

	"github.com/grokify/structured-tasks/tasks"
)

// Edge represents a dependency relationship between two tasks.
type Edge struct {
	From, To string
}

// DepsResult contains the dependency graph data.
type DepsResult struct {
	Edges   []Edge
	TaskMap map[string]tasks.Task
}

// BuildDependencyGraph extracts dependency edges from a task list.
func BuildDependencyGraph(tl *tasks.TaskList) DepsResult {
	var edges []Edge
	taskMap := make(map[string]tasks.Task)

	for _, task := range tl.Tasks {
		taskMap[task.ID] = task
		for _, dep := range task.DependsOn {
			edges = append(edges, Edge{From: dep, To: task.ID})
		}
	}

	return DepsResult{
		Edges:   edges,
		TaskMap: taskMap,
	}
}

// RenderMermaid renders a dependency graph in Mermaid format.
func RenderMermaid(w io.Writer, tl *tasks.TaskList, deps DepsResult) {
	fmt.Fprintln(w, "```mermaid")
	fmt.Fprintln(w, "graph TD")

	// Define nodes with labels
	seen := make(map[string]bool)
	for _, e := range deps.Edges {
		if !seen[e.From] {
			task := deps.TaskMap[e.From]
			label := sanitizeMermaid(task.Title)
			shape := StatusShape(task.Status)
			fmt.Fprintf(w, "    %s%s%s\n", e.From, shape[0], label)
			fmt.Fprintf(w, "    %s%s\n", e.From, shape[1])
			seen[e.From] = true
		}
		if !seen[e.To] {
			task := deps.TaskMap[e.To]
			label := sanitizeMermaid(task.Title)
			shape := StatusShape(task.Status)
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
func RenderDOT(w io.Writer, tl *tasks.TaskList, deps DepsResult) {
	fmt.Fprintf(w, "digraph \"%s\" {\n", tl.Project)
	fmt.Fprintln(w, "    rankdir=LR;")
	fmt.Fprintln(w, "    node [shape=box];")
	fmt.Fprintln(w)

	// Define nodes
	seen := make(map[string]bool)
	for _, e := range deps.Edges {
		for _, id := range []string{e.From, e.To} {
			if !seen[id] {
				task := deps.TaskMap[id]
				color := StatusColor(task.Status)
				fmt.Fprintf(w, "    %s [label=\"%s\" color=\"%s\"];\n", id, sanitizeDOT(task.Title), color)
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
func StatusShape(status tasks.Status) [2]string {
	switch status {
	case tasks.StatusCompleted:
		return [2]string{"([", "])"} // Stadium/rounded
	case tasks.StatusInProgress:
		return [2]string{"{{", "}}"} // Hexagon
	case tasks.StatusPlanned:
		return [2]string{"[", "]"} // Rectangle
	default:
		return [2]string{"((", "))"} // Circle
	}
}

// StatusColor returns the DOT node color for a status.
func StatusColor(status tasks.Status) string {
	switch status {
	case tasks.StatusCompleted:
		return "green"
	case tasks.StatusInProgress:
		return "orange"
	case tasks.StatusPlanned:
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
