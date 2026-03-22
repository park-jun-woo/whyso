//ff:func feature=output type=formatter control=iteration dimension=1
//ff:what 모든 이력을 출력 디렉토리에 기록 (기존 파일과 병합)

package output

import (
	"github.com/clari/whyso/pkg/history"
)

// WriteHistories writes all histories to the output directory, mirroring file paths.
// If an output file already exists, it merges new entries with existing ones.
func WriteHistories(histories map[string]*history.FileHistory, outputDir, format string) error {
	for relPath, h := range histories {
		if err := writeHistory(relPath, h, outputDir, format); err != nil {
			return err
		}
	}
	return nil
}
