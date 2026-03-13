package history

import (
	"strings"

	"github.com/clari/whyso/pkg/model"
)

const maxChainDepth = 100

// RecordIndex maps UUID to Record for fast lookup.
type RecordIndex map[string]*model.Record

// BuildIndex creates a UUID → Record map from a slice of records.
func BuildIndex(records []model.Record) RecordIndex {
	idx := make(RecordIndex, len(records))
	for i := range records {
		idx[records[i].UUID] = &records[i]
	}
	return idx
}

// FindUserRequest traces the parentUuid chain backwards to find the original user message.
func FindUserRequest(idx RecordIndex, recordUUID string) string {
	current, ok := idx[recordUUID]
	if !ok {
		return ""
	}

	for i := 0; i < maxChainDepth; i++ {
		if current.ParentUUID == nil {
			if current.IsUserMessage() {
				return current.UserContent()
			}
			return ""
		}

		parent, ok := idx[*current.ParentUUID]
		if !ok {
			return ""
		}

		if parent.IsUserMessage() {
			return parent.UserContent()
		}
		current = parent
	}

	return ""
}

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
		hasText := false
		hasToolUse := false
		var texts []string
		for _, b := range blocks {
			if b.Type == "text" && b.Text != "" {
				hasText = true
				texts = append(texts, b.Text)
			}
			if b.Type == "tool_use" {
				hasToolUse = true
			}
		}
		if hasText && !hasToolUse {
			return strings.Join(texts, "\n")
		}
	}
	return ""
}
