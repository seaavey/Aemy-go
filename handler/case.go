// Package utils provides helper functions and utilities for the bot.
// This file, case.go, specifically handles the routing of parsed commands
// to their corresponding implementation.
package handler

import (
	"botwa/types"
	"context"
	"fmt"
	"os"
	"runtime"
	"time"

	"go.mau.fi/libsignal/logger"
	"go.mau.fi/whatsmeow"
	waProto "go.mau.fi/whatsmeow/binary/proto"
	"go.mau.fi/whatsmeow/types/events"
	"google.golang.org/protobuf/proto"
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
	case "info":
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

		_, err := client.SendMessage(context.Background(), m.From, &waProto.Message{
			ExtendedTextMessage: &waProto.ExtendedTextMessage{
				Text: proto.String(infoMsg),
			},
		})
		if err != nil {
			logger.Error("Gagal kirim info server: " + err.Error())
		}

	}
}
