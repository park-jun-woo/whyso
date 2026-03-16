//ff:func feature=history type=util control=iteration dimension=1
//ff:what ContentBlock 배열에서 텍스트만 추출하여 tool_use 없는 순수 텍스트 여부와 함께 반환
package history

import "github.com/clari/whyso/pkg/model"

func extractTextBlocks(blocks []model.ContentBlock) (texts []string, hasToolUse bool) {
	for _, b := range blocks {
		if b.Type == "text" && b.Text != "" {
			texts = append(texts, b.Text)
		}
		if b.Type == "tool_use" {
			hasToolUse = true
		}
	}
	return
}
