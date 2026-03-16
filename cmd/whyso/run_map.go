//ff:func feature=cli type=command control=sequence
//ff:what map 서브커맨드: 키워드 맵을 생성하고 출력

package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/clari/whyso/pkg/codemap"
)

func runMap() error {
	target, outputFile, force := parseMapArgs()

	if target == "" {
		target = "."
	}

	absTarget, err := filepath.Abs(target)
	if err != nil {
		return err
	}

	// default output: .whyso/_map.md
	if outputFile == "" {
		cwd, err := os.Getwd()
		if err != nil {
			return err
		}
		defaultDir := filepath.Join(cwd, ".whyso")
		if err := os.MkdirAll(defaultDir, 0755); err != nil {
			return err
		}
		outputFile = filepath.Join(defaultDir, "_map.md")
	}

	// skip if no source files are newer than _map.md
	if !force && !codemap.NeedsUpdate(absTarget, outputFile) {
		fmt.Fprintln(os.Stderr, "up to date")
		return nil
	}

	sections, err := codemap.BuildMap(absTarget)
	if err != nil {
		return err
	}

	if len(sections) == 0 {
		fmt.Println("No keywords found.")
		return nil
	}

	// write to file
	f, err := os.Create(outputFile)
	if err != nil {
		return err
	}
	defer f.Close()
	codemap.FormatMap(f, sections)

	// also stdout
	codemap.FormatMap(os.Stdout, sections)
	return nil
}
