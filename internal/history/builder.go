package history

import (
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/clari/whyso/internal/parser"
)

// FileHistory holds the complete change history for a single file.
type FileHistory struct {
	File    string        `json:"file" yaml:"file"`
	Created time.Time     `json:"created" yaml:"created"`
	History []ChangeEntry `json:"history" yaml:"history"`
}

// Source points to a specific line in a JSONL file.
type Source struct {
	File string `json:"file" yaml:"file"`
	Line int    `json:"line" yaml:"line"`
}

// ChangeEntry represents one change event (possibly grouped from multiple tool_uses).
type ChangeEntry struct {
	Timestamp   time.Time `json:"timestamp" yaml:"timestamp"`
	Session     string    `json:"session" yaml:"session"`
	UserRequest string    `json:"user_request" yaml:"user_request"`
	Answer      string    `json:"answer,omitempty" yaml:"answer,omitempty"`
	Tool        string    `json:"tool" yaml:"tool"`
	Subagent    bool      `json:"subagent,omitempty" yaml:"subagent,omitempty"`
	Sources     []Source  `json:"sources" yaml:"sources"`
}

// BuildHistories builds file histories from all sessions in the directory.
func BuildHistories(sessionsDir, projectRoot string, filter func(string) bool) (map[string]*FileHistory, error) {
	return buildHistories(sessionsDir, projectRoot, time.Time{}, filter)
}

// BuildHistoriesIncremental builds file histories only from sessions modified after since.
func BuildHistoriesIncremental(sessionsDir, projectRoot string, since time.Time, filter func(string) bool) (map[string]*FileHistory, error) {
	return buildHistories(sessionsDir, projectRoot, since, filter)
}

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

func lastEntry(h *FileHistory) *ChangeEntry {
	if len(h.History) == 0 {
		return nil
	}
	return &h.History[len(h.History)-1]
}

