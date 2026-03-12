# 완료
# Phase007: 프로젝트 이름 변경 (whylog → whyso)

## 목표

프로젝트 이름을 whylog에서 whyso로 변경한다. whylog 이름이 선점되어 있으므로 whyso로 전환.

## 변경 대상

### Go 모듈/임포트
- go.mod: module path 변경
- cmd/whylog/ → cmd/whyso/ (디렉토리 이름)
- 모든 Go 파일의 import path 변경 (github.com/clari/whylog → github.com/clari/whyso)
- CLI usage 문자열 변경

### 문서
- CLAUDE.md: 프로젝트명, 루트 경로
- README.md: 전체
- NOTICE: 프로젝트명
- files/file-history-from-sessions.md: 전체 참조
- files/CLAUDE.md.yaml: 전체 참조
- specs/plans/ 내 Phase001, Phase003, Phase006

### 기타
- .gitignore: 바이너리명
- 기존 whylog 바이너리 삭제
