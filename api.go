package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
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

type DataItem struct {
	JsmonId       string         `json:"jsmonId"`
	URL           string         `json:"url"`
	ModuleName    []string       `json:"moduleName"`
	DetectedWords []DetectedWord `json:"detectedWords"`
	CreatedAt     string         `json:"createdAt"`
}

type DetectedWord struct {
	Name  string   `json:"name"`
	Words []string `json:"words"`
}

type DomainResponse struct {
	Domains []string `json:"domains"`
}

type AutomateScanDomainRequest struct {
	Domain string   `json:"domain"`
	Words  []string `json:"words"`
}



type URLEntry struct {
	URL string `json:"url"`
}

// Function :

func uploadFileEndpoint(filePath string, headers []string, wkspId string) {
	endpoint := fmt.Sprintf("%s/uploadFile?wkspId=%s", apiBaseURL, wkspId)

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

	responseBody, err := io.ReadAll(resp.Body)
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
	getAutomationResultsByFileId(fileID, wkspId)

}

func automateScanDomain(domain string, words []string, wkspId string) error {
	// fmt.Printf("automateScanDomain called with domain: %s and words: %v\n", domain, words)
	endpoint := fmt.Sprintf("%s/automateScanDomain?wkspId=%s", apiBaseURL, wkspId)

	requestBody := AutomateScanDomainRequest{
		Domain: domain,
		Words:  words,
	}
	body, err := json.Marshal(requestBody)
	if err != nil {
		return fmt.Errorf("failed to marshal request body: %v", err)
	}

	client := &http.Client{}
	req, err := http.NewRequest("POST", endpoint, bytes.NewBuffer(body))
	if err != nil {
		return fmt.Errorf("error creating request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Jsmon-Key", strings.TrimSpace(getAPIKey()))

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %v", err)
	}
	var response map[string]interface{}
	if err := json.Unmarshal(responseBody, &response); err != nil {
		return fmt.Errorf("failed to parse response body: %v", err)
	}

	if resp.StatusCode == 200 {
		fmt.Printf("[INF] %s scanned successfully\n", domain)
	} else if resp.StatusCode == 401 {
		fmt.Printf("[ERR] Wrong API Key\n")
	} else {
		fmt.Printf("[INF] %s, error in scanning\n", domain)
	}
	return nil
}
