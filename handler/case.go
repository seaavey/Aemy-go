// Package utils provides helper functions and utilities for the bot.
// This file, case.go, specifically handles the routing of parsed commands
// to their corresponding implementation.
package handler

import (
	"aemy/types"
	"aemy/utils"
	"encoding/json"
	"fmt"
	"os"
	"runtime"
	"time"

	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/types/events"
)

var startTime = time.Now()


// HandleCommand processes a serialized message to execute a command.
// It checks if a command prefix was used and then uses a switch statement
// to route the command to the appropriate function.
//
// Parameters:
//   client: A pointer to the active whatsmeow.Client instance, used for sending replies.
//   m: The serialized message object (types.Messages) containing parsed command info.
//   evt: The original raw message event (*events.Message) from whatsmeow.
func HandleCommand(client *whatsmeow.Client, m types.Messages, evt *events.Message) {
	// If no prefix is detected, it's not a command, so we do nothing.
	if m.Prefix == "" {
		return
	}

	// Switch statement to handle different commands.
	switch m.Command {
		// Handle the "ping" command.
	case "stats":
		hostname, _ := os.Hostname()
		uptime := time.Since(startTime) 
		var mem runtime.MemStats
		runtime.ReadMemStats(&mem)

		infoMsg := fmt.Sprintf(
			"*Server Info*\n\n"+
				"• Hostname: %s\n"+
				"• OS: %s/%s\n"+
				"• Arch: %s\n"+
				"• Uptime: %s\n"+
				"• RAM Usage: %.2f MB\n"+
				"• Goroutines: %d",
			hostname,
			runtime.GOOS, runtime.GOARCH,
			runtime.GOARCH,
			uptime.Truncate(time.Second).String(),
			float64(mem.Alloc)/1024/1024,
			runtime.NumGoroutine(),
		)
		m.Reply(infoMsg)

	case "tiktok", "ttdl", "tiktokdl", "tiktokslide":
	url := m.Text

	if url == "" {
		_ = m.Reply("Kirim link Tiktoknya dulu dong.")
		return
	}

	if !utils.TiktokRegex.MatchString(url) {
		_ = m.Reply("Linknya gak valid atau bukan link Tiktok.")
		return
	}

	// Debug log
	fmt.Printf("Processing TikTok URL: %s\n", url)

	res, err := utils.SeaaveyAPIs("downloader/tiktok", map[string]string{
		"url": url,
	})

	if err != nil {
		fmt.Printf("API Error: %v\n", err)
		_ = m.Reply("Fitur error atau server mati.")
		return
	}

	// Debug response
	fmt.Printf("API Response: %s\n", string(res.Body))

	// Validasi response body kosong
	if len(res.Body) == 0 {
		_ = m.Reply("Server mengembalikan response kosong.")
		return
	}

	var tiktokResp types.TiktokResponse
	if err := json.Unmarshal(res.Body, &tiktokResp); err != nil {
		fmt.Printf("JSON Parse Error: %v\n", err)
		_ = m.Reply(fmt.Sprintf("Gagal parsing data dari server: %v", err))
		return
	}

	if tiktokResp.Status != 200 {
		fmt.Printf("API Status Error: %d\n", tiktokResp.Status)
		_ = m.Reply(fmt.Sprintf("Server error: status %d", tiktokResp.Status))
		return
	}

	// Validasi data response
	if tiktokResp.Data.Video.NoWatermark == "" {
		_ = m.Reply("Video tidak tersedia atau link tidak valid.")
		return
	}

	videoURL := tiktokResp.Data.Video.NoWatermark
	caption := tiktokResp.Data.Title

	if caption == "" {
		caption = "Video TikTok"
	}

		m.SendMedia(videoURL, caption, nil)
	}


}
