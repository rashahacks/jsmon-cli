package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	// "time"
)

func urlsmultipleResponse(wkspId string) {
	endpoint := fmt.Sprintf("%s/urlWithMultipleResponse?wkspId=%s", apiBaseURL, wkspId)
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
		Message string `json:"message"`
		Data    []struct {
			URL string `json:"url"`
		} `json:"data"`
	}

	err = json.Unmarshal(body, &response)
	if err != nil {
		fmt.Println("Error parsing JSON:", err)
		return
	}

	if len(response.Data) > 0 {
		for _, url := range response.Data {
			fmt.Println(url.URL)
		}
	}
}
