# Enhancement: `mas` CLI

**Status**: Implemented in cmd/mas v0.1.0

## Problem

The Go SDK provides `Renderer.Render()` for box-format output, but there's no standalone CLI to render a `TeamReport` JSON file. Teams that produce reports from LLM agents (which output JSON, not Go structs) need a way to render without writing Go code.

Currently, consumers either:

- Write their own rendering scripts (e.g., Python) that duplicate the SDK logic
- Have agents generate box-formatted text directly (inconsistent, error-prone)

## Solution

Added `cmd/mas/` — a Cobra-based CLI with subcommands for multi-agent-spec operations. The `render` subcommand renders TeamReport JSON to box format and/or Pandoc-friendly Markdown narrative.

```
multi-agent-spec/
├── sdk/go/
│   ├── renderer.go     # Box format renderer
│   └── narrative.go    # Narrative Markdown renderer
├── cmd/
│   └── mas/
│       ├── main.go
│       ├── go.mod
│       └── cmd/
│           ├── root.go
│           └── render.go
```

### Usage

```bash
# Box format to stdout (default)
mas render report.json

# Narrative format to stdout
mas render --format=narrative report.json

# Both formats to separate files
mas render --box-out=report.txt --narrative-out=report.md report.json

# Validate JSON against schema before rendering
mas render --validate report.json

# Read from stdin
cat report.json | mas render --format=narrative

# Version
mas version
```

### Flags

| Flag | Default | Description |
|---|---|---|
| `--format` | `box` | Output format for stdout: `box` or `narrative` |
| `--box-out` | - | Write box format to file |
| `--narrative-out` | - | Write narrative format to file |
| `--validate` | false | Validate JSON against schema before rendering |
| `--schema` | (remote) | JSON Schema URL or file path for validation |

### Narrative Format

The narrative format is Pandoc-friendly Markdown designed for PDF generation:

```bash
mas-render --format=narrative report.json > report.md
pandoc report.md -o report.pdf --pdf-engine=xelatex \
  -V mainfont="Helvetica Neue" \
  -V geometry:margin=1in
```

Features:

- YAML frontmatter with title and date
- No emojis — status rendered as PASS/WARNING/FAIL/SKIP
- Proper heading hierarchy (H1 title, H2 sections, H3 teams, H4 subsections)
- Markdown tables for tasks
- Structured narrative sections (Problem, Analysis, Recommendation)

### Narrative Fields

New fields in `TeamReport` and `TeamSection` for narrative content:

```go
// NarrativeSection holds prose content for narrative reports.
type NarrativeSection struct {
    Problem        string `json:"problem,omitempty"`
    Analysis       string `json:"analysis,omitempty"`
    Recommendation string `json:"recommendation,omitempty"`
}

// Added to TeamSection:
Narrative *NarrativeSection `json:"narrative,omitempty"`

// Added to TeamReport:
Summary    string `json:"summary,omitempty"`    // Executive summary
Conclusion string `json:"conclusion,omitempty"` // Closing section
```

### JSON Schema Validation

Uses `github.com/santhosh-tekuri/jsonschema/v6` for validation against the multi-agent-spec schema. The `--validate` flag checks the input JSON before rendering.

### Build

```bash
cd cmd/mas
go build -o mas .

# Or install globally
go install github.com/plexusone/multi-agent-spec/cmd/mas@latest
```

### Integration with Agent Teams

Agent coordinators write `TeamReport` JSON, then shell out to `mas`:

```bash
# In agent steering: "write JSON to report.json, then render"
mas render report.json
mas render --format=narrative report.json

# Or generate both at once
mas render --box-out=report.txt --narrative-out=report.md report.json
```

## Files Changed

| File | Change |
|------|--------|
| `sdk/go/narrative.go` | New NarrativeRenderer and NarrativeSection |
| `sdk/go/narrative_test.go` | Tests for narrative rendering |
| `sdk/go/report.go` | Added Narrative, Summary, Conclusion fields |
| `cmd/mas/main.go` | CLI entry point |
| `cmd/mas/cmd/root.go` | Root command and version |
| `cmd/mas/cmd/render.go` | Render subcommand |
| `cmd/mas/go.mod` | CLI module with Cobra |
| `schema/report/team-report.schema.json` | Regenerated with narrative fields |

## Motivation

The Saviynt Custom Extension Analysis Team currently uses a Python script (`render_report.py`) to render reports because there's no CLI for the Go SDK. Other teams building on multi-agent-spec will hit the same gap. A shared `mas-render` binary eliminates duplicate rendering implementations.
