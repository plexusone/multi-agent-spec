// Package multiagentspec provides Go types for Multi-Agent Spec definitions.
//
// This package provides structs and utilities for defining multi-agent systems
// with full JSON serialization support.
//
// Example:
//
//	agent := multiagentspec.Agent{
//	    Name:        "my-agent",
//	    Description: "A helpful agent",
//	    Model:       multiagentspec.ModelSonnet,
//	    Tools:       []string{"Read", "Write"},
//	}
//	data, _ := json.MarshalIndent(agent, "", "  ")
package multiagentspec

// Model represents the model capability tier.
type Model string

const (
	ModelHaiku  Model = "haiku"
	ModelSonnet Model = "sonnet"
	ModelOpus   Model = "opus"
)

// Tool represents canonical tool names available to agents.
type Tool string

const (
	ToolWebSearch Tool = "WebSearch"
	ToolWebFetch  Tool = "WebFetch"
	ToolRead      Tool = "Read"
	ToolWrite     Tool = "Write"
	ToolGlob      Tool = "Glob"
	ToolGrep      Tool = "Grep"
	ToolBash      Tool = "Bash"
	ToolEdit      Tool = "Edit"
	ToolTask      Tool = "Task"
)

// TaskType represents how a task is executed.
type TaskType string

const (
	TaskTypeCommand TaskType = "command"
	TaskTypePattern TaskType = "pattern"
	TaskTypeFile    TaskType = "file"
	TaskTypeManual  TaskType = "manual"
)

// Task represents a task that an agent can perform.
type Task struct {
	// ID is the unique task identifier within this agent.
	ID string `json:"id"`

	// Description describes what this task validates or accomplishes.
	Description string `json:"description,omitempty"`

	// Type is how the task is executed (command, pattern, file, manual).
	Type TaskType `json:"type,omitempty"`

	// Command is the shell command to execute (for type: command).
	Command string `json:"command,omitempty"`

	// Pattern is the regex pattern to search for (for type: pattern).
	Pattern string `json:"pattern,omitempty"`

	// File is the file path to check (for type: file).
	File string `json:"file,omitempty"`

	// Files is a glob pattern for files to check (for type: pattern).
	Files string `json:"files,omitempty"`

	// Required indicates if task failure causes agent to report NO-GO.
	Required *bool `json:"required,omitempty"`

	// ExpectedOutput describes what constitutes success.
	ExpectedOutput string `json:"expected_output,omitempty"`

	// HumanInLoop describes when to prompt for human intervention.
	HumanInLoop string `json:"human_in_loop,omitempty"`
}

// DelegationConfig defines delegation permissions for an agent.
type DelegationConfig struct {
	// AllowDelegation enables this agent to delegate work to others.
	AllowDelegation bool `json:"allow_delegation,omitempty" yaml:"allow_delegation,omitempty"`

	// CanDelegateTo lists agent names this agent can delegate to.
	// Empty means no restrictions (can delegate to any agent).
	CanDelegateTo []string `json:"can_delegate_to,omitempty" yaml:"can_delegate_to,omitempty"`

	// CanReceiveFrom lists agent names this agent can receive delegations from.
	// Empty means no restrictions (can receive from any agent).
	CanReceiveFrom []string `json:"can_receive_from,omitempty" yaml:"can_receive_from,omitempty"`
}

// Agent represents an agent definition.
type Agent struct {
	// Name is the unique identifier for the agent (lowercase, hyphenated).
	Name string `json:"name" yaml:"name"`

	// Namespace is the optional namespace for organizing agents.
	// Derived from subdirectory path if not explicitly set in frontmatter.
	Namespace string `json:"namespace,omitempty" yaml:"namespace,omitempty"`

	// Description is a brief summary of what the agent does.
	Description string `json:"description,omitempty" yaml:"description,omitempty"`

	// Icon is the icon identifier for visual representation.
	// Formats: 'brandkit:name' (from brandkit repo), 'lucide:name' (Lucide icon),
	// or plain name for inference.
	Icon string `json:"icon,omitempty" yaml:"icon,omitempty"`

	// Model is the capability tier (haiku, sonnet, opus).
	Model Model `json:"model,omitempty" yaml:"model,omitempty"`

	// Tools are the tools available to this agent.
	Tools []string `json:"tools,omitempty" yaml:"tools,omitempty"`

	// AllowedTools are tools that can execute without user confirmation.
	AllowedTools []string `json:"allowedTools,omitempty" yaml:"allowedTools,omitempty"`

	// Skills are capabilities the agent can invoke.
	Skills []string `json:"skills,omitempty" yaml:"skills,omitempty"`

	// Dependencies are other agents this agent depends on.
	Dependencies []string `json:"dependencies,omitempty" yaml:"dependencies,omitempty"`

	// Requires lists external tools or binaries required (e.g., go, git).
	Requires []string `json:"requires,omitempty" yaml:"requires,omitempty"`

	// Instructions is the system prompt for the agent.
	Instructions string `json:"instructions,omitempty" yaml:"instructions,omitempty"`

	// Tasks are the tasks this agent can perform.
	Tasks []Task `json:"tasks,omitempty" yaml:"tasks,omitempty"`

	// Role-based fields for self-directed workflows

	// Role is the agent's role title (e.g., "Security Analyst").
	Role string `json:"role,omitempty" yaml:"role,omitempty"`

	// Goal describes what the agent aims to achieve.
	Goal string `json:"goal,omitempty" yaml:"goal,omitempty"`

	// Backstory provides context and background for the agent's role.
	Backstory string `json:"backstory,omitempty" yaml:"backstory,omitempty"`

	// Delegation defines delegation permissions for self-directed workflows.
	Delegation *DelegationConfig `json:"delegation,omitempty" yaml:"delegation,omitempty"`
}

// NewAgent creates a new Agent with the given name and description.
func NewAgent(name, description string) *Agent {
	return &Agent{
		Name:        name,
		Description: description,
		Model:       ModelSonnet,
	}
}

// WithModel sets the agent's model and returns the agent for chaining.
func (a *Agent) WithModel(model Model) *Agent {
	a.Model = model
	return a
}

// WithTools sets the agent's tools and returns the agent for chaining.
func (a *Agent) WithTools(tools ...string) *Agent {
	a.Tools = tools
	return a
}

// WithInstructions sets the agent's instructions and returns the agent for chaining.
func (a *Agent) WithInstructions(instructions string) *Agent {
	a.Instructions = instructions
	return a
}

// WithNamespace sets the agent's namespace and returns the agent for chaining.
func (a *Agent) WithNamespace(namespace string) *Agent {
	a.Namespace = namespace
	return a
}

// QualifiedName returns the fully qualified agent name.
// Returns "namespace/name" if namespace is set, otherwise just "name".
func (a *Agent) QualifiedName() string {
	if a.Namespace == "" {
		return a.Name
	}
	return a.Namespace + "/" + a.Name
}

// ParseQualifiedName splits a qualified agent name into namespace and name parts.
// Returns empty namespace if no "/" is present.
//
// Examples:
//
//	ParseQualifiedName("agent-name")        → ("", "agent-name")
//	ParseQualifiedName("prd/lead")          → ("prd", "lead")
//	ParseQualifiedName("shared/review")     → ("shared", "review")
func ParseQualifiedName(qualifiedName string) (namespace, name string) {
	for i := 0; i < len(qualifiedName); i++ {
		if qualifiedName[i] == '/' {
			return qualifiedName[:i], qualifiedName[i+1:]
		}
	}
	return "", qualifiedName
}

// CanDelegate returns true if this agent can delegate work.
func (a *Agent) CanDelegate() bool {
	return a.Delegation != nil && a.Delegation.AllowDelegation
}

// CanDelegateTo returns true if this agent can delegate to the target agent.
func (a *Agent) CanDelegateTo(target string) bool {
	if !a.CanDelegate() {
		return false
	}
	if len(a.Delegation.CanDelegateTo) == 0 {
		return true // No restrictions
	}
	for _, name := range a.Delegation.CanDelegateTo {
		if name == target {
			return true
		}
	}
	return false
}

// CanReceiveFrom returns true if this agent can receive delegations from the source agent.
func (a *Agent) CanReceiveFrom(source string) bool {
	if a.Delegation == nil {
		return true // No delegation config means can receive from anyone
	}
	if len(a.Delegation.CanReceiveFrom) == 0 {
		return true // No restrictions
	}
	for _, name := range a.Delegation.CanReceiveFrom {
		if name == source {
			return true
		}
	}
	return false
}

// WithRole sets the agent's role and returns the agent for chaining.
func (a *Agent) WithRole(role string) *Agent {
	a.Role = role
	return a
}

// WithGoal sets the agent's goal and returns the agent for chaining.
func (a *Agent) WithGoal(goal string) *Agent {
	a.Goal = goal
	return a
}

// WithBackstory sets the agent's backstory and returns the agent for chaining.
func (a *Agent) WithBackstory(backstory string) *Agent {
	a.Backstory = backstory
	return a
}

// WithDelegation sets the agent's delegation config and returns the agent for chaining.
func (a *Agent) WithDelegation(delegation *DelegationConfig) *Agent {
	a.Delegation = delegation
	return a
}
