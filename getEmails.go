package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

// Function to get emails
func getEmails(domains []string) {
	// Prepare request data
	endpoint := fmt.Sprintf("%s/getEmails", apiBaseURL)
	requestBody, err := json.Marshal(map[string]interface{}{
		"domains": domains,
	})
	if err != nil {
		fmt.Println("Error creating request body:", err)
		return
	}

	req, err := http.NewRequest("POST", endpoint, bytes.NewBuffer(requestBody))
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Jsmon-Key", strings.TrimSpace(getAPIKey()))

	// Send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return
	}
	defer resp.Body.Close()

	// Read the response
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response:", err)
		return
	}

	// Parse and print JSON response
	var result map[string]interface{}
	err = json.Unmarshal(body, &result)
	if err != nil {
		fmt.Println("Error parsing JSON:", err)
		return
	}

	// Pretty print JSON
	emails, ok := result["emails"].([]interface{})
	if !ok {
		fmt.Println("Error: 'email' field not found or not in expected format")
		return
	}
	for _, path := range emails {
		if pathStr, ok := path.(string); ok {
			fmt.Println(pathStr)
		} else {
			fmt.Println("Error: Invalid type in 'emails'")
		}
	}
}
