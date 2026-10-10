package handlers

import (
	"fmt"
	"time"

	"github.com/celestix/gotgproto/ext"
	"github.com/celestix/gotgproto/parsemode"
	"github.com/p1shiA/ezLink/internal/bot/keyboards"
	"github.com/p1shiA/ezLink/internal/bot/messages"
	"github.com/p1shiA/ezLink/internal/core/domain"
	"go.uber.org/zap"
)

func formatBytes(bytes int64) string {
	const (
		mb = 1024 * 1024
		gb = 1024 * mb
	)
	if bytes >= gb {
		return fmt.Sprintf("%.2f گیگابایت", float64(bytes)/float64(gb))
	}
	return fmt.Sprintf("%.2f مگابایت", float64(bytes)/float64(mb))
}

func (bh *BotHandler) onStartCommand(ctx *ext.Context, u *ext.Update) error {
	user := u.EffectiveUser()

	if bh.services != nil {
		// Register or update user
		_, err := bh.services.User.GetOrRegister(ctx, user.ID, user.Username, user.FirstName, user.LastName)
		if err != nil {
			bh.logger.Error("Failed to register user", zap.Error(err))
		}

		// Ensure active subscription (lazy renewal of daily free plan)
		_, err = bh.services.Subscription.GetOrCreateActiveSubscription(ctx, user.ID)
		if err != nil {
			bh.logger.Error("Failed to ensure active subscription", zap.Error(err))
		}
	}

	_, err := ctx.Reply(u, ext.ReplyTextString(messages.Welcome(user.FirstName)), &ext.ReplyOpts{
		Markup: keyboards.MainMenuKeyboard(),
	})

	bh.logger.Info("onStart Handler", zap.Int64("User", user.ID))
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

	sub, err := bh.services.Subscription.GetOrCreateActiveSubscription(ctx, user.ID)
	if err != nil {
		bh.logger.Error("Failed to get subscription in profile", zap.Error(err))
		_, err := ctx.Reply(u, ext.ReplyTextString("خطا در بارگذاری اطلاعات پروفایل."), nil)
		return err
	}

	totalFiles, totalBytes, _ := bh.services.FileLog.GetUserStats(ctx, user.ID)

	loc, err := time.LoadLocation("Asia/Tehran")
	if err != nil {
		loc = time.FixedZone("Asia/Tehran", 3*3600+30*60)
	}
	expiryTehran := sub.ExpiresAt.In(loc).Format("15:04:05 2006-01-02")

	planTitle := "رایگان روزانه"
	if sub.PlanType == domain.PlanTypePaid {
		planTitle = "اشتراک ویژه (پولی)"
	}

	remainingFormatted := formatBytes(sub.RemainingTrafficBytes())
	usedFormatted := formatBytes(sub.TrafficUsedBytes)
	totalLimitFormatted := formatBytes(sub.TrafficLimitBytes)
	extraTrafficFormatted := formatBytes(sub.ExtraTrafficLimitBytes)

	extraInfo := ""
	if sub.ExtraTrafficLimitBytes > 0 {
		extraInfo = fmt.Sprintf("\n➕ *ترافیک اضافی:* %s", extraTrafficFormatted)
	}

	profileMsg := fmt.Sprintf(
		"👤 *پروفایل کاربری*\n\n"+
			"🆔 *شناسه کاربری:* `%d`\n"+
			"🏷 *نوع پلن فعلی:* %s\n"+
			"🔋 *ترافیک باقی‌مانده:* %s\n"+
			"📊 *حجم مصرف‌شده:* %s\n"+
			"📦 *سقف پلن:* %s%s\n"+
			"⏳ *تاریخ انقضا:* `%s`\n"+
			"📁 *تعداد فایل‌های تبدیل‌شده:* %d (%s)",
		user.ID,
		planTitle,
		remainingFormatted,
		usedFormatted,
		totalLimitFormatted,
		extraInfo,
		expiryTehran,
		totalFiles,
		formatBytes(totalBytes),
	)

	_, err = ctx.Reply(u, ext.ReplyTextStyledTextArray(parsemode.StylizeText(profileMsg)), nil)
	return err
}

func (bh *BotHandler) onPlansCommand(ctx *ext.Context, u *ext.Update) error {
	if bh.services == nil {
		_, err := ctx.Reply(u, ext.ReplyTextString("سرویس در دسترس نیست."), nil)
		return err
	}

	plans, err := bh.services.Plan.ListPlans(ctx)
	if err != nil {
		bh.logger.Error("Failed to list plans", zap.Error(err))
		_, err := ctx.Reply(u, ext.ReplyTextString("خطا در بارگذاری لیست پلن‌ها."), nil)
		return err
	}

	text := "🛍 *پلن‌های اشتراک موجود*\n\nبرای خرید هر پلن می‌توانید با ارسال شناسه پلن اقدام کنید:\n\n"
	for _, p := range plans {
		gb := p.TrafficBytes / (1024 * 1024 * 1024)
		text += fmt.Sprintf("🔹 *%s*\n💰 قیمت: %d تومان\n⚡️ حجم: %d گیگابایت\n⏳ مدت اعتبار: %d روز\n\n", p.Name, p.PriceToman, gb, p.DurationDays)
	}

	_, err = ctx.Reply(u, ext.ReplyTextStyledTextArray(parsemode.StylizeText(text)), nil)
	return err
}

func (bh *BotHandler) onExtraCommand(ctx *ext.Context, u *ext.Update) error {
	user := u.EffectiveUser()
	if bh.services == nil {
		_, err := ctx.Reply(u, ext.ReplyTextString("سرویس در دسترس نیست."), nil)
		return err
	}

	// Verify user has active paid subscription
	sub, err := bh.services.Subscription.GetOrCreateActiveSubscription(ctx, user.ID)
	if err != nil || sub.PlanType != domain.PlanTypePaid {
		text := "⚠️ *توجه:*\nخرید ترافیک اضافی تنها روی پلن‌های پولی فعال امکان‌پذیر است.\nشما در حال حاضر پلن پولی فعال ندارید."
		_, err := ctx.Reply(u, ext.ReplyTextStyledTextArray(parsemode.StylizeText(text)), nil)
		return err
	}

	packages, err := bh.services.Plan.ListExtraPackages(ctx)
	if err != nil {
		bh.logger.Error("Failed to list extra packages", zap.Error(err))
		_, err := ctx.Reply(u, ext.ReplyTextString("خطا در دریافت بسته‌های ترافیک اضافی."), nil)
		return err
	}

	text := "➕ *بسته‌های ترافیک اضافی*\n\n"
	for _, ep := range packages {
		gb := ep.TrafficBytes / (1024 * 1024 * 1024)
		text += fmt.Sprintf("🔹 بسته %d گیگابایتی: %d تومان\n", gb, ep.PriceToman)
	}
	text += "\nنکته: ترافیک اضافی به حجم پلن فعلی شما اضافه می‌شود و همراه با آن منقضی می‌گردد."

	_, err = ctx.Reply(u, ext.ReplyTextStyledTextArray(parsemode.StylizeText(text)), nil)
	return err
}