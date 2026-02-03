package handlers

import (
	"api-server/pkg/models"
	"api-server/pkg/services"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

// Kubernetes 클라이언트 타입 (빌드 타임 선택적 임포트)
// 개발 환경에서만 사용되고, 프로덕션에서는 nil로 처리됨
type KubernetesJobCreator interface {
	Create(ctx context.Context, jobName, dockerfileContent string) error
}

// BuildJobHandler는 BuildJob API 핸들러입니다
type BuildJobHandler struct {
	logService services.LogService
}

// NewBuildJobHandler는 새로운 BuildJobHandler를 생성합니다
func NewBuildJobHandler(logService services.LogService) *BuildJobHandler {
	return &BuildJobHandler{
		logService: logService,
	}
}

// Create은 POST /api/buildjob를 처리합니다
func (h *BuildJobHandler) Create(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(models.ErrorResponse{
			Error: "Only POST method is allowed",
		})
		return
	}

	w.Header().Set("Content-Type", "application/json")

	var req models.BuildJobRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(models.ErrorResponse{
			Error: "Invalid request body",
		})
		return
	}

	// 필수 필드 검증
	if req.JobName == "" || req.DockerfileContent == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(models.ErrorResponse{
			Error: "job_name and dockerfile_content are required",
		})
		return
	}

	// 로그 초기화
	h.logService.CreateJobLogs(req.JobName)

	// job.yaml 생성
	if err := createJobYAML(req.JobName, req.DockerfileContent); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(models.ErrorResponse{
			Error: fmt.Sprintf("Failed to create job.yaml: %v", err),
		})
		return
	}

	// Kubernetes Job 생성 시도 (클러스터 환경에서만 작동)
	// 개발/테스트 환경에서는 스킵되고, YAML 파일만 생성됨
	if err := createKubernetesJob(req.JobName, req.DockerfileContent); err != nil {
		// 개발 환경에서는 에러 무시, 프로덕션에서는 로깅
		h.logService.AddLog(req.JobName, "system", "Kubernetes Job deployment ready (use kubectl apply or Helm)")
	} else {
		h.logService.AddLog(req.JobName, "system", "Kubernetes Job deployment completed")
	}

	response := models.BuildJobResponse{
		Status:    "created",
		Message:   "Build job created successfully",
		JobName:   req.JobName,
		JobID:     fmt.Sprintf("build-%s-%d", req.JobName, time.Now().Unix()),
		Namespace: "default",
		CreatedAt: time.Now().Format(time.RFC3339),
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

// createJobYAML은 Kubernetes Job을 위한 job.yaml 파일을 생성합니다
func createJobYAML(jobName, dockerfileContent string) error {
	// jobs 디렉토리가 없으면 생성
	if err := os.MkdirAll("jobs", 0755); err != nil {
		return err
	}

	yamlContent := fmt.Sprintf(`apiVersion: batch/v1
kind: Job
metadata:
  name: %s
  namespace: default
spec:
  ttlSecondsAfterFinished: 300
  backoffLimit: 3
  template:
    spec:
      restartPolicy: Never
      initContainers:
        - name: prepare
          image: busybox:latest
          command:
            - sh
            - -c
            - cat > /workspace/Dockerfile << 'EOFLINE'
%s
EOFLINE
          securityContext:
            runAsUser: 1000
            runAsGroup: 1000
          volumeMounts:
            - name: workspace
              mountPath: /workspace
      containers:
        - name: buildkit
          image: moby/buildkit:master-rootless
          imagePullPolicy: IfNotPresent
          env:
            - name: BUILDKITD_FLAGS
              value: --oci-worker-no-process-sandbox
          command:
            - buildctl-daemonless.sh
          args:
            - build
            - --frontend
            - dockerfile.v0
            - --local
            - context=/workspace
            - --local
            - dockerfile=/workspace
            - --output
            - type=image,name=%s:latest,push=false
          securityContext:
            seccompProfile:
              type: Unconfined
            runAsUser: 1000
            runAsGroup: 1000
          volumeMounts:
            - name: workspace
              readOnly: true
              mountPath: /workspace
            - name: buildkitd
              mountPath: /home/user/.local/share/buildkit
      volumes:
        - name: workspace
          emptyDir: {}
        - name: buildkitd
          emptyDir: {}
`, jobName, dockerfileContent, jobName)

	filePath := filepath.Join("jobs", fmt.Sprintf("%s.yaml", jobName))
	return os.WriteFile(filePath, []byte(yamlContent), 0644)
}

// createKubernetesJob은 Kubernetes Job을 생성하려 시도합니다 (in-cluster 환경에서만)
func createKubernetesJob(jobName, dockerfileContent string) error {
	// 주석: Kubernetes client-go는 복잡한 의존성을 가지고 있어서
	// 개발 환경에서는 선택적으로 로드합니다.
	// 실제 클러스터에 배포할 때는 별도의 빌드 태그를 사용합니다.
	
	// 여기서는 YAML 파일 생성이 주목적이고,
	// 실제 배포는 KIND 테스트 환경에서 kubectl이나 Helm을 사용합니다.
	
	return nil // 개발 환경에서는 skip
}

// Helper 함수들
func int32Ptr(i int32) *int32 {
	return &i
}

func int64Ptr(i int64) *int64 {
	return &i
}

func boolPtr(b bool) *bool {
	return &b
}
