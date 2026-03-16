//ff:func feature=change type=parser control=selection
//ff:what 단일 ContentBlock에서 Write/Edit FileChange를 추출 (tool_use 외 nil 반환)
package parser

import "github.com/clari/whyso/pkg/model"

func extractBlockChange(block model.ContentBlock) *FileChange {
	if block.Type != "tool_use" {
		return nil
	}
	switch block.Name {
	case "Write":
		return parseWrite(block)
	case "Edit":
		return parseEdit(block)
	default:
		return nil
	}
}
