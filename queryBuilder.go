package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type QueryBuilderResponse struct {
    PaginatedResults []interface{} `json:"paginatedResults"`
    URLs            []string      `json:"urls"`
    IsReverseSearchAvailable bool `json:"isReverseSearchAvailable"`
}

func queryBuilder(wkspId, query string) {
    endpoint := fmt.Sprintf("%s/queryBuilder?wkspId=%s", apiBaseURL, wkspId)

    // Create request body with the query parameter
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

    // Handle URLs if present (for jsUrls field)
    if len(result.URLs) > 0 {
        for _, url := range result.URLs {
            fmt.Println(url)
        }
        return
    }

    // Handle paginated results (for exposures and other fields)
    if len(result.PaginatedResults) > 0 {
        // Pretty print the results
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
