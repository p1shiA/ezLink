package handlers

import (
	"github.com/celestix/gotgproto/dispatcher"
	"github.com/celestix/gotgproto/ext"
	"github.com/p1shiA/ezLink/internal/bot/messages"
)

func (bh *BotHandler) onMenuTap(ctx *ext.Context, u *ext.Update) error {
	text := u.EffectiveMessage.Text

	switch text {
	case messages.BtnSupport:
		return bh.onHelpCommand(ctx, u)
	case messages.BtnMyProfile:
		return bh.onProfileCommand(ctx, u)
	case messages.BtnBuySubscription:
		return bh.onPlansCommand(ctx, u)
	case messages.BtnBuyExtraTraffic:
		return bh.onExtraCommand(ctx, u)
	default:
		return dispatcher.ContinueGroups
	}
}