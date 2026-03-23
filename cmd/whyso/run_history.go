//ff:func feature=cli type=command control=sequence
//ff:what history 서브커맨드: 파일 변경 이력을 빌드하고 출력

package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/clari/whyso/internal/output"
	"github.com/clari/whyso/pkg/history"
)

func runHistory() error {
	if err := checkDeps(); err != nil {
		return err
	}

	sessionsDir, err := getSessionsDir()
	if err != nil {
		return err
	}

	projectRoot, err := os.Getwd()
	if err != nil {
		return err
	}

	target, format, outputDir, all, quiet, reset := parseHistoryArgs()

	// always cache to .whyso/ unless --output overrides
	cacheDir := filepath.Join(projectRoot, ".whyso")
	if outputDir == "" {
		outputDir = cacheDir
	}

	if target == "" {
		return fmt.Errorf("usage: whyso history <file|dir> [--all] [--format yaml|json] [--output dir]")
	}

	// resolve target to absolute path for matching
	absTarget, err := filepath.Abs(target)
	if err != nil {
		return err
	}

	targetInfo, _ := os.Stat(absTarget)

	filter := makeFilter(targetInfo, all, projectRoot, absTarget)

	if reset {
		clearCache(outputDir, format)
	}

	since := resolveSince(targetInfo, projectRoot, absTarget, outputDir, format)
	var histories map[string]*history.FileHistory
	var buildErr error
	if since.IsZero() {
		histories, buildErr = history.BuildHistories(sessionsDir, projectRoot, filter)
	} else {
		histories, buildErr = history.BuildHistoriesIncremental(sessionsDir, projectRoot, since, filter)
	}
	if buildErr != nil {
		return buildErr
	}

	isSingleFile := targetInfo == nil || !targetInfo.IsDir()
	if len(histories) == 0 && isSingleFile {
		histories = retryWithRename(histories, sessionsDir, projectRoot, absTarget)
	}

	if len(histories) > 0 {
		if err := output.WriteHistories(histories, outputDir, format); err != nil {
			return err
		}
	}
	if !quiet && outputDir == filepath.Join(projectRoot, ".whyso") && isSingleFile {
		printHistoryOutput(histories, format, outputDir, projectRoot, absTarget)
	}
	return nil
}
