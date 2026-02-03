package models

// BuildJobRequest는 POST /api/buildjob 요청 구조입니다
type BuildJobRequest struct {
	JobName           string `json:"job_name"`
	DockerfileContent string `json:"dockerfile_content"`
	ImageName         string `json:"image_name,omitempty"`
	PushRegistry      bool   `json:"push_registry,omitempty"`
}

// BuildJobResponse는 POST /api/buildjob 응답 구조입니다
type BuildJobResponse struct {
	Status    string `json:"status"`
	Message   string `json:"message"`
	JobName   string `json:"job_name"`
	JobID     string `json:"job_id"`
	Namespace string `json:"namespace"`
	CreatedAt string `json:"created_at"`
}

// LogEntry는 로그 항목 구조입니다
type LogEntry struct {
	Timestamp string `json:"timestamp"`
	Container string `json:"container"`
	Message   string `json:"message"`
	Level     string `json:"level"`
}

// LogsResponse는 GET /api/buildjob/{job_name}/logs 응답 구조입니다
type LogsResponse struct {
	JobName    string     `json:"job_name"`
	Status     string     `json:"status"`
	Logs       []LogEntry `json:"logs"`
	TotalLines int        `json:"total_lines"`
}

// ErrorResponse는 에러 응답 구조입니다
type ErrorResponse struct {
	Error string `json:"error"`
}

// BuildRequest는 POST /api/build/create 요청 구조입니다 (레거시)
type BuildRequest struct {
	ProjectName string `json:"project_name"`
	Environment string `json:"environment"`
}

// BuildResponse는 GET /api/build 응답 구조입니다 (레거시)
type BuildResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Version string `json:"version"`
}

// BuildCreateResponse는 POST /api/build/create 응답 구조입니다 (레거시)
type BuildCreateResponse struct {
	Status      string `json:"status"`
	Message     string `json:"message"`
	ProjectName string `json:"project_name"`
	BuildID     string `json:"build_id"`
	Timestamp   string `json:"timestamp"`
}
