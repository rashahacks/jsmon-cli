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

func createWorkspace(workspace string) {
	// POST request to create workspace endpoint
	endpoint := fmt.Sprintf("%s/createWorkspace", apiBaseURL)

	requestBody, err := json.Marshal(map[string]interface{}{
		"name": workspace,
	})
	if err != nil {
		fmt.Printf("Failed to marshal request body: %v\n", err)
		os.Exit(1)
	}

	req, err := http.NewRequest("POST", endpoint, bytes.NewBuffer(requestBody))
	if err != nil {
		fmt.Printf("Failed to create request: %v\n", err)
		os.Exit(1)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Jsmon-Key", strings.TrimSpace(getAPIKey()))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Failed to send request: %v\n", err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Failed to read response body: %v\n", err)
		os.Exit(1)
	}

	var response map[string]interface{}
	if err := json.Unmarshal(body, &response); err != nil {
		fmt.Printf("Error parsing JSON response: %v\n", err)
		os.Exit(1)
	}

	// Print the entire response for debugging
	// fmt.Printf("Response: %+v\n", response)

	// Extract and print workspace details
	if workspaceID, ok := response["workspaceId"].(string); ok {
		fmt.Printf("Workspace created successfully! ID: %s\n", workspaceID)
	} else if message, ok := response["message"].(string); ok {
		fmt.Println(message)
	} else {
		fmt.Println("Unexpected response format:", response)
	}
}
