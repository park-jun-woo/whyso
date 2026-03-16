# 완료
# Phase003: 히스토리 빌더 + CLI 완성

## 목표

파일별 변경 히스토리를 집계하고, `whyso history` CLI로 yaml/json 출력한다.

## 근거 문서

- files/file-history-from-sessions.md
- specs/plans/Phase001-session-parser-cli.md
- specs/plans/Phase002-file-change-extraction.md

## 설계 결정

### 출력 디렉토리 구조

- `--output` 디렉토리 아래에 원본 파일 경로를 미러링
- 예: `whyso history internal/parser/ --output histories/`
  - `internal/parser/jsonl.go` → `histories/internal/parser/jsonl.go.yaml`

### 증분 갱신 전략

- 출력 히스토리 파일의 mtime을 기준 시각으로 사용
- JSONL 파일 중 mtime이 기준 이후인 것만 재파싱 (해당 파일은 처음부터 전체 파싱)
- 첫 실행이면 전체 파싱
- 별도 상태 파일 불필요 — `os.Stat` 비교로 충분

## 작업 단계

### Step 1. SSOT 설계

- [ ] `specs/backend/service/history.go` — 파일별 히스토리 집계 시퀀스

### Step 2. 히스토리 빌더

- [ ] `internal/history/builder.go`
  - 파일 경로 → []ChangeEntry 맵
  - ChangeEntry: timestamp, session, action, user_request, tool, diff
  - 시간순 정렬

### Step 3. 출력 포매터

- [ ] `internal/output/formatter.go` — yaml/json 출력
  - 기획서의 YAML 출력 형식 준수

### Step 4. CLI 확장

- [ ] `whyso history <file>` — 단일 파일 히스토리
- [ ] `whyso history <dir> --all` — 디렉토리 전체
- [ ] 옵션: `--format`, `--since`, `--output`

### Step 5. 검증

- [ ] whyso 프로젝트 자체 세션으로 end-to-end 검증
- [ ] 출력 형식이 기획서 예시와 일치하는지 확인

## 산출물

| 파일 | 설명 |
|---|---|
| `specs/backend/service/history.go` | 히스토리 SSOT |
| `internal/history/builder.go` | 히스토리 빌더 |
| `internal/output/formatter.go` | 출력 포매터 |
| `cmd/whyso/main.go` | CLI 확장 |
