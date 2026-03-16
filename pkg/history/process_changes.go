//ff:func feature=history type=builder control=iteration dimension=1
//ff:what 단일 세션의 변경 목록을 순회하며 파일 이력에 추가
package history

import (
	"path/filepath"
	"strings"

	"github.com/clari/whyso/pkg/model"
	"github.com/clari/whyso/pkg/parser"
)

func processChanges(changes []parser.FileChange, idx RecordIndex, histories map[string]*FileHistory, projectRoot string, sessionID string, records []model.Record, filter func(string) bool) {
	for _, fc := range changes {
		relPath, err := filepath.Rel(projectRoot, fc.FilePath)
		if err != nil || strings.HasPrefix(relPath, "..") {
			continue
		}

		if filter != nil && !filter(relPath) {
			continue
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

		// group consecutive entries with same (session, user_request, answer, tool)
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
	}
}
