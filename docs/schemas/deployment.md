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

| Platform | Description |
|----------|-------------|
| `claude-code` | Claude Code plugins |
| `kiro-cli` | Kiro CLI agents |
| `gemini-cli` | Gemini CLI extensions |
| `aws-agentcore` | AWS AgentCore |
| `crewai` | CrewAI Python framework |
| `autogen` | Microsoft AutoGen |
| `kubernetes` | Kubernetes deployment |
| `docker-compose` | Docker Compose |

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
    "format": "markdown"
  }
}
```

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

The `prefix` field is applied to all agent names, filenames, and steering files for namespace isolation when multiple teams share the same directory.

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

## Example

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
    },
    {
      "name": "prod-aws",
      "platform": "aws-agentcore",
      "priority": "p2",
      "output": "cdk",
      "awsAgentCore": {
        "region": "us-east-1",
        "foundationModel": "anthropic.claude-3-sonnet-20240229-v1:0",
        "iac": "cdk",
        "lambdaRuntime": "python3.11"
      }
    }
  ]
}
```

## Go SDK

```go
import mas "github.com/agentplexus/multi-agent-spec/sdk/go"

deployment := mas.NewDeployment("release-team").
    AddTarget(mas.Target{
        Name:     "local-claude",
        Platform: mas.PlatformClaudeCode,
        Output:   ".claude/agents",
        ClaudeCode: &mas.ClaudeCodeConfig{
            AgentDir: ".claude/agents",
            Format:   "markdown",
        },
    }).
    AddTarget(mas.Target{
        Name:     "local-kiro",
        Platform: mas.PlatformKiroCLI,
        Output:   "~/.kiro/agents",
        KiroCLI: &mas.KiroCLIConfig{
            PluginDir: "~/.kiro/agents",
            Format:    "json",
            Prefix:    "rel_",
        },
    })
```

## See Also

- [Team Schema](team.md) — Define teams to deploy
- [AssistantKit Deployment Guide](https://agentplexus.github.io/assistantkit/deployments/overview/) — Generation tool
