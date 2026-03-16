//ff:func feature=codemap type=parser control=sequence
//ff:what OpenAPI YAML에서 operationId를 추출

package codemap

import "github.com/smacker/go-tree-sitter/yaml"

func parseOpenAPI(src []byte) []string {
	// Extract operationId values from YAML
	return parseOpenAPIFromYAML(yaml.GetLanguage(), src)
}
