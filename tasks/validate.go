package tasks

import (
	"fmt"

	"github.com/grokify/structured-changelog/changelog"
)

// ValidationError represents a validation error.
type ValidationError struct {
	Field   string
	Message string
}

func (e ValidationError) Error() string {
	return fmt.Sprintf("%s: %s", e.Field, e.Message)
}

// ValidationResult holds the results of validation.
type ValidationResult struct {
	Valid  bool
	Errors []ValidationError
}

// Validate checks a TaskList for validity.
func Validate(tl *TaskList) ValidationResult {
	result := ValidationResult{Valid: true}

	// Required fields
	if tl.IRVersion == "" {
		result.addError("ir_version", "required field is missing")
	} else if tl.IRVersion != "1.0" {
		result.addError("ir_version", fmt.Sprintf("unsupported version: %s", tl.IRVersion))
	}

	if tl.Project == "" {
		result.addError("project", "required field is missing")
	}

	// Validate tasks
	taskIDs := make(map[string]bool)
	for i, task := range tl.Tasks {
		prefix := fmt.Sprintf("tasks[%d]", i)

		if task.ID == "" {
			result.addError(prefix+".id", "required field is missing")
		} else if taskIDs[task.ID] {
			result.addError(prefix+".id", fmt.Sprintf("duplicate ID: %s", task.ID))
		} else {
			taskIDs[task.ID] = true
		}

		if task.Title == "" {
			result.addError(prefix+".title", "required field is missing")
		}

		if task.Status == "" {
			result.addError(prefix+".status", "required field is missing")
		} else if !isValidStatus(task.Status) {
			result.addError(prefix+".status", fmt.Sprintf("invalid status: %s", task.Status))
		}

		// Validate phase is non-negative
		if task.Phase < 0 {
			result.addError(prefix+".phase", "phase must be non-negative")
		}

		// Validate type against structured-changelog change types
		if task.Type != "" {
			if !changelog.DefaultRegistry.IsValidName(task.Type) {
				result.addError(prefix+".type", fmt.Sprintf("invalid change type: %s (see structured-changelog for valid types)", task.Type))
			}
		}

		// Validate subtasks
		for j, subtask := range task.Subtasks {
			subtaskPrefix := fmt.Sprintf("%s.subtasks[%d]", prefix, j)
			if subtask.Description == "" {
				result.addError(subtaskPrefix+".description", "required field is missing")
			}
		}
	}

	// Validate depends_on references
	for i, task := range tl.Tasks {
		for _, dep := range task.DependsOn {
			if !taskIDs[dep] {
				result.addError(fmt.Sprintf("tasks[%d].depends_on", i), fmt.Sprintf("references unknown task: %s", dep))
			}
		}
	}

	// Validate areas
	areaIDs := make(map[string]bool)
	for i, area := range tl.Areas {
		prefix := fmt.Sprintf("areas[%d]", i)
		if area.ID == "" {
			result.addError(prefix+".id", "required field is missing")
		} else if areaIDs[area.ID] {
			result.addError(prefix+".id", fmt.Sprintf("duplicate ID: %s", area.ID))
		} else {
			areaIDs[area.ID] = true
		}
		if area.Name == "" {
			result.addError(prefix+".name", "required field is missing")
		}
	}

	// Validate task area references
	for i, task := range tl.Tasks {
		if task.Area != "" && len(tl.Areas) > 0 && !areaIDs[task.Area] {
			result.addError(fmt.Sprintf("tasks[%d].area", i), fmt.Sprintf("references unknown area: %s", task.Area))
		}
	}

	return result
}

func (r *ValidationResult) addError(field, message string) {
	r.Errors = append(r.Errors, ValidationError{Field: field, Message: message})
	r.Valid = false
}

func isValidStatus(s Status) bool {
	switch s {
	case StatusCompleted, StatusInProgress, StatusPlanned, StatusFuture:
		return true
	}
	return false
}
