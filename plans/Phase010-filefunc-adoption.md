# 완료
# Phase010: filefunc 도입

## 목표

whyso 코드베이스에 filefunc 규칙을 적용한다. 파일 분리(1파일 1func/1type), 어노테이션, codebook 설정을 완료하고 `filefunc validate` 위반 0건을 달성한다.

## 현황 분석

| 항목 | 현재 |
|---|---|
| Go 파일 수 | 12 |
| //ff: 어노테이션 | 0개 |
| codebook.yaml | 없음 |
| 1func/1file 위반 | 11/12 파일 |
| 1type/1file 위반 | 5/12 파일 |
| 총 func 수 | ~70 |
| 총 type 수 | ~15 |

### 파일별 위반 상세

| 파일 | func 수 | type 수 | 비고 |
|---|---|---|---|
| cmd/whyso/main.go | 9 | 0 | main + 8 핸들러 (411줄) |
| internal/output/formatter.go | 3 | 0 | + const Version |
| internal/output/reader.go | 4 | 0 | |
| pkg/codemap/codemap.go | 10 | 1 (Section) | 253줄, + const Version |
| pkg/codemap/treesitter.go | 15 | 0 | 190줄 |
| pkg/codemap/regex.go | 4 | 0 | + 5개 모듈 var |
| pkg/history/builder.go | 4 | 3 (FileHistory, Source, ChangeEntry) | |
| pkg/history/chain.go | 3 | 1 (RecordIndex) | + const maxChainDepth |
| pkg/history/merge.go | 2 | 0 | |
| pkg/model/record.go | 0 func, 3 method | 4 (Record, Message, ContentBlock, SessionInfo) | |
| pkg/parser/extract.go | 3 | 3 (FileChange, writeInput, editInput) | |
| pkg/parser/jsonl.go | 5 | 0 | + const maxLineSize |

## control= 및 depth 분석

filefunc 규칙:
- A10-A14: depth 1에 switch/loop 존재 여부로 control= 결정
- Q1: nesting depth ≤ 2 (ERROR)
- filepath.Walk 콜백은 익명 함수이므로 enclosing func의 depth에 포함

### control= 할당표

분리 후 각 func 파일의 control= 값. A10-A14 기준으로 depth 1에 switch/loop 유무를 검사한다.

| 함수 | depth 1 구조 | control= | Q1 depth | 비고 |
|---|---|---|---|---|
| **cmd/whyso** | | | | |
| main | switch | selection | 2 | ✓ |
| printUsage | 없음 | sequence | 1 | ✓ |
| runSessions | for (sessions 순회) | iteration | 2 | ✓ |
| runChanges | for (entries 순회) | iteration | 3 ✗ | for→for→if |
| runHistory | for (arg 파싱) | iteration | 3 ✗ | for→switch→if. 추출 후 sequence |
| runMap | for (arg 파싱) | iteration | 3 ✗ | for→switch→if. 추출 후 sequence |
| oldestOutputMtime | 없음 (Walk는 함수 호출) | sequence | 2 | ✓ |
| clearCache | 없음 (Walk는 함수 호출) | sequence | 2 | ✓ |
| getSessionsDir | for (arg 검색) | iteration | 2 | ✓ |
| **internal/output** | | | | |
| FormatYAML | for (history 순회) | iteration | 3 ✗ | for→else→for(sources) |
| FormatJSON | 없음 | sequence | 1 | ✓ |
| WriteHistories | for (histories 순회) | iteration | 3 ✗ | for→if(yaml)→if(err==nil) |
| ReadYAML | for (scanner) | iteration | 3 ✗ | for→if(timestamp)→if(current!=nil) |
| parseSource | 없음 | sequence | 1 | ✓ |
| unquote | 없음 | sequence | 1 | ✓ |
| OutputPath | 없음 | sequence | 1 | ✓ |
| **pkg/codemap** | | | | |
| NeedsUpdate | 없음 (Walk는 함수 호출) | sequence | 3 ✗ | Walk 콜백 내 if→if |
| BuildMap | 없음 (Walk는 함수 호출) | sequence | 3 ✗ | Walk 콜백 내 if→if |
| FormatMap | for (sections 순회) | iteration | 2 | ✓ |
| buildSections | for (order 순회) | iteration | 3 ✗ | for→for→if |
| shouldSkipDir | for (skip 순회) | iteration | 2 | ✓ |
| detectParser | switch (ext) | selection | 2 | ✓ |
| detectGroup | switch (lang) | selection | 1 | ✓ |
| extractGoPackage | for (lines 순회) | iteration | 3 ✗ | for→if→if |
| sortedKeys | for (keys 순회) | iteration | 1 | ✓ |
| dedupe | for (ss 순회) | iteration | 2 | ✓ |
| runQuery | for (matches 순회) | iteration | 3 ✗ | for→for→if |
| parseGo~parseRust | 없음 (runQuery 위임) | sequence | 1 | ✓ |
| parseRego | 없음 (parseRegoRegex 위임) | sequence | 1 | ✓ |
| parseOpenAPI | 없음 (위임) | sequence | 1 | ✓ |
| parseSTML | 없음 (위임) | sequence | 1 | ✓ |
| parseOpenAPIFromYAML | 없음 | sequence | 1 | ✓ |
| findOperationIDs | for (children 순회) | iteration | 3 ✗ | for→if→if |
| parseSTMLFromHTML | 없음 | sequence | 1 | ✓ |
| findDataAttributes | for (children 순회) | iteration | 3 ✗ | for 내부 if→for→switch |
| trimQuotes | 없음 | sequence | 1 | ✓ |
| parseGherkin | for (matches) | iteration | 2 | ✓ |
| parseMermaid | for (matches) | iteration | 2 | 2개 for loop ✓ |
| parseSQL | for (matches) | iteration | 2 | ✓ |
| parseRegoRegex | for (matches) | iteration | 2 | ✓ |
| **pkg/history** | | | | |
| BuildHistories | 없음 (위임) | sequence | 1 | ✓ |
| BuildHistoriesIncremental | 없음 (위임) | sequence | 1 | ✓ |
| buildHistories | for (entries 순회) | iteration | 3 ✗ | for 3중첩 |
| lastEntry | 없음 | sequence | 1 | ✓ |
| BuildIndex | for (records) | iteration | 1 | ✓ |
| FindUserRequest | for (chain 추적) | iteration | 3 ✗ | for→if→if |
| FindAnswer | for (records) | iteration | 3 ✗ | for→for→if |
| Merge | for (existing 순회) | iteration | 2 | 2개 for loop ✓ |
| entryKey | 없음 | sequence | 1 | ✓ |
| **pkg/model** | | | | |
| Record.IsUserMessage | 없음 | sequence | 1 | ✓ |
| Record.UserContent | 없음 | sequence | 1 | ✓ |
| Record.ContentBlocks | 없음 | sequence | 1 | ✓ |
| **pkg/parser** | | | | |
| ParseSession | for (subEntries) | iteration | 2 | ✓ |
| parseJSONL | for (scanner) | iteration | 2 | ✓ |
| DetectSessionsDir | 없음 | sequence | 1 | ✓ |
| toSlug | 없음 | sequence | 1 | ✓ |
| ListSessions | for (entries) | iteration | 3 ✗ | for→for→if |
| ExtractChanges | for (records) | iteration | 3 ✗ | for→for→switch |
| parseWrite | 없음 | sequence | 1 | ✓ |
| parseEdit | 없음 | sequence | 1 | ✓ |

### Q1 depth 3 위반 목록 (총 18건)

| # | 함수 | depth 경로 | 해결 방법 | 새 파일 |
|---|---|---|---|---|
| 1 | runChanges | for→for→if | inner loop body를 func 추출 | +1 |
| 2 | runHistory | for→switch→if | arg 파싱을 func 추출, if/continue 패턴으로 변경 | +1 |
| 3 | runMap | for→switch→if | arg 파싱을 func 추출, if/continue 패턴으로 변경 | +1 |
| 4 | FormatYAML | for→else→for | sources 포맷을 func 추출 | +1 |
| 5 | WriteHistories | for→if→if | 조건 병합: `existing, _ := ReadYAML(); if yaml && existing != nil` | 0 |
| 6 | ReadYAML | for→if→if | 조건 병합: `if isTimestamp && current != nil` 로 분리 | 0 |
| 7 | NeedsUpdate | Walk 콜백 if→if | 조건 병합: `if IsDir() && shouldSkipDir() { SkipDir }; if IsDir() { nil }` | 0 |
| 8 | BuildMap | Walk 콜백 if→if | 위와 동일 패턴 | 0 |
| 9 | buildSections | for→for→if | map 룩업으로 inner for 제거 | 0 |
| 10 | extractGoPackage | for→if→if | 조건 병합 또는 early continue | 0 |
| 11 | runQuery | for→for→if | early continue: `if name == "" { continue }` | 0 |
| 12 | findOperationIDs | for→if→if | early continue: `if type != pair { recurse; continue }` | 0 |
| 13 | findDataAttributes | for→for→switch | attribute 처리를 func 추출 | +1 |
| 14 | buildHistories | for→for→for | changes 처리를 func 추출 | +1 |
| 15 | FindUserRequest | for→if→if | 조건 병합: `parentUUID == nil && isUserMessage` | 0 |
| 16 | FindAnswer | for→for→if | blocks 텍스트 추출을 func 추출 | +1 |
| 17 | ListSessions | for→for→if | 첫 사용자 메시지 찾기를 func 추출 | +1 |
| 18 | ExtractChanges | for→for→switch | block 처리를 func 추출 | +1 |

조건 병합/early continue로 해결: 9건 (새 파일 불필요)
helper func 추출 필요: 9건 (+9 파일)

## 구현 단계

### Step 1. codebook.yaml 생성

프로젝트 루트에 `codebook.yaml` 작성. whyso 도메인에 맞는 키와 값을 정의한다.

```yaml
required:
  feature: [session, change, history, codemap, output, cli]
  type: [command, parser, builder, merger, formatter, reader, model, util]

optional:
  pattern: [file-visitor, incremental-update, parent-chain]
```

### Step 2. .ffignore 생성

```
vendor/
*.pb.go
*_gen.go
```

### Step 3. 파일 분리 — pkg/model/

record.go (4 type + 3 method) → 7파일:

| 새 파일 | 내용 | control= |
|---|---|---|
| record.go | Record struct | (type) |
| record_is_user_message.go | IsUserMessage method | sequence |
| record_user_content.go | UserContent method | sequence |
| record_content_blocks.go | ContentBlocks method | sequence |
| message.go | Message struct | (type) |
| content_block.go | ContentBlock struct | (type) |
| session_info.go | SessionInfo struct | (type) |

주의: Record가 Message를 필드로 참조, ContentBlocks()가 ContentBlock을 반환 — 같은 패키지이므로 문제없음.

### Step 4. 파일 분리 — pkg/parser/

**jsonl.go** (5 func + 1 const) → 5파일:

| 새 파일 | 내용 | control= |
|---|---|---|
| parse_session.go | ParseSession | iteration |
| parse_jsonl.go | parseJSONL + const maxLineSize | iteration |
| detect_sessions_dir.go | DetectSessionsDir | sequence |
| to_slug.go | toSlug | sequence |
| list_sessions.go | ListSessions | iteration |

**extract.go** (3 func + 3 type) → 6파일:

| 새 파일 | 내용 | control= |
|---|---|---|
| file_change.go | FileChange type | (type) |
| write_input.go | writeInput type | (type) |
| edit_input.go | editInput type | (type) |
| extract_changes.go | ExtractChanges | iteration |
| parse_write.go | parseWrite | sequence |
| parse_edit.go | parseEdit | sequence |

리팩토링:
- `ExtractChanges` depth 3: block 처리를 `extractBlockChange` func으로 추출
- `ListSessions` depth 3: 첫 사용자 메시지 찾기를 `findFirstUserMessage` func으로 추출

→ 추가: `extract_block_change.go` (selection), `find_first_user_message.go` (iteration)

### Step 5. 파일 분리 — pkg/history/

**builder.go** (4 func + 3 type) → 7파일:

| 새 파일 | 내용 | control= |
|---|---|---|
| file_history.go | FileHistory type | (type) |
| source.go | Source type | (type) |
| change_entry.go | ChangeEntry type | (type) |
| build_histories.go | BuildHistories (1줄 위임) | sequence |
| build_histories_incremental.go | BuildHistoriesIncremental (1줄 위임) | sequence |
| build_histories_core.go | buildHistories (private core) | iteration |
| last_entry.go | lastEntry | sequence |

리팩토링: `buildHistories` depth 3 → changes 처리를 `processChanges` func으로 추출.

→ 추가: `process_changes.go` (iteration)

**chain.go** (3 func + 1 type + 1 const) → 4파일:

| 새 파일 | 내용 | control= |
|---|---|---|
| record_index.go | RecordIndex type | (type) |
| build_index.go | BuildIndex | iteration |
| find_user_request.go | FindUserRequest + const maxChainDepth | iteration |
| find_answer.go | FindAnswer | iteration |

리팩토링:
- `FindUserRequest` depth 3: 조건 병합으로 해결. 새 파일 불필요.
- `FindAnswer` depth 3: blocks 텍스트 추출을 `extractTextBlocks` func으로 추출.

→ 추가: `extract_text_blocks.go` (iteration)

**merge.go** (2 func) → 2파일:

| 새 파일 | 내용 | control= |
|---|---|---|
| merge.go | Merge | iteration |
| entry_key.go | entryKey | sequence |

### Step 6. 파일 분리 — pkg/codemap/

**codemap.go** (10 func + 1 type + 1 const) → 12파일:

| 새 파일 | 내용 | control= |
|---|---|---|
| section.go | Section type | (type) |
| version.go | const Version | (const 전용, 어노테이션 불필요) |
| needs_update.go | NeedsUpdate | sequence |
| build_map.go | BuildMap | sequence |
| format_map.go | FormatMap | iteration |
| build_sections.go | buildSections | iteration |
| should_skip_dir.go | shouldSkipDir | iteration |
| detect_parser.go | detectParser | selection |
| detect_group.go | detectGroup | selection |
| extract_go_package.go | extractGoPackage | iteration |
| sorted_keys.go | sortedKeys | iteration |
| dedupe.go | dedupe | iteration |

리팩토링 (조건 병합, 새 파일 불필요):
- `NeedsUpdate` depth 3: `if IsDir() && shouldSkipDir() { SkipDir }; if IsDir() { nil }`
- `BuildMap` depth 3: 위와 동일
- `buildSections` depth 3: inner for(order 검색)를 map 룩업으로 대체
- `extractGoPackage` depth 3: early continue + 조건 분리

**treesitter.go** (15 func) → 15파일:

| 새 파일 | 내용 | control= |
|---|---|---|
| run_query.go | runQuery | iteration |
| parse_go.go | parseGo | sequence |
| parse_ssac.go | parseSSaC | sequence |
| parse_typescript.go | parseTypeScript | sequence |
| parse_javascript.go | parseJavaScript | sequence |
| parse_python.go | parsePython | sequence |
| parse_rust.go | parseRust | sequence |
| parse_rego.go | parseRego | sequence |
| parse_openapi.go | parseOpenAPI | sequence |
| parse_stml.go | parseSTML | sequence |
| parse_openapi_from_yaml.go | parseOpenAPIFromYAML | sequence |
| find_operation_ids.go | findOperationIDs | iteration |
| parse_stml_from_html.go | parseSTMLFromHTML | sequence |
| find_data_attributes.go | findDataAttributes | iteration |
| trim_quotes.go | trimQuotes | sequence |

리팩토링:
- `runQuery` depth 3: early continue `if name == "" { continue }`. 새 파일 불필요.
- `findOperationIDs` depth 3: early continue. 새 파일 불필요.
- `findDataAttributes` depth 3: attribute 처리를 `extractAttribute` func으로 추출.

→ 추가: `extract_attribute.go` (iteration)

**regex.go** (4 func + 5 var) → 5파일:

| 새 파일 | 내용 | control= |
|---|---|---|
| regex_vars.go | 5개 regexp 변수 | (var 전용, 어노테이션 불필요) |
| parse_gherkin.go | parseGherkin | iteration |
| parse_mermaid.go | parseMermaid | iteration |
| parse_sql.go | parseSQL | iteration |
| parse_rego_regex.go | parseRegoRegex | iteration |

### Step 7. 파일 분리 — internal/output/

**formatter.go** (3 func + 1 const) → 4파일:

| 새 파일 | 내용 | control= |
|---|---|---|
| version.go | const Version | (const 전용, 어노테이션 불필요) |
| format_yaml.go | FormatYAML | iteration |
| format_json.go | FormatJSON | sequence |
| write_histories.go | WriteHistories | iteration |

리팩토링:
- `FormatYAML` depth 3 (for→else→for sources): sources 포맷을 `formatSources` func으로 추출.
- `WriteHistories` depth 3 (for→if→if): 조건 병합으로 해결 — `existing, _ := ReadYAML(outPath); if format == "yaml" && existing != nil { h = Merge(existing, h) }`. 새 파일 불필요.

→ 추가: `format_sources.go` (iteration)

**reader.go** (4 func) → 4파일:

| 새 파일 | 내용 | control= |
|---|---|---|
| read_yaml.go | ReadYAML | iteration |
| parse_source.go | parseSource | sequence |
| unquote.go | unquote | sequence |
| output_path.go | OutputPath | sequence |

리팩토링: `ReadYAML` depth 3 (for→if→if): 조건 병합으로 해결 — `isNewEntry := HasPrefix("  - timestamp: "); if isNewEntry && current != nil { append }; if isNewEntry { current = ... }`. 새 파일 불필요.

### Step 8. 파일 분리 — cmd/whyso/

**main.go** (9 func, 411줄) → 9파일:

| 새 파일 | 내용 | control= |
|---|---|---|
| main.go | main (entry point) | selection |
| print_usage.go | printUsage | sequence |
| run_sessions.go | runSessions | iteration |
| run_changes.go | runChanges | iteration |
| run_history.go | runHistory → 리팩토링 후 | sequence |
| run_map.go | runMap → 리팩토링 후 | sequence |
| oldest_output_mtime.go | oldestOutputMtime | sequence |
| clear_cache.go | clearCache | sequence |
| get_sessions_dir.go | getSessionsDir | iteration |

리팩토링:
- `runHistory` depth 3 (for→switch→if): arg 파싱을 `parseHistoryArgs` func으로 추출. 추출된 func은 switch 대신 if/continue 패턴으로 작성하여 depth 2 유지. runHistory는 sequence.
- `runMap` depth 3: 위와 동일. `parseMapArgs` func 추출, if/continue 패턴.
- `runChanges` depth 3 (for→for→if): inner loop body를 `formatChangeRow` func으로 추출.

→ 추가: `parse_history_args.go` (iteration), `parse_map_args.go` (iteration), `format_change_row.go` (sequence)

runHistory는 arg 파싱 추출 후 ~127줄 (Q3 sequence 권장 100줄 WARNING, ERROR 아님).

### Step 9. 어노테이션 추가

모든 분리된 파일에 `//ff:func` 또는 `//ff:type` + `//ff:what` 어노테이션을 작성한다.

func 파일 예시:
```go
//ff:func feature=session type=parser control=sequence
//ff:what JSONL 세션 파일을 파싱하여 Record 슬라이스로 변환
package parser
```

type 파일 예시:
```go
//ff:type feature=session type=model
//ff:what JSONL 한 줄에 대응하는 레코드
package model
```

method 파일 예시:
```go
//ff:func feature=session type=model control=sequence
//ff:what Record가 사용자 메시지인지 판별
package model
```

### Step 10. 검증

```bash
filefunc validate
go build ./...
```

위반 0건 + 빌드 성공 확인. 위반이 있으면 수정 후 재검증.

## 분리 결과 예상

| 패키지 | 현재 파일 수 | 분리 파일 수 | 리팩토링 추가 | 합계 |
|---|---|---|---|---|
| cmd/whyso | 1 | 9 | +3 | 12 |
| internal/output | 2 | 8 | +1 | 9 |
| pkg/codemap | 3 | 32 | +1 | 33 |
| pkg/history | 3 | 13 | +2 | 15 |
| pkg/model | 1 | 7 | 0 | 7 |
| pkg/parser | 2 | 11 | +2 | 13 |
| **합계** | **12** | **80** | **+9** | **89** |

## 작업 순서 요약

1. codebook.yaml + .ffignore 생성
2. 패키지별 파일 분리 (Step 3~8) — 각 패키지마다 분리 → `go build ./...` 확인
3. Q1 depth 3 리팩토링 — 조건 병합/early continue 9건 + helper func 추출 9건
4. 어노테이션 추가 (Step 9)
5. `filefunc validate` + `go build ./...` 최종 확인
