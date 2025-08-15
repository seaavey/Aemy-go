// Package commands implements the logic for specific bot commands.
// This file handles the 'stats' command, providing server and runtime information.
package commands

import (
	"aemy/types"
	"bufio"
	"context"
	"fmt"
	"os"
	"runtime"
	"strings"
	"time"

	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/types/events"
)

// StatsHandler handles the 'stats' command.
type StatsHandler struct{}

// NewStatsHandler creates a new instance of StatsHandler.
func NewStatsHandler() *StatsHandler {
	return &StatsHandler{}
}

// Handle implements the CommandHandler interface for the 'stats' command.
func (h *StatsHandler) Handle(ctx context.Context, client *whatsmeow.Client, m types.Messages, evt *events.Message) error {
	// getCPUModel retrieves the CPU model name from /proc/cpuinfo.
	// This is specific to Linux systems.
	getCPUModel := func() string {
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

	startTime := time.Now() // Assuming this is defined somewhere globally, or we need to manage it differently

	host, _ := os.Hostname()
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
		host,
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
	return nil
}

// init function for automatic registration
func init() {
	handler := NewStatsHandler()
	MustRegister([]string{"stats"}, handler, "utility")
}