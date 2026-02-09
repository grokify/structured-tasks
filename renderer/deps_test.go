package renderer

import (
	"bytes"
	"strings"
	"testing"

	"github.com/grokify/structured-tasks/tasks"
)

func TestBuildDependencyGraph(t *testing.T) {
	tl := &tasks.TaskList{
		Tasks: []tasks.Task{
			{ID: "task-1", Title: "First Task", Status: tasks.StatusCompleted},
			{ID: "task-2", Title: "Second Task", Status: tasks.StatusPlanned, DependsOn: []string{"task-1"}},
			{ID: "task-3", Title: "Third Task", Status: tasks.StatusFuture, DependsOn: []string{"task-1", "task-2"}},
		},
	}

	deps := BuildDependencyGraph(tl)

	if len(deps.Edges) != 3 {
		t.Errorf("expected 3 edges, got %d", len(deps.Edges))
	}

	if len(deps.TaskMap) != 3 {
		t.Errorf("expected 3 tasks in map, got %d", len(deps.TaskMap))
	}

	// Verify edges
	expectedEdges := []Edge{
		{From: "task-1", To: "task-2"},
		{From: "task-1", To: "task-3"},
		{From: "task-2", To: "task-3"},
	}

	for i, expected := range expectedEdges {
		if deps.Edges[i] != expected {
			t.Errorf("edge %d: expected %v, got %v", i, expected, deps.Edges[i])
		}
	}
}

func TestBuildDependencyGraphEmpty(t *testing.T) {
	tl := &tasks.TaskList{
		Tasks: []tasks.Task{
			{ID: "task-1", Title: "First Task"},
			{ID: "task-2", Title: "Second Task"},
		},
	}

	deps := BuildDependencyGraph(tl)

	if len(deps.Edges) != 0 {
		t.Errorf("expected 0 edges, got %d", len(deps.Edges))
	}
}

func TestRenderMermaid(t *testing.T) {
	tl := &tasks.TaskList{
		Project: "test-project",
		Tasks: []tasks.Task{
			{ID: "task1", Title: "First Task", Status: tasks.StatusCompleted},
			{ID: "task2", Title: "Second Task", Status: tasks.StatusPlanned, DependsOn: []string{"task1"}},
		},
	}

	deps := BuildDependencyGraph(tl)
	var buf bytes.Buffer
	RenderMermaid(&buf, tl, deps)

	output := buf.String()

	// Check basic structure
	if !strings.Contains(output, "```mermaid") {
		t.Error("expected mermaid code fence")
	}
	if !strings.Contains(output, "graph TD") {
		t.Error("expected graph TD directive")
	}
	if !strings.Contains(output, "task1 --> task2") {
		t.Error("expected edge from task1 to task2")
	}
	if !strings.Contains(output, "```") {
		t.Error("expected closing code fence")
	}
}

func TestRenderDOT(t *testing.T) {
	tl := &tasks.TaskList{
		Project: "test-project",
		Tasks: []tasks.Task{
			{ID: "task1", Title: "First Task", Status: tasks.StatusCompleted},
			{ID: "task2", Title: "Second Task", Status: tasks.StatusInProgress, DependsOn: []string{"task1"}},
		},
	}

	deps := BuildDependencyGraph(tl)
	var buf bytes.Buffer
	RenderDOT(&buf, tl, deps)

	output := buf.String()

	// Check basic structure
	if !strings.Contains(output, `digraph "test-project"`) {
		t.Error("expected digraph declaration")
	}
	if !strings.Contains(output, "rankdir=LR") {
		t.Error("expected rankdir directive")
	}
	if !strings.Contains(output, "task1 -> task2") {
		t.Error("expected edge from task1 to task2")
	}
	if !strings.Contains(output, `color="green"`) {
		t.Error("expected green color for completed task")
	}
	if !strings.Contains(output, `color="orange"`) {
		t.Error("expected orange color for in-progress task")
	}
}

func TestStatusShape(t *testing.T) {
	tests := []struct {
		status   tasks.Status
		expected [2]string
	}{
		{tasks.StatusCompleted, [2]string{"([", "])"}},
		{tasks.StatusInProgress, [2]string{"{{", "}}"}},
		{tasks.StatusPlanned, [2]string{"[", "]"}},
		{tasks.StatusFuture, [2]string{"((", "))"}},
		{"unknown", [2]string{"((", "))"}},
	}

	for _, tt := range tests {
		result := StatusShape(tt.status)
		if result != tt.expected {
			t.Errorf("StatusShape(%s): expected %v, got %v", tt.status, tt.expected, result)
		}
	}
}

func TestStatusColor(t *testing.T) {
	tests := []struct {
		status   tasks.Status
		expected string
	}{
		{tasks.StatusCompleted, "green"},
		{tasks.StatusInProgress, "orange"},
		{tasks.StatusPlanned, "blue"},
		{tasks.StatusFuture, "gray"},
		{"unknown", "gray"},
	}

	for _, tt := range tests {
		result := StatusColor(tt.status)
		if result != tt.expected {
			t.Errorf("StatusColor(%s): expected %s, got %s", tt.status, tt.expected, result)
		}
	}
}

func TestSanitizeMermaid(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"Simple text", `["Simple text"]`},
		{`Text with "quotes"`, `["Text with 'quotes'"]`},
		{"Text with [brackets]", `["Text with (brackets)"]`},
	}

	for _, tt := range tests {
		result := sanitizeMermaid(tt.input)
		if result != tt.expected {
			t.Errorf("sanitizeMermaid(%q): expected %q, got %q", tt.input, tt.expected, result)
		}
	}
}

func TestSanitizeDOT(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"Simple text", "Simple text"},
		{`Text with "quotes"`, `Text with \"quotes\"`},
	}

	for _, tt := range tests {
		result := sanitizeDOT(tt.input)
		if result != tt.expected {
			t.Errorf("sanitizeDOT(%q): expected %q, got %q", tt.input, tt.expected, result)
		}
	}
}
