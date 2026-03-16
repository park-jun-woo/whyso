//ff:func feature=cli type=util control=sequence
//ff:what 디렉토리 타겟의 캐시 mtime으로 since를 결정
package main

import (
	"path/filepath"
	"time"
)

func resolveSinceDir(projectRoot, absTarget, outputDir, format string) time.Time {
	targetRel, _ := filepath.Rel(projectRoot, absTarget)
	if targetRel == "." {
		return oldestOutputMtime(outputDir, format)
	}
	return oldestOutputMtime(filepath.Join(outputDir, targetRel), format)
}
