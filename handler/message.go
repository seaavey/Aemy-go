// Package handler provides functions for processing events received from the WhatsApp client.
// It acts as the central point for routing different event types to their respective handlers.
package handler

import (
	"aemy/commands"
	"aemy/config"
	"aemy/utils"
	"context"
	"fmt"
	"strings"

	"go.mau.fi/whatsmeow"
	waTypes "go.mau.fi/whatsmeow/types"
	"go.mau.fi/whatsmeow/types/events"
)

// EventHandler is the primary event handler for the WhatsApp client.
// It receives all events from the whatsmeow client, determines their type,
// and delegates them to the appropriate logic.
//
// Parameters:
//   evt: The event interface{} received from the client. This can be any event type
//        defined by the whatsmeow library (e.g., *events.Message, *events.Receipt, etc.).
//   client: A pointer to the active whatsmeow.Client instance, used to perform actions
//           like sending messages or marking them as read.
func EventHandler(evt interface{}, client *whatsmeow.Client) {
	switch v := evt.(type) {
	// Case for handling incoming messages.
	case *events.Message:
		// Serialize the raw message event into a more manageable custom format.
		m := utils.Serialize(v, client)
		
		
		// Automatically mark status updates as read. 
		if m.From.String() == "status@broadcast" && config.ReadStatus && !m.FromMe {
			err := client.MarkRead(
				[]waTypes.MessageID{m.ID},
				m.Timestamp,
				m.From,
				m.Sender,
			)
			if err != nil {
				// Log if marking the status as read fails.
				fmt.Println("Failed to mark status as read:", err)
			}
		}
		
		
		// Ignore certain messages:
		// - Skip messages from newsletters to avoid processing channel-type messages (like WhatsApp Channels).
		// - If 'Self' mode is enabled, only allow commands from the bot owner.
		if m.FromServer == "newsletter" || (config.Self && !m.IsOwner) {
			return
		}
		
		// Pass the serialized message to the command handler for further processing.
		if m.Prefix == "" || !strings.HasPrefix(m.Body, m.Prefix) {
			return
		}

		cmd := strings.ToLower(m.Command)

		// Lookup the handler for the command from the automatic registry.
		if handler, ok := commands.Get(cmd); ok {
			// Create a context for the command execution.
			ctx := context.Background()
			
			// Execute the command handler.
			// In a more complex application, you might want to handle errors differently,
			// perhaps by logging them or sending an error message back to the user.
			// For now, we'll just ignore the error.
			_ = handler.Handle(ctx, client, m, v)
		}
		// If no handler is found, the command is silently ignored.
		// You could add a default handler here if desired.
	}
}
