package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/clari/whyso/internal/codemap"
	"github.com/clari/whyso/internal/history"
	"github.com/clari/whyso/internal/output"
	"github.com/clari/whyso/internal/parser"
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
	case "map":
		if err := runMap(); err != nil {
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
	var all, quiet, reset bool
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
		case "--quiet", "-q":
			quiet = true
		case "--all":
			all = true
		case "--reset":
			reset = true
		case "--sessions-dir":
			i++ // already handled
		default:
			if target == "" {
				target = os.Args[i]
			}
		}
	}

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

	since := oldestOutputMtime(outputDir, format)
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

	if len(histories) == 0 {
		return nil
	}

	// always write cache
	if err := output.WriteHistories(histories, outputDir, format); err != nil {
		return err
	}

	// stdout: single file only, unless -q or --output suppresses
	if !quiet && outputDir == filepath.Join(projectRoot, ".whyso") && !targetInfo.IsDir() {
		for _, h := range histories {
			switch format {
			case "json":
				output.FormatJSON(os.Stdout, h)
			default:
				output.FormatYAML(os.Stdout, h)
			}
			fmt.Println("---")
		}
	}
	return nil
}

func runMap() error {
	var target, outputFile string
	var force bool

	for i := 2; i < len(os.Args); i++ {
		switch os.Args[i] {
		case "-o":
			if i+1 < len(os.Args) {
				outputFile = os.Args[i+1]
				i++
			}
		case "-f", "--force":
			force = true
		default:
			if target == "" {
				target = os.Args[i]
			}
		}
	}

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

func clearCache(dir, format string) {
	filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return nil
		}
		if strings.HasSuffix(path, "."+format) {
			os.Remove(path)
		}
		return nil
	})
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
