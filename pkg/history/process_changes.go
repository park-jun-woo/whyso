//ff:func feature=history type=builder control=iteration dimension=1
//ff:what 단일 세션의 변경 목록을 순회하며 파일 이력에 추가
package history

import (
	"github.com/clari/whyso/pkg/model"
	"github.com/clari/whyso/pkg/parser"
)

func processChanges(changes []parser.FileChange, idx RecordIndex, histories map[string]*FileHistory, projectRoot string, sessionID string, records []model.Record, filter func(string) bool) {
	for _, fc := range changes {
		processChange(fc, idx, histories, projectRoot, sessionID, records, filter)
	}
}
