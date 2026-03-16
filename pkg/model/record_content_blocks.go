//ff:func feature=session type=model control=sequence
//ff:what message.content를 ContentBlock 배열로 파싱
package model

import "encoding/json"

// ContentBlocks parses message.content as an array of ContentBlock.
func (r *Record) ContentBlocks() []ContentBlock {
	if len(r.Message.Content) == 0 || r.Message.Content[0] != '[' {
		return nil
	}
	var blocks []ContentBlock
	if err := json.Unmarshal(r.Message.Content, &blocks); err != nil {
		return nil
	}
	return blocks
}
