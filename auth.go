package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

var apiKey string

func loadAPIKey() error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get user home directory: %w", err)
	}

	credPath := filepath.Join(homeDir, ".jsmon", "credentials")
	data, err := os.ReadFile(credPath)
	if err != nil {
		return fmt.Errorf("failed to read credentials file: %w", err)
	}

	apiKey = strings.TrimSpace(string(data))
	return nil
}

func setAPIKey(key string) {
	apiKey = strings.TrimSpace(key)
}

func getAPIKey() string {
	return apiKey
}
