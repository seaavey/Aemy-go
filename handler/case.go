// Package handler provides functions for processing events received from the WhatsApp client.
// This file, case.go, specifically handles the routing of parsed commands
// to their corresponding implementation using a command handler pattern.
package handler

import (
	"context"
	"strings"

	"aemy/commands" // Import the new commands package
	"aemy/types"

	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/types/events"
)

// commandHandlers maps command names to their respective handlers.
// Each handler is initialized with its New function.
var commandHandlers = map[string]commands.CommandHandler{
	"stats":         commands.NewStatsHandler(),
	"tiktok":        commands.NewTiktokHandler(),
	"ttdl":          commands.NewTiktokHandler(),
	"tiktokdl":      commands.NewTiktokHandler(),
	"tiktokslide":   commands.NewTiktokHandler(),
	"exec":          commands.NewExecHandler(),
}

// HandleCommand processes a serialized message to execute a command.
// It checks if a command prefix was used and then routes the command
// to the appropriate handler based on a map lookup.
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

	// Lookup the handler for the command.
	if handler, ok := commandHandlers[cmd]; ok {
		// Create a context for the command execution.
		ctx := context.Background()
		
		// Execute the command handler.
		// In a more complex application, you might want to handle errors differently,
		// perhaps by logging them or sending an error message back to the user.
		// For now, we'll just ignore the error.
		_ = handler.Handle(ctx, client, m, evt)
	}
	// If no handler is found, the command is silently ignored.
	// You could add a default handler here if desired.
}