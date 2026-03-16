//ff:func feature=session type=util control=sequence
//ff:what 작업 디렉토리에서 Claude Code 세션 디렉토리 경로를 탐지
package parser

import (
	"fmt"
	"os"
	"path/filepath"
)

// DetectSessionsDir returns the Claude Code sessions directory for the given working directory.
func DetectSessionsDir(cwd string) (string, error) {
	abs, err := filepath.Abs(cwd)
	if err != nil {
		return "", err
	}

	slug := toSlug(abs)
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	dir := filepath.Join(home, ".claude", "projects", slug)
	info, err := os.Stat(dir)
	if err != nil {
		return "", fmt.Errorf("sessions directory not found: %s", dir)
	}
	if !info.IsDir() {
		return "", fmt.Errorf("not a directory: %s", dir)
	}
	return dir, nil
}
