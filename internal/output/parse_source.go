//ff:func feature=output type=parser control=sequence
//ff:what "파일:라인" 형식 문자열을 Source 구조체로 파싱

package output

import (
	"strconv"
	"strings"

	"github.com/clari/whyso/pkg/history"
)

func parseSource(s string) history.Source {
	idx := strings.LastIndex(s, ":")
	if idx < 0 {
		return history.Source{File: s}
	}
	line, err := strconv.Atoi(s[idx+1:])
	if err != nil {
		return history.Source{File: s}
	}
	return history.Source{File: s[:idx], Line: line}
}
