package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	// "time"
)

func callViewProfile() error {
	endpoint := fmt.Sprintf("%s/viewProfile", apiBaseURL)

	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		return fmt.Errorf("error creating request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Jsmon-Key", strings.TrimSpace(getAPIKey()))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error making request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusForbidden {
		return fmt.Errorf("invalid API Key")
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("error reading response body: %v", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return fmt.Errorf("error unmarshalling response: %v", err)
	}

	if errMsg, ok := result["error"].(string); ok && errMsg != "" {
		return fmt.Errorf("invalid API Key")
	}

	if status, ok := result["status"].(string); ok && status != "success" {
		if message, ok := result["message"].(string); ok {
			return fmt.Errorf("%s", message)
		}
		return fmt.Errorf("invalid API Key")
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
		return nil
	}

	return fmt.Errorf("invalid API Key")
}
