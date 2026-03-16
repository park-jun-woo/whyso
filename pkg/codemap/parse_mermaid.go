//ff:func feature=codemap type=parser control=iteration dimension=1
//ff:what Mermaid stateDiagram에서 상태명을 추출

package codemap

import "strings"

func parseMermaid(src []byte) []string {
	text := string(src)
	// Only parse within ```mermaid blocks
	if !strings.Contains(text, "stateDiagram") {
		return nil
	}

	seen := map[string]bool{}
	var states []string

	for _, m := range reMermaidState.FindAllStringSubmatch(text, -1) {
		s := m[1]
		if s != "[*]" && !seen[s] {
			seen[s] = true
			states = append(states, s)
		}
	}
	for _, m := range reMermaidStateTo.FindAllStringSubmatch(text, -1) {
		s := m[1]
		if s != "[*]" && !seen[s] {
			seen[s] = true
			states = append(states, s)
		}
	}
	return states
}
