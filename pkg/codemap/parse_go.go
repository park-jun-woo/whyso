//ff:func feature=codemap type=parser control=sequence
//ff:what Go 소스에서 함수와 메서드명을 추출

package codemap

import "github.com/smacker/go-tree-sitter/golang"

func parseGo(src []byte) []string {
	query := `
		(function_declaration name: (identifier) @name)
		(method_declaration name: (field_identifier) @name)
	`
	return runQuery(golang.GetLanguage(), src, query)
}
