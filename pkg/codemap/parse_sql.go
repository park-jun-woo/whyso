//ff:func feature=codemap type=parser control=iteration dimension=1
//ff:what SQL 파일에서 쿼리명을 추출

package codemap

func parseSQL(src []byte) []string {
	matches := reSQLName.FindAllSubmatch(src, -1)
	var names []string
	for _, m := range matches {
		name := string(m[1])
		if name != "" {
			names = append(names, name)
		}
	}
	return names
}
