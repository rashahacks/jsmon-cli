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

func uploadUrlEndpoint(url string, customHeaders []string, wkspId string) {
	endpoint := fmt.Sprintf("%s/uploadUrl?wkspId=%s", apiBaseURL, wkspId)

	headerObjects := make([]map[string]string, 0)
	for _, header := range customHeaders {
		parts := strings.SplitN(header, ":", 2)
		if len(parts) == 2 {
			headerObjects = append(headerObjects, map[string]string{
				strings.TrimSpace(parts[0]): strings.TrimSpace(parts[1]),
			})
		}
	}

	requestBody, err := json.Marshal(map[string]interface{}{
		"url":     url,
		"headers": headerObjects,
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

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response:", err)
		return
	}

	var result struct {
		Message   string `json:"message"`
		JsmonID   string `json:"jsmonId"`
		Hash      string `json:"hash"`
		CreatedAt int64  `json:"createdAt"`
		URL       string `json:"url"`
	}

	err = json.Unmarshal(body, &result)
	if err != nil {
		fmt.Println("Error parsing JSON:", err)
		return
	}

	fmt.Printf("{\n")
	fmt.Printf("    \"message\": \"%s\",\n", result.Message)
	fmt.Printf("    \"jsmonId\": \"%s\",\n", result.JsmonID)
	fmt.Printf("    \"hash\": \"%s\",\n", result.Hash)
	fmt.Printf("    \"createdAt\": %d,\n", result.CreatedAt)
	fmt.Printf("    \"url\": \"%s\"\n", result.URL)
	fmt.Printf("}\n")

	if result.JsmonID != "" {
		getAutomationResultsByJsmonId(result.JsmonID, wkspId)
	}
}
