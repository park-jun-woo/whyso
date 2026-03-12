# whyso

Extract the **"why"** behind every file change from Claude Code sessions, and map your codebase for precise AI navigation.

`git blame` shows *who* changed *what* and *when*.
`whyso` shows **why** — the original user intent and AI's reasoning behind each change.

## How it works

```
whyso map    → tree-sitter parse → keyword map (functions, endpoints, rules, states)
whyso history → JSONL parse → parentUuid chain → per-file change history
```

## Install

```bash
go install github.com/park-jun-woo/whyso/cmd/whyso@latest
```

Or build from source:

```bash
git clone https://github.com/park-jun-woo/whyso.git
cd whyso
go build ./cmd/whyso/
```

Requires Go 1.22+ and a C compiler (CGO for tree-sitter).

## Usage

### Generate keyword map

```bash
# Current directory (stdout + .whyso/_map.md)
whyso map

# Specific path
whyso map internal/

# Custom output file
whyso map -o custom.md

# Force regeneration
whyso map -f
```

Example output:

```
# whyso/v1

## go
[parser]ParseSession,ExtractChanges,ListSessions
[history]BuildHistories,BuildIndex,FindUserRequest

## ssac
[service/gig]CreateGig,UpdateGig,PublishGig

## openapi
[api]CreateGig,UpdateGig,ListGigs
```

Supported languages: Go, TypeScript, JavaScript, Python, Rust, SSaC, OpenAPI, SQL, Rego, Gherkin, STML (HTML data-*), Mermaid stateDiagram.

### Show file change history

```bash
# Single file (stdout + .whyso/ cache)
whyso history README.md

# All files in a directory (.whyso/ cache only)
whyso history . --all

# Quiet mode (cache only, no stdout)
whyso history README.md -q

# Custom output directory
whyso history . --all --output custom-dir/

# JSON format
whyso history README.md --format json

# Clear cache and rebuild
whyso history README.md --reset
```

Example output:

```yaml
apiVersion: whyso/v1
file: CLAUDE.md
created: 2026-03-12T01:22:43Z
history:
  - timestamp: 2026-03-12T01:26:32Z
    session: 441b6643-d001-45df-811a-8ec138e73894
    user_request: "Add plan document rules to CLAUDE.md"
    answer: "Added specs/plans/ directory, plan-first workflow, and PhaseNNN naming convention."
    tool: Edit
    source: ~/.claude/projects/-home-user-project/441b6643.jsonl:79
```

### List sessions

```bash
whyso sessions
```

### Options

| Flag | Description |
|---|---|
| `-o <file>` | Map output file (default: `.whyso/_map.md`) |
| `-f, --force` | Force map regeneration (ignore mtime) |
| `--output <dir>` | History output directory (default: `.whyso/`) |
| `--format <yaml\|json>` | History output format (default: yaml) |
| `-q, --quiet` | Suppress stdout output |
| `--reset` | Clear history cache and rebuild |
| `--all` | Include all files in directory |
| `--sessions-dir <path>` | Override Claude Code sessions directory |

## Features

- **Keyword map** — tree-sitter powered extraction of functions, endpoints, queries, rules, states
- **User intent tracking** — traces `parentUuid` chain to find the original user request
- **AI answer extraction** — captures Claude's explanation of what it did
- **Grouped changes** — consecutive edits from the same request are merged
- **Subagent support** — includes changes made by subagent sessions
- **Incremental updates** — caches to `.whyso/`, only re-parses new sessions
- **Directory mirroring** — output structure mirrors source file paths

## License

MIT License — see [LICENSE](LICENSE).
