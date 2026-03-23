//ff:func feature=history type=parser control=iteration dimension=1
//ff:what git log --summary 출력에서 rename 이전 경로를 파싱

package history

import (
	"path/filepath"
	"regexp"
	"strings"
)

var renamePattern = regexp.MustCompile(`rename (.*)\{([^}]+) => ([^}]+)\}(.*)`)
var renameSimplePattern = regexp.MustCompile(`rename (.+) => (.+) \(`)

func parseRenameOutput(out string) string {
	for _, line := range strings.Split(out, "\n") {
		line = strings.TrimSpace(line)
		if !strings.HasPrefix(line, "rename ") {
			continue
		}
		if m := renamePattern.FindStringSubmatch(line); m != nil {
			suffix := strings.Split(m[4], " (")[0]
			return filepath.Join(m[1], m[2], suffix)
		}
		if m := renameSimplePattern.FindStringSubmatch(line); m != nil {
			return m[1]
		}
	}
	return ""
}
