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
		m.SendVideo("https://v16m-default.tiktokcdn.com/d732ed1ef06de4ab4256fb78a215c7fb/689b14e1/video/tos/alisg/tos-alisg-pve-0037c001/oEAVqRBaKBIhsi5AISlksYBEXEq3aAQGCi2va/?a=0&bti=OHYpOTY0Zik3OjlmOm01MzE6ZDQ0MDo%3D&ch=0&cr=13&dr=0&er=0&lr=all&net=0&cd=0%7C0%7C0%7C&cv=1&br=1830&bt=915&cs=0&ds=6&ft=EeF4ntZWD0Bb12NvoaNTWIxRYH0F8q_45SY&mime_type=video_mp4&qs=4&rc=aDk0Zzw6ZjRkOjc4ZWk8O0Bpamw5PHM5cnFkNTMzODczNEBiYzY0X14vNl8xM2MuMF8uYSNkXi9hMmRraC1hLS1kMTFzcw%3D%3D&vvpl=1&l=2025081212175139AE3E77EF03F30D8C31&btag=e000b8000", types.Options{})
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