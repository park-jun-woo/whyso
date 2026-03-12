# 완료
# Phase008: 기본 출력 디렉토리 `.whyso/` 도입

## 목표

`--output`을 지정하지 않아도 프로젝트 루트에 `.whyso/` 디렉토리를 기본 출력 경로로 사용한다. 증분 업데이트(mtime 기반)가 항상 작동하여 세션 순회를 최소화한다.

## 현재 동작

- `--output` 미지정 → stdout 출력, 캐시 없음
- `--output dir` 지정 → dir에 저장, mtime 기반 증분 업데이트

## 변경 후 동작

- `--output` 미지정 → `.whyso/`에 저장 (기본값), 증분 업데이트 작동
- `--output dir` 지정 → 기존과 동일, 지정 경로 우선
- `--stdout` 플래그 추가 → 명시적으로 stdout 출력 (기존 동작 유지용)

## 변경 대상

### cmd/whyso/main.go
- `outputDir` 기본값을 프로젝트 루트 기준 `.whyso/`로 설정
- `--stdout` 플래그 추가: outputDir 무시하고 stdout 출력
- `outputDir == ""` 분기 제거 → 항상 파일 출력이 기본

### .gitignore
- `.whyso/` 추가

### 문서
- README.md: 기본 출력 경로 설명 갱신

## 디렉토리 구조

```
project-root/
├── .whyso/
│   ├── internal/
│   │   └── parser/
│   │       └── jsonl.yaml
│   └── cmd/
│       └── whyso/
│           └── main.yaml
```

- 프로젝트 파일 경로를 미러링 (기존 `--output` 동작과 동일)

## 핵심 로직

1. `outputDir`이 비어있으면 `filepath.Join(projectRoot, ".whyso")`로 설정
2. `.whyso/` 디렉토리 없으면 자동 생성
3. `oldestOutputMtime`이 `.whyso/` 내 파일 확인 → 이후 세션만 파싱
4. `--stdout` 지정 시 기존 stdout 출력 경로 실행
