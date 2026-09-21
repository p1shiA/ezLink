package bot

import (
	"fmt"

	"github.com/celestix/gotgproto"
	"github.com/celestix/gotgproto/dispatcher"
	"github.com/celestix/gotgproto/dispatcher/handlers"
	"github.com/celestix/gotgproto/ext"
	"github.com/gotd/td/telegram/message/markup"
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

	bh.dp.AddHandler(handlers.NewCommand("start", bh.onStart))
	
}


func (bh *BotHandler) onStart (ctx *ext.Context, u *ext.Update) error {
	user := u.EffectiveUser().FirstName

	message := fmt.Sprintf("welcome %s", user)

	_, err := ctx.Reply(u, ext.ReplyTextString(message), &ext.ReplyOpts{
		Markup: markup.BuildKeyboard().Resize().Build(
			markup.Row(markup.Button("Test 1")),
			markup.Row(markup.Button("Test 2")),
		),
	})

	bh.logger.Info("onStart Handler", zap.String("User", user))
	return err
}

