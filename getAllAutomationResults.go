package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

func getAllAutomationResults(input string, size int, wkspId string) error {
	endpoint := fmt.Sprintf("%s/getAllAutomationResults?wkspId=%s", apiBaseURL, wkspId)

	url := fmt.Sprintf("%s?showonly=all&inputType=domain&input=%s&size=%d", endpoint, input, size)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return err
	}

	req.Header.Set("X-Jsmon-Key", strings.TrimSpace(getAPIKey()))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response:", err)
		return err
	}

	var result map[string]interface{}
	err = json.Unmarshal(body, &result)
	if err != nil {
		fmt.Println("Error parsing JSON:", err)
		return err
	}

	// Extract and print the message field
	if message, ok := result["message"].(string); ok {
		fmt.Println(message)
	} else {
		fmt.Println("Message not found in response.")
	}

	return nil
}
