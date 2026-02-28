# Team Schema

Defines multi-agent teams with workflow orchestration.

## Schema URL

```
https://raw.githubusercontent.com/plexusone/multi-agent-spec/main/schema/orchestration/team.schema.json
```

## Structure

```json
{
  "$schema": "...",
  "name": "string",
  "version": "string",
  "description": "string",
  "agents": ["string"],
  "workflow": {
    "type": "chain|scatter|graph|crew|swarm|council",
    "steps": [Step]
  },
  "collaboration": CollaborationConfig
}
```

## Fields

### Required Fields

| Field | Type | Description |
|-------|------|-------------|
| `name` | string | Team identifier |
| `version` | string | Semantic version |
| `agents` | string[] | List of agent names in this team |

### Optional Fields

| Field | Type | Description |
|-------|------|-------------|
| `description` | string | Team purpose |
| `orchestrator` | string | Orchestrator agent name |
| `workflow` | Workflow | Workflow definition |
| `context` | string | Shared background for all agents |
| `collaboration` | CollaborationConfig | Self-directed workflow config |
| `self_claim` | boolean | Enable task self-claiming (swarm) |
| `plan_approval` | boolean | Require plan approval (crew) |

## Workflow Categories

Multi-Agent Spec supports two workflow paradigms:

| Category | Description | Control |
|----------|-------------|---------|
| **Deterministic** | Execution paths defined in schema | Schema controls |
| **Self-directed** | Agents decide execution paths | Agents control |

## Workflow Types

### Deterministic Workflows

Schema defines exactly which agent runs when.

| Type | Pattern | Use Case |
|------|---------|----------|
| `chain` | A → B → C | Sequential pipeline |
| `scatter` | A → [B,C,D] → E | Parallel with fan-out/fan-in |
| `graph` | DAG | Complex dependencies |

#### Chain Workflow

Steps execute sequentially:

```
┌─────────┐     ┌─────────┐     ┌─────────┐
│ analyze │────▶│ review  │────▶│ report  │
└─────────┘     └─────────┘     └─────────┘
```

```json
{
  "workflow": {
    "type": "chain",
    "steps": [
      {"name": "analyze", "agent": "analyst"},
      {"name": "review", "agent": "reviewer"},
      {"name": "report", "agent": "reporter"}
    ]
  }
}
```

#### Scatter Workflow

Parallel execution with aggregation:

```
                ┌─────────┐
           ┌───▶│ test-1  │───┐
┌─────────┐│    └─────────┘   │┌─────────┐
│ prepare │┼───▶│ test-2  │───┼│ collect │
└─────────┘│    └─────────┘   │└─────────┘
           └───▶│ test-3  │───┘
                └─────────┘
```

```json
{
  "workflow": {
    "type": "scatter",
    "steps": [
      {"name": "prepare", "agent": "coordinator"},
      {"name": "test-1", "agent": "tester", "depends_on": ["prepare"]},
      {"name": "test-2", "agent": "tester", "depends_on": ["prepare"]},
      {"name": "test-3", "agent": "tester", "depends_on": ["prepare"]},
      {"name": "collect", "agent": "coordinator", "depends_on": ["test-1", "test-2", "test-3"]}
    ]
  }
}
```

#### Graph Workflow

Directed acyclic graph with dependencies:

```
┌─────────┐     ┌─────────┐
│ analyze │────▶│ review  │
└─────────┘     └────┬────┘
                     │
┌─────────┐          │
│ security│──────────┼────▶┌─────────┐
└─────────┘          └────▶│ report  │
                           └─────────┘
```

```json
{
  "workflow": {
    "type": "graph",
    "steps": [
      {"name": "analyze", "agent": "analyst"},
      {"name": "security", "agent": "security"},
      {"name": "review", "agent": "reviewer", "depends_on": ["analyze"]},
      {"name": "report", "agent": "reporter", "depends_on": ["review", "security"]}
    ]
  }
}
```

### Self-Directed Workflows

Agents decide execution paths dynamically.

| Type | Pattern | Use Case |
|------|---------|----------|
| `crew` | Lead → Specialists | Delegation hierarchy |
| `swarm` | Shared queue | Self-organizing teams |
| `council` | Peer debate | Consensus decisions |

#### Crew Workflow

Lead agent delegates to specialists:

```
         ┌─────────────┐
         │   lead      │
         └──────┬──────┘
                │ delegates
    ┌───────────┼───────────┐
    ▼           ▼           ▼
┌───────┐   ┌───────┐   ┌───────┐
│ spec1 │   │ spec2 │   │ spec3 │
└───────┘   └───────┘   └───────┘
```

```json
{
  "name": "development-team",
  "version": "1.0.0",
  "agents": ["architect", "frontend", "backend", "qa"],
  "workflow": {"type": "crew"},
  "collaboration": {
    "lead": "architect",
    "specialists": ["frontend", "backend", "qa"]
  },
  "plan_approval": true
}
```

#### Swarm Workflow

Agents self-claim from shared task queue:

```
┌─────────────────────────────┐
│        Task Queue           │
│  [task1] [task2] [task3]    │
└──────────────┬──────────────┘
               │ self-claim
    ┌──────────┼──────────┐
    ▼          ▼          ▼
┌───────┐  ┌───────┐  ┌───────┐
│agent-1│  │agent-2│  │agent-3│
└───────┘  └───────┘  └───────┘
```

```json
{
  "name": "bug-triage",
  "version": "1.0.0",
  "agents": ["triager-1", "triager-2", "triager-3"],
  "workflow": {"type": "swarm"},
  "collaboration": {
    "task_queue": true
  }
}
```

#### Council Workflow

Peer debate with consensus voting:

```
┌───────┐  ┌───────┐  ┌───────┐
│agent-1│◀▶│agent-2│◀▶│agent-3│
└───┬───┘  └───┬───┘  └───┬───┘
    │          │          │
    └──────────┼──────────┘
               ▼
         ┌──────────┐
         │ consensus│
         └──────────┘
```

```json
{
  "name": "architecture-review",
  "version": "1.0.0",
  "agents": ["senior-1", "senior-2", "senior-3"],
  "workflow": {"type": "council"},
  "collaboration": {
    "consensus": {
      "required_agreement": 0.66,
      "max_rounds": 3,
      "tie_breaker": "senior-1"
    }
  }
}
```

## Collaboration Config

Configuration for self-directed workflows:

```json
{
  "collaboration": {
    "lead": "string",
    "specialists": ["string"],
    "task_queue": true,
    "consensus": {
      "required_agreement": 0.66,
      "max_rounds": 3,
      "tie_breaker": "lead"
    },
    "channels": [
      {"name": "team", "type": "broadcast", "participants": ["*"]}
    ]
  }
}
```

### Collaboration Fields

| Field | Type | Description |
|-------|------|-------------|
| `lead` | string | Lead agent name (crew workflow) |
| `specialists` | string[] | Non-delegating specialists |
| `task_queue` | boolean | Enable shared task queue (swarm) |
| `consensus` | ConsensusRules | Consensus config (council) |
| `channels` | Channel[] | Communication channels |

### Consensus Rules

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `required_agreement` | number | 0.5 | Fraction required (0.0-1.0) |
| `max_rounds` | integer | 3 | Max debate rounds |
| `tie_breaker` | string | - | Agent to break ties |

### Channel Types

| Type | Description |
|------|-------------|
| `direct` | Point-to-point between two agents |
| `broadcast` | Send to all participants |
| `pub-sub` | Subscribe to topics |

## Step Definition

```json
{
  "name": "string",
  "agent": "string",
  "depends_on": ["string"],
  "inputs": [Port],
  "outputs": [Port]
}
```

### Step Fields

| Field | Type | Description |
|-------|------|-------------|
| `name` | string | Step identifier |
| `agent` | string | Agent to execute |
| `depends_on` | string[] | Steps that must complete first |
| `inputs` | Port[] | Data inputs |
| `outputs` | Port[] | Data outputs |

### Port Definition

```json
{
  "name": "analysis_result",
  "type": "object",
  "from": "analyze.findings",
  "required": true
}
```

## Go SDK

```go
import mas "github.com/plexusone/multi-agent-spec/sdk/go"

// Deterministic workflow
team := mas.NewTeam("release-team", "1.0.0").
    WithAgents("pm", "qa", "security", "docs").
    WithWorkflow(&mas.Workflow{
        Type: mas.WorkflowGraph,
        Steps: []mas.Step{
            {Name: "pm-review", Agent: "pm"},
            {Name: "qa-validation", Agent: "qa", DependsOn: []string{"pm-review"}},
            {Name: "security-audit", Agent: "security", DependsOn: []string{"pm-review"}},
            {Name: "docs-review", Agent: "docs", DependsOn: []string{"qa-validation", "security-audit"}},
        },
    })

// Self-directed workflow
crewTeam := mas.NewTeam("dev-team", "1.0.0").
    WithAgents("architect", "frontend", "backend").
    WithWorkflow(&mas.Workflow{Type: mas.WorkflowCrew}).
    WithCollaboration(&mas.CollaborationConfig{
        Lead:        "architect",
        Specialists: []string{"frontend", "backend"},
    }).
    WithPlanApproval(true)

// Check workflow category
if team.IsDeterministic() {
    // Schema controls execution
}
if crewTeam.IsSelfDirected() {
    // Agents control execution
}

// Validate configuration
if err := crewTeam.Validate(); err != nil {
    log.Fatal(err)
}
```

## See Also

- [Agent Schema](agent.md) - Define individual agents
- [Deployment Schema](deployment.md) - Deploy teams to platforms
- [Report Schema](report.md) - Team execution results
- [Message Schema](../releases/v0.6.0.md#message-schema) - Inter-agent messaging
