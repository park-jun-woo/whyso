//ff:type feature=history type=model
//ff:what JSONL 파일 내 특정 라인 위치
package history

// Source points to a specific line in a JSONL file.
type Source struct {
	File string `json:"file" yaml:"file"`
	Line int    `json:"line" yaml:"line"`
}
