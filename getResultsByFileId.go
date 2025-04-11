package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// Function to fetch automation results for a given jsmonId
func getAutomationResultsByFileId(fileId string, wkspId string) {
	// Define the API base URL and endpoint, appending the jsmonId as a query parameter
	endpoint := fmt.Sprintf("%s/getAllAutomationResults?inputType=fileid&input=%s&showonly=all&wkspId=%s", apiBaseURL, fileId, wkspId)

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
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Failed to read response: %v\n", err)
		return
	}

	// Check if the response is successful (Status Code: 2xx)
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		// Print the response with all fields related to the jsmonId
		var result map[string]interface{}
		err = json.Unmarshal(body, &result)
		if err != nil {
			fmt.Println("Error parsing JSON:", err)
			return
		}

		// Assert that result is of type map[string]interface{}
		if val, ok := result["results"]; ok {
			switch v := val.(type) {
			case []interface {}:
			// Access the "results" key
				if len(v) > 0 {
					prettyJSON, err := json.MarshalIndent(v[0], "", "  ")
					if err != nil {
						fmt.Println("Error formatting JSON:", err)
						return
					}

					fmt.Println(string(prettyJSON))
				} else {
					fmt.Println("Results array is empty.")
				}
			default:
				fmt.Printf("[ERR] Unexpected type for 'results': %T\n", v)
			} 
		} else {
			fmt.Println("No 'results' key found in response.")
		}
	} else {
		fmt.Printf("Error: Received status code %d\n", resp.StatusCode)
		fmt.Println("Response:", string(body))
	}
}
