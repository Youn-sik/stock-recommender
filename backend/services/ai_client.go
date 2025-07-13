package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"stock-recommender/backend/config"
	"stock-recommender/backend/models"
	"time"
)

type AIClient struct {
	baseURL string
	client  *http.Client
}

func NewAIClient(cfg *config.Config) *AIClient {
	return &AIClient{
		baseURL: cfg.API.AIServiceURL,
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (c *AIClient) GetDecision(request models.AIDecisionRequest) (*models.AIDecisionResponse, error) {
	url := fmt.Sprintf("%s/api/v1/decision", c.baseURL)
	
	// Convert to JSON
	jsonData, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}
	
	// Create HTTP request
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	
	req.Header.Set("Content-Type", "application/json")
	
	// Make request
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("AI service returned status %d", resp.StatusCode)
	}
	
	// Parse response
	var aiResponse models.AIDecisionResponse
	if err := json.NewDecoder(resp.Body).Decode(&aiResponse); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}
	
	return &aiResponse, nil
}

func (c *AIClient) HealthCheck() error {
	url := fmt.Sprintf("%s/health", c.baseURL)
	
	resp, err := c.client.Get(url)
	if err != nil {
		return fmt.Errorf("failed to check health: %w", err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("AI service health check failed with status %d", resp.StatusCode)
	}
	
	return nil
}

// GetModelStatus returns current AI model status
func (c *AIClient) GetModelStatus() (map[string]interface{}, error) {
	url := fmt.Sprintf("%s/api/v1/models/status", c.baseURL)
	
	resp, err := c.client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to get model status: %w", err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get model status, status: %d", resp.StatusCode)
	}
	
	var status map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&status); err != nil {
		return nil, fmt.Errorf("failed to decode status response: %w", err)
	}
	
	return status, nil
}