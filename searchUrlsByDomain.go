package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	// "time"
)

type SearchUrlsByDomainResponse struct {
	Message   string     `json:"message"`
	TotalUrls int        `json:"totalUrls"`
	URLs      []URLEntry `json:"urls"`
}

func searchUrlsByDomain(domain string, wkspId string) {
	endpoint := fmt.Sprintf("%s/searchUrlbyDomain?domain=%s&wkspId=%s", apiBaseURL, url.QueryEscape(domain), wkspId)

	req, err := http.NewRequest("POST", endpoint, nil)
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

	if resp.StatusCode == http.StatusUnauthorized {
		fmt.Println("[ERR] Wrong API key")
		return 
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response:", err)
		return
	}

	var result SearchUrlsByDomainResponse
	err = json.Unmarshal(body, &result)
	if err != nil {
		fmt.Println("Error parsing JSON:", err)
		return
	}

	// fmt.Printf("Message: %s\n", result.Message)
	// fmt.Printf("Total URLs: %d\n", result.TotalUrls)
	fmt.Println("URLs:")
	for _, entry := range result.URLs {
		fmt.Printf("- %s\n", entry.URL)
	}

}
