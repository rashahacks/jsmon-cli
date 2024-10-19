package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

type Response struct {
	Message string `json:"message"`
	Data    string `json:"data"`
}

type Domain struct {
	Domain string `json:"domain"`
	Notify bool   `json:"notify"`
}

func StartCron(cronNotification string, cronTime int64, cronType string, cronDomain string, cronDomainNotify string) error {
	notification := strings.TrimSpace(cronNotification)
	vulnerabilitiesType := strings.Split(cronType, ",")
	cronDomains := strings.Split(cronDomain, ",")
	cronDomainsNotify := strings.Split(cronDomainNotify, ",")

	if len(cronDomains) != len(cronDomainsNotify) {
		return fmt.Errorf("invalid format for cronDomains and cronDomainsNotify. Use: domain1,domain2,domain3 domainNotify1,domainNotify2,domainNotify3")
	}

	domains := make([]Domain, len(cronDomains))
	for i := range cronDomains {
		domains[i] = Domain{
			Domain: strings.TrimSpace(cronDomains[i]),
			Notify: strings.EqualFold(strings.TrimSpace(cronDomainsNotify[i]), "true"),
		}
	}

	data := map[string]interface{}{
		"notificationChannel": notification,
		"vulnerabilitiesType": vulnerabilitiesType,
		"time":                cronTime,
		"domains":             domains,
	}

	return sendRequest("PUT", "/startCron", data)
}

func StopCron() error {
	return sendRequest("PUT", "/stopCron", nil)
}

func UpdateCron(cronNotification string, cronType string, cronDomain string, cronDomainNotify string, cronTime int64) error {
	data := make(map[string]interface{})

	if cronNotification != "" {
		data["notificationChannel"] = cronNotification
	}
	if cronType != "" {
		data["vulnerabilitiesType"] = strings.Split(cronType, ",")
	}
	if cronTime != 0 {
		data["time"] = cronTime
	}
	if cronDomain != "" && cronDomainNotify != "" {
		cronDomains := strings.Split(cronDomain, ",")
		cronDomainsNotify := strings.Split(cronDomainNotify, ",")
		if len(cronDomains) != len(cronDomainsNotify) {
			return fmt.Errorf("invalid format for cronDomains and cronDomainsNotify. Use: domain1,domain2,domain3 domainNotify1,domainNotify2,domainNotify3")
		}

		domains := make([]Domain, len(cronDomains))
		for i := range cronDomains {
			domains[i] = Domain{
				Domain: strings.TrimSpace(cronDomains[i]),
				Notify: strings.EqualFold(strings.TrimSpace(cronDomainsNotify[i]), "true"),
			}
		}
		data["domains"] = domains
	}

	return sendRequest("PUT", "/updateCron", data)
}

func sendRequest(method, endpoint string, data interface{}) error {
	apiKey := strings.TrimSpace(getAPIKey())
	url := apiBaseURL + endpoint

	var body io.Reader
	if data != nil {
		jsonData, err := json.Marshal(data)
		if err != nil {
			return fmt.Errorf("failed to marshal JSON: %w", err)
		}
		body = bytes.NewReader(jsonData)
	}

	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("X-Jsmon-Key", apiKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	var response Response
	if err := json.Unmarshal(responseBody, &response); err != nil {
		return fmt.Errorf("failed to unmarshal JSON response: %w", err)
	}

	fmt.Println("Message:", response.Message)
	return nil
}
