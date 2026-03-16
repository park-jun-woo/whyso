//ff:type feature=history type=model
//ff:what 하나의 변경 이벤트 (여러 tool_use가 그룹화될 수 있음)
package history

import "time"

// ChangeEntry represents one change event (possibly grouped from multiple tool_uses).
type ChangeEntry struct {
	Timestamp   time.Time `json:"timestamp" yaml:"timestamp"`
	Session     string    `json:"session" yaml:"session"`
	UserRequest string    `json:"user_request" yaml:"user_request"`
	Answer      string    `json:"answer,omitempty" yaml:"answer,omitempty"`
	Tool        string    `json:"tool" yaml:"tool"`
	Subagent    bool      `json:"subagent,omitempty" yaml:"subagent,omitempty"`
	Sources     []Source  `json:"sources" yaml:"sources"`
}
