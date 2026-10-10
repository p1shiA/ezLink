package keyboards

import (
	"fmt"

	"github.com/gotd/td/telegram/message/markup"
	"github.com/gotd/td/tg"
	"github.com/p1shiA/ezLink/internal/bot/messages"
	"github.com/p1shiA/ezLink/internal/core/domain"
)

// MainMenuKeyboard returns the persistent bottom reply keyboard.
func MainMenuKeyboard() tg.ReplyMarkupClass {
	return markup.BuildKeyboard().Resize().Build(
		markup.Row(markup.Button(messages.BtnBuySubscription), markup.Button(messages.BtnBuyExtraTraffic)),
		markup.Row(markup.Button(messages.BtnMyProfile)),
		markup.Row(markup.Button(messages.BtnSupport)),
	)
}

// PlansInlineKeyboard constructs the inline keyboard for choosing a subscription plan.
func PlansInlineKeyboard(plans []domain.Plan, backToAccount bool) tg.ReplyMarkupClass {
	var rows []tg.KeyboardButtonRow
	for _, p := range plans {
		gb := p.TrafficBytes / (1024 * 1024 * 1024)
		text := fmt.Sprintf("🔹 %s (%dGB - %dت)", p.Name, gb, p.PriceToman)
		data := fmt.Sprintf("plan:select:%d", p.ID)
		rows = append(rows, markup.Row(markup.Callback(text, []byte(data))))
	}

	if backToAccount {
		rows = append(rows, markup.Row(markup.Callback("🔙 بازگشت به حساب کاربری", []byte("account:view"))))
	}

	return markup.InlineKeyboard(rows...)
}

// PlanConfirmInlineKeyboard returns controls for confirming a plan purchase or going back.
func PlanConfirmInlineKeyboard(planID int) tg.ReplyMarkupClass {
	return markup.InlineKeyboard(
		markup.Row(
			markup.Callback("✅ پرداخت و تأیید", []byte(fmt.Sprintf("plan:confirm:%d", planID))),
			markup.Callback("🔙 بازگشت به لیست پلن‌ها", []byte("plan:list")),
		),
	)
}

// PaymentUrlInlineKeyboard returns a button opening the external payment gateway.
func PaymentUrlInlineKeyboard(paymentURL string) tg.ReplyMarkupClass {
	return markup.InlineKeyboard(
		markup.Row(
			markup.URL("💳 ورود به درگاه پرداخت", paymentURL),
		),
	)
}

// ExtraPackagesInlineKeyboard returns inline selection for extra traffic packages.
func ExtraPackagesInlineKeyboard(packages []domain.ExtraTrafficPackage) tg.ReplyMarkupClass {
	var rows []tg.KeyboardButtonRow
	for _, ep := range packages {
		gb := ep.TrafficBytes / (1024 * 1024 * 1024)
		text := fmt.Sprintf("➕ بسته %d گیگابایت (%d تومان)", gb, ep.PriceToman)
		data := fmt.Sprintf("extra:select:%d", ep.ID)
		rows = append(rows, markup.Row(markup.Callback(text, []byte(data))))
	}

	rows = append(rows, markup.Row(markup.Callback("🔙 بازگشت به حساب کاربری", []byte("account:view"))))

	return markup.InlineKeyboard(rows...)
}

// ExtraPackageConfirmInlineKeyboard returns controls for confirming extra package purchase.
func ExtraPackageConfirmInlineKeyboard(pkgID int) tg.ReplyMarkupClass {
	return markup.InlineKeyboard(
		markup.Row(
			markup.Callback("✅ پرداخت و تأیید", []byte(fmt.Sprintf("extra:confirm:%d", pkgID))),
			markup.Callback("🔙 بازگشت به بسته‌ها", []byte("extra:list")),
		),
	)
}

// AccountMenuInlineKeyboard returns the 3 sub-actions for the user profile.
func AccountMenuInlineKeyboard() tg.ReplyMarkupClass {
	return markup.InlineKeyboard(
		markup.Row(
			markup.Callback("🛍️ لیست پلن‌ها", []byte("account:plans")),
		),
		markup.Row(
			markup.Callback("📜 تاریخچه تراکنش‌ها", []byte("account:txs")),
			markup.Callback("📁 تاریخچه فایل‌ها", []byte("account:files")),
		),
	)
}

// BackToAccountInlineKeyboard returns a single return button to the account view.
func BackToAccountInlineKeyboard() tg.ReplyMarkupClass {
	return markup.InlineKeyboard(
		markup.Row(
			markup.Callback("🔙 بازگشت به حساب کاربری", []byte("account:view")),
		),
	)
}

// FileConfirmInlineKeyboard returns confirmation controls for file link generation.
func FileConfirmInlineKeyboard(pendingID string) tg.ReplyMarkupClass {
	return markup.InlineKeyboard(
		markup.Row(
			markup.Callback("✅ تأیید و دریافت لینک", []byte(fmt.Sprintf("file:confirm:%s", pendingID))),
			markup.Callback("❌ انصراف", []byte(fmt.Sprintf("file:cancel:%s", pendingID))),
		),
	)
}

// BuyTrafficRedirectInlineKeyboard directs a user with insufficient traffic to purchase plans.
func BuyTrafficRedirectInlineKeyboard() tg.ReplyMarkupClass {
	return markup.InlineKeyboard(
		markup.Row(
			markup.Callback("🛍️ خرید اشتراک", []byte("account:plans")),
			markup.Callback("➕ خرید ترافیک اضافی", []byte("extra:list")),
		),
	)
}