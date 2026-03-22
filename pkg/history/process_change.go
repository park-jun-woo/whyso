//ff:func feature=history type=builder control=sequence
//ff:what 단일 파일 변경을 이력 맵에 추가 (중복 그룹핑 포함)

package history

import (
	"path/filepath"
	"strings"

	"github.com/clari/whyso/pkg/model"
	"github.com/clari/whyso/pkg/parser"
)

func processChange(fc parser.FileChange, idx RecordIndex, histories map[string]*FileHistory, projectRoot string, sessionID string, records []model.Record, filter func(string) bool) bool {
	relPath, err := filepath.Rel(projectRoot, fc.FilePath)
	if err != nil || strings.HasPrefix(relPath, "..") {
		return false
	}

	if filter != nil && !filter(relPath) {
		return false
	}

	userReq := FindUserRequest(idx, fc.RecordUUID)
	answer := FindAnswer(records, fc.RecordIndex)
	isSubagent := strings.Contains(fc.SourceFile, "/subagents/")

	src := Source{File: fc.SourceFile, Line: fc.SourceLine}

	h, ok := histories[relPath]
	if !ok {
		h = &FileHistory{File: relPath}
		histories[relPath] = h
	}

	if last := lastEntry(h); last != nil &&
		last.Session == sessionID &&
		last.UserRequest == userReq &&
		last.Answer == answer &&
		last.Tool == fc.Tool {
		last.Sources = append(last.Sources, src)
	} else {
		h.History = append(h.History, ChangeEntry{
			Timestamp:   fc.Timestamp,
			Session:     sessionID,
			UserRequest: userReq,
			Answer:      answer,
			Tool:        fc.Tool,
			Subagent:    isSubagent,
			Sources:     []Source{src},
		})
	}

	return true
}
