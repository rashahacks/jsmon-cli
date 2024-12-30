package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	// "time"
)

type addWordlistRequest struct {
	Domains []string `json:"domains"`
}

func createWordList(domains []string, wkspId string) {
	endpoint := fmt.Sprintf("%s/createWordList?wkspId=%s", apiBaseURL, wkspId)

	requestBody := addWordlistRequest{
		Domains: domains,
	}
	body, err := json.Marshal(requestBody)
	if err != nil {
		fmt.Printf("failed to marshal request body: %v\n", err)
		return
	}

	// Create HTTP request
	client := &http.Client{}
	req, err := http.NewRequest("POST", endpoint, bytes.NewBuffer(body))
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Jsmon-Key", strings.TrimSpace(getAPIKey()))

	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("failed to send request: %v\n", err)
		return
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("failed to read response body: %v\n", err)
		return
	}

	fmt.Printf("Word list:\n%s\n", string(responseBody))
}
