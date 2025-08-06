package utils

import (
	"botwa/config"
	"botwa/types"
	"strings"
	"time"

	"go.mau.fi/whatsmeow/types/events"
)

func Serialize(ctx *events.Message) types.Messages {
	loc, _ := time.LoadLocation("Asia/Jakarta")

	body := GetText(ctx)
	prefix := GetPrefix(body)
	words := strings.Fields(body)

	cmd, args := "", []string{}
	if len(words) > 0 && strings.HasPrefix(words[0], prefix) {
		cmd = strings.TrimPrefix(words[0], prefix)
		args = words[1:]
	}

	return types.Messages{
		From:         ctx.Info.Chat,
		FromUser:     ctx.Info.Chat.User,
		FromServer:   ctx.Info.Chat.Server,
		FromMe:       ctx.Info.IsFromMe,
		ID:           ctx.Info.ID,
		IsGroup:      ctx.Info.IsGroup,
		IsOwner:      isOwner(ctx.Info.Chat.User),
		Sender:       ctx.Info.Sender.User + ctx.Info.Sender.Server,
		SenderUser:   ctx.Info.Sender.User,
		SenderServer: ctx.Info.Sender.Server,
		Pushname:     ctx.Info.PushName,
		Timestamp:    ctx.Info.Timestamp.In(loc).Format("2006-01-02 15:04:05"),
		Prefix:       prefix,
		Command:      cmd,
		Args:         args,
		Text:         strings.Join(args, " "),
		Body:         body,
	}
}

func isOwner(user string) bool {
	for _, owner := range config.Owners {
		if owner == user {
			return true
		}
	}
	return false
}
