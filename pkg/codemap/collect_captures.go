//ff:func feature=codemap type=parser control=iteration dimension=1
//ff:what 쿼리 매치에서 캡처된 이름을 수집하여 슬라이스에 추가
package codemap

import sitter "github.com/smacker/go-tree-sitter"

func collectCaptures(names []string, match *sitter.QueryMatch, src []byte) []string {
	for _, capture := range match.Captures {
		name := capture.Node.Content(src)
		if name == "" {
			continue
		}
		names = append(names, name)
	}
	return names
}
