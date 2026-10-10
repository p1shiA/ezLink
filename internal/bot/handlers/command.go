package handlers

import (
	"time"

	"github.com/celestix/gotgproto/ext"
	"github.com/celestix/gotgproto/parsemode"
	"github.com/p1shiA/ezLink/internal/bot/keyboards"
	"github.com/p1shiA/ezLink/internal/bot/messages"
	"github.com/p1shiA/ezLink/internal/core/domain"
	"go.uber.org/zap"
)

func (bh *BotHandler) onStartCommand(ctx *ext.Context, u *ext.Update) error {
	user := u.EffectiveUser()

	if bh.services != nil {
		// 1. Ensure user record exists / updated
		_, err := bh.services.User.GetOrRegister(ctx, user.ID, user.Username, user.FirstName, user.LastName)
		if err != nil {
			bh.logger.Error("Failed to register user", zap.Error(err))
		}

		// 2. Ensure active subscription (lazy renewal of 2GB daily free plan until Tehran midnight)
		_, err = bh.services.Subscription.GetOrCreateActiveSubscription(ctx, user.ID)
		if err != nil {
			bh.logger.Error("Failed to ensure active subscription", zap.Error(err))
		}
	}

	_, err := ctx.Reply(u, ext.ReplyTextString(messages.Welcome(user.FirstName)), &ext.ReplyOpts{
		Markup: keyboards.MainMenuKeyboard(),
	})

	bh.logger.Info("onStart Handler", zap.Int64("user_id", user.ID))
	return err
}

func (bh *BotHandler) onHelpCommand(ctx *ext.Context, u *ext.Update) error {
	_, err := ctx.Reply(u, ext.ReplyTextStyledTextArray(parsemode.StylizeText(messages.Help)), nil)
	return err
}

func (bh *BotHandler) onProfileCommand(ctx *ext.Context, u *ext.Update) error {
	user := u.EffectiveUser()
	if bh.services == nil {
		_, err := ctx.Reply(u, ext.ReplyTextString("سرویس در دسترس نیست."), nil)
		return err
	}

	// 1. Update user profile
	_, _ = bh.services.User.GetOrRegister(ctx, user.ID, user.Username, user.FirstName, user.LastName)

	// 2. Load active subscription (auto renews free if expired)
	sub, err := bh.services.Subscription.GetOrCreateActiveSubscription(ctx, user.ID)
	if err != nil {
		bh.logger.Error("Failed to get subscription in profile", zap.Error(err))
		_, err := ctx.Reply(u, ext.ReplyTextString("خطا در بارگذاری اطلاعات پروفایل."), nil)
		return err
	}

	loc, err := time.LoadLocation("Asia/Tehran")
	if err != nil {
		loc = time.FixedZone("Asia/Tehran", 3*3600+30*60)
	}

	profileMsg := bh.buildProfileText(user, sub, loc)

	_, err = ctx.Reply(u, ext.ReplyTextStyledTextArray(parsemode.StylizeText(profileMsg)), &ext.ReplyOpts{
		Markup: keyboards.AccountMenuInlineKeyboard(),
	})
	return err
}

func (bh *BotHandler) onPlansCommand(ctx *ext.Context, u *ext.Update) error {
	if bh.services == nil {
		_, err := ctx.Reply(u, ext.ReplyTextString("سرویس در دسترس نیست."), nil)
		return err
	}

	plans, err := bh.services.Plan.ListPlans(ctx)
	if err != nil || len(plans) == 0 {
		bh.logger.Error("Failed to list plans or empty", zap.Error(err))
		_, err := ctx.Reply(u, ext.ReplyTextString(messages.PlansListEmpty), nil)
		return err
	}

	text := messages.PlansListHeader
	_, err = ctx.Reply(u, ext.ReplyTextStyledTextArray(parsemode.StylizeText(text)), &ext.ReplyOpts{
		Markup: keyboards.PlansInlineKeyboard(plans, false),
	})
	return err
}

func (bh *BotHandler) onExtraCommand(ctx *ext.Context, u *ext.Update) error {
	user := u.EffectiveUser()
	if bh.services == nil {
		_, err := ctx.Reply(u, ext.ReplyTextString("سرویس در دسترس نیست."), nil)
		return err
	}

	// 1. Verify user has active paid subscription
	sub, err := bh.services.Subscription.GetOrCreateActiveSubscription(ctx, user.ID)
	if err != nil || sub.PlanType != domain.PlanTypePaid {
		_, err := ctx.Reply(u, ext.ReplyTextStyledTextArray(parsemode.StylizeText(messages.ExtraTrafficPaidOnlyWarning)), &ext.ReplyOpts{
			Markup: keyboards.BuyTrafficRedirectInlineKeyboard(),
		})
		return err
	}

	// 2. Load extra traffic packages
	packages, err := bh.services.Plan.ListExtraPackages(ctx)
	if err != nil || len(packages) == 0 {
		bh.logger.Error("Failed to list extra packages", zap.Error(err))
		_, err := ctx.Reply(u, ext.ReplyTextString(messages.ExtraPackagesListEmpty), nil)
		return err
	}

	text := messages.ExtraPackagesListHeader
	_, err = ctx.Reply(u, ext.ReplyTextStyledTextArray(parsemode.StylizeText(text)), &ext.ReplyOpts{
		Markup: keyboards.ExtraPackagesInlineKeyboard(packages),
	})
	return err
}