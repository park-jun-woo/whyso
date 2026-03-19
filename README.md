# whyso

See why Claude actually changed that file — intent extracted from session JSONL.

## Why I built this

You use Claude Code, commits pile up. Weeks later you revisit the code — "why is this structured like this?" You run `git blame`. Commit message: `Refactor parser logic`. OK, but *why*?

The answer was already on my machine.

I opened a Claude Code session log (JSONL) and every file change was chained:

```
User: "Fix the bug where parser ignores empty sessions"
  → Claude: "Empty session filtering was missing, adding early return to ParseSession"
    → tool_use: Edit parser.go
```

Request → reasoning → execution. All recorded.

`whyso` traces this chain and links every file change to the **original request** and **Claude's reasoning**.

That refactoring you asked Claude to do 3 months ago — why it ended up that way — you can see it right now:

```bash
go install github.com/park-jun-woo/whyso/cmd/whyso@latest
whyso history CLAUDE.md   # < 1s, even for month-long projects
```

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

- **User intent tracking** — traces `parentUuid` chain to find the original user request
- **AI reasoning** — captures why Claude made each change, in its own words
- **Keyword map** — tree-sitter powered extraction of functions, endpoints, queries, rules, states
- **Grouped changes** — consecutive edits from the same request are merged
- **Subagent support** — includes changes made by subagent sessions
- **Incremental updates** — caches to `.whyso/`, only re-parses new sessions
- **Directory mirroring** — output structure mirrors source file paths

## License

MIT License — see [LICENSE](LICENSE).
