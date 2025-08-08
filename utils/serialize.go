package utils

import (
	"aemy/config"
	"aemy/types"
	"context"
	"fmt"
	"io"
	"mime"
	"net/http"
	"strings"
	"time"

	"go.mau.fi/whatsmeow"
	waProto "go.mau.fi/whatsmeow/binary/proto"
	"go.mau.fi/whatsmeow/proto/waE2E"
	"go.mau.fi/whatsmeow/types/events"
	"google.golang.org/protobuf/proto"
)

// Serialize converts a raw WhatsApp message event into a structured, application-specific message format.
// This function serves as the primary data transformation layer, extracting relevant information from
// the WhatsApp protocol buffer and mapping it to our internal message structure.
//
// Parameters:
//   ctx - The raw WhatsApp message event containing all message metadata and content
//   client - The active WhatsApp client instance used for sending replies and media
//
// Returns:
//   types.Messages - A fully populated message structure with parsed command data and utility functions
func Serialize(ctx *events.Message, client *whatsmeow.Client) types.Messages {
    // Load Jakarta timezone for consistent timestamp formatting across all messages
    loc, _ := time.LoadLocation("Asia/Jakarta")
    
    // Extract the primary text content from the message, handling various message types
    body := GetText(ctx)
    
    // Determine the command prefix based on configuration settings
    prefix := GetPrefix(body)
    
    // Tokenize the message body into individual words for command parsing
    words := strings.Fields(body)
    
    // Access the core message metadata from the event
    info := ctx.Info

    // Initialize command parsing variables
    cmd, args := "", []string{}
    
    // Parse command structure if the message starts with the configured prefix
    if len(words) > 0 && strings.HasPrefix(words[0], prefix) {
        cmd = strings.TrimPrefix(words[0], prefix)
        args = words[1:]
    }

    // Construct and return the complete message structure
    return types.Messages {
        From: info.Chat,
        FromUser: info.Chat.User,
        FromServer: info.Chat.Server,
        FromMe: info.IsFromMe,
        ID: info.ID,
        IsGroup: info.IsGroup,
        IsOwner: info.IsFromMe || isOwner(info.Sender.User),
        Sender: info.Sender,
        SenderUser: info.Sender.User,
        SenderServer: info.Sender.Server,
        Pushname: info.PushName,
        Timestamp: info.Timestamp.In(loc),
        Prefix: prefix,
        Command: cmd,
        Args: args,
        Text: strings.Join(args, " "),
        Body: body,

        Reply: func(text string) error {
            _, err := client.SendMessage(context.Background(), info.Chat, &waProto.Message {
                ExtendedTextMessage: &waProto.ExtendedTextMessage {
                    Text: proto.String(text),
                    ContextInfo: &waProto.ContextInfo {
                        StanzaID: &info.ID,
                        Participant: proto.String(info.Sender.String()),
                        QuotedMessage: ctx.Message,
                    },
                },
            })
            return err
        },

        SendMedia: func(url string, caption string, options *waE2E.ContextInfo) error {
            resp, err := http.Get(url)
            if err != nil {
                return err
            }
            defer resp.Body.Close()

            data, err := io.ReadAll(resp.Body)
            if err != nil {
                return err
            }

            contentType := resp.Header.Get("Content-Type")
            if contentType == "" {
                contentType = http.DetectContentType(data)
            }

            exts, _ := mime.ExtensionsByType(contentType)
            ext := ".bin"
            if len(exts) > 0 {
                ext = exts[0]
            }
            filename := "file" + ext

            var msg * waProto.Message
            var mediaType whatsmeow.MediaType

            switch {
                case strings.HasPrefix(contentType, "image/"):
                    msg = &waProto.Message {
                        ImageMessage: &waProto.ImageMessage {
                            Caption: proto.String(caption),
                            Mimetype: proto.String(contentType),
                            JPEGThumbnail: [] byte {},
							ContextInfo: &waE2E.ContextInfo{
								 StanzaID: &info.ID,
								Participant: proto.String(info.Sender.String()),
								QuotedMessage: ctx.Message,
							},
                        },
                    }
                    mediaType = whatsmeow.MediaImage

                case strings.HasPrefix(contentType, "video/"):
                    msg = &waProto.Message {
                        VideoMessage: &waProto.VideoMessage {
                            Caption: proto.String(caption),
                            Mimetype: proto.String(contentType),
							ContextInfo: &waE2E.ContextInfo{
								StanzaID: &info.ID,
								Participant: proto.String(info.Sender.String()),
								QuotedMessage: ctx.Message,
							},
                        },
                    }
                    mediaType = whatsmeow.MediaVideo

                case strings.HasPrefix(contentType, "audio/"):
                    msg = &waProto.Message {
                        AudioMessage: &waProto.AudioMessage {
                            Mimetype: proto.String(contentType),
                            PTT: proto.Bool(false),
							ContextInfo: &waE2E.ContextInfo{
								StanzaID: &info.ID,
								Participant: proto.String(info.Sender.String()),
								QuotedMessage: ctx.Message,
							},
                        },
                    }
                    mediaType = whatsmeow.MediaAudio

                case strings.HasPrefix(contentType, "application/"):
                    msg = &waProto.Message {
                        DocumentMessage: &waProto.DocumentMessage {
                            Title: proto.String(filename),
                            Mimetype: proto.String(contentType),
                            FileName: proto.String(filename),
                            Caption: proto.String(caption),
							ContextInfo: &waE2E.ContextInfo{
								StanzaID: &info.ID,
								Participant: proto.String(info.Sender.String()),
								QuotedMessage: ctx.Message,
							},
                        },
                    }
                    mediaType = whatsmeow.MediaDocument

                default:
                    return fmt.Errorf("unsupported content type: %s", contentType)
            }

            uploaded, err := client.Upload(context.Background(), data, mediaType)
            if err != nil {
                return err
            }

            switch {
                case msg.ImageMessage != nil:
                    msg.ImageMessage.URL = &uploaded.URL
                    msg.ImageMessage.DirectPath = &uploaded.DirectPath
                    msg.ImageMessage.MediaKey = uploaded.MediaKey
                    msg.ImageMessage.FileEncSHA256 = uploaded.FileEncSHA256
                    msg.ImageMessage.FileSHA256 = uploaded.FileSHA256
                    msg.ImageMessage.FileLength = &uploaded.FileLength

                case msg.VideoMessage != nil:
                    msg.VideoMessage.URL = &uploaded.URL
                    msg.VideoMessage.DirectPath = &uploaded.DirectPath
                    msg.VideoMessage.MediaKey = uploaded.MediaKey
                    msg.VideoMessage.FileEncSHA256 = uploaded.FileEncSHA256
                    msg.VideoMessage.FileSHA256 = uploaded.FileSHA256
                    msg.VideoMessage.FileLength = &uploaded.FileLength

                case msg.AudioMessage != nil:
                    msg.AudioMessage.URL = &uploaded.URL
                    msg.AudioMessage.DirectPath = &uploaded.DirectPath
                    msg.AudioMessage.MediaKey = uploaded.MediaKey
                    msg.AudioMessage.FileEncSHA256 = uploaded.FileEncSHA256
                    msg.AudioMessage.FileSHA256 = uploaded.FileSHA256
                    msg.AudioMessage.FileLength = &uploaded.FileLength

                case msg.DocumentMessage != nil:
                    msg.DocumentMessage.URL = &uploaded.URL
                    msg.DocumentMessage.DirectPath = &uploaded.DirectPath
                    msg.DocumentMessage.MediaKey = uploaded.MediaKey
                    msg.DocumentMessage.FileEncSHA256 = uploaded.FileEncSHA256
                    msg.DocumentMessage.FileSHA256 = uploaded.FileSHA256
                    msg.DocumentMessage.FileLength = &uploaded.FileLength
            }

            // Kirim message
            _, err = client.SendMessage(context.Background(), ctx.Info.Chat, msg)
            return err
        },
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