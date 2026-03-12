package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/clari/whylog/internal/history"
	"github.com/clari/whylog/internal/output"
	"github.com/clari/whylog/internal/parser"
)

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	switch os.Args[1] {
	case "sessions":
		if err := runSessions(); err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}
	case "changes":
		if err := runChanges(); err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}
	case "history":
		if err := runHistory(); err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}
	default:
		fmt.Fprintf(os.Stderr, "unknown command: %s\n", os.Args[1])
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Fprintln(os.Stderr, "Usage: whylog <command>")
	fmt.Fprintln(os.Stderr, "")
	fmt.Fprintln(os.Stderr, "Commands:")
	fmt.Fprintln(os.Stderr, "  sessions    List sessions for the current project")
	fmt.Fprintln(os.Stderr, "  changes     List file changes across all sessions")
	fmt.Fprintln(os.Stderr, "  history     Show file change history")
	fmt.Fprintln(os.Stderr, "")
	fmt.Fprintln(os.Stderr, "History options:")
	fmt.Fprintln(os.Stderr, "  --format <yaml|json>     Output format (default: yaml)")
	fmt.Fprintln(os.Stderr, "  --output <dir>           Write to directory (mirrors file structure)")
	fmt.Fprintln(os.Stderr, "  --all                    All files in directory")
	fmt.Fprintln(os.Stderr, "")
	fmt.Fprintln(os.Stderr, "Options:")
	fmt.Fprintln(os.Stderr, "  --sessions-dir <path>    Override sessions directory")
}

func runSessions() error {
	sessionsDir, err := getSessionsDir()
	if err != nil {
		return err
	}

	sessions, err := parser.ListSessions(sessionsDir)
	if err != nil {
		return err
	}

	if len(sessions) == 0 {
		fmt.Println("No sessions found.")
		return nil
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintf(w, "TIMESTAMP\tSESSION ID\tFIRST MESSAGE\n")
	for _, s := range sessions {
		ts := s.Timestamp.Local().Format("2006-01-02 15:04")
		fmt.Fprintf(w, "%s\t%s\t%s\n", ts, s.ID, s.FirstMessage)
	}
	w.Flush()

	return nil
}

func runChanges() error {
	sessionsDir, err := getSessionsDir()
	if err != nil {
		return err
	}

	entries, err := os.ReadDir(sessionsDir)
	if err != nil {
		return err
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintf(w, "TIMESTAMP\tTOOL\tFILE\tUSER REQUEST\n")

	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".jsonl") {
			continue
		}
		path := filepath.Join(sessionsDir, entry.Name())
		records, err := parser.ParseSession(path)
		if err != nil {
			continue
		}

		changes := parser.ExtractChanges(records)
		idx := history.BuildIndex(records)

		for _, fc := range changes {
			ts := fc.Timestamp.Local().Format("2006-01-02 15:04")
			userReq := history.FindUserRequest(idx, fc.RecordUUID)
			if len(userReq) > 60 {
				userReq = userReq[:60] + "..."
			}
			fmt.Fprintf(w, "%s\t%s\t%s\t%s\n", ts, fc.Tool, fc.FilePath, userReq)
		}
	}
	w.Flush()
	return nil
}

func runHistory() error {
	sessionsDir, err := getSessionsDir()
	if err != nil {
		return err
	}

	projectRoot, err := os.Getwd()
	if err != nil {
		return err
	}

	// parse args after "history"
	var target, format, outputDir string
	var all bool
	format = "yaml"

	for i := 2; i < len(os.Args); i++ {
		switch os.Args[i] {
		case "--format":
			if i+1 < len(os.Args) {
				format = os.Args[i+1]
				i++
			}
		case "--output":
			if i+1 < len(os.Args) {
				outputDir = os.Args[i+1]
				i++
			}
		case "--all":
			all = true
		case "--sessions-dir":
			i++ // already handled
		default:
			if target == "" {
				target = os.Args[i]
			}
		}
	}

	if target == "" {
		return fmt.Errorf("usage: whylog history <file|dir> [--all] [--format yaml|json] [--output dir]")
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

	var histories map[string]*history.FileHistory
	var buildErr error

	if outputDir != "" {
		since := oldestOutputMtime(outputDir, format)
		if since.IsZero() {
			histories, buildErr = history.BuildHistories(sessionsDir, projectRoot, filter)
		} else {
			histories, buildErr = history.BuildHistoriesIncremental(sessionsDir, projectRoot, since, filter)
		}
	} else {
		histories, buildErr = history.BuildHistories(sessionsDir, projectRoot, filter)
	}
	if buildErr != nil {
		return buildErr
	}

	if len(histories) == 0 {
		if outputDir != "" {
			// incremental: no new changes, existing files are up to date
			return nil
		}
		fmt.Println("No history found.")
		return nil
	}

	if outputDir != "" {
		return output.WriteHistories(histories, outputDir, format)
	}

	// stdout
	for _, h := range histories {
		switch format {
		case "json":
			output.FormatJSON(os.Stdout, h)
		default:
			output.FormatYAML(os.Stdout, h)
		}
		fmt.Println("---")
	}
	return nil
}

// oldestOutputMtime returns the oldest mtime among output files in the directory.
// Using the oldest ensures no session is skipped — merge handles duplicates.
func oldestOutputMtime(dir, format string) time.Time {
	var oldest time.Time
	filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return nil
		}
		if !strings.HasSuffix(path, "."+format) {
			return nil
		}
		if oldest.IsZero() || info.ModTime().Before(oldest) {
			oldest = info.ModTime()
		}
		return nil
	})
	return oldest
}

func getSessionsDir() (string, error) {
	for i := 2; i < len(os.Args); i++ {
		if os.Args[i] == "--sessions-dir" && i+1 < len(os.Args) {
			return os.Args[i+1], nil
		}
	}
	cwd, err := os.Getwd()
	if err != nil {
		return "", err
	}
	return parser.DetectSessionsDir(cwd)
}
