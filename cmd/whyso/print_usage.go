//ff:func feature=cli type=command control=sequence
//ff:what CLI 사용법을 stderr에 출력

package main

import (
	"fmt"
	"os"
)

func printUsage() {
	fmt.Fprintln(os.Stderr, "Usage: whyso <command>")
	fmt.Fprintln(os.Stderr, "")
	fmt.Fprintln(os.Stderr, "Commands:")
	fmt.Fprintln(os.Stderr, "  sessions    List sessions for the current project")
	fmt.Fprintln(os.Stderr, "  changes     List file changes across all sessions")
	fmt.Fprintln(os.Stderr, "  history     Show file change history")
	fmt.Fprintln(os.Stderr, "  map         Generate keyword map (functions, endpoints, etc.)")
	fmt.Fprintln(os.Stderr, "")
	fmt.Fprintln(os.Stderr, "")
	fmt.Fprintln(os.Stderr, "Map options:")
	fmt.Fprintln(os.Stderr, "  -o <file>                Output file (default: .whyso/_map.md)")
	fmt.Fprintln(os.Stderr, "  -f, --force              Force regeneration")
	fmt.Fprintln(os.Stderr, "")
	fmt.Fprintln(os.Stderr, "History options:")
	fmt.Fprintln(os.Stderr, "  --format <yaml|json>     Output format (default: yaml)")
	fmt.Fprintln(os.Stderr, "  --output <dir>           Write to directory only, no stdout (default: .whyso/)")
	fmt.Fprintln(os.Stderr, "  -q, --quiet              Suppress stdout output")
	fmt.Fprintln(os.Stderr, "  --reset                  Clear cache and rebuild from scratch")
	fmt.Fprintln(os.Stderr, "  --all                    All files in directory")
	fmt.Fprintln(os.Stderr, "")
	fmt.Fprintln(os.Stderr, "Options:")
	fmt.Fprintln(os.Stderr, "  --sessions-dir <path>    Override sessions directory")
}
