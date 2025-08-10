package types

import "net/http"

// ResponseAPIs represents a general HTTP response structure.
// It stores status code, raw response body, and headers for any API call.
type ResponseAPIs struct {
	// Status contains the HTTP status code returned by the server.
	Status int

	// Body contains the raw response body as a byte slice.
	Body []byte

	// Headers contains all HTTP headers returned in the response.
	Headers http.Header
}

// Image represents an image metadata, commonly used in responses with media.
type Image struct {
	// URL is the direct link to the image resource.
	URL string `json:"url"`

	// Width is the image width in pixels.
	Width int `json:"width"`

	// Height is the image height in pixels.
	Height int `json:"height"`
}

// TiktokResponse models the JSON response from a TikTok API call.
type TiktokResponse struct {
	// Status is the HTTP-like status code returned by the TikTok API.
	Status int `json:"status"`

	// Data contains the main payload of the TikTok response.
	Data struct {
		// Title is the video title or description.
		Title string `json:"title"`

		// Video holds video-specific data, e.g., links without watermark.
		Video *struct {
			// NoWatermark is the URL to the video without watermark.
			NoWatermark string `json:"noWatermark"`
		} `json:"video"`

		// Music contains metadata about the background music.
		Music struct {
			// PlayURL is the direct URL to the music file.
			PlayURL string `json:"play_url"`
		} `json:"music"`

		// Images is a list of related images or thumbnails.
		Images []Image `json:"images"`
	} `json:"data"`
}
