//ff:func feature=cli type=util control=sequence
//ff:what 출력 디렉토리에서 가장 오래된 출력 파일의 mtime을 반환

package main

import (
	"os"
	"path/filepath"
	"strings"
	"time"
)

// oldestOutputMtime returns the oldest mtime among output files in the directory.
// Using the oldest ensures no session is skipped — merge handles duplicates.
func oldestOutputMtime(dir, format string) time.Time {
	var oldest time.Time
	filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return nil
		}
		if !strings.HasSuffix(path, "."+format) {
			return nil
		}
		if oldest.IsZero() || info.ModTime().Before(oldest) {
			oldest = info.ModTime()
		}
		return nil
	})
	return oldest
}
