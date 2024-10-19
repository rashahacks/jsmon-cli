package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// IPAddresses represents the structure of the IP addresses in the response
type IPAddresses struct {
	IPv4 []string `json:"ipv4Addresses"`
	IPv6 []string `json:"ipv6Addresses"`
}

// Response represents the structure of the API response
type Response struct {
	IPAddresses IPAddresses `json:"ipAddresses"`
}

// getAllIps retrieves IP addresses based on given domains
func getAllIps(domains []string) {
	endpoint := fmt.Sprintf("%s/getIps", apiBaseURL)
	requestBody, err := json.Marshal(map[string][]string{"domains": domains})
	if err != nil {
		fmt.Printf("Failed to marshal request body: %v\n", err)
		return
	}

	req, err := http.NewRequest(http.MethodPost, endpoint, bytes.NewBuffer(requestBody))
	if err != nil {
		fmt.Printf("Failed to create request: %v\n", err)
		return
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Jsmon-Key", strings.TrimSpace(getAPIKey()))

	client := &http.Client{}
	resp, err := client.Do(req)
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

	var response Response
	if err := json.Unmarshal(body, &response); err != nil {
		fmt.Printf("Failed to unmarshal JSON response: %v\n", err)
		return
	}

	printIPAddresses(response.IPAddresses)
}

// printIPAddresses prints the IPv4 and IPv6 addresses
func printIPAddresses(ipAddresses IPAddresses) {
	if len(ipAddresses.IPv4) > 0 {
		fmt.Println("IPv4 Addresses:")
		for _, ip := range ipAddresses.IPv4 {
			fmt.Println(ip)
		}
	} else {
		fmt.Println("No IPv4 addresses found")
	}

	if len(ipAddresses.IPv6) > 0 {
		fmt.Println("IPv6 Addresses:")
		for _, ip := range ipAddresses.IPv6 {
			fmt.Println(ip)
		}
	} else {
		fmt.Println("No IPv6 addresses found")
	}
}
