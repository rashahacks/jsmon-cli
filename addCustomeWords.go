package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	// "time"
)

type addCustomWordsRequest struct {
	Words []string `json:"words"`
}

func addCustomWordUser(words []string, wkspId string) {
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
	endpoint := fmt.Sprintf("%s/addCustomWords?operation=%s&wkspId=%s", apiBaseURL, operation, wkspId)

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
