package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	// "time"
)

type ScannerResult struct {
	Message string     `json:"message"`
	Data    []DataItem `json:"data"`
}

func getScannerResults(wkspId string) {
	endpoint := fmt.Sprintf("%s/getScannerResults?wkspId=%s", apiBaseURL, wkspId)

	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}

	req.Header.Set("X-Jsmon-Key", strings.TrimSpace(getAPIKey()))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response:", err)
		return
	}

	var result ScannerResult

	err = json.Unmarshal(body, &result)
	if err != nil {
		fmt.Println("Error parsing JSON:", err)
		return
	}

	fmt.Println("Message:", result.Message)
	prettyJSON, err := json.MarshalIndent(result.Data, "", "  ")
	if err != nil {
		fmt.Println("Error creating JSON:", err)
		return
	}

	// Print the pretty JSON output
	fmt.Printf("Data:\n%s\n", prettyJSON)
}
