//ff:func feature=cli type=formatter control=sequence
//ff:what 단일 파일 변경을 테이블 행으로 포맷하여 출력

package main

import (
	"fmt"
	"text/tabwriter"

	"github.com/clari/whyso/pkg/history"
	"github.com/clari/whyso/pkg/parser"
)

func formatChangeRow(w *tabwriter.Writer, fc parser.FileChange, idx history.RecordIndex) {
	ts := fc.Timestamp.Local().Format("2006-01-02 15:04")
	userReq := history.FindUserRequest(idx, fc.RecordUUID)
	if len(userReq) > 60 {
		userReq = userReq[:60] + "..."
	}
	fmt.Fprintf(w, "%s\t%s\t%s\t%s\n", ts, fc.Tool, fc.FilePath, userReq)
}
