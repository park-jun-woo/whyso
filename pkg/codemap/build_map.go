//ff:func feature=codemap type=builder control=sequence
//ff:what 디렉토리를 스캔하여 지원 파일에서 키워드를 추출

package codemap

import (
	"os"
	"path/filepath"
)

// BuildMap scans the directory and extracts keywords from all supported files.
func BuildMap(root string) ([]Section, error) {
	sections := map[string]map[string][]string{} // language -> group -> keywords

	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		if info.IsDir() && shouldSkipDir(info.Name()) {
			return filepath.SkipDir
		}
		if info.IsDir() {
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
