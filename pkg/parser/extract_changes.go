//ff:func feature=change type=parser control=iteration dimension=1
//ff:what Write/Edit tool_use 블록에서 파일 변경 목록을 추출
package parser

import (
	"github.com/clari/whyso/pkg/model"
)

// ExtractChanges extracts file changes from Write/Edit tool_use blocks.
func ExtractChanges(records []model.Record) []FileChange {
	var changes []FileChange

	for i := range records {
		rec := &records[i]
		if rec.Type != "assistant" {
			continue
		}

		blocks := rec.ContentBlocks()
		for _, block := range blocks {
			fc := extractBlockChange(block)
			if fc == nil {
				continue
			}
			fc.RecordUUID = rec.UUID
			fc.RecordIndex = i
			fc.Timestamp = rec.Timestamp
			fc.SessionID = rec.SessionID
			fc.SourceFile = rec.SourceFile
			fc.SourceLine = rec.SourceLine
			changes = append(changes, *fc)
		}
	}

	return changes
}
