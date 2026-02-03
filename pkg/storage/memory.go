package storage

import (
	"api-server/pkg/models"
	"sync"
	"time"
)

// MemoryStorage는 메모리 기반 로그 저장소 구현입니다
type MemoryStorage struct {
	mu    sync.RWMutex
	logs  map[string][]models.LogEntry
}

// NewMemoryStorage는 새로운 메모리 저장소를 생성합니다
func NewMemoryStorage() LogStorage {
	return &MemoryStorage{
		logs: make(map[string][]models.LogEntry),
	}
}

// SaveLog는 새로운 로그를 저장합니다
func (s *MemoryStorage) SaveLog(jobName, container, message, level string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.logs[jobName]; !exists {
		s.logs[jobName] = []models.LogEntry{}
	}

	entry := models.LogEntry{
		Timestamp: time.Now().Format(time.RFC3339),
		Container: container,
		Message:   message,
		Level:     level,
	}

	s.logs[jobName] = append(s.logs[jobName], entry)
}

// GetLogs는 특정 Job의 모든 로그를 조회합니다
func (s *MemoryStorage) GetLogs(jobName string) ([]models.LogEntry, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	logs, exists := s.logs[jobName]
	return logs, exists
}

// DeleteLogs는 특정 Job의 로그를 삭제합니다
func (s *MemoryStorage) DeleteLogs(jobName string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.logs, jobName)
}

// Exists는 특정 Job의 로그가 존재하는지 확인합니다
func (s *MemoryStorage) Exists(jobName string) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()

	_, exists := s.logs[jobName]
	return exists
}
