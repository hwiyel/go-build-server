package main

import (
	"encoding/json"
	"log"
	"net/http"
)

type BuildResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Version string `json:"version"`
}

type BuildRequest struct {
	ProjectName string `json:"project_name"`
	Environment string `json:"environment"`
}

type BuildCreateResponse struct {
	Status      string `json:"status"`
	Message     string `json:"message"`
	ProjectName string `json:"project_name"`
	BuildID     string `json:"build_id"`
	Timestamp   string `json:"timestamp"`
}

func buildHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	response := BuildResponse{
		Status:  "success",
		Message: "Build API is working!",
		Version: "1.0.0",
	}

	json.NewEncoder(w).Encode(response)
}

func buildCreateHandler(w http.ResponseWriter, r *http.Request) {
	// POST 메서드만 허용
	if r.Method != http.MethodPost {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Only POST method is allowed",
		})
		return
	}

	w.Header().Set("Content-Type", "application/json")

	var req BuildRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Invalid request body",
		})
		return
	}

	// 필수 필드 확인
	if req.ProjectName == "" || req.Environment == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "project_name and environment are required",
		})
		return
	}

	// 빌드 생성 응답
	response := BuildCreateResponse{
		Status:      "success",
		Message:     "Build created successfully!",
		ProjectName: req.ProjectName,
		BuildID:     "build-" + req.ProjectName + "-001",
		Timestamp:   "2026-02-03T14:36:00Z",
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

func main() {
	http.HandleFunc("/api/build", buildHandler)
	http.HandleFunc("/api/build/create", buildCreateHandler)

	log.Println("Server starting on http://localhost:8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
