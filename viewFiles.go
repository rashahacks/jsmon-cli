package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type FileResponse struct {
	Message string     `json:"message"`
	Data    []FileItem `json:"data"`
}

type FileItem struct {
	FileID    string  `json:"fileId"`
	FileSize  float64 `json:"fileSize"`
	FileName  string  `json:"fileName"`
	FileKey   string  `json:"fileKey"`
	Urls      int     `json:"urls"`
	CreatedAt string  `json:"createdAt"`
}

func viewFiles(wkspId string) {
	endpoint := fmt.Sprintf("%s/viewFiles?wkspId=%s", apiBaseURL, wkspId)
	
	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		fmt.Printf("[ERR] Error creating request: %v\n", err)
		return
	}

	req.Header.Set("X-Jsmon-Key", strings.TrimSpace(getAPIKey()))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("[ERR] Failed to send request: %v\n", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized {
		fmt.Println("[ERR] Wrong API key")
		return
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("[ERR] Failed to read response body: %v\n", err)
		return
	}

	var response FileResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		fmt.Printf("[ERR] Failed to unmarshal JSON response: %v\n", err)
		return
	}

	if len(response.Data) > 0 {
		jsonResponse, err := json.MarshalIndent(response.Data, "", "  ")
		if err != nil {
			fmt.Printf("[ERR] Failed to marshal JSON: %v\n", err)
			return
		}
		fmt.Println(string(jsonResponse))
	} else {
		fmt.Println("[INF] No files found.")
	}
}
