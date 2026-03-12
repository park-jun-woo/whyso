package history

import (
	"fmt"
	"sort"
)

// Merge combines existing and new histories, removing duplicates.
// Duplicate is defined by (timestamp, first source file:line).
func Merge(existing, incoming *FileHistory) *FileHistory {
	if existing == nil {
		return incoming
	}
	if incoming == nil {
		return existing
	}

	seen := make(map[string]bool)
	var merged []ChangeEntry

	for _, e := range existing.History {
		key := entryKey(e)
		if !seen[key] {
			seen[key] = true
			merged = append(merged, e)
		}
	}

	for _, e := range incoming.History {
		key := entryKey(e)
		if !seen[key] {
			seen[key] = true
			merged = append(merged, e)
		}
	}

	sort.Slice(merged, func(i, j int) bool {
		return merged[i].Timestamp.Before(merged[j].Timestamp)
	})

	result := &FileHistory{
		File:    incoming.File,
		History: merged,
	}
	if len(merged) > 0 {
		result.Created = merged[0].Timestamp
	}
	return result
}

func entryKey(e ChangeEntry) string {
	ts := e.Timestamp.Unix()
	if len(e.Sources) > 0 {
		return fmt.Sprintf("%d:%s:%d", ts, e.Sources[0].File, e.Sources[0].Line)
	}
	return fmt.Sprintf("%d:%s:%s", ts, e.Session, e.Tool)
}
