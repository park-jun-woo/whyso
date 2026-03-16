//ff:func feature=cli type=util control=sequence
//ff:what 출력 디렉토리의 캐시 파일을 삭제

package main

import (
	"os"
	"path/filepath"
	"strings"
)

func clearCache(dir, format string) {
	filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return nil
		}
		if strings.HasSuffix(path, "."+format) {
			os.Remove(path)
		}
		return nil
	})
}
