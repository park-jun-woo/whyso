//ff:type feature=change type=model
//ff:what tool_use에서 추출된 단일 파일 변경 정보
package parser

import "time"

// FileChange represents a single file modification extracted from a tool_use.
type FileChange struct {
	FilePath    string
	Tool        string // "Write" | "Edit"
	ToolUseID   string
	OldString   string // Edit only
	NewString   string // Edit only
	Content     string // Write only
	RecordUUID  string
	RecordIndex int // index in records slice (for FindAnswer)
	Timestamp   time.Time
	SessionID   string
	SourceFile  string // JSONL file path
	SourceLine  int    // line number in JSONL
}
