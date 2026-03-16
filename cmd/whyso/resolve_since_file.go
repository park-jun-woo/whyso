//ff:func feature=cli type=util control=sequence
//ff:what 단일 파일 타겟의 캐시 mtime으로 since를 결정
package main

import (
	"os"
	"path/filepath"
	"time"

	"github.com/clari/whyso/internal/output"
)

func resolveSinceFile(projectRoot, absTarget, outputDir, format string) time.Time {
	targetRel, _ := filepath.Rel(projectRoot, absTarget)
	cachedPath := output.OutputPath(outputDir, targetRel, format)
	info, err := os.Stat(cachedPath)
	if err != nil {
		return time.Time{}
	}
	return info.ModTime()
}
