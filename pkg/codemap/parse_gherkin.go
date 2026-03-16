//ff:func feature=codemap type=parser control=iteration dimension=1
//ff:what Gherkin 파일에서 시나리오명을 추출

package codemap

import "strings"

func parseGherkin(src []byte) []string {
	matches := reGherkinScenario.FindAllSubmatch(src, -1)
	var names []string
	for _, m := range matches {
		name := strings.TrimSpace(string(m[1]))
		if name != "" {
			names = append(names, name)
		}
	}
	return names
}
