package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

// Function to fetch automation results for a given jsmonId
func getAutomationResultsByJsmonId(jsmonId string, wkspId string) {
	// Define the API base URL and endpoint, appending the jsmonId as a query parameter
	endpoint := fmt.Sprintf("%s/getAllAutomationResults?inputType=jsmonid&input=%s&showonly=all&wkspId=%s", apiBaseURL, jsmonId, wkspId)

	// Create a new HTTP request with the GET method
	req, err := http.NewRequest("GET", endpoint, nil) // No need for request body in GET
	if err != nil {
		fmt.Printf("Failed to create request: %v\n", err)
		return
	}

	// Set necessary headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Jsmon-Key", strings.TrimSpace(getAPIKey())) // Trim any whitespace from the API key

	// Create an HTTP client and make the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Failed to send request: %v\n", err)
		return
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Failed to read response: %v\n", err)
		return
	}

	// Check if the response is successful (Status Code: 2xx)
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		// Print the response with all fields related to the jsmonId
		var result interface{}
		err = json.Unmarshal(body, &result)
		if err != nil {
			fmt.Println("Error parsing JSON:", err)
			return
		}

		// Assert that result is of type map[string]interface{}
		if resMap, ok := result.(map[string]interface{}); ok {
			// Access the "results" key
			if results, ok := resMap["results"].([]interface{}); ok {
				if len(results) > 0 {
					// Access the first element
					firstElement := results[0]

					// Print the first element as a pretty JSON string
					prettyJSON, err := json.MarshalIndent(firstElement, "", "  ")
					if err != nil {
						fmt.Println("Error formatting JSON:", err)
						return
					}

					fmt.Println(string(prettyJSON))
				} else {
					fmt.Println("Results array is empty.")
				}
			} else {
				fmt.Println("results is not of type []interface{}")
			}
		} else {
			fmt.Println("result is not of type map[string]interface{}")
		}
	} else {
		fmt.Printf("Error: Received status code %d\n", resp.StatusCode)
		fmt.Println("Response:", string(body)) // Print the response even if it's an error
	}
}
