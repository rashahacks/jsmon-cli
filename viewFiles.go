package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
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
	client := &http.Client{}
	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}

	req.Header.Set("X-Jsmon-Key", strings.TrimSpace(getAPIKey()))

	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Failed to send request: %v", err)
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Failed to read response body: %v", err)
		return
	}

	var response FileResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		fmt.Printf("Failed to unmarshal JSON response: %v", err)
		return
	}

	fmt.Println(response.Message)
	for _, fileItem := range response.Data {
		fmt.Printf("File Name: %s\nFile Size: %.3f MB\nFile ID: %s\nFile Key: %s\nNumber of URLs: %d\nCreated At: %s\n\n",
			fileItem.FileName, fileItem.FileSize, fileItem.FileID, fileItem.FileKey, fileItem.Urls, fileItem.CreatedAt)
	}
}
