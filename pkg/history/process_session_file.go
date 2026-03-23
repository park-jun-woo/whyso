//ff:func feature=history type=builder control=sequence
//ff:what 단일 세션 파일의 레코드에서 이력과 힌트를 처리

package history

import (
	"path/filepath"
	"strings"

	"github.com/clari/whyso/pkg/model"
	"github.com/clari/whyso/pkg/parser"
)

func processSessionFile(path string, records []model.Record, histories map[string]*FileHistory, projectRoot string, filter func(string) bool) {
	changes := parser.ExtractChanges(records)
	idx := BuildIndex(records)
	sessionID := strings.TrimSuffix(filepath.Base(path), ".jsonl")

	processChanges(changes, idx, histories, projectRoot, sessionID, records, filter)

	hints := parser.ExtractBashHints(records)
	processHints(hints, histories, projectRoot)
}
