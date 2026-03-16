//ff:func feature=codemap type=util control=iteration dimension=1
//ff:what 디렉토리명이 스킵 대상인지 확인

package codemap

func shouldSkipDir(name string) bool {
	skip := []string{".git", ".whyso", ".claude", "node_modules", "vendor", "__pycache__", ".venv", "target", "dist", "build", ".idea", ".vscode"}
	for _, s := range skip {
		if name == s {
			return true
		}
	}
	return false
}
