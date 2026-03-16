//ff:type feature=history type=model
//ff:what 단일 파일의 전체 변경 이력
package history

import "time"

// FileHistory holds the complete change history for a single file.
type FileHistory struct {
	File    string        `json:"file" yaml:"file"`
	Created time.Time     `json:"created" yaml:"created"`
	History []ChangeEntry `json:"history" yaml:"history"`
}
