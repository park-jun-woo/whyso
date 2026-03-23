//ff:func feature=history type=util control=sequence
//ff:what git log --follow로 파일의 이전 경로(rename)를 감지

package history

import (
	"os/exec"
)

// DetectRename uses git log --follow to find the previous path of a file.
// Returns empty string if no rename is found.
func DetectRename(projectRoot, relPath string) string {
	cmd := exec.Command("git", "log", "--follow", "--diff-filter=R", "--summary", "--format=", "--", relPath)
	cmd.Dir = projectRoot
	out, err := cmd.Output()
	if err != nil {
		return ""
	}
	return parseRenameOutput(string(out))
}
