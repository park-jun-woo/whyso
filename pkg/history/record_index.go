//ff:type feature=history type=model
//ff:what UUID→Record 빠른 조회용 인덱스
package history

import "github.com/clari/whyso/pkg/model"

// RecordIndex maps UUID to Record for fast lookup.
type RecordIndex map[string]*model.Record
