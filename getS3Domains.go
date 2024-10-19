package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// getS3Domains retrieves API paths based on domains
func getS3Domains(domains []string) {
	endpoint := fmt.Sprintf("%s/getS3Domains", apiBaseURL)
	
	requestBody, err := json.Marshal(map[string]interface{}{
		"domains": domains,
	})
	if err != nil {
		fmt.Printf("Failed to marshal request body: %v\n", err)
		return
	}

	req, err := createRequest(endpoint, requestBody)
	if err != nil {
		fmt.Printf("Failed to create request: %v\n", err)
		return
	}

	resp, err := sendRequest(req)
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

	var response map[string]interface{}
	if err := json.Unmarshal(body, &response); err != nil {
		fmt.Printf("Failed to unmarshal JSON response: %v\n", err)
		return
	}

	printS3Domains(response)
}

func createRequest(endpoint string, requestBody []byte) (*http.Request, error) {
	req, err := http.NewRequest(http.MethodPost, endpoint, bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Jsmon-Key", strings.TrimSpace(getAPIKey()))
	return req, nil
}

func sendRequest(req *http.Request) (*http.Response, error) {
	client := &http.Client{}
	return client.Do(req)
}

func printS3Domains(response map[string]interface{}) {
	s3Domains, ok := response["s3Domains"].([]interface{})
	if !ok || len(s3Domains) == 0 {
		fmt.Println("No S3 domains found in the response")
		return
	}

	for _, domain := range s3Domains {
		if domainStr, ok := domain.(string); ok {
			fmt.Println(domainStr)
		}
	}
}
