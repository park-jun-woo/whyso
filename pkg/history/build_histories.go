//ff:func feature=history type=builder control=sequence
//ff:what 세션 디렉토리의 모든 세션에서 파일 이력을 빌드
package history

import "time"

// BuildHistories builds file histories from all sessions in the directory.
func BuildHistories(sessionsDir, projectRoot string, filter func(string) bool) (map[string]*FileHistory, error) {
	return buildHistories(sessionsDir, projectRoot, time.Time{}, filter)
}
