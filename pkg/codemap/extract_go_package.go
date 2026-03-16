//ff:func feature=codemap type=parser control=iteration dimension=1
//ff:what Go 소스에서 package 선언을 추출

package codemap

import "strings"

func extractGoPackage(src []byte) string {
	for _, line := range strings.Split(string(src), "\n") {
		line = strings.TrimSpace(line)
		if !strings.HasPrefix(line, "package ") {
			continue
		}
		parts := strings.Fields(line)
		if len(parts) >= 2 {
			return parts[1]
		}
	}
	return "unknown"
}
