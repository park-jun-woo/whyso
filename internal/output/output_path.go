//ff:func feature=output type=util control=sequence
//ff:what 상대 경로와 포맷에서 출력 파일 경로를 생성

package output

import "fmt"

// OutputPath returns the output file path for a given relative path.
func OutputPath(outputDir, relPath, format string) string {
	return fmt.Sprintf("%s/%s.%s", outputDir, relPath, format)
}
