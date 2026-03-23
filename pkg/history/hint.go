//ff:type feature=history type=model
//ff:what Bash cp/mv 명령에서 추정된 파일 이동 힌트

package history

import "time"

// Hint represents a cp/mv command detected in a Bash tool_use.
type Hint struct {
	Timestamp time.Time `json:"timestamp" yaml:"timestamp"`
	Session   string    `json:"session" yaml:"session"`
	Command   string    `json:"command" yaml:"command"`
	Source    Source    `json:"source" yaml:"source"`
}
