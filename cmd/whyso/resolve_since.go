//ff:func feature=cli type=util control=sequence
//ff:what 타겟 경로와 캐시에서 증분 갱신 기준 시각(since)을 결정
package main

import (
	"os"
	"time"
)

func resolveSince(targetInfo os.FileInfo, projectRoot, absTarget, outputDir, format string) time.Time {
	if !targetInfo.IsDir() {
		return resolveSinceFile(projectRoot, absTarget, outputDir, format)
	}
	return resolveSinceDir(projectRoot, absTarget, outputDir, format)
}
