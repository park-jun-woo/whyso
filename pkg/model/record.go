package model

import (
	"encoding/json"
	"time"
)

// Record represents a single JSONL line from a Claude Code session.
type Record struct {
	Type       string          `json:"type"`
	UUID       string          `json:"uuid"`
	ParentUUID *string         `json:"parentUuid"`
	Timestamp  time.Time       `json:"timestamp"`
	SessionID  string          `json:"sessionId"`
	Message    Message         `json:"message"`
	CWD        string          `json:"cwd"`
	UserType   string          `json:"userType"`
	SourceFile string          `json:"-"` // JSONL file path (set by parser)
	SourceLine int             `json:"-"` // line number in JSONL (set by parser)
}

// Message represents the message field in a record.
type Message struct {
	Role    string          `json:"role"`
	Content json.RawMessage `json:"content"`
}

// ContentBlock represents a single block in an assistant message content array.
type ContentBlock struct {
	Type  string          `json:"type"`
	Text  string          `json:"text,omitempty"`
	ID    string          `json:"id,omitempty"`
	Name  string          `json:"name,omitempty"`
	Input json.RawMessage `json:"input,omitempty"`
}

// IsUserMessage returns true if this record is a real user request (not a tool_result container).
func (r *Record) IsUserMessage() bool {
	if r.Type != "user" {
		return false
	}
	if len(r.Message.Content) == 0 {
		return false
	}
	// string starts with '"', array starts with '['
	return r.Message.Content[0] == '"'
}

// UserContent returns the user message text. Returns empty string if not a user message.
func (r *Record) UserContent() string {
	if !r.IsUserMessage() {
		return ""
	}
	var s string
	if err := json.Unmarshal(r.Message.Content, &s); err != nil {
		return ""
	}
	return s
}

// ContentBlocks parses message.content as an array of ContentBlock.
func (r *Record) ContentBlocks() []ContentBlock {
	if len(r.Message.Content) == 0 || r.Message.Content[0] != '[' {
		return nil
	}
	var blocks []ContentBlock
	if err := json.Unmarshal(r.Message.Content, &blocks); err != nil {
		return nil
	}
	return blocks
}

// SessionInfo holds summary information about a session.
type SessionInfo struct {
	ID           string
	Timestamp    time.Time
	FirstMessage string
}
