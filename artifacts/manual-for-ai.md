# whyso — Manual for AI Agents

## What it does

`whyso` extracts the **why** behind file changes from Claude Code sessions. Use it before modifying code to understand prior intent.

## Commands

```bash
# View history of a single file
whyso history <file>

# View history of all files in a directory
whyso history <dir> --all

# Save to directory (enables incremental updates via mtime)
whyso history . --all --output .whyso/

# JSON output
whyso history <file> --format json

# List sessions for current project
whyso sessions
```

## Recommended workflow

1. `whyso history <file>` — read why the file was changed before editing
2. `Grep "functionName"` — locate the exact code
3. Edit with full context

## Output format (YAML)

```yaml
file: internal/parser/jsonl.go
created: 2026-03-12T01:22:43Z
history:
  - timestamp: 2026-03-12T01:26:32Z
    session: 441b6643-...
    user_request: "Add JSONL parser"
    answer: "Implemented streaming JSONL parser with..."
    tool: Write
    source: ~/.claude/projects/.../441b6643.jsonl:79
```

## Caching

When using `--output <dir>`, only new sessions are parsed on subsequent runs. Use `--output .whyso/` to maintain a persistent cache in the project root.
