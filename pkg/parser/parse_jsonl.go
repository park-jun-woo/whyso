//ff:func feature=session type=parser control=iteration dimension=1
//ff:what 단일 JSONL 파일을 파싱하여 Record 슬라이스 반환 (소스 추적 포함)
package parser

import (
	"bufio"
	"encoding/json"
	"os"

	"github.com/clari/whyso/pkg/model"
)

const maxLineSize = 10 * 1024 * 1024 // 10MB

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
