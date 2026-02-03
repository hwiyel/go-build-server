package handlers

import (
	"api-server/pkg/models"
	"api-server/pkg/services"
	"encoding/json"
	"net/http"
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
