//ff:func feature=codemap type=parser control=iteration dimension=1
//ff:what HTML AST에서 data-fetch/data-action 어트리뷰트를 재귀 탐색

package codemap

import sitter "github.com/smacker/go-tree-sitter"

func findDataAttributes(node *sitter.Node, src []byte, names *[]string) {
	if node.Type() == "attribute" {
		if val := extractAttribute(node, src); val != "" {
			*names = append(*names, val)
		}
	}
	for i := 0; i < int(node.ChildCount()); i++ {
		findDataAttributes(node.Child(i), src, names)
	}
}
