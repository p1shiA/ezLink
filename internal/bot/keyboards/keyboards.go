package keyboards

import (
	"github.com/gotd/td/telegram/message/markup"
	"github.com/gotd/td/tg"
	"github.com/p1shiA/ezLink/internal/bot/messages"
)

func MainMenuKeyboard() tg.ReplyMarkupClass {
	return markup.BuildKeyboard().Resize().Build(
		markup.Row(markup.Button(messages.BtnBuySubscription), markup.Button(messages.BtnBuyExtraTraffic)),
		markup.Row(markup.Button(messages.BtnMyProfile)),
		markup.Row(markup.Button(messages.BtnSupport)),
	)
}