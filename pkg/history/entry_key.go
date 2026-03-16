//ff:func feature=history type=util control=sequence
//ff:what ChangeEntry의 중복 판별용 키를 생성
package history

import "fmt"

func entryKey(e ChangeEntry) string {
	ts := e.Timestamp.Unix()
	if len(e.Sources) > 0 {
		return fmt.Sprintf("%d:%s:%d", ts, e.Sources[0].File, e.Sources[0].Line)
	}
	return fmt.Sprintf("%d:%s:%s", ts, e.Session, e.Tool)
}
