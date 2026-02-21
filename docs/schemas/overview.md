# Schema Overview

Multi-Agent Spec uses JSON Schema (draft 2020-12) for all definitions. This enables validation, IDE support, and code generation.

## Core Schemas

| Schema | Purpose | File |
|--------|---------|------|
| [Agent](agent.md) | Individual agent definitions | `agent.schema.json` |
| [Team](team.md) | Multi-agent workflows | `team.schema.json` |
| [Deployment](deployment.md) | Platform-specific configs | `deployment.schema.json` |
| [Report](report.md) | Execution results | `team-report.schema.json` |

## Schema URLs

All schemas are hosted on GitHub:

```
https://raw.githubusercontent.com/agentplexus/multi-agent-spec/main/schema/{category}/{name}.schema.json
```

## Using Schemas

### In JSON Files

Add `$schema` to enable validation:

```json
{
  "$schema": "https://raw.githubusercontent.com/agentplexus/multi-agent-spec/main/schema/agent/agent.schema.json",
  "name": "my-agent",
  "description": "My agent description"
}
```

### In VS Code

Add to `.vscode/settings.json`:

```json
{
  "json.schemas": [
    {
      "fileMatch": ["**/agents/*.json"],
      "url": "https://raw.githubusercontent.com/agentplexus/multi-agent-spec/main/schema/agent/agent.schema.json"
    },
    {
      "fileMatch": ["**/teams/*.json"],
      "url": "https://raw.githubusercontent.com/agentplexus/multi-agent-spec/main/schema/orchestration/team.schema.json"
    }
  ]
}
```

## Go SDK

The Go SDK provides typed structs for all schemas:

```go
import mas "github.com/agentplexus/multi-agent-spec/sdk/go"

agent := &mas.Agent{
    Name:        "analyzer",
    Description: "Analyzes code",
    Model:       "sonnet",
}

team := &mas.Team{
    Name:   "analysis-team",
    Agents: []string{"analyzer"},
}
```

## Design Principles

1. **Go-first** — Schemas are generated from Go types for perfect alignment
2. **Portable** — Platform-agnostic core with platform-specific extensions
3. **Validated** — All fields have types and constraints
4. **Extensible** — `metadata` fields allow custom data without schema changes
