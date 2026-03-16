//ff:type feature=codemap type=model
//ff:what 언어별 그룹화된 키워드 섹션

package codemap

// Section represents a language section with grouped keywords.
type Section struct {
	Language string
	Groups   map[string][]string // group name -> keywords
}
