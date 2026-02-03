package main

import (
	"api-server/pkg/handlers"
	"api-server/pkg/services"
	"log"
	"net/http"
)

func main() {
	// 의존성 주입
	logService := services.NewInMemoryLogService()

	// 핸들러 생성
	jobHandler := handlers.NewBuildJobHandler(logService)
	logsHandler := handlers.NewLogsHandler(logService)

	// BuildJob API 라우팅
	http.HandleFunc("/api/buildjob", jobHandler.Create)
	http.HandleFunc("/api/buildjob/", logsHandler.Get)

	log.Println("Server starting on http://localhost:8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}

