//ff:func feature=codemap type=util control=selection
//ff:what 언어와 경로에서 그룹명을 결정

package codemap

import "path/filepath"

func detectGroup(lang, relPath string, src []byte) string {
	switch lang {
	case "go":
		return extractGoPackage(src)
	case "ssac":
		dir := filepath.Dir(relPath)
		return dir
	case "openapi":
		return "api"
	case "sql":
		return filepath.Dir(relPath)
	case "rego":
		return "policy"
	case "gherkin":
		return filepath.Dir(relPath)
	case "stml":
		return filepath.Dir(relPath)
	case "mermaid":
		return filepath.Dir(relPath)
	default:
		return filepath.Dir(relPath)
	}
}
