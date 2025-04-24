package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// Response structure for JSON unmarshalling
type QueryBuilderResponse struct {
	JSUrls           []string      `json:"urls"`
	PaginatedResults []interface{} `json:"paginatedResults"`
}

func queryBuilder(wkspId, query string) {
	endpoint := fmt.Sprintf("%s/queryBuilder?wkspId=%s", apiBaseURL, wkspId)

	requestBody, err := json.Marshal(map[string]interface{}{
		"query": query,
	})
	if err != nil {
		fmt.Printf("Failed to marshal request body: %v\n", err)
		return
	}

	req, err := http.NewRequest("POST", endpoint, bytes.NewBuffer(requestBody))
	if err != nil {
		fmt.Printf("Failed to create request: %v\n", err)
		return
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Jsmon-Key", strings.TrimSpace(apiKey))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Failed to send request: %v\n", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized {
		fmt.Println("[ERR] Wrong API key")
		return
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Failed to read response body: %v\n", err)
		return
	}

	var result QueryBuilderResponse
	if err := json.Unmarshal(body, &result); err != nil {
		fmt.Printf("Failed to parse JSON: %v\n", err)
		return
	}

	// Case 1: Show JS URLs if present
	if len(result.JSUrls) > 0 {
		for _, url := range result.JSUrls {
			fmt.Println(url)
		}
		return
	}

	// Case 2: Show paginated results if present (each on new line, no commas)
	if len(result.PaginatedResults) > 0 {
		for _, item := range result.PaginatedResults {
			if val, ok := item.(string); ok {
				fmt.Println(val)
			}
		}
		return
	}

	// Case 3: Neither present
	fmt.Println("No JS URLs or paginated results found.")
}
