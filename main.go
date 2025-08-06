package main

import (
	"botwa/client"
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	client.Init()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop

	fmt.Println("Bot dimatikan.")
	client.WhatsAppClient.Disconnect()
}
