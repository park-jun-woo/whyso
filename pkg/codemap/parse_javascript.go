//ff:func feature=codemap type=parser control=sequence
//ff:what JavaScript 소스에서 함수명을 추출

package codemap

import "github.com/smacker/go-tree-sitter/javascript"

func parseJavaScript(src []byte) []string {
	query := `
		(function_declaration name: (identifier) @name)
		(method_definition name: (property_identifier) @name)
		(lexical_declaration
			(variable_declarator
				name: (identifier) @name
				value: (arrow_function)))
	`
	return runQuery(javascript.GetLanguage(), src, query)
}
