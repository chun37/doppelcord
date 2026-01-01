package llm

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

// Config はLLMクライアントの設定
type Config struct {
	APIURL string
	APIKey string
	Model  string
}

// Client はOpenAI互換APIクライアント
type Client struct {
	config     Config
	httpClient *http.Client
}

// NewClient は新しいLLMクライアントを生成
func NewClient(config Config) *Client {
	return &Client{
		config: config,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// Chat はチャットリクエストを送信し、レスポンスを返す
func (c *Client) Chat(ctx context.Context, prompt string) (string, error) {
	log.Printf("[LLM Request] Model: %s", c.config.Model)
	log.Printf("[LLM Request] User: %s", prompt)

	req := ChatRequest{
		Model: c.config.Model,
		Messages: []ChatMessage{
			{Role: "user", Content: prompt},
		},
	}

	body, err := json.Marshal(req)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", c.config.APIURL, bytes.NewReader(body))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	if c.config.APIKey != "" {
		httpReq.Header.Set("Authorization", "Bearer "+c.config.APIKey)
	}

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return "", fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("API returned status %d: %s", resp.StatusCode, string(respBody))
	}

	var chatResp ChatResponse
	if err := json.Unmarshal(respBody, &chatResp); err != nil {
		return "", fmt.Errorf("failed to unmarshal response: %w", err)
	}

	if chatResp.Error != nil {
		return "", errors.New(chatResp.Error.Message)
	}

	if len(chatResp.Choices) == 0 {
		return "", errors.New("no response from API")
	}

	return chatResp.Choices[0].Message.Content, nil
}

// ChatWithSystem はsystemプロンプト付きでチャットリクエストを送信
func (c *Client) ChatWithSystem(ctx context.Context, systemPrompt, userPrompt string) (string, error) {
	log.Printf("[LLM Request] Model: %s", c.config.Model)
	log.Printf("[LLM Request] System: %s", systemPrompt)
	log.Printf("[LLM Request] User: %s", userPrompt)

	req := ChatRequest{
		Model: c.config.Model,
		Messages: []ChatMessage{
			{Role: "system", Content: systemPrompt},
			{Role: "user", Content: userPrompt},
		},
	}

	body, err := json.Marshal(req)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", c.config.APIURL, bytes.NewReader(body))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	if c.config.APIKey != "" {
		httpReq.Header.Set("Authorization", "Bearer "+c.config.APIKey)
	}

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return "", fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("API returned status %d: %s", resp.StatusCode, string(respBody))
	}

	var chatResp ChatResponse
	if err := json.Unmarshal(respBody, &chatResp); err != nil {
		return "", fmt.Errorf("failed to unmarshal response: %w", err)
	}

	if chatResp.Error != nil {
		return "", errors.New(chatResp.Error.Message)
	}

	if len(chatResp.Choices) == 0 {
		return "", errors.New("no response from API")
	}

	return chatResp.Choices[0].Message.Content, nil
}
