package handlers

import (
	"github.com/celestix/gotgproto/ext"
	"github.com/celestix/gotgproto/parsemode"
	"github.com/p1shiA/ezLink/internal/bot/keyboards"
	"github.com/p1shiA/ezLink/internal/bot/messages"
	"go.uber.org/zap"
)

func (bh *BotHandler) onStartCommand(ctx *ext.Context, u *ext.Update) error {
	fn := u.EffectiveUser().FirstName

	_, err := ctx.Reply(u, ext.ReplyTextString(messages.Welcome(fn)), &ext.ReplyOpts{
		Markup: keyboards.MainMenuKeyboard(),
	})

	bh.logger.Info("onStart Handler", zap.Int64("User", u.EffectiveUser().ID))
	return err
}

func (bh *BotHandler) onHelpCommand(ctx *ext.Context, u *ext.Update) error {
	_, err := ctx.Reply(u, ext.ReplyTextStyledTextArray(parsemode.StylizeText(messages.Help)), nil)
	return err
}