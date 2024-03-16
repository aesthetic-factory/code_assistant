package http_client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

// Request struct represents the input data for TextGenerateRemote request
type TextGenRequest struct {
	URL         string  `json:"url"`
	Model       string  `json:"model"`
	Temperature float64 `json:"temperature"`
	Prompt      string  `json:"prompt"`
	Stream      bool    `json:"stream"`
}

// Response struct represents the output data from TextGenerateRemote response
type TextGenResponse struct {
	Result string `json:"response"`
	Token  int32  `json:"eval_count"`
}

// GenerateRemote sends a POST request to the remote server with the provided data
// It takes a TextGenRequest struct as input and returns a TextGenResponse struct
func TextGenerateRemote(req TextGenRequest) (TextGenResponse, error) {
	// Convert TextGenRequest struct to JSON
	reqBody, err := json.Marshal(req)
	if err != nil {
		return TextGenResponse{}, fmt.Errorf("failed to marshal request JSON: %v", err)
	}

	// Create a context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel() // Ensure cancel is called to release resources

	// Create HTTP client with timeout
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	// Create a new request with the context
	httpRequest, err := http.NewRequestWithContext(ctx, "POST", req.URL, bytes.NewBuffer(reqBody))
	if err != nil {
		return TextGenResponse{}, fmt.Errorf("failed to create HTTP request: %v", err)
	}

	// Set content type header
	httpRequest.Header.Set("Content-Type", "application/json")

	// Send the HTTP request
	resp, err := client.Do(httpRequest)
	if err != nil {
		// Check if the error is due to a timeout
		if ctx.Err() == context.DeadlineExceeded {
			return TextGenResponse{}, fmt.Errorf("request timed out: %v", err)
		}
		return TextGenResponse{}, fmt.Errorf("failed to send HTTP request: %v", err)
	}
	defer resp.Body.Close()

	// Read response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return TextGenResponse{}, fmt.Errorf("failed to read response body: %v", err)
	}

	// Unmarshal response JSON into Response struct
	var response TextGenResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return TextGenResponse{}, fmt.Errorf("failed to unmarshal response JSON: %v", err)
	}

	return response, nil
}

// Request struct represents the input data for EmbeddingGenerateRemote request
type EmbeddingRequest struct {
	URL    string `json:"url"`
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
}

// Response struct represents the output data from EmbeddingGenerateRemote response
type EmbeddingResponse struct {
	Result []float32 `json:"embedding"`
}

// EmbeddingGenerateRemote sends a POST request to the remote server with the provided data
// It takes a EmbeddingRequest struct as input and returns a EmbeddingResponse struct
func EmbeddingGenerateRemote(req EmbeddingRequest) (EmbeddingResponse, error) {
	// Convert EmbeddingRequest struct to JSON
	reqBody, err := json.Marshal(req)
	if err != nil {
		return EmbeddingResponse{}, fmt.Errorf("failed to marshal request JSON: %v", err)
	}

	// Create a context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel() // Ensure cancel is called to release resources

	// Create HTTP client with timeout
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	// Create a new request with the context
	httpRequest, err := http.NewRequestWithContext(ctx, "POST", req.URL, bytes.NewBuffer(reqBody))
	if err != nil {
		return EmbeddingResponse{}, fmt.Errorf("failed to create HTTP request: %v", err)
	}

	// Set content type header
	httpRequest.Header.Set("Content-Type", "application/json")

	// Send the HTTP request
	resp, err := client.Do(httpRequest)
	if err != nil {
		// Check if the error is due to a timeout
		if ctx.Err() == context.DeadlineExceeded {
			return EmbeddingResponse{}, fmt.Errorf("request timed out: %v", err)
		}
		return EmbeddingResponse{}, fmt.Errorf("failed to send HTTP request: %v", err)
	}
	defer resp.Body.Close()

	// Read response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return EmbeddingResponse{}, fmt.Errorf("failed to read response body: %v", err)
	}

	// Unmarshal response JSON into Response struct
	var response EmbeddingResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return EmbeddingResponse{}, fmt.Errorf("failed to unmarshal response JSON: %v", err)
	}

	return response, nil
}
