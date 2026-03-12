# whyso — Manual for AI Agents

## What it does

`whyso` gives AI agents two things before touching code:

- **`whyso map`** — what exists (functions, endpoints, rules, states)
- **`whyso history`** — why it was changed (user intent + AI reasoning)

## Commands

```bash
# Generate keyword map (stdout + .whyso/_map.md)
whyso map [path]

# Save map to custom file
whyso map [path] -o custom.md

# Force map regeneration
whyso map -f

# Single file history (stdout + .whyso/ cache)
whyso history <file>

# Directory history (.whyso/ cache only)
whyso history <dir> --all

# Quiet mode (cache only, no stdout)
whyso history <file> -q

# JSON output
whyso history <file> --format json

# Clear history cache and rebuild
whyso history <file> --reset

# List sessions
whyso sessions
```

## Recommended workflow

1. `Read .whyso/_map.md` — scan all keywords (functions, endpoints, queries, states)
2. `whyso history <file>` — read why the file was changed before editing
3. `Grep "keyword"` — locate exact code with precise keyword from the map
4. Edit with full context

## Map output format

```
# whyso/v1

## go
[parser]ParseSession,ExtractChanges,ListSessions

## ssac
[service/gig]CreateGig,UpdateGig,PublishGig

## openapi
[api]CreateGig,UpdateGig,ListGigs
```

Supported: Go, TypeScript, JavaScript, Python, Rust, SSaC, OpenAPI, SQL, Rego, Gherkin, STML, Mermaid.

## History output format (YAML)

```yaml
apiVersion: whyso/v1
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

All output is cached to `.whyso/` by default:

- `.whyso/_map.md` — keyword map
- `.whyso/*.yaml` — file history (incremental, mtime-based)

Subsequent runs only parse new sessions.
