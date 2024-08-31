package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

func getGqlOps(domains []string) {
	endpoint := fmt.Sprintf("%s/getGqlOps", apiBaseURL)
	requestBody, err := json.Marshal(map[string]interface{}{
		"domains": domains,
	})
	if err != nil {
		fmt.Printf("Failed to marshal request body: %v\n", err)
		return
	}

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
		fmt.Printf("Failed to read response: %v\n", err)
		return
	}

	// Print the raw response body for debugging
	//fmt.Println("Raw response body:")
	//fmt.Println(string(body))

	// Parse response
	var response map[string]interface{}
	err = json.Unmarshal(body, &response)
	if err != nil {
		fmt.Printf("Failed to unmarshal: %v\n", err)
		return
	}

	// Check if 'gqlOps' exists in the response
	if gqlOpsData, ok := response["gqlOps"].(map[string]interface{}); ok {
		// Extract gqlQueries 
		if gqlQuery, ok := gqlOpsData["gqlQuery"].([]interface{}); ok {
			fmt.Println("gqlQueries:")
			for _, gql := range gqlQuery {
				if gqlStr, ok := gql.(string); ok {
					fmt.Println(gqlStr)
				} else {
					fmt.Println("Error: Invalid type in 'gqlQueries'")
				}
			}
		} else {
			fmt.Println("gqlQueries not found")
		}

		// Extract gqlMutations 
		if gqlMutation, ok := gqlOpsData["gqlMutation"].([]interface{}); ok {
			fmt.Println("gqlMutations:")
			for _, gql := range gqlMutation {
				if gqlStr, ok := gql.(string); ok {
					fmt.Println(gqlStr)
				} else {
					fmt.Println("Error: Invalid type in 'gqlMutations'")
				}
			}
		} else {
			fmt.Println("gqlMutations not found")
		}
	} else {
		fmt.Println("Error: 'gqlOps' field not found")
	}
}

