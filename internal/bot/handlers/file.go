package handlers

import (
	"fmt"
	"strconv"

	"github.com/celestix/gotgproto/ext"
	"github.com/celestix/gotgproto/parsemode"
	"github.com/gotd/td/tg"
	"github.com/p1shiA/ezLink/internal/bot"
	"github.com/p1shiA/ezLink/internal/bot/keyboards"
	"github.com/p1shiA/ezLink/internal/bot/messages"
	"go.uber.org/zap"
)

type extractedMedia struct {
	fileID       string
	fileName     string
	fileType     string
	fileSizeByte int64
}

// OnMediaMessage handles incoming file uploads (documents, videos, audio, photos).
func (bh *BotHandler) onMediaMessage(ctx *ext.Context, u *ext.Update) error {
	msg := u.EffectiveMessage
	if msg == nil || msg.Media == nil {
		return nil
	}

	media, ok := bh.extractMediaDetails(msg.Media)
	if !ok || media.fileSizeByte <= 0 {
		return nil
	}

	user := u.EffectiveUser()
	bh.logger.Info("Received media file",
		zap.Int64("user_id", user.ID),
		zap.String("file_name", media.fileName),
		zap.Int64("size", media.fileSizeByte),
		zap.String("type", media.fileType),
	)

	// 1. Check user's available traffic
	sub, err := bh.services.Subscription.GetOrCreateActiveSubscription(ctx, user.ID)
	if err != nil {
		bh.logger.Error("Failed to fetch subscription on media receive", zap.Error(err))
		_, err := ctx.Reply(u, ext.ReplyTextString("خطا در بررسی سهمیه ترافیک."), nil)
		return err
	}

	availableTraffic := sub.RemainingTrafficBytes()

	// 2. Insufficient traffic check
	if !sub.HasSufficientTraffic(media.fileSizeByte) {
		text := messages.InsufficientTraffic(media.fileSizeByte, availableTraffic)
		_, err := ctx.Reply(u, ext.ReplyTextStyledTextArray(parsemode.StylizeText(text)), &ext.ReplyOpts{
			Markup: keyboards.BuyTrafficRedirectInlineKeyboard(),
		})
		return err
	}

	// 3. Cache pending upload with short UUID (15-min TTL)
	pending := &bot.PendingFile{
		TelegramID:   user.ID,
		ChannelID:    u.EffectiveChat().GetID(),
		FileID:       media.fileID,
		FileName:     media.fileName,
		FileType:     media.fileType,
		FileSizeByte: media.fileSizeByte,
	}
	pendingID := bot.GetPendingUploadCache().Set(pending)

	// 4. Display confirmation card
	text := messages.FileConfirmCard(media.fileName, media.fileType, media.fileSizeByte, availableTraffic)
	_, err = ctx.Reply(u, ext.ReplyTextStyledTextArray(parsemode.StylizeText(text)), &ext.ReplyOpts{
		Markup: keyboards.FileConfirmInlineKeyboard(pendingID),
	})

	return err
}

func (bh *BotHandler) extractMediaDetails(rawMedia tg.MessageMediaClass) (*extractedMedia, bool) {
	switch m := rawMedia.(type) {
	case *tg.MessageMediaDocument:
		doc, ok := m.Document.(*tg.Document)
		if !ok || doc == nil {
			return nil, false
		}

		fileName := ""
		fileType := doc.MimeType
		if fileType == "" {
			fileType = "application/octet-stream"
		}

		for _, attr := range doc.Attributes {
			switch a := attr.(type) {
			case *tg.DocumentAttributeFilename:
				fileName = a.FileName
			case *tg.DocumentAttributeVideo:
				if fileName == "" {
					fileName = "video.mp4"
				}
			case *tg.DocumentAttributeAudio:
				if fileName == "" {
					if a.Title != "" {
						fileName = fmt.Sprintf("%s.mp3", a.Title)
					} else {
						fileName = "audio.mp3"
					}
				}
			}
		}

		if fileName == "" {
			fileName = fmt.Sprintf("file_%d", doc.ID)
		}

		return &extractedMedia{
			fileID:       strconv.FormatInt(doc.ID, 10),
			fileName:     fileName,
			fileType:     fileType,
			fileSizeByte: doc.Size,
		}, true

	case *tg.MessageMediaPhoto:
		photo, ok := m.Photo.(*tg.Photo)
		if !ok || photo == nil {
			return nil, false
		}

		var maxSize int64
		for _, s := range photo.Sizes {
			switch size := s.(type) {
			case *tg.PhotoSize:
				if int64(size.Size) > maxSize {
					maxSize = int64(size.Size)
				}
			case *tg.PhotoSizeProgressive:
				for _, w := range size.Sizes {
					if int64(w) > maxSize {
						maxSize = int64(w)
					}
				}
			}
		}

		if maxSize == 0 {
			maxSize = 250 * 1024 // Estimated ~250KB fallback if size not provided
		}

		return &extractedMedia{
			fileID:       strconv.FormatInt(photo.ID, 10),
			fileName:     fmt.Sprintf("photo_%d.jpg", photo.ID),
			fileType:     "image/jpeg",
			fileSizeByte: maxSize,
		}, true
	}

	return nil, false
}
