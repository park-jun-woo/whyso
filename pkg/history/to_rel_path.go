//ff:func feature=history type=util control=sequence
//ff:what 절대 경로를 프로젝트 루트 기준 상대 경로로 변환

package history

import (
	"path/filepath"
	"strings"
)

func toRelPath(projectRoot, path string) string {
	if filepath.IsAbs(path) {
		rel, err := filepath.Rel(projectRoot, path)
		if err != nil || strings.HasPrefix(rel, "..") {
			return path
		}
		return rel
	}
	return path
}
