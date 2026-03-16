//ff:func feature=codemap type=parser control=sequence
//ff:what Python 소스에서 함수명을 추출

package codemap

import "github.com/smacker/go-tree-sitter/python"

func parsePython(src []byte) []string {
	query := `(function_definition name: (identifier) @name)`
	return runQuery(python.GetLanguage(), src, query)
}
