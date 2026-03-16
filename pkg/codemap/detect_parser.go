//ff:func feature=codemap type=util control=selection
//ff:what 파일 확장자와 경로로 언어와 파서를 결정

package codemap

import (
	"os"
	"strings"
)

func detectParser(ext, relPath string) (string, func([]byte) []string) {
	switch ext {
	case ".go":
		return "go", parseGo
	case ".ssac":
		return "ssac", parseSSaC
	case ".ts", ".tsx":
		return "typescript", parseTypeScript
	case ".js", ".jsx":
		return "javascript", parseJavaScript
	case ".py":
		return "python", parsePython
	case ".rs":
		return "rust", parseRust
	case ".rego":
		return "rego", parseRego
	case ".feature":
		return "gherkin", parseGherkin
	case ".sql":
		return "sql", parseSQL
	case ".html":
		return "stml", parseSTML
	case ".yaml", ".yml":
		if strings.Contains(relPath, "api/") || strings.Contains(relPath, "openapi") {
			return "openapi", parseOpenAPI
		}
		return "", nil
	case ".md":
		if strings.HasPrefix(relPath, "states/") || strings.HasPrefix(relPath, "states"+string(os.PathSeparator)) {
			return "mermaid", parseMermaid
		}
		return "", nil
	}
	return "", nil
}
