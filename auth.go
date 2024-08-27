package main

import (
	"io/ioutil"
	"os"
	"path/filepath"
)

var apiKey string

func loadAPIKey() error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	credPath := filepath.Join(homeDir, ".jsmon", "credentials")
	data, err := ioutil.ReadFile(credPath)
	if err != nil {
		return err
	}

	apiKey = string(data)
	return nil
}

func setAPIKey(key string) {
	apiKey = key
}

func getAPIKey() string {
	return apiKey
}
