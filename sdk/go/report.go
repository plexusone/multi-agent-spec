package multiagentspec

import (
	"encoding/json"
	"time"
)

// Status represents the validation status following NASA Go/No-Go terminology.
type Status string

const (
	StatusGo   Status = "GO"
	StatusWarn Status = "WARN"
	StatusNoGo Status = "NO-GO"
	StatusSkip Status = "SKIP"
)

// Icon returns the UTF-8 icon for the status.
func (s Status) Icon() string {
	switch s {
	case StatusGo:
		return "\U0001F7E2" // 🟢
	case StatusWarn:
		return "\U0001F7E1" // 🟡
	case StatusNoGo:
		return "\U0001F534" // 🔴
	case StatusSkip:
		return "\u26AA" // ⚪
	default:
		return "?"
	}
}

// AgentResult is the JSON-serializable output from each validation agent.
// This is the intermediate representation that agents produce and the
// coordinator consumes to build the final TeamReport.
type AgentResult struct {
	// Schema is the JSON schema URL for validation
	Schema string `json:"$schema,omitempty"`

	// AgentID identifies the agent (e.g., "pm", "qa", "documentation")
	AgentID string `json:"agent_id"`

	// StepID is the workflow step name (e.g., "pm-validation", "qa-validation")
	StepID string `json:"step_id"`

	// Inputs are values received from upstream agents in the DAG
	Inputs map[string]interface{} `json:"inputs,omitempty"`

	// Outputs are values produced by this agent for downstream agents
	Outputs map[string]interface{} `json:"outputs,omitempty"`

	// Checks are the individual validation checks performed
	Checks []Check `json:"checks"`

	// Status is the overall status for this agent (computed from checks)
	Status Status `json:"status"`

	// ExecutedAt is when the agent completed execution
	ExecutedAt time.Time `json:"executed_at"`

	// AgentModel is the LLM model used (e.g., "sonnet", "haiku", "opus")
	AgentModel string `json:"agent_model,omitempty"`

	// Duration is how long the agent took to execute
	Duration string `json:"duration,omitempty"`

	// Error is set if the agent failed to execute
	Error string `json:"error,omitempty"`
}

// Check represents a single validation check result.
type Check struct {
	// ID is the check identifier (e.g., "build", "tests", "version-recommendation")
	ID string `json:"id"`

	// Status is GO, WARN, NO-GO, or SKIP
	Status Status `json:"status"`

	// Detail is optional additional information
	Detail string `json:"detail,omitempty"`

	// Metadata allows checks to include structured data
	Metadata map[string]interface{} `json:"metadata,omitempty"`
}

// TeamSection represents a team/agent section in the report.
type TeamSection struct {
	// ID is the workflow step ID (e.g., "pm-validation")
	ID string `json:"id"`

	// Name is the agent name (e.g., "pm")
	Name string `json:"name"`

	// AgentID matches the agent definition in team.json
	AgentID string `json:"agent_id,omitempty"`

	// Model is the LLM model used
	Model string `json:"model,omitempty"`

	// Checks are the validation checks for this team
	Checks []Check `json:"checks"`

	// Status is the overall status (computed from checks)
	Status Status `json:"status"`
}

// TeamReport is the complete JSON-serializable report.
// This is what the coordinator produces by aggregating AgentResults.
type TeamReport struct {
	// Schema is the JSON schema URL for validation
	Schema string `json:"$schema,omitempty"`

	// Project is the repository identifier
	Project string `json:"project"`

	// Version is the target release version
	Version string `json:"version"`

	// Target is a human-readable target description
	Target string `json:"target,omitempty"`

	// Phase is the workflow phase (e.g., "PHASE 1: REVIEW")
	Phase string `json:"phase"`

	// Teams are the validation teams/agents
	Teams []TeamSection `json:"teams"`

	// Status is the overall status (computed from teams)
	Status Status `json:"status"`

	// GeneratedAt is when the report was generated
	GeneratedAt time.Time `json:"generated_at"`

	// GeneratedBy identifies the coordinator
	GeneratedBy string `json:"generated_by,omitempty"`
}

// ComputeStatus computes the overall status from checks.
func (a *AgentResult) ComputeStatus() Status {
	return computeStatusFromChecks(a.Checks)
}

// ToTeamSection converts an AgentResult to a TeamSection for the report.
func (a *AgentResult) ToTeamSection() TeamSection {
	return TeamSection{
		ID:      a.StepID,
		Name:    a.AgentID,
		AgentID: a.AgentID,
		Model:   a.AgentModel,
		Checks:  a.Checks,
		Status:  a.ComputeStatus(),
	}
}

// OverallStatus computes the overall status for a team section.
func (t *TeamSection) OverallStatus() Status {
	return computeStatusFromChecks(t.Checks)
}

// ComputeOverallStatus computes the overall status from all teams.
func (r *TeamReport) ComputeOverallStatus() Status {
	hasNoGo := false
	hasWarn := false

	for _, t := range r.Teams {
		switch t.Status {
		case StatusNoGo:
			hasNoGo = true
		case StatusWarn:
			hasWarn = true
		}
	}

	if hasNoGo {
		return StatusNoGo
	}
	if hasWarn {
		return StatusWarn
	}
	return StatusGo
}

// IsGo returns true if all teams pass validation.
func (r *TeamReport) IsGo() bool {
	for _, t := range r.Teams {
		if t.Status == StatusNoGo {
			return false
		}
	}
	return true
}

// FinalMessage returns the final status message for display.
func (r *TeamReport) FinalMessage() string {
	if r.IsGo() {
		return "\U0001F680 TEAM: GO for " + r.Version + " \U0001F680" // 🚀 TEAM: GO for vX.Y.Z 🚀
	}
	return "\U0001F6D1 TEAM: NO-GO for " + r.Version + " \U0001F6D1" // 🛑 TEAM: NO-GO for vX.Y.Z 🛑
}

// ToJSON serializes the report to JSON.
func (r *TeamReport) ToJSON() ([]byte, error) {
	return json.MarshalIndent(r, "", "  ")
}

// AggregateResults combines multiple AgentResults into a TeamReport.
func AggregateResults(results []AgentResult, project, version, phase string) *TeamReport {
	teams := make([]TeamSection, 0, len(results))
	for _, r := range results {
		teams = append(teams, r.ToTeamSection())
	}

	report := &TeamReport{
		Schema:      "https://raw.githubusercontent.com/agentplexus/multi-agent-spec/main/schema/report/team-report.schema.json",
		Project:     project,
		Version:     version,
		Target:      version,
		Phase:       phase,
		Teams:       teams,
		GeneratedAt: time.Now().UTC(),
		GeneratedBy: "release-coordinator",
	}

	report.Status = report.ComputeOverallStatus()

	return report
}

// ParseAgentResult parses JSON into an AgentResult.
func ParseAgentResult(data []byte) (*AgentResult, error) {
	var result AgentResult
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// ParseTeamReport parses JSON into a TeamReport.
func ParseTeamReport(data []byte) (*TeamReport, error) {
	var report TeamReport
	if err := json.Unmarshal(data, &report); err != nil {
		return nil, err
	}
	return &report, nil
}

// computeStatusFromChecks is a helper to compute status from a slice of checks.
func computeStatusFromChecks(checks []Check) Status {
	hasNoGo := false
	hasWarn := false
	allSkipped := true

	for _, c := range checks {
		if c.Status != StatusSkip {
			allSkipped = false
		}
		switch c.Status {
		case StatusNoGo:
			hasNoGo = true
		case StatusWarn:
			hasWarn = true
		}
	}

	if allSkipped {
		return StatusSkip
	}
	if hasNoGo {
		return StatusNoGo
	}
	if hasWarn {
		return StatusWarn
	}
	return StatusGo
}
