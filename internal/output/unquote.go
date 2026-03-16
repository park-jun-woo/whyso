//ff:func feature=output type=util control=sequence
//ff:what 문자열 앞뒤 큰따옴표를 해석하여 제거

package output

import "strconv"

func unquote(s string) string {
	if len(s) >= 2 && s[0] == '"' && s[len(s)-1] == '"' {
		u, err := strconv.Unquote(s)
		if err == nil {
			return u
		}
	}
	return s
}
