# Enhancement: Content Blocks in Team Reports

## Problem

The current report schema supports only flat task lists (one status line per task). Complex analysis workflows produce richer output — findings lists, key-value summaries, action items, metrics — that don't fit the `tasks` array model.

Today, agents either:

- Cram everything into the `detail` string field (loses structure)
- Generate box-formatted text directly (inconsistent, unparseable)
- Produce a separate narrative outside the report schema (two disconnected outputs)

## Proposal

Add an optional `content_blocks` array to `team` and `task` objects, plus a top-level `footer_blocks` array on the report. This keeps backward compatibility — existing reports with only `tasks` still render identically.

### Block Types

| Type | Fields | Use Case |
|---|---|---|
| `text` | `title?`, `body` | Descriptions, narratives (auto-wrapped) |
| `kv_pairs` | `title?`, `pairs: [{key, value, icon?}]` | Metadata, summary stats, config |
| `list` | `title?`, `items: [{text, icon?}]` | Findings, action items, recommendations |
| `table` | `title?`, `headers`, `rows` | Comparison matrices, coverage reports |
| `metric` | `label`, `value`, `status`, `target?` | Coverage %, scores |

### Schema Changes

`team-report.schema.json` additions:

```json
{
  "definitions": {
    "content_block": {
      "oneOf": [
        {
          "type": "object",
          "properties": {
            "type": { "const": "text" },
            "title": { "type": "string" },
            "body": { "type": "string" }
          },
          "required": ["type", "body"]
        },
        {
          "type": "object",
          "properties": {
            "type": { "const": "kv_pairs" },
            "title": { "type": "string" },
            "pairs": {
              "type": "array",
              "items": {
                "type": "object",
                "properties": {
                  "key": { "type": "string" },
                  "value": { "type": "string" },
                  "icon": { "type": "string" }
                },
                "required": ["key", "value"]
              }
            }
          },
          "required": ["type", "pairs"]
        },
        {
          "type": "object",
          "properties": {
            "type": { "const": "list" },
            "title": { "type": "string" },
            "items": {
              "type": "array",
              "items": {
                "type": "object",
                "properties": {
                  "text": { "type": "string" },
                  "icon": { "type": "string" }
                },
                "required": ["text"]
              }
            }
          },
          "required": ["type", "items"]
        },
        {
          "type": "object",
          "properties": {
            "type": { "const": "table" },
            "title": { "type": "string" },
            "headers": { "type": "array", "items": { "type": "string" } },
            "rows": { "type": "array", "items": { "type": "array", "items": { "type": "string" } } }
          },
          "required": ["type", "headers", "rows"]
        },
        {
          "type": "object",
          "properties": {
            "type": { "const": "metric" },
            "label": { "type": "string" },
            "value": { "type": "string" },
            "status": { "enum": ["GO", "WARN", "NO-GO", "SKIP"] },
            "target": { "type": "string" }
          },
          "required": ["type", "label", "value", "status"]
        }
      ]
    }
  }
}
```

Add to `team` object:

```json
"content_blocks": {
  "type": "array",
  "items": { "$ref": "#/definitions/content_block" }
}
```

Add to report root:

```json
"footer_blocks": {
  "type": "array",
  "items": { "$ref": "#/definitions/content_block" }
}
```

### Go SDK Changes

Extend `sdk/go` with:

1. `ContentBlock` type (union via `Type` field discriminator)
2. Add `ContentBlocks []ContentBlock` to `TeamResult` struct
3. Add `FooterBlocks []ContentBlock` to `TeamReport` struct
4. Extend `Renderer.Render()` to emit content blocks after task lines
5. Add `Renderer.RenderNarrative()` for Markdown prose output

### Box Format Rendering

Content blocks render inside the box after the team's task lines:

```
║ security-analysis (security)                                                 ║
║   dependency-scan          WARN  5 findings (2 HIGH, 3 MEDIUM)               ║
║   StrictHostKeyChecking disabled (HIGH)                                      ║
║   Outdated JSch 0.1.55 - CVE-2016-5725 (HIGH)                               ║
║   Outdated AWS SDK 1.12.192 (MEDIUM)                                         ║
```

Footer blocks render before the final status line:

```
║ ACTION ITEMS                                                                 ║
║   1: Upgrade JSch to 0.2.18                                                  ║
║   2: Upgrade AWS SDK to 2.25.x                                               ║
```

### Narrative Rendering

New `RenderNarrative()` method produces Markdown:

```markdown
## Security Analysis

**Status**: WARN

### Tasks

| Task | Status | Detail |
|---|---|---|
| dependency-scan | WARN | 5 findings (2 HIGH, 3 MEDIUM) |

### Findings

- StrictHostKeyChecking disabled (HIGH)
- Outdated JSch 0.1.55 - CVE-2016-5725 (HIGH)
- Outdated AWS SDK 1.12.192 (MEDIUM)
```

## Backward Compatibility

- `content_blocks` and `footer_blocks` are optional
- Existing reports without these fields render identically
- Existing consumers that don't read these fields are unaffected
- No changes to `agent-result.schema.json` required (agents produce `content_blocks` in their team output, coordinator passes them through)

## Reference Implementation

A Python prototype renderer exists at:

`gitlab.com/saviynt/product/agents/custom-extension-analysis-team/render_report.py`

This should be replaced by the Go SDK implementation once this enhancement ships.

## Motivation

The Saviynt Custom Extension Analysis Team (6 agents) produces reports with security findings, performance issues, migration blockers, action items, and Aha Ideas. The current schema forces all of this into `detail` strings, losing structure and making programmatic consumption impossible.
