package codemap

import (
	"context"

	sitter "github.com/smacker/go-tree-sitter"
	"github.com/smacker/go-tree-sitter/golang"
	"github.com/smacker/go-tree-sitter/html"
	"github.com/smacker/go-tree-sitter/javascript"
	"github.com/smacker/go-tree-sitter/python"
	"github.com/smacker/go-tree-sitter/rust"
	ts "github.com/smacker/go-tree-sitter/typescript/typescript"
	"github.com/smacker/go-tree-sitter/yaml"
)

func runQuery(lang *sitter.Language, src []byte, query string) []string {
	parser := sitter.NewParser()
	parser.SetLanguage(lang)

	tree, err := parser.ParseCtx(context.Background(), nil, src)
	if err != nil {
		return nil
	}
	root := tree.RootNode()

	q, err := sitter.NewQuery([]byte(query), lang)
	if err != nil {
		return nil
	}

	cursor := sitter.NewQueryCursor()
	cursor.Exec(q, root)

	var names []string
	for {
		match, ok := cursor.NextMatch()
		if !ok {
			break
		}
		for _, capture := range match.Captures {
			name := capture.Node.Content(src)
			if name != "" {
				names = append(names, name)
			}
		}
	}
	return names
}

func parseGo(src []byte) []string {
	query := `
		(function_declaration name: (identifier) @name)
		(method_declaration name: (field_identifier) @name)
	`
	return runQuery(golang.GetLanguage(), src, query)
}

func parseSSaC(src []byte) []string {
	// SSaC uses Go syntax with .ssac extension
	query := `(function_declaration name: (identifier) @name)`
	return runQuery(golang.GetLanguage(), src, query)
}

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

func parsePython(src []byte) []string {
	query := `(function_definition name: (identifier) @name)`
	return runQuery(python.GetLanguage(), src, query)
}

func parseRust(src []byte) []string {
	query := `(function_item name: (identifier) @name)`
	return runQuery(rust.GetLanguage(), src, query)
}

func parseRego(src []byte) []string {
	// Rego: tree-sitter grammar not bundled, use regex fallback
	return parseRegoRegex(src)
}

func parseOpenAPI(src []byte) []string {
	// Extract operationId values from YAML
	return parseOpenAPIFromYAML(yaml.GetLanguage(), src)
}

func parseSTML(src []byte) []string {
	return parseSTMLFromHTML(html.GetLanguage(), src)
}

func parseOpenAPIFromYAML(lang *sitter.Language, src []byte) []string {
	parser := sitter.NewParser()
	parser.SetLanguage(lang)

	tree, err := parser.ParseCtx(context.Background(), nil, src)
	if err != nil {
		return nil
	}
	root := tree.RootNode()

	var names []string
	findOperationIDs(root, src, &names)
	return names
}

func findOperationIDs(node *sitter.Node, src []byte, names *[]string) {
	for i := 0; i < int(node.ChildCount()); i++ {
		child := node.Child(i)
		if child.Type() == "block_mapping_pair" {
			key := child.ChildByFieldName("key")
			value := child.ChildByFieldName("value")
			if key != nil && value != nil && key.Content(src) == "operationId" {
				val := value.Content(src)
				// strip quotes if present
				val = trimQuotes(val)
				if val != "" {
					*names = append(*names, val)
				}
			}
		}
		findOperationIDs(child, src, names)
	}
}

func parseSTMLFromHTML(lang *sitter.Language, src []byte) []string {
	parser := sitter.NewParser()
	parser.SetLanguage(lang)

	tree, err := parser.ParseCtx(context.Background(), nil, src)
	if err != nil {
		return nil
	}
	root := tree.RootNode()

	var names []string
	findDataAttributes(root, src, &names)
	return names
}

func findDataAttributes(node *sitter.Node, src []byte, names *[]string) {
	if node.Type() == "attribute" {
		var attrName, attrVal string
		for i := 0; i < int(node.ChildCount()); i++ {
			child := node.Child(i)
			switch child.Type() {
			case "attribute_name":
				attrName = child.Content(src)
			case "quoted_attribute_value", "attribute_value":
				attrVal = trimQuotes(child.Content(src))
			}
		}
		if (attrName == "data-fetch" || attrName == "data-action") && attrVal != "" {
			*names = append(*names, attrName+":"+attrVal)
		}
	}
	for i := 0; i < int(node.ChildCount()); i++ {
		findDataAttributes(node.Child(i), src, names)
	}
}

func trimQuotes(s string) string {
	if len(s) >= 2 {
		if (s[0] == '"' && s[len(s)-1] == '"') || (s[0] == '\'' && s[len(s)-1] == '\'') {
			return s[1 : len(s)-1]
		}
	}
	return s
}
