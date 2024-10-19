package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	urlPkg "net/url"
	"path/filepath"
)

// Common function to make API requests
func makeAPIRequest(endpoint, showonly, input string, size int, sortBy string) ([]interface{}, error) {
	url := fmt.Sprintf("%s/getAllAutomationResults?showonly=%s&inputType=domain&input=%s&size=%d&start=0&sortOrder=desc&sortBy=%s",
		apiBaseURL, showonly, input, size, sortBy)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Jsmon-Key", strings.TrimSpace(getAPIKey()))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to get data: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}

	var response map[string]interface{}
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON response: %v", err)
	}

	results, ok := response["results"].([]interface{})
	if !ok {
		return nil, fmt.Errorf("data not found or not in expected format")
	}

	return results, nil
}

// Common function to extract URLs from results
func extractURLs(results []interface{}, field string) []string {
	var urls []string
	for _, result := range results {
		if resultMap, ok := result.(map[string]interface{}); ok {
			if fieldUrls, ok := resultMap[field].([]interface{}); ok {
				for _, urlStr := range fieldUrls {
					if url, ok := urlStr.(string); ok {
						urls = append(urls, url)
					}
				}
			}
		}
	}
	return urls
}

func getAllFileExtensionUrls(input string, extensions []string, size int) {
	results, err := makeAPIRequest(apiBaseURL, "fileExtensionUrls", input, size, "fileExtensionUrls")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	urls := extractURLs(results, "fileExtensionUrls")

	if len(extensions) > 0 {
		extensionsMap := make(map[string]struct{})
		for _, ext := range extensions {
			extensionsMap["."+ext] = struct{}{}
		}

		var filteredUrls []string
		for _, urlStr := range urls {
			parsedURL, err := urlPkg.Parse(urlStr)
			if err != nil {
				continue
			}
			ext := filepath.Ext(parsedURL.Path)
			if _, ok := extensionsMap[ext]; ok {
				filteredUrls = append(filteredUrls, urlStr)
			}
		}
		urls = filteredUrls
	}

	sendOutputToStdout(urls)
}

func getAllSocialMediaUrls(input string, size int) {
	results, err := makeAPIRequest(apiBaseURL, "socialMediaUrls", input, size, "createdAt")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	urls := extractURLs(results, "socialMediaUrls")
	sendOutputToStdout(urls)
}

func getAllQueryParamsUrls(input string, size int) {
	results, err := makeAPIRequest(apiBaseURL, "queryParamsUrls", input, size, "createdAt")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	urls := extractURLs(results, "queryParamsUrls")
	sendOutputToStdout(urls)
}

func getAllLocalhostUrls(input string, size int) {
	results, err := makeAPIRequest(apiBaseURL, "localhostUrls", input, size, "localhostUrls")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	urls := extractURLs(results, "localhostUrls")
	sendOutputToStdout(urls)
}

func getAllFilteredPortUrls(input string, size int) {
	results, err := makeAPIRequest(apiBaseURL, "filteredPortUrls", input, size, "filteredPortUrls")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	urls := extractURLs(results, "filteredPortUrls")
	sendOutputToStdout(urls)
}

func getAllS3DomainsInvalid(input string, size int) {
	results, err := makeAPIRequest(apiBaseURL, "s3DomainsInvalid", input, size, "s3DomainsInvalid")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	urls := extractURLs(results, "s3DomainsInvalid")
	sendOutputToStdout(urls)
}

func sendOutputToStdout(output []string) {
	for _, urlStr := range output {
		fmt.Println(urlStr)
	}
}
