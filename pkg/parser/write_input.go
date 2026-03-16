//ff:type feature=change type=model
//ff:what Write tool_use의 input 필드 구조
package parser

type writeInput struct {
	FilePath string `json:"file_path"`
	Content  string `json:"content"`
}
