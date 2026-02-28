# Quick Start

This guide walks through creating a simple multi-agent team.

## 1. Create Project Structure

```bash
mkdir my-team && cd my-team
mkdir -p specs/{agents,teams,deployments}
```

## 2. Define an Agent

Create `specs/agents/analyzer.md`:

```markdown
---
name: analyzer
description: Analyzes code for patterns and issues
model: sonnet
tools:
  - Read
  - Grep
  - Glob
tasks:
  - id: find-todos
    description: Find TODO comments in code
    type: pattern
    pattern: "TODO|FIXME|HACK"
  - id: check-imports
    description: Check for unused imports
    type: command
    command: "go vet ./..."
---

You are a code analyzer agent.

## Your Role

Analyze the codebase for:

1. TODO comments that need attention
2. Code quality issues
3. Potential bugs

Report findings with severity levels.
```

## 3. Define a Team Workflow

Create `specs/teams/analysis-team.json`:

```json
{
  "$schema": "https://raw.githubusercontent.com/plexusone/multi-agent-spec/main/schema/orchestration/team.schema.json",
  "name": "analysis-team",
  "description": "Code analysis team",
  "agents": ["analyzer"],
  "workflow": {
    "steps": [
      {
        "name": "analyze",
        "agent": "analyzer",
        "depends_on": []
      }
    ]
  }
}
```

## 4. Configure Deployment

Create `specs/deployments/local.json`:

```json
{
  "$schema": "https://raw.githubusercontent.com/plexusone/multi-agent-spec/main/schema/deployment/deployment.schema.json",
  "team": "analysis-team",
  "targets": [
    {
      "name": "local-claude",
      "platform": "claude-code",
      "output": ".claude/agents"
    }
  ]
}
```

## 5. Generate Platform Configs

Use AssistantKit to generate platform-specific configurations:

```bash
assistantkit generate --specs=specs --target=local
```

## 6. Use the Agent

```bash
# Start Claude Code
claude

# The analyzer agent is now available
> Use the analyzer agent to check this codebase
```

## Next Steps

- [Agent Schema](../schemas/agent.md) — Full agent definition reference
- [Team Schema](../schemas/team.md) — Multi-agent workflow reference
- [Deployment Schema](../schemas/deployment.md) — Platform configuration reference
