package storage

import "api-server/pkg/models"

// LogStorage는 로그 저장소 인터페이스입니다
type LogStorage interface {
	// SaveLog는 새로운 로그를 저장합니다
	SaveLog(jobName, container, message, level string)

	// GetLogs는 특정 Job의 모든 로그를 조회합니다
	GetLogs(jobName string) ([]models.LogEntry, bool)

	// DeleteLogs는 특정 Job의 로그를 삭제합니다
	DeleteLogs(jobName string)

	// Exists는 특정 Job의 로그가 존재하는지 확인합니다
	Exists(jobName string) bool
}
