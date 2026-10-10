package handlers

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/celestix/gotgproto/ext"
	"github.com/celestix/gotgproto/parsemode"
	"github.com/gotd/td/telegram/message/entity"
	"github.com/gotd/td/telegram/message/styling"
	"github.com/gotd/td/tg"
	"github.com/p1shiA/ezLink/internal/bot"
	"github.com/p1shiA/ezLink/internal/bot/keyboards"
	"github.com/p1shiA/ezLink/internal/bot/messages"
	"github.com/p1shiA/ezLink/internal/core/domain"
	"go.uber.org/zap"
)

// OnCallbackQuery handles all inline button interactions.
func (bh *BotHandler) onCallbackQuery(ctx *ext.Context, u *ext.Update) error {
	cbq := u.CallbackQuery
	if cbq == nil {
		return nil
	}

	data := string(cbq.Data)
	bh.logger.Debug("Callback query received",
		zap.String("data", data),
		zap.Int64("user_id", u.EffectiveUser().ID),
	)

	parts := strings.Split(data, ":")
	action := parts[0]

	switch action {
	case "plan":
		return bh.handlePlanCallback(ctx, u, parts)
	case "extra":
		return bh.handleExtraCallback(ctx, u, parts)
	case "account":
		return bh.handleAccountCallback(ctx, u, parts)
	case "file":
		return bh.handleFileCallback(ctx, u, parts)
	default:
		bh.answerCallback(ctx, cbq.QueryID, "", false)
		return nil
	}
}

// -------------------------------------------------------------
// Subscription Plan Callbacks
// -------------------------------------------------------------

func (bh *BotHandler) handlePlanCallback(ctx *ext.Context, u *ext.Update, parts []string) error {
	queryID := u.CallbackQuery.QueryID
	if len(parts) < 2 {
		bh.answerCallback(ctx, queryID, "", false)
		return nil
	}

	subAction := parts[1]

	switch subAction {
	case "list":
		plans, err := bh.services.Plan.ListPlans(ctx)
		if err != nil || len(plans) == 0 {
			bh.answerCallback(ctx, queryID, messages.PlansListEmpty, true)
			return nil
		}
		text := messages.PlansListHeader
		markup := keyboards.PlansInlineKeyboard(plans, false)
		bh.editMessage(ctx, u, text, markup)
		bh.answerCallback(ctx, queryID, "", false)

	case "select":
		if len(parts) < 3 {
			bh.answerCallback(ctx, queryID, "", false)
			return nil
		}
		planID, err := strconv.Atoi(parts[2])
		if err != nil {
			bh.answerCallback(ctx, queryID, messages.PlanNotFoundAlert, true)
			return nil
		}

		plan, err := bh.services.Plan.GetPlanByID(ctx, planID)
		if err != nil {
			bh.answerCallback(ctx, queryID, messages.PlanNotFoundAlert, true)
			return nil
		}

		gb := plan.TrafficBytes / (1024 * 1024 * 1024)
		text := messages.PlanConfirm(plan.Name, plan.PriceToman, gb, plan.DurationDays)
		markup := keyboards.PlanConfirmInlineKeyboard(plan.ID)
		bh.editMessage(ctx, u, text, markup)
		bh.answerCallback(ctx, queryID, "", false)

	case "confirm":
		if len(parts) < 3 {
			bh.answerCallback(ctx, queryID, "", false)
			return nil
		}
		planID, err := strconv.Atoi(parts[2])
		if err != nil {
			bh.answerCallback(ctx, queryID, messages.PlanNotFoundAlert, true)
			return nil
		}

		plan, err := bh.services.Plan.GetPlanByID(ctx, planID)
		if err != nil {
			bh.answerCallback(ctx, queryID, messages.PlanNotFoundAlert, true)
			return nil
		}

		user := u.EffectiveUser()
		_, paymentURL, err := bh.services.Transaction.InitiatePlanPurchase(ctx, user.ID, plan.ID)
		if err != nil {
			bh.logger.Error("Failed to initiate plan payment", zap.Error(err))
			bh.answerCallback(ctx, queryID, messages.PaymentInitiateFailedAlert, true)
			return nil
		}

		text := messages.PaymentRedirect(plan.Name, plan.PriceToman)
		markup := keyboards.PaymentUrlInlineKeyboard(paymentURL)
		bh.editMessage(ctx, u, text, markup)
		bh.answerCallback(ctx, queryID, "", false)
	}

	return nil
}

// -------------------------------------------------------------
// Extra Traffic Callbacks
// -------------------------------------------------------------

func (bh *BotHandler) handleExtraCallback(ctx *ext.Context, u *ext.Update, parts []string) error {
	queryID := u.CallbackQuery.QueryID
	user := u.EffectiveUser()

	if len(parts) < 2 {
		bh.answerCallback(ctx, queryID, "", false)
		return nil
	}

	subAction := parts[1]

	switch subAction {
	case "list":
		// Verify active paid subscription
		sub, err := bh.services.Subscription.GetOrCreateActiveSubscription(ctx, user.ID)
		if err != nil || sub.PlanType != domain.PlanTypePaid {
			bh.answerCallback(ctx, queryID, "خرید ترافیک اضافی تنها روی پلن پولی فعال امکان‌پذیر است.", true)
			return nil
		}

		packages, err := bh.services.Plan.ListExtraPackages(ctx)
		if err != nil || len(packages) == 0 {
			bh.answerCallback(ctx, queryID, messages.ExtraPackagesListEmpty, true)
			return nil
		}

		text := messages.ExtraPackagesListHeader
		markup := keyboards.ExtraPackagesInlineKeyboard(packages)
		bh.editMessage(ctx, u, text, markup)
		bh.answerCallback(ctx, queryID, "", false)

	case "select":
		if len(parts) < 3 {
			bh.answerCallback(ctx, queryID, "", false)
			return nil
		}
		pkgID, err := strconv.Atoi(parts[2])
		if err != nil {
			bh.answerCallback(ctx, queryID, messages.ExtraPackageNotFoundAlert, true)
			return nil
		}

		pkg, err := bh.services.Plan.GetExtraPackageByID(ctx, pkgID)
		if err != nil {
			bh.answerCallback(ctx, queryID, messages.ExtraPackageNotFoundAlert, true)
			return nil
		}

		gb := pkg.TrafficBytes / (1024 * 1024 * 1024)
		text := messages.ExtraPackageConfirm(gb, pkg.PriceToman)
		markup := keyboards.ExtraPackageConfirmInlineKeyboard(pkg.ID)
		bh.editMessage(ctx, u, text, markup)
		bh.answerCallback(ctx, queryID, "", false)

	case "confirm":
		if len(parts) < 3 {
			bh.answerCallback(ctx, queryID, "", false)
			return nil
		}
		pkgID, err := strconv.Atoi(parts[2])
		if err != nil {
			bh.answerCallback(ctx, queryID, messages.ExtraPackageNotFoundAlert, true)
			return nil
		}

		pkg, err := bh.services.Plan.GetExtraPackageByID(ctx, pkgID)
		if err != nil {
			bh.answerCallback(ctx, queryID, messages.ExtraPackageNotFoundAlert, true)
			return nil
		}

		_, paymentURL, err := bh.services.Transaction.InitiateExtraTrafficPurchase(ctx, user.ID, pkg.ID)
		if err != nil {
			bh.logger.Error("Failed to initiate extra traffic payment", zap.Error(err))
			if err == domain.ErrOnlyPaidCanAddTraffic {
				bh.answerCallback(ctx, queryID, "خرید ترافیک اضافی تنها روی پلن‌های پولی امکان‌پذیر است.", true)
				return nil
			}
			bh.answerCallback(ctx, queryID, messages.PaymentInitiateFailedAlert, true)
			return nil
		}

		gb := pkg.TrafficBytes / (1024 * 1024 * 1024)
		text := messages.PaymentRedirect(fmt.Sprintf("ترافیک اضافی %d گیگابایت", gb), pkg.PriceToman)
		markup := keyboards.PaymentUrlInlineKeyboard(paymentURL)
		bh.editMessage(ctx, u, text, markup)
		bh.answerCallback(ctx, queryID, "", false)
	}

	return nil
}

// -------------------------------------------------------------
// Account Callbacks
// -------------------------------------------------------------

func (bh *BotHandler) handleAccountCallback(ctx *ext.Context, u *ext.Update, parts []string) error {
	queryID := u.CallbackQuery.QueryID
	user := u.EffectiveUser()

	if len(parts) < 2 {
		bh.answerCallback(ctx, queryID, "", false)
		return nil
	}

	subAction := parts[1]
	loc, _ := time.LoadLocation("Asia/Tehran")
	if loc == nil {
		loc = time.FixedZone("Asia/Tehran", 3*3600+30*60)
	}

	switch subAction {
	case "view":
		sub, err := bh.services.Subscription.GetOrCreateActiveSubscription(ctx, user.ID)
		if err != nil {
			bh.answerCallback(ctx, queryID, "خطا در بارگذاری اطلاعات حساب کاربری.", true)
			return nil
		}

		text := bh.buildProfileText(user, sub, loc)
		markup := keyboards.AccountMenuInlineKeyboard()
		bh.editMessage(ctx, u, text, markup)
		bh.answerCallback(ctx, queryID, "", false)

	case "plans":
		plans, err := bh.services.Plan.ListPlans(ctx)
		if err != nil || len(plans) == 0 {
			bh.answerCallback(ctx, queryID, messages.PlansListEmpty, true)
			return nil
		}

		text := messages.PlansListHeader
		markup := keyboards.PlansInlineKeyboard(plans, true) // true attaches "بازگشت به حساب کاربری"
		bh.editMessage(ctx, u, text, markup)
		bh.answerCallback(ctx, queryID, "", false)

	case "txs":
		txs, err := bh.services.Transaction.GetUserTransactions(ctx, user.ID, 10)
		if err != nil || len(txs) == 0 {
			// Spec: If no records exist, displays an alert instead of editing the message
			bh.answerCallback(ctx, queryID, messages.NoTransactionsAlert, true)
			return nil
		}

		text := messages.TransactionHistory(txs, loc)
		markup := keyboards.BackToAccountInlineKeyboard()
		bh.editMessage(ctx, u, text, markup)
		bh.answerCallback(ctx, queryID, "", false)

	case "files":
		files, err := bh.services.FileLog.GetUserFileLogs(ctx, user.ID, 10, 0)
		if err != nil || len(files) == 0 {
			// Spec: If no records exist, displays an alert instead of editing the message
			bh.answerCallback(ctx, queryID, messages.NoFilesAlert, true)
			return nil
		}

		text := messages.FileHistory(files, loc)
		markup := keyboards.BackToAccountInlineKeyboard()
		bh.editMessage(ctx, u, text, markup)
		bh.answerCallback(ctx, queryID, "", false)
	}

	return nil
}

// -------------------------------------------------------------
// File Confirmation & Cancellation Callbacks
// -------------------------------------------------------------

func (bh *BotHandler) handleFileCallback(ctx *ext.Context, u *ext.Update, parts []string) error {
	queryID := u.CallbackQuery.QueryID
	if len(parts) < 3 {
		bh.answerCallback(ctx, queryID, "", false)
		return nil
	}

	actionType := parts[1]
	pendingID := parts[2]
	cache := bot.GetPendingUploadCache()

	switch actionType {
	case "confirm":
		file, exists := cache.Get(pendingID)
		if !exists {
			bh.answerCallback(ctx, queryID, messages.FileExpiredAlert, true)
			return nil
		}

		// Download link format
		dlLink := fmt.Sprintf("https://ezlink.download/file/%s", file.FileID)

		// Deduct traffic 1:1 and save log
		_, sub, err := bh.services.FileLog.ProcessFileLink(
			ctx, file.TelegramID, file.ChannelID,
			file.FileID, file.FileName, file.FileType,
			file.FileSizeByte, dlLink,
		)
		if err != nil {
			bh.logger.Error("Failed to process file link", zap.Error(err))
			if err == domain.ErrInsufficientQuota {
				bh.answerCallback(ctx, queryID, "ترافیک شما برای دریافت این فایل کافی نیست.", true)
				return nil
			}
			bh.answerCallback(ctx, queryID, "خطا در پردازش فایل: "+err.Error(), true)
			return nil
		}

		cache.Delete(pendingID)

		text := messages.FileLinkSuccess(file.FileName, dlLink, file.FileSizeByte, sub.RemainingTrafficBytes())
		bh.editMessage(ctx, u, text, nil)
		bh.answerCallback(ctx, queryID, "لینک با موفقیت تولید شد!", false)

	case "cancel":
		cache.Delete(pendingID)
		// Spec: deletes the confirmation message and sends a cancellation callback response
		_ = ctx.DeleteMessages(u.EffectiveChat().GetID(), []int{u.CallbackQuery.MsgID})
		bh.answerCallback(ctx, queryID, messages.FileCancelSuccess, false)
	}

	return nil
}

// -------------------------------------------------------------
// Helpers for message editing and answering callbacks
// -------------------------------------------------------------

func (bh *BotHandler) editMessage(ctx *ext.Context, u *ext.Update, text string, markup tg.ReplyMarkupClass) {
	var eb entity.Builder
	_ = styling.Perform(&eb, parsemode.StylizeText(text)...)
	cleanText, entities := eb.Complete()

	req := &tg.MessagesEditMessageRequest{
		Peer:        u.EffectiveChat().GetInputPeer(),
		ID:          u.CallbackQuery.MsgID,
		Message:     cleanText,
		ReplyMarkup: markup,
	}
	if len(entities) > 0 {
		req.SetEntities(entities)
	}

	_, err := ctx.EditMessage(u.EffectiveChat().GetID(), req)
	if err != nil {
		bh.logger.Warn("Failed to edit message in callback", zap.Error(err))
	}
}

func (bh *BotHandler) answerCallback(ctx *ext.Context, queryID int64, text string, alert bool) {
	_, err := ctx.AnswerCallback(&tg.MessagesSetBotCallbackAnswerRequest{
		QueryID: queryID,
		Message: text,
		Alert:   alert,
	})
	if err != nil {
		bh.logger.Warn("Failed to answer callback", zap.Error(err))
	}
}

func (bh *BotHandler) buildProfileText(user *tg.User, sub *domain.Subscription, loc *time.Location) string {
	expiryTehran := sub.ExpiresAt.In(loc).Format("15:04 2006/01/02")

	planTitle := "رایگان روزانه"
	if sub.PlanType == domain.PlanTypePaid {
		planTitle = "اشتراک ویژه (پولی)"
	}

	remainingFormatted := messages.FormatBytes(sub.RemainingTrafficBytes())
	usedFormatted := messages.FormatBytes(sub.TrafficUsedBytes)
	totalLimitFormatted := messages.FormatBytes(sub.TrafficLimitBytes)
	extraTrafficFormatted := messages.FormatBytes(sub.ExtraTrafficLimitBytes)

	extraInfo := ""
	if sub.ExtraTrafficLimitBytes > 0 {
		extraInfo = fmt.Sprintf("\n➕ *ترافیک اضافی:* %s", extraTrafficFormatted)
	}

	return messages.ProfileMessage(
		user.ID,
		user.FirstName,
		planTitle,
		remainingFormatted,
		usedFormatted,
		totalLimitFormatted,
		extraInfo,
		expiryTehran,
	)
}
