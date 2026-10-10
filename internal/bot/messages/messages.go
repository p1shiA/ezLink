package messages

import (
	"fmt"
	"time"

	"github.com/p1shiA/ezLink/internal/core/domain"
)

// Main menu button labels for the persistent reply keyboard.
const (
	BtnBuySubscription = "🛍️ خرید اشتراک"
	BtnBuyExtraTraffic = "➕ خرید ترافیک اضاف"
	BtnMyProfile       = "👤 حساب کاربری من"
	BtnSupport         = "📞 پشتیبانی و راهنما"
)

// FormatBytes formats a byte count into a human-readable Persian string.
func FormatBytes(bytes int64) string {
	const (
		kb = 1024
		mb = 1024 * kb
		gb = 1024 * mb
	)

	if bytes >= gb {
		return fmt.Sprintf("%.2f گیگابایت", float64(bytes)/float64(gb))
	} else if bytes >= mb {
		return fmt.Sprintf("%.2f مگابایت", float64(bytes)/float64(mb))
	} else if bytes >= kb {
		return fmt.Sprintf("%.2f کیلوبایت", float64(bytes)/float64(kb))
	}
	return fmt.Sprintf("%d بایت", bytes)
}

// Welcome returns the greeting message displayed when a user starts the bot.
func Welcome(firstName string) string {
	return fmt.Sprintf(
		"سلام %s عزیز! 👋\n\nبرای شروع، یک فایل ارسال کن تا لینک نیم‌بهای دانلود دریافت کنی.\nهمچنین می‌تونی از منوی زیر برای مدیریت حساب یا خرید اشتراک استفاده کنی.",
		firstName,
	)
}

// Help contains the text response for the /help command and support button.
const Help = `
📖 *راهنمای استفاده از ربات*

*1️⃣ ارسال فایل:*
فایل خود را برای ربات ارسال کنید. ربات اطلاعات فایل را بررسی کرده و پس از تأیید شما، لینک دانلود تولید می‌کند.

*2️⃣ حساب کاربری:*
با انتخاب دکمه «حساب کاربری من» می‌توانید:
- ترافیک باقی‌مانده و مصرف‌شده را ببینید
- تاریخ انقضای پلن را بررسی کنید
- تاریخچه ۱۰ تراکنش و ۱۰ فایل اخیر را مشاهده کنید

*3️⃣ خرید اشتراک:*
- پلن‌های ماهانه را مشاهده کنید
- با پرداخت آنلاین امن (زرین‌پال) اشتراک خود را فعال کنید

*4️⃣ ترافیک اضافی:*
- بسته‌های ترافیک اضافی برای افزایش حجم پلن پولی فعلی
- فقط برای کاربران دارای پلن پولی فعال قابل خرید است

*⚡️ قوانین و نکات مهم:*
- هر روز یک سهمیه رایگان ۲ گیگابایتی تا بامداد روز بعد دریافت می‌کنید.
- با خرید پلن جدید، پلن قبلی منقضی شده و حجم پلن جدید جایگزین می‌شود.
- ترافیک اضافی همراه با پلن فعال منقضی می‌شود.
- حجم مصرفی اینترنت برای دانلود با لینک‌های تولیدشده، به صورت نیم‌بها محاسبه خواهد شد.
`

// Plan messages
const (
	PlansListHeader  = "🛍 *پلن‌های اشتراک موجود*\n\nبرای خرید، پلن مورد نظر خود را از دکمه‌های زیر انتخاب کنید:\n"
	PlansListEmpty   = "⚠️ در حال حاضر هیچ پلن اشتراکی برای خرید در دسترس نیست."
	PlanNotFoundAlert = "پلن انتخاب‌شده یافت نشد یا غیرفعال شده است."
)

func PlanConfirm(name string, price int64, gb int64, days int) string {
	return fmt.Sprintf(
		"🛒 *تأیید خرید اشتراک*\n\n"+
			"🔹 *نام پلن:* %s\n"+
			"⚡️ *حجم ترافیک:* %d گیگابایت\n"+
			"⏳ *مدت اعتبار:* %d روز\n"+
			"💰 *مبلغ قابل پرداخت:* %d تومان\n\n"+
			"آیا برای انتقال به درگاه پرداخت و تکمیل خرید اطمینان دارید؟",
		name, gb, days, price,
	)
}

func PaymentRedirect(itemName string, price int64) string {
	return fmt.Sprintf(
		"💳 *آماده پرداخت*\n\n"+
			"مورد: *%s*\n"+
			"مبلغ: *%d تومان*\n\n"+
			"روی دکمه پرداخت زیر کلیک کنید تا به درگاه هدایت شوید.\nپس از پرداخت، اشتراک شما بلافاصله فعال خواهد شد.",
		itemName, price,
	)
}

// Extra traffic messages
const (
	ExtraTrafficPaidOnlyWarning = "⚠️ *توجه:*\nخرید ترافیک اضافی تنها برای کاربرانی که *اشتراک ویژه (پولی)* فعال دارند امکان‌پذیر است.\n\nشما در حال حاضر از پلن رایگان استفاده می‌کنید. ابتدا یک اشتراک خریداری فرمایید."
	ExtraPackagesListHeader     = "➕ *بسته‌های ترافیک اضافی*\n\nترافیک خریداری‌شده به حجم پلن فعال شما اضافه می‌شود و تا پایان مدت اعتبار آن معتبر است.\nبسته مورد نظر را انتخاب کنید:\n"
	ExtraPackagesListEmpty      = "⚠️ در حال حاضر بسته ترافیک اضافی فعالی یافت نشد."
	ExtraPackageNotFoundAlert   = "بسته ترافیک انتخاب‌شده یافت نشد."
)

func ExtraPackageConfirm(gb int64, price int64) string {
	return fmt.Sprintf(
		"🛒 *تأیید خرید ترافیک اضافی*\n\n"+
			"➕ *حجم ترافیک:* %d گیگابایت\n"+
			"💰 *مبلغ قابل پرداخت:* %d تومان\n\n"+
			"این حجم به پلن پولی فعال شما افزوده خواهد شد.\nآیا تمایل به پرداخت دارید؟",
		gb, price,
	)
}

// Account messages
func ProfileMessage(userID int64, firstName, planTitle, remaining, used, totalLimit, extraInfo, expiryTehran string) string {
	return fmt.Sprintf(
		"👤 *حساب کاربری من*\n\n"+
			"🆔 *شناسه کاربری:* `%d`\n"+
			"👤 *نام:* %s\n"+
			"🏷 *وضعیت اشتراک:* %s\n\n"+
			"🔋 *ترافیک باقی‌مانده:* %s\n"+
			"📊 *حجم مصرف‌شده:* %s\n"+
			"📦 *سقف پلن:* %s%s\n"+
			"⏳ *اعتبار تا:* `%s`\n",
		userID,
		firstName,
		planTitle,
		remaining,
		used,
		totalLimit,
		extraInfo,
		expiryTehran,
	)
}

const (
	NoTransactionsAlert = "شما هنوز هیچ تراکنشی ثبت نکرده‌اید."
	NoFilesAlert        = "هنوز هیچ فایلی توسط شما تبدیل نشده است."
)

func TransactionHistory(txs []domain.Transaction, loc *time.Location) string {
	text := "📜 *تاریخچه ۱۰ تراکنش اخیر شما:*\n\n"
	for i, t := range txs {
		dateStr := t.CreatedAt.In(loc).Format("15:04 2006/01/02")
		statusFa := "⏳ در انتظار"
		if t.Status == domain.TxStatusCompleted {
			statusFa = "✅ موفق"
		} else if t.Status == domain.TxStatusFailed {
			statusFa = "❌ ناموفق"
		}

		typeFa := "خرید اشتراک"
		if t.TransactionType == domain.TxTypeExtraTrafficPurchase {
			typeFa = "ترافیک اضافی"
		}

		text += fmt.Sprintf(
			"%d. *%s* (%s)\n   مبلغ: %d تومان | وضعیت: %s\n   تاریخ: `%s`\n\n",
			i+1, typeFa, statusFa, t.PriceToman, statusFa, dateStr,
		)
	}
	return text
}

func FileHistory(files []domain.FileLog, loc *time.Location) string {
	text := "📁 *تاریخچه ۱۰ فایل اخیر شما:*\n\n"
	for i, f := range files {
		dateStr := f.CreatedAt.In(loc).Format("15:04 2006/01/02")
		name := f.FileName
		if name == "" {
			name = "فایل بدون نام"
		}
		text += fmt.Sprintf(
			"%d. 📄 *%s*\n   حجم: %s | فرمت: %s\n   تاریخ: `%s`\n   🔗 وضعیت لینک: تولید شده\n\n",
			i+1, name, FormatBytes(f.FileSizeByte), f.FileType, dateStr,
		)
	}
	return text
}

// File upload messages
func FileConfirmCard(fileName, fileType string, fileSize, availableTraffic int64) string {
	return fmt.Sprintf(
		"📥 *فایل دریافت شد!*\n\n"+
			"📄 *نام فایل:* %s\n"+
			"📦 *حجم فایل:* %s\n"+
			"📁 *نوع فایل:* %s\n\n"+
			"🔋 *ترافیک در دسترس شما:* %s\n\n"+
			"آیا مایل به کسر حجم و تولید لینک دانلود نیم‌بها هستید؟",
		fileName,
		FormatBytes(fileSize),
		fileType,
		FormatBytes(availableTraffic),
	)
}

func InsufficientTraffic(fileSize, availableTraffic int64) string {
	deficit := fileSize - availableTraffic
	if deficit < 0 {
		deficit = 0
	}
	return fmt.Sprintf(
		"⚠️ *ترافیک ناکافی است!*\n\n"+
			"📦 *حجم فایل:* %s\n"+
			"🔋 *ترافیک باقی‌مانده شما:* %s\n"+
			"❌ *کسری حجم:* %s\n\n"+
			"برای دانلود این فایل، لطفاً اشتراک جدید یا بسته ترافیک اضافی تهیه فرمایید.",
		FormatBytes(fileSize),
		FormatBytes(availableTraffic),
		FormatBytes(deficit),
	)
}

func FileLinkSuccess(fileName, link string, fileSize, remainingTraffic int64) string {
	return fmt.Sprintf(
		"✅ *لینک دانلود نیم‌بها با موفقیت تولید شد!*\n\n"+
			"📄 *فایل:* %s\n"+
			"📦 *حجم کسرشده:* %s\n"+
			"🔋 *ترافیک باقی‌مانده:* %s\n\n"+
			"🔗 *لینک مستقیم دانلود:*\n`%s`\n\n"+
			"⚡️ اینترنت مصرفی این لینک به صورت نیم‌بها محاسبه خواهد شد.",
		fileName,
		FormatBytes(fileSize),
		FormatBytes(remainingTraffic),
		link,
	)
}

const (
	FileCancelSuccess          = "عملیات تولید لینک لغو شد و ترافیکی کسر نگردید."
	FileExpiredAlert           = "زمان تأیید این فایل به پایان رسیده است. لطفاً فایل را مجدداً ارسال کنید."
	PaymentInitiateFailedAlert = "خطا در ایجاد درخواست پرداخت. لطفاً دوباره تلاش کنید."
)