//ff:func feature=history type=builder control=iteration dimension=1
//ff:what 세션 파일을 순회하며 파일별 변경 이력을 구축하는 코어 로직
package history

import (
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/clari/whyso/pkg/parser"
)

func buildHistories(sessionsDir, projectRoot string, since time.Time, filter func(string) bool) (map[string]*FileHistory, error) {
	entries, err := os.ReadDir(sessionsDir)
	if err != nil {
		return nil, err
	}

	histories := make(map[string]*FileHistory)

	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".jsonl") {
			continue
		}

		if !since.IsZero() {
			info, err := entry.Info()
			if err != nil {
				continue
			}
			if !info.ModTime().After(since) {
				continue
			}
		}

		path := filepath.Join(sessionsDir, entry.Name())
		records, err := parser.ParseSession(path)
		if err != nil {
			continue
		}

		changes := parser.ExtractChanges(records)
		idx := BuildIndex(records)
		sessionID := strings.TrimSuffix(entry.Name(), ".jsonl")

		processChanges(changes, idx, histories, projectRoot, sessionID, records, filter)
	}

	for _, h := range histories {
		sort.Slice(h.History, func(i, j int) bool {
			return h.History[i].Timestamp.Before(h.History[j].Timestamp)
		})
		if len(h.History) > 0 {
			h.Created = h.History[0].Timestamp
		}
	}

	return histories, nil
}
