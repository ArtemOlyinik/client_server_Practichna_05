package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/ArtemOlyinik/client_server_Practichna_05/internal/domain"
	"github.com/ArtemOlyinik/client_server_Practichna_05/internal/dto"
)

var client = &http.Client{Timeout: 5 * time.Second}

func main() {
	baseURL := strings.TrimRight(env("BASE_URL", "http://localhost:8080"), "/")
	email := fmt.Sprintf("engineer%d@company.internal", time.Now().Unix())

	fmt.Println("=== Перевірка Incident Tracker API ===")
	fmt.Printf("Тестовий інженер: %s\n", email)

	do(baseURL, http.MethodPost, "/api/v1/auth/register", map[string]string{
		"email":    email,
		"password": "SuperSecretPassword123",
	}, "")

	_, loginBody := do(baseURL, http.MethodPost, "/api/v1/auth/login", map[string]string{
		"email":    email,
		"password": "SuperSecretPassword123",
	}, "")

	var login dto.LoginResponse
	if err := json.Unmarshal(loginBody, &login); err != nil {
		log.Fatalf("не вдалося розібрати токен: %v", err)
	}

	_, createdBody := do(baseURL, http.MethodPost, "/api/v1/incidents", map[string]string{
		"title":       "Overheating on Gateway-01",
		"description": "Core temperature reached 85C",
		"priority":    "critical",
	}, login.Token)

	var incident domain.Incident
	if err := json.Unmarshal(createdBody, &incident); err != nil {
		log.Fatalf("не вдалося розібрати створений інцидент: %v", err)
	}

	do(baseURL, http.MethodGet, "/api/v1/incidents", nil, "")
	do(baseURL, http.MethodGet, fmt.Sprintf("/api/v1/incidents/%d", incident.ID), nil, "")
	do(baseURL, http.MethodPatch, fmt.Sprintf("/api/v1/incidents/%d/status", incident.ID), map[string]string{
		"status": "in_progress",
	}, login.Token)
	do(baseURL, http.MethodDelete, fmt.Sprintf("/api/v1/incidents/%d", incident.ID), nil, login.Token)
	do(baseURL, http.MethodGet, fmt.Sprintf("/api/v1/incidents/%d", incident.ID), nil, "")

	do(baseURL, http.MethodDelete, fmt.Sprintf("/api/v1/incidents/%d", incident.ID), nil, "")
	do(baseURL, http.MethodPost, "/api/v1/incidents", map[string]string{
		"title":       "Overheating on Gateway-01",
		"description": "Core temperature reached 85C",
		"priority":    "urgent",
	}, login.Token)

	fmt.Println("\nПеревірку завершено.")
}

func do(baseURL, method, path string, payload any, token string) (*http.Response, []byte) {
	var body io.Reader
	if payload != nil {
		encoded, err := json.Marshal(payload)
		if err != nil {
			log.Fatalf("не вдалося закодувати JSON: %v", err)
		}
		body = bytes.NewReader(encoded)
	}

	request, err := http.NewRequest(method, baseURL+path, body)
	if err != nil {
		log.Fatalf("не вдалося створити запит: %v", err)
	}
	if payload != nil {
		request.Header.Set("Content-Type", "application/json")
	}
	if token != "" {
		request.Header.Set("Authorization", "Bearer "+token)
	}

	response, err := client.Do(request)
	if err != nil {
		log.Fatalf("помилка виконання запиту %s %s: %v", method, path, err)
	}
	defer response.Body.Close()

	data, err := io.ReadAll(response.Body)
	if err != nil {
		log.Fatalf("не вдалося прочитати відповідь: %v", err)
	}

	fmt.Printf("\n[HTTP] %s %s -> %d %s\n", method, path, response.StatusCode, http.StatusText(response.StatusCode))
	if len(data) > 0 {
		var pretty bytes.Buffer
		if err := json.Indent(&pretty, data, "", "  "); err == nil {
			fmt.Println(pretty.String())
		} else {
			fmt.Println(string(data))
		}
	}

	return response, data
}

func env(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
