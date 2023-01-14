package main

import (
	"io"
	"net/http"
)

// for control over HTTP client headers,
// redirect policy, and other settings,
// A Client is an HTTP client
func httpGet(client *http.Client, url string, params map[string]string) (string, error) {
	// Create a new request using http
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}

	// add parameters to request
	q := req.URL.Query()
	for key, value := range params {
		q.Add(key, value)
	}
	req.URL.RawQuery = q.Encode()

	// Send req using http Client
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// Call ReadAll to get the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}
