package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// getDomainUrls fetches and prints domain URLs for the given domains
func getDomainUrls(domains []string) {
	endpoint := fmt.Sprintf("%s/getDomainsUrls", apiBaseURL)
	
	requestBody, err := json.Marshal(map[string]interface{}{
		"domains": domains,
	})
	if err != nil {
		fmt.Printf("Failed to marshal request body: %v\n", err)
		return
	}

	req, err := http.NewRequest(http.MethodPost, endpoint, bytes.NewBuffer(requestBody))
	if err != nil {
		fmt.Printf("Failed to create request: %v\n", err)
		return
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Jsmon-Key", strings.TrimSpace(getAPIKey()))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Failed to send request: %v\n", err)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Failed to read response body: %v\n", err)
		return
	}

	var response struct {
		Data struct {
			ExtractedDomains []string `json:"extractedDomains"`
			ExtractedUrls    []string `json:"extractedUrls"`
		} `json:"data"`
	}

	if err := json.Unmarshal(body, &response); err != nil {
		fmt.Printf("Failed to unmarshal JSON response: %v\n", err)
		return
	}

	// Print extracted domains
	for _, domain := range response.Data.ExtractedDomains {
		fmt.Println(domain)
	}

	// Print extracted URLs
	for _, url := range response.Data.ExtractedUrls {
		fmt.Println(url)
	}
}
