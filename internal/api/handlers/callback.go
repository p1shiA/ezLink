package handlers

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/p1shiA/ezLink/internal/core/domain"
	"go.uber.org/zap"
)

type TelegramNotifier interface {
	SendNotification(ctx context.Context, telegramID int64, message string) error
}

type CallbackHandler struct {
	txService   domain.TransactionService
	notifier    TelegramNotifier
	botUsername string
	logger      *zap.Logger
}

func NewCallbackHandler(txService domain.TransactionService, notifier TelegramNotifier, botUsername string, logger *zap.Logger) *CallbackHandler {
	return &CallbackHandler{
		txService:   txService,
		notifier:    notifier,
		botUsername: botUsername,
		logger:      logger.Named("PaymentCallback"),
	}
}

func (h *CallbackHandler) Handle(c *gin.Context) {
	authority := c.Query("Authority")
	status := c.Query("Status")

	if authority == "" {
		c.Data(http.StatusBadRequest, "text/html; charset=utf-8", []byte(renderResultHTML(false, "خطا در پارامترهای درخواست", "شناسه تراکنش نامعتبر است", "", 0, h.botUsername)))
		return
	}

	tx, sub, err := h.txService.VerifyAndCompletePayment(c.Request.Context(), authority, status)
	if err != nil {
		h.logger.Warn("Payment verification failed", zap.String("authority", authority), zap.Error(err))

		// If transaction was found, notify user of failure
		if tx != nil && h.notifier != nil {
			_ = h.notifier.SendNotification(
				context.Background(),
				tx.TelegramID,
				"❌ تراکنش شما ناموفق بود یا توسط شما لغو شد.\nدر صورت کسر وجه از حساب، مبلغ حداکثر تا ۷۲ ساعت آینده برگشت داده می‌شود.",
			)
		}

		c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(renderResultHTML(false, "پرداخت ناموفق", "تراکنش تکمیل نشد یا لغو گردید.", authority, 0, h.botUsername)))
		return
	}

	refID := ""
	if tx.RefID != nil {
		refID = *tx.RefID
	}

	// Send success notification to user in Telegram
	if h.notifier != nil {
		msg := fmt.Sprintf(
			"✅ *پرداخت با موفقیت انجام شد!*\n\n💳 کد پیگیری: `%s`\n💰 مبلغ: %d تومان\n⚡️ اشتراک شما فعال شد.",
			refID, tx.PriceToman,
		)
		if sub != nil && sub.PlanType == domain.PlanTypePaid {
			msg += fmt.Sprintf("\n📦 حجم اشتراک: %d گیگابایت", (sub.TrafficLimitBytes+sub.ExtraTrafficLimitBytes)/(1024*1024*1024))
		}
		_ = h.notifier.SendNotification(context.Background(), tx.TelegramID, msg)
	}

	c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(renderResultHTML(true, "پرداخت موفق", "تراکنش با موفقیت انجام و اشتراک شما فعال شد.", refID, tx.PriceToman, h.botUsername)))
}

func renderResultHTML(success bool, title, message, refID string, amount int64, botUsername string) string {
	statusClass := "success"
	icon := "✓"
	statusColor := "#10b981"
	if !success {
		statusClass = "error"
		icon = "✕"
		statusColor = "#ef4444"
	}

	botLink := "https://t.me/" + botUsername
	if botUsername == "" {
		botLink = "tg://resolve"
	}

	refHtml := ""
	if refID != "" {
		refHtml = fmt.Sprintf(`<div class="item"><span class="label">کد پیگیری:</span><span class="value" dir="ltr">%s</span></div>`, refID)
	}

	amountHtml := ""
	if amount > 0 {
		amountHtml = fmt.Sprintf(`<div class="item"><span class="label">مبلغ پرداختی:</span><span class="value">%d تومان</span></div>`, amount)
	}

	return fmt.Sprintf(`<!DOCTYPE html>
<html lang="fa" dir="rtl">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>%s</title>
    <style>
        * { box-sizing: border-box; margin: 0; padding: 0; font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, Tahoma, sans-serif; }
        body {
            background: linear-gradient(135deg, #0f172a 0%%, #1e293b 100%%);
            min-height: 100vh;
            display: flex;
            align-items: center;
            justify-content: center;
            padding: 20px;
            color: #f8fafc;
        }
        .card {
            background: rgba(30, 41, 59, 0.7);
            backdrop-filter: blur(16px);
            border: 1px solid rgba(255, 255, 255, 0.1);
            border-radius: 24px;
            padding: 40px 32px;
            width: 100%%;
            max-width: 440px;
            text-align: center;
            box-shadow: 0 25px 50px -12px rgba(0, 0, 0, 0.5);
        }
        .icon-wrapper {
            width: 80px;
            height: 80px;
            margin: 0 auto 24px;
            border-radius: 50%%;
            display: flex;
            align-items: center;
            justify-content: center;
            font-size: 38px;
            font-weight: bold;
            background: rgba(%s, 0.15);
            color: %s;
            border: 2px solid %s;
        }
        h1 { font-size: 24px; margin-bottom: 12px; font-weight: 700; }
        p { color: #94a3b8; font-size: 15px; line-height: 1.6; margin-bottom: 24px; }
        .details {
            background: rgba(15, 23, 42, 0.6);
            border-radius: 14px;
            padding: 16px;
            margin-bottom: 28px;
            border: 1px solid rgba(255, 255, 255, 0.05);
        }
        .item {
            display: flex;
            justify-content: space-between;
            padding: 8px 0;
            font-size: 14px;
            border-bottom: 1px solid rgba(255, 255, 255, 0.05);
        }
        .item:last-child { border-bottom: none; }
        .label { color: #94a3b8; }
        .value { color: #f1f5f9; font-weight: 600; }
        .btn {
            display: inline-block;
            width: 100%%;
            padding: 14px;
            background: #3b82f6;
            color: white;
            text-decoration: none;
            border-radius: 12px;
            font-size: 16px;
            font-weight: 600;
            transition: all 0.2s ease;
            box-shadow: 0 4px 12px rgba(59, 130, 246, 0.3);
        }
        .btn:hover { background: #2563eb; transform: translateY(-1px); }
    </style>
</head>
<body>
    <div class="card">
        <div class="icon-wrapper %s">%s</div>
        <h1>%s</h1>
        <p>%s</p>
        <div class="details">
            %s
            %s
        </div>
        <a href="%s" class="btn">بازگشت به ربات تلگرام</a>
    </div>
</body>
</html>`,
		title,
		statusColor, statusColor, statusColor,
		statusClass, icon,
		title,
		message,
		refHtml,
		amountHtml,
		botLink,
	)
}
