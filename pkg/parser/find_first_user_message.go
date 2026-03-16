//ff:func feature=session type=parser control=iteration dimension=1
//ff:what Record 슬라이스에서 첫 번째 사용자 메시지의 타임스탬프와 내용을 반환
package parser

import (
	"time"

	"github.com/clari/whyso/pkg/model"
)

func findFirstUserMessage(records []model.Record) (time.Time, string) {
	for _, rec := range records {
		if rec.IsUserMessage() {
			content := rec.UserContent()
			if len(content) > 80 {
				content = content[:80] + "..."
			}
			return rec.Timestamp, content
		}
	}
	return time.Time{}, ""
}
