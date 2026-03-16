//ff:func feature=cli type=parser control=iteration dimension=1
//ff:what history 서브커맨드의 CLI 인자를 파싱

package main

import "os"

func parseHistoryArgs() (target, format, outputDir string, all, quiet, reset bool) {
	format = "yaml"
	for i := 2; i < len(os.Args); i++ {
		arg := os.Args[i]
		if arg == "--format" && i+1 < len(os.Args) {
			format = os.Args[i+1]
			i++
			continue
		}
		if arg == "--output" && i+1 < len(os.Args) {
			outputDir = os.Args[i+1]
			i++
			continue
		}
		if arg == "--quiet" || arg == "-q" {
			quiet = true
			continue
		}
		if arg == "--all" {
			all = true
			continue
		}
		if arg == "--reset" {
			reset = true
			continue
		}
		if arg == "--sessions-dir" {
			i++ // already handled by getSessionsDir
			continue
		}
		if target == "" {
			target = arg
		}
	}
	return
}
