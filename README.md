# whylog

Extract the **"why"** behind every file change from Claude Code sessions.

`git blame` shows *who* changed *what* and *when*.
`whylog` shows **why** — the original user intent and AI's reasoning behind each change.

## How it works

Claude Code stores every conversation as JSONL files in `~/.claude/projects/`. whylog parses these sessions, extracts Write/Edit tool invocations, and traces the `parentUuid` chain back to the original user request.

```
~/.claude/projects/**/*.jsonl
  → whylog parse (JSONL parsing + tool_use extraction)
  → whylog history (parentUuid chain → per-file history)
  → stdout or file output
```

## Install

```bash
go install github.com/park-jun-woo/whylog/cmd/whylog@latest
```

Or build from source:

```bash
git clone https://github.com/park-jun-woo/whylog.git
cd whylog
go build ./cmd/whylog/
```

## Usage

### List sessions

```bash
whylog sessions
```

### Show file change history

```bash
# Single file
whylog history README.md

# All files in a directory
whylog history internal/ --all

# Output to directory (mirrors file structure)
whylog history . --all --output .file-histories/

# JSON format
whylog history README.md --format json
```

### Example output

```yaml
file: CLAUDE.md
created: 2026-03-12T01:22:43Z
history:
  - timestamp: 2026-03-12T01:26:32Z
    session: 441b6643-d001-45df-811a-8ec138e73894
    user_request: "Add plan document rules to CLAUDE.md"
    answer: "Added specs/plans/ directory, plan-first workflow, and PhaseNNN naming convention."
    tool: Edit
    source: ~/.claude/projects/-home-user-project/441b6643.jsonl:79
  - timestamp: 2026-03-12T01:32:09Z
    session: 441b6643-d001-45df-811a-8ec138e73894
    user_request: "Use Go CLI only, no frontend"
    answer: "Reduced tech stack to Go CLI, removed API/DB/UI/infra SSOT entries."
    tool: Edit
    sources:
      - ~/.claude/projects/-home-user-project/441b6643.jsonl:132
      - ~/.claude/projects/-home-user-project/441b6643.jsonl:135
```

### Options

| Flag | Description |
|---|---|
| `--sessions-dir <path>` | Override Claude Code sessions directory (default: auto-detect) |
| `--output <dir>` | Write to directory, mirroring file structure |
| `--format <yaml\|json>` | Output format (default: yaml) |
| `--all` | Include all files in directory |

## Features

- **User intent tracking** — traces `parentUuid` chain to find the original user request
- **AI answer extraction** — captures Claude's explanation of what it did
- **Grouped changes** — consecutive edits from the same request are merged
- **Subagent support** — includes changes made by subagent sessions
- **Incremental updates** — only re-parses modified sessions when using `--output`
- **Directory mirroring** — output structure mirrors source file paths

## Scope

whylog tracks all **text-based files** modified via Claude Code's Write/Edit tools. Binary files are excluded as they are not targets of these tools.

## License

MIT License — see [LICENSE](LICENSE).
