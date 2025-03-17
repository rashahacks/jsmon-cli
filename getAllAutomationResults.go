package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

func getAllAutomationResults(input string, size int, wkspId string) error {
	// Constructing the API endpoint
	endpoint := fmt.Sprintf("%s/getAllAutomationResults?wkspId=%s", apiBaseURL, wkspId)
	url := fmt.Sprintf("%s&showonly=all&inputType=domain&input=%s&size=%d", endpoint, input, size)

	// Creating the HTTP request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Jsmon-Key", strings.TrimSpace(getAPIKey()))

	// Sending the HTTP request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return err
	}
	defer resp.Body.Close()

	// Reading the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response:", err)
		return err
	}

	// Log the entire response body for debugging
	fmt.Println(string(body))

	// Parse the JSON response
	var result map[string]interface{}
	err = json.Unmarshal(body, &result)
	if err != nil {
		fmt.Println("Error parsing JSON:", err)
		return err
	}

	if message, ok := result["message"].(string); ok {
		fmt.Println(message)
	} 

	return nil
}
