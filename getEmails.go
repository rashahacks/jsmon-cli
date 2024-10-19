package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// getEmails retrieves emails for the given domains
func getEmails(domains []string) error {
	endpoint := fmt.Sprintf("%s/getEmails", apiBaseURL)
	requestBody, err := json.Marshal(map[string]interface{}{
		"domains": domains,
	})
	if err != nil {
		return fmt.Errorf("error creating request body: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, endpoint, bytes.NewBuffer(requestBody))
	if err != nil {
		return fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Jsmon-Key", strings.TrimSpace(getAPIKey()))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error sending request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("error reading response: %w", err)
	}

	var result struct {
		Emails []string `json:"emails"`
	}
	if err := json.Unmarshal(body, &result); err != nil {
		return fmt.Errorf("error parsing JSON: %w", err)
	}

	for _, email := range result.Emails {
		fmt.Println(email)
	}

	return nil
}
