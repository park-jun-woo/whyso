package parser

import (
	"encoding/json"
	"time"

	"github.com/clari/whyso/internal/model"
)

// FileChange represents a single file modification extracted from a tool_use.
type FileChange struct {
	FilePath   string
	Tool       string // "Write" | "Edit"
	ToolUseID  string
	OldString   string // Edit only
	NewString   string // Edit only
	Content     string // Write only
	RecordUUID  string
	RecordIndex int    // index in records slice (for FindAnswer)
	Timestamp  time.Time
	SessionID  string
	SourceFile string // JSONL file path
	SourceLine int    // line number in JSONL
}

type writeInput struct {
	FilePath string `json:"file_path"`
	Content  string `json:"content"`
}

type editInput struct {
	FilePath  string `json:"file_path"`
	OldString string `json:"old_string"`
	NewString string `json:"new_string"`
}

// ExtractChanges extracts file changes from Write/Edit tool_use blocks.
func ExtractChanges(records []model.Record) []FileChange {
	var changes []FileChange

	for i := range records {
		rec := &records[i]
		if rec.Type != "assistant" {
			continue
		}

		blocks := rec.ContentBlocks()
		for _, block := range blocks {
			if block.Type != "tool_use" {
				continue
			}

			var fc *FileChange
			switch block.Name {
			case "Write":
				fc = parseWrite(block)
			case "Edit":
				fc = parseEdit(block)
			default:
				continue
			}

			if fc == nil {
				continue
			}
			fc.RecordUUID = rec.UUID
			fc.RecordIndex = i
			fc.Timestamp = rec.Timestamp
			fc.SessionID = rec.SessionID
			fc.SourceFile = rec.SourceFile
			fc.SourceLine = rec.SourceLine
			changes = append(changes, *fc)
		}
	}

	return changes
}


func parseWrite(block model.ContentBlock) *FileChange {
	var input writeInput
	if err := json.Unmarshal(block.Input, &input); err != nil {
		return nil
	}
	return &FileChange{
		FilePath: input.FilePath,
		Tool:     "Write",
		ToolUseID: block.ID,
		Content:   input.Content,
	}
}

func parseEdit(block model.ContentBlock) *FileChange {
	var input editInput
	if err := json.Unmarshal(block.Input, &input); err != nil {
		return nil
	}
	return &FileChange{
		FilePath: input.FilePath,
		Tool:     "Edit",
		ToolUseID: block.ID,
		OldString: input.OldString,
		NewString: input.NewString,
	}
}
