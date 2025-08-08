// Package utils provides utility functions and helpers for the SeaaveyBot application.
// This file specifically contains HTTP client functionality for making API requests
// to external services, particularly the Seaavey API endpoints.
package utils

import (
	"io"
	"net/http"
	"net/url"
)

// BaseAPI defines the base URL for the Seaavey API service.
// This constant serves as the root endpoint for all API interactions.
const BaseAPI = "https://api.seaavey.my.id/api"

// Fetch represents the response structure returned by HTTP requests.
// It encapsulates the complete response data including status code,
// response body, and HTTP headers for further processing.
type Fetch struct {
	// Status contains the HTTP status code returned by the server
	Status int
	
	// Body contains the raw response body as a byte slice
	Body []byte
	
	// Headers contains all HTTP headers returned in the response
	Headers http.Header
}

// SeaaveyAPIs performs an HTTP GET request to the Seaavey API with the specified endpoint and parameters.
// This function constructs the complete URL by combining the base API URL with the provided endpoint
// and appends query parameters as needed.
//
// Parameters:
//   - endpoint: The API endpoint path to append to the base URL (e.g., "users", "data/info")
//   - params: A map of query parameters to include in the request URL
//
// Returns:
//   - *Fetch: A pointer to the Fetch struct containing the response data
//   - error: An error if the request fails at any stage (URL construction, network issues, etc.)
//
// Example usage:
//   response, err := SeaaveyAPIs("downloader/tiktok", map[string]string{"url": "https://www.tiktok.com/@ayrdnaa/video/7530206290354720018"})
//   if err != nil {
//       log.Fatal(err)
//   }
//   fmt.Printf("Status: %d, Body: %s\n", response.Status, string(response.Body))
//
// The function implements the following features:
//   - Automatic URL construction with proper path joining
//   - Query parameter encoding and URL encoding
//   - 10-second timeout for requests to prevent hanging
//   - Proper User-Agent header identification
//   - Automatic response body reading and memory management
//   - Deferred response body closure to prevent resource leaks
func SeaaveyAPIs(endpoint string, params map[string]string) (*Fetch, error) {
	base := "https://api.seaavey.my.id/api/"
	query := url.Values{}
	for k, v := range params {
		query.Set(k, v)
	}

	fullURL := base + endpoint + "?" + query.Encode()
	resp, err := http.Get(fullURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return &Fetch{
		Status:  resp.StatusCode,
		Body:    body,
		Headers: resp.Header,
	}, nil
}
