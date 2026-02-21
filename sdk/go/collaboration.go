package multiagentspec

// CollaborationConfig defines how agents interact in self-directed workflows.
type CollaborationConfig struct {
	// Lead is the lead agent name (required for crew workflow).
	Lead string `json:"lead,omitempty"`

	// Specialists are non-delegating specialist agent names.
	Specialists []string `json:"specialists,omitempty"`

	// TaskQueue enables shared task queue for self-claiming (swarm workflow).
	TaskQueue bool `json:"task_queue,omitempty"`

	// Consensus defines consensus rules (council workflow).
	Consensus *ConsensusRules `json:"consensus,omitempty"`

	// Channels define communication pathways between agents.
	Channels []Channel `json:"channels,omitempty"`
}

// ConsensusRules defines how agents reach agreement in council workflows.
type ConsensusRules struct {
	// RequiredAgreement is the fraction of agents that must agree (0.0-1.0).
	// Default is 0.5 (simple majority).
	RequiredAgreement float64 `json:"required_agreement,omitempty"`

	// MaxRounds is the maximum debate rounds before forcing decision.
	// Default is 3.
	MaxRounds int `json:"max_rounds,omitempty"`

	// TieBreaker is the agent name to break ties, or "lead" for lead agent.
	TieBreaker string `json:"tie_breaker,omitempty"`
}

// ChannelType represents the communication pattern of a channel.
type ChannelType string

const (
	// ChannelDirect is point-to-point communication between two agents.
	ChannelDirect ChannelType = "direct"

	// ChannelBroadcast sends messages to all participants.
	ChannelBroadcast ChannelType = "broadcast"

	// ChannelPubSub allows agents to subscribe to topics.
	ChannelPubSub ChannelType = "pub-sub"
)

// Channel defines a communication pathway between agents.
type Channel struct {
	// Name is the channel identifier.
	Name string `json:"name"`

	// Type is the channel type: direct, broadcast, or pub-sub.
	Type ChannelType `json:"type"`

	// Participants are agent names. Use "*" for all agents.
	Participants []string `json:"participants,omitempty"`
}

// HasChannel returns true if a channel with the given name exists.
func (c *CollaborationConfig) HasChannel(name string) bool {
	if c == nil {
		return false
	}
	for _, ch := range c.Channels {
		if ch.Name == name {
			return true
		}
	}
	return false
}

// GetChannel returns the channel with the given name, or nil if not found.
func (c *CollaborationConfig) GetChannel(name string) *Channel {
	if c == nil {
		return nil
	}
	for i := range c.Channels {
		if c.Channels[i].Name == name {
			return &c.Channels[i]
		}
	}
	return nil
}

// EffectiveRequiredAgreement returns the required agreement, defaulting to 0.5.
func (r *ConsensusRules) EffectiveRequiredAgreement() float64 {
	if r == nil || r.RequiredAgreement == 0 {
		return 0.5
	}
	return r.RequiredAgreement
}

// EffectiveMaxRounds returns the max rounds, defaulting to 3.
func (r *ConsensusRules) EffectiveMaxRounds() int {
	if r == nil || r.MaxRounds == 0 {
		return 3
	}
	return r.MaxRounds
}
