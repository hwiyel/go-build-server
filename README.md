# API Server - Go Project

기본적인 Golang API 서버 프로젝트입니다.

## 프로젝트 구조

```
.
├── main.go              # 메인 API 서버
├── main_test.go         # 테스트 케이스
├── go.mod               # Go 모듈 정의
├── Dockerfile           # Docker 이미지 빌드 정의
├── docker-compose.yml   # Docker Compose 개발 환경
└── .dockerignore        # Docker 빌드 제외 파일
```

## 기능

- `/api/build` - API 엔드포인트 (GET)
  - 응답: JSON 형식으로 status, message, version 반환

## 로컬 실행

### Go 직접 실행
```bash
go run main.go
```
서버가 `http://localhost:8080`에서 시작됩니다.

### Docker 실행
```bash
docker build -t api-server:1.0 .
docker run -p 8080:8080 api-server:1.0
```

### Docker Compose 실행
```bash
docker-compose up
```
개발 환경에서 서버가 `http://localhost:8080`에서 시작됩니다.

## 테스트

### 단위 테스트 실행
```bash
go test -v
```

### 테스트 커버리지 확인
```bash
go test -cover
```

### 상세 커버리지 리포트
```bash
go test -coverprofile=coverage.out
go tool cover -html=coverage.out
```

## API 테스트

### cURL을 사용한 테스트
```bash
curl http://localhost:8080/api/build
```

### 응답 예시
```json
{
  "status": "success",
  "message": "Build API is working!",
  "version": "1.0.0"
}
```

## 개발 명령어

### 테스트 실행 및 빌드
```bash
go test -v && go build -o api-server main.go
```

### Docker 컨테이너 관리
```bash
# 컨테이너 시작
docker-compose up -d

# 로그 확인
docker-compose logs -f

# 컨테이너 중지
docker-compose down
```

## 버전
- Go: 1.23
- API Version: 1.0.0
