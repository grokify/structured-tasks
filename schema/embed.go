// Package schema provides embedded JSON schema for task list validation.
package schema

import (
	_ "embed"
)

// SchemaV1 contains the embedded JSON schema for task list v1.0.
//
//go:embed tasks.v1.schema.json
var SchemaV1 []byte

// SchemaVersion returns the current schema version.
func SchemaVersion() string {
	return "1.0"
}
