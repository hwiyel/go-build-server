package handlers

import (
	"api-server/pkg/models"
	"api-server/pkg/services"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

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

	response := models.BuildJobResponse{
		Status:    "created",
		Message:   "Build job created successfully",
		JobName:   req.JobName,
		JobID:     "build-" + req.JobName + "-001",
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
    metadata:
      name: %s
    spec:
      serviceAccountName: default
      restartPolicy: Never
      containers:
      - name: builder
        image: docker:latest
        imagePullPolicy: IfNotPresent
        securityContext:
          privileged: true
        workingDir: /workspace
        command:
          - sh
          - -c
          - |
            cat > Dockerfile << 'EOFLINE'
%s
EOFLINE
            docker build -t %s:latest -f Dockerfile .
        volumeMounts:
        - name: workspace
          mountPath: /workspace
        - name: docker-sock
          mountPath: /var/run/docker.sock
      volumes:
      - name: workspace
        emptyDir: {}
      - name: docker-sock
        hostPath:
          path: /var/run/docker.sock
`, jobName, jobName, dockerfileContent, jobName)

	filePath := filepath.Join("jobs", fmt.Sprintf("%s.yaml", jobName))
	return os.WriteFile(filePath, []byte(yamlContent), 0644)
}

