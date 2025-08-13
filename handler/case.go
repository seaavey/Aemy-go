// Package utils provides helper functions and utilities for the bot.
// This file, case.go, specifically handles the routing of parsed commands
// to their corresponding implementation.
package handler

import (
	"aemy/types"
	"aemy/utils"
	"bufio"
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

// getCPUModel retrieves the CPU model name from /proc/cpuinfo.
// This is specific to Linux systems.
func getCPUModel() string {
	file, err := os.Open("/proc/cpuinfo")
	if err != nil {
		return "N/A"
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "model name") {
			parts := strings.Split(line, ":")
			if len(parts) > 1 {
				return strings.TrimSpace(parts[1])
			}
		}
	}
	return "N/A"
}


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
		cpuModel := getCPUModel()

		infoMsg := fmt.Sprintf(`*Server Info*

• Hostname: %s
• OS: %s
• Arch: %s
• Go Version: %s
• CPU: %s
• CPU Core: %d
• Uptime: %s

*Memory Usage*

• RAM Usage: %.2f MB
• Total Allocated: %.2f MB
• System Memory: %.2f MB
• Heap Allocated: %.2f MB
• Mallocs: %d
• Frees: %d

*Goroutine & GC*

• Goroutines: %d
• GC Count: %d
• Last GC: %s`,
		hostname,
		runtime.GOOS,
		runtime.GOARCH,
		runtime.Version(),
		cpuModel,
		runtime.NumCPU(),
		uptime.Truncate(time.Second).String(),
		float64(mem.Alloc)/1024/1024,
		float64(mem.TotalAlloc)/1024/1024,
		float64(mem.Sys)/1024/1024,
		float64(mem.HeapAlloc)/1024/1024,
		mem.Mallocs,
		mem.Frees,
		runtime.NumGoroutine(),
		mem.NumGC,
		time.Unix(0, int64(mem.LastGC)).Format("2006-01-02 15:04:05"),
	)
		_ = m.Reply(infoMsg)

	case "test":
		m.SendImage("https://avatars.githubusercontent.com/u/121863865?v=4", types.Options{})
		m.React("👍")
	case "tiktok", "ttdl", "tiktokdl", "tiktokslide":
		url := strings.TrimSpace(m.Text)
		if url == "" {
			m.Reply("Kirim link Tiktoknya dulu dong.")
			return
		}
		if !utils.TiktokRegex.MatchString(url) {
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
					m.Reply(fmt.Sprintf("Gagal mengirim foto: %v", err))
				}
			}
		} else if data.Data.Video != nil && data.Data.Video.NoWatermark != "" {
			_, err := m.SendVideo(data.Data.Video.NoWatermark, types.Options{
				Caption: data.Data.Title,
			})
			if err != nil {
				m.Reply("Gagal mengirim video.")
			}
		} else {
			m.Reply("Tidak ada media untuk dikirim.")
		}

	case "exec":
		if !m.IsOwner {
			return
		}

		output, err := utils.ExecuteShell(m.Text)
		if err != nil {
			m.Reply(fmt.Sprintf("Error: %v", err))
			return
		}
		_ = m.Reply(output)
	
	}
	
}