package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

const (
	inputTypeParam = "jsmonid"
	showOnlyParam  = "all"
)

type APIResponse struct {
	Results []map[string]interface{} `json:"results"`
}

// Function to fetch automation results for a given jsmonId
func getAutomationResultsByJsmonId(jsmonId string) {
	endpoint := fmt.Sprintf("%s/getAllAutomationResults?inputType=%s&input=%s&showonly=%s", apiBaseURL, inputTypeParam, jsmonId, showOnlyParam)

	req, err := http.NewRequest(http.MethodGet, endpoint, nil)
	if err != nil {
		fmt.Printf("Failed to create request: %v\n", err)
		return
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Jsmon-Key", strings.TrimSpace(getAPIKey()))

	resp, err := httpClient.Do(req)
	if err != nil {
		fmt.Printf("Failed to send request: %v\n", err)
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Failed to read response: %v\n", err)
		return
	}

	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		printFirstResult(body)
	} else {
		fmt.Printf("Error: Received status code %d\n", resp.StatusCode)
		fmt.Println("Response:", string(body))
	}
}

func printFirstResult(body []byte) {
	var apiResp APIResponse
	if err := json.Unmarshal(body, &apiResp); err != nil {
		fmt.Println("Error parsing JSON:", err)
		return
	}

	if len(apiResp.Results) == 0 {
		fmt.Println("Results array is empty.")
		return
	}

	prettyJSON, err := json.MarshalIndent(apiResp.Results[0], "", "  ")
	if err != nil {
		fmt.Println("Error formatting JSON:", err)
		return
	}

	fmt.Println(string(prettyJSON))
}

