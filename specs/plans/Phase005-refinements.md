# Phase005: 출력 정제 + 증분 갱신 + 서브에이전트 표시

## 작업 단계

### Step 1. 동일 요청 그룹핑

동일한 user_request + answer에서 연속 발생한 tool_use를 하나의 항목으로 묶는다.
source만 리스트로 확장.

변경 전:
```yaml
  - timestamp: 2026-03-12T01:32:09Z
    user_request: "go cli로 할꺼야"
    answer: "완료. 변경 내용:..."
    tool: Edit
    source: ...441b6643.jsonl:132
  - timestamp: 2026-03-12T01:32:13Z
    user_request: "go cli로 할꺼야"
    answer: "완료. 변경 내용:..."
    tool: Edit
    source: ...441b6643.jsonl:135
  - timestamp: 2026-03-12T01:32:18Z
    user_request: "go cli로 할꺼야"
    answer: "완료. 변경 내용:..."
    tool: Edit
    source: ...441b6643.jsonl:138
```

변경 후:
```yaml
  - timestamp: 2026-03-12T01:32:09Z
    user_request: "go cli로 할꺼야"
    answer: "완료. 변경 내용:..."
    tool: Edit
    sources:
      - ...441b6643.jsonl:132
      - ...441b6643.jsonl:135
      - ...441b6643.jsonl:138
```

- [ ] ChangeEntry의 source를 sources []Source로 변경
- [ ] 빌더에서 연속된 동일 (user_request, answer, tool, session) 항목을 병합
- [ ] 타임스탬프는 그룹의 첫 항목 기준

### Step 2. 증분 갱신

- [ ] 대상 출력 파일(예: `histories/CLAUDE.md.yaml`)의 mtime을 기준으로 사용
- [ ] mtime 이후에 수정된 JSONL만 재파싱 (해당 파일은 처음부터 전체 파싱)
- [ ] 기존 히스토리 파일을 읽어서 새로 추출한 항목과 병합 (timestamp 기준 정렬, 중복 제거)
- [ ] 출력 파일이 없으면 전체 파싱

### Step 3. 서브에이전트 표시

- [ ] 서브에이전트 JSONL에서 발생한 변경에 `subagent: true` 표시
- [ ] 서브에이전트 파일명에서 agent ID 추출하여 기록
- [ ] 가능하면 skill 이름 추출 — 서브에이전트 JSONL 내 첫 user 메시지 또는 메타데이터에서 skill 정보 탐색 필요
