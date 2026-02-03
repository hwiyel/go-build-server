package handlers

import (
	"api-server/pkg/models"
	"api-server/pkg/services"
	"encoding/json"
	"net/http"
	"strings"
)

// LogsHandler는 로그 조회 핸들러입니다
type LogsHandler struct {
	logService services.LogService
}

// NewLogsHandler는 새로운 LogsHandler를 생성합니다
func NewLogsHandler(logService services.LogService) *LogsHandler {
	return &LogsHandler{
		logService: logService,
	}
}

// Get은 GET /api/buildjob/{job_name}/logs를 처리합니다
func (h *LogsHandler) Get(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(models.ErrorResponse{
			Error: "Only GET method is allowed",
		})
		return
	}

	jobName := strings.TrimPrefix(r.URL.Path, "/api/buildjob/")
	jobName = strings.TrimSuffix(jobName, "/logs")

	w.Header().Set("Content-Type", "application/json")

	logs, exists := h.logService.GetJobLogs(jobName)
	if !exists {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(models.ErrorResponse{
			Error: "Job not found",
		})
		return
	}

	response := models.LogsResponse{
		JobName:    jobName,
		Status:     "running",
		Logs:       logs,
		TotalLines: len(logs),
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
