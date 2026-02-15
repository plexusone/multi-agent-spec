# Enhancement: Content Blocks in Team Reports

**Status**: Implemented in sdk/go v0.6.0

## Problem

The current report schema supports only flat task lists (one status line per task). Complex analysis workflows produce richer output — findings lists, key-value summaries, action items, metrics — that don't fit the `tasks` array model.

Today, agents either:

- Cram everything into the `detail` string field (loses structure)
- Generate box-formatted text directly (inconsistent, unparseable)
- Produce a separate narrative outside the report schema (two disconnected outputs)

## Solution

Added optional content block arrays to TeamReport and TeamSection, plus constructors and rendering support in the Go SDK. Existing reports with only `tasks` still render identically (backward compatible).

### Block Types

| Type | Fields | Use Case |
|---|---|---|
| `text` | `title?`, `content` | Descriptions, narratives (auto-wrapped) |
| `kv_pairs` | `title?`, `pairs: [{key, value, icon?}]` | Metadata, summary stats, config |
| `list` | `title?`, `items: [{text, icon?, status?}]` | Findings, action items, recommendations |
| `table` | `title?`, `headers`, `rows` | Comparison matrices, coverage reports |
| `metric` | `label`, `value`, `status`, `target?` | Coverage %, scores |

### Go SDK Types

```go
// ContentBlockType discriminates content block variants.
type ContentBlockType string

const (
    ContentBlockKVPairs ContentBlockType = "kv_pairs"
    ContentBlockList    ContentBlockType = "list"
    ContentBlockTable   ContentBlockType = "table"
    ContentBlockText    ContentBlockType = "text"
    ContentBlockMetric  ContentBlockType = "metric"
)

// ContentBlock represents rich content within a report section.
type ContentBlock struct {
    Type    ContentBlockType `json:"type"`
    Title   string           `json:"title,omitempty"`
    Pairs   []KVPair         `json:"pairs,omitempty"`   // for kv_pairs
    Items   []ListItem       `json:"items,omitempty"`   // for list
    Headers []string         `json:"headers,omitempty"` // for table
    Rows    [][]string       `json:"rows,omitempty"`    // for table
    Content string           `json:"content,omitempty"` // for text
    Label   string           `json:"label,omitempty"`   // for metric
    Value   string           `json:"value,omitempty"`   // for metric
    Status  Status           `json:"status,omitempty"`  // for metric
    Target  string           `json:"target,omitempty"`  // for metric
}

// KVPair is a key-value pair with optional icon.
type KVPair struct {
    Key   string `json:"key"`
    Value string `json:"value"`
    Icon  string `json:"icon,omitempty"`
}

// ListItem is a list entry with optional icon and status.
type ListItem struct {
    Text   string `json:"text"`
    Icon   string `json:"icon,omitempty"`
    Status Status `json:"status,omitempty"`
}
```

### Constructor Functions

```go
NewKVPairsBlock(title string, pairs ...KVPair) ContentBlock
NewListBlock(title string, items ...ListItem) ContentBlock
NewTextBlock(title, content string) ContentBlock
NewTableBlock(title string, headers []string, rows [][]string) ContentBlock
NewMetricBlock(label, value string, status Status, target string) ContentBlock
```

### TeamReport Extensions

```go
type TeamReport struct {
    // ... existing fields ...

    // Title is the report title (defaults to "TEAM STATUS REPORT").
    Title string `json:"title,omitempty"`

    // SummaryBlocks appear after the header, before the phase.
    SummaryBlocks []ContentBlock `json:"summary_blocks,omitempty"`

    // FooterBlocks appear after all teams, before the final message.
    FooterBlocks []ContentBlock `json:"footer_blocks,omitempty"`
}

type TeamSection struct {
    // ... existing fields ...

    // Tasks is now optional (can use ContentBlocks instead).
    Tasks []TaskResult `json:"tasks,omitempty"`

    // ContentBlocks holds rich content for this team section.
    ContentBlocks []ContentBlock `json:"content_blocks,omitempty"`
}
```

### Box Format Rendering

Content blocks render inside the box after the team's task lines:

```
╠══════════════════════════════════════════════════════════════════════════════╣
║ Security Analysis (security-analysis)                                        ║
║   dependency-scan          🟡 WARN  5 findings (2 HIGH, 3 MEDIUM)            ║
║ 🔴 StrictHostKeyChecking disabled (HIGH)                                     ║
║ 🔴 Outdated JSch 0.1.55 - CVE-2016-5725 (HIGH)                               ║
║ 🟡 Outdated AWS SDK 1.12.192 (MEDIUM)                                        ║
```

Footer blocks render before the final status line:

```
╠══════════════════════════════════════════════════════════════════════════════╣
║ ACTION ITEMS                                                                 ║
║ 🔴 1: Upgrade JSch to 0.2.18                                                 ║
║ 🟡 2: Upgrade AWS SDK to 2.25.x                                              ║
╠══════════════════════════════════════════════════════════════════════════════╣
║                       🚀 TEAM: GO for v1.0.0 🚀                              ║
╚══════════════════════════════════════════════════════════════════════════════╝
```

Metrics with targets render as:

```
║ 🟢 Coverage: 85% (target: 80%)                                               ║
```

## Schema

The JSON schema is auto-generated from Go types. See `schema/report/team-report.schema.json`.

Key definitions:

- `ContentBlock` - union type with `type` discriminator
- `ContentBlockType` - enum: `kv_pairs`, `list`, `table`, `text`, `metric`
- `KVPair` - key/value with optional icon
- `ListItem` - text with optional icon/status

## Backward Compatibility

- `content_blocks`, `summary_blocks`, and `footer_blocks` are optional
- `tasks` field is now optional on TeamSection
- Existing reports without these fields render identically
- Existing consumers that don't read these fields are unaffected

## Files Changed

| File | Change |
|------|--------|
| `sdk/go/content_block.go` | New types and constructors |
| `sdk/go/content_block_test.go` | Tests |
| `sdk/go/report.go` | Extended TeamReport, TeamSection, AgentResult |
| `sdk/go/renderer.go` | Content block rendering |
| `sdk/go/jsonschema.go` | ContentBlockType schema |
| `schema/report/team-report.schema.json` | Regenerated |

## Future Work

- `RenderNarrative()` method for Markdown output
- Content blocks on TaskResult (if needed)

## Motivation

The Saviynt Custom Extension Analysis Team (6 agents) produces reports with security findings, performance issues, migration blockers, action items, and Aha Ideas. The current schema forces all of this into `detail` strings, losing structure and making programmatic consumption impossible.
