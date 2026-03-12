package codemap

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// Section represents a language section with grouped keywords.
type Section struct {
	Language string
	Groups   map[string][]string // group name -> keywords
}

// BuildMap scans the directory and extracts keywords from all supported files.
func BuildMap(root string) ([]Section, error) {
	sections := map[string]map[string][]string{} // language -> group -> keywords

	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		if info.IsDir() {
			base := info.Name()
			if shouldSkipDir(base) {
				return filepath.SkipDir
			}
			return nil
		}

		rel, err := filepath.Rel(root, path)
		if err != nil {
			return nil
		}

		ext := filepath.Ext(path)
		lang, parser := detectParser(ext, rel)
		if parser == nil {
			return nil
		}

		src, err := os.ReadFile(path)
		if err != nil {
			return nil
		}

		group := detectGroup(lang, rel, src)
		keywords := parser(src)
		if len(keywords) == 0 {
			return nil
		}

		if sections[lang] == nil {
			sections[lang] = map[string][]string{}
		}
		sections[lang][group] = append(sections[lang][group], keywords...)
		return nil
	})
	if err != nil {
		return nil, err
	}

	return buildSections(sections), nil
}

// FormatMap writes the map to w in the standard format.
func FormatMap(w io.Writer, sections []Section) {
	for i, sec := range sections {
		if i > 0 {
			fmt.Fprintln(w)
		}
		fmt.Fprintf(w, "## %s\n", sec.Language)

		groups := sortedKeys(sec.Groups)
		for _, g := range groups {
			keywords := sec.Groups[g]
			sort.Strings(keywords)
			keywords = dedupe(keywords)
			fmt.Fprintf(w, "[%s]%s\n", g, strings.Join(keywords, ","))
		}
	}
}

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
	for lang, groups := range data {
		found := false
		for _, o := range order {
			if lang == o {
				found = true
				break
			}
		}
		if !found {
			result = append(result, Section{Language: lang, Groups: groups})
		}
	}
	return result
}

func shouldSkipDir(name string) bool {
	skip := []string{".git", ".whyso", ".claude", "node_modules", "vendor", "__pycache__", ".venv", "target", "dist", "build", ".idea", ".vscode"}
	for _, s := range skip {
		if name == s {
			return true
		}
	}
	return false
}

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

func extractGoPackage(src []byte) string {
	for _, line := range strings.Split(string(src), "\n") {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "package ") {
			parts := strings.Fields(line)
			if len(parts) >= 2 {
				return parts[1]
			}
		}
	}
	return "unknown"
}

func sortedKeys(m map[string][]string) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

func dedupe(ss []string) []string {
	if len(ss) == 0 {
		return ss
	}
	result := []string{ss[0]}
	for i := 1; i < len(ss); i++ {
		if ss[i] != ss[i-1] {
			result = append(result, ss[i])
		}
	}
	return result
}
