# 완료
# Phase011: filefunc dimension= 속성 추가

## 목표

filefunc 매뉴얼 업데이트(A15, A16)에 따라 모든 `control=iteration` 어노테이션에 `dimension=` 속성을 추가한다.

## 배경

filefunc 규칙 변경:
- A15: `control=iteration` requires `dimension=` (NEW)
- A16: `dimension=` value must be a positive integer (NEW)
- Q1: iteration의 depth 제한이 `dimension+1`로 변경 (기존: 일률 ≤ 2)

dimension=1 → depth ≤ 2 (flat list 순회)
dimension=2 → depth ≤ 3 (named type 중첩 필요)

## 현황

- `control=iteration` 파일: 35개
- Phase010에서 모든 iteration func을 depth ≤ 2로 리팩토링 완료
- 따라서 전부 `dimension=1`

## 대상 파일 (35개)

### cmd/whyso (5개)
- run_sessions.go
- run_changes.go
- get_sessions_dir.go
- parse_history_args.go
- parse_map_args.go

### internal/output (4개)
- format_yaml.go
- write_histories.go
- read_yaml.go
- format_sources.go

### pkg/codemap (14개)
- format_map.go
- build_sections.go
- should_skip_dir.go
- extract_go_package.go
- sorted_keys.go
- dedupe.go
- run_query.go
- find_operation_ids.go
- find_data_attributes.go
- extract_attribute.go
- parse_gherkin.go
- parse_mermaid.go
- parse_sql.go
- parse_rego_regex.go

### pkg/history (7개)
- build_histories_core.go
- build_index.go
- find_user_request.go
- find_answer.go
- extract_text_blocks.go
- merge.go
- process_changes.go

### pkg/parser (5개)
- parse_session.go
- parse_jsonl.go
- list_sessions.go
- find_first_user_message.go
- extract_changes.go

## 구현

각 파일의 1행 `control=iteration`을 `control=iteration dimension=1`로 변경.

변경 전: `//ff:func feature=xxx type=xxx control=iteration`
변경 후: `//ff:func feature=xxx type=xxx control=iteration dimension=1`

## 검증

```bash
filefunc validate
go build ./...
```
