# Example: Stats Agent Team

A complete example of a multi-agent team for statistical research and validation.

## Project Structure

```
examples/stats-agent-team/
├── agents/
│   ├── researcher.json
│   └── validator.json
├── team.json
└── deployment.json
```

## Agent: Researcher

`agents/researcher.json`:

```json
{
  "$schema": "../../schema/agent/agent.schema.json",
  "name": "researcher",
  "description": "Researches statistics on a given topic",
  "model": "sonnet",
  "tools": ["WebSearch", "WebFetch", "Read", "Write"],
  "tasks": [
    {
      "id": "find-sources",
      "description": "Find credible sources for statistics",
      "type": "manual"
    },
    {
      "id": "extract-stats",
      "description": "Extract relevant statistics from sources",
      "type": "manual"
    },
    {
      "id": "format-output",
      "description": "Format statistics as structured JSON",
      "type": "file",
      "file": "stats.json"
    }
  ]
}
```

## Agent: Validator

`agents/validator.json`:

```json
{
  "$schema": "../../schema/agent/agent.schema.json",
  "name": "validator",
  "description": "Validates statistics by cross-referencing sources",
  "model": "sonnet",
  "tools": ["WebFetch", "Read"],
  "tasks": [
    {
      "id": "verify-sources",
      "description": "Verify each statistic appears in claimed source",
      "type": "manual"
    },
    {
      "id": "check-recency",
      "description": "Check statistics are from recent sources",
      "type": "manual"
    },
    {
      "id": "cross-reference",
      "description": "Cross-reference statistics across multiple sources",
      "type": "manual"
    }
  ]
}
```

## Team Definition

`team.json`:

```json
{
  "$schema": "../schema/orchestration/team.schema.json",
  "name": "stats-agent-team",
  "description": "Multi-agent team for statistics research and validation",
  "agents": ["researcher", "validator"],
  "workflow": {
    "steps": [
      {
        "name": "research",
        "agent": "researcher",
        "depends_on": [],
        "outputs": [
          {"name": "statistics"}
        ]
      },
      {
        "name": "validate",
        "agent": "validator",
        "depends_on": ["research"],
        "inputs": [
          {"name": "statistics", "from": "research.statistics"}
        ]
      }
    ]
  }
}
```

## Deployment

`deployment.json`:

```json
{
  "$schema": "../schema/deployment/deployment.schema.json",
  "team": "stats-agent-team",
  "targets": [
    {
      "name": "local-claude",
      "platform": "claude-code",
      "priority": "p1",
      "output": ".claude/agents",
      "claudeCode": {
        "agentDir": ".claude/agents",
        "format": "markdown"
      }
    },
    {
      "name": "local-kiro",
      "platform": "kiro-cli",
      "priority": "p1",
      "output": "plugins/kiro/agents",
      "kiroCli": {
        "pluginDir": "plugins/kiro/agents",
        "format": "json",
        "prefix": "stats_"
      }
    }
  ]
}
```

## Workflow Execution

```
┌────────────┐     ┌────────────┐
│ researcher │────▶│ validator  │
└────────────┘     └────────────┘
   Phase 1            Phase 2
```

1. **Research phase**: Researcher agent finds and extracts statistics
2. **Validation phase**: Validator agent verifies the statistics

## Sample Report

```json
{
  "title": "STATISTICS VALIDATION REPORT",
  "project": "market-research",
  "version": "v1.0.0",
  "phase": "RESEARCH & VALIDATION",
  "tags": {
    "topic": "cloud-computing",
    "year": "2024"
  },
  "teams": [
    {
      "id": "research",
      "name": "researcher",
      "status": "GO",
      "tasks": [
        {"id": "find-sources", "status": "GO", "detail": "Found 5 credible sources"},
        {"id": "extract-stats", "status": "GO", "detail": "Extracted 12 statistics"},
        {"id": "format-output", "status": "GO", "detail": "Output written to stats.json"}
      ]
    },
    {
      "id": "validate",
      "name": "validator",
      "status": "WARN",
      "tasks": [
        {"id": "verify-sources", "status": "GO", "detail": "All sources verified"},
        {"id": "check-recency", "status": "WARN", "severity": "low", "detail": "2 stats from 2022"},
        {"id": "cross-reference", "status": "GO", "detail": "8/12 stats cross-referenced"}
      ]
    }
  ],
  "status": "WARN",
  "generated_at": "2024-01-15T10:30:00Z"
}
```

## Rendering

```bash
# Box format
mas render report.json

# Narrative format
mas render report.json --format=narrative -o report.md
```

## See Also

- [Team Schema](../schemas/team.md) — Workflow definition
- [Report Schema](../schemas/report.md) — Output format
- [Deployment Schema](../schemas/deployment.md) — Platform configs
