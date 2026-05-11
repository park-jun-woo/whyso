//ff:func feature=output type=formatter control=iteration dimension=1
//ff:what ChangeEntry의 sources를 YAML 포맷으로 출력

package output

import (
	"fmt"
	"io"

	"github.com/park-jun-woo/whyso/pkg/history"
)

func formatSources(w io.Writer, sources []history.Source) {
	if len(sources) == 1 {
		fmt.Fprintf(w, "    source: %s:%d\n", sources[0].File, sources[0].Line)
		return
	}
	fmt.Fprintf(w, "    sources:\n")
	for _, s := range sources {
		fmt.Fprintf(w, "      - %s:%d\n", s.File, s.Line)
	}
}
