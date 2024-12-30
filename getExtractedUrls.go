package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	urlPkg "net/url"
	"path/filepath"
	"strings"
)

func getAllFileExtensionUrls(input string, extensions []string, size int, wkspId string) {
	endpoint := fmt.Sprintf("%s/getAllAutomationResults?wkspId=%s", apiBaseURL, wkspId)
	url := fmt.Sprintf("%s?showonly=%s&inputType=domain&input=%s&size=%d&start=%d&sortOrder=%s&sortBy=%s", endpoint, "fileExtensionUrls", input, size, 0, "desc", "fileExtensionUrls")

	// create request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Printf("Failed to create request: %v\n", err)
		return
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Jsmon-Key", strings.TrimSpace(getAPIKey()))

	// send request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Failed to get data")
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

	// maping the extensions to a map, while converting pdf to .pdf
	extensionsMap := make(map[string]struct{}, 0)
	for _, ext := range extensions {
		extensionsMap["."+ext] = struct{}{}
	}

	// fmt.Println(response, "HELLP")
	// access urls map
	results, err := extractResultsFromResponse(response)
	if err != nil {
		return
	}

	urls := captureResultsFromResponse(results, "fileExtensionUrls")
	new_urls := make([]string, 0)

	if len(extensions) > 0 {
		for _, urlStr := range urls {
			parsedURL, err := urlPkg.Parse(urlStr)
			if err != nil {
				continue
			}
			ext := filepath.Ext(parsedURL.Path)

			if _, ok := extensionsMap[ext]; ok {
				new_urls = append(new_urls, urlStr)
			}
		}
		urls = new_urls
	}
	sendOutputToStdout(urls)
}

func getAllSocialMediaUrls(input string, size int, wkspId string) {
	endpoint := fmt.Sprintf("%s/getAllAutomationResults?wkspId=%s", apiBaseURL, wkspId)
	url := fmt.Sprintf("%s?showonly=%s&inputType=domain&input=%s&size=%d&start=%d&sortOrder=%s&sortBy=%s", endpoint, "socialMediaUrls", input, size, 0, "desc", "createdAt")

	// fmt.Println("URL :", url)
	// create request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Printf("Failed to create request: %v\n", err)
		return
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Jsmon-Key", strings.TrimSpace(getAPIKey()))

	// send request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Failed to get data")
		return
	}
	defer resp.Body.Close()

	// read response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Failed to read data")
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

	results, err := extractResultsFromResponse(response)
	if err != nil {
		return
	}

	urls := captureResultsFromResponse(results, "socialMediaUrls")

	sendOutputToStdout(urls)
}

func getAllQueryParamsUrls(input string, size int, wkspId string) {
	endpoint := fmt.Sprintf("%s/getAllAutomationResults?wkspId=%s", apiBaseURL, wkspId)
	url := fmt.Sprintf("%s?showonly=%s&inputType=domain&input=%s&size=%d&start=%d&sortOrder=%s&sortBy=%s", endpoint, "queryParamsUrls", input, size, 0, "desc", "createdAt")

	// create request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Printf("Failed to create request: %v\n", err)
		return
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Jsmon-Key", strings.TrimSpace(getAPIKey()))

	// send request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Failed to get data")
		return
	}
	defer resp.Body.Close()

	// read response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Failed to read data")
		return
	}

	// parse response
	var response map[string]interface{}
	err = json.Unmarshal(body, &response)
	if err != nil {
		fmt.Printf("Failed to parse data")
		return
	}

	// parse the response from the API call
	results, err := extractResultsFromResponse(response)
	if err != nil {
		return
	}

	urls := captureResultsFromResponse(results, "queryParamsUrls")
	sendOutputToStdout(urls)
}

func getAllLocalhostUrls(input string, size int, wkspId string) {
	endpoint := fmt.Sprintf("%s/getAllAutomationResults?wkspId=%s", apiBaseURL, wkspId)
	url := fmt.Sprintf("%s?showonly=%s&inputType=domain&input=%s&size=%d&start=%d&sortOrder=%s&sortBy=%s", endpoint, "localhostUrls", input, size, 0, "desc", "localhostUrls")

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
		fmt.Printf("Failed to get data")
		return
	}
	defer resp.Body.Close()

	// read response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Failed to read data")
		return
	}

	// parse response
	var response map[string]interface{}
	err = json.Unmarshal(body, &response)
	if err != nil {
		fmt.Printf("Failed to parse data")
		return
	}

	// parse the response from the API call
	// parse the response from the API call
	results, err := extractResultsFromResponse(response)
	if err != nil {
		return
	}

	urls := captureResultsFromResponse(results, "localhostUrls")
	sendOutputToStdout(urls)
}

func getAllFilteredPortUrls(input string, size int, wkspId string) {
	endpoint := fmt.Sprintf("%s/getAllAutomationResults?wkspId=%s", apiBaseURL, wkspId)
	url := fmt.Sprintf("%s?showonly=%s&inputType=domain&input=%s&size=%d&start=%d&sortOrder=%s&sortBy=%s", endpoint, "filteredPortUrls", input, size, 0, "desc", "filteredPortUrls")

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
		fmt.Printf("Failed to get data")
		return
	}
	defer resp.Body.Close()

	// read response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Failed to read data")
		return
	}

	// parse response
	var response map[string]interface{}
	err = json.Unmarshal(body, &response)
	if err != nil {
		fmt.Printf("Failed to parse data")
		return
	}

	// parse the response from the API call
	results, err := extractResultsFromResponse(response)
	if err != nil {
		return
	}

	urls := captureResultsFromResponse(results, "filteredPortUrls")
	sendOutputToStdout(urls)
}

func getAllS3DomainsInvalid(input string, size int, wkspId string) {
	endpoint := fmt.Sprintf("%s/getAllAutomationResults?wkspId=%s", apiBaseURL, wkspId)
	url := fmt.Sprintf("%s?showonly=%s&inputType=domain&input=%s&size=%d&start=%d&sortOrder=%s&sortBy=%s", endpoint, "s3DomainsInvalid", input, size, 0, "desc", "s3DomainsInvalid")

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
		fmt.Printf("Failed to get data")
		return
	}
	defer resp.Body.Close()

	// read response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Failed to read data")
		return
	}

	// parse response
	var response map[string]interface{}
	err = json.Unmarshal(body, &response)
	if err != nil {
		fmt.Printf("Failed to parse data")
		return
	}

	// parse the response from the API call
	results, err := extractResultsFromResponse(response)
	if err != nil {
		return
	}

	urls := captureResultsFromResponse(results, "s3DomainsInvalid")
	sendOutputToStdout(urls)
}

func extractResultsFromResponse(response map[string]interface{}) ([]interface{}, error) {
	if results, ok := response["results"].([]interface{}); ok {
		return results, nil
	}
	return nil, fmt.Errorf("data not found or not in expected format")
}

func captureResultsFromResponse(results []interface{}, field string) []string {
	output := make([]string, 0)
	for _, result := range results {
		if resultMap, ok := result.(map[string]interface{}); ok {
			// Access 'extractedDomainsStatus'
			if fieldUrls, ok := resultMap[field].([]interface{}); ok {
				for _, urlStr := range fieldUrls {
					if url, ok := urlStr.(string); ok {
						// fmt.Println(url)
						output = append(output, url)
					} else {
						fmt.Println("Found Invalid URL format")
					}
				}
			}
		}
	}
	return output
}

func sendOutputToStdout(output []string) {
	for _, urlStr := range output {
		fmt.Println(urlStr)
	}
}
