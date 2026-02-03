# 빠른 시작 가이드

이 문서는 프로젝트를 빠르게 시작하고 실행하는 방법을 설명합니다.

## 전제 조건

- Docker Desktop (Windows)
- WSL 2
- kubectl
- Helm 3.x
- kind

## 로컬 실행 (Go 직접)

```bash
# 의존성 설치
go mod download

# 서버 실행
go run main.go

# 다른 터미널에서 테스트
curl http://localhost:8080/api/build
```

## Docker로 실행

```bash
# 이미지 빌드
docker build -t api-server:1.3 .

# 컨테이너 실행
docker run -p 8080:8080 api-server:1.3

# 테스트
curl http://localhost:8080/api/build
```

## KIND 클러스터에 배포

### WSL 터미널에서 실행

```bash
# 1. 클러스터 생성
kind create cluster --name api-server-cluster --config kind-config.yaml

# 2. 이미지 로드
kind load docker-image api-server:1.3 --name api-server-cluster

# 3. Helm 배포
helm install api-server ./helm/api-server -n default

# 4. 배포 확인
kubectl get pods -n default

# 5. 포트 포워딩
kubectl port-forward service/api-server 8080:8080 -n default
```

### 별도 터미널에서 API 테스트

```bash
# GET 요청
curl http://localhost:8080/api/build

# POST 요청
curl -X POST http://localhost:8080/api/build/create \
  -H "Content-Type: application/json" \
  -d '{
    "project_name": "my-project",
    "environment": "production"
  }'
```

## 테스트 실행

```bash
# WSL 터미널에서
cd /mnt/c/Users/hwiyel/project#1

# 모든 테스트 실행
go test -v

# 커버리지 확인
go test -cover
```

## Helm 관리

```bash
# 릴리스 확인
helm list -n default

# 업그레이드
helm upgrade api-server ./helm/api-server --set replicaCount=3 -n default

# 롤백
helm rollback api-server 1 -n default

# 제거
helm uninstall api-server -n default
```

## 클러스터 정리

```bash
# KIND 클러스터 삭제
kind delete cluster --name api-server-cluster

# 모든 클러스터 확인
kind get clusters
```

## 문제 해결

### Pod가 CrashLoopBackOff 상태

```bash
# Pod 로그 확인
kubectl logs <pod-name> -n default

# Pod 상세 정보 확인
kubectl describe pod <pod-name> -n default
```

### 이미지를 찾을 수 없음

```bash
# KIND에 이미지 재로드
kind load docker-image api-server:1.3 --name api-server-cluster

# 로드된 이미지 확인
docker exec -it api-server-cluster-control-plane crictl images
```

### 포트가 이미 사용 중

```bash
# 다른 포트로 포워딩
kubectl port-forward service/api-server 9080:8080 -n default
curl http://localhost:9080/api/build
```

## 다음 단계

- [개발 기록 보기](./DEVELOPMENT_LOG.md)
- [KIND 배포 가이드 보기](../KIND_DEPLOYMENT_GUIDE.md)
- [API 엔드포인트 문서](./API.md)
