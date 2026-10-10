package bot

import (
	"context"

	"github.com/celestix/gotgproto"
	"github.com/celestix/gotgproto/parsemode"
	"github.com/p1shiA/ezLink/internal/api/handlers"
	"go.uber.org/zap"
)

type BotNotifier struct {
	client *gotgproto.Client
	logger *zap.Logger
}

func NewBotNotifier(client *gotgproto.Client, logger *zap.Logger) handlers.TelegramNotifier {
	return &BotNotifier{
		client: client,
		logger: logger.Named("BotNotifier"),
	}
}

func (n *BotNotifier) SendNotification(ctx context.Context, telegramID int64, message string) error {
	botCtx := n.client.CreateContext()
	peer, err := botCtx.ResolveInputPeerById(telegramID)
	if err != nil {
		n.logger.Warn("Failed to resolve telegram ID peer for notification", zap.Int64("telegram_id", telegramID), zap.Error(err))
		return err
	}

	_, err = botCtx.Sender.To(peer).StyledText(ctx, parsemode.StylizeText(message)...)
	if err != nil {
		n.logger.Error("Failed to send notification message to user", zap.Int64("telegram_id", telegramID), zap.Error(err))
		return err
	}

	return nil
}
