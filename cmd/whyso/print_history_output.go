//ff:func feature=cli type=formatter control=sequence
//ff:what history 결과를 stdout에 출력 (새 결과 또는 캐시에서 읽기)
package main

import (
	"path/filepath"

	"github.com/clari/whyso/internal/output"
	"github.com/clari/whyso/pkg/history"
)

func printHistoryOutput(histories map[string]*history.FileHistory, format, outputDir, projectRoot, absTarget string) {
	if len(histories) > 0 {
		printHistories(histories, format)
		return
	}
	targetRel, _ := filepath.Rel(projectRoot, absTarget)
	cachedPath := output.OutputPath(outputDir, targetRel, format)
	cached, err := output.ReadYAML(cachedPath)
	if err != nil {
		return
	}
	printSingleHistory(cached, format)
}
