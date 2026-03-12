package output

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/clari/whylog/internal/history"
)

// ReadYAML reads an existing YAML history file back into a FileHistory.
func ReadYAML(path string) (*history.FileHistory, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	h := &history.FileHistory{}
	var current *history.ChangeEntry
	inSources := false

	scanner := bufio.NewScanner(f)
	scanner.Buffer(make([]byte, 0, 64*1024), 10*1024*1024)

	for scanner.Scan() {
		line := scanner.Text()

		if strings.HasPrefix(line, "file: ") {
			h.File = strings.TrimPrefix(line, "file: ")
			continue
		}
		if strings.HasPrefix(line, "created: ") {
			t, err := time.Parse(time.RFC3339, strings.TrimPrefix(line, "created: "))
			if err == nil {
				h.Created = t
			}
			continue
		}
		if line == "history:" {
			continue
		}
		if strings.HasPrefix(line, "  - timestamp: ") {
			if current != nil {
				h.History = append(h.History, *current)
			}
			current = &history.ChangeEntry{}
			inSources = false
			t, err := time.Parse(time.RFC3339, strings.TrimPrefix(line, "  - timestamp: "))
			if err == nil {
				current.Timestamp = t
			}
			continue
		}

		if current == nil {
			continue
		}

		trimmed := strings.TrimSpace(line)

		if strings.HasPrefix(trimmed, "session: ") {
			current.Session = strings.TrimPrefix(trimmed, "session: ")
			inSources = false
			continue
		}
		if strings.HasPrefix(trimmed, "user_request: ") {
			current.UserRequest = unquote(strings.TrimPrefix(trimmed, "user_request: "))
			inSources = false
			continue
		}
		if strings.HasPrefix(trimmed, "answer: ") {
			current.Answer = unquote(strings.TrimPrefix(trimmed, "answer: "))
			inSources = false
			continue
		}
		if strings.HasPrefix(trimmed, "tool: ") {
			current.Tool = strings.TrimPrefix(trimmed, "tool: ")
			inSources = false
			continue
		}
		if strings.HasPrefix(trimmed, "subagent: ") {
			current.Subagent = strings.TrimPrefix(trimmed, "subagent: ") == "true"
			inSources = false
			continue
		}
		if strings.HasPrefix(trimmed, "source: ") {
			src := parseSource(strings.TrimPrefix(trimmed, "source: "))
			current.Sources = append(current.Sources, src)
			inSources = false
			continue
		}
		if trimmed == "sources:" {
			inSources = true
			continue
		}
		if inSources && strings.HasPrefix(trimmed, "- ") {
			src := parseSource(strings.TrimPrefix(trimmed, "- "))
			current.Sources = append(current.Sources, src)
			continue
		}
	}

	if current != nil {
		h.History = append(h.History, *current)
	}

	return h, scanner.Err()
}

func parseSource(s string) history.Source {
	idx := strings.LastIndex(s, ":")
	if idx < 0 {
		return history.Source{File: s}
	}
	line, err := strconv.Atoi(s[idx+1:])
	if err != nil {
		return history.Source{File: s}
	}
	return history.Source{File: s[:idx], Line: line}
}

func unquote(s string) string {
	if len(s) >= 2 && s[0] == '"' && s[len(s)-1] == '"' {
		u, err := strconv.Unquote(s)
		if err == nil {
			return u
		}
	}
	return s
}

// OutputPath returns the output file path for a given relative path.
func OutputPath(outputDir, relPath, format string) string {
	return fmt.Sprintf("%s/%s.%s", outputDir, relPath, format)
}
