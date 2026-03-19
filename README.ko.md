# whyso

Claude가 그 파일을 왜 바꿨는지 보여준다 — 세션 JSONL에서 의도를 추출한다.

## 왜 만들었나

Claude Code로 작업하다 보면 커밋이 쌓인다. 몇 주 지나고 코드를 다시 보면 — "이거 왜 이렇게 돼 있지?" `git blame` 쳐본다. 커밋 메시지: `Refactor parser logic`. 그래서? 왜?

답은 이미 내 로컬에 있었다.

Claude Code 세션 로그(JSONL)를 열어봤더니, 모든 파일 변경이 체인으로 연결돼 있었다:

```
사용자: "파서가 빈 세션을 무시하는 버그 고쳐줘"
  → Claude: "빈 세션 필터링 로직이 누락되어 있어서 ParseSession에 early return을 추가합니다"
    → tool_use: Edit parser.go
```

요청 → 판단 → 실행. 전부 기록돼 있다.

`whyso`는 이 체인을 추적해서, 모든 파일 변경에 **원래 요청**과 **Claude의 판단 근거**를 붙여준다.

3개월 전에 Claude한테 시킨 리팩토링이 왜 그런 구조가 됐는지, 지금 바로 볼 수 있다:

```bash
go install github.com/park-jun-woo/whyso/cmd/whyso@latest
whyso history CLAUDE.md   # < 1s, 한 달 프로젝트도
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

## 동작 방식

```
whyso map    → tree-sitter 파싱 → 키워드 맵 (함수, 엔드포인트, 규칙, 상태)
whyso history → JSONL 파싱 → parentUuid 체인 → 파일별 변경 이력
```

## 설치

```bash
go install github.com/park-jun-woo/whyso/cmd/whyso@latest
```

소스에서 빌드:

```bash
git clone https://github.com/park-jun-woo/whyso.git
cd whyso
go build ./cmd/whyso/
```

Go 1.22+ 및 C 컴파일러 필요 (tree-sitter용 CGO).

## 사용법

### 키워드 맵 생성

```bash
# 현재 디렉토리 (stdout + .whyso/_map.md)
whyso map

# 특정 경로
whyso map internal/

# 출력 파일 지정
whyso map -o custom.md

# 강제 재생성
whyso map -f
```

출력 예시:

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

지원 언어: Go, TypeScript, JavaScript, Python, Rust, SSaC, OpenAPI, SQL, Rego, Gherkin, STML (HTML data-*), Mermaid stateDiagram.

### 파일 변경 이력 조회

```bash
# 단일 파일 (stdout + .whyso/ 캐시)
whyso history README.md

# 디렉토리 내 전체 파일 (.whyso/ 캐시만)
whyso history . --all

# 조용 모드 (캐시만, stdout 없음)
whyso history README.md -q

# 출력 디렉토리 지정
whyso history . --all --output custom-dir/

# JSON 형식
whyso history README.md --format json

# 캐시 초기화 후 재빌드
whyso history README.md --reset
```

### 세션 목록

```bash
whyso sessions
```

### 옵션

| 플래그 | 설명 |
|---|---|
| `-o <file>` | 맵 출력 파일 (기본: `.whyso/_map.md`) |
| `-f, --force` | 맵 강제 재생성 (mtime 무시) |
| `--output <dir>` | 이력 출력 디렉토리 (기본: `.whyso/`) |
| `--format <yaml\|json>` | 이력 출력 형식 (기본: yaml) |
| `-q, --quiet` | stdout 출력 억제 |
| `--reset` | 이력 캐시 초기화 후 재빌드 |
| `--all` | 디렉토리 내 전체 파일 포함 |
| `--sessions-dir <path>` | Claude Code 세션 디렉토리 경로 지정 |

## 기능

- **사용자 의도 추적** — `parentUuid` 체인을 따라 원래 사용자 요청까지 역추적
- **AI 판단 근거** — Claude가 왜 그렇게 했는지, 그 설명을 캡처
- **키워드 맵** — tree-sitter 기반 함수, 엔드포인트, 쿼리, 규칙, 상태 추출
- **변경 그룹화** — 같은 요청에서 연속된 편집은 하나로 병합
- **서브에이전트 지원** — 서브에이전트 세션의 변경사항 포함
- **증분 업데이트** — `.whyso/`에 캐시, 새 세션만 재파싱
- **디렉토리 미러링** — 출력 구조가 소스 파일 경로를 그대로 반영

## 라이선스

MIT License — [LICENSE](LICENSE) 참조.
