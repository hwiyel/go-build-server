package utils

import "strings"

// DetectLogLevel은 메시지 내용으로부터 로그 레벨을 추론합니다
func DetectLogLevel(message string) string {
	lowerMsg := strings.ToLower(message)

	if strings.Contains(lowerMsg, "error") || strings.Contains(lowerMsg, "failed") {
		return "error"
	}

	if strings.Contains(lowerMsg, "warn") {
		return "warn"
	}

	return "info"
}
