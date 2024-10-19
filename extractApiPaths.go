package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// Function to get API paths based on domains
func getApiPaths(domains []string) {
	// Prepare request data
	endpoint := fmt.Sprintf("%s/apiPathfromDomain", apiBaseURL)
	requestBody, err := json.Marshal(map[string]interface{}{
		"domains": domains,
	})
	if err != nil {
		fmt.Printf("Failed to marshal request body: %v\n", err)
		return
	}

	// Create request
	req, err := http.NewRequest("POST", endpoint, bytes.NewBuffer(requestBody))
	if err != nil {
		fmt.Printf("Failed to create request: %v\n", err)
		return
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Jsmon-Key", strings.TrimSpace(getAPIKey()))

	// Send request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Failed to send request: %v\n", err)
		return
	}
	defer resp.Body.Close()

	// Read response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Failed to read response body: %v\n", err)
		return
	}

	// Parse response
	var response struct {
		ApiPaths []string `json:"apiPaths"`
	}
	err = json.Unmarshal(body, &response)
	if err != nil {
		fmt.Printf("Failed to unmarshal JSON response: %v\n", err)
		return
	}

	// Access and print API paths
	for _, path := range response.ApiPaths {
		fmt.Println(path)
	}
}
