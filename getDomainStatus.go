package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

func getAllDomainsStatus(input string, size int, wkspId string) {
	endpoint := fmt.Sprintf("%s/getAllAutomationResults?wkspId=%s", apiBaseURL, wkspId)
	url := fmt.Sprintf("%s?showonly=%s&inputType=domain&input=%s&size=%d&start=%d&sortOrder=%s&sortBy=%s", endpoint, "domainStatus", input, size, 0, "desc", "createdAt")

	// create request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Printf("Failed to create request")
		return
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Jsmon-Key", strings.TrimSpace(getAPIKey()))

	// send request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Failed to send request: %v\n", err)
		return
	}
	defer resp.Body.Close()

	// read response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Failed to read response body: %v\n", err)
		return
	}

	// parse response
	var response map[string]interface{}
	err = json.Unmarshal(body, &response)
	if err != nil {
		fmt.Printf("Failed to unmarshal JSON response: %v\n", err)
		return
	}

	// fmt.Println(response, "HELLP")
	// access urls map
	results, err := extractResults(response)
	if err != nil {
		fmt.Println(err)
		return
	}

	// fmt.Println("SOME OUTPUT FROM RESULTS", results)

	domainStatus := processExtractedDomains(results)
	// fmt.Println("domain Status", domainStatus)
	for _, domain := range domainStatus {
		fmt.Println(domain)
	}

}
func extractResults(response map[string]interface{}) ([]interface{}, error) {
	if results, ok := response["results"].([]interface{}); ok {
		return results, nil
	}
	return nil, fmt.Errorf("no data found in response")
}

func processExtractedDomains(results []interface{}) []string {
	output := make([]string, 0)

	for _, result := range results {
		if resultMap, ok := result.(map[string]interface{}); ok {
			// Access 'extractedDomainsStatus'

			if extractedDomainsStatus, ok := resultMap["extractedDomainsStatus"].([]interface{}); ok {
				domainStatusResult := processDomainStatus(extractedDomainsStatus)
				output = append(output, domainStatusResult...)
			}
		}
	}
	return output
}

func processDomainStatus(extractedDomainsStatus []interface{}) []string {
	output := make([]string, 0)

	for _, domainArray := range extractedDomainsStatus {
		if domains, ok := domainArray.(map[string]interface{}); ok {

			record := fmt.Sprintf("%s | %s | %s", domains["domainName"], domains["status"], domains["expiryDate"])
			fmt.Println(record)
		}
	}
	return output
}
