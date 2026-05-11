//ff:func feature=cli type=formatter control=selection
//ff:what 단일 FileHistory를 format에 따라 stdout에 출력
package main

import (
	"fmt"
	"os"

	"github.com/park-jun-woo/whyso/internal/output"
	"github.com/park-jun-woo/whyso/pkg/history"
)

func printSingleHistory(h *history.FileHistory, format string) {
	switch format {
	case "json":
		output.FormatJSON(os.Stdout, h)
	default:
		output.FormatYAML(os.Stdout, h)
	}
	fmt.Println("---")
}
