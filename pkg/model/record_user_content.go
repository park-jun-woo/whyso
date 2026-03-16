//ff:func feature=session type=model control=sequence
//ff:what Record에서 사용자 메시지 텍스트를 추출
package model

import "encoding/json"

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
