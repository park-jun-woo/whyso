//ff:func feature=codemap type=parser control=sequence
//ff:what Rego 소스에서 규칙명을 추출

package codemap

func parseRego(src []byte) []string {
	// Rego: tree-sitter grammar not bundled, use regex fallback
	return parseRegoRegex(src)
}
