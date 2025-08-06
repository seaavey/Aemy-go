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

func HandleCommand(client *whatsmeow.Client, m types.Messages, evt *events.Message) {
	if m.Prefix == "" {
		return
	}

	switch m.Command {
	case "ping":
		jid := evt.Info.Chat
		_, err := client.SendMessage(context.Background(), jid, &waProto.Message{
			Conversation: proto.String("Pong 🏓"),
		})
		if err != nil {
			logger.Error("Gagal kirim balasan !ping: " + err.Error())
		}
	}

	
}


