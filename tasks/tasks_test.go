package tasks

import (
	"errors"
	"strings"
	"testing"
)

func TestParse(t *testing.T) {
	tests := []struct {
		name    string
		json    string
		wantErr bool
	}{
		{
			name: "minimal valid task list",
			json: `{"ir_version": "1.0", "project": "test-project"}`,
		},
		{
			name: "task list with areas",
			json: `{
				"ir_version": "1.0",
				"project": "test-project",
				"areas": [
					{"id": "core", "name": "Core"}
				]
			}`,
		},
		{
			name: "task list with tasks",
			json: `{
				"ir_version": "1.0",
				"project": "test-project",
				"tasks": [
					{"id": "task-1", "title": "Feature 1", "status": "completed"},
					{"id": "task-2", "title": "Feature 2", "status": "planned"}
				]
			}`,
		},
		{
			name:    "invalid json",
			json:    `{invalid}`,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tl, err := Parse([]byte(tt.json))
			if (err != nil) != tt.wantErr {
				t.Errorf("Parse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && tl == nil {
				t.Error("Parse() returned nil task list for valid input")
			}
		})
	}
}

func TestValidate(t *testing.T) {
	tests := []struct {
		name      string
		taskList  *TaskList
		wantValid bool
	}{
		{
			name: "valid minimal task list",
			taskList: &TaskList{
				IRVersion: "1.0",
				Project:   "test",
			},
			wantValid: true,
		},
		{
			name: "missing ir_version",
			taskList: &TaskList{
				Project: "test",
			},
			wantValid: false,
		},
		{
			name: "missing project",
			taskList: &TaskList{
				IRVersion: "1.0",
			},
			wantValid: false,
		},
		{
			name: "unsupported ir_version",
			taskList: &TaskList{
				IRVersion: "2.0",
				Project:   "test",
			},
			wantValid: false,
		},
		{
			name: "valid task list with tasks",
			taskList: &TaskList{
				IRVersion: "1.0",
				Project:   "test",
				Tasks: []Task{
					{ID: "task-1", Title: "Feature", Status: StatusCompleted},
				},
			},
			wantValid: true,
		},
		{
			name: "task missing id",
			taskList: &TaskList{
				IRVersion: "1.0",
				Project:   "test",
				Tasks: []Task{
					{Title: "Feature", Status: StatusCompleted},
				},
			},
			wantValid: false,
		},
		{
			name: "task missing title",
			taskList: &TaskList{
				IRVersion: "1.0",
				Project:   "test",
				Tasks: []Task{
					{ID: "task-1", Status: StatusCompleted},
				},
			},
			wantValid: false,
		},
		{
			name: "task missing status",
			taskList: &TaskList{
				IRVersion: "1.0",
				Project:   "test",
				Tasks: []Task{
					{ID: "task-1", Title: "Feature"},
				},
			},
			wantValid: false,
		},
		{
			name: "duplicate task ids",
			taskList: &TaskList{
				IRVersion: "1.0",
				Project:   "test",
				Tasks: []Task{
					{ID: "task-1", Title: "Feature 1", Status: StatusCompleted},
					{ID: "task-1", Title: "Feature 2", Status: StatusPlanned},
				},
			},
			wantValid: false,
		},
		{
			name: "invalid status",
			taskList: &TaskList{
				IRVersion: "1.0",
				Project:   "test",
				Tasks: []Task{
					{ID: "task-1", Title: "Feature", Status: "invalid"},
				},
			},
			wantValid: false,
		},
		{
			name: "valid type from structured-changelog",
			taskList: &TaskList{
				IRVersion: "1.0",
				Project:   "test",
				Tasks: []Task{
					{ID: "task-1", Title: "Feature", Status: StatusCompleted, Type: "Added"},
				},
			},
			wantValid: true,
		},
		{
			name: "invalid type",
			taskList: &TaskList{
				IRVersion: "1.0",
				Project:   "test",
				Tasks: []Task{
					{ID: "task-1", Title: "Feature", Status: StatusCompleted, Type: "InvalidType"},
				},
			},
			wantValid: false,
		},
		{
			name: "valid dependency reference",
			taskList: &TaskList{
				IRVersion: "1.0",
				Project:   "test",
				Tasks: []Task{
					{ID: "task-1", Title: "Feature 1", Status: StatusCompleted},
					{ID: "task-2", Title: "Feature 2", Status: StatusPlanned, DependsOn: []string{"task-1"}},
				},
			},
			wantValid: true,
		},
		{
			name: "invalid dependency reference",
			taskList: &TaskList{
				IRVersion: "1.0",
				Project:   "test",
				Tasks: []Task{
					{ID: "task-1", Title: "Feature 1", Status: StatusCompleted},
					{ID: "task-2", Title: "Feature 2", Status: StatusPlanned, DependsOn: []string{"nonexistent"}},
				},
			},
			wantValid: false,
		},
		{
			name: "valid area reference",
			taskList: &TaskList{
				IRVersion: "1.0",
				Project:   "test",
				Areas:     []Area{{ID: "core", Name: "Core"}},
				Tasks: []Task{
					{ID: "task-1", Title: "Feature", Status: StatusCompleted, Area: "core"},
				},
			},
			wantValid: true,
		},
		{
			name: "invalid area reference",
			taskList: &TaskList{
				IRVersion: "1.0",
				Project:   "test",
				Areas:     []Area{{ID: "core", Name: "Core"}},
				Tasks: []Task{
					{ID: "task-1", Title: "Feature", Status: StatusCompleted, Area: "nonexistent"},
				},
			},
			wantValid: false,
		},
		{
			name: "valid phase number",
			taskList: &TaskList{
				IRVersion: "1.0",
				Project:   "test",
				Tasks: []Task{
					{ID: "task-1", Title: "Feature", Status: StatusCompleted, Phase: 1},
				},
			},
			wantValid: true,
		},
		{
			name: "task with subtasks",
			taskList: &TaskList{
				IRVersion: "1.0",
				Project:   "test",
				Tasks: []Task{
					{
						ID:     "task-1",
						Title:  "Task with subtasks",
						Status: StatusInProgress,
						Subtasks: []Subtask{
							{Description: "Subtask 1", Completed: true},
							{Description: "Subtask 2", Completed: false},
						},
					},
				},
			},
			wantValid: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Validate(tt.taskList)
			if result.Valid != tt.wantValid {
				t.Errorf("Validate() valid = %v, want %v", result.Valid, tt.wantValid)
				if len(result.Errors) > 0 {
					for _, e := range result.Errors {
						t.Logf("  Error: %s: %s", e.Field, e.Message)
					}
				}
			}
		})
	}
}

func TestStats(t *testing.T) {
	tl := &TaskList{
		IRVersion: "1.0",
		Project:   "test",
		Tasks: []Task{
			{ID: "1", Title: "Task 1", Status: StatusCompleted, Area: "core", Type: "Added"},
			{ID: "2", Title: "Task 2", Status: StatusCompleted, Area: "core", Type: "Added"},
			{ID: "3", Title: "Task 3", Status: StatusInProgress, Area: "api", Type: "Changed"},
			{ID: "4", Title: "Task 4", Status: StatusPlanned, Area: "api", Type: "Added"},
			{ID: "5", Title: "Task 5", Status: StatusFuture},
		},
	}

	stats := tl.Stats()

	if stats.Total != 5 {
		t.Errorf("Total = %d, want 5", stats.Total)
	}
	if stats.ByStatus[StatusCompleted] != 2 {
		t.Errorf("ByStatus[completed] = %d, want 2", stats.ByStatus[StatusCompleted])
	}
	if stats.ByStatus[StatusInProgress] != 1 {
		t.Errorf("ByStatus[in_progress] = %d, want 1", stats.ByStatus[StatusInProgress])
	}
	if stats.ByArea["core"] != 2 {
		t.Errorf("ByArea[core] = %d, want 2", stats.ByArea["core"])
	}
	if stats.ByType["Added"] != 3 {
		t.Errorf("ByType[Added] = %d, want 3", stats.ByType["Added"])
	}
	if stats.CompletedCount() != 2 {
		t.Errorf("CompletedCount() = %d, want 2", stats.CompletedCount())
	}
}

func TestTasksBy(t *testing.T) {
	tl := &TaskList{
		IRVersion: "1.0",
		Project:   "test",
		Tasks: []Task{
			{ID: "1", Title: "Task 1", Status: StatusCompleted, Area: "core", Phase: 1},
			{ID: "2", Title: "Task 2", Status: StatusPlanned, Area: "api", Phase: 1},
			{ID: "3", Title: "Task 3", Status: StatusPlanned, Area: "core", Phase: 2},
		},
	}

	t.Run("TasksByArea", func(t *testing.T) {
		byArea := tl.TasksByArea()
		if len(byArea["core"]) != 2 {
			t.Errorf("TasksByArea[core] = %d tasks, want 2", len(byArea["core"]))
		}
		if len(byArea["api"]) != 1 {
			t.Errorf("TasksByArea[api] = %d tasks, want 1", len(byArea["api"]))
		}
	})

	t.Run("TasksByPhase", func(t *testing.T) {
		byPhase := tl.TasksByPhase()
		if len(byPhase[1]) != 2 {
			t.Errorf("TasksByPhase[1] = %d tasks, want 2", len(byPhase[1]))
		}
		if len(byPhase[2]) != 1 {
			t.Errorf("TasksByPhase[2] = %d tasks, want 1", len(byPhase[2]))
		}
	})

	t.Run("TasksByStatus", func(t *testing.T) {
		byStatus := tl.TasksByStatus()
		if len(byStatus[StatusCompleted]) != 1 {
			t.Errorf("TasksByStatus[completed] = %d tasks, want 1", len(byStatus[StatusCompleted]))
		}
		if len(byStatus[StatusPlanned]) != 2 {
			t.Errorf("TasksByStatus[planned] = %d tasks, want 2", len(byStatus[StatusPlanned]))
		}
	})
}

func TestDefaultLegend(t *testing.T) {
	legend := DefaultLegend()
	if legend[StatusCompleted].Emoji != "✅" {
		t.Errorf("DefaultLegend[completed].Emoji = %q, want ✅", legend[StatusCompleted].Emoji)
	}
	if legend[StatusInProgress].Emoji != "🚧" {
		t.Errorf("DefaultLegend[in_progress].Emoji = %q, want 🚧", legend[StatusInProgress].Emoji)
	}
}

func TestGetLegend(t *testing.T) {
	t.Run("uses default when no legend", func(t *testing.T) {
		tl := &TaskList{IRVersion: "1.0", Project: "test"}
		legend := tl.GetLegend()
		if legend[StatusCompleted].Emoji != "✅" {
			t.Error("Expected default legend")
		}
	})

	t.Run("merges custom legend", func(t *testing.T) {
		tl := &TaskList{
			IRVersion: "1.0",
			Project:   "test",
			Legend: map[Status]LegendEntry{
				StatusCompleted: {Emoji: "✓", Description: "Done"},
			},
		}
		legend := tl.GetLegend()
		if legend[StatusCompleted].Emoji != "✓" {
			t.Errorf("Expected custom emoji, got %q", legend[StatusCompleted].Emoji)
		}
		// Should still have defaults for other statuses
		if legend[StatusInProgress].Emoji != "🚧" {
			t.Errorf("Expected default emoji for in_progress, got %q", legend[StatusInProgress].Emoji)
		}
	})
}

func TestSentinelErrors(t *testing.T) {
	t.Run("Parse returns ErrParseJSON for invalid JSON", func(t *testing.T) {
		_, err := Parse([]byte(`{invalid}`))
		if err == nil {
			t.Fatal("Expected error, got nil")
		}
		if !errors.Is(err, ErrParseJSON) {
			t.Errorf("Expected error to wrap ErrParseJSON, got %v", err)
		}
	})

	t.Run("ParseFile returns ErrReadFile for missing file", func(t *testing.T) {
		_, err := ParseFile("/nonexistent/file.json")
		if err == nil {
			t.Fatal("Expected error, got nil")
		}
		if !errors.Is(err, ErrReadFile) {
			t.Errorf("Expected error to wrap ErrReadFile, got %v", err)
		}
	})
}

func TestFieldError(t *testing.T) {
	err := NewFieldError("tasks[0].id", "required field is missing", ErrMissingRequiredField)

	if err.Field != "tasks[0].id" {
		t.Errorf("Field = %q, want %q", err.Field, "tasks[0].id")
	}
	if err.Message != "required field is missing" {
		t.Errorf("Message = %q, want %q", err.Message, "required field is missing")
	}
	if !errors.Is(err, ErrMissingRequiredField) {
		t.Error("Expected error to wrap ErrMissingRequiredField")
	}

	expectedStr := "tasks[0].id: required field is missing"
	if err.Error() != expectedStr {
		t.Errorf("Error() = %q, want %q", err.Error(), expectedStr)
	}
}

func TestGetStatusEmoji(t *testing.T) {
	tl := &TaskList{IRVersion: "1.0", Project: "test"}

	tests := []struct {
		status Status
		want   string
	}{
		{StatusCompleted, "✅"},
		{StatusInProgress, "🚧"},
		{StatusPlanned, "📋"},
		{StatusFuture, "💡"},
		{"unknown", ""},
	}

	for _, tt := range tests {
		t.Run(string(tt.status), func(t *testing.T) {
			got := tl.GetStatusEmoji(tt.status)
			if got != tt.want {
				t.Errorf("GetStatusEmoji(%q) = %q, want %q", tt.status, got, tt.want)
			}
		})
	}
}

func TestTasksByType(t *testing.T) {
	tl := &TaskList{
		IRVersion: "1.0",
		Project:   "test",
		Tasks: []Task{
			{ID: "1", Title: "Task 1", Status: StatusCompleted, Type: "Added"},
			{ID: "2", Title: "Task 2", Status: StatusCompleted, Type: "Added"},
			{ID: "3", Title: "Task 3", Status: StatusCompleted, Type: "Changed"},
			{ID: "4", Title: "Task 4", Status: StatusCompleted}, // No type
		},
	}

	byType := tl.TasksByType()

	if len(byType["Added"]) != 2 {
		t.Errorf("TasksByType[Added] = %d tasks, want 2", len(byType["Added"]))
	}
	if len(byType["Changed"]) != 1 {
		t.Errorf("TasksByType[Changed] = %d tasks, want 1", len(byType["Changed"]))
	}
	if len(byType["_unspecified"]) != 1 {
		t.Errorf("TasksByType[_unspecified] = %d tasks, want 1", len(byType["_unspecified"]))
	}
}

func TestStatusOrder(t *testing.T) {
	order := StatusOrder()
	if len(order) != 4 {
		t.Errorf("StatusOrder() returned %d items, want 4", len(order))
	}
	if order[0] != StatusInProgress {
		t.Errorf("StatusOrder()[0] = %q, want %q", order[0], StatusInProgress)
	}
	if order[3] != StatusCompleted {
		t.Errorf("StatusOrder()[3] = %q, want %q", order[3], StatusCompleted)
	}
}

func TestToJSON(t *testing.T) {
	tl := &TaskList{
		IRVersion: "1.0",
		Project:   "test-project",
		Tasks: []Task{
			{ID: "task-1", Title: "Test Task", Status: StatusCompleted},
		},
	}

	data, err := ToJSON(tl)
	if err != nil {
		t.Fatalf("ToJSON() error = %v", err)
	}

	// Parse it back to verify
	tl2, err := Parse(data)
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}

	if tl2.Project != "test-project" {
		t.Errorf("Project = %q, want %q", tl2.Project, "test-project")
	}
	if len(tl2.Tasks) != 1 {
		t.Errorf("Tasks = %d, want 1", len(tl2.Tasks))
	}
}

func TestWriteFile(t *testing.T) {
	tl := &TaskList{
		IRVersion: "1.0",
		Project:   "test-project",
	}

	// Write to temp file
	tmpFile := t.TempDir() + "/tasks.json"
	err := WriteFile(tmpFile, tl)
	if err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	// Read it back
	tl2, err := ParseFile(tmpFile)
	if err != nil {
		t.Fatalf("ParseFile() error = %v", err)
	}

	if tl2.Project != "test-project" {
		t.Errorf("Project = %q, want %q", tl2.Project, "test-project")
	}
}

func TestParseError(t *testing.T) {
	underlying := errors.New("connection refused")
	parseErr := &ParseError{
		Op:  "read",
		Err: underlying,
	}

	expectedStr := "read: connection refused"
	if parseErr.Error() != expectedStr {
		t.Errorf("Error() = %q, want %q", parseErr.Error(), expectedStr)
	}

	if !errors.Is(parseErr, underlying) {
		t.Error("Expected Unwrap to return underlying error")
	}

	unwrapped := parseErr.Unwrap()
	if unwrapped != underlying {
		t.Errorf("Unwrap() = %v, want %v", unwrapped, underlying)
	}
}

func TestValidationError(t *testing.T) {
	err := ValidationError{
		Field:   "ir_version",
		Message: "required field is missing",
	}

	errStr := err.Error()
	if !strings.Contains(errStr, "ir_version") {
		t.Error("Expected error to contain ir_version")
	}
	if !strings.Contains(errStr, "required field is missing") {
		t.Error("Expected error to contain message")
	}
	expected := "ir_version: required field is missing"
	if errStr != expected {
		t.Errorf("Error() = %q, want %q", errStr, expected)
	}
}

func TestValidationResultWithErrors(t *testing.T) {
	result := ValidationResult{
		Valid: false,
		Errors: []ValidationError{
			{Field: "ir_version", Message: "required field is missing"},
			{Field: "project", Message: "required field is missing"},
		},
	}

	if result.Valid {
		t.Error("Expected result.Valid to be false")
	}
	if len(result.Errors) != 2 {
		t.Errorf("Expected 2 errors, got %d", len(result.Errors))
	}
}

func TestValidationResultValid(t *testing.T) {
	result := ValidationResult{
		Valid:  true,
		Errors: nil,
	}

	if !result.Valid {
		t.Error("Expected result.Valid to be true")
	}
	if len(result.Errors) != 0 {
		t.Errorf("Expected 0 errors, got %d", len(result.Errors))
	}
}

func TestTasksByAreaUnspecified(t *testing.T) {
	tl := &TaskList{
		IRVersion: "1.0",
		Project:   "test",
		Tasks: []Task{
			{ID: "1", Title: "Task 1", Status: StatusCompleted, Area: "core"},
			{ID: "2", Title: "Task 2", Status: StatusCompleted}, // No area
		},
	}

	byArea := tl.TasksByArea()
	if len(byArea["core"]) != 1 {
		t.Errorf("TasksByArea[core] = %d tasks, want 1", len(byArea["core"]))
	}
	if len(byArea["_unspecified"]) != 1 {
		t.Errorf("TasksByArea[_unspecified] = %d tasks, want 1", len(byArea["_unspecified"]))
	}
}

func TestPhaseNumbers(t *testing.T) {
	tl := &TaskList{
		IRVersion: "1.0",
		Project:   "test",
		Tasks: []Task{
			{ID: "1", Title: "Task 1", Status: StatusCompleted, Phase: 2},
			{ID: "2", Title: "Task 2", Status: StatusCompleted, Phase: 1},
			{ID: "3", Title: "Task 3", Status: StatusCompleted, Phase: 3},
			{ID: "4", Title: "Task 4", Status: StatusCompleted}, // Phase 0
		},
	}

	phases := tl.PhaseNumbers()
	if len(phases) != 3 {
		t.Errorf("PhaseNumbers() returned %d phases, want 3", len(phases))
	}
	// Should be sorted
	if phases[0] != 1 || phases[1] != 2 || phases[2] != 3 {
		t.Errorf("PhaseNumbers() = %v, want [1, 2, 3]", phases)
	}
}

func TestValidateMoreCases(t *testing.T) {
	tests := []struct {
		name      string
		taskList  *TaskList
		wantValid bool
	}{
		{
			name: "duplicate area ids",
			taskList: &TaskList{
				IRVersion: "1.0",
				Project:   "test",
				Areas: []Area{
					{ID: "core", Name: "Core"},
					{ID: "core", Name: "Core 2"}, // Duplicate
				},
			},
			wantValid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Validate(tt.taskList)
			if result.Valid != tt.wantValid {
				t.Errorf("Validate() valid = %v, want %v", result.Valid, tt.wantValid)
				for _, e := range result.Errors {
					t.Logf("  Error: %s: %s", e.Field, e.Message)
				}
			}
		})
	}
}
