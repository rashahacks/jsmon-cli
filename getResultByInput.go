package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"
)

// AutomationResult represents the structure of the API response
type AutomationResult struct {
	// Add fields based on the actual response structure
	// For example:
	// JsmonID string `json:"jsmonId"`
	// URL     string `json:"url"`
	// ... other fields
}

func getAutomationResultsByInput(ctx context.Context, inputType, value string) ([]AutomationResult, error) {
	endpoint := fmt.Sprintf("%s/getAllJsUrlsResults?inputType=%s&input=%s", apiBaseURL, inputType, value)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Jsmon-Key", strings.TrimSpace(getAPIKey()))

	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("unexpected status code: %d, response: %s", resp.StatusCode, string(body))
	}

	var results []AutomationResult
	if err := json.Unmarshal(body, &results); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return results, nil
}

// Usage example
func main() {
	ctx := context.Background()
	results, err := getAutomationResultsByInput(ctx, "someInputType", "someValue")
	if err != nil {
		log.Fatalf("Error getting automation results: %v", err)
	}

	for _, result := range results {
		log.Printf("Result: %+v", result)
	}
}
