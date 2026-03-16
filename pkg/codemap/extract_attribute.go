//ff:func feature=codemap type=parser control=iteration dimension=1
//ff:what 단일 attribute 노드에서 name-value를 추출하여 data-fetch/data-action 결과를 반환

package codemap

import sitter "github.com/smacker/go-tree-sitter"

func extractAttribute(node *sitter.Node, src []byte) string {
	var attrName, attrVal string
	for i := 0; i < int(node.ChildCount()); i++ {
		child := node.Child(i)
		switch child.Type() {
		case "attribute_name":
			attrName = child.Content(src)
		case "quoted_attribute_value", "attribute_value":
			attrVal = trimQuotes(child.Content(src))
		}
	}
	if (attrName == "data-fetch" || attrName == "data-action") && attrVal != "" {
		return attrName + ":" + attrVal
	}
	return ""
}
