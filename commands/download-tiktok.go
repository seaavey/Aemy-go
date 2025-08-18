// Package commands implements the logic for specific bot commands.
// This file handles the 'tiktok' (and related) commands for downloading TikTok content.
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

// TiktokHandler handles the 'tiktok' command.
type TiktokHandler struct{}

// NewTiktokHandler creates a new instance of TiktokHandler.
func NewTiktokHandler() *TiktokHandler {
	return &TiktokHandler{}
}

// Handle implements the CommandHandler interface for the 'tiktok' command.
func (h *TiktokHandler) Handle(ctx context.Context, client *whatsmeow.Client, m types.Messages, evt *events.Message) error {
	url := strings.TrimSpace(m.Text)
	if url == "" {
		m.Reply("Please send a TikTok link first.")
		return nil // Return nil as this is a user input error, not a system error
	}
	if !utils.TiktokRegex.MatchString(url) {
		m.Reply("Invalid link or not a TikTok link.")
		return nil // Return nil as this is a user input error, not a system error
	}

	res, err := utils.SeaaveyAPIs("downloader/tiktok", map[string]string{"url": url})
	if err != nil || len(res.Body) == 0 {
		m.Reply("Feature error or server is down.")
		return err // Return the error to indicate a system issue
	}

	var data types.TiktokResponse
	if err := json.Unmarshal(res.Body, &data); err != nil || data.Status != 200 {
		m.Reply("Failed to get data from server.")
		return err // Return the error to indicate a system issue
	}

	if len(data.Data.Images) > 0 {
		for i, img := range data.Data.Images {
			caption := ""
			if i == 0 {
				caption = data.Data.Title
			}
			_, err := m.SendImage(img.URL, types.Options{
				Caption: caption,
			})
			if err != nil {
				m.Reply(fmt.Sprintf("Failed to send image: %v", err))
				// Continue sending other images even if one fails
			}
		}
	} else if data.Data.Video != nil && data.Data.Video.NoWatermark != "" {
		_, err := m.SendVideo(data.Data.Video.NoWatermark, types.Options{
			Caption: data.Data.Title,
		})
		if err != nil {
			m.Reply("Failed to send video.")
			return err // Return the error to indicate a system issue
		}
	} else {
		m.Reply("No media to send.")
		// Not an error, just no media found
	}

	return nil
}

// init function for automatic registration
func init() {
	handler := NewTiktokHandler()
	MustRegister([]string{"tiktok", "ttdl", "tiktokdl", "tiktokslide"}, handler, "downloader")
}