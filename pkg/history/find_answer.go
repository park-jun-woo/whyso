//ff:func feature=history type=util control=iteration dimension=1
//ff:what recordIndex 이후의 텍스트 전용 어시스턴트 메시지(답변)를 찾음
package history

import (
	"strings"

	"github.com/clari/whyso/pkg/model"
)

// FindAnswer looks forward from recordIndex in the records slice
// for the next text-only assistant message (the answer after tool_use).
func FindAnswer(records []model.Record, recordIndex int) string {
	for i := recordIndex + 1; i < len(records) && i < recordIndex+10; i++ {
		rec := &records[i]
		if rec.Type != "assistant" {
			continue
		}
		blocks := rec.ContentBlocks()
		if blocks == nil {
			continue
		}
		texts, hasToolUse := extractTextBlocks(blocks)
		if len(texts) > 0 && !hasToolUse {
			return strings.Join(texts, "\n")
		}
	}
	return ""
}
