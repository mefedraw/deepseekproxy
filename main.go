package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

type PromptRequest struct {
	Prompt string `json:"prompt"`
}

type PromptResponse struct {
	Response string `json:"response"`
}

func main() {
	http.HandleFunc("/prompt", handlePrompt)
	http.HandleFunc("/health", handleHealthCheck)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s...\n", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func handlePrompt(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	var req PromptRequest
	err = json.Unmarshal(body, &req)
	if err != nil {
		http.Error(w, "Error parsing JSON", http.StatusBadRequest)
		return
	}

	// Логируем полученный промпт в консоль
	log.Printf("Received prompt: %s\n", req.Prompt)

	// Здесь можно добавить обработку промпта (например, отправить в AI API)
	// Для примера просто возвращаем модифицированный промпт
	response := fmt.Sprintf("Processed: %s", req.Prompt)

	resp := PromptResponse{Response: response}
	jsonResp, err := json.Marshal(resp)
	if err != nil {
		http.Error(w, "Error creating response", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResp)

	// Выводим ответ в консоль сервера
	log.Printf("Response: %s\n", response)
}

func handleHealthCheck(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}
