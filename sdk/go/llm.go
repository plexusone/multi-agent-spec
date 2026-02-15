package multiagentspec

// EvaluationType discriminates between rule-based and LLM evaluation.
type EvaluationType string

const (
	// EvaluationTypeRule indicates deterministic rule-based evaluation.
	EvaluationTypeRule EvaluationType = "rule"

	// EvaluationTypeLLM indicates LLM-based semantic evaluation.
	EvaluationTypeLLM EvaluationType = "llm"

	// EvaluationTypeCombined indicates both rule-based and LLM evaluation.
	EvaluationTypeCombined EvaluationType = "combined"
)

// LLMEvaluation contains LLM-based evaluation results.
// This is the nested "llm" object in evaluation results when
// EvaluationType is "llm" or "combined".
type LLMEvaluation struct {
	// Score is the numeric score from LLM evaluation (typically 0-10).
	Score float64 `json:"score"`

	// MaxScore is the maximum possible score (default: 10).
	MaxScore float64 `json:"maxScore,omitempty"`

	// Confidence is the LLM's confidence in the evaluation (0-1 scale).
	Confidence float64 `json:"confidence,omitempty"`

	// Reasoning is the LLM's explanation of the score.
	Reasoning string `json:"reasoning,omitempty"`

	// Strengths are positive aspects identified by the LLM.
	Strengths []string `json:"strengths,omitempty"`

	// Concerns are issues or problems identified by the LLM.
	Concerns []string `json:"concerns,omitempty"`

	// Suggestions are actionable recommendations for improvement.
	Suggestions []string `json:"suggestions,omitempty"`

	// Model is the LLM model used (e.g., "claude-sonnet-4", "gpt-4").
	Model string `json:"model"`

	// Provider is the LLM provider (e.g., "anthropic", "openai", "bedrock").
	Provider string `json:"provider,omitempty"`

	// TokensUsed is the total tokens consumed for this evaluation.
	TokensUsed int `json:"tokensUsed,omitempty"`

	// LatencyMs is the LLM API response time in milliseconds.
	LatencyMs int `json:"latencyMs,omitempty"`

	// PromptVersion is a version identifier for the evaluation prompt.
	// Used for reproducibility and prompt iteration tracking.
	PromptVersion string `json:"promptVersion,omitempty"`
}

// CombinedWeights specifies weighting for combined rule + LLM evaluation.
type CombinedWeights struct {
	// Rule is the weight for rule-based score (0-1).
	Rule float64 `json:"rule"`

	// LLM is the weight for LLM score (0-1).
	LLM float64 `json:"llm"`
}

// Issue represents a specific problem identified in evaluation.
// Used in narrative reports for detailed fix guidance.
type Issue struct {
	// ID is the issue identifier (e.g., "ISS-001").
	ID string `json:"id"`

	// Category is the evaluation category this issue belongs to.
	Category string `json:"category"`

	// Severity indicates how serious the issue is.
	Severity Severity `json:"severity"`

	// Problem describes what the issue is.
	Problem string `json:"problem"`

	// Location indicates where in the document the issue occurs
	// (e.g., "requirements.functional[2]", "executiveSummary.problemStatement").
	Location string `json:"location,omitempty"`

	// Analysis explains why this is a problem.
	Analysis string `json:"analysis,omitempty"`

	// Recommendation describes how to fix the issue.
	Recommendation string `json:"recommendation,omitempty"`

	// Example provides sample improved text or structure.
	Example string `json:"example,omitempty"`

	// Effort estimates the work required to fix this issue.
	Effort Effort `json:"effort,omitempty"`

	// RelatedIssues lists IDs of related issues.
	RelatedIssues []string `json:"relatedIssues,omitempty"`
}

// Severity indicates how serious an issue is.
type Severity string

const (
	// SeverityCritical indicates a blocking issue that must be fixed.
	SeverityCritical Severity = "critical"

	// SeverityMajor indicates a significant issue that should be fixed.
	SeverityMajor Severity = "major"

	// SeverityMinor indicates a small issue that could be fixed.
	SeverityMinor Severity = "minor"

	// SeveritySuggestion indicates an optional improvement.
	SeveritySuggestion Severity = "suggestion"
)

// Effort estimates the work required to address an issue.
type Effort string

const (
	// EffortTrivial indicates minimal effort (< 15 minutes).
	EffortTrivial Effort = "trivial"

	// EffortLow indicates small effort (< 1 hour).
	EffortLow Effort = "low"

	// EffortMedium indicates moderate effort (1-4 hours).
	EffortMedium Effort = "medium"

	// EffortHigh indicates significant effort (> 4 hours).
	EffortHigh Effort = "high"
)
