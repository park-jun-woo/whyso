# 해결됨
# BUG-001: whyso history 출력 없음

## 증상

`whyso history <file>` 및 `whyso history <dir> --all` 실행 시 stdout 출력 없음, .whyso/ 캐시도 생성되지 않음.

## 재현

```bash
# fullend 프로젝트 루트에서 실행
cd ~/.clari/repos/fullend

# 단일 파일 — 출력 없음
whyso history internal/crosscheck/authz.go

# 디렉토리 — 출력 없음
whyso history internal/crosscheck/

# --all 옵션 — 출력 없음
whyso history internal/ --all
```

## 정상 동작하는 명령

```bash
whyso sessions    # 세션 18개 정상 출력
whyso changes     # 4,539건 변경 기록 정상 출력
whyso map         # 함수명 맵 정상 출력
```

## 단서

- `whyso changes | grep "authz.go"` → 변경 기록 5건 존재 확인
- 세션 탐지, 파일 변경 감지는 정상
- history 서브커맨드만 출력 안 됨
- 리브랜딩 전 `whylog history` 로는 동일 파일에서 정상 출력됨
- `whyso history internal/crosscheck/ssac_openapi.go` 는 이전에 1건 캐시 성공 (.whyso/ 에 존재)

## 추정 원인

리브랜딩(whylog → whyso) 과정에서 history 서브커맨드의 파일 경로 매칭 또는 출력 로직에 regression 발생 가능.

## 실제 원인

증분 모드에서 새 변경이 없으면 `histories`가 빈 맵으로 반환되어 stdout 출력과 캐시 기록이 모두 스킵됨. 기존 캐시 파일이 있어도 읽지 않았음.

## 수정

`len(histories) == 0`일 때 기존 캐시 파일(`OutputPath`)에서 읽어서 stdout에 출력하도록 수정.
