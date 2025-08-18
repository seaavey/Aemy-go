package commands

import (
	"aemy/types"
	"aemy/utils"
	"context"
	"fmt"
	"os"
	"sort"
	"time"

	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/proto/waE2E"
	"go.mau.fi/whatsmeow/types/events"
	"google.golang.org/protobuf/proto"
)

type MenuHandler struct{}

func NewMenuHandler() *MenuHandler {
	return &MenuHandler{}
}

func (h *MenuHandler) Handle(ctx context.Context, client *whatsmeow.Client, m types.Messages, evt *events.Message) error {
	// Ambil waktu server sekarang
	currentTime := time.Now().Format("02-Jan-2006 15:04:05")
	hostname, _ := os.Hostname()

	uptime := time.Since(startTime).Round(time.Second)

	// Buat header info server
	txt := "*Server Info*\n"
	txt += fmt.Sprintf("• Hostname: %s\n", hostname)
	txt += fmt.Sprintf("• Time: %s\n", currentTime)
	txt += fmt.Sprintf("• Uptime: %s\n\n", uptime)

	// Get all registered commands grouped by category
	commandsByCategory := ByCategory()
	categories := make([]string, 0, len(commandsByCategory))
	for category := range commandsByCategory {
		categories = append(categories, category)
	}
	sort.Strings(categories)

	for _, category := range categories {
		commands := commandsByCategory[category]
		txt += fmt.Sprintf("*%s:*\n", utils.TitleCaser(category))

		commandNames := make([]string, 0, len(commands))
		for name := range commands {
			commandNames = append(commandNames, name)
		}
		sort.Strings(commandNames)

		for _, name := range commandNames {
			// info := commands[name]
			txt += fmt.Sprintf("  • *%s*\n", name)
		}
		txt += "\n"
	}

	thumbnail, err := os.ReadFile("config/thumbnail.png")
	if err != nil {
		thumbnail, _ = utils.FetchBuffer("https://raw.githubusercontent.com/seaavey/Aemy-go/refs/heads/main/config/thumbnail.png", nil) 
	}
	_ = m.ReplyContext(txt, &waE2E.ContextInfo{
		StanzaID:      &m.ID,
		Participant:   proto.String(m.Sender.String()),
		QuotedMessage: m.Message,
		ExternalAdReply: &waE2E.ContextInfo_ExternalAdReplyInfo{
			Title: proto.String(fmt.Sprintf("Hello, %s", utils.Ucapan())),
			Body: proto.String("Hello Everyone, I Am Seaavey Bot"),
			MediaType: (*waE2E.ContextInfo_ExternalAdReplyInfo_MediaType)(proto.Int32(1)),
			Thumbnail: thumbnail,
			SourceURL: proto.String("https://github.com/seaavey/Aemy-go"),
			RenderLargerThumbnail: proto.Bool(true),
		},
	})
	return nil
}

var startTime time.Time

func init() {
	startTime = time.Now() // buat hitung uptime
	handler := NewMenuHandler()
	MustRegister([]string{"menu", "help"}, handler, "main")
}
