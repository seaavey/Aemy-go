package handler

import (
	"botwa/utils"
	"fmt"
	"time"

	"go.mau.fi/whatsmeow"
	waTypes "go.mau.fi/whatsmeow/types"
	"go.mau.fi/whatsmeow/types/events"
)

func EventHandler(evt interface{}, client *whatsmeow.Client) {
	switch v := evt.(type) {
	case *events.Message:
		m := utils.Serialize(v)

		fmt.Println(m.From, m.Text)
		if m.FromServer == "newsletter" {
			return
		}

		if m.From.String() == "status@broadcast" {
			err := client.MarkRead(
				[]waTypes.MessageID{m.ID},
				time.Now(),
				m.From,
				*client.Store.ID,
				waTypes.ReceiptTypeRead,
			)
			if err != nil {
				fmt.Println("Gagal read status:", err)
			}
		}

		utils.HandleCommand(client, m, v)
	}
}
