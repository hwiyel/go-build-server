# API 엔드포인트 문서

## 기본 정보

- **Base URL**: `http://localhost:8080`
- **Content-Type**: `application/json`

---

## 엔드포인트

### 1. GET /api/build

API 서버의 상태를 확인합니다.

**요청:**
```bash
curl http://localhost:8080/api/build
```

**응답 (200 OK):**
```json
{
  "status": "success",
  "message": "Build API is working!",
  "version": "1.0.0"
}
```

**응답 필드:**
| 필드 | 타입 | 설명 |
|------|------|------|
| status | string | 작업 상태 (success/error) |
| message | string | 상태 메시지 |
| version | string | API 버전 |

---

### 2. POST /api/build/create

새로운 빌드를 생성합니다.

**요청:**
```bash
curl -X POST http://localhost:8080/api/build/create \
  -H "Content-Type: application/json" \
  -d '{
    "project_name": "my-project",
    "environment": "production"
  }'
```

**요청 본문:**
| 필드 | 타입 | 필수 | 설명 |
|------|------|------|------|
| project_name | string | ✅ | 프로젝트 이름 |
| environment | string | ✅ | 배포 환경 (development/staging/production) |

**응답 (201 Created):**
```json
{
  "status": "success",
  "message": "Build created successfully!",
  "project_name": "my-project",
  "build_id": "build-my-project-001",
  "timestamp": "2026-02-03T14:36:00Z"
}
```

**응답 필드:**
| 필드 | 타입 | 설명 |
|------|------|------|
| status | string | 작업 상태 |
| message | string | 상태 메시지 |
| project_name | string | 생성된 프로젝트 이름 |
| build_id | string | 생성된 빌드 ID |
| timestamp | string | 생성 시간 (ISO 8601) |

---

## 에러 응답

### 400 Bad Request

필수 필드가 누락된 경우:

```json
{
  "error": "project_name and environment are required"
}
```

### 405 Method Not Allowed

지원하지 않는 HTTP 메서드를 사용한 경우:

```json
{
  "error": "Only POST method is allowed"
}
```

---

## 테스트 예제

### cURL

```bash
# GET 요청
curl http://localhost:8080/api/build

# POST 요청 (성공)
curl -X POST http://localhost:8080/api/build/create \
  -H "Content-Type: application/json" \
  -d '{"project_name":"test","environment":"dev"}'

# POST 요청 (필드 누락)
curl -X POST http://localhost:8080/api/build/create \
  -H "Content-Type: application/json" \
  -d '{"project_name":"test"}'

# POST 요청 (메서드 오류)
curl -X GET http://localhost:8080/api/build/create
```

### PowerShell

```powershell
# GET 요청
curl.exe http://localhost:8080/api/build | ConvertFrom-Json

# POST 요청
$body = @{
    project_name = "my-project"
    environment = "production"
} | ConvertTo-Json

curl.exe -X POST http://localhost:8080/api/build/create `
  -Headers @{"Content-Type"="application/json"} `
  -Body $body | ConvertFrom-Json
```

### Bash/Git Bash

```bash
# GET 요청 (pretty print)
curl http://localhost:8080/api/build | jq .

# POST 요청
curl -X POST http://localhost:8080/api/build/create \
  -H "Content-Type: application/json" \
  -d '{
    "project_name": "my-project",
    "environment": "production"
  }' | jq .
```

---

## HTTP 상태 코드

| 코드 | 설명 |
|------|------|
| 200 | OK - 요청 성공 |
| 201 | Created - 리소스 생성 성공 |
| 400 | Bad Request - 잘못된 요청 |
| 405 | Method Not Allowed - 지원하지 않는 메서드 |

---

## 주의사항

1. `project_name`과 `environment`는 필수 필드입니다.
2. POST 요청 시 Content-Type을 `application/json`으로 설정해야 합니다.
3. 모든 응답은 JSON 형식입니다.
