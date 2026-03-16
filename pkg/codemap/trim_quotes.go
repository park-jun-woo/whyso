//ff:func feature=codemap type=util control=sequence
//ff:what 문자열 앞뒤 따옴표를 제거

package codemap

func trimQuotes(s string) string {
	if len(s) >= 2 {
		if (s[0] == '"' && s[len(s)-1] == '"') || (s[0] == '\'' && s[len(s)-1] == '\'') {
			return s[1 : len(s)-1]
		}
	}
	return s
}
