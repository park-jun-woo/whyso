//ff:func feature=codemap type=parser control=iteration dimension=1
//ff:what Rego 파일에서 allow 규칙명을 추출 (정규식)

package codemap

func parseRegoRegex(src []byte) []string {
	matches := reRegoAllow.FindAllSubmatch(src, -1)
	var names []string
	for _, m := range matches {
		name := string(m[1])
		if name != "" {
			names = append(names, name)
		}
	}
	return names
}
