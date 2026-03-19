# 완료
# Phase001: JSONL 파서 + 데이터 모델

## 목표

Claude Code 세션 JSONL 파일을 파싱하여 Go 구조체로 변환한다. CLI 뼈대를 잡고 `whyso sessions` 명령으로 세션 목록을 출력하는 것까지.

## 근거 문서

- files/file-history-from-sessions.md

## 작업 단계

### Step 1. SSOT 설계

- [ ] `specs/backend/service/parser.go` — JSONL 레코드 → Go 구조체 변환 시퀀스

### Step 2. Go 프로젝트 초기화

- [ ] `cmd/whyso/` — main.go (cobra CLI)
- [ ] `internal/model/` — 데이터 모델
- [ ] `internal/parser/` — JSONL 파서
- [ ] `go mod init github.com/clari/whyso`

### Step 3. 데이터 모델 정의

- [ ] `internal/model/record.go` — JSONL 레코드 구조체
  - Record: type, uuid, parentUuid, timestamp, sessionId, message
  - message.content: string(user) 또는 []ContentBlock(assistant)
  - ContentBlock: type(text/tool_use/tool_result), tool name, input

### Step 4. JSONL 스트리밍 파서

- [ ] `internal/parser/jsonl.go`
  - bufio.Scanner로 라인 단위 스트리밍
  - json.Unmarshal → Record
  - 세션 디렉토리 자동 감지 (cwd → slug → `~/.claude/projects/<slug>/`)

### Step 5. CLI — `whyso sessions`

- [ ] `cmd/whyso/main.go` — cobra 기반
  - `whyso sessions` — 현재 프로젝트의 세션 목록 (ID, 시간, 첫 사용자 메시지)
  - `--sessions-dir` 옵션

### Step 6. 검증

- [ ] whyso 프로젝트 자체 세션 JSONL로 `whyso sessions` 동작 확인

## 산출물

| 파일 | 설명 |
|---|---|
| `specs/backend/service/parser.go` | 파서 SSOT |
| `cmd/whyso/main.go` | CLI 진입점 |
| `internal/model/record.go` | JSONL 레코드 모델 |
| `internal/parser/jsonl.go` | 스트리밍 파서 |
