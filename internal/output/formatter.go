package output

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/clari/whyso/internal/history"
)

// FormatYAML writes a FileHistory as YAML to w.
func FormatYAML(w io.Writer, h *history.FileHistory) error {
	fmt.Fprintf(w, "file: %s\n", h.File)
	fmt.Fprintf(w, "created: %s\n", h.Created.Format(time.RFC3339))
	fmt.Fprintf(w, "history:\n")
	for _, e := range h.History {
		fmt.Fprintf(w, "  - timestamp: %s\n", e.Timestamp.Format(time.RFC3339))
		fmt.Fprintf(w, "    session: %s\n", e.Session)
		fmt.Fprintf(w, "    user_request: %q\n", e.UserRequest)
		if e.Answer != "" {
			fmt.Fprintf(w, "    answer: %q\n", e.Answer)
		}
		fmt.Fprintf(w, "    tool: %s\n", e.Tool)
		if e.Subagent {
			fmt.Fprintf(w, "    subagent: true\n")
		}
		if len(e.Sources) == 1 {
			fmt.Fprintf(w, "    source: %s:%d\n", e.Sources[0].File, e.Sources[0].Line)
		} else {
			fmt.Fprintf(w, "    sources:\n")
			for _, s := range e.Sources {
				fmt.Fprintf(w, "      - %s:%d\n", s.File, s.Line)
			}
		}
	}
	return nil
}

// FormatJSON writes a FileHistory as JSON to w.
func FormatJSON(w io.Writer, h *history.FileHistory) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(h)
}

// WriteHistories writes all histories to the output directory, mirroring file paths.
// If an output file already exists, it merges new entries with existing ones.
func WriteHistories(histories map[string]*history.FileHistory, outputDir, format string) error {
	for relPath, h := range histories {
		outPath := filepath.Join(outputDir, relPath+"."+format)

		// merge with existing file if present
		if format == "yaml" {
			existing, err := ReadYAML(outPath)
			if err == nil {
				h = history.Merge(existing, h)
			}
		}

		if err := os.MkdirAll(filepath.Dir(outPath), 0o755); err != nil {
			return err
		}

		f, err := os.Create(outPath)
		if err != nil {
			return err
		}

		var writeErr error
		switch format {
		case "json":
			writeErr = FormatJSON(f, h)
		default:
			writeErr = FormatYAML(f, h)
		}

		f.Close()
		if writeErr != nil {
			return writeErr
		}
	}
	return nil
}
