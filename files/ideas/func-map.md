# whyso map: tree-sitter 기반 키워드 맵

## 핵심 아이디어

AI 에이전트가 코드베이스를 탐색할 때 grep으로 추측 검색하는 대신, 키워드 목록을 먼저 보고 정확한 이름으로 grep한다.

tree-sitter로 소스코드와 SSOT를 파싱하여 함수명, endpoint, 규칙명 등 키워드를 추출한다. 경로/라인 번호는 포함하지 않는다 — grep이 알려주므로 중복 정보다.


## 출력 포맷

언어별 섹션으로 구분하고, 그룹 단위로 함수명을 나열한다.

- Go: `package` 선언 그대로
- TypeScript/JavaScript: 상대 디렉토리 경로
- Python: 상대 디렉토리 경로
- Rust: 상대 디렉토리 경로

```
## go
[authz]CheckOwnership,CheckRequest,LoadPolicy,NewAuthorizer
[parser]ParseSession,ExtractChanges,ListSessions
[history]BuildHistories,BuildIndex,FindUserRequest

## typescript
[src/components]Button,Modal,Sidebar,Toast
[src/utils]formatDate,parseQuery,debounce

## python
[app/models]User,Project,Session
[app/services]authenticate,authorize,validate

## ssac
[service/gig]CreateGig,UpdateGig,DeleteGig,PublishGig
[service/auth]Register,Login,RefreshToken

## openapi
[api]CreateGig,UpdateGig,ListGigs,PublishGig,Register,Login

## sql
[db/queries/gigs]GigCreate,GigFindByID,GigList,GigUpdateStatus
[db/queries/users]UserCreate,UserFindByID,UserFindByEmail

## rego
[policy]allow

## gherkin
[scenario/gig]CreateAndPublishGig,BidOnGig
[scenario/auth]RegisterAndLogin

## stml
[frontend/gigs]data-fetch:ListGigs,data-action:CreateGig
[frontend/auth]data-action:Login,data-action:Register

## mermaid
[states/gig]unpublished,published,deleted
[states/course]unpublished,published
```

## 왜 이 포맷인가

- **경로/라인 제거** — grep 결과가 알려준다. 중복 정보는 토큰 낭비
- **언어별 섹션** — AI가 특정 언어만 선택적으로 grep 가능
- **한 줄 한 그룹** — 10만 파일, 함수 50만 개여도 몇 MB 텍스트 파일 하나
- **마이크로서비스 단위면 한눈에** — 서비스당 함수 수백 개 = 토큰 몇천. Read 한 번으로 전체 구조 파악
- **맵 자체도 grep 대상** — `grep "CheckOwnership" .code-map` → 어떤 패키지인지 즉시

## AI 에이전트 작업 흐름

1. `Read .code-map` → 키워드 전체 목록 파악 (함수, endpoint, 규칙, 상태)
2. `Read .whyso/해당파일.yaml` → 왜 이렇게 됐는지 파악
3. `Grep "키워드"` → SSaC, OpenAPI, SQL, Go 어디든 한방에 찍기
4. 수정

grep 전에 이미 전체 그림이 잡혀있다. `CreateGig`을 grep하면 SSaC 선언, OpenAPI endpoint, sqlc 쿼리, Gherkin 시나리오, STML 바인딩이 전부 나온다.

## 구현

### 명령어

```
whyso map [path]                  # stdout 출력
whyso map [path] -o .code-map    # 파일 저장
```

### tree-sitter 파싱 대상

#### 범용 언어

| 언어 | 추출 대상 | 그룹 기준 |
|---|---|---|
| Go | func, method (receiver 포함) | `package` 선언 |
| TypeScript/JavaScript | function, arrow function, class method | 상대 디렉토리 경로 |
| Python | def, class method | 상대 디렉토리 경로 |
| Rust | fn, impl method | 상대 디렉토리 경로 |

#### fullend SSOT

| 포맷 | 확장자 | 추출 대상 | 그룹 기준 |
|---|---|---|---|
| SSaC | `.ssac` | `func` 이름 | `service/<domain>/` |
| OpenAPI | `.yaml` | `operationId` | `api/` |
| SQL | `.sql` | `-- name:` 쿼리명 | `db/queries/` |
| Rego | `.rego` | `allow` 규칙명 | `policy/` |
| Gherkin | `.feature` | `Scenario` 이름 | `scenario/` |
| STML | `.html` | `data-fetch`, `data-action` 값 | `frontend/` |
| Mermaid | `.md` | state 이름 | `states/` |

### 갱신

- `whyso map`은 호출 시점에 tree-sitter로 파싱하여 생성
- 캐싱은 선택사항 — 파일 mtime 비교로 변경분만 재파싱 가능

## whyso 통합

whyso는 두 가지 질문에 답한다:

- `whyso map` — **뭐가 있는지** (구조: 함수, endpoint, 규칙, 상태)
- `whyso history` — **왜 그런지** (이유: 변경 의도와 맥락)

grep은 AI가 이미 잘 한다. whyso는 grep의 입력을 정확하게 만들어주는 나침반이다.

fullend가 SSOT를 선언하고, whyso가 그 구조를 맵으로 제공하고, Claude Code가 맵을 보고 정확하게 탐색/수정한다.
