package utils

import (
	"aemy/config"
	local "aemy/types"
	"bytes"
	"context"
	"fmt"
	"image"
	"image/jpeg"
	_ "image/png" // Import for decoding PNGs
	"net/http"
	"strings"
	"time"

	_ "golang.org/x/image/webp"

	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/proto/waCommon"
	"go.mau.fi/whatsmeow/proto/waE2E"
	"go.mau.fi/whatsmeow/types"
	"go.mau.fi/whatsmeow/types/events"
	"google.golang.org/protobuf/proto"
)

// Serialize converts a raw WhatsApp incoming message event (ctx) into a structured Messages object.
// It extracts command info, sender details, mentions, timestamps, and provides helper functions
// for replying, reacting, and sending images.
//
// Parameters:
//   - ctx: the WhatsApp message event received from whatsmeow
//   - client: the whatsmeow client instance to send messages or reactions
//
// Returns:
//   - local.Messages: a fully parsed and ready-to-use message struct with methods to interact with WhatsApp
//
// Details:
//   - Parses message text to find command and args based on a prefix (e.g., '!' or '.').
//   - Identifies if sender is bot owner based on configured owners.
//   - Extracts mentioned users if any in ExtendedTextMessage context.
//   - Provides Reply(text) function to send a quoted reply to the message.
//   - Provides React(emoji) function to react with an emoji.
//   - Provides SendImage(url, opts) function to download, upload, create thumbnail, and send an image message with optional caption.

func Serialize(ctx *events.Message, client *whatsmeow.Client) local.Messages {
	loc, _ := time.LoadLocation("Asia/Jakarta")
	body := GetText(ctx)
	prefix := GetPrefix(body)
	words := strings.Fields(body)
	info := ctx.Info
	cmd, args := "", []string{}

	if len(words) > 0 && strings.HasPrefix(words[0], prefix) {
		cmd = strings.TrimPrefix(words[0], prefix)
		args = words[1:]
	}

	mentionedJIDs := []types.JID{}
	if ctx.Message.ExtendedTextMessage != nil && ctx.Message.ExtendedTextMessage.ContextInfo != nil {
		for _, jid := range ctx.Message.ExtendedTextMessage.ContextInfo.MentionedJID {
			parsedJid, _ := types.ParseJID(jid)
			mentionedJIDs = append(mentionedJIDs, parsedJid)
		}
	}

	var quotedMsg *local.Messages
	if ctx.Message.ExtendedTextMessage != nil && ctx.Message.ExtendedTextMessage.ContextInfo.GetQuotedMessage() != nil {
		quotedInfo := ctx.Message.ExtendedTextMessage.ContextInfo
		quotedSenderJID, _ := types.ParseJID(quotedInfo.GetParticipant())

		quotedMsg = &local.Messages{
			ID:           quotedInfo.GetStanzaID(),
			From:         info.Chat, // The chat is the same
			Sender:       quotedSenderJID,
			SenderUser:   quotedSenderJID.User,
			SenderServer: quotedSenderJID.Server,
			Body:         GetQuotedText(quotedInfo.GetQuotedMessage()),
			Message:      quotedInfo.GetQuotedMessage(),
		}
	}

	return local.Messages{
		From:         info.Chat,
		FromUser:     info.Chat.User,
		FromServer:   info.Chat.Server,
		FromMe:       info.IsFromMe,
		ID:           info.ID,
		IsGroup:      info.IsGroup,
		IsOwner:      info.IsFromMe || isOwner(info.Sender.User),
		Sender:       info.Sender,
		SenderUser:   info.Sender.User,
		SenderServer: info.Sender.Server,
		Pushname:     info.PushName,
		Timestamp:    info.Timestamp.In(loc),
		Prefix:       prefix,
		Command:      cmd,
		Args:         args,
		Text:         strings.Join(args, " "),
		Body:         body,
		Mentioned:    mentionedJIDs,
		Message:      ctx.Message,
		Quoted:       quotedMsg,

		Reply: func(text string) error {
				_, err := client.SendMessage(context.Background(), info.Chat, &waE2E.Message{
					ExtendedTextMessage: &waE2E.ExtendedTextMessage{
						Text: proto.String(text),
						ContextInfo: &waE2E.ContextInfo{
							StanzaID:      &info.ID,
							Participant:   proto.String(info.Sender.String()),
							QuotedMessage: ctx.Message,
						},
					},
				})
				return err
			},

			ReplyContext: func(text string, contextInfo *waE2E.ContextInfo) error {
				// If no context info provided, use default quoted message context
				if contextInfo == nil {
					contextInfo = &waE2E.ContextInfo{
						StanzaID:      &info.ID,
						Participant:   proto.String(info.Sender.String()),
						QuotedMessage: ctx.Message,
					}
				}
				
				_, err := client.SendMessage(context.Background(), info.Chat, &waE2E.Message{
					ExtendedTextMessage: &waE2E.ExtendedTextMessage{
						Text:        proto.String(text),
						ContextInfo: contextInfo,
					},
				})
				return err
			},

		React: func(emoji string) error {
			_, err := client.SendMessage(context.Background(), info.Chat, &waE2E.Message{
				ReactionMessage: &waE2E.ReactionMessage{
					Key: &waCommon.MessageKey{
						FromMe:    proto.Bool(info.IsFromMe),
						ID:        proto.String(info.ID),
						Participant: proto.String(info.Sender.String()),
						RemoteJID: proto.String(info.Chat.String()),
					},
					Text: proto.String(emoji),
				},
			})
			return err
		},

		SendImage: func(url string, opts local.Options) (whatsmeow.SendResponse, error) {
			// Fetch file from URL
			data, err := FetchBuffer(url, nil)
			if err != nil {
				return whatsmeow.SendResponse{}, fmt.Errorf("fetch error: %s", err)
			}

			// Upload to WhatsApp
			uploaded, err := client.Upload(context.Background(), data, whatsmeow.MediaImage)
			if err != nil {
				return whatsmeow.SendResponse{}, fmt.Errorf("upload error: %s", err)
			}

			img, _, err := image.Decode(bytes.NewReader(data))

			if err != nil {
				return whatsmeow.SendResponse{}, fmt.Errorf("decode image error: %s", err)
			}
			
			var thumbnail bytes.Buffer
			if err := jpeg.Encode(&thumbnail, img, &jpeg.Options{Quality: 20}); err != nil {
				return whatsmeow.SendResponse{}, fmt.Errorf("encode thumbnail error: %s", err)
			}
			// Send message
			msg := &waE2E.Message{
				ImageMessage: &waE2E.ImageMessage{
					URL:           proto.String(uploaded.URL),
					DirectPath:    proto.String(uploaded.DirectPath),
					MediaKey:      uploaded.MediaKey,
					Caption:       proto.String(opts.Caption),
					Mimetype:      proto.String(http.DetectContentType(data)),
					FileEncSHA256: uploaded.FileEncSHA256,
					FileSHA256:    uploaded.FileSHA256,
					FileLength:    proto.Uint64(uint64(len(data))),
					JPEGThumbnail: thumbnail.Bytes(),
					ContextInfo:   &waE2E.ContextInfo{
						StanzaID:      &info.ID,
						Participant:   proto.String(info.Sender.String()),
						QuotedMessage: ctx.Message,	
					},
				},
			}
			
			ok, err := client.SendMessage(context.Background(), info.Chat, msg)
			if err != nil {
				return whatsmeow.SendResponse{}, fmt.Errorf("error send message %s", err)
			}
			
			return ok, nil
		},

		// BETA: Send a video
		SendVideo: func(url string, opts local.Options) (whatsmeow.SendResponse, error) {
			// Fetch file from URL
			data, err := FetchBuffer(url, nil)
			if err != nil {
				return whatsmeow.SendResponse{}, fmt.Errorf("fetch error: %s", err)
			}

			// Upload to WhatsApp
			uploaded, err := client.Upload(context.Background(), data, whatsmeow.MediaVideo)
			if err != nil {
				return whatsmeow.SendResponse{}, fmt.Errorf("upload error: %s", err)
			}

			// Generate thumbnail from first frame
			var thumbnail bytes.Buffer
			if img, _, err := image.Decode(bytes.NewReader(data)); err == nil {
				if err := jpeg.Encode(&thumbnail, img, &jpeg.Options{Quality: 20}); err != nil {
					return whatsmeow.SendResponse{}, fmt.Errorf("encode thumbnail error: %s", err)
				}
			}

			// Send message
			msg := &waE2E.Message{
				VideoMessage: &waE2E.VideoMessage{
					URL:           proto.String(uploaded.URL), 
					DirectPath:    proto.String(uploaded.DirectPath),
					MediaKey:      uploaded.MediaKey,
					Caption:       proto.String(opts.Caption),
					Mimetype:      proto.String(http.DetectContentType(data)),
					FileEncSHA256: uploaded.FileEncSHA256,
					FileSHA256:    uploaded.FileSHA256,
					FileLength:    proto.Uint64(uint64(len(data))),
					JPEGThumbnail: thumbnail.Bytes(),
					ContextInfo: &waE2E.ContextInfo{
						StanzaID:      &info.ID,
						Participant:   proto.String(info.Sender.String()),
						QuotedMessage: ctx.Message,
					},
				},
			}

			ok, err := client.SendMessage(context.Background(), info.Chat, msg)
			if err != nil {
				return whatsmeow.SendResponse{}, fmt.Errorf("error send message %s", err)
			}

			return ok, nil
		},
	}
}

// isOwner checks if a given user ID is listed as an owner in the config.
// Returns true if user is an owner, false otherwise.

func isOwner(user string) bool {
	for _, owner := range config.Owners {
		if owner == user {
			return true
		}
	}
	return false
}