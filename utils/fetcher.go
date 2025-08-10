// Package utils provides utility functions and helpers for the SeaaveyBot application.
// This file specifically contains HTTP client functionality for making API requests
// to external services, particularly the Seaavey API endpoints.
package utils

import (
	"aemy/types"
	"io"
	"net/http"
	"net/url"
	"time"
)

// BaseAPI defines the base URL for the Seaavey API service.
// This constant serves as the root endpoint for all API interactions.
const BaseAPI = "https://api.seaavey.my.id/api"

// httpClient is a shared HTTP client with a 30-second timeout.
var httpClient = &http.Client{
	Timeout: 30 * time.Second,
}

// SeaaveyAPIs performs an HTTP GET request to the Seaavey API with the specified endpoint and parameters.
// It automatically builds the full URL with query parameters, sends the request, and returns
// a ResponseAPIs struct with the response data.
//
// Parameters:
//   - endpoint: path after base URL, e.g. "downloader/tiktok"
//   - params: map of query parameters to encode into URL
//
// Returns:
//   - *types.ResponseAPIs: struct containing status code, body bytes, and headers
//   - error: if request creation, network call, or reading response fails
//
// Example:
//   resp, err := SeaaveyAPIs("downloader/tiktok", map[string]string{"url": "https://www.tiktok.com/..."})
//   if err != nil { ... }
//   fmt.Println(resp.Status, string(resp.Body))
func SeaaveyAPIs(endpoint string, params map[string]string) (*types.ResponseAPIs, error) {
	base := BaseAPI + "/"
	query := url.Values{}
	for k, v := range params {
		query.Set(k, v)
	}

	fullURL := base + endpoint + "?" + query.Encode()
	req, err := http.NewRequest("GET", fullURL, nil)
	if err != nil {
		return nil, err
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return &types.ResponseAPIs{
		Status:  resp.StatusCode,
		Body:    body,
		Headers: resp.Header,
	}, nil
}

// FetchBuffer performs a generic HTTP GET request to the specified URL with optional headers,
// returning the response body as a byte slice.
//
// Parameters:
//   - url: the full URL to fetch
//   - headers: optional map of HTTP headers to set on the request
//
// Returns:
//   - []byte: response body bytes
//   - error: if request creation, network call, or reading response fails
func FetchBuffer(url string, headers map[string]string) ([]byte, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	for key, value := range headers {
		req.Header.Set(key, value)
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return io.ReadAll(resp.Body)
}
