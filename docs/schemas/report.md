# Report Schema

Defines structured output from multi-agent team executions using Go/No-Go semantics.

## Schema URL

```
https://raw.githubusercontent.com/agentplexus/multi-agent-spec/main/schema/report/team-report.schema.json
```

## Structure

```json
{
  "$schema": "...",
  "title": "string",
  "project": "string",
  "version": "string",
  "phase": "string",
  "tags": {},
  "teams": [TeamSection],
  "status": "Status",
  "generated_at": "datetime"
}
```

## Status Values

Inspired by NASA mission control Go/No-Go terminology:

| Status | Icon | Description |
|--------|------|-------------|
| `GO` | 🟢 | All checks passed |
| `WARN` | 🟡 | Passed with warnings |
| `NO-GO` | 🔴 | Critical issues found |
| `SKIP` | ⚪ | Check was skipped |

## TeamReport Fields

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `project` | string | Yes | Project identifier |
| `version` | string | Yes | Version being validated |
| `phase` | string | Yes | Workflow phase |
| `teams` | TeamSection[] | Yes | Team/agent results |
| `status` | Status | Yes | Overall status |
| `generated_at` | datetime | Yes | Report timestamp |
| `title` | string | No | Report title |
| `tags` | map[string]string | No | Aggregation tags |
| `summary_blocks` | ContentBlock[] | No | Header content |
| `footer_blocks` | ContentBlock[] | No | Footer content |

### Tags

Tags enable filtering and aggregation across reports:

```json
{
  "tags": {
    "customer": "acme",
    "environment": "production",
    "use_case": "provisioning",
    "target_system": "ad"
  }
}
```

## TeamSection Fields

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `id` | string | Yes | Section identifier |
| `name` | string | Yes | Agent name |
| `status` | Status | Yes | Section status |
| `tasks` | TaskResult[] | No | Task results |
| `verdict` | string | No | Domain-specific verdict |
| `content_blocks` | ContentBlock[] | No | Rich content |

### Verdict

The `verdict` field provides domain-specific labels richer than the 4-value Status:

```json
{
  "status": "NO-GO",
  "verdict": "BLOCKED_PENDING_ENHANCEMENT"
}
```

Examples: `COMPLIANT`, `NON_COMPLIANT`, `NEEDS_WORK`, `APPROVED`, `REJECTED`

## TaskResult Fields

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `id` | string | Yes | Task identifier |
| `status` | Status | Yes | Task status |
| `severity` | string | No | Impact level |
| `detail` | string | No | Result details |
| `duration_ms` | integer | No | Execution time |
| `metadata` | object | No | Custom data |

### Severity

The `severity` field indicates impact level, orthogonal to status:

| Severity | Description |
|----------|-------------|
| `critical` | Must be fixed immediately |
| `high` | Should be fixed before release |
| `medium` | Should be addressed |
| `low` | Nice to fix |
| `info` | Informational only |

Status answers "did it pass?" — Severity answers "how bad is it?"

```json
{
  "id": "sql-injection",
  "status": "NO-GO",
  "severity": "critical",
  "detail": "Found SQL injection in auth module"
}
```

## Example

```json
{
  "$schema": "https://raw.githubusercontent.com/agentplexus/multi-agent-spec/main/schema/report/team-report.schema.json",
  "title": "RELEASE VALIDATION REPORT",
  "project": "my-app",
  "version": "v1.2.0",
  "phase": "PHASE 1: REVIEW",
  "tags": {
    "environment": "staging",
    "release_type": "minor"
  },
  "teams": [
    {
      "id": "security-audit",
      "name": "security",
      "status": "NO-GO",
      "verdict": "BLOCKED_SECURITY_ISSUES",
      "tasks": [
        {
          "id": "hardcoded-secrets",
          "status": "GO",
          "detail": "No hardcoded secrets found"
        },
        {
          "id": "sql-injection",
          "status": "NO-GO",
          "severity": "critical",
          "detail": "SQL injection in UserRepository.findByName()"
        }
      ]
    },
    {
      "id": "qa-validation",
      "name": "qa",
      "status": "GO",
      "tasks": [
        {
          "id": "test-coverage",
          "status": "GO",
          "detail": "Coverage: 87%"
        }
      ]
    }
  ],
  "status": "NO-GO",
  "generated_at": "2024-01-15T10:30:00Z",
  "generated_by": "release-coordinator"
}
```

## Rendering

Use the `mas` CLI to render reports:

```bash
# Box format (terminal)
mas render report.json --format=box

# Narrative format (markdown)
mas render report.json --format=narrative
```

## Go SDK

```go
import mas "github.com/agentplexus/multi-agent-spec/sdk/go"

report := &mas.TeamReport{
    Project: "my-app",
    Version: "v1.2.0",
    Phase:   "PHASE 1: REVIEW",
    Tags: map[string]string{
        "environment": "staging",
    },
    Teams: []mas.TeamSection{
        {
            ID:     "security-audit",
            Name:   "security",
            Status: mas.StatusNoGo,
            Verdict: "BLOCKED_SECURITY_ISSUES",
            Tasks: []mas.TaskResult{
                {
                    ID:       "sql-injection",
                    Status:   mas.StatusNoGo,
                    Severity: "critical",
                    Detail:   "SQL injection found",
                },
            },
        },
    },
    Status:      mas.StatusNoGo,
    GeneratedAt: time.Now(),
}

// Render to box format
renderer := mas.NewBoxRenderer()
output := renderer.Render(report)
```

## See Also

- [mas CLI](../cli/mas.md) — Render reports from command line
- [Team Schema](team.md) — Define team workflows
