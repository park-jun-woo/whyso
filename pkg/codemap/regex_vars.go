package codemap

import "regexp"

var (
	reGherkinScenario = regexp.MustCompile(`(?m)^\s*Scenario(?:\s+Outline)?:\s*(.+)$`)
	reMermaidState    = regexp.MustCompile(`(?m)^\s*(\w+)\s+-->`)
	reMermaidStateTo  = regexp.MustCompile(`(?m)-->\s+(\w+)`)
	reSQLName         = regexp.MustCompile(`(?m)^--\s*name:\s*(\w+)`)
	reRegoAllow       = regexp.MustCompile(`(?m)^(\w+)\s+if\s*\{`)
)
