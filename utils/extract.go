package utils

import (
	"botwa/config"

	"go.mau.fi/whatsmeow/types/events"
)

func GetText(ctx *events.Message) string {
	if ctx == nil || ctx.Message == nil {
		return ""
	}

	msg := ctx.Message

	switch {
	case msg.ImageMessage != nil && msg.ImageMessage.Caption != nil:
		return *msg.ImageMessage.Caption
	case msg.VideoMessage != nil && msg.VideoMessage.Caption != nil:
		return *msg.VideoMessage.Caption
	case msg.ExtendedTextMessage != nil && msg.ExtendedTextMessage.Text != nil:
		return *msg.ExtendedTextMessage.Text
	case msg.DocumentMessage != nil && msg.DocumentMessage.Caption != nil:
		return *msg.DocumentMessage.Caption
	case msg.Conversation != nil:
		return *msg.Conversation
	default:
		return ""
	}
}

func GetPrefix(text string) string {
	if len(text) == 0 {
		return ""
	}

	first := rune(text[0])
	for _, prefix := range config.Prefixes {
		if first == prefix {
			return string(first)
		}
	}

	return ""
}
