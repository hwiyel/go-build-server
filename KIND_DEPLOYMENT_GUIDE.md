# KinD 환경 배포 가이드

kind(Kubernetes IN Docker) 환경에서 API 서버를 배포하기 위한 가이드입니다.

## 전제 조건

- Docker 설치
- kind 설치
- kubectl 설치
- Helm 설치

## KinD 클러스터 생성

```bash
# 기본 클러스터 생성
kind create cluster --name api-server-cluster

# 포트 매핑을 포함한 클러스터 생성 (localhost:80 -> 80)
cat > kind-config.yaml << EOF
kind: Cluster
apiVersion: kind.x-k8s.io/v1alpha4
nodes:
  - role: control-plane
    extraPortMappings:
      - containerPort: 80
        hostPort: 80
        protocol: TCP
      - containerPort: 443
        hostPort: 443
        protocol: TCP
EOF

kind create cluster --name api-server-cluster --config kind-config.yaml
```

## Docker 이미지를 KinD로 로드

```bash
# 로컬 Docker 이미지를 kind 클러스터에 로드
kind load docker-image api-server:1.1 --name api-server-cluster
```

## Helm 차트로 배포

### 1. 기본 배포
```bash
helm install api-server ./helm/api-server -n default
```

### 2. 커스텀 값으로 배포
```bash
# 복제 수 변경
helm install api-server ./helm/api-server \
  --set replicaCount=3 \
  -n default

# 리소스 요청 변경
helm install api-server ./helm/api-server \
  --set resources.requests.cpu=200m \
  --set resources.limits.cpu=500m \
  -n default

# 여러 값 변경
helm install api-server ./helm/api-server \
  --values custom-values.yaml \
  -n default
```

### 3. 특정 namespace에 배포
```bash
kubectl create namespace api-server
helm install api-server ./helm/api-server -n api-server
```

## 배포 확인

```bash
# 릴리스 확인
helm list -n default

# 생성된 리소스 확인
kubectl get all -n default

# Pod 로그 확인
kubectl logs -l app.kubernetes.io/name=api-server -n default

# Pod 접근
kubectl port-forward service/api-server 8080:8080 -n default
```

## API 테스트

```bash
# 포트 포워딩 설정 후
kubectl port-forward service/api-server 8080:8080

# GET 요청 테스트
curl http://localhost:8080/api/build

# POST 요청 테스트
curl -X POST http://localhost:8080/api/build/create \
  -H "Content-Type: application/json" \
  -d '{
    "project_name": "test-app",
    "environment": "production"
  }'
```

## Helm 업데이트 및 롤백

```bash
# 값 업데이트
helm upgrade api-server ./helm/api-server \
  --set replicaCount=2 \
  -n default

# 버전 확인
helm history api-server -n default

# 이전 버전으로 롤백
helm rollback api-server 1 -n default

# 릴리스 제거
helm uninstall api-server -n default
```

## KinD 클러스터 정리

```bash
# 클러스터 삭제
kind delete cluster --name api-server-cluster

# 모든 kind 클러스터 확인
kind get clusters
```

## 커스텀 values.yaml 예시

```yaml
# custom-values.yaml
replicaCount: 3

image:
  tag: "1.1"
  pullPolicy: IfNotPresent

service:
  type: ClusterIP
  port: 8080

resources:
  requests:
    cpu: 100m
    memory: 64Mi
  limits:
    cpu: 200m
    memory: 128Mi

ingress:
  enabled: true
  className: nginx
  hosts:
    - host: api-server.local
      paths:
        - path: /
          pathType: Prefix
```

## 문제 해결

### Pod가 Pending 상태인 경우
```bash
kubectl describe pod <pod-name> -n default
```

### 이미지를 찾을 수 없는 경우
```bash
# 이미지 로드 확인
docker images | grep api-server

# kind에 이미지 로드
kind load docker-image api-server:1.1 --name api-server-cluster
```

### 서비스 접근이 안 되는 경우
```bash
# 포트 포워딩 사용
kubectl port-forward service/api-server 8080:8080

# 또는 Ingress 활성화
helm upgrade api-server ./helm/api-server \
  --set ingress.enabled=true \
  -n default
```
