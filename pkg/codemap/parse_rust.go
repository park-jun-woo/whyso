//ff:func feature=codemap type=parser control=sequence
//ff:what Rust 소스에서 함수명을 추출

package codemap

import "github.com/smacker/go-tree-sitter/rust"

func parseRust(src []byte) []string {
	query := `(function_item name: (identifier) @name)`
	return runQuery(rust.GetLanguage(), src, query)
}
