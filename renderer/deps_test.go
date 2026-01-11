package renderer

import (
	"bytes"
	"strings"
	"testing"

	"github.com/grokify/structured-roadmap/roadmap"
)

func TestBuildDependencyGraph(t *testing.T) {
	r := &roadmap.Roadmap{
		Items: []roadmap.Item{
			{ID: "item-1", Title: "First Item", Status: roadmap.StatusCompleted},
			{ID: "item-2", Title: "Second Item", Status: roadmap.StatusPlanned, DependsOn: []string{"item-1"}},
			{ID: "item-3", Title: "Third Item", Status: roadmap.StatusFuture, DependsOn: []string{"item-1", "item-2"}},
		},
	}

	deps := BuildDependencyGraph(r)

	if len(deps.Edges) != 3 {
		t.Errorf("expected 3 edges, got %d", len(deps.Edges))
	}

	if len(deps.ItemMap) != 3 {
		t.Errorf("expected 3 items in map, got %d", len(deps.ItemMap))
	}

	// Verify edges
	expectedEdges := []Edge{
		{From: "item-1", To: "item-2"},
		{From: "item-1", To: "item-3"},
		{From: "item-2", To: "item-3"},
	}

	for i, expected := range expectedEdges {
		if deps.Edges[i] != expected {
			t.Errorf("edge %d: expected %v, got %v", i, expected, deps.Edges[i])
		}
	}
}

func TestBuildDependencyGraphEmpty(t *testing.T) {
	r := &roadmap.Roadmap{
		Items: []roadmap.Item{
			{ID: "item-1", Title: "First Item"},
			{ID: "item-2", Title: "Second Item"},
		},
	}

	deps := BuildDependencyGraph(r)

	if len(deps.Edges) != 0 {
		t.Errorf("expected 0 edges, got %d", len(deps.Edges))
	}
}

func TestRenderMermaid(t *testing.T) {
	r := &roadmap.Roadmap{
		Project: "test-project",
		Items: []roadmap.Item{
			{ID: "item1", Title: "First Item", Status: roadmap.StatusCompleted},
			{ID: "item2", Title: "Second Item", Status: roadmap.StatusPlanned, DependsOn: []string{"item1"}},
		},
	}

	deps := BuildDependencyGraph(r)
	var buf bytes.Buffer
	RenderMermaid(&buf, r, deps)

	output := buf.String()

	// Check basic structure
	if !strings.Contains(output, "```mermaid") {
		t.Error("expected mermaid code fence")
	}
	if !strings.Contains(output, "graph TD") {
		t.Error("expected graph TD directive")
	}
	if !strings.Contains(output, "item1 --> item2") {
		t.Error("expected edge from item1 to item2")
	}
	if !strings.Contains(output, "```") {
		t.Error("expected closing code fence")
	}
}

func TestRenderDOT(t *testing.T) {
	r := &roadmap.Roadmap{
		Project: "test-project",
		Items: []roadmap.Item{
			{ID: "item1", Title: "First Item", Status: roadmap.StatusCompleted},
			{ID: "item2", Title: "Second Item", Status: roadmap.StatusInProgress, DependsOn: []string{"item1"}},
		},
	}

	deps := BuildDependencyGraph(r)
	var buf bytes.Buffer
	RenderDOT(&buf, r, deps)

	output := buf.String()

	// Check basic structure
	if !strings.Contains(output, `digraph "test-project"`) {
		t.Error("expected digraph declaration")
	}
	if !strings.Contains(output, "rankdir=LR") {
		t.Error("expected rankdir directive")
	}
	if !strings.Contains(output, "item1 -> item2") {
		t.Error("expected edge from item1 to item2")
	}
	if !strings.Contains(output, `color="green"`) {
		t.Error("expected green color for completed item")
	}
	if !strings.Contains(output, `color="orange"`) {
		t.Error("expected orange color for in-progress item")
	}
}

func TestStatusShape(t *testing.T) {
	tests := []struct {
		status   roadmap.Status
		expected [2]string
	}{
		{roadmap.StatusCompleted, [2]string{"([", "])"}},
		{roadmap.StatusInProgress, [2]string{"{{", "}}"}},
		{roadmap.StatusPlanned, [2]string{"[", "]"}},
		{roadmap.StatusFuture, [2]string{"((", "))"}},
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
		status   roadmap.Status
		expected string
	}{
		{roadmap.StatusCompleted, "green"},
		{roadmap.StatusInProgress, "orange"},
		{roadmap.StatusPlanned, "blue"},
		{roadmap.StatusFuture, "gray"},
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
