//ff:func feature=cli type=util control=sequence
//ff:what git repo와 Claude Code 세션 디렉토리 존재 여부를 검증

package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

func checkDeps() error {
	if err := exec.Command("git", "rev-parse", "--git-dir").Run(); err != nil {
		return fmt.Errorf("not a git repository. whyso requires a git project root")
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	claudeDir := filepath.Join(home, ".claude")
	if _, err := os.Stat(claudeDir); err != nil {
		return fmt.Errorf("~/.claude/ not found. install Claude Code first: https://claude.ai/download")
	}

	return nil
}
