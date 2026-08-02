// Package openai is a minimal client for the OpenAI completions HTTP API.
package openai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const completionsURL = "https://api.openai.com/v1/completions"

// Client calls the OpenAI REST API with a bearer token.
type Client struct {
	apiKey     string
	httpClient *http.Client
}

// NewClient returns a Client authenticated with the given API key.
func NewClient(apiKey string) *Client {
	return &Client{
		apiKey:     apiKey,
		httpClient: &http.Client{Timeout: 30 * time.Second},
	}
}

// CompletionRequest is the request body for the completions endpoint.
type CompletionRequest struct {
	Model       string  `json:"model"`
	Prompt      string  `json:"prompt"`
	MaxTokens   int     `json:"max_tokens"`
	Temperature float64 `json:"temperature"`
}

// CompletionResponse is the subset of the completions response we use.
type CompletionResponse struct {
	ID      string `json:"id"`
	Model   string `json:"model"`
	Choices []struct {
		Text         string `json:"text"`
		FinishReason string `json:"finish_reason"`
	} `json:"choices"`
	Usage struct {
		PromptTokens     int `json:"prompt_tokens"`
		CompletionTokens int `json:"completion_tokens"`
		TotalTokens      int `json:"total_tokens"`
	} `json:"usage"`
}

// Complete sends prompt to the completions endpoint and returns the first
// generated text choice.
func (c *Client) Complete(ctx context.Context, prompt string) (string, error) {
	body, err := json.Marshal(CompletionRequest{
		Model:       "gpt-3.5-turbo-instruct",
		Prompt:      prompt,
		MaxTokens:   256,
		Temperature: 0.3,
	})
	if err != nil {
		return "", fmt.Errorf("marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, completionsURL, bytes.NewReader(body))
	if err != nil {
		return "", fmt.Errorf("build request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+c.apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("call completions API: %w", err)
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("read response: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("completions API returned %s: %s", resp.Status, data)
	}

	var completion CompletionResponse
	if err := json.Unmarshal(data, &completion); err != nil {
		return "", fmt.Errorf("decode response: %w", err)
	}
	if len(completion.Choices) == 0 {
		return "", fmt.Errorf("completions API returned no choices")
	}
	return completion.Choices[0].Text, nil
}
