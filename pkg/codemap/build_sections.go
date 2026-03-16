//ff:func feature=codemap type=builder control=iteration dimension=1
//ff:what 언어→그룹→키워드 맵을 정렬된 Section 슬라이스로 변환

package codemap

func buildSections(data map[string]map[string][]string) []Section {
	order := []string{"go", "typescript", "javascript", "python", "rust", "ssac", "openapi", "sql", "rego", "gherkin", "stml", "mermaid"}
	var result []Section
	for _, lang := range order {
		groups, ok := data[lang]
		if !ok {
			continue
		}
		result = append(result, Section{Language: lang, Groups: groups})
	}
	// any remaining languages not in order
	orderSet := make(map[string]bool, len(order))
	for _, o := range order {
		orderSet[o] = true
	}
	for lang, groups := range data {
		if !orderSet[lang] {
			result = append(result, Section{Language: lang, Groups: groups})
		}
	}
	return result
}
