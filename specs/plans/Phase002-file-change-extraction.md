# 완료
# Phase002: 파일 변경 추출 + 인과관계 추적

## 목표

파싱된 레코드에서 Write/Edit tool_use를 추출하고, parentUuid 체인을 역추적하여 원래 사용자 요청을 찾는다. 대상은 텍스트 기반 파일에 한정 (바이너리는 Write/Edit의 대상이 아니므로 원천 제외).

## 근거 문서

- files/file-history-from-sessions.md
- specs/plans/Phase001-session-parser-cli.md

## 작업 단계

### Step 1. SSOT 설계

- [ ] `specs/backend/service/extractor.go` — tool_use → FileChange 추출 시퀀스
- [ ] `specs/backend/service/chain.go` — parentUuid 역추적 시퀀스

### Step 2. 파일 변경 추출기

- [ ] `internal/parser/extract.go`
  - Write → file_path, action=create
  - Edit → file_path, old_string, new_string, action=edit
  - Bash → rm/mv/cp/mkdir 패턴 매칭 (1차에서는 스킵 가능)

### Step 3. 인과관계 추적기

- [ ] `internal/history/chain.go`
  - uuid → Record 인덱스 맵 구축
  - parentUuid 체인 역추적 → 최초 user message 탐색
  - 서브에이전트 JSONL 포함

### Step 4. 검증

- [ ] 실제 세션에서 tool_use 추출 + user_request 매핑 확인

## 산출물

| 파일 | 설명 |
|---|---|
| `specs/backend/service/extractor.go` | 추출기 SSOT |
| `specs/backend/service/chain.go` | 체인 추적 SSOT |
| `internal/parser/extract.go` | 변경 추출기 |
| `internal/history/chain.go` | 인과관계 추적기 |
