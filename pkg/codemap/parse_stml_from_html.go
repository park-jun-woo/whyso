//ff:func feature=codemap type=parser control=sequence
//ff:what HTML tree-sitter로 data 어트리뷰트를 파싱

package codemap

import (
	"context"

	sitter "github.com/smacker/go-tree-sitter"
)

func parseSTMLFromHTML(lang *sitter.Language, src []byte) []string {
	parser := sitter.NewParser()
	parser.SetLanguage(lang)

	tree, err := parser.ParseCtx(context.Background(), nil, src)
	if err != nil {
		return nil
	}
	root := tree.RootNode()

	var names []string
	findDataAttributes(root, src, &names)
	return names
}
