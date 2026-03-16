//ff:func feature=session type=model control=sequence
//ff:what Record가 실제 사용자 요청인지 판별 (tool_result 컨테이너 제외)
package model

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
