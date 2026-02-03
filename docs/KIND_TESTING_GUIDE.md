# KIND 클러스터에서 테스트하기

이 가이드는 KIND (Kubernetes in Docker)를 사용하여 api-server v3.0을 배포하고 테스트하는 방법을 설명합니다.

## 필수 요구사항

- Docker Desktop (또는 Docker Engine)
- KIND 설치
- kubectl 설치
- Helm 설치 (v3.x)
- Go 1.23 (로컬 개발용)

## 1단계: KIND 클러스터 생성

### 1-1. 클러스터 생성

```bash
kind create cluster --name api-server-cluster --config kind-config.yaml
```

`kind-config.yaml` 파일의 구성:
```yaml
kind: Cluster
apiVersion: kind.x-k8s.io/v1alpha4
nodes:
  - role: control-plane
    ports:
      - containerPort: 8080
        hostPort: 8080
```

### 1-2. 클러스터 확인

```bash
kubectl cluster-info --context kind-api-server-cluster
kubectl get nodes
```

## 2단계: Docker 이미지 빌드 및 로드

### 2-1. 로컬에서 이미지 빌드

```bash
cd c:\Users\hwiyel\project#1

# Docker 이미지 빌드 (다단계 빌드에서 Go 컴파일 포함)
docker build -t api-server:3.0 .
```

### 2-2. KIND 클러스터에 이미지 로드

```bash
kind load docker-image api-server:3.0 --name api-server-cluster
```

이미지가 클러스터에 로드되었는지 확인:
```bash
kubectl get images --context kind-api-server-cluster
# 또는
docker exec api-server-cluster-control-plane crictl images | grep api-server
```

## 3단계: Helm으로 배포

### 3-1. Helm 차트 배포

```bash
helm install api-server ./helm/api-server \
  --namespace default \
  --set image.tag=3.0 \
  --set image.pullPolicy=Never
```

### 3-2. 배포 상태 확인

```bash
# Pod 상태 확인
kubectl get pods -n default
kubectl describe pod -l app=api-server -n default

# Service 상태 확인
kubectl get svc -n default
kubectl get endpoints -n default

# Logs 확인
kubectl logs -l app=api-server -n default
```

## 4단계: API 테스트

### 4-1. 포트 포워딩 설정

```bash
kubectl port-forward svc/api-server 8080:8080 -n default
```

### 4-2. API 엔드포인트 테스트

**BuildJob 생성 (POST /api/buildjob)**

```bash
curl -X POST http://localhost:8080/api/buildjob \
  -H "Content-Type: application/json" \
  -d '{
    "job_name": "test-build-001",
    "dockerfile_content": "FROM alpine\nRUN echo Hello World"
  }'
```

응답 예시:
```json
{
  "status": "created",
  "message": "Build job created successfully",
  "job_name": "test-build-001",
  "job_id": "build-test-build-001-001",
  "namespace": "default",
  "created_at": "2026-02-04T10:30:45Z"
}
```

**로그 조회 (GET /api/buildjob/{job_name}/logs)**

```bash
curl -X GET http://localhost:8080/api/buildjob/test-build-001/logs \
  -H "Content-Type: application/json"
```

응답 예시:
```json
{
  "job_name": "test-build-001",
  "status": "running",
  "logs": [
    {
      "timestamp": "2026-02-04T10:30:45Z",
      "container": "system",
      "message": "Build job created successfully",
      "level": "info"
    }
  ],
  "total_lines": 1
}
```

### 4-3. PowerShell에서 테스트 (Windows)

```powershell
# BuildJob 생성
$body = @{
    job_name = "powershell-test-001"
    dockerfile_content = "FROM alpine`nRUN apk add curl"
} | ConvertTo-Json

Invoke-WebRequest -Uri "http://localhost:8080/api/buildjob" `
  -Method POST `
  -ContentType "application/json" `
  -Body $body

# 로그 조회
Invoke-WebRequest -Uri "http://localhost:8080/api/buildjob/powershell-test-001/logs" `
  -Method GET `
  -ContentType "application/json" | ForEach-Object { $_.Content | ConvertFrom-Json | ConvertTo-Json }
```

## 5단계: 내부 테스트 (클러스터 내부)

### 5-1. 클러스터 내부에서 API 호출

```bash
# 테스트용 Pod 생성
kubectl run -it debug --image=alpine --restart=Never -- sh

# 클러스터 내부에서 API 호출
apk add curl
curl -X POST http://api-server:8080/api/buildjob \
  -H "Content-Type: application/json" \
  -d '{"job_name":"internal-test","dockerfile_content":"FROM alpine"}'
```

## 6단계: 배포 업데이트

이미지를 변경하고 다시 배포하려면:

```bash
# 1. 코드 수정
# 2. 로컬 테스트
go test ./...

# 3. Docker 이미지 빌드
docker build -t api-server:3.0 .

# 4. KIND에 로드
kind load docker-image api-server:3.0 --name api-server-cluster

# 5. Pod 재시작 (새 이미지 사용)
kubectl rollout restart deployment/api-server -n default

# 6. 배포 상태 확인
kubectl rollout status deployment/api-server -n default
```

## 7단계: 정리 및 삭제

### 7-1. Helm 릴리스 삭제

```bash
helm uninstall api-server -n default
```

### 7-2. KIND 클러스터 삭제

```bash
kind delete cluster --name api-server-cluster
```

## 문제 해결

### Pod가 시작되지 않음

```bash
# Pod 상세 정보 확인
kubectl describe pod <pod-name> -n default

# Pod 로그 확인
kubectl logs <pod-name> -n default
```

### 이미지를 찾을 수 없음

```bash
# 이미지가 클러스터에 로드되었는지 확인
kind load docker-image api-server:3.0 --name api-server-cluster

# 확인
docker exec api-server-cluster-control-plane crictl images
```

### Service에 접근할 수 없음

```bash
# Service IP와 엔드포인트 확인
kubectl get svc api-server -n default
kubectl get endpoints api-server -n default

# 포트 포워딩 재설정
kubectl port-forward svc/api-server 8080:8080 -n default
```

## Skaffold를 사용한 자동 배포 (선택사항)

Skaffold를 사용하면 코드 변경 시 자동으로 빌드, 로드, 배포합니다:

```bash
# 개발 모드 (hot reload)
skaffold dev

# 한 번만 배포
skaffold run
```

## 참고 자료

- [KIND 공식 문서](https://kind.sigs.k8s.io/)
- [Kubernetes 공식 문서](https://kubernetes.io/docs/)
- [Helm 공식 문서](https://helm.sh/docs/)
- [Skaffold 공식 문서](https://skaffold.dev/)
