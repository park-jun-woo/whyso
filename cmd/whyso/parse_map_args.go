//ff:func feature=cli type=parser control=iteration dimension=1
//ff:what map 서브커맨드의 CLI 인자를 파싱

package main

import "os"

func parseMapArgs() (target, outputFile string, force bool) {
	for i := 2; i < len(os.Args); i++ {
		arg := os.Args[i]
		if arg == "-o" && i+1 < len(os.Args) {
			outputFile = os.Args[i+1]
			i++
			continue
		}
		if arg == "-f" || arg == "--force" {
			force = true
			continue
		}
		if target == "" {
			target = arg
		}
	}
	return
}
