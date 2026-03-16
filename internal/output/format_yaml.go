//ff:func feature=output type=formatter control=iteration dimension=1
//ff:what FileHistory를 YAML 포맷으로 출력

package output

import (
	"fmt"
	"io"
	"time"

	"github.com/clari/whyso/pkg/history"
)

// FormatYAML writes a FileHistory as YAML to w.
func FormatYAML(w io.Writer, h *history.FileHistory) error {
	fmt.Fprintf(w, "apiVersion: whyso/v1\n")
	fmt.Fprintf(w, "file: %s\n", h.File)
	fmt.Fprintf(w, "created: %s\n", h.Created.Format(time.RFC3339))
	fmt.Fprintf(w, "history:\n")
	for _, e := range h.History {
		fmt.Fprintf(w, "  - timestamp: %s\n", e.Timestamp.Format(time.RFC3339))
		fmt.Fprintf(w, "    session: %s\n", e.Session)
		fmt.Fprintf(w, "    user_request: %q\n", e.UserRequest)
		if e.Answer != "" {
			fmt.Fprintf(w, "    answer: %q\n", e.Answer)
		}
		fmt.Fprintf(w, "    tool: %s\n", e.Tool)
		if e.Subagent {
			fmt.Fprintf(w, "    subagent: true\n")
		}
		formatSources(w, e.Sources)
	}
	return nil
}
