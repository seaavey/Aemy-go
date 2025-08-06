// Package utils provides helper functions and utilities for the bot.
// This file, case.go, specifically handles the routing of parsed commands
// to their corresponding implementation.
package utils

import (
	"botwa/types"
	"context"

	"go.mau.fi/libsignal/logger"
	"go.mau.fi/whatsmeow"
	waProto "go.mau.fi/whatsmeow/binary/proto"
	"go.mau.fi/whatsmeow/types/events"
	"google.golang.org/protobuf/proto"
)

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
	case "ping":
		// The "ping" command sends back a "Pong" message to indicate the bot is responsive.
		jid := evt.Info.Chat
		_, err := client.SendMessage(context.Background(), jid, &waProto.Message{
			Conversation: proto.String("Pong 🏓"),
		})
		if err != nil {
			logger.Error("Failed to send ping reply: " + err.Error())
		}
	}
}
