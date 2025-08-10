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
	"strings"
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
	if m.Prefix == "" || !strings.HasPrefix(m.Body, m.Prefix) {
		return
	}

	cmd := strings.ToLower(m.Command)

	// Switch statement to handle different commands.
	switch cmd {
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
		 _ = m.Reply(infoMsg)

	case "test":
		m.SendImage("https://avatars.githubusercontent.com/u/121863865?v=4", types.Options{
		
		})
	case "tiktok", "ttdl", "tiktokdl", "tiktokslide":
	url := strings.TrimSpace(m.Text)

	switch {
		case url == "":
			m.Reply("Kirim link Tiktoknya dulu dong.")
			return
		case !utils.TiktokRegex.MatchString(url):
			m.Reply("Linknya gak valid atau bukan link Tiktok.")
			return
		}

		res, err := utils.SeaaveyAPIs("downloader/tiktok", map[string]string{"url": url})
		if err != nil || len(res.Body) == 0 {
			m.Reply("Fitur error atau server mati.")
			return
		}

		var data types.TiktokResponse
		if err := json.Unmarshal(res.Body, &data); err != nil || data.Status != 200 {
			m.Reply("Gagal ambil data dari server.")
			return
		}

	}

	
}
