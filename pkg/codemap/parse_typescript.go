//ff:func feature=codemap type=parser control=sequence
//ff:what TypeScript 소스에서 함수명을 추출

package codemap

import ts "github.com/smacker/go-tree-sitter/typescript/typescript"

func parseTypeScript(src []byte) []string {
	query := `
		(function_declaration name: (identifier) @name)
		(method_definition name: (property_identifier) @name)
		(lexical_declaration
			(variable_declarator
				name: (identifier) @name
				value: (arrow_function)))
	`
	return runQuery(ts.GetLanguage(), src, query)
}
