// Package commands defines the interface for command handlers and provides
// a registry for mapping command names to their respective handlers.
package types

import (
	"context"

	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/types/events"
)

// CommandHandler defines the interface that all command handlers must implement.
// Each command handler is responsible for processing a specific command.
type CommandHandler interface {
	// Handle processes the command with the given context, client, message, and event.
	// It returns an error if the command processing fails.
	Handle(ctx context.Context, client *whatsmeow.Client, m Messages, evt *events.Message) error
}