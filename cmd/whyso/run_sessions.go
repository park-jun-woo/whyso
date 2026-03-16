//ff:func feature=cli type=command control=iteration dimension=1
//ff:what sessions 서브커맨드: 세션 목록을 테이블로 출력

package main

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/clari/whyso/pkg/parser"
)

func runSessions() error {
	sessionsDir, err := getSessionsDir()
	if err != nil {
		return err
	}

	sessions, err := parser.ListSessions(sessionsDir)
	if err != nil {
		return err
	}

	if len(sessions) == 0 {
		fmt.Println("No sessions found.")
		return nil
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintf(w, "TIMESTAMP\tSESSION ID\tFIRST MESSAGE\n")
	for _, s := range sessions {
		ts := s.Timestamp.Local().Format("2006-01-02 15:04")
		fmt.Fprintf(w, "%s\t%s\t%s\n", ts, s.ID, s.FirstMessage)
	}
	w.Flush()

	return nil
}
