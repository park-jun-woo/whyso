package parser

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/clari/whyso/internal/model"
)

const maxLineSize = 10 * 1024 * 1024 // 10MB

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

// parseJSONL reads a single JSONL file and returns records with source tracking.
func parseJSONL(path string) ([]model.Record, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var records []model.Record
	scanner := bufio.NewScanner(f)
	scanner.Buffer(make([]byte, 0, 64*1024), maxLineSize)

	lineNum := 0
	for scanner.Scan() {
		lineNum++
		line := scanner.Bytes()
		if len(line) == 0 {
			continue
		}
		var rec model.Record
		if err := json.Unmarshal(line, &rec); err != nil {
			continue // skip malformed lines
		}
		rec.SourceFile = path
		rec.SourceLine = lineNum
		records = append(records, rec)
	}
	if err := scanner.Err(); err != nil {
		return records, err
	}
	return records, nil
}

// DetectSessionsDir returns the Claude Code sessions directory for the given working directory.
func DetectSessionsDir(cwd string) (string, error) {
	abs, err := filepath.Abs(cwd)
	if err != nil {
		return "", err
	}

	slug := toSlug(abs)
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	dir := filepath.Join(home, ".claude", "projects", slug)
	info, err := os.Stat(dir)
	if err != nil {
		return "", fmt.Errorf("sessions directory not found: %s", dir)
	}
	if !info.IsDir() {
		return "", fmt.Errorf("not a directory: %s", dir)
	}
	return dir, nil
}

// toSlug converts an absolute path to a Claude Code project slug.
// e.g. /home/user/.clari/project → -home-user--clari-project
func toSlug(abs string) string {
	s := strings.TrimPrefix(abs, "/")
	s = strings.ReplaceAll(s, "/", "-")
	s = strings.ReplaceAll(s, ".", "-")
	return "-" + s
}

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

		for _, rec := range records {
			if rec.IsUserMessage() {
				info.Timestamp = rec.Timestamp
				content := rec.UserContent()
				if len(content) > 80 {
					content = content[:80] + "..."
				}
				info.FirstMessage = content
				break
			}
		}

		sessions = append(sessions, info)
	}

	sort.Slice(sessions, func(i, j int) bool {
		return sessions[i].Timestamp.Before(sessions[j].Timestamp)
	})

	return sessions, nil
}
