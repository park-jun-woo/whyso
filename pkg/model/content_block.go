//ff:type feature=session type=model
//ff:what 어시스턴트 메시지 content 배열의 단일 블록 (text/tool_use)
package model

import "encoding/json"

// ContentBlock represents a single block in an assistant message content array.
type ContentBlock struct {
	Type  string          `json:"type"`
	Text  string          `json:"text,omitempty"`
	ID    string          `json:"id,omitempty"`
	Name  string          `json:"name,omitempty"`
	Input json.RawMessage `json:"input,omitempty"`
}
