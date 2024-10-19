package main

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	"golang.org/x/time/rate"
)

const (
	defaultTimeout = 30 * time.Second
	maxRetries     = 3
	retryDelay     = 1 * time.Second
)

var (
	httpClient *http.Client
	limiter    *rate.Limiter
)

func init() {
	httpClient = &http.Client{
		Timeout: defaultTimeout,
	}
	limiter = rate.NewLimiter(rate.Every(time.Second), 10) // 10 requests per second
}

type DiffItem struct {
	Added   bool   `json:"added"`
	Removed bool   `json:"removed"`
	Value   string `json:"value"`
}

type DomainResponse struct {
	Domains []string `json:"domains"`
}

type AutomateScanDomainRequest struct {
	Domain string   `json:"domain"`
	Words  []string `json:"words"`
}

type addCustomWordsRequest struct {
	Words []string `json:"words"`
}

type addWordlistRequest struct {
	Domains []string `json:"domains"`
}

type AnalysisResult struct {
	Message     string `json:"message"`
	TotalChunks int    `json:"totalChunks"`
}

type ModuleScanResult struct {
	Message string `json:"message"`
	Data    []struct {
		ModuleName string `json:"ModuleName"`
		URL        string `json:"URL"`
	} `json:"data"`
}

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

type ScanResponse struct {
	Message          string           `json:"message"`
	AnalysisResult   AnalysisResult   `json:"analysis_result"`
	ModuleScanResult ModuleScanResult `json:"modulescan_result"`
}

type AutomateScanDomainResponse struct {
	Message       string       `json:"message"`
	FileId        string       `json:"fileId"`
	TrimmedDomain string       `json:"trimmedDomain"`
	ScanResponse  ScanResponse `json:"scanResponse"`
}

type URLEntry struct {
	URL string `json:"url"`
}

type SearchUrlsByDomainResponse struct {
	Message   string     `json:"message"`
	TotalUrls int        `json:"totalUrls"`
	URLs      []URLEntry `json:"urls"`
}

type RescanDomainResponse struct {
	Message   string `json:"message"`
	TotalUrls int    `json:"totalUrls"`
}

func makeRequest(ctx context.Context, method, url string, body io.Reader, headers map[string]string) (*http.Response, error) {
	if err := limiter.Wait(ctx); err != nil {
		return nil, fmt.Errorf("rate limit exceeded: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	for key, value := range headers {
		req.Header.Set(key, value)
	}

	req.Header.Set("X-Jsmon-Key", strings.TrimSpace(getAPIKey()))

	var resp *http.Response
	var attempt int
	for attempt = 0; attempt < maxRetries; attempt++ {
		resp, err = httpClient.Do(req)
		if err == nil {
			break
		}
		time.Sleep(retryDelay)
	}

	if err != nil {
		return nil, fmt.Errorf("error sending request after %d attempts: %w", maxRetries, err)
	}

	return resp, nil
}

func rescanDomain(domain string) {
	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	endpoint := fmt.Sprintf("%s/rescanDomain", apiBaseURL)

	requestBody, err := json.Marshal(map[string]string{
		"domain": domain,
	})
	if err != nil {
		log.Printf("Error creating request body: %v", err)
		return
	}

	headers := map[string]string{
		"Content-Type": "application/json",
	}

	resp, err := makeRequest(ctx, http.MethodPost, endpoint, bytes.NewBuffer(requestBody), headers)
	if err != nil {
		log.Printf("Error making request: %v", err)
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error reading response: %v", err)
		return
	}

	var result RescanDomainResponse
	if err := json.Unmarshal(body, &result); err != nil {
		log.Printf("Error parsing JSON: %v", err)
		return
	}

	fmt.Printf("Message: %s\n", result.Message)
	fmt.Printf("Total URLs submitted for scanning: %d\n", result.TotalUrls)
}

func totalAnalysisData() {
	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	endpoint := fmt.Sprintf("%s/totalCountAnalysisData", apiBaseURL)

	headers := map[string]string{
		"Content-Type": "application/json",
	}

	resp, err := makeRequest(ctx, http.MethodGet, endpoint, nil, headers)
	if err != nil {
		log.Printf("Error making request: %v", err)
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error reading response: %v", err)
		return
	}

	var results AnalysisData
	if err := json.Unmarshal(body, &results); err != nil {
		log.Printf("Error parsing JSON: %v", err)
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

func searchUrlsByDomain(domain string) {
	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	endpoint := fmt.Sprintf("%s/searchUrlbyDomain?domain=%s", apiBaseURL, url.QueryEscape(domain))

	headers := map[string]string{
		"Content-Type": "application/json",
	}

	resp, err := makeRequest(ctx, http.MethodPost, endpoint, nil, headers)
	if err != nil {
		log.Printf("Error making request: %v", err)
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error reading response: %v", err)
		return
	}

	var result SearchUrlsByDomainResponse
	if err := json.Unmarshal(body, &result); err != nil {
		log.Printf("Error parsing JSON: %v", err)
		return
	}

	fmt.Printf("Message: %s\n", result.Message)
	fmt.Printf("Total URLs: %d\n", result.TotalUrls)
	fmt.Println("URLs:")
	for _, entry := range result.URLs {
		fmt.Printf("- %s\n", entry.URL)
	}
}

func uploadUrlEndpoint(url string, customHeaders []string) {
	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	endpoint := fmt.Sprintf("%s/uploadUrl", apiBaseURL)

	headerObjects := make([]map[string]string, 0)
	for _, header := range customHeaders {
		parts := strings.SplitN(header, ":", 2)
		if len(parts) == 2 {
			headerObjects = append(headerObjects, map[string]string{
				strings.TrimSpace(parts[0]): strings.TrimSpace(parts[1]),
			})
		}
	}

	requestBody, err := json.Marshal(map[string]interface{}{
		"url":     url,
		"headers": headerObjects,
	})
	if err != nil {
		log.Printf("Error creating request body: %v", err)
		return
	}

	headers := map[string]string{
		"Content-Type": "application/json",
	}

	resp, err := makeRequest(ctx, http.MethodPost, endpoint, bytes.NewBuffer(requestBody), headers)
	if err != nil {
		log.Printf("Error making request: %v", err)
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error reading response: %v", err)
		return
	}

	var result struct {
		Message   string `json:"message"`
		JsmonID   string `json:"jsmonId"`
		Hash      string `json:"hash"`
		CreatedAt int64  `json:"createdAt"`
		URL       string `json:"url"`
	}

	if err := json.Unmarshal(body, &result); err != nil {
		log.Printf("Error parsing JSON: %v", err)
		return
	}

	fmt.Printf("{\n")
	fmt.Printf("    \"message\": \"%s\",\n", result.Message)
	fmt.Printf("    \"jsmonId\": \"%s\",\n", result.JsmonID)
	fmt.Printf("    \"hash\": \"%s\",\n", result.Hash)
	fmt.Printf("    \"createdAt\": %d,\n", result.CreatedAt)
	fmt.Printf("    \"url\": \"%s\"\n", result.URL)
	fmt.Printf("}\n")

	if result.JsmonID != "" {
		getAutomationResultsByJsmonId(result.JsmonID)
	}
}

func rescanUrlEndpoint(scanId string) {
	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	endpoint := fmt.Sprintf("%s/rescanURL/%s", apiBaseURL, scanId)

	headers := map[string]string{
		"Content-Type": "application/json",
	}

	resp, err := makeRequest(ctx, http.MethodPost, endpoint, nil, headers)
	if err != nil {
		log.Printf("Error making request: %v", err)
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error reading response: %v", err)
		return
	}

	var result interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		log.Printf("Error parsing JSON: %v", err)
		return
	}

	prettyJSON, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		log.Printf("Error formatting JSON: %v", err)
		return
	}

	fmt.Println(string(prettyJSON))
}

func getDomains() {
	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	endpoint := fmt.Sprintf("%s/getDomains", apiBaseURL)

	headers := map[string]string{
		"Content-Type": "application/json",
	}

	resp, err := makeRequest(ctx, http.MethodGet, endpoint, nil, headers)
	if err != nil {
		log.Printf("Error making request: %v", err)
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error reading response: %v", err)
		return
	}

	var domains []string
	if err := json.Unmarshal(body, &domains); err != nil {
		log.Printf("Error parsing JSON: %v", err)
		return
	}

	for _, domain := range domains {
		fmt.Println(domain)
	}
}

func createWordList(domains []string) {
	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	endpoint := fmt.Sprintf("%s/createWordList", apiBaseURL)

	requestBody := addWordlistRequest{
		Domains: domains,
	}
	body, err := json.Marshal(requestBody)
	if err != nil {
		log.Printf("failed to marshal request body: %v", err)
		return
	}

	headers := map[string]string{
		"Content-Type": "application/json",
	}

	resp, err := makeRequest(ctx, http.MethodPost, endpoint, bytes.NewBuffer(body), headers)
	if err != nil {
		log.Printf("failed to send request: %v", err)
		return
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("failed to read response body: %v", err)
		return
	}

	fmt.Printf("Word list:\n%s\n", string(responseBody))
}

func scanFileEndpoint(fileId string) {
	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	endpoint := fmt.Sprintf("%s/scanFile/%s", apiBaseURL, fileId)

	headers := map[string]string{
		"Content-Type": "application/json",
	}

	resp, err := makeRequest(ctx, http.MethodPost, endpoint, nil, headers)
	if err != nil {
		log.Printf("Error making request: %v", err)
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error reading response: %v", err)
		return
	}

	var result interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		log.Printf("Error parsing JSON: %v", err)
		return
	}

	prettyJSON, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		log.Printf("Error formatting JSON: %v", err)
		return
	}

	fmt.Println(string(prettyJSON))
}

func addCustomWordUser(words []string) {
	cleanedWords := []string{}
	for _, word := range words {
		if strings.TrimSpace(word) != "" {
			cleanedWords = append(cleanedWords, word)
		}
	}

	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Do you want to append or overwrite the custom words?")
	fmt.Println("1. Append")
	fmt.Println("2. Overwrite")
	fmt.Print("Select option (1 or 2): ")

	operationChoice, _ := reader.ReadString('\n')
	operationChoice = strings.TrimSpace(operationChoice)

	var operation string
	if operationChoice == "1" {
		operation = "append"
	} else if operationChoice == "2" {
		operation = "overwrite"
	} else {
		fmt.Println("Invalid option selected. Exiting.")
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	endpoint := fmt.Sprintf("%s/addCustomWords?operation=%s", apiBaseURL, operation)

	requestBody := addCustomWordsRequest{
		Words: cleanedWords,
	}
	body, err := json.Marshal(requestBody)
	if err != nil {
		log.Printf("failed to marshal request body: %v", err)
		return
	}

	headers := map[string]string{
		"Content-Type": "application/json",
	}

	resp, err := makeRequest(ctx, http.MethodPost, endpoint, bytes.NewBuffer(body), headers)
	if err != nil {
		log.Printf("failed to send request: %v", err)
		return
	}
	defer resp.Body.Close()

	var response map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		log.Printf("failed to unmarshal JSON response: %v", err)
		return
	}

	jsonData, err := json.MarshalIndent(response, "", "  ")
	if err != nil {
		log.Printf("failed to marshal response for pretty print: %v", err)
		return
	}

	fmt.Printf("%s\n", jsonData)
}

func urlsmultipleResponse() {
	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	endpoint := fmt.Sprintf("%s/urlWithMultipleResponse", apiBaseURL)

	headers := map[string]string{
		"Content-Type": "application/json",
	}

	resp, err := makeRequest(ctx, http.MethodGet, endpoint, nil, headers)
	if err != nil {
		log.Printf("Error making request: %v", err)
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error reading response: %v", err)
		return
	}

	var response struct {
		Message string   `json:"message"`
		Data    []string `json:"data"`
	}

	if err := json.Unmarshal(body, &response); err != nil {
		log.Printf("Error parsing JSON: %v", err)
		return
	}

	if len(response.Data) > 0 {
		for _, url := range response.Data {
			fmt.Println(url)
		}
	}
}

func uploadFileEndpoint(filePath string, headers []string) {
	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	endpoint := fmt.Sprintf("%s/uploadFile", apiBaseURL)

	headerMaps := []map[string]string{}
	for _, header := range headers {
		parts := strings.SplitN(header, ":", 2)
		if len(parts) == 2 {
			key := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])
			headerMaps = append(headerMaps, map[string]string{key: value})
		}
	}

	headersJSON, err := json.Marshal(headerMaps)
	if err != nil {
		log.Printf("Error marshaling headers to JSON: %v", err)
		return
	}

	queryParams := url.Values{}
	queryParams.Add("headers", string(headersJSON))
	endpoint = fmt.Sprintf("%s?%s", endpoint, queryParams.Encode())

	file, err := os.Open(filePath)
	if err != nil {
		log.Printf("Error opening file: %v", err)
		return
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		log.Printf("Error getting file info: %v", err)
		return
	}
	if fileInfo.Size() > 10*1024*1024 {
		log.Printf("File size exceeds limit")
		return
	}

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	h := make(textproto.MIMEHeader)
	h.Set("Content-Disposition", fmt.Sprintf(`form-data; name="%s"; filename="%s"`, "file", filepath.Base(filePath)))
	h.Set("Content-Type", "text/plain")
	part, err := writer.CreatePart(h)
	if err != nil {
		log.Printf("Error creating form part: %v", err)
		return
	}

	if _, err := io.Copy(part, file); err != nil {
		log.Printf("Error copying file content: %v", err)
		return
	}

	if err := writer.Close(); err != nil {
		log.Printf("Error closing multipart writer: %v", err)
		return
	}

	headers = map[string]string{
		"Content-Type":     writer.FormDataContentType(),
		"Accept-Encoding":  "gzip, deflate, br",
	}

	resp, err := makeRequest(ctx, http.MethodPost, endpoint, body, headers)
	if err != nil {
		log.Printf("Error making request: %v", err)
		return
	}
	defer resp.Body.Close()

	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error reading response: %v", err)
		return
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		log.Printf("Upload failed with status code: %d", resp.StatusCode)
		return
	}

	var result map[string]interface{}
	if err := json.Unmarshal(responseBody, &result); err != nil {
		log.Printf("Failed to parse JSON response: %v", err)
		return
	}

	fileID, ok := result["fileId"].(string)
	if !ok {
		log.Printf("Error: fileId is not a string")
		return
	}
	getAutomationResultsByFileId(fileID)
}

func getAllAutomationResults(input string, size int) {
	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	endpoint := fmt.Sprintf("%s/getAllAutomationResults", apiBaseURL)

	url := fmt.Sprintf("%s?showonly=all&inputType=domain&input=%s&size=%d", endpoint, input, size)

	headers := map[string]string{}

	resp, err := makeRequest(ctx, http.MethodGet, url, nil, headers)
	if err != nil {
		log.Printf("Error making request: %v", err)
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error reading response: %v", err)
		return
	}

	var result interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		log.Printf("Error parsing JSON: %v", err)
		return
	}

	prettyJSON, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		log.Printf("Error formatting JSON: %v", err)
		return
	}

	fmt.Println(string(prettyJSON))
}

func getScannerResults() {
	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	endpoint := fmt.Sprintf("%s/getScannerResults", apiBaseURL)

	headers := map[string]string{}

	resp, err := makeRequest(ctx, http.MethodGet, endpoint, nil, headers)
	if err != nil {
		log.Printf("Error making request: %v", err)
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error reading response: %v", err)
		return
	}

	var result struct {
		Message string `json:"message"`
		Data    struct {
			ModuleName []string `json:"moduleName"`
			URL        string   `json:"url"`
		} `json:"data"`
	}

	if err := json.Unmarshal(body, &result); err != nil {
		log.Printf("Error parsing JSON: %v", err)
		return
	}

	fmt.Println("Message:", result.Message)
	fmt.Println("URL:", result.Data.URL)
	fmt.Println("Modules found:")
	for _, module := range result.Data.ModuleName {
		fmt.Printf("- %s\n", module)
	}
}

func automateScanDomain(domain string, words []string) {
	fmt.Printf("automateScanDomain called with domain: %s and words: %v\n", domain, words)
	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	endpoint := fmt.Sprintf("%s/automateScanDomain", apiBaseURL)

	requestBody := AutomateScanDomainRequest{
		Domain: domain,
		Words:  words,
	}
	body, err := json.Marshal(requestBody)
	if err != nil {
		log.Printf("failed to marshal request body: %v", err)
		return
	}

	headers := map[string]string{
		"Content-Type": "application/json",
	}

	resp, err := makeRequest(ctx, http.MethodPost, endpoint, bytes.NewBuffer(body), headers)
	if err != nil {
		log.Printf("failed to send request: %v", err)
		return
	}
	defer resp.Body.Close()

	decoder := json.NewDecoder(resp.Body)
	for {
		var response map[string]interface{}
		err := decoder.Decode(&response)
		if err != nil {
			if err == io.EOF {
				break
			}
			log.Printf("failed to unmarshal JSON response: %v", err)
			return
		}

		printFormattedResponse(response)
	}
}

func printFormattedResponse(response map[string]interface{}) {
	fmt.Println("Message:", response["message"])
	fmt.Println("File ID:", response["fileId"])
	fmt.Println("Trimmed Domain:", response["trimmedDomain"])

	scanResponse, ok := response["scanResponse"].(map[string]interface{})
	if ok {
		fmt.Println("\nResult")
		fmt.Println("", scanResponse["message"])

		analysisResult, ok := scanResponse["analysis_result"].(map[string]interface{})
		if ok {
			fmt.Println("\n  Analysis Result:")
			fmt.Println("    ", analysisResult["message"])
			fmt.Println("    Total Chunks:", analysisResult["totalChunks"])
			fmt.Println("    Use -automationData flag to view all automation data for this domain")
		}

		moduleScanResult, ok := scanResponse["modulescan_result"].(map[string]interface{})
		if ok {
			fmt.Println("\n  Module Scan Result:")
			fmt.Println("    ", moduleScanResult["message"])
			fmt.Println("    Use -scannerData flag to view module scanner data for this domain")
		}
	}
}

func callViewProfile() {
	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	endpoint := fmt.Sprintf("%s/viewProfile", apiBaseURL)

	headers := map[string]string{
		"Content-Type": "application/json",
	}

	resp, err := makeRequest(ctx, http.MethodGet, endpoint, nil, headers)
	if err != nil {
		log.Printf("Error making request: %v", err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error reading response body: %v", err)
		os.Exit(1)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		log.Printf("Error unmarshalling response: %v", err)
		os.Exit(1)
	}

	if data, ok := result["data"].(map[string]interface{}); ok {
		var accountType string
		if orgFound, ok := data["orgFound"].(bool); ok && orgFound {
			accountType = "org"
		} else if personalProfile, ok := data["personalProfile"].(bool); ok && personalProfile {
			accountType = "user"
		} else {
			accountType = "unknown"
		}

		filteredResult := map[string]interface{}{
			"limits": data["apiCallLimits"],
			"type":   accountType,
		}
		filteredData, _ := json.MarshalIndent(filteredResult, "", "  ")
		fmt.Println(string(filteredData))
	} else {
		fmt.Println("Error: Invalid response format")
	}
}

func compareEndpoint(id1, id2 string) {
	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	endpoint := fmt.Sprintf("%s/compare", apiBaseURL)

	requestBody, err := json.Marshal(map[string]string{
		"id1": id1,
		"id2": id2,
	})
	if err != nil {
		log.Printf("Error creating request body: %v", err)
		return
	}

	headers := map[string]string{
		"Content-Type": "application/json",
	}

	resp, err := makeRequest(ctx, http.MethodPost, endpoint, bytes.NewBuffer(requestBody), headers)
	if err != nil {
		log.Printf("Error making request: %v", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("Unexpected status code: %d", resp.StatusCode)
		body, _ := ioutil.ReadAll(resp.Body)
		fmt.Printf("Response: %s\n", string(body))
		return
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error reading response: %v", err)
		return
	}

	var diffItems []DiffItem
	if err := json.Unmarshal(body, &diffItems); err != nil {
		log.Printf("Error parsing JSON: %v", err)
		fmt.Printf("Response: %s\n", string(body))
		return
	}

	addedCount := 0
	removedCount := 0

	fmt.Println("Summary of changes:")
	for _, item := range diffItems {
		if item.Added {
			addedCount++
			if addedCount <= 20 {
				fmt.Printf("+ %s\n", item.Value)
			}
		} else if item.Removed {
			removedCount++
			if removedCount <= 20 {
				fmt.Printf("- %s\n", item.Value)
			}
		}
	}

	fmt.Printf("\nTotal additions: %d\n", addedCount)
	fmt.Printf("Total removals: %d\n", removedCount)
	if addedCount > 5 {
		fmt.Printf("(Only the first 20 additions are shown)\n")
	}
	if removedCount > 5 {
		fmt.Printf("(Only the first 20 removals are shown)\n")
	}
}