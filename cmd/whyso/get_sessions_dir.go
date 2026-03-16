//ff:func feature=cli type=util control=iteration dimension=1
//ff:what CLI 인자에서 세션 디렉토리를 결정 (--sessions-dir 또는 자동 탐지)

package main

import (
	"os"

	"github.com/clari/whyso/pkg/parser"
)

func getSessionsDir() (string, error) {
	for i := 2; i < len(os.Args); i++ {
		if os.Args[i] == "--sessions-dir" && i+1 < len(os.Args) {
			return os.Args[i+1], nil
		}
	}
	cwd, err := os.Getwd()
	if err != nil {
		return "", err
	}
	return parser.DetectSessionsDir(cwd)
}
