//ff:func feature=cli type=command control=sequence
//ff:what history 서브커맨드: 파일 변경 이력을 빌드하고 출력

package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/clari/whyso/internal/output"
	"github.com/clari/whyso/pkg/history"
)

func runHistory() error {
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

	targetInfo, err := os.Stat(absTarget)
	if err != nil {
		return err
	}

	filter := func(relPath string) bool {
		if targetInfo.IsDir() {
			if !all {
				return false
			}
			targetRel, err := filepath.Rel(projectRoot, absTarget)
			if err != nil {
				return false
			}
			if targetRel == "." {
				return true
			}
			return strings.HasPrefix(relPath, targetRel+"/") || relPath == targetRel
		}
		targetRel, err := filepath.Rel(projectRoot, absTarget)
		if err != nil {
			return false
		}
		return relPath == targetRel
	}

	if reset {
		clearCache(outputDir, format)
	}

	// 타겟 범위 내 캐시 mtime만 확인
	var since time.Time
	if !targetInfo.IsDir() {
		targetRel, _ := filepath.Rel(projectRoot, absTarget)
		cachedPath := output.OutputPath(outputDir, targetRel, format)
		if info, err := os.Stat(cachedPath); err == nil {
			since = info.ModTime()
		}
	} else {
		targetRel, _ := filepath.Rel(projectRoot, absTarget)
		if targetRel == "." {
			since = oldestOutputMtime(outputDir, format)
		} else {
			since = oldestOutputMtime(filepath.Join(outputDir, targetRel), format)
		}
	}
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

	if len(histories) > 0 {
		// 새 변경 사항 캐시에 기록
		if err := output.WriteHistories(histories, outputDir, format); err != nil {
			return err
		}
	}

	// stdout: 단일 파일만, -q 아니고 기본 캐시 경로일 때
	if !quiet && outputDir == filepath.Join(projectRoot, ".whyso") && !targetInfo.IsDir() {
		// 새 결과가 있으면 그걸 출력, 없으면 기존 캐시에서 읽기
		if len(histories) > 0 {
			for _, h := range histories {
				switch format {
				case "json":
					output.FormatJSON(os.Stdout, h)
				default:
					output.FormatYAML(os.Stdout, h)
				}
				fmt.Println("---")
			}
		} else {
			// 기존 캐시 파일에서 읽어서 출력
			targetRel, _ := filepath.Rel(projectRoot, absTarget)
			cachedPath := output.OutputPath(outputDir, targetRel, format)
			if cached, err := output.ReadYAML(cachedPath); err == nil {
				switch format {
				case "json":
					output.FormatJSON(os.Stdout, cached)
				default:
					output.FormatYAML(os.Stdout, cached)
				}
				fmt.Println("---")
			}
		}
	}
	return nil
}
