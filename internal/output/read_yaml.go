//ff:func feature=output type=reader control=iteration dimension=1
//ff:what 기존 YAML 이력 파일을 FileHistory로 읽기

package output

import (
	"bufio"
	"os"
	"strings"
	"time"

	"github.com/clari/whyso/pkg/history"
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
			h.Created, _ = time.Parse(time.RFC3339, strings.TrimPrefix(line, "created: "))
			continue
		}
		if line == "history:" {
			continue
		}

		isNewEntry := strings.HasPrefix(line, "  - timestamp: ")
		if isNewEntry && current != nil {
			h.History = append(h.History, *current)
		}
		if isNewEntry {
			current = &history.ChangeEntry{}
			inSources = false
			current.Timestamp, _ = time.Parse(time.RFC3339, strings.TrimPrefix(line, "  - timestamp: "))
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
