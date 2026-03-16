//ff:func feature=codemap type=util control=iteration dimension=1
//ff:what 정렬된 문자열 슬라이스에서 중복을 제거

package codemap

func dedupe(ss []string) []string {
	if len(ss) == 0 {
		return ss
	}
	result := []string{ss[0]}
	for i := 1; i < len(ss); i++ {
		if ss[i] != ss[i-1] {
			result = append(result, ss[i])
		}
	}
	return result
}
