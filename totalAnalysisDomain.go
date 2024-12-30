package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	// "time"
)

type AnalysisData struct {
	TotalDocuments         int `json:"totalDocuments"`
	TotalUrls              int `json:"totalUrls"`
	TotalDomains           int `json:"totalDomains"`
	TotalS3Domains         int `json:"totalS3Domains"`
	TotalEmails            int `json:"totalEmails"`
	TotalApiPaths          int `json:"totalApiPaths"`
	TotalJwtTokens         int `json:"totalJwtTokens"`
	TotalNodeModules       int `json:"totalNodeModules"`
	TotalGuids             int `json:"totalGuids"`
	TotalQueryParamsUrls   int `json:"totalQueryParamsUrls"`
	TotalS3DomainsInvalid  int `json:"totalS3DomainsInvalid"`
	TotalSocialMediaUrls   int `json:"totalSocialMediaUrls"`
	TotalLocalhostUrls     int `json:"totalLocalhostUrls"`
	TotalFilteredPortUrls  int `json:"totalFilteredPortUrls"`
	TotalFileExtensionUrls int `json:"totalFileExtensionUrls"`
	TotalVulnerabilities   int `json:"totalVulnerabilities"`
	TotalIpAddresses       int `json:"totalIpAddresses"`
	TotalGql               int `json:"totalGql"`
}

func totalAnalysisData(wkspId string) {
	endpoint := fmt.Sprintf("%s/totalCountAnalysisData?wkspId=%s", apiBaseURL, wkspId)

	req, err := http.NewRequest("GET", endpoint, nil)
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

	var results AnalysisData
	err = json.Unmarshal(body, &results)
	if err != nil {
		fmt.Println("Error parsing JSON:", err)
		return
	}

	fmt.Printf("Total Documents: %d\n", results.TotalDocuments)
	fmt.Printf("Total URLs: %d\n", results.TotalUrls)
	fmt.Printf("Total Domains: %d\n", results.TotalDomains)
	fmt.Printf("Total S3 Domains: %d\n", results.TotalS3Domains)
	fmt.Printf("Total Emails: %d\n", results.TotalEmails)
	fmt.Printf("Total API Paths: %d\n", results.TotalApiPaths)
	fmt.Printf("Total JWT Tokens: %d\n", results.TotalJwtTokens)
	fmt.Printf("Total Node Modules: %d\n", results.TotalNodeModules)
	fmt.Printf("Total GUIDs: %d\n", results.TotalGuids)
	fmt.Printf("Total Query Params URLs: %d\n", results.TotalQueryParamsUrls)
	fmt.Printf("Total S3 Domains Invalid: %d\n", results.TotalS3DomainsInvalid)
	fmt.Printf("Total Social Media URLs: %d\n", results.TotalSocialMediaUrls)
	fmt.Printf("Total Localhost URLs: %d\n", results.TotalLocalhostUrls)
	fmt.Printf("Total Filtered Port URLs: %d\n", results.TotalFilteredPortUrls)
	fmt.Printf("Total File Extension URLs: %d\n", results.TotalFileExtensionUrls)
	fmt.Printf("Total Vulnerabilities: %d\n", results.TotalVulnerabilities)
	fmt.Printf("Total IP Addresses: %d\n", results.TotalIpAddresses)
	fmt.Printf("Total GraphQL Queries: %d\n", results.TotalGql)
}
