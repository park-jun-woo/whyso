# 함수 단위 히스토리 설계 논의

## 배경

현재 whyso는 파일 단위 히스토리를 생성한다. 파일-히스토리도 강력하지만, 함수 단위로 세분화하면 컨텍스트를 더 날카롭게 좁힐 수 있다.

## 목표

컨텍스트를 날카롭게 좁혀서 LLM에게 쓸데없는 컨텍스트 오염을 먹이지 않는다.

## 핵심 아이디어: 심볼릭 recall + RAG precision

### 파이프라인

1. **심볼릭 recall** — 정적 call graph + 히스토리 co-change로 넓게 후보군 확보 (누락 방지)
2. **RAG precision** — 그 안에서 LLM이 정밀 필터링 (노이즈 제거)

넓혔다가 좁히는 것과 애초에 넓히지 않는 것은 다르다. 넓히지 않으면 누락이 발생한다.

### 심볼릭으로 넓히되 최소한으로 넓히고, 가능한 소형 LLM으로 정밀하게 깎아낸다

사실상 RAG인데, 일반 RAG와 다른 점:
- 일반 RAG: 전체 코드베이스 → 임베딩 → 유사도 검색 → 노이즈 범벅
- whyso 방식: 심볼릭(call graph + co-change + history)으로 후보군을 먼저 확 좁힘 → 그 좁은 범위 안에서 RAG

whyso는 RAG의 retrieval 앞단에 붙는 심볼릭 필터다.

## 연관 그래프: 두 가지 소스

### 1. 정적 분석 (tree-sitter)

명시적 호출 관계. 확실한 엣지.

- funcB가 funcC를 호출 → 연결
- funcA가 funcB를 호출 → 연결

### 2. 히스토리 co-change

호출 관계 없지만 함께 변경된 함수들. 숨은 엣지 보강.

- funcB와 funcD가 같은 user_request로 함께 생성/수정됨 → 연결
- 정적 분석만으로는 funcD를 절대 못 찾음

둘 다 쓴다. 정적 분석이 기본 그래프, 히스토리가 엣지를 보강한다.

## 데이터 구조

파일-히스토리 위에 func 레이어를 얹는 구조.

```yaml
file: internal/auth/verify.go
functions:
  funcB:
    history:
      - "인증 로직 수정" → Edit
    co_changed: [funcA, funcC]
  funcC:
    history:
      - "토큰 검증 추가" → Edit
file_history:
  - "인증 모듈 생성" → Write
```

- 파일 단위로 보고 싶으면 file_history
- 함수 단위로 좁히고 싶으면 functions
- 둘 다 유지

## --grep과의 관계

`--grep`과 함수 단위 히스토리는 같은 주제다.

- 함수 단위 히스토리로 데이터를 구조화
- `--grep`으로 키워드 매칭되는 함수 목록을 추출
- 함수 이름 자체도 매칭 대상

```bash
whyso history internal/auth/verify.go --grep "인증"
# → funcB (user_request에 "인증" 포함)

whyso history internal/auth/verify.go --grep "funcC"
# → funcC (함수 이름 자체 매칭)
```

## 구현 핵심 포인트

### Edit 변경 → 함수 매핑

Edit의 old_string/new_string이 어떤 함수 범위에 속하는지 tree-sitter로 매핑.
변경된 줄 번호 → 해당 줄을 감싸는 함수 노드 찾기.

### Go 우선

- Go가 심볼릭 분석에 가장 쉬움 (거의 정적, 동적 디스패치 거의 없음, 패키지 구조 명확)
- Go에서 먼저 증명하고 → TypeScript, Python 순으로 확장
- Go에서 안 되면 다른 언어에서도 안 됨

## SILK와의 관계

SILK (Symbolic Index for LLM Knowledge) 아키텍처의 코드 도메인 적용.

| SILK | whyso |
|---|---|
| 심볼릭 (코드북, 비트 AND) | 심볼릭 (tree-sitter, call graph, co-change) |
| 뉴럴 (LLM 태깅/분류) | RAG (LLM 정제) |
| 기록 시점에 구조 부여 | 세션 시점에 히스토리 구조화 |
| 검색이 공짜 | Grep이 공짜 |

SILK의 열화 구현체. 탐색 엔진 관점에서 SIDX 64비트 SIMD 비트마스크 vs 텍스트 grep이라 성능 차이는 명확하나, whyso 도메인(코드베이스 히스토리)은 데이터 규모가 수천~수만 건이라 grep으로 충분.

공통 철학: "검색이 어려운 건 구조 없는 데이터의 증상이다."
코드베이스 탐색이 어려운 것도 컨텍스트에 구조가 없어서다.
