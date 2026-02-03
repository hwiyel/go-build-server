package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestBuildHandler(t *testing.T) {
	req, err := http.NewRequest("GET", "/api/build", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(buildHandler)
	handler.ServeHTTP(rr, req)

	// 상태 코드 확인
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Content-Type 확인
	expected := "application/json"
	if ct := rr.Header().Get("Content-Type"); ct != expected {
		t.Errorf("handler returned wrong content type: got %v want %v",
			ct, expected)
	}

	// 응답 본문 확인
	var response BuildResponse
	if err := json.NewDecoder(rr.Body).Decode(&response); err != nil {
		t.Errorf("failed to decode response: %v", err)
	}

	if response.Status != "success" {
		t.Errorf("handler returned wrong status: got %v want %v",
			response.Status, "success")
	}

	if response.Version != "1.0.0" {
		t.Errorf("handler returned wrong version: got %v want %v",
			response.Version, "1.0.0")
	}
}

func TestBuildHandlerResponseStructure(t *testing.T) {
	req, err := http.NewRequest("GET", "/api/build", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(buildHandler)
	handler.ServeHTTP(rr, req)

	var response BuildResponse
	if err := json.NewDecoder(rr.Body).Decode(&response); err != nil {
		t.Errorf("failed to decode response: %v", err)
	}

	// 모든 필드가 비어있지 않은지 확인
	if response.Status == "" {
		t.Error("Status field is empty")
	}
	if response.Message == "" {
		t.Error("Message field is empty")
	}
	if response.Version == "" {
		t.Error("Version field is empty")
	}
}

func TestBuildCreateHandler(t *testing.T) {
	payload := BuildRequest{
		ProjectName: "test-project",
		Environment: "development",
	}

	body, err := json.Marshal(payload)
	if err != nil {
		t.Fatal(err)
	}

	req, err := http.NewRequest("POST", "/api/build/create", bytes.NewReader(body))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(buildCreateHandler)
	handler.ServeHTTP(rr, req)

	// 상태 코드 확인 (201 Created)
	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusCreated)
	}

	// 응답 본문 확인
	var response BuildCreateResponse
	if err := json.NewDecoder(rr.Body).Decode(&response); err != nil {
		t.Errorf("failed to decode response: %v", err)
	}

	if response.Status != "success" {
		t.Errorf("handler returned wrong status: got %v want %v",
			response.Status, "success")
	}

	if response.ProjectName != "test-project" {
		t.Errorf("handler returned wrong project name: got %v want %v",
			response.ProjectName, "test-project")
	}

	if response.BuildID == "" {
		t.Error("BuildID field is empty")
	}
}

func TestBuildCreateHandlerMissingFields(t *testing.T) {
	// 필드가 없는 요청
	payload := BuildRequest{
		ProjectName: "test-project",
		// Environment 누락
	}

	body, err := json.Marshal(payload)
	if err != nil {
		t.Fatal(err)
	}

	req, err := http.NewRequest("POST", "/api/build/create", bytes.NewReader(body))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(buildCreateHandler)
	handler.ServeHTTP(rr, req)

	// 상태 코드 확인 (400 Bad Request)
	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusBadRequest)
	}
}

func TestBuildCreateHandlerInvalidMethod(t *testing.T) {
	// GET 요청으로 POST 핸들러 호출
	req, err := http.NewRequest("GET", "/api/build/create", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(buildCreateHandler)
	handler.ServeHTTP(rr, req)

	// 상태 코드 확인 (405 Method Not Allowed)
	if status := rr.Code; status != http.StatusMethodNotAllowed {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusMethodNotAllowed)
	}
}
