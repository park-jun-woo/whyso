//ff:type feature=session type=model
//ff:what 레코드 내 메시지 (role + content)
package model

import "encoding/json"

// Message represents the message field in a record.
type Message struct {
	Role    string          `json:"role"`
	Content json.RawMessage `json:"content"`
}
