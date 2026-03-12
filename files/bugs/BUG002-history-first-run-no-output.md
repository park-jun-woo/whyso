# 해결됨
# BUG-002: whyso history 첫 조회 시 출력 없음

## 증상

캐시가 없는 파일에 대해 `whyso history <file>` 첫 실행 시 stdout 출력 없음, .whyso/ 캐시도 생성되지 않음.

BUG-001 수정으로 "캐시 있는데 새 변경 없을 때" 는 해결됐으나, "캐시 자체가 없는 첫 조회"에서 여전히 출력 안 됨.

## 재현

```bash
cd ~/.clari/repos/fullend

# 캐시 없는 파일 — 출력 없음
whyso history pkg/authz/authz.go
whyso history internal/gluegen/authzgen.go

# 캐시 있는 파일 (이전에 조회한 적 있음) — 정상 출력
whyso history internal/crosscheck/authz.go
```

## 확인

```bash
# 변경 기록은 존재함
whyso changes | grep "authz.go"
# → 5건 출력됨
```

## 실제 원인

`oldestOutputMtime`이 `.whyso/` 전체를 순회하여 다른 파일 캐시의 mtime을 반환. 이 mtime 이후 세션에 해당 파일 변경이 없으면 빈 결과 반환.

## 수정

단일 파일 조회 시 해당 파일의 캐시 mtime만 확인하도록 변경. 캐시가 없으면 `since`가 zero → 전체 빌드.
