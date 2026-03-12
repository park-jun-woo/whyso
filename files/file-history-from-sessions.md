# whyso: Claude Code 세션에서 파일별 변경 히스토리 추출

## 핵심 아이디어

Claude Code는 모든 대화 세션을 `~/.claude/projects/<project-id>/<session-id>.jsonl`에 저장한다.
각 JSONL 라인에는 `type`(user/assistant), `parentUuid` 체인, tool_use(Write/Edit/Bash) 정보가 포함된다.

**이 데이터로 프로젝트 내 모든 텍스트 기반 파일(소스코드, 기획서, 설계 문서 등)에 "왜 이 파일이 생겼고, 누가 어떤 요청으로 수정했는지" 히스토리를 자동 구축할 수 있다.** Write/Edit tool_use는 텍스트 파일만 대상으로 하므로 바이너리 파일은 추적 대상에서 제외된다.

git blame이 "누가 언제 뭘 바꿨는지"를 보여준다면,
이 도구는 **"왜 바꿨는지"** — 사용자의 원래 의도와 AI의 판단 근거를 보여준다.

## 데이터 구조 (검증 완료)

### 세션 파일 위치
```
~/.claude/projects/<project-path-slug>/<session-id>.jsonl
```

### JSONL 레코드 타입

| type | 역할 | 핵심 필드 |
|---|---|---|
| `user` | 사용자 메시지 | `message.content` (string), `uuid`, `timestamp` |
| `assistant` | AI 응답 + tool_use | `message.content[]` (tool_use 배열), `parentUuid` |
| `tool_result` | 도구 실행 결과 | `content` (성공/실패 메시지) |

### 파일 변경 감지 가능한 tool_use

| tool | input 필드 | 감지 방법 |
|---|---|---|
| `Write` | `file_path`, `content` | 파일 생성/덮어쓰기 |
| `Edit` | `file_path`, `old_string`, `new_string` | 부분 수정 |
| `Bash` | `command` | `rm`, `mv`, `cp`, `mkdir` 패턴 매칭 |

### 인과관계 추적

`parentUuid` 체인을 역방향으로 탐색하면, 어떤 tool_use가 어떤 사용자 요청에서 비롯되었는지 추적 가능:

```
user message (uuid: A)
  → assistant tool_use Write (parentUuid: A)
    → tool_result (parentUuid: B)
      → assistant tool_use Edit (parentUuid: C)  ← 같은 턴의 후속 작업
```

중간에 `tool_result` 타입의 user 메시지가 끼어 있으므로, `type == "user"` && `content`가 문자열인 것만 추출하면 진짜 사용자 요청을 찾을 수 있다.

## 실측 데이터 (검증용: whyso 프로젝트)

| 항목 | 수치 |
|---|---|
| 세션 수 | 17 |
| Write 호출 | 660 |
| Edit 호출 | 1,415 |
| Bash 파일 조작 | 206 |
| 변경된 고유 파일 수 | 480 |

## 출력 형식 (제안)

### 파일별 히스토리 (예: 소스코드, 기획서, 설계 문서 등 모든 파일 대상)

```yaml
file: internal/crosscheck/crosscheck.go
created: 2026-03-08T14:23:01Z
history:
  - timestamp: 2026-03-08T14:23:01Z
    session: 09351222-d7be-41fe-994f-87c2d7067e5d
    action: create
    user_request: "whyso validate 명령어 구현 시작해"
    tool: Write
    summary: 초기 crosscheck 패키지 생성 — SSaC↔OpenAPI operationId 매칭 검증

  - timestamp: 2026-03-09T10:15:33Z
    session: 4e9b4e5e-3a50-43f2-be6e-e5db228ecc3b
    action: edit
    user_request: "x-sort 컬럼이 DDL에 있는지 검증 추가해"
    tool: Edit
    diff: |
      + // x-sort allowed columns must exist in DDL
      + func checkSortColumns(...)
    summary: x-sort/x-filter 컬럼 → DDL 교차 검증 추가

  - timestamp: 2026-03-10T09:41:22Z
    session: b2e43b4f-cb21-4286-975d-1eb9de8a16c0
    action: edit
    user_request: "Func 스펙 교차 검증도 추가해"
    tool: Edit
    diff: |
      + func checkFuncSpec(...)
    summary: @call ↔ Func spec 인자 수/타입 교차 검증 추가
```

### summary 생성

`summary`는 tool_use 전후의 assistant 텍스트 메시지에서 추출하거나,
diff(old_string→new_string)를 보고 한 줄 요약을 생성.

## CLI 설계 (안)

```bash
# 특정 파일의 히스토리 추출
whyso history internal/crosscheck/crosscheck.go

# 디렉토리 전체 파일의 히스토리 일괄 생성
whyso history internal/crosscheck/ --output histories/

# 특정 프로젝트의 모든 파일 히스토리 (소스코드, 문서 등 전체)
whyso history . --all --output .file-histories/

# 세션 목록 조회
whyso sessions

# 특정 세션의 변경 파일 목록
whyso sessions <session-id> --files
```

### 옵션

| 플래그 | 설명 |
|---|---|
| `--sessions-dir` | Claude Code 세션 디렉토리 (기본: 자동 감지) |
| `--output` | 출력 디렉토리 (기본: stdout) |
| `--format` | `yaml` / `json` / `markdown` |
| `--since` | 특정 날짜 이후만 |
| `--all` | 모든 파일의 히스토리 생성 |

## 구현 고려사항

### 1. 세션 디렉토리 자동 감지

현재 작업 디렉토리의 절대 경로에서 `/`와 `.`을 `-`로 치환하여 프로젝트 슬러그를 생성:
```
/home/parkjunwoo/.clari/repos/whyso
→ -home-parkjunwoo--clari-repos-whyso
```

### 2. Bash 명령어에서 파일 조작 감지

Write/Edit는 명시적이지만, Bash는 휴리스틱이 필요:
- `rm file` → delete
- `mv old new` → rename
- `cp src dst` → copy
- `mkdir -p dir` → create directory
- `sed -i`, `echo > file` → edit (Write/Edit를 쓰라는 지침이 있으므로 드묾)

### 3. 서브에이전트 추적

`~/.claude/projects/<project>/<session-id>/subagents/` 디렉토리에 서브에이전트 세션이 있다.
서브에이전트의 tool_use도 포함해야 완전한 히스토리가 된다.

### 4. 성능

검증 프로젝트(whyso): 17개 세션, 2,281개 tool_use. 파싱 시간은 무시할 수준.
대규모 프로젝트(수백 세션)에서도 JSONL 스트리밍 파싱이므로 메모리 문제 없음.

### 5. 개인정보

세션 데이터에는 사용자 메시지 전문이 포함된다.
히스토리 출력에 포함할 범위를 `--redact` 옵션으로 제어.

## 프로젝트 구조

whyso는 독립 프로젝트로 개발한다.

### 주요 컴포넌트

| 컴포넌트 | 역할 |
|---|---|
| CLI (`whyso`) | 세션 파싱, 히스토리 추출, 로컬 조회 |
| 백엔드 (Go) | 파싱된 데이터를 PostgreSQL에 저장, API 제공 |
| 앱 프론트엔드 (React) | 파일별 히스토리 탐색 UI |
| 콘텐츠 프론트엔드 (Hugo) | 프로젝트 소개 페이지 |

### 핵심 데이터 흐름

```
~/.claude/projects/**/*.jsonl
  → whyso parse (JSONL 파싱 + tool_use 추출)
  → whyso history (parentUuid 체인 역추적 → 파일별 히스토리 구축)
  → DB 저장 or stdout 출력
```
