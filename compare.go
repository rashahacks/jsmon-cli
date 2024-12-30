package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	// "time"
)

type DiffItem struct {
	Added   bool   `json:"added"`
	Removed bool   `json:"removed"`
	Value   string `json:"value"`
}

func compareEndpoint(id1, id2 string, wkspId string) {
	endpoint := fmt.Sprintf("%s/compare?wkspId=%s", apiBaseURL, wkspId)

	requestBody, err := json.Marshal(map[string]string{
		"id1": id1,
		"id2": id2,
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

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("Unexpected status code: %d\n", resp.StatusCode)
		body, _ := ioutil.ReadAll(resp.Body)
		fmt.Printf("Response: %s\n", string(body))
		return
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response:", err)
		return
	}

	var diffItems []DiffItem
	err = json.Unmarshal(body, &diffItems)
	if err != nil {
		fmt.Println("Error parsing JSON:", err)
		fmt.Printf("Response: %s\n", string(body))
		return
	}

	addedCount := 0
	removedCount := 0

	fmt.Println("Summary of changes:")
	for _, item := range diffItems {
		if item.Added {
			addedCount++
			if addedCount <= 20 { // Print the first 20 additions
				fmt.Printf("+ %s\n", item.Value)
			}
		} else if item.Removed {
			removedCount++
			if removedCount <= 20 { // Print the first 20 removals
				fmt.Printf("- %s\n", item.Value)
			}
		}
	}

	fmt.Printf("\nTotal additions: %d\n", addedCount)
	fmt.Printf("Total removals: %d\n", removedCount)
	if addedCount > 5 {
		fmt.Printf("(Only the first 20 additions are shown)\n")
	}
	if removedCount > 5 {
		fmt.Printf("(Only the first 20 removals are shown)\n")
	}
}
