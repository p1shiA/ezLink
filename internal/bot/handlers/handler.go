package handlers

import (
	"github.com/celestix/gotgproto"
	"github.com/celestix/gotgproto/dispatcher"
	"github.com/celestix/gotgproto/dispatcher/handlers"
	"github.com/celestix/gotgproto/dispatcher/handlers/filters"
	"github.com/celestix/gotgproto/ext"
	"go.uber.org/zap"
)

type BotHandler struct {
	dp dispatcher.Dispatcher
	logger *zap.Logger
}

func NewBotHandler(client *gotgproto.Client, logger *zap.Logger) *BotHandler {
	return &BotHandler{
		dp: client.Dispatcher,
		logger: logger.Named("Bot Handler"),
	}
}

func RegisterBotHandlers(bh *BotHandler) {	
	// filters only private chats
	bh.dp.AddHandlerToGroup(handlers.NewAnyUpdate(func(ctx *ext.Context, u *ext.Update) error {
		if !u.EffectiveChat().IsAUser() {
			return dispatcher.EndGroups
		}
		return dispatcher.ContinueGroups
	}), -1)

	bh.dp.AddHandler(handlers.NewCommand("start", bh.onStartCommand))
	bh.dp.AddHandler(handlers.NewCommand("help", bh.onHelpCommand))

	bh.dp.AddHandler(handlers.NewMessage(filters.Message.Text, bh.onMenuTap))
	
}




