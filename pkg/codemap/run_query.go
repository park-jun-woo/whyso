//ff:func feature=codemap type=parser control=iteration dimension=1
//ff:what tree-sitter 쿼리를 실행하여 캡처된 이름 목록을 반환

package codemap

import (
	"context"

	sitter "github.com/smacker/go-tree-sitter"
)

func runQuery(lang *sitter.Language, src []byte, query string) []string {
	parser := sitter.NewParser()
	parser.SetLanguage(lang)

	tree, err := parser.ParseCtx(context.Background(), nil, src)
	if err != nil {
		return nil
	}
	root := tree.RootNode()

	q, err := sitter.NewQuery([]byte(query), lang)
	if err != nil {
		return nil
	}

	cursor := sitter.NewQueryCursor()
	cursor.Exec(q, root)

	var names []string
	for {
		match, ok := cursor.NextMatch()
		if !ok {
			break
		}
		for _, capture := range match.Captures {
			name := capture.Node.Content(src)
			if name == "" {
				continue
			}
			names = append(names, name)
		}
	}
	return names
}
