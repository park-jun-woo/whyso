//ff:type feature=session type=model
//ff:what 세션 목록 표시용 요약 정보
package model

import "time"

// SessionInfo holds summary information about a session.
type SessionInfo struct {
	ID           string
	Timestamp    time.Time
	FirstMessage string
}
