package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

type URLResponse struct {
	Urls    []URLItem `json:"urls"`
	Message string    `json:"Message"`
}

type URLItem struct {
	URL string `json:"url"`
}

func viewUrls(size int) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	endpoint := fmt.Sprintf("%s/searchAllUrls?size=%d&start=0", apiBaseURL, size)
	
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
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

	var response URLResponse
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		fmt.Printf("Failed to decode JSON response: %v\n", err)
		return
	}

	for _, urlItem := range response.Urls {
		fmt.Println(urlItem.URL)
	}
}
