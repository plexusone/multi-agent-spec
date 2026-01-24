package multiagentspec

import (
	"os"
	"path/filepath"
	"testing"
)

func TestParseAgentMarkdown(t *testing.T) {
	input := `---
name: test-agent
description: A test agent for unit testing
model: sonnet
tools: [Read, Write, Bash]
skills: [coding, debugging]
dependencies: [helper-agent]
requires: [go, git]
tasks:
  - id: run-tests
    description: Execute test suite
    type: command
    command: go test ./...
    required: true
    expected_output: PASS
---

# Test Agent

You are a test agent for validating the loader.

## Instructions

Follow these steps:
1. Read the code
2. Run tests
3. Report results
`

	agent, err := ParseAgentMarkdown([]byte(input))
	if err != nil {
		t.Fatalf("ParseAgentMarkdown failed: %v", err)
	}

	if agent.Name != "test-agent" {
		t.Errorf("Name = %q, want %q", agent.Name, "test-agent")
	}

	if agent.Description != "A test agent for unit testing" {
		t.Errorf("Description = %q, want %q", agent.Description, "A test agent for unit testing")
	}

	if agent.Model != ModelSonnet {
		t.Errorf("Model = %q, want %q", agent.Model, ModelSonnet)
	}

	if len(agent.Tools) != 3 {
		t.Errorf("Tools count = %d, want 3", len(agent.Tools))
	}

	if len(agent.Skills) != 2 {
		t.Errorf("Skills count = %d, want 2", len(agent.Skills))
	}

	if len(agent.Dependencies) != 1 {
		t.Errorf("Dependencies count = %d, want 1", len(agent.Dependencies))
	}

	if len(agent.Requires) != 2 {
		t.Errorf("Requires count = %d, want 2", len(agent.Requires))
	}

	if len(agent.Tasks) != 1 {
		t.Errorf("Tasks count = %d, want 1", len(agent.Tasks))
	}

	if agent.Tasks[0].ID != "run-tests" {
		t.Errorf("Task ID = %q, want %q", agent.Tasks[0].ID, "run-tests")
	}

	if agent.Instructions == "" {
		t.Error("Instructions should not be empty")
	}

	if agent.Instructions[:12] != "# Test Agent" {
		t.Errorf("Instructions should start with '# Test Agent', got %q", agent.Instructions[:12])
	}
}

func TestLoadAgentsFromDir(t *testing.T) {
	// Create temp directory with test files
	tmpDir := t.TempDir()

	agent1 := `---
name: agent-one
description: First agent
model: haiku
tools: [Read]
---

Instructions for agent one.
`

	agent2 := `---
name: agent-two
description: Second agent
model: opus
tools: [Write, Bash]
---

Instructions for agent two.
`

	if err := os.WriteFile(filepath.Join(tmpDir, "agent-one.md"), []byte(agent1), 0600); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(tmpDir, "agent-two.md"), []byte(agent2), 0600); err != nil {
		t.Fatal(err)
	}
	// Add a non-md file that should be ignored
	if err := os.WriteFile(filepath.Join(tmpDir, "readme.txt"), []byte("ignore me"), 0600); err != nil {
		t.Fatal(err)
	}

	agents, err := LoadAgentsFromDir(tmpDir)
	if err != nil {
		t.Fatalf("LoadAgentsFromDir failed: %v", err)
	}

	if len(agents) != 2 {
		t.Errorf("Agent count = %d, want 2", len(agents))
	}

	// Check both agents are loaded (order may vary)
	names := make(map[string]bool)
	for _, a := range agents {
		names[a.Name] = true
	}

	if !names["agent-one"] {
		t.Error("agent-one not found")
	}
	if !names["agent-two"] {
		t.Error("agent-two not found")
	}
}

func TestLoadTeamFromFile(t *testing.T) {
	tmpDir := t.TempDir()

	teamJSON := `{
  "name": "test-team",
  "version": "1.0.0",
  "description": "A test team",
  "agents": ["agent-one", "agent-two"],
  "orchestrator": "agent-one",
  "workflow": {
    "type": "dag",
    "steps": [
      {
        "name": "step-one",
        "agent": "agent-one"
      },
      {
        "name": "step-two",
        "agent": "agent-two",
        "depends_on": ["step-one"]
      }
    ]
  }
}`

	path := filepath.Join(tmpDir, "team.json")
	if err := os.WriteFile(path, []byte(teamJSON), 0600); err != nil {
		t.Fatal(err)
	}

	team, err := LoadTeamFromFile(path)
	if err != nil {
		t.Fatalf("LoadTeamFromFile failed: %v", err)
	}

	if team.Name != "test-team" {
		t.Errorf("Name = %q, want %q", team.Name, "test-team")
	}

	if team.Version != "1.0.0" {
		t.Errorf("Version = %q, want %q", team.Version, "1.0.0")
	}

	if len(team.Agents) != 2 {
		t.Errorf("Agents count = %d, want 2", len(team.Agents))
	}

	if team.Orchestrator != "agent-one" {
		t.Errorf("Orchestrator = %q, want %q", team.Orchestrator, "agent-one")
	}

	if team.Workflow == nil {
		t.Fatal("Workflow should not be nil")
	}

	if team.Workflow.Type != WorkflowDAG {
		t.Errorf("Workflow.Type = %q, want %q", team.Workflow.Type, WorkflowDAG)
	}

	if len(team.Workflow.Steps) != 2 {
		t.Errorf("Steps count = %d, want 2", len(team.Workflow.Steps))
	}
}

func TestLoadAgentsFromDirNested(t *testing.T) {
	// Create temp directory with nested structure
	tmpDir := t.TempDir()

	// Create subdirectories
	sharedDir := filepath.Join(tmpDir, "shared")
	prdDir := filepath.Join(tmpDir, "prd")
	if err := os.MkdirAll(sharedDir, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(prdDir, 0755); err != nil {
		t.Fatal(err)
	}

	// Root-level agent (no namespace)
	rootAgent := `---
name: orchestrator
description: Root orchestrator
model: opus
---

Root orchestrator instructions.
`

	// Shared namespace agent
	sharedAgent := `---
name: review-board
description: Shared review board
model: sonnet
---

Review board instructions.
`

	// PRD namespace agent
	prdAgent := `---
name: lead
description: PRD lead agent
model: sonnet
---

PRD lead instructions.
`

	// PRD namespace agent with explicit namespace (should not be overwritten)
	prdAgentExplicit := `---
name: requirements
namespace: custom
description: Requirements agent with explicit namespace
model: haiku
---

Requirements instructions.
`

	if err := os.WriteFile(filepath.Join(tmpDir, "orchestrator.md"), []byte(rootAgent), 0600); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(sharedDir, "review-board.md"), []byte(sharedAgent), 0600); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(prdDir, "lead.md"), []byte(prdAgent), 0600); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(prdDir, "requirements.md"), []byte(prdAgentExplicit), 0600); err != nil {
		t.Fatal(err)
	}

	agents, err := LoadAgentsFromDir(tmpDir)
	if err != nil {
		t.Fatalf("LoadAgentsFromDir failed: %v", err)
	}

	if len(agents) != 4 {
		t.Errorf("Agent count = %d, want 4", len(agents))
	}

	// Build map for easier testing
	agentMap := make(map[string]*Agent)
	for _, a := range agents {
		agentMap[a.QualifiedName()] = a
	}

	// Test root-level agent (no namespace)
	if a, ok := agentMap["orchestrator"]; !ok {
		t.Error("orchestrator not found")
	} else {
		if a.Namespace != "" {
			t.Errorf("orchestrator namespace = %q, want empty", a.Namespace)
		}
	}

	// Test shared namespace agent
	if a, ok := agentMap["shared/review-board"]; !ok {
		t.Error("shared/review-board not found")
	} else {
		if a.Namespace != "shared" {
			t.Errorf("review-board namespace = %q, want %q", a.Namespace, "shared")
		}
	}

	// Test prd namespace agent
	if a, ok := agentMap["prd/lead"]; !ok {
		t.Error("prd/lead not found")
	} else {
		if a.Namespace != "prd" {
			t.Errorf("lead namespace = %q, want %q", a.Namespace, "prd")
		}
	}

	// Test explicit namespace (should NOT be overwritten by directory)
	if a, ok := agentMap["custom/requirements"]; !ok {
		t.Error("custom/requirements not found")
	} else {
		if a.Namespace != "custom" {
			t.Errorf("requirements namespace = %q, want %q", a.Namespace, "custom")
		}
	}
}

func TestLoadAgentsFromDirFlat(t *testing.T) {
	// Create temp directory with nested structure
	tmpDir := t.TempDir()

	// Create a subdirectory that should be ignored
	subDir := filepath.Join(tmpDir, "subdir")
	if err := os.MkdirAll(subDir, 0755); err != nil {
		t.Fatal(err)
	}

	rootAgent := `---
name: root-agent
description: Root agent
---

Root instructions.
`

	subAgent := `---
name: sub-agent
description: Sub agent
---

Sub instructions.
`

	if err := os.WriteFile(filepath.Join(tmpDir, "root-agent.md"), []byte(rootAgent), 0600); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(subDir, "sub-agent.md"), []byte(subAgent), 0600); err != nil {
		t.Fatal(err)
	}

	// LoadAgentsFromDirFlat should only load root-level agents
	agents, err := LoadAgentsFromDirFlat(tmpDir)
	if err != nil {
		t.Fatalf("LoadAgentsFromDirFlat failed: %v", err)
	}

	if len(agents) != 1 {
		t.Errorf("Agent count = %d, want 1 (should ignore subdirectories)", len(agents))
	}

	if agents[0].Name != "root-agent" {
		t.Errorf("Agent name = %q, want %q", agents[0].Name, "root-agent")
	}
}

func TestParseQualifiedName(t *testing.T) {
	tests := []struct {
		input     string
		namespace string
		name      string
	}{
		{"agent-name", "", "agent-name"},
		{"prd/lead", "prd", "lead"},
		{"shared/review-board", "shared", "review-board"},
		{"deep/nested/agent", "deep", "nested/agent"},
	}

	for _, tt := range tests {
		ns, name := ParseQualifiedName(tt.input)
		if ns != tt.namespace || name != tt.name {
			t.Errorf("ParseQualifiedName(%q) = (%q, %q), want (%q, %q)",
				tt.input, ns, name, tt.namespace, tt.name)
		}
	}
}

func TestAgentQualifiedName(t *testing.T) {
	tests := []struct {
		agent *Agent
		want  string
	}{
		{&Agent{Name: "agent-name"}, "agent-name"},
		{&Agent{Name: "lead", Namespace: "prd"}, "prd/lead"},
		{&Agent{Name: "review-board", Namespace: "shared"}, "shared/review-board"},
	}

	for _, tt := range tests {
		got := tt.agent.QualifiedName()
		if got != tt.want {
			t.Errorf("QualifiedName() = %q, want %q", got, tt.want)
		}
	}
}

func TestLoadDeploymentFromFile(t *testing.T) {
	tmpDir := t.TempDir()

	deployJSON := `{
  "team": "test-team",
  "targets": [
    {
      "name": "local-kiro",
      "platform": "kiro-cli",
      "mode": "single-process",
      "output": "plugins/kiro"
    }
  ]
}`

	path := filepath.Join(tmpDir, "deployment.json")
	if err := os.WriteFile(path, []byte(deployJSON), 0600); err != nil {
		t.Fatal(err)
	}

	deployment, err := LoadDeploymentFromFile(path)
	if err != nil {
		t.Fatalf("LoadDeploymentFromFile failed: %v", err)
	}

	if deployment.Team != "test-team" {
		t.Errorf("Team = %q, want %q", deployment.Team, "test-team")
	}

	if len(deployment.Targets) != 1 {
		t.Errorf("Targets count = %d, want 1", len(deployment.Targets))
	}

	if deployment.Targets[0].Platform != PlatformKiroCLI {
		t.Errorf("Platform = %q, want %q", deployment.Targets[0].Platform, PlatformKiroCLI)
	}
}
