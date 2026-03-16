//ff:func feature=output type=formatter control=sequence
//ff:what FileHistory를 JSON 포맷으로 출력

package output

import (
	"encoding/json"
	"io"

	"github.com/clari/whyso/pkg/history"
)

// FormatJSON writes a FileHistory as JSON to w.
func FormatJSON(w io.Writer, h *history.FileHistory) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(h)
}
