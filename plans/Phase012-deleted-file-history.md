# 완료
# Phase012: 삭제된 파일 이력 조회 및 이동 추정

## 목표

1. 삭제된 파일도 경로만으로 이력 제공
2. git rename 감지로 `moved_from` 표시
3. Bash cp/mv 명령을 참고 힌트로 표시 (이력 병합 안 함)

## 배경

파일이 `cp`/`mv`로 이동되거나 삭제된 경우, 현재 `whyso history`는 `os.Stat` 실패로 에러를 반환한다. 세션 기록은 원본 경로로 존재하지만 접근할 수 없다.

## 단계

### Step 0: 의존성 사전 검증

`whyso history` 실행 초기에 필수 의존성을 검사하고, 누락 시 명확한 에러 메시지로 안내.

- git 저장소 여부: `git rev-parse --git-dir` 실행 → 실패 시 `"not a git repository. whyso requires a git project root"`
- Claude Code 세션 디렉토리: `~/.claude/` 존재 여부 → 없으면 `"~/.claude/ not found. install Claude Code first: https://claude.ai/download"`
- 기존 `getSessionsDir()` 에러보다 선행하여 사용자 친화적 메시지 제공

**영향 파일:**
- 신규 `cmd/whyso/check_deps.go` — 의존성 검증 함수
- `cmd/whyso/run_history.go` — 함수 초기에 호출

### Step 1: 삭제된 파일 경로 지원

`run_history.go`에서 `os.Stat` 실패 시 에러 대신 경로 기반 검색 분기.

- `os.Stat` 실패 → 파일 미존재 모드 (디렉토리 모드 불가)
- `absTarget`을 상대경로로 변환하여 filter 직접 구성
- `makeFilter`, `resolveSince`의 `targetInfo` 의존 제거 (nil 허용 분기)
- `printHistoryOutput`도 `targetInfo` 없이 동작하도록 수정

**영향 파일:**
- `cmd/whyso/run_history.go` — os.Stat 분기
- `cmd/whyso/make_filter.go` — targetInfo nil 처리
- `cmd/whyso/resolve_since.go` — targetInfo nil 처리

### Step 2: git rename 감지 및 moved_from 표시

이력 빌드 후, 결과가 비어 있으면 git log --follow로 이전 경로를 탐색.

- 신규 `pkg/history/detect_rename.go` — `git log --follow --diff-filter=R --summary -- <path>` 실행, 이전 경로 반환
- `FileHistory`에 `MovedFrom string` 필드 추가
- `run_history.go`에서 이력 빈 경우 detect_rename 호출 → 이전 경로로 재검색 → `MovedFrom` 설정
- `format_yaml.go`에서 `moved_from` 필드 출력

**영향 파일:**
- 신규 `pkg/history/detect_rename.go`
- `pkg/history/file_history.go` — MovedFrom 필드
- `cmd/whyso/run_history.go` — rename 감지 분기
- `internal/output/format_yaml.go` — moved_from 출력

### Step 3: Bash cp/mv 참고 힌트

세션 JSONL에서 Bash tool_use의 command 필드를 파싱하여, 대상 파일 경로가 포함된 cp/mv 명령을 감지. 이력 병합은 하지 않고 별도 `hints` 필드로 표시.

- 신규 `pkg/parser/parse_bash_hint.go` — Bash tool_use에서 `cp`/`mv` 패턴 매칭, 단순 정규식 (`^(cp|mv)\s+.+\s+.+`)
- 신규 `pkg/parser/bash_hint.go` — BashHint 타입 정의
- `extract_changes.go` 또는 별도 함수에서 BashHint 수집
- `FileHistory`에 `Hints []BashHint` 필드 추가
- `format_yaml.go`에서 `hints` 섹션 출력

**힌트 출력 예시:**
```yaml
hints:
  - timestamp: 2026-03-23T00:30:00Z
    session: f126f7f9-...
    command: "mv internal/analyzer/collect_chain.go pkg/analyzer/"
    source: ...f126f7f9.jsonl:1500
```

**제한**: 파이프, 변수치환, glob 등 복잡한 명령은 미지원. 단순 `cp src dst`, `mv src dst` 패턴만 감지.

## 구현 순서

| 순서 | 항목 | 의존 |
|---|---|---|
| 0 | Step 0: 의존성 사전 검증 | 없음 (최우선) |
| 1 | Step 1: 삭제된 파일 경로 지원 | Step 0 |
| 2 | Step 2: git rename 감지 | Step 1 (이력 빈 경우 분기) |
| 3 | Step 3: Bash cp/mv 힌트 | 없음 (독립) |

## 검증

```bash
# 삭제된 파일 이력 조회
whyso history internal/analyzer/collect_chain.go
# → 에러 없이 이전 세션 기록 출력

# 이동된 파일의 moved_from 표시
whyso history pkg/analyzer/collect_chain.go
# → moved_from: internal/analyzer/collect_chain.go

# Bash cp/mv 힌트 표시
# → hints 섹션에 mv 명령 참고 정보

# 기존 파일 동작 변경 없음
whyso history pkg/toulmin/validate_graph_def.go
# → 기존과 동일

filefunc validate
go build ./...
go test ./...
```
