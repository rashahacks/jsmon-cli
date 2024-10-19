package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

const (
	timeout         = 30 * time.Second
	endpointFormat  = "%s/getAllAutomationResults?inputType=fileid&input=%s&showonly=all"
	contentTypeJSON = "application/json"
)

type APIResponse struct {
	Results []map[string]interface{} `json:"results"`
}

// getAutomationResultsByFileId fetches automation results for a given fileId
func getAutomationResultsByFileId(fileId string) error {
	if fileId == "" {
		return fmt.Errorf("fileId cannot be empty")
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	endpoint := fmt.Sprintf(endpointFormat, apiBaseURL, fileId)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", contentTypeJSON)
	req.Header.Set("X-Jsmon-Key", strings.TrimSpace(getAPIKey()))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("received non-success status code %d: %s", resp.StatusCode, string(body))
	}

	var apiResp APIResponse
	if err := json.Unmarshal(body, &apiResp); err != nil {
		return fmt.Errorf("error parsing JSON: %w", err)
	}

	if len(apiResp.Results) == 0 {
		fmt.Println("Results array is empty.")
		return nil
	}

	prettyJSON, err := json.MarshalIndent(apiResp.Results[0], "", "  ")
	if err != nil {
		return fmt.Errorf("error formatting JSON: %w", err)
	}

	fmt.Println(string(prettyJSON))
	return nil
}
