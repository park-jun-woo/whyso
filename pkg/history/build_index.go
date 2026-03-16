//ff:func feature=history type=builder control=iteration dimension=1
//ff:what Record 슬라이스에서 UUID→Record 인덱스를 생성
package history

import "github.com/clari/whyso/pkg/model"

// BuildIndex creates a UUID -> Record map from a slice of records.
func BuildIndex(records []model.Record) RecordIndex {
	idx := make(RecordIndex, len(records))
	for i := range records {
		idx[records[i].UUID] = &records[i]
	}
	return idx
}
