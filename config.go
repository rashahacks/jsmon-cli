package main

import (
	"net/http"
	"time"
)

const (
	apiBaseURL = "https://api.jsmon.sh/api/v2"
	credFile   = "~/.jsmon/credentials"
	timeout    = 10 * time.Second //
)

var httpClient *http.Client

func init() {
	httpClient = &http.Client{
		Timeout: timeout,
	}
}

// I've added a new constant timeout to store the timeout duration.
// This timeout is used to set the Timeout field of the http.Client.
// The http.Client is then used in various parts of the code to make HTTP requests.
// This ensures that all HTTP requests have a consistent timeout setting.
// The init() function is a special function in Go that is automatically called when the package is initialized.
// It's used here to initialize the httpClient with a timeout.
