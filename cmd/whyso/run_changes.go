//ff:func feature=cli type=command control=iteration dimension=1
//ff:what changes 서브커맨드: 파일 변경 내역을 테이블로 출력

package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/tabwriter"

	"github.com/clari/whyso/pkg/history"
	"github.com/clari/whyso/pkg/parser"
)

func runChanges() error {
	sessionsDir, err := getSessionsDir()
	if err != nil {
		return err
	}

	entries, err := os.ReadDir(sessionsDir)
	if err != nil {
		return err
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintf(w, "TIMESTAMP\tTOOL\tFILE\tUSER REQUEST\n")

	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".jsonl") {
			continue
		}
		path := filepath.Join(sessionsDir, entry.Name())
		records, err := parser.ParseSession(path)
		if err != nil {
			continue
		}

		changes := parser.ExtractChanges(records)
		idx := history.BuildIndex(records)

		for _, fc := range changes {
			formatChangeRow(w, fc, idx)
		}
	}
	w.Flush()
	return nil
}
