# Release Notes v0.2.0

**Release Date:** 2026-01-11

## Overview

This release adds **Homebrew distribution** and multi-platform binary releases via GoReleaser. It also improves the overview table with smarter sorting and more concise priority labels.

## Highlights

- Homebrew and multi-platform binary distribution via GoReleaser
- Overview table now sorted by completion percentage (descending) then priority
- Priority labels use concise single-letter format in table cells

## New: Homebrew Distribution

Install `sroadmap` via Homebrew with a single command:

```bash
brew tap grokify/tap
brew install structured-roadmap
```

The formula installs both command names:

- `sroadmap` - Short form for daily use
- `structured-roadmap` - Long form matching the project name

### Multi-Platform Binaries

Pre-built binaries are available for:

| OS | Architecture |
|----|--------------|
| Linux | amd64, arm64 |
| macOS | amd64 (Intel), arm64 (Apple Silicon) |
| Windows | amd64, arm64 |

### Version Command

New `version` command displays build information:

```bash
$ sroadmap version
sroadmap 0.2.0
  commit: abc1234
  built:  2026-01-11T12:00:00Z
```

## Changes

### Overview Table Sorting

The overview table now sorts items more intuitively:

1. **Primary sort:** Completion percentage (descending) - items closest to done appear first
2. **Secondary sort:** Priority (P0 → P1 → P2 → P3) - higher priority items appear first within same completion

This makes it easier to see what's almost done and what needs attention.

### Concise Priority Labels

Priority labels in table cells now use a compact format for better table readability:

| Before | After |
|--------|-------|
| `P0 - Critical` | `P0` |
| `P1 - High` | `P1` |
| `P2 - Medium` | `P2` |
| `P3 - Low` | `P3` |

The full labels remain in item details; only table cells use the concise format.

## Installation

```bash
# Via Go
go install github.com/grokify/structured-roadmap/cmd/sroadmap@v0.2.0

# Via Homebrew
brew upgrade structured-roadmap
```

## Links

- [Full Changelog](CHANGELOG.md)
- [Documentation](https://pkg.go.dev/github.com/grokify/structured-roadmap)
