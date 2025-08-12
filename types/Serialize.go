// Package types defines custom data structures used throughout the application.
// These types help standardize data formats, making message handling cleaner
// and easier when working with WhatsMeow.
package types

import (
	"time"

	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/proto/waE2E"
	"go.mau.fi/whatsmeow/types"
)

// Options defines optional parameters for sending media or messages.
type Options struct {
	// Caption is the text that will accompany the sent media.
	Caption string

	// ContextInfo contains additional message context such as mentions,
	// quoted messages, or forwarded info.
	ContextInfo *waE2E.ContextInfo
}

// Messages represents a parsed and structured WhatsApp message.
// It abstracts the raw WhatsMeow message event into a cleaner format
// with easy access to sender, chat, and content information.
type Messages struct {
	// From is the JID (Jabber ID) of the chat where the message was sent.
	// This can be either a group JID or an individual user JID.
	From types.JID

	// FromUser is the username part of the From JID (e.g., phone number without @server).
	FromUser string

	// FromServer is the server domain part of the From JID (e.g., "s.whatsapp.net").
	FromServer string

	// FromMe indicates whether the message was sent by the bot's own account.
	FromMe bool

	// ID is the unique identifier of the message.
	ID types.MessageID

	// IsGroup is true if the message came from a group chat.
	IsGroup bool

	// IsOwner is true if the message sender is listed as a bot owner in the configuration.
	IsOwner bool

	// Sender is the JID of the actual message sender.
	// In groups, this is the participant’s JID; in private chats, it’s the same as From.
	Sender types.JID

	// SenderUser is the username part of the Sender JID.
	SenderUser string

	// SenderServer is the server domain part of the Sender JID.
	SenderServer string

	// Pushname is the display name set by the sender in their WhatsApp profile.
	Pushname string

	// Timestamp is the exact time when the message was sent.
	Timestamp time.Time

	// Prefix is the bot command prefix used (e.g., "!", ".", "/").
	Prefix string

	// Command is the extracted bot command from the message text (without the prefix).
	Command string

	// Args contains any arguments following the command, split by spaces.
	Args []string

	// Text is the complete raw text content of the message.
	Text string

	// Body is an alias for Text, representing the main content of the message.
	Body string

	// Mentioned contains a list of JIDs mentioned in the message.
	Mentioned []types.JID

	// React sends a reaction (emoji) to the message.
	// Example: m.React("👍")
	React func(emoji string) error

	// Reply sends a text reply to the message.
	// Example: m.Reply("Hello!")
	Reply func(text string) error

	// SendImage sends an image to the chat.
	// url: direct URL to the image file.
	// opts: optional parameters such as Caption and ContextInfo.
	SendImage func(url string, opts Options) (whatsmeow.SendResponse, error)

	SendVideo func(url string, opts Options) (whatsmeow.SendResponse, error)
	
	// Quoted contains the serialized data of the message being replied to.
	// It is nil if the message is not a reply.
	Quoted *Messages

	// Message is the raw *waE2E.Message from whatsmeow.
	Message *waE2E.Message
}
