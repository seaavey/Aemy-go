// Package commands implements the logic for specific bot commands.
// This file handles the 'exec' command, allowing owners to execute shell commands.
package commands

import (
	"aemy/types"
	"aemy/utils"
	"context"
	"fmt"

	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/types/events"
)

// ExecHandler handles the 'exec' command.
type ExecHandler struct{}

// NewExecHandler creates a new instance of ExecHandler.
func NewExecHandler() *ExecHandler {
	return &ExecHandler{}
}

// Handle implements the CommandHandler interface for the 'exec' command.
func (h *ExecHandler) Handle(ctx context.Context, client *whatsmeow.Client, m types.Messages, evt *events.Message) error {
	if !m.IsOwner {
		// Silently ignore if not owner, as per original logic
		return nil
	}

	output, err := utils.ExecuteShell(m.Text)
	if err != nil {
		m.Reply(fmt.Sprintf("Error: %v", err))
		return err // Return the error to indicate a system issue
	}
	_ = m.Reply(output)
	return nil
}