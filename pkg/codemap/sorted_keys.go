//ff:func feature=codemap type=util control=iteration dimension=1
//ff:what 맵의 키를 정렬된 슬라이스로 반환

package codemap

import "sort"

func sortedKeys(m map[string][]string) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
