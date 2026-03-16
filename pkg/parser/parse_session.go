//ff:func feature=session type=parser control=iteration dimension=1
//ff:what JSONL 세션 파일과 subagent 파일을 파싱하여 전체 Record 슬라이스 반환
package parser

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/clari/whyso/pkg/model"
)

// ParseSession reads a JSONL file and its subagent files, returning all records.
func ParseSession(path string) ([]model.Record, error) {
	records, err := parseJSONL(path)
	if err != nil {
		return nil, err
	}

	// check for subagent files: <session-id>/subagents/*.jsonl
	sessionID := strings.TrimSuffix(filepath.Base(path), ".jsonl")
	subagentDir := filepath.Join(filepath.Dir(path), sessionID, "subagents")
	subEntries, err := os.ReadDir(subagentDir)
	if err != nil {
		return records, nil // no subagents directory, that's fine
	}

	for _, entry := range subEntries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".jsonl") {
			continue
		}
		subPath := filepath.Join(subagentDir, entry.Name())
		subRecords, err := parseJSONL(subPath)
		if err != nil {
			continue
		}
		records = append(records, subRecords...)
	}

	return records, nil
}
