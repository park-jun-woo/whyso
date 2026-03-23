//ff:type feature=change type=model
//ff:what Bash tool_use에서 감지된 cp/mv 명령 힌트

package parser

import "time"

// BashHint represents a cp/mv command detected in a Bash tool_use.
type BashHint struct {
	Command    string
	SrcPath    string
	DstPath    string
	Timestamp  time.Time
	SessionID  string
	SourceFile string
	SourceLine int
}
