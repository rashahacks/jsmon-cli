package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

// Function to get API paths based on domains
func getS3Domains(domains []string) {
	// Prepare request data
	endpoint := fmt.Sprintf("%s/getS3Domains", apiBaseURL)
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

	// Extract and print the data in the desired format
	message, ok := response["message"].(string)
	if ok {
		fmt.Println(message)
	}

	s3Domains, ok := response["s3Domains"].([]interface{})
	if ok && len(s3Domains) > 0 {
		fmt.Println("S3 Domains:")
		for _, domain := range s3Domains {
			if domainStr, ok := domain.(string); ok {
				fmt.Println(domainStr)
			}
		}
	} else {
		fmt.Println("No S3 Domains found")
	}
}
