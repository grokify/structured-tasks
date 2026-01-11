package renderer

import (
	"strings"
	"testing"

	"github.com/grokify/structured-roadmap/roadmap"
)

func TestRender(t *testing.T) {
	r := &roadmap.Roadmap{
		IRVersion: "1.0",
		Project:   "Test Project",
		Areas: []roadmap.Area{
			{ID: "core", Name: "Core Features", Priority: 1},
			{ID: "api", Name: "API", Priority: 2},
		},
		Items: []roadmap.Item{
			{ID: "item-1", Title: "Feature 1", Status: roadmap.StatusCompleted, Area: "core", Type: "Added"},
			{ID: "item-2", Title: "Feature 2", Status: roadmap.StatusPlanned, Area: "core", Type: "Added"},
			{ID: "item-3", Title: "Feature 3", Status: roadmap.StatusInProgress, Area: "api", Type: "Changed"},
		},
	}

	t.Run("default options", func(t *testing.T) {
		opts := DefaultOptions()
		output := Render(r, opts)

		if !strings.Contains(output, "# Test Project Roadmap") {
			t.Error("Expected title in output")
		}
		if !strings.Contains(output, "## Core Features") {
			t.Error("Expected area heading in output")
		}
		if !strings.Contains(output, "[x]") {
			t.Error("Expected completed checkbox in output")
		}
		if !strings.Contains(output, "[ ]") {
			t.Error("Expected uncompleted checkbox in output")
		}
	})

	t.Run("no checkboxes", func(t *testing.T) {
		opts := DefaultOptions().WithCheckboxes(false)
		output := Render(r, opts)

		if strings.Contains(output, "[x]") || strings.Contains(output, "[ ]") {
			t.Error("Should not contain checkboxes")
		}
	})

	t.Run("no emoji", func(t *testing.T) {
		opts := DefaultOptions().WithEmoji(false)
		output := Render(r, opts)

		if strings.Contains(output, "✅") || strings.Contains(output, "🚧") || strings.Contains(output, "📋") {
			t.Error("Should not contain emoji")
		}
	})

	t.Run("with legend", func(t *testing.T) {
		opts := DefaultOptions().WithLegend(true)
		output := Render(r, opts)

		if !strings.Contains(output, "## Legend") {
			t.Error("Expected legend section in output")
		}
	})

	t.Run("with numbered items", func(t *testing.T) {
		opts := DefaultOptions().WithNumberedItems(true)
		output := Render(r, opts)

		if !strings.Contains(output, "1.") {
			t.Error("Expected numbered items in output")
		}
	})
}

func TestRenderGroupBy(t *testing.T) {
	r := &roadmap.Roadmap{
		IRVersion: "1.0",
		Project:   "Test",
		Areas: []roadmap.Area{
			{ID: "core", Name: "Core", Priority: 1},
		},
		Phases: []roadmap.Phase{
			{ID: "phase-1", Name: "Phase 1", Status: roadmap.StatusCompleted, Order: 1},
		},
		Items: []roadmap.Item{
			{ID: "1", Title: "Item 1", Status: roadmap.StatusCompleted, Area: "core", Phase: "phase-1", Priority: roadmap.PriorityHigh, TargetQuarter: "Q1 2026", Type: "Added"},
			{ID: "2", Title: "Item 2", Status: roadmap.StatusPlanned, Area: "core", Phase: "phase-1", Priority: roadmap.PriorityMedium, TargetQuarter: "Q2 2026", Type: "Changed"},
		},
	}

	tests := []struct {
		name    string
		groupBy GroupBy
		expect  string
	}{
		{"by area", GroupByArea, "## Core"},
		{"by type", GroupByType, "## Added"},
		{"by phase", GroupByPhase, "## Phase 1"},
		{"by status", GroupByStatus, "## ✅ Completed"},
		{"by quarter", GroupByQuarter, "## Q1 2026"},
		{"by priority", GroupByPriority, "## High Priority"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opts := DefaultOptions().WithGroupBy(tt.groupBy)
			output := Render(r, opts)

			if !strings.Contains(output, tt.expect) {
				t.Errorf("Expected %q in output for group by %s", tt.expect, tt.groupBy)
			}
		})
	}
}

func TestRenderPhaseWithAreaSubheadings(t *testing.T) {
	r := &roadmap.Roadmap{
		IRVersion: "1.0",
		Project:   "Test",
		Areas: []roadmap.Area{
			{ID: "core", Name: "Core Package", Priority: 1},
			{ID: "api", Name: "API Layer", Priority: 2},
		},
		Phases: []roadmap.Phase{
			{ID: "phase-1", Name: "Phase 1: Foundation", Status: roadmap.StatusCompleted, Order: 1},
		},
		Items: []roadmap.Item{
			{ID: "1", Title: "interfaces.go", Status: roadmap.StatusCompleted, Area: "core", Phase: "phase-1"},
			{ID: "2", Title: "client.go", Status: roadmap.StatusCompleted, Area: "api", Phase: "phase-1"},
		},
	}

	opts := DefaultOptions().WithGroupBy(GroupByPhase)
	opts.ShowAreaSubheadings = true
	output := Render(r, opts)

	if !strings.Contains(output, "## Phase 1: Foundation") {
		t.Error("Expected phase heading")
	}
	if !strings.Contains(output, "### Core Package") {
		t.Error("Expected area subheading")
	}
	if !strings.Contains(output, "### API Layer") {
		t.Error("Expected area subheading")
	}
}

func TestRenderTOC(t *testing.T) {
	r := &roadmap.Roadmap{
		IRVersion: "1.0",
		Project:   "Test",
		Areas: []roadmap.Area{
			{ID: "core", Name: "Core", Priority: 1},
		},
		Items: []roadmap.Item{
			{ID: "1", Title: "Item 1", Status: roadmap.StatusCompleted, Area: "core"},
			{ID: "2", Title: "Item 2", Status: roadmap.StatusPlanned, Area: "core"},
		},
	}

	opts := DefaultOptions()
	opts.ShowTOC = true
	output := Render(r, opts)

	if !strings.Contains(output, "## Table of Contents") {
		t.Error("Expected TOC section")
	}
	if !strings.Contains(output, "(1/2)") {
		t.Error("Expected progress count in TOC")
	}
}

func TestRenderOverviewTable(t *testing.T) {
	r := &roadmap.Roadmap{
		IRVersion: "1.0",
		Project:   "Test",
		Areas: []roadmap.Area{
			{ID: "core", Name: "Core", Priority: 1},
		},
		Items: []roadmap.Item{
			{ID: "1", Title: "Item 1", Status: roadmap.StatusCompleted, Area: "core", Priority: roadmap.PriorityHigh},
			{ID: "2", Title: "Item 2", Status: roadmap.StatusPlanned, Area: "core", Priority: roadmap.PriorityMedium},
		},
	}

	opts := DefaultOptions()
	opts.ShowOverviewTable = true
	output := Render(r, opts)

	if !strings.Contains(output, "## Overview") {
		t.Error("Expected overview section")
	}
	if !strings.Contains(output, "| Item | Status | Priority | Area |") {
		t.Error("Expected table headers")
	}
}

func TestRenderContentBlocks(t *testing.T) {
	r := &roadmap.Roadmap{
		IRVersion: "1.0",
		Project:   "Test",
		Items: []roadmap.Item{
			{
				ID:     "1",
				Title:  "Item with content",
				Status: roadmap.StatusCompleted,
				Content: []roadmap.ContentBlock{
					{Type: roadmap.ContentTypeText, Value: "Some text content"},
					{Type: roadmap.ContentTypeCode, Value: "fmt.Println(\"hello\")", Language: "go"},
					{Type: roadmap.ContentTypeBlockquote, Value: "A blockquote"},
					{Type: roadmap.ContentTypeList, Items: []string{"Item A", "Item B"}},
					{Type: roadmap.ContentTypeTable, Headers: []string{"Col1", "Col2"}, Rows: [][]string{{"A", "B"}}},
				},
			},
		},
	}

	opts := DefaultOptions()
	output := Render(r, opts)

	if !strings.Contains(output, "Some text content") {
		t.Error("Expected text content")
	}
	if !strings.Contains(output, "```go") {
		t.Error("Expected code block with language")
	}
	if !strings.Contains(output, "> A blockquote") {
		t.Error("Expected blockquote")
	}
	if !strings.Contains(output, "- Item A") {
		t.Error("Expected list items")
	}
	if !strings.Contains(output, "| Col1 | Col2 |") {
		t.Error("Expected table")
	}
}

func TestRenderSections(t *testing.T) {
	r := &roadmap.Roadmap{
		IRVersion: "1.0",
		Project:   "Test",
		Sections: []roadmap.Section{
			{
				ID:    "overview",
				Title: "Overview",
				Order: 1,
				Content: []roadmap.ContentBlock{
					{Type: roadmap.ContentTypeText, Value: "Project overview text."},
				},
			},
		},
	}

	opts := DefaultOptions()
	output := Render(r, opts)

	if !strings.Contains(output, "## Overview") {
		t.Error("Expected section heading")
	}
	if !strings.Contains(output, "Project overview text.") {
		t.Error("Expected section content")
	}
}

func TestRenderVersionHistory(t *testing.T) {
	r := &roadmap.Roadmap{
		IRVersion: "1.0",
		Project:   "Test",
		VersionHistory: []roadmap.VersionEntry{
			{Version: "1.0.0", Date: "2026-01-01", Status: roadmap.StatusCompleted, Summary: "Initial release"},
			{Version: "2.0.0", Status: roadmap.StatusPlanned, Summary: "Major update"},
		},
	}

	opts := DefaultOptions()
	output := Render(r, opts)

	if !strings.Contains(output, "## Version History") {
		t.Error("Expected version history section")
	}
	if !strings.Contains(output, "1.0.0") {
		t.Error("Expected version in output")
	}
}

func TestRenderTasks(t *testing.T) {
	r := &roadmap.Roadmap{
		IRVersion: "1.0",
		Project:   "Test",
		Items: []roadmap.Item{
			{
				ID:     "1",
				Title:  "Feature with tasks",
				Status: roadmap.StatusInProgress,
				Tasks: []roadmap.Task{
					{Description: "Completed task", Completed: true},
					{Description: "Pending task", Completed: false},
				},
			},
		},
	}

	opts := DefaultOptions()
	output := Render(r, opts)

	if !strings.Contains(output, "[x] Completed task") {
		t.Error("Expected completed task")
	}
	if !strings.Contains(output, "[ ] Pending task") {
		t.Error("Expected pending task")
	}
}

func TestDefaultOptions(t *testing.T) {
	opts := DefaultOptions()

	if opts.GroupBy != GroupByArea {
		t.Errorf("Default GroupBy = %s, want %s", opts.GroupBy, GroupByArea)
	}
	if !opts.ShowCompleted {
		t.Error("ShowCompleted should default to true")
	}
	if !opts.UseCheckboxes {
		t.Error("UseCheckboxes should default to true")
	}
	if !opts.UseEmoji {
		t.Error("UseEmoji should default to true")
	}
	if !opts.ShowIntro {
		t.Error("ShowIntro should default to true")
	}
	if opts.ShowLegend {
		t.Error("ShowLegend should default to false")
	}
	if opts.ShowTOC {
		t.Error("ShowTOC should default to false")
	}
	if !opts.HorizontalRules {
		t.Error("HorizontalRules should default to true")
	}
}
