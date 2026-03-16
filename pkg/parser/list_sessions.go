//ff:func feature=session type=reader control=iteration dimension=1
//ff:what 디렉토리 내 모든 JSONL 세션의 요약 정보를 반환
package parser

import (
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/clari/whyso/pkg/model"
)

// ListSessions returns session info for all JSONL files in the given directory.
func ListSessions(dir string) ([]model.SessionInfo, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	var sessions []model.SessionInfo
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".jsonl") {
			continue
		}

		id := strings.TrimSuffix(entry.Name(), ".jsonl")
		path := filepath.Join(dir, entry.Name())

		info := model.SessionInfo{ID: id}
		records, err := ParseSession(path)
		if err != nil {
			continue
		}

		info.Timestamp, info.FirstMessage = findFirstUserMessage(records)

		sessions = append(sessions, info)
	}

	sort.Slice(sessions, func(i, j int) bool {
		return sessions[i].Timestamp.Before(sessions[j].Timestamp)
	})

	return sessions, nil
}
