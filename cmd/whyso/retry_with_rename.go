//ff:func feature=cli type=util control=sequence
//ff:what 이력이 비어있을 때 git rename 감지 후 이전 경로로 재검색

package main

import (
	"path/filepath"

	"github.com/clari/whyso/pkg/history"
)

func retryWithRename(histories map[string]*history.FileHistory, sessionsDir, projectRoot, absTarget string) map[string]*history.FileHistory {
	targetRel, err := filepath.Rel(projectRoot, absTarget)
	if err != nil {
		return histories
	}

	oldPath := history.DetectRename(projectRoot, targetRel)
	if oldPath == "" {
		return histories
	}

	oldFilter := func(relPath string) bool { return relPath == oldPath }
	oldHistories, err := history.BuildHistories(sessionsDir, projectRoot, oldFilter)
	if err != nil || len(oldHistories) == 0 {
		return histories
	}

	h, ok := oldHistories[oldPath]
	if !ok {
		return histories
	}

	h.File = targetRel
	h.MovedFrom = oldPath
	if histories == nil {
		histories = make(map[string]*history.FileHistory)
	}
	histories[targetRel] = h
	return histories
}
