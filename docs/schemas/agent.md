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
  "namespace": "string",
  "description": "string",
  "icon": "string",
  "model": "haiku|sonnet|opus",
  "tools": ["string"],
  "allowedTools": ["string"],
  "skills": ["string"],
  "dependencies": ["string"],
  "requires": ["string"],
  "instructions": "string",
  "tasks": [Task],
  "role": "string",
  "goal": "string",
  "backstory": "string",
  "delegation": DelegationConfig
}
```

## Fields

### Required Fields

| Field | Type | Description |
|-------|------|-------------|
| `name` | string | Unique agent identifier (lowercase, hyphenated) |

### Core Fields

| Field | Type | Description |
|-------|------|-------------|
| `description` | string | What the agent does |
| `namespace` | string | Namespace for organizing agents (derived from subdirectory) |
| `icon` | string | Icon identifier (`brandkit:name`, `lucide:name`, or plain name) |
| `model` | string | LLM capability tier: `haiku`, `sonnet`, `opus` |
| `instructions` | string | System prompt for the agent |

### Tool Fields

| Field | Type | Description |
|-------|------|-------------|
| `tools` | string[] | Available tools (Read, Write, Bash, Grep, Glob, etc.) |
| `allowedTools` | string[] | Tools that execute without user confirmation |
| `skills` | string[] | Referenced skill names the agent can invoke |

### Dependency Fields

| Field | Type | Description |
|-------|------|-------------|
| `dependencies` | string[] | Other agents this agent depends on |
| `requires` | string[] | External tools/binaries required (e.g., `go`, `git`) |

### Task Fields

| Field | Type | Description |
|-------|------|-------------|
| `tasks` | Task[] | Validation tasks the agent performs |

### Self-Directed Workflow Fields

These fields enable agents to participate in self-directed workflows (crew, swarm, council).

| Field | Type | Description |
|-------|------|-------------|
| `role` | string | Agent's role title (e.g., "Security Analyst") |
| `goal` | string | What the agent aims to achieve |
| `backstory` | string | Context and background for the role |
| `delegation` | DelegationConfig | Delegation permissions |

## Role-Based Agent Definition

For self-directed workflows, agents need role context to make autonomous decisions:

### Role

The agent's job title or function within the team.

```yaml
role: Security Analyst
```

### Goal

What the agent is trying to accomplish. Guides autonomous decision-making.

```yaml
goal: Identify security vulnerabilities and recommend mitigations
```

### Backstory

Background context that informs the agent's perspective and expertise.

```yaml
backstory: |
  Senior security engineer with 10 years of experience in
  application security, penetration testing, and secure code review.
  Known for finding subtle authentication bypasses and injection flaws.
```

## Delegation Config

Controls how agents can delegate work to each other in self-directed workflows.

```json
{
  "delegation": {
    "allow_delegation": true,
    "can_delegate_to": ["frontend", "backend", "qa"],
    "can_receive_from": ["architect", "lead"]
  }
}
```

### Delegation Fields

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `allow_delegation` | boolean | `false` | Whether this agent can delegate work |
| `can_delegate_to` | string[] | `[]` | Agent names to delegate to (empty = no restrictions) |
| `can_receive_from` | string[] | `[]` | Agent names to receive from (empty = no restrictions) |

### Delegation Patterns

**Lead Agent (can delegate to anyone):**

```yaml
delegation:
  allow_delegation: true
```

**Lead with Restricted Specialists:**

```yaml
delegation:
  allow_delegation: true
  can_delegate_to: [frontend, backend, qa]
```

**Specialist (receives but doesn't delegate):**

```yaml
delegation:
  allow_delegation: false
  can_receive_from: [architect]
```

**Peer Agent (can delegate to specific peers):**

```yaml
delegation:
  allow_delegation: true
  can_delegate_to: [peer-reviewer]
  can_receive_from: [peer-reviewer]
```

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
  "files": "string",
  "required": true,
  "expected_output": "string",
  "human_in_loop": "string"
}
```

### Task Types

| Type | Description | Key Field |
|------|-------------|-----------|
| `pattern` | Regex search in files | `pattern`, `files` |
| `command` | Execute shell command | `command` |
| `file` | Check file existence | `file` |
| `manual` | Human verification | `human_in_loop` |

## Namespace

Agents can be organized into namespaces using subdirectories:

```
specs/agents/
├── pm.md                  # namespace: "" (root)
├── security/
│   ├── scanner.md         # namespace: "security"
│   └── auditor.md         # namespace: "security"
└── qa/
    ├── unit.md            # namespace: "qa"
    └── integration.md     # namespace: "qa"
```

Reference namespaced agents in teams:

```json
{
  "agents": ["pm", "security/scanner", "security/auditor", "qa/unit"]
}
```

## Examples

### Basic Agent (JSON)

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
    }
  ]
}
```

### Self-Directed Agent (Markdown)

```markdown
---
name: security-reviewer
description: Reviews code for security vulnerabilities
model: sonnet
tools:
  - Read
  - Grep
  - Glob
  - Bash

# Self-directed workflow fields
role: Security Analyst
goal: Identify security vulnerabilities and recommend mitigations
backstory: |
  Senior security engineer with 10 years of experience in
  application security, penetration testing, and secure code review.

delegation:
  allow_delegation: false
  can_receive_from:
    - architect
    - lead
---

# Security Review Instructions

When reviewing code:

1. Check for injection vulnerabilities (SQL, XSS, command injection)
2. Validate authentication and authorization logic
3. Look for hardcoded secrets or credentials
4. Assess input validation and sanitization

## Communication

- Share findings via the team broadcast channel
- Challenge other reviewers if you see security implications they missed
- Vote on final recommendations when consensus is requested
```

### Lead Agent with Delegation

```markdown
---
name: architect
description: Technical architect who coordinates development work
model: opus
tools:
  - Read
  - Write
  - Glob
  - Grep
  - Task

role: Technical Architect
goal: Design solutions and coordinate implementation across specialists
backstory: |
  Principal engineer with 15 years of experience designing
  large-scale distributed systems. Expert in breaking down
  complex problems into manageable tasks.

delegation:
  allow_delegation: true
  can_delegate_to:
    - frontend
    - backend
    - qa
    - security
---

# Architecture Instructions

You are the lead architect for this project.

## Responsibilities

1. Break down requirements into technical tasks
2. Delegate implementation work to appropriate specialists
3. Review and approve specialist work before integration
4. Ensure architectural consistency across components

## Delegation Guidelines

- Delegate frontend work to `frontend` specialist
- Delegate backend work to `backend` specialist
- Request security review from `security` for sensitive changes
- Request QA validation from `qa` before merging
```

### Specialist Agent

```markdown
---
name: frontend
description: Frontend specialist for UI implementation
model: sonnet
tools:
  - Read
  - Write
  - Glob
  - Bash

role: Frontend Developer
goal: Implement responsive, accessible user interfaces
backstory: |
  Senior frontend developer specializing in React and TypeScript.
  Passionate about UX and accessibility. Expert in modern CSS
  and component architecture.

delegation:
  allow_delegation: false
  can_receive_from:
    - architect
---

# Frontend Development Instructions

You implement frontend features as delegated by the architect.

## Workflow

1. Receive task assignment from architect
2. Implement the feature following project patterns
3. Write unit tests for new components
4. Report completion back to architect
```

## Go SDK

```go
import mas "github.com/agentplexus/multi-agent-spec/sdk/go"

// Basic agent
agent := mas.NewAgent("security-scanner", "Scans code for vulnerabilities").
    WithModel(mas.ModelSonnet).
    WithTools("Read", "Grep", "Glob")

// Self-directed agent with role
reviewer := mas.NewAgent("security-reviewer", "Reviews code for security").
    WithModel(mas.ModelSonnet).
    WithTools("Read", "Grep", "Glob").
    WithRole("Security Analyst").
    WithGoal("Identify vulnerabilities and recommend mitigations").
    WithBackstory("Senior security engineer with 10 years experience").
    WithDelegation(&mas.DelegationConfig{
        AllowDelegation: false,
        CanReceiveFrom:  []string{"architect"},
    })

// Lead agent with delegation
architect := mas.NewAgent("architect", "Technical architect").
    WithModel(mas.ModelOpus).
    WithTools("Read", "Write", "Task").
    WithRole("Technical Architect").
    WithGoal("Design solutions and coordinate implementation").
    WithDelegation(&mas.DelegationConfig{
        AllowDelegation: true,
        CanDelegateTo:   []string{"frontend", "backend", "qa"},
    })

// Check delegation permissions
if architect.CanDelegateTo("frontend") {
    // Architect can delegate to frontend
}

if reviewer.CanReceiveFrom("architect") {
    // Reviewer can receive work from architect
}
```

## Workflow Type Compatibility

| Field | Deterministic | Crew | Swarm | Council |
|-------|---------------|------|-------|---------|
| `name` | Required | Required | Required | Required |
| `description` | Recommended | Recommended | Recommended | Recommended |
| `model` | Optional | Optional | Optional | Optional |
| `tools` | Required | Required | Required | Required |
| `role` | - | Required | Optional | Required |
| `goal` | - | Required | Optional | Required |
| `backstory` | - | Optional | Optional | Recommended |
| `delegation` | - | Required for lead | - | - |

## See Also

- [Team Schema](team.md) - Orchestrate multiple agents
- [Deployment Schema](deployment.md) - Deploy agents to platforms
- [v0.6.0 Release Notes](../releases/v0.6.0.md) - Self-directed workflow details
