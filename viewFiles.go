package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

type FileResponse struct {
	Message string     `json:"message"`
	Data    []FileItem `json:"data"`
}

type FileItem struct {
	FileID    string    `json:"fileId"`
	FileSize  float64   `json:"fileSize"`
	FileName  string    `json:"fileName"`
	FileKey   string    `json:"fileKey"`
	Urls      int       `json:"urls"`
	CreatedAt time.Time `json:"createdAt"`
}

func viewFiles() {
	endpoint := fmt.Sprintf("%s/viewFiles", apiBaseURL)
	
	req, err := http.NewRequest(http.MethodGet, endpoint, nil)
	if err != nil {
		fmt.Printf("Error creating request: %v\n", err)
		return
	}

	req.Header.Set("X-Jsmon-Key", strings.TrimSpace(getAPIKey()))

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Failed to send request: %v\n", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("Unexpected status code: %d\n", resp.StatusCode)
		return
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Failed to read response body: %v\n", err)
		return
	}

	var response FileResponse
	if err := json.Unmarshal(body, &response); err != nil {
		fmt.Printf("Failed to unmarshal JSON response: %v\n", err)
		return
	}

	fmt.Println(response.Message)
	for _, fileItem := range response.Data {
		fmt.Printf("File Name: %s\nFile Size: %.3f MB\nFile ID: %s\nFile Key: %s\nNumber of URLs: %d\nCreated At: %s\n\n",
			fileItem.FileName, fileItem.FileSize, fileItem.FileID, fileItem.FileKey, fileItem.Urls, fileItem.CreatedAt.Format(time.RFC3339))
	}
}
