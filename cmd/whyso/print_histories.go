//ff:func feature=cli type=formatter control=iteration dimension=1
//ff:what FileHistory 맵의 모든 이력을 stdout에 출력
package main

import (
	"github.com/clari/whyso/pkg/history"
)

func printHistories(histories map[string]*history.FileHistory, format string) {
	for _, h := range histories {
		printSingleHistory(h, format)
	}
}
