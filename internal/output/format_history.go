//ff:func feature=output type=formatter control=selection
//ff:what 포맷에 따라 이력을 파일에 기록 (JSON 또는 YAML)

package output

import (
	"io"

	"github.com/clari/whyso/pkg/history"
)

func formatHistory(w io.Writer, h *history.FileHistory, format string) error {
	switch format {
	case "json":
		return FormatJSON(w, h)
	default:
		return FormatYAML(w, h)
	}
}
