package multiagentspec

import (
	"encoding/json"
	"testing"
)

func TestCollaborationConfig_HasChannel(t *testing.T) {
	config := &CollaborationConfig{
		Channels: []Channel{
			{Name: "findings", Type: ChannelBroadcast},
			{Name: "direct", Type: ChannelDirect},
		},
	}

	if !config.HasChannel("findings") {
		t.Error("HasChannel should return true for existing channel")
	}
	if !config.HasChannel("direct") {
		t.Error("HasChannel should return true for existing channel")
	}
	if config.HasChannel("nonexistent") {
		t.Error("HasChannel should return false for nonexistent channel")
	}

	// Test nil config
	var nilConfig *CollaborationConfig
	if nilConfig.HasChannel("any") {
		t.Error("HasChannel on nil config should return false")
	}
}

func TestCollaborationConfig_GetChannel(t *testing.T) {
	config := &CollaborationConfig{
		Channels: []Channel{
			{Name: "findings", Type: ChannelBroadcast, Participants: []string{"*"}},
		},
	}

	ch := config.GetChannel("findings")
	if ch == nil {
		t.Fatal("GetChannel should return channel")
	}
	if ch.Type != ChannelBroadcast {
		t.Errorf("Type = %q, want %q", ch.Type, ChannelBroadcast)
	}

	if config.GetChannel("nonexistent") != nil {
		t.Error("GetChannel should return nil for nonexistent channel")
	}
}

func TestConsensusRules_Defaults(t *testing.T) {
	var nilRules *ConsensusRules
	if nilRules.EffectiveRequiredAgreement() != 0.5 {
		t.Errorf("EffectiveRequiredAgreement() = %f, want 0.5", nilRules.EffectiveRequiredAgreement())
	}
	if nilRules.EffectiveMaxRounds() != 3 {
		t.Errorf("EffectiveMaxRounds() = %d, want 3", nilRules.EffectiveMaxRounds())
	}

	emptyRules := &ConsensusRules{}
	if emptyRules.EffectiveRequiredAgreement() != 0.5 {
		t.Errorf("EffectiveRequiredAgreement() = %f, want 0.5", emptyRules.EffectiveRequiredAgreement())
	}
	if emptyRules.EffectiveMaxRounds() != 3 {
		t.Errorf("EffectiveMaxRounds() = %d, want 3", emptyRules.EffectiveMaxRounds())
	}

	customRules := &ConsensusRules{
		RequiredAgreement: 0.66,
		MaxRounds:         5,
	}
	if customRules.EffectiveRequiredAgreement() != 0.66 {
		t.Errorf("EffectiveRequiredAgreement() = %f, want 0.66", customRules.EffectiveRequiredAgreement())
	}
	if customRules.EffectiveMaxRounds() != 5 {
		t.Errorf("EffectiveMaxRounds() = %d, want 5", customRules.EffectiveMaxRounds())
	}
}

func TestCollaborationConfigSerialization(t *testing.T) {
	config := &CollaborationConfig{
		Lead:        "lead-agent",
		Specialists: []string{"security", "performance"},
		TaskQueue:   true,
		Consensus: &ConsensusRules{
			RequiredAgreement: 0.66,
			MaxRounds:         3,
			TieBreaker:        "lead",
		},
		Channels: []Channel{
			{Name: "findings", Type: ChannelBroadcast, Participants: []string{"*"}},
		},
	}

	data, err := json.Marshal(config)
	if err != nil {
		t.Fatalf("json.Marshal failed: %v", err)
	}

	var decoded CollaborationConfig
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("json.Unmarshal failed: %v", err)
	}

	if decoded.Lead != "lead-agent" {
		t.Errorf("Lead = %q, want %q", decoded.Lead, "lead-agent")
	}
	if len(decoded.Specialists) != 2 {
		t.Errorf("len(Specialists) = %d, want 2", len(decoded.Specialists))
	}
	if !decoded.TaskQueue {
		t.Error("TaskQueue should be true")
	}
	if decoded.Consensus == nil {
		t.Fatal("Consensus should not be nil")
	}
	if decoded.Consensus.RequiredAgreement != 0.66 {
		t.Errorf("RequiredAgreement = %f, want 0.66", decoded.Consensus.RequiredAgreement)
	}
	if len(decoded.Channels) != 1 {
		t.Errorf("len(Channels) = %d, want 1", len(decoded.Channels))
	}
}

func TestChannelTypeConstants(t *testing.T) {
	tests := []struct {
		ct   ChannelType
		want string
	}{
		{ChannelDirect, "direct"},
		{ChannelBroadcast, "broadcast"},
		{ChannelPubSub, "pub-sub"},
	}

	for _, tt := range tests {
		if string(tt.ct) != tt.want {
			t.Errorf("ChannelType %v = %q, want %q", tt.ct, string(tt.ct), tt.want)
		}
	}
}
