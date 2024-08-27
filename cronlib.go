package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

type Response struct {
	Message string `json:"message"`
	Data    string `json:"data"`
}

func StartCron(cronNotification string, cronTime int64, cronType string, cronDomain string, cronDomainNotify string) {

	notification := strings.TrimSpace(cronNotification)
	//split vulnerabilities
	vulnerabilitiesType := strings.Split(cronType, ",")
	cronDomains := strings.Split(cronDomain, ",")
	cronDomainsNotify := strings.Split(cronDomainNotify, ",")
	if len(cronDomains) != len(cronDomainsNotify) {
		fmt.Println("Invalid format for cronDomains and cronDomainsNotify. Use: domain1,domain2,domain3 domainNotify1,domainNotify2,domainNotify3")
		return
	}

	//trim domains and domainsNotify
	for i := 0; i < len(cronDomains); i++ {
		cronDomains[i] = strings.TrimSpace(cronDomains[i])
		cronDomainsNotify[i] = strings.TrimSpace(cronDomainsNotify[i])
	}
	//create domains map
	var domains []map[string]interface{}
	for i := 0; i < len(cronDomains); i++ {
		notify := strings.EqualFold(cronDomainsNotify[i], "true")
		domain := map[string]interface{}{
			"domain": cronDomains[i],
			"notify": notify,
		}
		domains = append(domains, domain)
	}

	apiKey := strings.TrimSpace(getAPIKey())
	baseUrl := apiBaseURL
	client := &http.Client{}

	var method = "PUT"
	var url = baseUrl + "/startCron"
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		fmt.Printf("failed to create request: %v", err)
		return
	}
	req.Header.Set("X-Jsmon-Key", apiKey)
	req.Header.Set("Content-Type", "application/json")

	data := map[string]interface{}{
		"notificationChannel": notification,
		"vulnerabilitiesType": vulnerabilitiesType,
		"time":                cronTime,
		"domains":             domains,
	}
	jsonData, err := json.Marshal(data)
	if err != nil {
		fmt.Printf("failed to marshal JSON: %v", err)
		return
	}

	req.Body = ioutil.NopCloser(bytes.NewReader(jsonData))

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
	var response Response
	err = json.Unmarshal(body, &response)
	if err != nil {
		fmt.Printf("failed to unmarshal JSON response: %v", err)
		return
	}

	fmt.Println("Message:", response.Message)

}

func StopCron() {
	apiKey := strings.TrimSpace(getAPIKey())
	baseUrl := apiBaseURL
	client := &http.Client{}
	var method = "PUT"
	var url = baseUrl + "/stopCron"
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		fmt.Printf("failed to create request: %v", err)
		return
	}
	req.Header.Set("X-Jsmon-Key", apiKey)
	req.Header.Set("Content-Type", "application/json")

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

	var response Response
	err = json.Unmarshal(body, &response)
	if err != nil {
		fmt.Printf("failed to unmarshal JSON response: %v", err)
		return
	}

	fmt.Println("Message:", response.Message)

}

func UpdateCron(cronNotification string, cronType string, cronDomain string, cronDomainNotify string, cronTime int64) {
	apiKey := strings.TrimSpace(getAPIKey())
	baseUrl := apiBaseURL

	client := &http.Client{}
	var method = "PUT"
	var url = baseUrl + "/updateCron"

	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		fmt.Printf("failed to create request: %v", err)
		return
	}

	req.Header.Set("X-Jsmon-Key", apiKey)
	req.Header.Set("Content-Type", "application/json")

	data := map[string]interface{}{}
	if cronNotification != "" {
		data["notificationChannel"] = cronNotification
	}
	if cronType != "" {
		vulnerabilitiesType := strings.Split(cronType, ",")
		data["vulnerabilitiesType"] = vulnerabilitiesType
	}
	if cronTime != 0 {
		data["time"] = cronTime
	}
	if cronDomain != "" && cronDomainNotify != "" {
		cronDomains := strings.Split(cronDomain, ",")
		cronDomainsNotify := strings.Split(cronDomainNotify, ",")
		if len(cronDomains) != len(cronDomainsNotify) {
			fmt.Println("Invalid format for cronDomains and cronDomainsNotify. Use: domain1,domain2,domain3 domainNotify1,domainNotify2,domainNotify3")
			return
		}
		//trim domains and domainsNotify
		for i := 0; i < len(cronDomains); i++ {
			cronDomains[i] = strings.TrimSpace(cronDomains[i])
			cronDomainsNotify[i] = strings.TrimSpace(cronDomainsNotify[i])
		}
		//create domains map
		var domains []map[string]interface{}
		for i := 0; i < len(cronDomains); i++ {
			notify := strings.EqualFold(cronDomainsNotify[i], "true")
			domain := map[string]interface{}{
				"domain": cronDomains[i],
				"notify": notify,
			}
			domains = append(domains, domain)
		}
		data["domains"] = domains
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		fmt.Printf("failed to marshal JSON: %v", err)
		return
	}

	req.Body = ioutil.NopCloser(bytes.NewReader(jsonData))

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
	var response Response
	err = json.Unmarshal(body, &response)
	if err != nil {
		fmt.Printf("failed to unmarshal JSON response: %v", err)
		return
	}

	fmt.Println("Message:", response.Message)

}
