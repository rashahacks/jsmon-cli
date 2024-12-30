package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

type URLResponse struct {
	Urls    []URLItem `json:"urls"`
	Message string    `json:"Message"`
}

type URLItem struct {
	URL string `json:"url"`
}

func viewUrls(size int, wkspId string) {
	endpoint := fmt.Sprintf("%s/searchAllUrls?size=%d&start=0&wkspId=%s", apiBaseURL, size, wkspId) // Use the size parameter
	client := &http.Client{}
	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}

	req.Header.Set("X-Jsmon-Key", strings.TrimSpace(getAPIKey()))

	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("failed to send request: %v", err)
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("failed to read response body: %v", err)
		return
	}
	var response URLResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		fmt.Printf("failed to unmarshal JSON response: %v", err)
		return
	}

	for _, urlItem := range response.Urls {
		fmt.Println(urlItem.URL)
	}
}
