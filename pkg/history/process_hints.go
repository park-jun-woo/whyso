//ff:func feature=history type=builder control=iteration dimension=1
//ff:what BashHint를 순회하며 관련 FileHistory에 Hint로 추가

package history

import (
	"github.com/clari/whyso/pkg/parser"
)

func processHints(hints []parser.BashHint, histories map[string]*FileHistory, projectRoot string) {
	for _, bh := range hints {
		processHint(bh, histories, projectRoot)
	}
}
