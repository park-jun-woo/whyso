//ff:type feature=session type=model
//ff:what JSONL 한 줄에 대응하는 세션 레코드
package model

import "time"

// Record represents a single JSONL line from a Claude Code session.
type Record struct {
	Type       string    `json:"type"`
	UUID       string    `json:"uuid"`
	ParentUUID *string   `json:"parentUuid"`
	Timestamp  time.Time `json:"timestamp"`
	SessionID  string    `json:"sessionId"`
	Message    Message   `json:"message"`
	CWD        string    `json:"cwd"`
	UserType   string    `json:"userType"`
	SourceFile string    `json:"-"` // JSONL file path (set by parser)
	SourceLine int       `json:"-"` // line number in JSONL (set by parser)
}
