# whyso map: tree-sitter 기반 함수명 맵

## 핵심 아이디어

AI 에이전트가 코드베이스를 탐색할 때 grep으로 추측 검색하는 대신, 함수명 목록을 먼저 보고 정확한 이름으로 grep한다.

tree-sitter로 소스코드를 파싱하여 패키지별 함수명 목록을 생성한다. 경로/라인 번호는 포함하지 않는다 — grep이 알려주므로 중복 정보다.

## 출력 포맷

```
[패키지명]함수1,함수2,함수3,...
```

예시:
```
[authz]CheckOwnership,CheckRequest,LoadPolicy,NewAuthorizer
[crosscheck]ValidateSSaCDDL,ValidateStates,ValidateFunc,ValidateScenario
[gluegen]GenerateHandler,GenerateRouter,GenerateModel
```

## 왜 이 포맷인가

- **경로/라인 제거** — grep 결과가 알려준다. 중복 정보는 토큰 낭비
- **한 줄 한 패키지** — 10만 파일, 함수 50만 개여도 몇 MB 텍스트 파일 하나
- **마이크로서비스 단위면 한눈에** — 서비스당 함수 수백 개 = 토큰 몇천. Read 한 번으로 전체 구조 파악
- **맵 자체도 grep 대상** — `grep "CheckOwnership" .code-map` → 어떤 패키지인지 즉시

## AI 에이전트 작업 흐름

1. `Read .code-map` → 함수 전체 목록 파악
2. `Read .file-histories/해당파일.yaml` → 왜 이렇게 됐는지 파악
3. `Grep "함수명"` → 정확한 위치 찍기
4. 수정

grep 전에 이미 전체 그림이 잡혀있다.

## 구현

### 명령어

```
whyso map [path]                  # stdout 출력
whyso map [path] -o .code-map    # 파일 저장
```

### tree-sitter 파싱 대상

| 언어 | 추출 대상 |
|---|---|
| Go | func, method (receiver 포함) |
| TypeScript/JavaScript | function, arrow function, class method |
| Python | def, class method |
| Rust | fn, impl method |

### 패키지 결정 규칙

- Go: `package` 선언
- TypeScript/JavaScript: 디렉토리명 또는 export 모듈명
- Python: 디렉토리명 (\_\_init\_\_.py 기준)
- Rust: `mod` 선언

### 갱신

- `whyso map`은 호출 시점에 tree-sitter로 파싱하여 생성
- 캐싱은 선택사항 — 파일 mtime 비교로 변경분만 재파싱 가능

## whyso 통합

whyso는 두 가지 질문에 답한다:

- `whyso map` — **뭐가 있는지** (구조)
- `whyso history` — **왜 그런지** (이유)

grep은 AI가 이미 잘 한다. whyso는 grep의 입력을 정확하게 만들어주는 나침반이다.
