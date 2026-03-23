//ff:func feature=change type=parser control=iteration dimension=1
//ff:what 모든 Record에서 Bash cp/mv 힌트를 추출

package parser

import (
	"github.com/clari/whyso/pkg/model"
)

// ExtractBashHints extracts cp/mv hints from Bash tool_use blocks.
func ExtractBashHints(records []model.Record) []BashHint {
	var hints []BashHint
	for _, rec := range records {
		if rec.Type != "assistant" {
			continue
		}
		hints = collectBlockHints(hints, &rec)
	}
	return hints
}
