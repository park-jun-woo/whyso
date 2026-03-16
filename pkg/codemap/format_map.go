//ff:func feature=codemap type=formatter control=iteration dimension=1
//ff:what 섹션 목록을 표준 맵 포맷으로 출력

package codemap

import (
	"fmt"
	"io"
	"sort"
	"strings"
)

// FormatMap writes the map to w in the standard format.
func FormatMap(w io.Writer, sections []Section) {
	fmt.Fprintln(w, "# whyso/v1")
	for _, sec := range sections {
		fmt.Fprintln(w)
		fmt.Fprintf(w, "## %s\n", sec.Language)

		groups := sortedKeys(sec.Groups)
		for _, g := range groups {
			keywords := sec.Groups[g]
			sort.Strings(keywords)
			keywords = dedupe(keywords)
			fmt.Fprintf(w, "[%s]%s\n", g, strings.Join(keywords, ","))
		}
	}
}
