# Agent Schema

Defines individual AI agents with their capabilities, tools, and tasks.

## Schema URL

```
https://raw.githubusercontent.com/agentplexus/multi-agent-spec/main/schema/agent/agent.schema.json
```

## Structure

```json
{
  "$schema": "...",
  "name": "string",
  "description": "string",
  "model": "string",
  "tools": ["string"],
  "skills": ["string"],
  "tasks": [Task],
  "metadata": {}
}
```

## Fields

### Required Fields

| Field | Type | Description |
|-------|------|-------------|
| `name` | string | Unique agent identifier |
| `description` | string | What the agent does |

### Optional Fields

| Field | Type | Description |
|-------|------|-------------|
| `model` | string | LLM model (sonnet, opus, haiku) |
| `tools` | string[] | Available tools (Read, Write, Bash, etc.) |
| `skills` | string[] | Referenced skill names |
| `tasks` | Task[] | Validation tasks |
| `metadata` | object | Custom key-value data |

## Task Definition

Tasks define validation checks the agent performs.

```json
{
  "id": "string",
  "description": "string",
  "type": "pattern|command|file|manual",
  "pattern": "string",
  "command": "string",
  "file": "string",
  "files": ["string"],
  "required": true,
  "expected_output": "string",
  "human_in_loop": false
}
```

### Task Types

| Type | Description | Key Field |
|------|-------------|-----------|
| `pattern` | Regex search in files | `pattern` |
| `command` | Execute shell command | `command` |
| `file` | Check file existence | `file` or `files` |
| `manual` | Human verification | `human_in_loop: true` |

## Example

```json
{
  "$schema": "https://raw.githubusercontent.com/agentplexus/multi-agent-spec/main/schema/agent/agent.schema.json",
  "name": "security-scanner",
  "description": "Scans code for security vulnerabilities",
  "model": "sonnet",
  "tools": ["Read", "Grep", "Glob"],
  "tasks": [
    {
      "id": "hardcoded-secrets",
      "description": "Check for hardcoded secrets",
      "type": "pattern",
      "pattern": "(password|secret|api_key)\\s*=\\s*[\"'][^\"']+[\"']",
      "expected_output": "No hardcoded secrets found"
    },
    {
      "id": "sql-injection",
      "description": "Check for SQL injection",
      "type": "pattern",
      "pattern": "fmt\\.Sprintf.*SELECT.*%s",
      "expected_output": "No SQL injection patterns found"
    }
  ]
}
```

## Markdown Format

Agents can also be defined as Markdown with YAML frontmatter:

```markdown
---
name: security-scanner
description: Scans code for security vulnerabilities
model: sonnet
tools:
  - Read
  - Grep
  - Glob
tasks:
  - id: hardcoded-secrets
    description: Check for hardcoded secrets
    type: pattern
    pattern: "(password|secret|api_key)\\s*=\\s*[\"'][^\"']+[\"']"
---

You are a security scanner agent.

## Your Role

Scan the codebase for:

1. Hardcoded secrets and credentials
2. SQL injection vulnerabilities
3. XSS attack vectors

Report findings with severity levels.
```

## Go SDK

```go
import mas "github.com/agentplexus/multi-agent-spec/sdk/go"

agent := &mas.Agent{
    Name:        "security-scanner",
    Description: "Scans code for security vulnerabilities",
    Model:       "sonnet",
    Tools:       []string{"Read", "Grep", "Glob"},
    Tasks: []mas.Task{
        {
            ID:          "hardcoded-secrets",
            Description: "Check for hardcoded secrets",
            Type:        "pattern",
            Pattern:     `(password|secret|api_key)\s*=\s*["'][^"']+["']`,
        },
    },
}
```

## See Also

- [Team Schema](team.md) — Orchestrate multiple agents
- [Deployment Schema](deployment.md) — Deploy agents to platforms
