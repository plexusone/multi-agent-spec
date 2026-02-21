# Installation

## Go SDK

```bash
go get github.com/agentplexus/multi-agent-spec/sdk/go@latest
```

## CLI

```bash
go install github.com/agentplexus/multi-agent-spec/cmd/mas@latest
```

Verify installation:

```bash
mas --help
```

## JSON Schemas

Schemas are available at:

```
https://raw.githubusercontent.com/agentplexus/multi-agent-spec/main/schema/agent/agent.schema.json
https://raw.githubusercontent.com/agentplexus/multi-agent-spec/main/schema/orchestration/team.schema.json
https://raw.githubusercontent.com/agentplexus/multi-agent-spec/main/schema/deployment/deployment.schema.json
https://raw.githubusercontent.com/agentplexus/multi-agent-spec/main/schema/report/team-report.schema.json
```

Reference them in your JSON files:

```json
{
  "$schema": "https://raw.githubusercontent.com/agentplexus/multi-agent-spec/main/schema/agent/agent.schema.json",
  "name": "my-agent",
  ...
}
```

## Requirements

- Go 1.21+ (for SDK and CLI)
- Node.js 18+ (for TypeScript SDK, optional)
