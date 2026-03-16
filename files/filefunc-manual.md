# filefunc — Manual for AI Agents

For Go application-layer projects: backend services, CLI tools, code generators, SSOT validators.

## How to Navigate

1. Read `codebook.yaml` — project vocabulary (required/optional keys and allowed values)
2. `filefunc chain func <target> --chon 2` — trace call relationships before modifying
3. `rg '//ff:func feature=validate'` — grep with codebook values to find files
4. Read `//ff:what` to narrow down — skip body if what is sufficient
5. Full read only the files you need, then work

---

## Rules

### File structure

| Rule | Description | Severity |
|---|---|---|
| F1 | One func per file (filename = func name, snake_case) | ERROR |
| F2 | One type per file (filename = type name, snake_case) | ERROR |
| F3 | One method per file (`{receiver}_{method}.go`) | ERROR |
| F4 | init() must not exist alone (requires var or func) | ERROR |
| F5 | _test.go files may have multiple funcs | exception |
| F6 | Semantically grouped consts allowed in one file | exception |

### Code quality

| Rule | Description | Severity |
|---|---|---|
| Q1 | Nesting depth: sequence=2, selection=2, iteration=dimension+1 | ERROR |
| Q2 | Func max 1000 lines | ERROR |
| Q3 | Func recommended max: sequence/iteration 100, selection 300 | WARNING |

### Annotation

| Rule | Description | Severity |
|---|---|---|
| A1 | Func files require `//ff:func`, type files require `//ff:type` | ERROR |
| A2 | Annotation values must exist in codebook | ERROR |
| A3 | Func/type files require `//ff:what` | ERROR |
| A6 | Annotations must be at the top of the file (above package) | ERROR |
| A7 | `//ff:checked` hash mismatch → body changed after LLM verification | ERROR |
| A8 | Required codebook keys must be present in annotation | ERROR |
| A9 | Func files must have `control=` (sequence/selection/iteration) | ERROR |
| A10 | `control=selection` but no switch at depth 1 | ERROR |
| A11 | `control=iteration` but no loop at depth 1 | ERROR |
| A12 | `control=sequence` but switch/loop exists at depth 1 | ERROR |
| A13 | `control=selection` but loop exists at depth 1 | ERROR |
| A14 | `control=iteration` but switch exists at depth 1 | ERROR |
| A15 | `control=iteration` requires `dimension=` | ERROR |
| A16 | `dimension=` value must be a positive integer | ERROR |

### Code quality (Q3 control-specific)

| control | Q3 limit |
|---|---|
| sequence | 100 lines |
| iteration | 100 lines |
| selection | 300 lines |

### Codebook format

| Rule | Description | Severity |
|---|---|---|
| C1 | `required` section must have at least one key with at least one value | ERROR |
| C2 | No duplicate keys within the same section | ERROR |
| C3 | All keys lowercase + hyphens only (`[a-z][a-z0-9-]*`) | ERROR |
| C4 | Required values should have non-empty descriptions | WARNING |

Codebook is validated first. If codebook fails, code validation does not run.

### Exceptions (not violations)

- const-only and var-only files do not require annotations
- If no `//ff:checked` exists in the project, A7 is skipped entirely

---

## Annotations

Write at the **very top** of every func/type file (above package declaration):

```go
//ff:func feature=validate type=rule control=sequence
//ff:what F1: validates one func per file
//ff:checked llm=gpt-oss:20b hash=a3f8c1d2    (auto by llmc)
package validate
```

`control=` is required for all func files (A9). Values: `sequence`, `selection` (switch), `iteration` (loop). Böhm-Jacopini (1966). 1 func 1 control — no mixing.

`dimension=` is required for `control=iteration` (A15). Q1 depth limit = dimension + 1. dimension=1 for flat lists (depth ≤ 2). dimension ≥ 2 requires named type (struct/interface) nesting — raw `[][][]int` is not allowed.

| Annotation | Required | Description |
|---|---|---|
| `//ff:func` | func files | Metadata (feature, type, control). Values from codebook.yaml + control rule |
| `//ff:type` | type files | Metadata (feature, type). Values from codebook.yaml |
| `//ff:what` | func/type files | One-line description. What does this do? |
| `//ff:why` | optional | Why designed this way? User decisions only |
| `//ff:checked` | auto (llmc) | LLM verification signature. Do not write manually |

### Control-based read strategy

```
control=selection  → read entire body at once. Don't read cases partially.
control=iteration  → focus on loop body. Outside loop is initialization.
control=sequence   → read only the step you need. Other steps: what is enough.
```

### Naming

- Filenames: `snake_case`
- Variables/functions: `camelCase`
- Types: `PascalCase`
- gofmt compliance, early return pattern

---

## Codebook

`codebook.yaml` must exist in the project root (next to `go.mod`). `required` keys must be in every annotation (A8). `optional` keys are used when relevant.

```yaml
required:
  feature:
    validate: "code structure rule validation (F1,Q1,A1 etc.)"
    parse: "source code, annotation, codebook parsing"
  type:
    command: "cobra command entrypoint"
    rule: "individual validation rule"

optional:
  pattern:
    error-collection: "collect errors for batch reporting"
  level:
    error: ""
```

Each value has a description (`key: "description"`). Used by `filefunc context` for LLM feature selection.

Amend codebook.yaml when new values are needed.

---

## Commands

```bash
filefunc validate                                    # current dir as project root
filefunc validate /path/to/project                   # explicit project root
filefunc validate --format json
filefunc chain func RunAll --chon 2                  # call relationships
filefunc chain func RunAll --chon 2 --meta what      # with //ff:what annotations
filefunc chain func RunAll --chon 2 --meta all       # with all annotations
filefunc chain func RunAll --chon 2 --meta what \
  --prompt "nesting depth 수정" --rate 0.8            # reranker filtering
filefunc chain feature validate                      # feature-wide chain
filefunc chain func RunAll --root /path/to/project   # explicit project root
filefunc context "nesting depth 수정"                   # LLM 4-stage context search
filefunc llmc                                        # LLM what-body verification
filefunc llmc /path/to/project
filefunc llmc --model qwen3:8b --threshold 0.9
```

Project root must contain `go.mod` and `codebook.yaml`. Omit to use current directory.

`--prompt` requires vLLM server: `pip install vllm && vllm serve Qwen/Qwen3-Reranker-0.6B --task score --hf_overrides '{"architectures":["Qwen3ForSequenceClassification"],"classifier_from_token":["no","yes"],"is_original_qwen3_reranker":true}'`

Exit code 1 on violations. Zero violations required before committing.

### .ffignore

Place in project root. Same syntax as `.gitignore`. Excludes paths from all commands.

```
vendor/
*.pb.go
*_gen.go
```

---

## Common Mistakes

| Mistake | Fix |
|---|---|
| Two funcs in one file | Extract helper functions into separate files |
| depth 3 (for→switch→if, for→if→if) | Type assertions + early continue, merge conditions, or extract func |
| Missing //ff:what | Write annotations first when creating a file |
| Value not in codebook | Check codebook.yaml first. Amend if absent |
| //ff:checked hash mismatch | Run `filefunc llmc` to re-verify |
| "codebook.yaml required" | Create codebook.yaml next to go.mod |
