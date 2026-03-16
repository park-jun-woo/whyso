//ff:func feature=change type=parser control=iteration dimension=1
//ff:what 단일 Record의 ContentBlock에서 파일 변경을 추출하여 슬라이스에 추가
package parser

import (
	"github.com/clari/whyso/pkg/model"
)

func collectBlockChanges(changes []FileChange, rec *model.Record, recordIndex int) []FileChange {
	for _, block := range rec.ContentBlocks() {
		fc := extractBlockChange(block)
		if fc == nil {
			continue
		}
		fc.RecordUUID = rec.UUID
		fc.RecordIndex = recordIndex
		fc.Timestamp = rec.Timestamp
		fc.SessionID = rec.SessionID
		fc.SourceFile = rec.SourceFile
		fc.SourceLine = rec.SourceLine
		changes = append(changes, *fc)
	}
	return changes
}
