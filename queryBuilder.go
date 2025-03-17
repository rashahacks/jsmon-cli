package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

func queryBuilder(wkspId, query string) {
	// POST request to the queryBuilder endpoint
	endpoint := fmt.Sprintf("%s/queryBuilder?wkspId=%s", apiBaseURL, wkspId)

	// Create the request body
	requestBody, err := json.Marshal(map[string]interface{}{
		"query": query, // The query string is passed here
	})
	if err != nil {
		fmt.Printf("Failed to marshal request body: %v\n", err)
		os.Exit(1)
	}

	// Create the HTTP POST request
	req, err := http.NewRequest("POST", endpoint, bytes.NewBuffer(requestBody))
	if err != nil {
		fmt.Printf("Failed to create request: %v\n", err)
		os.Exit(1)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Jsmon-Key", strings.TrimSpace(getAPIKey()) )// Assuming a function to get the API key

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Failed to send request: %v\n", err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	// Read and handle the response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Failed to read response body: %v\n", err)
		os.Exit(1)
	}

	// Parse the response
	var response map[string]interface{}
	err = json.Unmarshal(body, &response)
	if err != nil {
		fmt.Printf("Failed to unmarshal JSON response: %v\n", err)
		os.Exit(1)
	}

	if urls, ok := response["urls"].([]interface{}); ok {
		for _, url := range urls {
			fmt.Println(url)
		}
	}

	if paginatedResults, ok := response["paginatedResults"].([]interface{}); ok {
		for _, result := range paginatedResults {
			fmt.Println(result)
		}
	} 
}