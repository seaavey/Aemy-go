// Package commands implements the logic for specific bot commands.
// This file handles the 'instagram' command for downloading Instagram content.
package commands

import (
	"aemy/types"
	"aemy/utils"
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/types/events"
)

// InstagramHandler handles the 'instagram' command.
type InstagramHandler struct{}

// NewInstagramHandler creates a new instance of InstagramHandler.
func NewInstagramHandler() *InstagramHandler {
	return &InstagramHandler{}
}

// Handle implements the CommandHandler interface for the 'instagram' command.
func (h *InstagramHandler) Handle(ctx context.Context, client *whatsmeow.Client, m types.Messages, evt *events.Message) error {
	url := strings.TrimSpace(m.Text)
	if url == "" {
		m.Reply("Please send an Instagram link first.")
		return nil // Return nil as this is a user input error, not a system error
	}
	if !utils.InstagramRegex.MatchString(url) {
		m.Reply("Invalid link or not an Instagram link.")
		return nil // Return nil as this is a user input error, not a system error
	}

	// Reply with waiting message
	m.Reply("Tunggu sebentar...")

	// Fetch data from Seaavey API
	res, err := utils.SeaaveyAPIs("downloader/instagram", map[string]string{"url": url})
	if err != nil || len(res.Body) == 0 {
		m.Reply("Feature error or server is down.")
		return err // Return the error to indicate a system issue
	}

	var data types.InstagramResponse
	if err := json.Unmarshal(res.Body, &data); err != nil || data.Status != 200 {
		m.Reply("Failed to get data from server.")
		return err // Return the error to indicate a system issue
	}

	if len(data.Data) == 0 {
		m.Reply("No media to send.")
		return nil
	}

	// Process each URL in the data array
	for _, mediaURL := range data.Data {
		// Determine content type
		contentType, err := utils.GetContentType(mediaURL)
		if err != nil {
			m.Reply(fmt.Sprintf("Failed to determine content type for URL: %s", mediaURL))
			continue // Continue with next URL
		}

		switch {
		case strings.HasPrefix(contentType, "video"):
			_, err := m.SendVideo(mediaURL, types.Options{})
			if err != nil {
				m.Reply(fmt.Sprintf("Failed to send video: %v", err))
				// Continue sending other media even if one fails
			}
		case strings.HasPrefix(contentType, "image"):
			_, err := m.SendImage(mediaURL, types.Options{})
			if err != nil {
				m.Reply(fmt.Sprintf("Failed to send image: %v", err))
				// Continue sending other media even if one fails
			}
		default:
			m.Reply(fmt.Sprintf("Unsupported content type: %s", contentType))
		}
	}

	return nil
}

// init function for automatic registration
func init() {
	handler := NewInstagramHandler()
	MustRegister([]string{"instagram", "igdl", "ig"}, handler, "downloader")
}