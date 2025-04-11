package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

func queryBuilder(wkspId, query string) {
	// POST request to the queryBuilder endpoint
	endpoint := fmt.Sprintf("%s/queryBuilder?wkspId=%s", apiBaseURL, wkspId)

	// Create the request body
	requestBody, err := json.Marshal(map[string]interface{}{
		"query": query,
	})
	if err != nil {
		fmt.Printf("Failed to marshal request body: %v\n", err)
		os.Exit(1)
	}

	// Create the HTTP POST request
	req, err := http.NewRequest("POST", endpoint, bytes.NewBuffer(requestBody))
	if err != nil {
		fmt.Printf("Failed to create request: %v\n", err)
		os.Exit(1)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Jsmon-Key", strings.TrimSpace(getAPIKey()))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Failed to send request: %v\n", err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized {
		fmt.Println("[ERR] Wrong API key")
		return 
	}

	// Read and handle the response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Failed to read response body: %v\n", err)
		os.Exit(1)
	}

	// Parse the response
	var response map[string]interface{}
	err = json.Unmarshal(body, &response)
	if err != nil {
		fmt.Printf("Failed to unmarshal JSON response: %v\n", err)
		os.Exit(1)
	}

	// Format output
	var formattedResults []map[string]interface{}

	if paginatedResults, ok := response["paginatedResults"].([]interface{}); ok {
		for _, result := range paginatedResults {
			if resultMap, ok := result.(map[string]interface{}); ok {
				formattedResult := make(map[string]interface{})

				// Get the URL and domainName fields
				if url, ok := resultMap["url"].(string); ok {
					formattedResult["url"] = url
				}
				if domainName, ok := resultMap["domainName"].(string); ok {
					formattedResult["domainName"] = domainName
				}

				// Format detectedWords
				if detectedWords, ok := resultMap["detectedWords"].([]interface{}); ok {
					var formattedDetectedWords []map[string]interface{}
					for _, item := range detectedWords {
						if wordMap, ok := item.(map[string]interface{}); ok {
							formattedWord := make(map[string]interface{})
							if name, ok := wordMap["name"].(string); ok {
								formattedWord["name"] = name
							}
							if words, ok := wordMap["words"].([]interface{}); ok {
								var wordList []string
								for _, word := range words {
									if wordStr, ok := word.(string); ok {
										wordList = append(wordList, wordStr)
									}
								}
								formattedWord["words"] = wordList
							}
							formattedDetectedWords = append(formattedDetectedWords, formattedWord)
						}
					}
					formattedResult["detectedWords"] = formattedDetectedWords
				}

				formattedResults = append(formattedResults, formattedResult)
			}
		}

		// Print the formatted result as JSON
		output, err := json.MarshalIndent(formattedResults, "", "    ")
		if err != nil {
			fmt.Printf("Failed to marshal formatted results: %v\n", err)
			os.Exit(1)
		}
		fmt.Println(string(output))
	} else {
		fmt.Println("No paginated results found.")
	}
}