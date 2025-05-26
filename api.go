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
	"unicode/utf8"
	// "net/url"
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


func uploadFileEndpoint(filePath string, headers []string, wkspId string) {
    endpoint := fmt.Sprintf("%s/uploadFile?wkspId=%s", apiBaseURL, wkspId)
    
    if len(headers) > 0 {
        headersJSON, err := json.Marshal(headers)
        if err != nil {
            log.Fatalf("Error marshaling headers: %v", err)
        }
        endpoint = fmt.Sprintf("%s&headers=%s", endpoint, string(headersJSON))
    }

    file, err := os.Open(filePath)
    if err != nil {
        log.Fatalf("Error opening file: %v", err)
    }
    defer file.Close()

    // Read and validate file content
    content, err := io.ReadAll(file)
    if err != nil {
        log.Fatalf("Error reading file: %v", err)
    }

    // Check if content is valid UTF-8
    if !utf8.Valid(content) {
        log.Fatalf("File content is not valid UTF-8")
    }

    // Count lines and validate URLs
    lines := strings.Split(string(content), "\n")
    validURLCount := 0
    for _, line := range lines {
        line = strings.TrimSpace(line)
        if line != "" && strings.HasPrefix(line, "http") {
            validURLCount++
        }
    }
    fmt.Printf("Found %d valid URLs in file\n", validURLCount)

    if validURLCount == 0 {
        log.Fatalf("No valid URLs found in file")
    }
    if validURLCount > 1000 {
        log.Fatalf("Too many URLs in file (max 1000)")
    }

    // Reset file pointer
    file.Seek(0, 0)

    // Create multipart form
    var requestBody bytes.Buffer
    writer := multipart.NewWriter(&requestBody)

    // Create form file part with explicit Content-Type
    h := make(textproto.MIMEHeader)
    h.Set("Content-Disposition", fmt.Sprintf(`form-data; name="file"; filename="%s"`, filepath.Base(filePath)))
    h.Set("Content-Type", "text/plain; charset=utf-8")
    
    part, err := writer.CreatePart(h)
    if err != nil {
        log.Fatalf("Error creating form file: %v", err)
    }

    _, err = io.Copy(part, bytes.NewReader(content))
    if err != nil {
        log.Fatalf("Error copying file data: %v", err)
    }

    err = writer.Close()
    if err != nil {
        log.Fatalf("Error closing writer: %v", err)
    }

    // Build HTTP request
    req, err := http.NewRequest("POST", endpoint, &requestBody)
    if err != nil {
        log.Fatalf("Error creating HTTP request: %v", err)
    }

    // Set required headers
    req.Header.Set("Content-Type", writer.FormDataContentType())
    req.Header.Set("X-Jsmon-Key", strings.TrimSpace(getAPIKey()))

    // Debug output
    // fmt.Printf("Sending request to: %s\n", endpoint)
    // fmt.Printf("Content-Type: %s\n", writer.FormDataContentType())
    // fmt.Printf("File size: %d bytes\n", len(content))
    // fmt.Printf("API Key: %s\n", maskAPIKey(strings.TrimSpace(getAPIKey())))

    // Execute request
    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        log.Fatalf("Upload request failed: %v", err)
    }
    defer resp.Body.Close()

    // Read and handle response
    bodyBytes, _ := io.ReadAll(resp.Body)
    
    if resp.StatusCode == http.StatusOK {
        fmt.Println("File uploaded successfully!")
        fmt.Println("Response:", string(bodyBytes))
    } else {
        fmt.Printf("Upload failed with status code: %d\n", resp.StatusCode)
        // fmt.Printf("Response headers: %v\n", resp.Header)
        // fmt.Println("Response body:", string(bodyBytes))
    }
}

// Helper function to mask API key for logging
func maskAPIKey(key string) string {
    if len(key) <= 8 {
        return "****"
    }
    return key[:4] + "..." + key[len(key)-4:]
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
