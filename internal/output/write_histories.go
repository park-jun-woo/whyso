//ff:func feature=output type=formatter control=iteration dimension=1
//ff:what 모든 이력을 출력 디렉토리에 기록 (기존 파일과 병합)

package output

import (
	"os"
	"path/filepath"

	"github.com/clari/whyso/pkg/history"
)

// WriteHistories writes all histories to the output directory, mirroring file paths.
// If an output file already exists, it merges new entries with existing ones.
func WriteHistories(histories map[string]*history.FileHistory, outputDir, format string) error {
	for relPath, h := range histories {
		outPath := filepath.Join(outputDir, relPath+"."+format)

		// merge with existing file if present
		existing, _ := ReadYAML(outPath)
		if format == "yaml" && existing != nil {
			h = history.Merge(existing, h)
		}

		if err := os.MkdirAll(filepath.Dir(outPath), 0o755); err != nil {
			return err
		}

		f, err := os.Create(outPath)
		if err != nil {
			return err
		}

		var writeErr error
		switch format {
		case "json":
			writeErr = FormatJSON(f, h)
		default:
			writeErr = FormatYAML(f, h)
		}

		f.Close()
		if writeErr != nil {
			return writeErr
		}
	}
	return nil
}
