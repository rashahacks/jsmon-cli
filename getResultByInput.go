package main
import (
	"fmt"
	"strings"
	"net/http"
	"io/ioutil"
)

func getAutomationResultsByInput(inputType, value string){
	endpoint := fmt.Sprintf("%s/getAllJsUrlsResults?inputType=%s&input=%s", apiBaseURL, inputType, value)


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
		fmt.Println(string(body))
	} else {
		fmt.Printf("Error: Received status code %d\n", resp.StatusCode)
		fmt.Println("Response:", string(body)) // Print the response even if it's an error
	}
    
}