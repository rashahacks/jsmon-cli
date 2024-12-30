package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

// Function to get IP addresses based on domains
func getAllIps(domains []string, wkspId string) {
	// Prepare request data
	endpoint := fmt.Sprintf("%s/getIps?wkspId=%s", apiBaseURL, wkspId)
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

	// Access ipAddresses map
	if ipData, ok := response["ipAddresses"].(map[string]interface{}); ok {
		// Extract and print IPv4 addresses
		if ipv4, ok := ipData["ipv4Addresses"].([]interface{}); ok {
			for _, ip := range ipv4 {
				if ipStr, ok := ip.(string); ok {
					fmt.Println(ipStr)
				} else {
					fmt.Println("Error: Invalid type in 'ipv4Addresses'")
				}
			}
		} else {
			fmt.Println("No 'ipv4Addresses' found or not in expected format")
		}

		// Extract and print IPv6 addresses
		if ipv6, ok := ipData["ipv6Addresses"].([]interface{}); ok {
			for _, ip := range ipv6 {
				if ipStr, ok := ip.(string); ok {
					fmt.Println(ipStr)
				} else {
					fmt.Println("Error: Invalid type in 'ipv6Addresses'")
				}
			}
		} else {
			fmt.Println("No 'ipv6Addresses' found or not in expected format")
		}
	} else {
		fmt.Println("Error: 'ipAddresses' field not found or not in expected format")
	}
}
