# 완료
# Phase009: whyso map — tree-sitter 기반 키워드 맵

## 목표

`whyso map` 명령어로 코드베이스의 함수명, endpoint, 규칙명 등 키워드를 추출하여 AI 에이전트가 정확한 grep 입력을 할 수 있게 한다.

## 근거 문서

- files/funcmap.md

## 출력 포맷

```
## go
[패키지명]함수1,함수2,...

## ssac
[service/도메인]함수1,함수2,...
```

## 파싱 대상

### tree-sitter 파싱

| 포맷 | 확장자 | 추출 대상 | 그룹 기준 |
|---|---|---|---|
| Go | `.go` | func, method | `package` 선언 |
| TypeScript/JavaScript | `.ts`,`.js`,`.tsx`,`.jsx` | function, arrow function, class method | 상대 디렉토리 |
| Python | `.py` | def, class method | 상대 디렉토리 |
| Rust | `.rs` | fn, impl method | 상대 디렉토리 |
| SSaC | `.ssac` | `func` 이름 | `service/<domain>/` |
| OpenAPI | `.yaml` | `operationId` 값 | `api/` |
| SQL | `.sql` | `-- name:` 쿼리명 | `db/queries/` |
| STML | `.html` | `data-fetch`, `data-action` 값 | `frontend/` |
| Rego | `.rego` | `allow` 규칙명 | `policy/` |

### 정규식 파싱 (tree-sitter 그래머 미지원)

| 포맷 | 확장자 | 추출 대상 | 그룹 기준 |
|---|---|---|---|
| Gherkin | `.feature` | `Scenario:` 이름 | `scenario/` |
| Mermaid | `.md` | state 이름 | `states/` |

## CLI 인터페이스

```
whyso map [path]                  # stdout 출력 + .whyso/_map.md 저장
whyso map [path] -o custom.md    # 지정 경로에 저장
```

- path 미지정 시 현재 디렉토리
- `.gitignore` 패턴 존중 (vendor/, node_modules/ 등 제외)

## 의존성

- `github.com/smacker/go-tree-sitter` — tree-sitter Go 바인딩 (CGO)
- 언어별 그래머: go, typescript, javascript, python, rust, yaml, html, rego

## 구현 단계

### Step 1. 프로젝트 의존성 추가

- [ ] `go-tree-sitter` 및 언어별 그래머 패키지 추가
- [ ] CGO 빌드 확인

### Step 2. 파서 인터페이스 정의

- [ ] `internal/codemap/` 패키지 생성
- [ ] 언어별 파서 인터페이스: `Extract(source []byte) []string`
- [ ] 그룹 결정 로직: 언어별 규칙 적용

### Step 3. tree-sitter 쿼리 구현

- [ ] Go: `(function_declaration name: (identifier) @name)`, `(method_declaration name: (field_identifier) @name)`
- [ ] TypeScript/JavaScript: function, arrow function, class method
- [ ] Python: function_definition
- [ ] Rust: function_item
- [ ] SSaC: Go 파서 재사용 (`.ssac` 확장자)
- [ ] OpenAPI: YAML 파서로 `operationId` 키 추출
- [ ] SQL: 코멘트에서 `-- name:` 패턴 추출 (정규식)
- [ ] STML: HTML 파서로 `data-fetch`, `data-action` 어트리뷰트 추출
- [ ] Rego: `allow` 규칙 추출

### Step 4. 정규식 파서 구현

- [ ] Gherkin: `Scenario:` 뒤 이름 추출
- [ ] Mermaid: stateDiagram 내 state 이름 추출

### Step 5. 파일 탐색 및 출력

- [ ] 디렉토리 재귀 탐색 (`.gitignore` 존중)
- [ ] 확장자별 파서 매핑
- [ ] 섹션별 `[그룹]키워드,...` 포맷 출력
- [ ] `-o` 옵션: 파일 저장

### Step 6. CLI 통합

- [ ] `cmd/whyso/main.go`에 `map` 서브커맨드 추가
- [ ] usage 메시지 갱신

### Step 7. 검증

- [ ] whyso 프로젝트 자체에서 `whyso map` 실행 확인
- [ ] fullend 프로젝트에서 SSOT 키워드 추출 확인

## 산출물

| 파일 | 설명 |
|---|---|
| `internal/codemap/codemap.go` | 코어 로직, 파일 탐색 |
| `internal/codemap/treesitter.go` | tree-sitter 파서 |
| `internal/codemap/regex.go` | 정규식 파서 (Gherkin, Mermaid) |
| `internal/codemap/queries/` | 언어별 tree-sitter 쿼리 |
| `cmd/whyso/main.go` | map 서브커맨드 추가 |
