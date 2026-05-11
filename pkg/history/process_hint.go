//ff:func feature=history type=builder control=sequence
//ff:what 단일 BashHint를 Hint로 변환하여 관련 FileHistory에 추가

package history

import (
	"github.com/park-jun-woo/whyso/pkg/parser"
)

func processHint(bh parser.BashHint, histories map[string]*FileHistory, projectRoot string) {
	relSrc := toRelPath(projectRoot, bh.SrcPath)
	relDst := toRelPath(projectRoot, bh.DstPath)

	hint := Hint{
		Timestamp: bh.Timestamp,
		Session:   bh.SessionID,
		Command:   bh.Command,
		Source:    Source{File: bh.SourceFile, Line: bh.SourceLine},
	}

	addHintIfExists(histories, relSrc, hint)
	addHintIfExists(histories, relDst, hint)
}
