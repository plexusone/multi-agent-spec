package multiagentspec

import (
	"crypto/rand"
	"encoding/hex"
	"time"
)

// MessageType represents the type of inter-agent message.
type MessageType string

const (
	// MsgDelegateWork assigns a task from lead to specialist.
	MsgDelegateWork MessageType = "delegate_work"

	// MsgAskQuestion requests information from another agent.
	MsgAskQuestion MessageType = "ask_question"

	// MsgShareFinding broadcasts a discovery to all agents.
	MsgShareFinding MessageType = "share_finding"

	// MsgRequestApproval asks lead to approve a plan.
	MsgRequestApproval MessageType = "request_approval"

	// MsgApproval confirms approval of a request.
	MsgApproval MessageType = "approval"

	// MsgRejection denies a request with reason.
	MsgRejection MessageType = "rejection"

	// MsgChallenge disputes another agent's finding.
	MsgChallenge MessageType = "challenge"

	// MsgVote casts a consensus vote.
	MsgVote MessageType = "vote"

	// MsgTaskClaimed indicates an agent claimed a task from queue.
	MsgTaskClaimed MessageType = "task_claimed"

	// MsgTaskCompleted indicates an agent completed a task.
	MsgTaskCompleted MessageType = "task_completed"

	// MsgShutdownRequest asks to terminate the team session.
	MsgShutdownRequest MessageType = "shutdown_request"

	// MsgShutdownApproved confirms team session termination.
	MsgShutdownApproved MessageType = "shutdown_approved"
)

// Message represents an inter-agent message in self-directed workflows.
type Message struct {
	// ID is a unique message identifier.
	ID string `json:"id"`

	// Type is the message type.
	Type MessageType `json:"type"`

	// From is the sender agent name.
	From string `json:"from"`

	// To is the recipient agent name, or "*" for broadcast.
	To string `json:"to,omitempty"`

	// Subject is an optional message subject line.
	Subject string `json:"subject,omitempty"`

	// Content is the message body.
	Content string `json:"content"`

	// Attachments are optional file or data attachments.
	Attachments []Attachment `json:"attachments,omitempty"`

	// Metadata holds additional key-value data.
	Metadata map[string]interface{} `json:"metadata,omitempty"`

	// Timestamp is when the message was created.
	Timestamp time.Time `json:"timestamp"`
}

// AttachmentType represents the type of message attachment.
type AttachmentType string

const (
	// AttachmentFile is a file path reference.
	AttachmentFile AttachmentType = "file"

	// AttachmentData is inline data.
	AttachmentData AttachmentType = "data"

	// AttachmentReference is a reference to external content.
	AttachmentReference AttachmentType = "reference"
)

// Attachment represents a message attachment.
type Attachment struct {
	// Name is the attachment identifier.
	Name string `json:"name"`

	// Type is the attachment type: file, data, or reference.
	Type AttachmentType `json:"type"`

	// Data is the attachment content (type-dependent).
	Data interface{} `json:"data,omitempty"`
}

// IsBroadcast returns true if this message is sent to all agents.
func (m *Message) IsBroadcast() bool {
	return m.To == "*" || m.To == ""
}

// NewMessage creates a new message with generated ID and timestamp.
func NewMessage(msgType MessageType, from, to, content string) *Message {
	return &Message{
		ID:        generateMessageID(),
		Type:      msgType,
		From:      from,
		To:        to,
		Content:   content,
		Timestamp: time.Now().UTC(),
	}
}

// NewBroadcast creates a new broadcast message to all agents.
func NewBroadcast(msgType MessageType, from, content string) *Message {
	return NewMessage(msgType, from, "*", content)
}

// WithSubject sets the message subject and returns the message for chaining.
func (m *Message) WithSubject(subject string) *Message {
	m.Subject = subject
	return m
}

// WithAttachment adds an attachment and returns the message for chaining.
func (m *Message) WithAttachment(name string, attachType AttachmentType, data interface{}) *Message {
	m.Attachments = append(m.Attachments, Attachment{
		Name: name,
		Type: attachType,
		Data: data,
	})
	return m
}

// WithMetadata sets a metadata key-value and returns the message for chaining.
func (m *Message) WithMetadata(key string, value interface{}) *Message {
	if m.Metadata == nil {
		m.Metadata = make(map[string]interface{})
	}
	m.Metadata[key] = value
	return m
}

// generateMessageID generates a random message ID.
func generateMessageID() string {
	b := make([]byte, 8)
	_, _ = rand.Read(b)
	return "msg-" + hex.EncodeToString(b)
}
