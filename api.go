package main

import (
	"bufio"
	"bytes"
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
	// "time"
)

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

func rescanDomain(domain string) {
	endpoint := fmt.Sprintf("%s/rescanDomain", apiBaseURL)

	requestBody, err := json.Marshal(map[string]string{
		"domain": domain,
	})
	if err != nil {
		fmt.Println("Error creating request body:", err)
		return
	}

	req, err := http.NewRequest("POST", endpoint, bytes.NewBuffer(requestBody))
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

	var result RescanDomainResponse
	err = json.Unmarshal(body, &result)
	if err != nil {
		fmt.Println("Error parsing JSON:", err)
		return
	}

	fmt.Printf("Message: %s\n", result.Message)
	fmt.Printf("Total URLs submitted for scanning: %d\n", result.TotalUrls)
}

func totalAnalysisData() {
	endpoint := fmt.Sprintf("%s/totalCountAnalysisData", apiBaseURL)

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
func searchUrlsByDomain(domain string) {
	endpoint := fmt.Sprintf("%s/searchUrlbyDomain?domain=%s", apiBaseURL, url.QueryEscape(domain))

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

	body, err := ioutil.ReadAll(resp.Body)
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

	fmt.Printf("Message: %s\n", result.Message)
	fmt.Printf("Total URLs: %d\n", result.TotalUrls)
	fmt.Println("URLs:")
	for _, entry := range result.URLs {
		fmt.Printf("- %s\n", entry.URL)
	}

}

func uploadUrlEndpoint(url string, customHeaders []string) {

	endpoint := fmt.Sprintf("%s/uploadUrl", apiBaseURL)

	// Call the function : function is in getResultsByJsmonID
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
		fmt.Println("Error creating request body:", err)
		return
	}

	req, err := http.NewRequest("POST", endpoint, bytes.NewBuffer(requestBody))
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

	var result map[string]interface{}
	err = json.Unmarshal(body, &result)
	if err != nil {
		fmt.Println("Error parsing JSON:", err)
		return
	}

	if jsmonId, ok := result["jsmonId"].(string); ok {
		getAutomationResultsByJsmonId(jsmonId)
	}
	// if hash, ok := result["hash"].(string); ok {
	// 	fmt.Printf("Hash: %s\n", hash)
	// }
	// if createdAt, ok := result["createdAt"].(float64); ok {
	// 	timestamp := time.Unix(int64(createdAt), 0)
	// 	fmt.Printf("Created At: %s\n", timestamp.Format(time.RFC3339))
	// }
	// if url, ok := result["url"].(string); ok {
	// 	fmt.Printf("url: %s\n", url)
	// }
	// if message, ok := result["message"].(string); ok {
	// 	fmt.Printf("Message: %s\n", message)
	// }
}

// Function :
func rescanUrlEndpoint(scanId string) {
	endpoint := fmt.Sprintf("%s/rescanURL/%s", apiBaseURL, scanId)

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

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response:", err)
		return
	}

	var result interface{}
	err = json.Unmarshal(body, &result)
	if err != nil {
		fmt.Println("Error parsing JSON:", err)
		return
	}

	// Pretty print JSON
	prettyJSON, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		fmt.Println("Error formatting JSON:", err)
		return
	}

	fmt.Println(string(prettyJSON))
}

func getDomains() {
	endpoint := fmt.Sprintf("%s/getDomains", apiBaseURL)

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

	var domains []string
	err = json.Unmarshal(body, &domains)
	if err != nil {
		fmt.Println("Error parsing JSON:", err)
		return
	}

	// Print each domain on a new line
	for _, domain := range domains {
		fmt.Println(domain)
	}
}

func createWordList(domains []string) {
	endpoint := fmt.Sprintf("%s/createWordList", apiBaseURL)

	requestBody := addWordlistRequest{
		Domains: domains,
	}
	body, err := json.Marshal(requestBody)
	if err != nil {
		fmt.Printf("failed to marshal request body: %v\n", err)
		return
	}

	// Create HTTP request
	client := &http.Client{}
	req, err := http.NewRequest("POST", endpoint, bytes.NewBuffer(body))
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Jsmon-Key", strings.TrimSpace(getAPIKey()))

	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("failed to send request: %v\n", err)
		return
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("failed to read response body: %v\n", err)
		return
	}

	fmt.Printf("Word list:\n%s\n", string(responseBody))
}

func scanFileEndpoint(fileId string) {
	endpoint := fmt.Sprintf("%s/scanFile/%s", apiBaseURL, fileId)

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

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response:", err)
		return
	}

	var result interface{}
	err = json.Unmarshal(body, &result)
	if err != nil {
		fmt.Println("Error parsing JSON:", err)
		return
	}

	prettyJSON, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		fmt.Println("Error formatting JSON:", err)
		return
	}

	fmt.Println(string(prettyJSON))
}

func addCustomWordUser(words []string) {
	// Remove empty strings from the words slice
	cleanedWords := []string{}
	for _, word := range words {
		if strings.TrimSpace(word) != "" {
			cleanedWords = append(cleanedWords, word)
		}
	}

	// Prompt user for operation: append or overwrite
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

	// Append the selected operation to the endpoint as a query parameter
	endpoint := fmt.Sprintf("%s/addCustomWords?operation=%s", apiBaseURL, operation)

	// Create request body
	requestBody := addCustomWordsRequest{
		Words: cleanedWords,
	}
	body, err := json.Marshal(requestBody)
	if err != nil {
		fmt.Printf("failed to marshal request body: %v\n", err)
		return
	}

	// Create HTTP request
	client := &http.Client{}
	req, err := http.NewRequest("POST", endpoint, bytes.NewBuffer(body))
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Jsmon-Key", strings.TrimSpace(getAPIKey()))

	// Send the request
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("failed to send request: %v\n", err)
		return
	}
	defer resp.Body.Close()

	// Decode and pretty-print the response
	var response map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		fmt.Printf("failed to unmarshal JSON response: %v\n", err)
		return
	}

	jsonData, err := json.MarshalIndent(response, "", "  ")
	if err != nil {
		fmt.Printf("failed to marshal response for pretty print: %v\n", err)
		return
	}

	fmt.Printf("%s\n", jsonData)
}
func urlsmultipleResponse() {
	endpoint := fmt.Sprintf("%s/urlWithMultipleResponse", apiBaseURL)
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

	var response struct {
		Message string   `json:"message"`
		Data    []string `json:"data"`
	}

	err = json.Unmarshal(body, &response)
	if err != nil {
		fmt.Println("Error parsing JSON:", err)
		return
	}

	if len(response.Data) > 0 {
		for _, url := range response.Data {
			fmt.Println(url)
		}
	}
}

func uploadFileEndpoint(filePath string, headers []string) {
	endpoint := fmt.Sprintf("%s/uploadFile", apiBaseURL)

	headerMaps := []map[string]string{}

	// Parse headers into the correct format
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
		log.Fatalf("Error marshaling headers to JSON: %v", err)
	}

	// Create query parameters
	queryParams := url.Values{}
	queryParams.Add("headers", string(headersJSON))

	// Append query parameters to the endpoint URL
	endpoint = fmt.Sprintf("%s?%s", endpoint, queryParams.Encode())

	// Log the final endpoint URL for debugging
	// log.Printf("Final endpoint URL: %s", endpoint)

	file, err := os.Open(filePath)
	if err != nil {
		log.Fatalf("Error opening file: %v", err)
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		log.Fatalf("Error getting file info: %v", err)
	}
	if fileInfo.Size() > 10*1024*1024 {
		log.Fatalf("File size exceeds limit")
	}

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	h := make(textproto.MIMEHeader)
	h.Set("Content-Disposition", fmt.Sprintf(`form-data; name="%s"; filename="%s"`, "file", filepath.Base(filePath)))
	h.Set("Content-Type", "text/plain")
	part, err := writer.CreatePart(h)
	if err != nil {
		log.Fatalf("Error creating form part: %v", err)
	}

	_, err = io.Copy(part, file)
	if err != nil {
		log.Fatalf("Error copying file content: %v", err)
	}

	err = writer.Close()
	if err != nil {
		log.Fatalf("Error closing multipart writer: %v", err)
	}

	req, err := http.NewRequest("POST", endpoint, body)
	if err != nil {
		log.Fatalf("Error creating request: %v", err)
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("X-Jsmon-Key", strings.TrimSpace(getAPIKey()))
	req.Header.Set("Accept-Encoding", "gzip, deflate, br")

	// log.Printf("File being uploaded: %s", filepath.Base(filePath))
	// log.Printf("Request body length: %d bytes", body.Len())
	// log.Printf("Request body content (first 200 bytes): %s", body.String()[:min(200, body.Len())])

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Error sending request: %v", err)
	}
	defer resp.Body.Close()

	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Error reading response: %v", err)
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		log.Fatalf("Upload failed with status code: %d", resp.StatusCode)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(responseBody, &result); err != nil {
		log.Fatalf("Failed to parse JSON response: %v", err)
	}

	fileID, ok := result["fileId"].(string)
	if !ok {
		fmt.Println("Error: fileId is not a string")
		return
	}
	getAutomationResultsByFileId(fileID)
	// Print the response in a more user-friendly format
	// fmt.Println("File ID: \n",result["fileId"])
	// if jsmonId, ok := result["jsmonId"].(string); ok {
	// 	fmt.Printf("JSMON ID: %s\n", jsmonId)
	// }
	// if hash, ok := result["hash"].(string); ok {
	// 	fmt.Printf("Hash: %s\n", hash)
	// }
	// if createdAt, ok := result["createdAt"].(float64); ok {
	// 	timestamp := time.Unix(int64(createdAt), 0)
	// 	fmt.Printf("Created At: %s\n", timestamp.Format(time.RFC3339))
	// }
	// if message, ok := result["message"].(string); ok {
	// 	fmt.Printf("Message: %s\n", message)
	// }
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func getAllAutomationResults(input string, size int) {
	endpoint := fmt.Sprintf("%s/getAllAutomationResults", apiBaseURL)

	url := fmt.Sprintf("%s?showonly=all&inputType=domain&input=%s&size=%d", endpoint, input, size)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}

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

	var result interface{}
	err = json.Unmarshal(body, &result)
	if err != nil {
		fmt.Println("Error parsing JSON:", err)
		return
	}

	prettyJSON, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		fmt.Println("Error formatting JSON:", err)
		return
	}

	fmt.Println(string(prettyJSON))
}
func getScannerResults() {
	endpoint := fmt.Sprintf("%s/getScannerResults", apiBaseURL)

	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}

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

	var result struct {
		Message string `json:"message"`
		Data    struct {
			ModuleName []string `json:"moduleName"`
			URL        string   `json:"url"`
		} `json:"data"`
	}

	err = json.Unmarshal(body, &result)
	if err != nil {
		fmt.Println("Error parsing JSON:", err)
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
	endpoint := fmt.Sprintf("%s/automateScanDomain", apiBaseURL)

	requestBody := AutomateScanDomainRequest{
		Domain: domain,
		Words:  words,
	}
	body, err := json.Marshal(requestBody)
	if err != nil {
		fmt.Printf("failed to marshal request body: %v\n", err)
		return
	}

	client := &http.Client{}
	req, err := http.NewRequest("POST", endpoint, bytes.NewBuffer(body))
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Jsmon-Key", strings.TrimSpace(getAPIKey()))

	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("failed to send request: %v\n", err)
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
			fmt.Printf("failed to unmarshal JSON response: %v\n", err)
			return
		}

		printFormattedResponse(response)
	}
}

// func printFormattedResponse(response map[string]interface{}) {
// 	fmt.Println("Message:", response["message"])
// 	fmt.Println("File ID:", response["fileId"])
// 	fmt.Println("Trimmed Domain:", response["trimmedDomain"])

// 	scanResponse, ok := response["scanResponse"].(map[string]interface{})
// 	if ok {
// 		fmt.Println("\nScan Response:")
// 		fmt.Println("  Message:", scanResponse["message"])

// 		analysisResult, ok := scanResponse["analysis_result"].(map[string]interface{})
// 		if ok {
// 			fmt.Println("\n  Analysis Result:")
// 			fmt.Println("    Message:", analysisResult["message"])
// 			fmt.Println("    Total Chunks:", analysisResult["totalChunks"])
// 		}

//			moduleScanResult, ok := scanResponse["modulescan_result"].(map[string]interface{})
//			if ok {
//				fmt.Println("\n  Module Scan Result:")
//				fmt.Println("    Message:", moduleScanResult["message"])
//				modules, ok := moduleScanResult["data"].([]interface{})
//				if ok {
//					for _, module := range modules {
//						m := module.(map[string]interface{})
//						fmt.Println("    Module Name:", m["moduleName"])
//						fmt.Println("    URL:", m["url"])
//						fmt.Println()
//					}
//				}
//			}
//		}
//	}
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
	endpoint := fmt.Sprintf("%s/viewProfile", apiBaseURL)

	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
		os.Exit(1)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Jsmon-Key", strings.TrimSpace(getAPIKey()))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error making request:", err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		os.Exit(1)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		fmt.Println("Error unmarshalling response:", err)
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
	endpoint := fmt.Sprintf("%s/compare", apiBaseURL)

	requestBody, err := json.Marshal(map[string]string{
		"id1": id1,
		"id2": id2,
	})
	if err != nil {
		fmt.Println("Error creating request body:", err)
		return
	}

	req, err := http.NewRequest("POST", endpoint, bytes.NewBuffer(requestBody))
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

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("Unexpected status code: %d\n", resp.StatusCode)
		body, _ := ioutil.ReadAll(resp.Body)
		fmt.Printf("Response: %s\n", string(body))
		return
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response:", err)
		return
	}

	var diffItems []DiffItem
	err = json.Unmarshal(body, &diffItems)
	if err != nil {
		fmt.Println("Error parsing JSON:", err)
		fmt.Printf("Response: %s\n", string(body))
		return
	}

	addedCount := 0
	removedCount := 0

	fmt.Println("Summary of changes:")
	for _, item := range diffItems {
		if item.Added {
			addedCount++
			if addedCount <= 20 { // Print the first 20 additions
				fmt.Printf("+ %s\n", item.Value)
			}
		} else if item.Removed {
			removedCount++
			if removedCount <= 20 { // Print the first 20 removals
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
