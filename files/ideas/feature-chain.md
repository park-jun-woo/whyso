# Feature Chain

## 정의

하나의 사용자 기능을 구현하기 위해 레이어를 관통하는 코드 경로.
프론트엔드부터 DB까지 하나의 기능에 엮인 모든 코드.

## 예시: "주문 생성"

```
[프론트엔드]
  OrderForm.tsx        — 폼 UI, 입력 검증
  useCreateOrder.ts    — API 호출 훅

[API 경계]
  POST /api/orders     — 엔드포인트 정의 (OpenAPI)

[백엔드]
  order_handler.go     — HTTP 요청 파싱, 응답
  order_service.go     — 비즈니스 로직 (재고 확인, 가격 계산, 할인 적용)
  order_repo.go        — DB 접근 (INSERT, SELECT)

[인프라]
  003_create_orders.sql — 테이블 스키마
  authz.rego           — "주문 생성 권한" 정책

[테스트]
  order_test.feature    — 시나리오: "재고 있는 상품을 주문하면 성공"
```

어느 하나를 수정하면 chain 전체에 영향이 갈 수 있고, 하나를 빠뜨리면 기능이 깨진다.

## 탐색 방법

### 정적 분석

호출 관계를 따라가면 chain의 일부를 추출할 수 있으나, 프론트엔드↔백엔드 경계, 코드↔SQL, 코드↔정책(Rego) 간 연결은 정적 분석으로 잡을 수 없다.

### 히스토리 co-change

Claude Code 세션에서 "주문 생성 API 만들어줘" 한 마디에 위 파일들이 전부 함께 생성/수정된다. 같은 user_request에 묶여 있으므로 co-change로 한 덩어리. 레이어를 관통하는 연결을 잡아낸다.

### fullend SSaC

SSaC가 feature chain을 선언으로 명시한다. 탐색할 필요 없이 chain이 코드에 보인다.

## whyso와 fullend의 역할 분리

- **fullend:** feature chain이 뭔지 (구조, 선언)
- **whyso:** 그 chain이 왜 바뀌었는지 (이력, 추적)

feature chain 탐색 자체는 fullend의 영역. whyso는 co-change 관계로 chain을 보강하고, chain의 변경 이력을 추적하는 역할.
