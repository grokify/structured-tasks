package tasks

import (
	"encoding/json"
	"fmt"
	"os"
)

// ParseFile reads and parses a TASKS.json file.
func ParseFile(path string) (*TaskList, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrReadFile, err)
	}
	return Parse(data)
}

// Parse parses JSON data into a TaskList.
func Parse(data []byte) (*TaskList, error) {
	var tl TaskList
	if err := json.Unmarshal(data, &tl); err != nil {
		return nil, fmt.Errorf("%w: %v", ErrParseJSON, err)
	}
	return &tl, nil
}

// WriteFile writes a TaskList to a JSON file.
func WriteFile(path string, tl *TaskList) error {
	data, err := json.MarshalIndent(tl, "", "  ")
	if err != nil {
		return fmt.Errorf("%w: %v", ErrWriteFile, err)
	}
	if err := os.WriteFile(path, data, 0600); err != nil {
		return fmt.Errorf("%w: %v", ErrWriteFile, err)
	}
	return nil
}

// ToJSON converts a TaskList to JSON bytes.
func ToJSON(tl *TaskList) ([]byte, error) {
	return json.MarshalIndent(tl, "", "  ")
}
