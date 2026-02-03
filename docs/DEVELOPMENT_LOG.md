# API Server 개발 기록

## 프로젝트 개요

Go 기반의 간단한 REST API 서버를 개발하고, Docker 및 Kubernetes(KIND)를 통해 컨테이너화하는 프로젝트입니다.

## 개발 단계별 기록

### 1단계: 기본 API 서버 구현 (v1.0)

**작업 내용:**
- Go HTTP 서버 기본 구현
- GET `/api/build` 엔드포인트 추가
- JSON 응답 포맷 정의 (status, message, version)

**주요 파일:**
- `main.go` - API 서버 메인 코드
- `go.mod` - Go 모듈 정의

**결과:**
```
GET /api/build
응답: {"status":"success","message":"Build API is working!","version":"1.0.0"}
```

---

### 2단계: Docker 컨테이너화 (v1.0)

**작업 내용:**
- 멀티 스테이지 빌드 Dockerfile 작성
- Builder 이미지와 Runtime 이미지 분리
- 최적화된 이미지 크기 (alpine 기반)

**주요 파일:**
- `Dockerfile` - 컨테이너 이미지 정의
- `.dockerignore` - 빌드 제외 파일 정의

**명령어:**
```bash
docker build -t api-server:1.0 .
docker run -p 8080:8080 api-server:1.0
```

---

### 3단계: 테스트 환경 세팅

**작업 내용:**
- Go 단위 테스트 작성 (2개)
- HTTP 핸들러 테스트
- 응답 검증 테스트

**테스트 케이스:**
1. `TestBuildHandler` - 기본 응답 검증
2. `TestBuildHandlerResponseStructure` - 응답 구조 검증

**테스트 실행:**
```bash
go test -v
go test -cover
```

---

### 4단계: POST API 추가 (v1.1)

**작업 내용:**
- POST `/api/build/create` 엔드포인트 추가
- 요청 본문 검증 (projectName, environment)
- 상태 코드별 응답 처리 (201 Created, 400 Bad Request, 405 Method Not Allowed)

**주요 기능:**
- JSON 요청 파싱
- 필드 검증
- 빌드 ID 생성

**테스트 추가:**
- `TestBuildCreateHandler` - 성공 케이스
- `TestBuildCreateHandlerMissingFields` - 필드 누락 검증
- `TestBuildCreateHandlerInvalidMethod` - 메서드 검증

---

### 5단계: Helm Chart 구성

**작업 내용:**
- Helm 차트 완전 구성
- Kubernetes 리소스 템플릿 작성

**생성 파일:**
- `helm/api-server/Chart.yaml` - 차트 메타데이터
- `helm/api-server/values.yaml` - 기본 설정값
- `helm/api-server/templates/deployment.yaml` - Deployment
- `helm/api-server/templates/service.yaml` - Service
- `helm/api-server/templates/serviceaccount.yaml` - ServiceAccount
- `helm/api-server/templates/ingress.yaml` - Ingress
- `helm/api-server/templates/_helpers.tpl` - 헬퍼 함수

**주요 기능:**
- 복제본 수 설정 (기본 2개)
- 리소스 제한 (100m CPU, 64Mi 메모리)
- Health Check (Liveness/Readiness Probe)
- 자동 스케일링 지원

---

### 6단계: Kubernetes (KIND) 배포

**작업 내용:**
- KIND 클러스터 설정 파일 작성
- 배포 가이드 문서화
- 포트 매핑 설정 (80, 443, 8080)

**생성 파일:**
- `kind-config.yaml` - KIND 클러스터 설정
- `KIND_DEPLOYMENT_GUIDE.md` - 배포 가이드

**배포 명령:**
```bash
kind create cluster --name api-server-cluster --config kind-config.yaml
kind load docker-image api-server:1.1 --name api-server-cluster
helm install api-server ./helm/api-server -n default
```

---

### 7단계: Dockerfile 최적화 (v1.2 → v1.3)

**이슈:**
```
permission denied: unknown
unable to start container process: exec: "./api-server": start ./api-server
```

**해결 방안:**
1. 실행 권한 추가 (`chmod +x`)
2. 절대 경로 사용 (`/app/api-server`)
3. `libc6-compat` 라이브러리 추가 (glibc 호환성)
4. `ENTRYPOINT` 사용 (CMD 대신)
5. 빌드 플래그 최적화 (`-a -installsuffix cgo`)

**최종 Dockerfile:**
```dockerfile
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o api-server main.go
COPY --from=builder /app/api-server /app/api-server
RUN chmod +x /app/api-server
ENTRYPOINT ["/app/api-server"]
```

---

### 8단계: API 테스트

**환경:** WSL의 KIND 클러스터

**포트 포워딩:**
```bash
kubectl port-forward service/api-server 8080:8080 -n default
```

**테스트 명령어:**
```bash
# GET 테스트
curl http://localhost:8080/api/build

# POST 테스트
curl -X POST http://localhost:8080/api/build/create \
  -H "Content-Type: application/json" \
  -d '{"project_name":"test-project","environment":"production"}'
```

**결과:** ✅ 모든 API 정상 작동

---

## 최종 프로젝트 구조

```
project#1/
├── main.go                      # API 서버 메인 코드
├── main_test.go                 # 단위 테스트 (5개)
├── go.mod                        # Go 모듈
├── Dockerfile                    # 컨테이너 이미지 (v1.3)
├── docker-compose.yml            # 개발 환경
├── kind-config.yaml              # KIND 클러스터 설정
├── README.md                      # 프로젝트 문서
├── KIND_DEPLOYMENT_GUIDE.md       # 배포 가이드
├── docs/                          # 개발 기록
│   └── DEVELOPMENT_LOG.md        # 이 파일
└── helm/
    └── api-server/
        ├── Chart.yaml
        ├── values.yaml
        └── templates/
            ├── deployment.yaml
            ├── service.yaml
            ├── serviceaccount.yaml
            ├── ingress.yaml
            └── _helpers.tpl
```

## 기술 스택

| 항목 | 버전 |
|------|------|
| Go | 1.23 |
| Docker | Latest |
| Kubernetes (KIND) | Latest |
| Helm | 3.x |
| Alpine Linux | Latest |

## 주요 학습점

1. **Go HTTP Server**: 기본 HTTP 핸들러 작성 및 JSON 응답 처리
2. **Docker**: 멀티 스테이지 빌드, 보안 및 최적화
3. **Kubernetes**: Deployment, Service, Pod 관리
4. **Helm**: 차트 작성 및 템플릿 활용
5. **Testing**: Go 단위 테스트 작성 및 httptest 활용
6. **Linux/WSL**: 컨테이너 환경에서의 권한 및 경로 관리

## 완료 체크리스트

- ✅ Go API 서버 (GET/POST)
- ✅ Docker 컨테이너화
- ✅ 단위 테스트 작성
- ✅ Helm Chart 구성
- ✅ KIND 클러스터 배포
- ✅ API 테스트 완료
- ✅ 문서화

## 향후 개선 사항

- [ ] 데이터베이스 연동 (PostgreSQL)
- [ ] 로깅 시스템 (structured logging)
- [ ] CI/CD 파이프라인 (GitHub Actions)
- [ ] Prometheus 메트릭
- [ ] API 문서 (Swagger/OpenAPI)
- [ ] 멀티 환경 Helm values (dev, staging, prod)
