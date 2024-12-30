package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	// "time"
)

func callViewProfile() {
	endpoint := fmt.Sprintf("%s/viewProfile", apiBaseURL)

	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
		os.Exit(1)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Jsmon-Key", strings.TrimSpace(getAPIKey()))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error making request:", err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		os.Exit(1)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		fmt.Println("Error unmarshalling response:", err)
		os.Exit(1)
	}

	if data, ok := result["data"].(map[string]interface{}); ok {
		var accountType string
		if orgFound, ok := data["orgFound"].(bool); ok && orgFound {
			accountType = "org"
		} else if personalProfile, ok := data["personalProfile"].(bool); ok && personalProfile {
			accountType = "user"
		} else {
			accountType = "unknown"
		}

		filteredResult := map[string]interface{}{
			"limits": data["apiCallLimits"],
			"type":   accountType,
		}
		filteredData, _ := json.MarshalIndent(filteredResult, "", "  ")
		fmt.Println(string(filteredData))
	} else {
		fmt.Println("Error: Invalid response format")
	}
}
