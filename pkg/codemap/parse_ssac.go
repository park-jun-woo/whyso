//ff:func feature=codemap type=parser control=sequence
//ff:what SSaC 소스에서 함수명을 추출

package codemap

import "github.com/smacker/go-tree-sitter/golang"

func parseSSaC(src []byte) []string {
	// SSaC uses Go syntax with .ssac extension
	query := `(function_declaration name: (identifier) @name)`
	return runQuery(golang.GetLanguage(), src, query)
}
