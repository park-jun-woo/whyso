//ff:func feature=cli type=util control=sequence
//ff:what history 서브커맨드용 파일 필터 함수를 생성
package main

import (
	"os"
	"path/filepath"
	"strings"
)

func makeFilter(targetInfo os.FileInfo, all bool, projectRoot, absTarget string) func(string) bool {
	if !targetInfo.IsDir() {
		targetRel, err := filepath.Rel(projectRoot, absTarget)
		if err != nil {
			return func(string) bool { return false }
		}
		return func(relPath string) bool { return relPath == targetRel }
	}
	if !all {
		return func(string) bool { return false }
	}
	targetRel, err := filepath.Rel(projectRoot, absTarget)
	if err != nil {
		return func(string) bool { return false }
	}
	if targetRel == "." {
		return func(string) bool { return true }
	}
	return func(relPath string) bool {
		return strings.HasPrefix(relPath, targetRel+"/") || relPath == targetRel
	}
}
