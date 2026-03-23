//ff:func feature=output type=formatter control=iteration dimension=1
//ff:what Hint 목록을 YAML hints 섹션으로 출력

package output

import (
	"fmt"
	"io"
	"time"

	"github.com/clari/whyso/pkg/history"
)

func formatHints(w io.Writer, hints []history.Hint) {
	if len(hints) == 0 {
		return
	}
	fmt.Fprintf(w, "hints:\n")
	for _, h := range hints {
		fmt.Fprintf(w, "  - timestamp: %s\n", h.Timestamp.Format(time.RFC3339))
		fmt.Fprintf(w, "    session: %s\n", h.Session)
		fmt.Fprintf(w, "    command: %q\n", h.Command)
		fmt.Fprintf(w, "    source: %s:%d\n", h.Source.File, h.Source.Line)
	}
}
