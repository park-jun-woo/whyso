//ff:func feature=codemap type=parser control=sequence
//ff:what STML HTML에서 data 어트리뷰트를 추출

package codemap

import "github.com/smacker/go-tree-sitter/html"

func parseSTML(src []byte) []string {
	return parseSTMLFromHTML(html.GetLanguage(), src)
}
