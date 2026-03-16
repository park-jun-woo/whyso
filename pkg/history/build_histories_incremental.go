//ff:func feature=history type=builder control=sequence
//ff:what since 이후 수정된 세션에서만 파일 이력을 빌드
package history

import "time"

// BuildHistoriesIncremental builds file histories only from sessions modified after since.
func BuildHistoriesIncremental(sessionsDir, projectRoot string, since time.Time, filter func(string) bool) (map[string]*FileHistory, error) {
	return buildHistories(sessionsDir, projectRoot, since, filter)
}
