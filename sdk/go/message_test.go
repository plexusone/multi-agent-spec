package multiagentspec

import (
	"encoding/json"
	"testing"
)

func TestMessageTypeConstants(t *testing.T) {
	tests := []struct {
		mt   MessageType
		want string
	}{
		{MsgDelegateWork, "delegate_work"},
		{MsgAskQuestion, "ask_question"},
		{MsgShareFinding, "share_finding"},
		{MsgRequestApproval, "request_approval"},
		{MsgApproval, "approval"},
		{MsgRejection, "rejection"},
		{MsgChallenge, "challenge"},
		{MsgVote, "vote"},
		{MsgTaskClaimed, "task_claimed"},
		{MsgTaskCompleted, "task_completed"},
		{MsgShutdownRequest, "shutdown_request"},
		{MsgShutdownApproved, "shutdown_approved"},
	}

	for _, tt := range tests {
		if string(tt.mt) != tt.want {
			t.Errorf("MessageType %v = %q, want %q", tt.mt, string(tt.mt), tt.want)
		}
	}
}

func TestNewMessage(t *testing.T) {
	msg := NewMessage(MsgDelegateWork, "lead", "specialist", "Please review this code")

	if msg.ID == "" {
		t.Error("ID should not be empty")
	}
	if msg.Type != MsgDelegateWork {
		t.Errorf("Type = %q, want %q", msg.Type, MsgDelegateWork)
	}
	if msg.From != "lead" {
		t.Errorf("From = %q, want %q", msg.From, "lead")
	}
	if msg.To != "specialist" {
		t.Errorf("To = %q, want %q", msg.To, "specialist")
	}
	if msg.Content != "Please review this code" {
		t.Errorf("Content = %q, want %q", msg.Content, "Please review this code")
	}
	if msg.Timestamp.IsZero() {
		t.Error("Timestamp should not be zero")
	}
}

func TestNewBroadcast(t *testing.T) {
	msg := NewBroadcast(MsgShareFinding, "researcher", "Found important data")

	if msg.To != "*" {
		t.Errorf("To = %q, want %q", msg.To, "*")
	}
	if !msg.IsBroadcast() {
		t.Error("IsBroadcast() should return true")
	}
}

func TestMessage_IsBroadcast(t *testing.T) {
	tests := []struct {
		to   string
		want bool
	}{
		{"*", true},
		{"", true},
		{"agent-a", false},
		{"specific-agent", false},
	}

	for _, tt := range tests {
		msg := &Message{To: tt.to}
		if got := msg.IsBroadcast(); got != tt.want {
			t.Errorf("IsBroadcast() with To=%q = %v, want %v", tt.to, got, tt.want)
		}
	}
}

func TestMessage_Chaining(t *testing.T) {
	msg := NewMessage(MsgDelegateWork, "lead", "specialist", "Task details").
		WithSubject("Code Review").
		WithAttachment("diff.txt", AttachmentFile, "/path/to/diff.txt").
		WithMetadata("priority", "high")

	if msg.Subject != "Code Review" {
		t.Errorf("Subject = %q, want %q", msg.Subject, "Code Review")
	}
	if len(msg.Attachments) != 1 {
		t.Errorf("len(Attachments) = %d, want 1", len(msg.Attachments))
	}
	if msg.Attachments[0].Name != "diff.txt" {
		t.Errorf("Attachment Name = %q, want %q", msg.Attachments[0].Name, "diff.txt")
	}
	if msg.Metadata["priority"] != "high" {
		t.Errorf("Metadata[priority] = %v, want %q", msg.Metadata["priority"], "high")
	}
}

func TestMessageSerialization(t *testing.T) {
	msg := NewMessage(MsgChallenge, "reviewer-a", "reviewer-b", "I disagree with this finding").
		WithSubject("Security Assessment").
		WithMetadata("severity", "high")

	data, err := json.Marshal(msg)
	if err != nil {
		t.Fatalf("json.Marshal failed: %v", err)
	}

	var decoded Message
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("json.Unmarshal failed: %v", err)
	}

	if decoded.Type != MsgChallenge {
		t.Errorf("Type = %q, want %q", decoded.Type, MsgChallenge)
	}
	if decoded.From != "reviewer-a" {
		t.Errorf("From = %q, want %q", decoded.From, "reviewer-a")
	}
	if decoded.To != "reviewer-b" {
		t.Errorf("To = %q, want %q", decoded.To, "reviewer-b")
	}
	if decoded.Subject != "Security Assessment" {
		t.Errorf("Subject = %q, want %q", decoded.Subject, "Security Assessment")
	}
}

func TestAttachmentTypeConstants(t *testing.T) {
	tests := []struct {
		at   AttachmentType
		want string
	}{
		{AttachmentFile, "file"},
		{AttachmentData, "data"},
		{AttachmentReference, "reference"},
	}

	for _, tt := range tests {
		if string(tt.at) != tt.want {
			t.Errorf("AttachmentType %v = %q, want %q", tt.at, string(tt.at), tt.want)
		}
	}
}
