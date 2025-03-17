package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	// "time"
)

func rescanUrlEndpoint(scanId string) {
	endpoint := fmt.Sprintf("%s/rescanURL/%s", apiBaseURL, scanId)

	req, err := http.NewRequest("POST", endpoint, nil)
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

	var result interface{}
	err = json.Unmarshal(body, &result)
	if err != nil {
		fmt.Println("Error parsing JSON:", err)
		return
	}

	// Pretty print JSON
	prettyJSON, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		fmt.Println("Error formatting JSON:", err)
		return
	}

	fmt.Println(string(prettyJSON))
}
