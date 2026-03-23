//ff:func feature=change type=parser control=iteration dimension=1
//ff:what 단일 Record의 ContentBlock에서 Bash cp/mv 힌트를 추출하여 슬라이스에 추가

package parser

import (
	"github.com/clari/whyso/pkg/model"
)

func collectBlockHints(hints []BashHint, rec *model.Record) []BashHint {
	for _, block := range rec.ContentBlocks() {
		if block.Type != "tool_use" || block.Name != "Bash" {
			continue
		}
		hint := parseBashHint(block)
		if hint == nil {
			continue
		}
		hint.Timestamp = rec.Timestamp
		hint.SessionID = rec.SessionID
		hint.SourceFile = rec.SourceFile
		hint.SourceLine = rec.SourceLine
		hints = append(hints, *hint)
	}
	return hints
}
