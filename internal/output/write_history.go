//ff:func feature=output type=formatter control=sequence
//ff:what 단일 파일 이력을 출력 경로에 기록 (기존 파일과 병합)

package output

import (
	"os"
	"path/filepath"

	"github.com/clari/whyso/pkg/history"
)

func writeHistory(relPath string, h *history.FileHistory, outputDir, format string) error {
	outPath := filepath.Join(outputDir, relPath+"."+format)

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

	writeErr := formatHistory(f, h, format)
	f.Close()

	return writeErr
}
