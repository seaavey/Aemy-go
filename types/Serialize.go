package types

import "go.mau.fi/whatsmeow/types"

type Messages struct {
	From         types.JID
	FromUser     string
	FromServer   string
	FromMe       bool
	ID           types.MessageID
	IsGroup      bool
	IsOwner      bool
	Sender       string
	SenderUser   string
	SenderServer string
	Pushname 	 string
	Timestamp    string
	Prefix       string
	Command      string
	Args         []string
	Text 		 string
	Body         string
}
