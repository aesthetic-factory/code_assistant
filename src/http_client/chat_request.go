package http_client

import (
	"bytes"
	"code_assistant/src/config"
	"code_assistant/src/llm_prompt"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// Chat struct represents the Messages in ChatRequest
type Chat struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// Request struct represents the input data for ChatGenerateRemote request
type ChatRequest struct {
	URL         string  `json:"url"`
	Model       string  `json:"model"`
	Temperature float64 `json:"temperature"`
	Top_p       float64 `json:"top_p"`
	Messages    []Chat  `json:"messages"`
	Stream      bool    `json:"stream"`
	System      string  `json:"system"`
}

// Response struct represents the output data from ChatGenerateRemote response
type ChatResponse struct {
	Result Chat  `json:"message"`
	Token  int32 `json:"eval_count"`
}

// Create a ChatGenerateRemote request with default value
func NewChatRequest() ChatRequest {
	req := ChatRequest{
		URL:         config.AppConfig.Ollama.BaseUrl + "/api/chat",
		Model:       config.AppConfig.Ollama.ChatModel,
		Temperature: 0.15,
		Top_p:       0.3,
		Stream:      false,
		Messages:    []Chat{},
		System:      llm_prompt.SystemPrompt(),
	}
	return req
}

// ChatGenerateRemote sends a POST request to the remote server with the provided data
// It takes a ChatRequest struct as input and returns a ChatResponse struct
func ChatGenerateRemote(req ChatRequest) (ChatResponse, error) {

	// Convert ChatRequest struct to JSON
	reqBody, err := json.MarshalIndent(req, "", "  ")
	if err != nil {
		return ChatResponse{}, fmt.Errorf("failed to marshal request JSON: %v", err)
	}

	// fmt.Println(string(reqBody)) //DEBUG

	// Create a context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), HTTP_TIMEOUT_SEC*time.Second)
	defer cancel() // Ensure cancel is called to release resources

	// Create HTTP client with timeout
	client := &http.Client{
		Timeout: HTTP_TIMEOUT_SEC * time.Second,
	}

	// Create a new request with the context
	httpRequest, err := http.NewRequestWithContext(ctx, "POST", req.URL, bytes.NewBuffer(reqBody))
	if err != nil {
		return ChatResponse{}, fmt.Errorf("failed to create HTTP request: %v", err)
	}

	// Set content type header
	httpRequest.Header.Set("Content-Type", "application/json")

	// Send the HTTP request
	resp, err := client.Do(httpRequest)
	if err != nil {
		// Check if the error is due to a timeout
		if ctx.Err() == context.DeadlineExceeded {
			return ChatResponse{}, fmt.Errorf("request timed out: %v", err)
		}
		return ChatResponse{}, fmt.Errorf("failed to send HTTP request: %v", err)
	}
	defer resp.Body.Close()

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return ChatResponse{}, fmt.Errorf("failed to read response body: %v", err)
	}

	// Unmarshal response JSON into Response struct
	var response ChatResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return ChatResponse{}, fmt.Errorf("failed to unmarshal response JSON: %v", err)
	}

	return response, nil
}
