//ff:func feature=change type=parser control=sequence
//ff:what Bash tool_use의 command에서 cp/mv 명령을 감지하여 BashHint 반환

package parser

import (
	"encoding/json"
	"regexp"

	"github.com/clari/whyso/pkg/model"
)

var bashCpMvPattern = regexp.MustCompile(`^\s*(cp|mv)\s+(?:-[a-zA-Z]+\s+)*(.+?)\s+(.+?)\s*$`)

func parseBashHint(block model.ContentBlock) *BashHint {
	var input struct {
		Command string `json:"command"`
	}
	if err := json.Unmarshal(block.Input, &input); err != nil {
		return nil
	}

	m := bashCpMvPattern.FindStringSubmatch(input.Command)
	if m == nil {
		return nil
	}

	return &BashHint{
		Command: input.Command,
		SrcPath: m[2],
		DstPath: m[3],
	}
}
