package services

import (
	"api-server/pkg/models"
	"api-server/pkg/storage"
	"api-server/pkg/utils"
)

// LogService는 로그 관련 비즈니스 로직을 담당합니다
type LogService interface {
	// CreateJobLogs는 새로운 Job의 로그를 초기화합니다
	CreateJobLogs(jobName string)

	// AddLog는 로그 엔트리를 추가합니다
	AddLog(jobName, container, message string)

	// GetJobLogs는 특정 Job의 모든 로그를 조회합니다
	GetJobLogs(jobName string) ([]models.LogEntry, bool)

	// DeleteJobLogs는 특정 Job의 로그를 삭제합니다
	DeleteJobLogs(jobName string)
}

// InMemoryLogService는 메모리 기반 로그 서비스 구현입니다
type InMemoryLogService struct {
	storage storage.LogStorage
}

// NewInMemoryLogService는 새로운 메모리 기반 로그 서비스를 생성합니다
func NewInMemoryLogService() LogService {
	return &InMemoryLogService{
		storage: storage.NewMemoryStorage(),
	}
}

// CreateJobLogs는 새로운 Job의 로그를 초기화합니다
func (s *InMemoryLogService) CreateJobLogs(jobName string) {
	s.storage.SaveLog(jobName, "system", "Build job created successfully", "info")
}

// AddLog는 로그 엔트리를 추가합니다
func (s *InMemoryLogService) AddLog(jobName, container, message string) {
	level := utils.DetectLogLevel(message)
	s.storage.SaveLog(jobName, container, message, level)
}

// GetJobLogs는 특정 Job의 모든 로그를 조회합니다
func (s *InMemoryLogService) GetJobLogs(jobName string) ([]models.LogEntry, bool) {
	return s.storage.GetLogs(jobName)
}

// DeleteJobLogs는 특정 Job의 로그를 삭제합니다
func (s *InMemoryLogService) DeleteJobLogs(jobName string) {
	s.storage.DeleteLogs(jobName)
}
