package handlers

import (
	"github.com/celestix/gotgproto"
	"github.com/celestix/gotgproto/dispatcher"
	"github.com/celestix/gotgproto/dispatcher/handlers"
	"github.com/celestix/gotgproto/dispatcher/handlers/filters"
	"github.com/celestix/gotgproto/ext"
	"github.com/p1shiA/ezLink/internal/core/domain"
	"go.uber.org/zap"
)

type Services struct {
	User         domain.UserService
	Plan         domain.PlanService
	Subscription domain.SubscriptionService
	Transaction  domain.TransactionService
	FileLog      domain.FileLogService
	Admin        domain.AdminService
}

type BotHandler struct {
	dp       dispatcher.Dispatcher
	services *Services
	logger   *zap.Logger
}

func NewBotHandler(client *gotgproto.Client, services *Services, logger *zap.Logger) *BotHandler {
	return &BotHandler{
		dp:       client.Dispatcher,
		services: services,
		logger:   logger.Named("Bot Handler"),
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
	bh.dp.AddHandler(handlers.NewCommand("profile", bh.onProfileCommand))
	bh.dp.AddHandler(handlers.NewCommand("plans", bh.onPlansCommand))
	bh.dp.AddHandler(handlers.NewCommand("extra", bh.onExtraCommand))

	// Reply menu button text taps
	bh.dp.AddHandler(handlers.NewMessage(filters.Message.Text, bh.onMenuTap))

	// Media file messages (documents, videos, audio, photos)
	bh.dp.AddHandler(handlers.NewMessage(filters.Message.Media, bh.onMediaMessage))

	// Inline keyboard callback queries
	bh.dp.AddHandler(handlers.NewCallbackQuery(filters.CallbackQuery.All, bh.onCallbackQuery))
}




