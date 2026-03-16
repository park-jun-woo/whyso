//ff:type feature=change type=model
//ff:what Edit tool_use의 input 필드 구조
package parser

type editInput struct {
	FilePath  string `json:"file_path"`
	OldString string `json:"old_string"`
	NewString string `json:"new_string"`
}
