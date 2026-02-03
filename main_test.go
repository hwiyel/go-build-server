package main

import (
	"api-server/pkg/handlers"
	"api-server/pkg/models"
	"api-server/pkg/services"
	"api-server/pkg/utils"
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

// === BuildJob Handler 테스트 ===

func TestCreateBuildJob(t *testing.T) {
	logService := services.NewInMemoryLogService()
	handler := handlers.NewBuildJobHandler(logService)

	payload := models.BuildJobRequest{
		JobName:           "test-build-job",
		DockerfileContent: "FROM alpine\nRUN apk add gcc",
	}

	body, _ := json.Marshal(payload)
	req, _ := http.NewRequest("POST", "/api/buildjob", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	handler.Create(rr, req)

	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusCreated)
	}

	var response models.BuildJobResponse
	json.NewDecoder(rr.Body).Decode(&response)

	if response.JobName != "test-build-job" {
		t.Errorf("handler returned wrong job name: got %v want %v", response.JobName, "test-build-job")
	}
}

// === Logs Handler 테스트 ===

func TestGetLogs(t *testing.T) {
	logService := services.NewInMemoryLogService()
	logService.CreateJobLogs("test-job")

	handler := handlers.NewLogsHandler(logService)
	req, _ := http.NewRequest("GET", "/api/buildjob/test-job/logs", nil)
	rr := httptest.NewRecorder()

	handler.Get(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	var response models.LogsResponse
	json.NewDecoder(rr.Body).Decode(&response)

	if len(response.Logs) == 0 {
		t.Error("expected logs but got none")
	}

	if response.Logs[0].Level != "info" {
		t.Errorf("expected log level 'info' but got %v", response.Logs[0].Level)
	}
}

func TestGetLogsNotFound(t *testing.T) {
	logService := services.NewInMemoryLogService()
	handler := handlers.NewLogsHandler(logService)

	req, _ := http.NewRequest("GET", "/api/buildjob/nonexistent-job/logs", nil)
	rr := httptest.NewRecorder()

	handler.Get(rr, req)

	if status := rr.Code; status != http.StatusNotFound {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusNotFound)
	}
}

// === Utility 테스트 ===

func TestDetectLogLevel(t *testing.T) {
	tests := []struct {
		message  string
		expected string
	}{
		{"Build succeeded", "info"},
		{"Warning: deprecated feature", "warn"},
		{"Error: failed to build image", "error"},
		{"ERROR in dockerfile", "error"},
	}

	for _, tt := range tests {
		result := utils.DetectLogLevel(tt.message)
		if result != tt.expected {
			t.Errorf("DetectLogLevel(%q) = %v, want %v", tt.message, result, tt.expected)
		}
	}
}

// === LogService 테스트 ===

func TestLogService(t *testing.T) {
	logService := services.NewInMemoryLogService()

	// Job 로그 생성
	logService.CreateJobLogs("test-job")
	logs, exists := logService.GetJobLogs("test-job")

	if !exists {
		t.Error("job logs should exist after creation")
	}

	if len(logs) == 0 {
		t.Error("job logs should have initial log entry")
	}

	// 새로운 로그 추가
	logService.AddLog("test-job", "buildkit", "Building image...")
	logs, _ = logService.GetJobLogs("test-job")

	if len(logs) != 2 {
		t.Errorf("expected 2 logs, got %d", len(logs))
	}

	if logs[1].Level != "info" {
		t.Errorf("expected info level for build message, got %v", logs[1].Level)
	}

	// 에러 로그 추가
	logService.AddLog("test-job", "buildkit", "Error: failed to compile")
	logs, _ = logService.GetJobLogs("test-job")

	if logs[len(logs)-1].Level != "error" {
		t.Errorf("expected error level for failure message, got %v", logs[len(logs)-1].Level)
	}
}

func TestLogServiceDeleteJob(t *testing.T) {
	logService := services.NewInMemoryLogService()

	logService.CreateJobLogs("delete-test-job")
	logService.AddLog("delete-test-job", "builder", "Test log")

	// 로그가 존재하는지 확인
	_, exists := logService.GetJobLogs("delete-test-job")
	if !exists {
		t.Error("job logs should exist before deletion")
	}

	// 로그 삭제
	logService.DeleteJobLogs("delete-test-job")

	// 로그가 삭제되었는지 확인
	_, exists = logService.GetJobLogs("delete-test-job")
	if exists {
		t.Error("job logs should not exist after deletion")
	}
}

// === Integration 테스트 ===

func TestBuildJobWorkflow(t *testing.T) {
	logService := services.NewInMemoryLogService()
	jobHandler := handlers.NewBuildJobHandler(logService)
	logsHandler := handlers.NewLogsHandler(logService)

	// 1. BuildJob 생성
	jobPayload := models.BuildJobRequest{
		JobName:           "workflow-test-job",
		DockerfileContent: "FROM alpine\nRUN echo test",
	}

	body, _ := json.Marshal(jobPayload)
	jobReq, _ := http.NewRequest("POST", "/api/buildjob", bytes.NewReader(body))
	jobRR := httptest.NewRecorder()
	jobHandler.Create(jobRR, jobReq)

	if jobRR.Code != http.StatusCreated {
		t.Errorf("failed to create build job: got status %d", jobRR.Code)
	}

	// 2. 로그 추가
	logService.AddLog("workflow-test-job", "builder", "Build step 1")
	logService.AddLog("workflow-test-job", "builder", "Build step 2 completed")

	// 3. 로그 조회
	logsReq, _ := http.NewRequest("GET", "/api/buildjob/workflow-test-job/logs", nil)
	logsRR := httptest.NewRecorder()
	logsHandler.Get(logsRR, logsReq)

	if logsRR.Code != http.StatusOK {
		t.Errorf("failed to get logs: got status %d", logsRR.Code)
	}

	var logsResponse models.LogsResponse
	json.NewDecoder(logsRR.Body).Decode(&logsResponse)

	// CreateJobLogs는 자동으로 시스템 로그를 추가하므로 3개 (system + 2 user logs)
	if logsResponse.TotalLines != 3 {
		t.Errorf("expected 3 logs, got %d", logsResponse.TotalLines)
	}
}
