package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

var fieldMapping = map[string]string{
	"urls":                   "extractedUrls",
	"domains":                "extractedDomains",
	"ipv4":                   "ipv4Addresses",
	"ipv6":                   "ipv6Addresses",
	"emails":                 "emails",
	"cloud-buckets":          "s3Domains",
	"apis":                   "apiPaths",
	"gql-queries":            "gqlQuery",
	"gql-mutations":          "gqlMutation",
	"node-modules-confusion": "invalidNodeModules",
	"node-modules":           "validNodeModules",
	"gql-fragments":          "gqlFragment",
	"vulnerabilities":        "vulnerabilities",
	"guids":                  "guids",
	"domains-status":         "extractedDomainsStatus",
	"urls-parameters":        "queryParamsUrls",
	"bucket-takeovers":       "invalidS3Domains",
	"urls-socialmedia":       "socialMediaUrls",
	"urls-localhost":         "localhostUrls",
	"urls-ports":             "filteredPortUrls",
	"exposures":              "exposures",
	"jsUrls":                 "jsUrls",
}

type QueryBuilderResponse struct {
	PaginatedResults         []map[string]interface{} `json:"paginatedResults"`
	URLs                     []string                 `json:"urls"`
	IsReverseSearchAvailable bool                     `json:"isReverseSearchAvailable"`
}

func queryBuilder(wkspId, query string) {

	if strings.HasPrefix(query, "field=") {
		fieldType := strings.TrimPrefix(query, "field=")
		if mappedField, exists := fieldMapping[fieldType]; exists {
			query = fmt.Sprintf("field:%s", mappedField)
		} else {
			fmt.Printf("Field type '%s' not found in mapping\n", fieldType)
		}
	}

	endpoint := fmt.Sprintf("%s/queryBuilder?wkspId=%s", apiBaseURL, wkspId)

	requestBody, err := json.Marshal(map[string]string{
		"query": query,
	})
	if err != nil {
		fmt.Printf("Failed to marshal request body: %v\n", err)
		return
	}
	req, err := http.NewRequest("POST", endpoint, bytes.NewBuffer(requestBody))
	if err != nil {
		fmt.Printf("Failed to create request: %v\n", err)
		return
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Jsmon-Key", strings.TrimSpace(apiKey))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Failed to send request: %v\n", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized {
		fmt.Println("[ERR] Wrong API key")
		return
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Failed to read response body: %v\n", err)
		return
	}

	var result QueryBuilderResponse
	if err := json.Unmarshal(body, &result); err != nil {
		fmt.Printf("Failed to parse JSON: %v\n", err)
		return
	}
	if len(result.URLs) > 0 {
		for _, url := range result.URLs {
			fmt.Println(url)
		}
		return
	}
	if len(result.PaginatedResults) > 0 {
		for _, item := range result.PaginatedResults {
			jsonData, err := json.MarshalIndent(item, "", "    ")
			if err != nil {
				fmt.Printf("Error formatting result: %v\n", err)
				continue
			}
			fmt.Println(string(jsonData))
		}
		return
	}

	fmt.Println("No results found")
}
