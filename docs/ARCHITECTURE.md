# 프로젝트 아키텍처

## 개요

api-server v3.0은 깔끔한 패키지 구조와 의존성 주입을 기반으로 한 Go 애플리케이션입니다.

## 계층 구조

```
┌─────────────────────────────────────┐
│     HTTP Layer (main.go)            │
│  ├─ GET /api/buildjob/{job}/logs   │
│  └─ POST /api/buildjob             │
└──────────────┬──────────────────────┘
               │
┌──────────────▼──────────────────────┐
│  Handlers Layer (pkg/handlers)      │
│  ├─ BuildJobHandler                 │
│  └─ LogsHandler                     │
└──────────────┬──────────────────────┘
               │
┌──────────────▼──────────────────────┐
│  Services Layer (pkg/services)      │
│  ├─ LogService (Interface)          │
│  └─ InMemoryLogService              │
└──────────────┬──────────────────────┘
               │
┌──────────────▼──────────────────────┐
│  Storage Layer (pkg/storage)        │
│  ├─ LogStorage (Interface)          │
│  └─ MemoryStorage                   │
└──────────────┬──────────────────────┘
               │
┌──────────────▼──────────────────────┐
│  Models (pkg/models)                │
│  ├─ BuildJobRequest                 │
│  ├─ BuildJobResponse                │
│  └─ LogEntry                        │
└─────────────────────────────────────┘
```

## 패키지 상세

### 1. pkg/models - 데이터 모델

**책임:** 애플리케이션의 모든 데이터 구조 정의

**구성원:**
```go
type BuildJobRequest struct {
    JobName           string
    DockerfileContent string
}

type LogEntry struct {
    Timestamp string
    Container string
    Message   string
    Level     string
}

type LogsResponse struct {
    JobName    string
    Status     string
    Logs       []LogEntry
    TotalLines int
}
```

**의존성:** 없음 (다른 패키지에 의존하지 않음)

---

### 2. pkg/utils - 유틸리티 함수

**책임:** 공통 로직 제공

**주요 함수:**
```go
func DetectLogLevel(message string) string {
    // "error", "warn", "info" 중 하나 반환
}
```

**의존성:** 없음

---

### 3. pkg/storage - 저장소 계층

**책임:** 데이터 저장소 추상화

**인터페이스:**
```go
type LogStorage interface {
    SaveLog(jobName, container, message, level string)
    GetLogs(jobName string) ([]models.LogEntry, bool)
    DeleteLogs(jobName string)
    Exists(jobName string) bool
}
```

**구현체:**
- `MemoryStorage` - 메모리 기반 저장소 (RWMutex로 스레드 안전성 보장)

**의존성:** models

**확장 가능성:**
```go
// Phase 2: ConfigMapStorage 구현
type ConfigMapStorage struct {
    clientset kubernetes.Interface
}

// Phase 3: DatabaseStorage 구현
type DatabaseStorage struct {
    db *gorm.DB
}
```

---

### 4. pkg/services - 비즈니스 로직 계층

**책임:** 로그 관련 비즈니스 로직 구현

**인터페이스:**
```go
type LogService interface {
    CreateJobLogs(jobName string)
    AddLog(jobName, container, message string)
    GetJobLogs(jobName string) ([]models.LogEntry, bool)
    DeleteJobLogs(jobName string)
}
```

**구현체:**
- `InMemoryLogService` - 메모리 기반 서비스

**의존성:** models, storage, utils

**역할:**
- 저장소와의 상호작용
- 로그 레벨 감지
- 비즈니스 로직 조정

---

### 5. pkg/handlers - HTTP 핸들러

**책임:** HTTP 요청 처리 및 응답

**핸들러:**

#### BuildJobHandler
```go
type BuildJobHandler struct {
    logService services.LogService
}

func (h *BuildJobHandler) Create(w http.ResponseWriter, r *http.Request)
```

- POST /api/buildjob 처리
- 요청 검증
- 로그 초기화

#### LogsHandler
```go
type LogsHandler struct {
    logService services.LogService
}

func (h *LogsHandler) Get(w http.ResponseWriter, r *http.Request)
```

- GET /api/buildjob/{job_name}/logs 처리
- 로그 조회 및 반환

**의존성:** models, services

---

### 6. main.go - 진입점

**책임:** 의존성 주입 및 라우팅 설정

```go
func main() {
    // 의존성 주입
    logService := services.NewInMemoryLogService()
    
    // 핸들러 생성
    jobHandler := handlers.NewBuildJobHandler(logService)
    logsHandler := handlers.NewLogsHandler(logService)
    
    // 라우팅
    http.HandleFunc("/api/buildjob", jobHandler.Create)
    http.HandleFunc("/api/buildjob/", logsHandler.Get)
    
    // 서버 시작
    http.ListenAndServe(":8080", nil)
}
```

---

## 의존성 흐름

```
main.go
  ├─> services.LogService (인터페이스)
  │     └─> storage.LogStorage (인터페이스)
  │           └─> models
  │
  ├─> handlers.BuildJobHandler
  │     ├─> services.LogService
  │     └─> models
  │
  └─> handlers.LogsHandler
        ├─> services.LogService
        └─> models

utils → models (독립적)
```

## 설계 패턴

### 1. Dependency Injection (의존성 주입)

```go
// 명시적 의존성 주입
func NewBuildJobHandler(logService services.LogService) *BuildJobHandler {
    return &BuildJobHandler{
        logService: logService,
    }
}

// 테스트에서 Mock 서비스 주입 가능
mockService := &MockLogService{}
handler := handlers.NewBuildJobHandler(mockService)
```

**이점:**
- 테스트 용이성
- 느슨한 결합
- 명시적 의존성

### 2. Interface-Based Design (인터페이스 기반 설계)

```go
// 구현이 아닌 인터페이스에 의존
type LogService interface {
    CreateJobLogs(jobName string)
    // ...
}

// 다양한 구현 가능
type InMemoryLogService struct { ... }      // Phase 1
type ConfigMapLogService struct { ... }    // Phase 2
type DatabaseLogService struct { ... }     // Phase 3
```

**이점:**
- 구현 교체 가능
- 단위 테스트 Mock 객체 생성 용이
- 향후 기능 확장 용이

### 3. Separation of Concerns (관심사의 분리)

| 패키지 | 책임 | 변경 이유 |
|--------|------|---------|
| models | 데이터 구조 | 요청/응답 형식 변경 |
| handlers | HTTP 처리 | API 엔드포인트 추가/삭제 |
| services | 비즈니스 로직 | 로직 변경 |
| storage | 데이터 저장 | 저장소 변경 |
| utils | 공통 함수 | 유틸리티 추가 |

## 데이터 흐름

### BuildJob 생성 흐름

```
HTTP Request (POST /api/buildjob)
    │
    ├─> handlers.BuildJobHandler.Create()
    │     │
    │     ├─> JSON 디코딩
    │     ├─> 입력값 검증
    │     │
    │     └─> services.LogService.CreateJobLogs()
    │           │
    │           └─> storage.MemoryStorage.SaveLog()
    │                 └─> 메모리에 "job created" 로그 저장
    │
    └─> HTTP Response (201 Created)
```

### 로그 조회 흐름

```
HTTP Request (GET /api/buildjob/{job_name}/logs)
    │
    ├─> handlers.LogsHandler.Get()
    │     │
    │     ├─> Job Name 추출
    │     │
    │     └─> services.LogService.GetJobLogs()
    │           │
    │           └─> storage.MemoryStorage.GetLogs()
    │                 └─> 메모리에서 로그 조회
    │
    └─> HTTP Response (200 OK with logs)
```

## 확장성 계획

### Phase 2: 지속성 저장소
```go
// ConfigMap 기반 저장소
type ConfigMapLogService struct {
    clientset kubernetes.Interface
}

// Kubernetes ConfigMap에 로그 저장
// 장점: 클러스터 내 로그 공유, Pod 재시작 후에도 데이터 유지
```

### Phase 3: 데이터베이스 연동
```go
// PostgreSQL 기반 저장소
type DatabaseLogService struct {
    db *gorm.DB
}

// 장점: 확장 가능한 저장소, 복잡한 쿼리 지원
```

## 테스트 전략

```
┌─ Unit Tests
│   ├─ Handler Tests (httptest 사용)
│   ├─ Service Tests (Mock Storage)
│   └─ Utility Tests
│
└─ Integration Tests
    └─ BuildJobWorkflow (전체 흐름)
```

## 성능 고려사항

### 메모리 저장소 (현재)
- **장점:** 빠른 응답, 간단한 구현
- **단점:** 메모리 제한, Pod 재시작 시 손실
- **적합:** 개발, 테스트, 단기 저장소

### 향후 개선
1. **캐싱 계층 추가** - Redis 활용
2. **비동기 처리** - Worker Queue 도입
3. **로그 압축** - 오래된 로그 아카이브
4. **모니터링** - Prometheus 메트릭 추가
