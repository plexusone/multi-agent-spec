# Deployment Schema

Defines platform-specific deployment configurations for multi-agent teams.

## Schema URL

```
https://raw.githubusercontent.com/agentplexus/multi-agent-spec/main/schema/deployment/deployment.schema.json
```

## Structure

```json
{
  "$schema": "...",
  "team": "string",
  "targets": [Target]
}
```

## Fields

| Field | Type | Description |
|-------|------|-------------|
| `team` | string | Team name to deploy |
| `targets` | Target[] | Deployment targets |

## Target Definition

```json
{
  "name": "string",
  "platform": "Platform",
  "mode": "DeploymentMode",
  "priority": "Priority",
  "output": "string",
  "runtime": RuntimeConfig,
  "claudeCode": ClaudeCodeConfig,
  "kiroCli": KiroCLIConfig,
  "crewai": CrewAIConfig,
  "awsAgentCore": AWSAgentCoreConfig,
  ...
}
```

### Target Fields

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `name` | string | Yes | Target identifier |
| `platform` | Platform | Yes | Target platform |
| `output` | string | No | Output directory |
| `mode` | DeploymentMode | No | Execution mode |
| `priority` | Priority | No | Deployment priority |

### Platforms

| Platform | Description | Workflow Support |
|----------|-------------|------------------|
| `claude-code` | Claude Code plugins | Deterministic + Self-directed |
| `kiro-cli` | Kiro CLI agents | Deterministic |
| `gemini-cli` | Gemini CLI extensions | Deterministic |
| `aws-agentcore` | AWS AgentCore | Deterministic |
| `crewai` | CrewAI Python framework | Self-directed (crew, swarm) |
| `autogen` | Microsoft AutoGen | Self-directed |
| `kubernetes` | Kubernetes deployment | All |
| `docker-compose` | Docker Compose | All |

### Deployment Modes

| Mode | Description |
|------|-------------|
| `single-process` | All agents in one process |
| `multi-process` | Agents in separate processes |
| `distributed` | Distributed across nodes |
| `serverless` | Serverless functions |

## Platform Configurations

### Claude Code

```json
{
  "claudeCode": {
    "agentDir": ".claude/agents",
    "format": "markdown",
    "team_mode": "team",
    "teammate_mode": "in-process",
    "enable_teams": true
  }
}
```

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `agentDir` | string | `.claude/agents` | Output directory |
| `format` | string | `markdown` | Output format |
| `team_mode` | string | `subagent` | `subagent` or `team` |
| `teammate_mode` | string | `auto` | `in-process`, `tmux`, `auto` |
| `enable_teams` | boolean | `false` | Enable agent teams |

#### Workflow Mapping

| Workflow Type | team_mode | Description |
|---------------|-----------|-------------|
| `chain`, `scatter`, `graph` | `subagent` | Deterministic subagent calls |
| `crew` | `team` | Lead spawns teammates |
| `swarm` | `team` | Self-claim + messaging |
| `council` | `team` | Broadcast + challenge |

### Kiro CLI

```json
{
  "kiroCli": {
    "pluginDir": "~/.kiro/agents",
    "format": "json",
    "prefix": "myteam_"
  }
}
```

| Field | Type | Description |
|-------|------|-------------|
| `pluginDir` | string | Output directory |
| `format` | string | Output format |
| `prefix` | string | Prefix for agent names and files |

The `prefix` field is applied to all agent names, filenames, and steering files for namespace isolation when multiple teams share the same directory.

### CrewAI

```json
{
  "crewai": {
    "model": "claude-3-5-sonnet",
    "processType": "hierarchical",
    "verbose": true,
    "memory": true,
    "allowDelegation": true,
    "managerLlm": "claude-3-opus"
  }
}
```

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `model` | string | - | Default model for agents |
| `processType` | string | `sequential` | `sequential`, `hierarchical`, `consensual` |
| `verbose` | boolean | `false` | Enable verbose logging |
| `memory` | boolean | `false` | Enable agent memory |
| `allowDelegation` | boolean | `true` | Enable agent delegation |
| `managerLlm` | string | - | Model for manager agent (hierarchical) |
| `maxIterations` | integer | - | Maximum crew iterations |

#### Workflow Mapping

| Workflow Type | processType | Description |
|---------------|-------------|-------------|
| `chain` | `sequential` | Sequential execution |
| `crew` | `hierarchical` | Manager delegates to agents |
| `swarm`, `council` | `consensual` | Peer-based with voting |

### AWS AgentCore

```json
{
  "awsAgentCore": {
    "region": "us-east-1",
    "foundationModel": "anthropic.claude-3-sonnet-20240229-v1:0",
    "iac": "cdk",
    "lambdaRuntime": "python3.11"
  }
}
```

| Field | Type | Description |
|-------|------|-------------|
| `region` | string | AWS region |
| `foundationModel` | string | Bedrock model ID |
| `iac` | string | Infrastructure as code (`cdk`, `terraform`, `pulumi`) |
| `lambdaRuntime` | string | Lambda runtime |

### Kubernetes

```json
{
  "kubernetes": {
    "namespace": "agents",
    "helmChart": true,
    "imageRegistry": "gcr.io/my-project",
    "resourceLimits": {
      "cpu": "500m",
      "memory": "512Mi"
    }
  }
}
```

## Examples

### Deterministic Workflow Deployment

```json
{
  "$schema": "https://raw.githubusercontent.com/agentplexus/multi-agent-spec/main/schema/deployment/deployment.schema.json",
  "team": "release-team",
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
      "output": "~/.kiro/agents",
      "kiroCli": {
        "pluginDir": "~/.kiro/agents",
        "format": "json",
        "prefix": "rel_"
      }
    }
  ]
}
```

### Self-Directed Workflow Deployment

```json
{
  "team": "code-review-council",
  "targets": [
    {
      "name": "claude-teams",
      "platform": "claude-code",
      "mode": "multi-process",
      "claudeCode": {
        "agentDir": ".claude/agents",
        "format": "markdown",
        "team_mode": "team",
        "teammate_mode": "in-process",
        "enable_teams": true
      }
    },
    {
      "name": "crewai-local",
      "platform": "crewai",
      "mode": "single-process",
      "crewai": {
        "model": "claude-3-5-sonnet",
        "processType": "consensual",
        "memory": true,
        "verbose": true
      }
    }
  ]
}
```

### Crew Workflow Deployment

```json
{
  "team": "dev-team",
  "targets": [
    {
      "name": "claude-crew",
      "platform": "claude-code",
      "claudeCode": {
        "team_mode": "team",
        "enable_teams": true
      }
    },
    {
      "name": "crewai-hierarchical",
      "platform": "crewai",
      "crewai": {
        "processType": "hierarchical",
        "allowDelegation": true,
        "managerLlm": "claude-3-opus"
      }
    }
  ]
}
```

## Go SDK

```go
import mas "github.com/agentplexus/multi-agent-spec/sdk/go"

// Deterministic deployment
deployment := mas.NewDeployment("release-team").
    AddTarget(mas.Target{
        Name:     "local-claude",
        Platform: mas.PlatformClaudeCode,
        Output:   ".claude/agents",
        ClaudeCode: &mas.ClaudeCodeConfig{
            AgentDir: ".claude/agents",
            Format:   "markdown",
        },
    })

// Self-directed deployment
selfDirectedDeployment := mas.NewDeployment("code-review-council").
    AddTarget(mas.Target{
        Name:     "claude-teams",
        Platform: mas.PlatformClaudeCode,
        ClaudeCode: &mas.ClaudeCodeConfig{
            AgentDir:     ".claude/agents",
            Format:       "markdown",
            TeamMode:     "team",
            TeammateMode: "in-process",
            EnableTeams:  true,
        },
    }).
    AddTarget(mas.Target{
        Name:     "crewai",
        Platform: mas.PlatformCrewAI,
        CrewAI: &mas.CrewAIConfig{
            Model:           "claude-3-5-sonnet",
            ProcessType:     "consensual",
            Memory:          true,
            AllowDelegation: true,
        },
    })
```

## See Also

- [Team Schema](team.md) - Define teams to deploy
- [Agent Schema](agent.md) - Agent definitions with role fields
- [v0.6.0 Release Notes](../releases/v0.6.0.md) - Self-directed workflow details
