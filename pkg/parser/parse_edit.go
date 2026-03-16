//ff:func feature=change type=parser control=sequence
//ff:what Edit tool_use ContentBlock에서 FileChange를 파싱
package parser

import (
	"encoding/json"

	"github.com/clari/whyso/pkg/model"
)

func parseEdit(block model.ContentBlock) *FileChange {
	var input editInput
	if err := json.Unmarshal(block.Input, &input); err != nil {
		return nil
	}
	return &FileChange{
		FilePath:  input.FilePath,
		Tool:      "Edit",
		ToolUseID: block.ID,
		OldString: input.OldString,
		NewString: input.NewString,
	}
}
