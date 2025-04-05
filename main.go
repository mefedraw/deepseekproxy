package main

import (
	"bytes"
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

// Структуры для DeepSeek API
type DeepSeekMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type DeepSeekRequest struct {
	Model    string            `json:"model"`
	Messages []DeepSeekMessage `json:"messages"`
}

type DeepSeekResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}

func main() {
	http.HandleFunc("/prompt", handlePrompt)
	http.HandleFunc("/health", handleHealthCheck)

	port := os.Getenv("PORT")
	if port == "" {
		port = "3232"
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

	log.Printf("Received prompt: %s\n", req.Prompt)

	deepSeekResponse, err := callDeepSeekAPI(req.Prompt)
	if err != nil {
		log.Printf("DeepSeek API error: %v\n", err)
		http.Error(w, "Error processing prompt with DeepSeek", http.StatusInternalServerError)
		return
	}

	resp := PromptResponse{Response: deepSeekResponse}
	jsonResp, err := json.Marshal(resp)
	if err != nil {
		http.Error(w, "Error creating response", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResp)

	log.Printf("Response: %s\n", deepSeekResponse)
}

func callDeepSeekAPI(prompt string) (string, error) {
	apiKey := os.Getenv("DEEPSEEK_API_KEY")
	if apiKey == "" {
		return "", fmt.Errorf("DEEPSEEK_API_KEY environment variable not set")
	}

	// Правильный формат запроса для DeepSeek API
	deepSeekReq := DeepSeekRequest{
		Model: "deepseek-chat",
		Messages: []DeepSeekMessage{
			{
				Role:    "user",
				Content: prompt,
			},
		},
	}

	requestBody, err := json.Marshal(deepSeekReq)
	if err != nil {
		return "", fmt.Errorf("error marshaling DeepSeek request: %v", err)
	}

	req, err := http.NewRequest(
		"POST",
		"https://api.deepseek.com/v1/chat/completions",
		bytes.NewBuffer(requestBody),
	)
	if err != nil {
		return "", fmt.Errorf("error creating request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("error sending request to DeepSeek: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error reading response body: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("DeepSeek API returned status %d: %s", resp.StatusCode, string(body))
	}

	var deepSeekResp DeepSeekResponse
	err = json.Unmarshal(body, &deepSeekResp)
	if err != nil {
		return "", fmt.Errorf("error parsing DeepSeek response: %v", err)
	}

	if len(deepSeekResp.Choices) == 0 {
		return "", fmt.Errorf("no choices in DeepSeek response")
	}

	return deepSeekResp.Choices[0].Message.Content, nil
}

func handleHealthCheck(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}
