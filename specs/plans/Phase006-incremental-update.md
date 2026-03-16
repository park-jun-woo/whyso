# 완료
# Phase006: 증분 갱신

## 목표

`--output` 사용 시, 이미 생성된 히스토리 파일을 기준으로 변경된 세션만 재파싱하고 기존 히스토리와 병합한다.

## 설계

### 동작 흐름

1. 대상 출력 파일(예: `histories/CLAUDE.md.yaml`)이 존재하는지 확인
2. 존재하면 → 해당 파일의 mtime을 `since`로 사용
3. mtime 이후에 수정된 JSONL만 재파싱 (해당 파일은 처음부터 전체 파싱)
4. 기존 히스토리 파일을 읽어서 기존 항목 로드
5. 새로 추출한 항목을 기존 항목과 병합 (timestamp + source로 중복 제거)
6. 출력 파일이 없으면 전체 파싱 (기존과 동일)

### 중복 제거 기준

ChangeEntry의 (timestamp, sources[0]) 조합이 같으면 동일 항목으로 판단.

### 구현 대상

- [ ] `internal/output/reader.go` — 기존 YAML/JSON 히스토리 파일 읽기
- [ ] `internal/history/merge.go` — 기존 + 신규 히스토리 병합 (중복 제거, 시간순 정렬)
- [ ] `internal/output/formatter.go` — WriteHistories에 증분 로직 통합
- [ ] `cmd/whyso/main.go` — 출력 파일별 mtime 기반 증분 호출
