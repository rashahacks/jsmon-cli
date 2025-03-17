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

type RescanDomainResponse struct {
	Message   string `json:"message"`
	TotalUrls int    `json:"totalUrls"`
}

func rescanDomain(domain string) {
	endpoint := fmt.Sprintf("%s/rescanDomain", apiBaseURL)

	requestBody, err := json.Marshal(map[string]string{
		"domain": domain,
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

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response:", err)
		return
	}

	var result RescanDomainResponse
	err = json.Unmarshal(body, &result)
	if err != nil {
		fmt.Println("Error parsing JSON:", err)
		return
	}

	fmt.Printf("Message: %s\n", result.Message)
	fmt.Printf("Total URLs submitted for scanning: %d\n", result.TotalUrls)
}
