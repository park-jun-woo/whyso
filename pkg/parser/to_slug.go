//ff:func feature=session type=util control=sequence
//ff:what 절대 경로를 Claude Code 프로젝트 slug로 변환
package parser

import "strings"

// toSlug converts an absolute path to a Claude Code project slug.
// e.g. /home/user/.clari/project → -home-user--clari-project
func toSlug(abs string) string {
	s := strings.TrimPrefix(abs, "/")
	s = strings.ReplaceAll(s, "/", "-")
	s = strings.ReplaceAll(s, ".", "-")
	return "-" + s
}
