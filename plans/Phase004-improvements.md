# 완료
# Phase004: 개선사항

## 작업 단계

### Step 1. Action 필드 제거

- [ ] FileChange, ChangeEntry에서 Action 필드 제거
- [ ] Tool 이름(Write/Edit)을 그대로 사용 — 별도 해석하지 않음

### Step 2. summary 필드 제거

- [ ] 기획서에 있지만 user_request + answer로 충분하므로 구현하지 않음
- [ ] 기존 코드에 summary 관련 구현이 있다면 제거

### Step 3. answer 필드 추가

- [ ] tool_use와 같은 assistant 응답에서 `type: "text"` 블록을 추출
- [ ] ChangeEntry에 `Answer string` 필드 추가
- [ ] Claude가 자기 행동을 설명하는 텍스트를 그대로 기록

### Step 4. diff → source 참조로 대체

- [ ] ChangeEntry에서 Diff 필드 제거
- [ ] 대신 Source 필드 추가: JSONL 파일 경로 + 라인 번호
- [ ] 파싱 시 레코드의 JSONL 라인 번호를 추적하여 FileChange에 기록
- 상세 diff가 필요하면 `sed -n '<line>p' <source>` 로 원본 참조

### Step 5. 증분 갱신 구현

- [ ] `--output` 디렉토리 내 파일 중 가장 최근 mtime을 기준 시각으로 사용
- [ ] CLI에서 `BuildHistoriesIncremental(since)` 호출로 변경 (함수는 이미 존재)
- Phase003 설계 결정 참조

### Step 6. 서브에이전트 추적

- 구조: `<session-id>/subagents/agent-<id>.jsonl`
- [ ] ParseSession 시 `<session-id>/subagents/*.jsonl`도 함께 파싱
- [ ] 서브에이전트 레코드를 메인 세션의 RecordIndex에 합쳐 parentUuid 역추적이 세션 전체에서 동작하도록 함
