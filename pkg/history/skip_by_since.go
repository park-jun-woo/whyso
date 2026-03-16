//ff:func feature=history type=util control=sequence
//ff:what since 이후 미수정된 디렉토리 엔트리를 스킵할지 판별
package history

import (
	"os"
	"time"
)

func skipBySince(entry os.DirEntry, since time.Time) bool {
	info, err := entry.Info()
	if err != nil {
		return true
	}
	return !info.ModTime().After(since)
}
