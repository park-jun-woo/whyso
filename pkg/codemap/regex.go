package codemap

import (
	"regexp"
	"strings"
)

var (
	reGherkinScenario = regexp.MustCompile(`(?m)^\s*Scenario(?:\s+Outline)?:\s*(.+)$`)
	reMermaidState    = regexp.MustCompile(`(?m)^\s*(\w+)\s+-->`)
	reMermaidStateTo  = regexp.MustCompile(`(?m)-->\s+(\w+)`)
	reSQLName         = regexp.MustCompile(`(?m)^--\s*name:\s*(\w+)`)
	reRegoAllow       = regexp.MustCompile(`(?m)^(\w+)\s+if\s*\{`)
)

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

func parseSQL(src []byte) []string {
	matches := reSQLName.FindAllSubmatch(src, -1)
	var names []string
	for _, m := range matches {
		name := string(m[1])
		if name != "" {
			names = append(names, name)
		}
	}
	return names
}

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
