//ff:func feature=codemap type=util control=sequence
//ff:what 소스 파일이 맵 파일보다 새로운지 확인

package codemap

import (
	"os"
	"path/filepath"
)

// NeedsUpdate returns true if any source file in root is newer than mapFile.
// Returns true if mapFile does not exist.
func NeedsUpdate(root, mapFile string) bool {
	info, err := os.Stat(mapFile)
	if err != nil {
		return true
	}
	mapMtime := info.ModTime()

	needsUpdate := false
	filepath.Walk(root, func(path string, fi os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		if fi.IsDir() && shouldSkipDir(fi.Name()) {
			return filepath.SkipDir
		}
		if fi.IsDir() {
			return nil
		}
		ext := filepath.Ext(path)
		rel, _ := filepath.Rel(root, path)
		lang, _ := detectParser(ext, rel)
		if lang == "" {
			return nil
		}
		if fi.ModTime().After(mapMtime) {
			needsUpdate = true
			return filepath.SkipAll
		}
		return nil
	})
	return needsUpdate
}
