# Go SDK

The Go SDK provides typed structs and utilities for working with multi-agent-spec.

## Installation

```bash
go get github.com/agentplexus/multi-agent-spec/sdk/go@latest
```

## Import

```go
import mas "github.com/agentplexus/multi-agent-spec/sdk/go"
```

## Core Types

### Agent

```go
type Agent struct {
    Name        string   `json:"name"`
    Description string   `json:"description"`
    Model       string   `json:"model,omitempty"`
    Tools       []string `json:"tools,omitempty"`
    Skills      []string `json:"skills,omitempty"`
    Tasks       []Task   `json:"tasks,omitempty"`
}
```

### Team

```go
type Team struct {
    Name        string    `json:"name"`
    Description string    `json:"description,omitempty"`
    Agents      []string  `json:"agents"`
    Workflow    *Workflow `json:"workflow"`
}

type Workflow struct {
    Steps []Step `json:"steps"`
}

type Step struct {
    Name      string   `json:"name"`
    Agent     string   `json:"agent"`
    DependsOn []string `json:"depends_on,omitempty"`
    Inputs    []Input  `json:"inputs,omitempty"`
    Outputs   []Output `json:"outputs,omitempty"`
}
```

### Deployment

```go
type Deployment struct {
    Team    string   `json:"team"`
    Targets []Target `json:"targets"`
}

type Target struct {
    Name         string            `json:"name"`
    Platform     Platform          `json:"platform"`
    Output       string            `json:"output,omitempty"`
    ClaudeCode   *ClaudeCodeConfig `json:"claudeCode,omitempty"`
    KiroCLI      *KiroCLIConfig    `json:"kiroCli,omitempty"`
    AWSAgentCore *AWSAgentCoreConfig `json:"awsAgentCore,omitempty"`
    // ... other platform configs
}
```

### TeamReport

```go
type TeamReport struct {
    Title       string            `json:"title,omitempty"`
    Project     string            `json:"project"`
    Version     string            `json:"version"`
    Phase       string            `json:"phase"`
    Tags        map[string]string `json:"tags,omitempty"`
    Teams       []TeamSection     `json:"teams"`
    Status      Status            `json:"status"`
    GeneratedAt time.Time         `json:"generated_at"`
}

type TeamSection struct {
    ID      string       `json:"id"`
    Name    string       `json:"name"`
    Status  Status       `json:"status"`
    Verdict string       `json:"verdict,omitempty"`
    Tasks   []TaskResult `json:"tasks,omitempty"`
}

type TaskResult struct {
    ID       string `json:"id"`
    Status   Status `json:"status"`
    Severity string `json:"severity,omitempty"`
    Detail   string `json:"detail,omitempty"`
}
```

## Status Constants

```go
const (
    StatusGo   Status = "GO"
    StatusWarn Status = "WARN"
    StatusNoGo Status = "NO-GO"
    StatusSkip Status = "SKIP"
)

// Get icon for status
status := mas.StatusGo
icon := status.Icon() // "🟢"
```

## Platform Constants

```go
const (
    PlatformClaudeCode   Platform = "claude-code"
    PlatformKiroCLI      Platform = "kiro-cli"
    PlatformGeminiCLI    Platform = "gemini-cli"
    PlatformAWSAgentCore Platform = "aws-agentcore"
    PlatformCrewAI       Platform = "crewai"
    PlatformAutoGen      Platform = "autogen"
    PlatformKubernetes   Platform = "kubernetes"
)
```

## Creating Reports

```go
report := &mas.TeamReport{
    Project: "my-app",
    Version: "v1.2.0",
    Phase:   "PHASE 1: REVIEW",
    Tags: map[string]string{
        "environment": "staging",
    },
    Teams: []mas.TeamSection{
        {
            ID:      "security",
            Name:    "security",
            Status:  mas.StatusNoGo,
            Verdict: "BLOCKED_SECURITY_ISSUES",
            Tasks: []mas.TaskResult{
                {
                    ID:       "sql-injection",
                    Status:   mas.StatusNoGo,
                    Severity: "critical",
                    Detail:   "Found SQL injection",
                },
            },
        },
    },
    Status:      mas.StatusNoGo,
    GeneratedAt: time.Now(),
}

// Compute overall status from teams
report.Status = report.ComputeOverallStatus()

// Check if all teams pass
if report.IsGo() {
    fmt.Println("Ready for release!")
}
```

## Aggregating Results

```go
// Aggregate multiple agent results into a report
results := []mas.AgentResult{
    {AgentID: "security", Tasks: securityTasks, Status: mas.StatusNoGo},
    {AgentID: "qa", Tasks: qaTasks, Status: mas.StatusGo},
}

report := mas.AggregateResults(results, "my-app", "v1.2.0", "REVIEW")
```

## Rendering Reports

### Box Format

```go
renderer := mas.NewBoxRenderer()
output := renderer.Render(report)
fmt.Println(output)
```

### Narrative Format

```go
renderer := mas.NewNarrativeRenderer()
markdown := renderer.Render(report)
os.WriteFile("report.md", []byte(markdown), 0644)
```

## DAG Sorting

```go
// Sort teams by dependency order
report.SortByDAG()

// Teams are now in topological order:
// 1. Teams with no dependencies
// 2. Teams whose dependencies are satisfied
// 3. Alphabetically within each level
```

## Parsing JSON

```go
// Parse agent result
data, _ := os.ReadFile("agent-result.json")
result, err := mas.ParseAgentResult(data)

// Parse team report
data, _ = os.ReadFile("report.json")
report, err := mas.ParseTeamReport(data)

// Serialize to JSON
jsonBytes, err := report.ToJSON()
```

## Builder Pattern

```go
deployment := mas.NewDeployment("release-team").
    AddTarget(mas.Target{
        Name:     "local-claude",
        Platform: mas.PlatformClaudeCode,
        Output:   ".claude/agents",
    }).
    AddTarget(mas.Target{
        Name:     "local-kiro",
        Platform: mas.PlatformKiroCLI,
        Output:   "~/.kiro/agents",
        KiroCLI: &mas.KiroCLIConfig{
            Prefix: "rel_",
        },
    })
```

## See Also

- [Schema Overview](../schemas/overview.md) — JSON Schema reference
- [mas CLI](../cli/mas.md) — Command-line interface
