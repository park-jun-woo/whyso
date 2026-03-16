//ff:func feature=codemap type=parser control=iteration dimension=1
//ff:what YAML AST에서 operationId 키-값을 재귀 탐색

package codemap

import sitter "github.com/smacker/go-tree-sitter"

func findOperationIDs(node *sitter.Node, src []byte, names *[]string) {
	for i := 0; i < int(node.ChildCount()); i++ {
		child := node.Child(i)
		if child.Type() != "block_mapping_pair" {
			findOperationIDs(child, src, names)
			continue
		}
		key := child.ChildByFieldName("key")
		value := child.ChildByFieldName("value")
		if key == nil || value == nil || key.Content(src) != "operationId" {
			findOperationIDs(child, src, names)
			continue
		}
		val := trimQuotes(value.Content(src))
		if val != "" {
			*names = append(*names, val)
		}
	}
}
