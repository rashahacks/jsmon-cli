package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// Define a struct for the domain status
type DomainStatus struct {
	DomainName string `json:"domainName"`
	Status     string `json:"status"`
	ExpiryDate string `json:"expiryDate"`
}

func getAllDomainsStatus(input string, size int) {
	endpoint := fmt.Sprintf("%s/getAllAutomationResults", apiBaseURL)
	url := fmt.Sprintf("%s?showonly=domainStatus&inputType=domain&input=%s&size=%d&start=0&sortOrder=desc&sortBy=createdAt", endpoint, input, size)

	// Create request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Printf("Failed to create request: %v\n", err)
		return
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Jsmon-Key", strings.TrimSpace(getAPIKey()))

	// Send request
	client := &http.Client{Timeout: 10 * time.Second}
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
	var response map[string]interface{}
	if err := json.Unmarshal(body, &response); err != nil {
		fmt.Printf("Failed to unmarshal JSON response: %v\n", err)
		return
	}

	results, err := extractResults(response)
	if err != nil {
		fmt.Println(err)
		return
	}

	domainStatuses := processExtractedDomains(results)
	for _, status := range domainStatuses {
		fmt.Println(status)
	}
}

func extractResults(response map[string]interface{}) ([]interface{}, error) {
	results, ok := response["results"].([]interface{})
	if !ok {
		return nil, fmt.Errorf("no data found in response")
	}
	return results, nil
}

func processExtractedDomains(results []interface{}) []string {
	var output []string

	for _, result := range results {
		resultMap, ok := result.(map[string]interface{})
		if !ok {
			continue
		}

		extractedDomainsStatus, ok := resultMap["extractedDomainsStatus"].([]interface{})
		if !ok {
			continue
		}

		output = append(output, processDomainStatus(extractedDomainsStatus)...)
	}

	return output
}

func processDomainStatus(extractedDomainsStatus []interface{}) []string {
	var output []string

	for _, domainArray := range extractedDomainsStatus {
		domains, ok := domainArray.(map[string]interface{})
		if !ok {
			continue
		}

		status := DomainStatus{
			DomainName: fmt.Sprintf("%v", domains["domainName"]),
			Status:     fmt.Sprintf("%v", domains["status"]),
			ExpiryDate: fmt.Sprintf("%v", domains["expiryDate"]),
		}

		record := fmt.Sprintf("%s | %s | %s", status.DomainName, status.Status, status.ExpiryDate)
		output = append(output, record)
	}

	return output
}
