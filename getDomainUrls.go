package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

// Function to get domain URLs
func getDomainUrls(domains []string, wkspId string) {
	endpoint := fmt.Sprintf("%s/getDomainsUrls?wkspId=%s", apiBaseURL, wkspId)
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
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Failed to read response body: %v\n", err)
		return
	}

	// Parse response
	var response map[string]interface{}
	err = json.Unmarshal(body, &response)
	if err != nil {
		fmt.Printf("Failed to unmarshal JSON response: %v\n", err)
		return
	}

	// Extract data
	if data, ok := response["data"].(map[string]interface{}); ok {

		// Print extracted domains
		if extractedDomains, ok := data["extractedDomains"].([]interface{}); ok {
			for _, domain := range extractedDomains {
				if domainStr, ok := domain.(string); ok {
					fmt.Println(domainStr)
				}
			}
		}

		// Print extracted URLs
		if extractedUrls, ok := data["extractedUrls"].([]interface{}); ok {
			for _, url := range extractedUrls {
				if urlStr, ok := url.(string); ok {
					fmt.Println(urlStr)
				}
			}
		}
	} else {
		fmt.Println("Error: 'data' field not found or not in expected format")
	}
}
