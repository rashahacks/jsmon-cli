package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

func uploadUrlEndpoint(url string, customHeaders []string, wkspId string) error {
	endpoint := fmt.Sprintf("%s/uploadUrl?wkspId=%s", apiBaseURL, wkspId)

	headerObjects := make([]map[string]string, 0)
	for _, header := range customHeaders {
		parts := strings.SplitN(header, ":", 2)
		if len(parts) == 2 {
			headerObjects = append(headerObjects, map[string]string{
				strings.TrimSpace(parts[0]): strings.TrimSpace(parts[1]),
			})
		}
	}

	requestBody, err := json.Marshal(map[string]interface{}{
		"url":     url,
		"headers": headerObjects,
	})
	if err != nil {
		return fmt.Errorf("error creating request body: %v", err)
	}

	req, err := http.NewRequest("POST", endpoint, bytes.NewBuffer(requestBody))
	if err != nil {
		return fmt.Errorf("error creating request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Jsmon-Key", strings.TrimSpace(getAPIKey()))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error sending request: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("error reading response: %v", err)
	}

	// Parse response into a generic map
	var response map[string]interface{}
	if err := json.Unmarshal(body, &response); err != nil {
		return fmt.Errorf("error parsing JSON response: %v", err)
	}

	// Pretty print the response
	// prettyJSON, err := json.MarshalIndent(response, "", "    ")
	// if err != nil {
	// 	return fmt.Errorf("error formatting JSON response: %v", err)
	// }
	if messageArr, ok := response["message"].([]interface{}); ok {
		for _, msg := range messageArr {
			fmt.Println(msg.(string))
			if msg.(string) == "You don't have access to this workspace or it doesn't exist." {
				showAvailableWorkspaces()
			}
		}
	} else if msgStr, ok := response["message"].(string); ok {
		fmt.Println(msgStr)
		if msgStr == "You don't have access to this workspace or it doesn't exist." {
			showAvailableWorkspaces()
		}
	} else {
		fmt.Println("[ERR] Unexpected message format in response")
	}

	// Check for jsmonId or fileId to determine if we need to get automation results
	if jsmonID, ok := response["jsmonId"].(string); ok && jsmonID != "" {
		getAutomationResultsByJsmonId(jsmonID, wkspId)
	} else if fileID, ok := response["fileId"].(string); ok && fileID != "" {
		// You can add handling for fileId here if needed
		fmt.Printf("File ID received: %s\n", fileID)
	}

	return nil
}
