//ff:func feature=change type=parser control=sequence
//ff:what Write tool_use ContentBlock에서 FileChange를 파싱
package parser

import (
	"encoding/json"

	"github.com/clari/whyso/pkg/model"
)

func parseWrite(block model.ContentBlock) *FileChange {
	var input writeInput
	if err := json.Unmarshal(block.Input, &input); err != nil {
		return nil
	}
	return &FileChange{
		FilePath:  input.FilePath,
		Tool:      "Write",
		ToolUseID: block.ID,
		Content:   input.Content,
	}
}
