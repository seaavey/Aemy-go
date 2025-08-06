package client

import (
	"botwa/handler"
	"context"
	"fmt"
	"os"

	_ "github.com/mattn/go-sqlite3"
	"github.com/skip2/go-qrcode"
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/store/sqlstore"
	waLog "go.mau.fi/whatsmeow/util/log"
)

var WhatsAppClient *whatsmeow.Client

func Init() {
	log := waLog.Stdout("Client", "INFO", true)

	container, err := sqlstore.New(context.Background(), "sqlite3", "file:session.db?_foreign_keys=on", log)
	if err != nil {
		log.Errorf("DB error: %v", err)
		return
	}

	device, err := container.GetFirstDevice(context.Background())
	if err != nil {
		log.Errorf("Device error: %v", err)
		return
	}

	WhatsAppClient = whatsmeow.NewClient(device, log)
	
	WhatsAppClient.AddEventHandler(func(evt interface{}) {
		handler.EventHandler(evt, WhatsAppClient)
	})

	if WhatsAppClient.Store.ID == nil {
		qrChan, _ := WhatsAppClient.GetQRChannel(context.Background())
		err := WhatsAppClient.Connect()
		if err != nil {
			log.Errorf("Connect error: %v", err)
			return
		}

		for evt := range qrChan {
			switch evt.Event {
			case "code":
				qrcode.WriteFile(evt.Code, qrcode.Medium, 256, "qrcode.png")
				fmt.Println("QR disimpan di qrcode.png")
			case "success", "timeout":
				os.Remove("qrcode.png")
			}
		}
	} else {
		err := WhatsAppClient.Connect()
		if err != nil {
			log.Errorf("Connect error: %v", err)
			return
		}
		os.Remove("qrcode.png")
	}
}
