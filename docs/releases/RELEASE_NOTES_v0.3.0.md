# Release Notes - v0.3.0

**Release Date:** 2026-02-09

## Overview

This release brings a **breaking change** to align the JSON schema with structured-changelog's camelCase convention, along with improvements to the status table display.

## Breaking Changes

### JSON Schema Field Names Changed to camelCase

All field names in the JSON schema have been migrated from snake_case to camelCase for consistency with structured-changelog.

#### Migration Guide

Update your TASKS.json files with the following replacements:

| Old (snake_case) | New (camelCase) |
|------------------|-----------------|
| `ir_version` | `irVersion` |
| `generated_at` | `generatedAt` |
| `version_history` | `versionHistory` |
| `target_quarter` | `targetQuarter` |
| `target_version` | `targetVersion` |
| `completed_date` | `completedDate` |
| `depends_on` | `dependsOn` |
| `in_progress` | `inProgress` (status value) |

**Before:**

```json
{
  "ir_version": "1.0",
  "project": "my-project",
  "items": [
    {
      "id": "feature-1",
      "status": "in_progress",
      "target_quarter": "Q2 2026",
      "depends_on": ["feature-0"]
    }
  ]
}
```

**After:**

```json
{
  "irVersion": "1.0",
  "project": "my-project",
  "items": [
    {
      "id": "feature-1",
      "status": "inProgress",
      "targetQuarter": "Q2 2026",
      "dependsOn": ["feature-0"]
    }
  ]
}
```

## Status Table Improvements

Several improvements have been made to the status table output:

- **Phases renumbered from 1** - Previously used 0-indexed numbering
- **Completed phases hidden** - Cleaner output by hiding finished phases
- **Phase column first** - Better visual hierarchy with Phase as the first column

## Dependencies

- Updated `github.com/grokify/structured-changelog` to v0.10.0

## AI Agent Improvements

Added improvements for AI agent usage including better status display for integration with LLM-based development workflows.

## Installation

### Homebrew (macOS/Linux)

```bash
brew upgrade structured-tasks
```

### Go Install

```bash
go install github.com/grokify/structured-tasks/cmd/stasks@v0.3.0
```

## Full Changelog

See [CHANGELOG.md](../../CHANGELOG.md) for the complete list of changes.
