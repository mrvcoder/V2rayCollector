# 🚀 GO V2rayCollector

<div dir="rtl">

## 📖 معرفی پروژه

**اولین پروژه در اینترنت** که به صورت خودکار کانفیگ‌های V2Ray را از کانال‌های تلگرام جمع‌آوری می‌کند! این پروژه هر 5 ساعت یکبار به صورت خودکار اجرا شده و لیست کانفیگ‌ها را به‌روزرسانی می‌کند.

### ✨ ویژگی‌های کلیدی

- 🔄 **به‌روزرسانی خودکار**: هر 5 ساعت یکبار با GitHub Actions
- 🌐 **پشتیبانی از پروتکل‌های متنوع**: VMess, ShadowSocks, Trojan, VLESS
- 📡 **جمع‌آوری از +170 کانال تلگرام**: پوشش گسترده منابع
- 🔍 **حذف تکراری‌ها**: فیلتر هوشمند کانفیگ‌های تکراری
- 📊 **مرتب‌سازی**: از جدیدترین به قدیمی‌ترین
- 🎯 **لینک اشتراک مستقیم**: آماده برای استفاده در کلاینت‌های V2Ray
- ⚡ **سریع و کارآمد**: نوشته شده با زبان Go
- 🔧 **قابل تنظیم**: پشتیبانی از تنظیمات مختلف برای هر کانال

## 📋 فهرست مطالب

- [نحوه استفاده](#-نحوه-استفاده)
- [نصب و راه‌اندازی](#-نصب-و-راه‌اندازی)
- [لینک‌های اشتراک](#-لینک‌های-اشتراک)
- [نحوه کار پروژه](#-نحوه-کار-پروژه)
- [تنظیمات پیشرفته](#-تنظیمات-پیشرفته)
- [مشارکت در پروژه](#-مشارکت-در-پروژه)
- [سوالات متداول](#-سوالات-متداول)
- [مشکلات رایج](#-مشکلات-رایج)

## 🎯 نحوه استفاده

### روش ساده (استفاده از لینک اشتراک)

1. یک **کلاینت V2Ray** که از لینک اشتراک پشتیبانی می‌کند را نصب کنید
2. یکی از لینک‌های زیر را در بخش اشتراک کلاینت خود اضافه کنید
3. لیست کانفیگ‌ها را به‌روزرسانی کنید
4. یک کانفیگ را انتخاب کنید و متصل شوید!

### کلاینت‌های پیشنهادی

#### 🖥️ ویندوز
- **V2RayN** - [دانلود](https://github.com/2dust/v2rayN/releases)
- **Clash for Windows** - [دانلود](https://github.com/Fndroid/clash_for_windows_pkg/releases)
- **Qv2ray** - [دانلود](https://github.com/Qv2ray/Qv2ray/releases)

#### 📱 اندروید
- **V2RayNG** - [دانلود](https://github.com/2dust/v2rayNG/releases)
- **SagerNet** - [دانلود](https://github.com/SagerNet/SagerNet/releases)
- **Clash for Android** - [دانلود](https://github.com/Kr328/ClashForAndroid/releases)

#### 🍎 iOS
- **Shadowrocket** (پولی) - [App Store](https://apps.apple.com/app/shadowrocket/id932747118)
- **V2Box** - [App Store](https://apps.apple.com/app/v2box-v2ray-client/id6446814690)
- **Stash** (پولی) - [App Store](https://apps.apple.com/app/stash/id1596063349)

#### 🐧 لینوکس
- **Qv2ray** - [دانلود](https://github.com/Qv2ray/Qv2ray/releases)
- **V2RayA** - [نصب](https://v2raya.org/en/docs/prologue/installation/debian/)
- **Clash** - [دانلود](https://github.com/Dreamacro/clash/releases)

#### 🍏 macOS
- **V2RayX** - [دانلود](https://github.com/Cenmrev/V2RayX/releases)
- **Qv2ray** - [دانلود](https://github.com/Qv2ray/Qv2ray/releases)
- **ClashX** - [دانلود](https://github.com/yichengchen/clashX/releases)

## 🔗 لینک‌های اشتراک

> ⚠️ **توجه**: GitHub Actions روی اکانت اصلی سازنده غیرفعال شده است. می‌توانید از لینک‌های زیر استفاده کنید یا این مخزن را Fork کرده و GitHub Actions را در حساب خود فعال کنید.

| نوع کانفیگ | لینک اشتراک | توضیحات |
|-----------|-------------|---------|
| **VMess** | `https://raw.githubusercontent.com/youfoundamin/V2rayCollector/main/vmess_iran.txt` | پروتکل VMess (رایج‌ترین) |
| **ShadowSocks** | `https://raw.githubusercontent.com/youfoundamin/V2rayCollector/main/ss_iran.txt` | پروتکل ShadowSocks |
| **Trojan** | `https://raw.githubusercontent.com/youfoundamin/V2rayCollector/main/trojan_iran.txt` | پروتکل Trojan |
| **VLESS** | `https://raw.githubusercontent.com/youfoundamin/V2rayCollector/main/vless_iran.txt` | پروتکل VLESS |
| **ترکیبی** | `https://raw.githubusercontent.com/youfoundamin/V2rayCollector/main/mixed_iran.txt` | ترکیب همه پروتکل‌ها |

### 📝 نحوه افزودن لینک اشتراک

#### در V2RayNG (اندروید):
1. روی آیکون `+` کلیک کنید
2. گزینه "Import subscription" را انتخاب کنید
3. لینک را وارد کنید
4. روی دکمه "OK" کلیک کنید
5. روی آیکون رفرش کلیک کنید تا کانفیگ‌ها دانلود شوند

#### در V2RayN (ویندوز):
1. منوی "Subscription" → "Subscription setting" را باز کنید
2. روی "Add" کلیک کنید
3. در قسمت "url" لینک را وارد کنید
4. روی "OK" کلیک کنید
5. منوی "Subscription" → "Update subscription" را انتخاب کنید

#### در Clash:
1. فایل کانفیگ Clash خود را باز کنید
2. در بخش `proxy-providers` لینک را اضافه کنید
3. فایل را ذخیره کنید و Clash را ریلود کنید

## 🛠️ نصب و راه‌اندازی

### پیش‌نیازها

- **Go 1.22.0** یا بالاتر
- **Git**
- اتصال به اینترنت برای دسترسی به کانال‌های تلگرام

### نصب

```bash
# کلون کردن مخزن
git clone https://github.com/mrvcoder/V2rayCollector.git
cd V2rayCollector

# نصب وابستگی‌ها
go mod download

# کامپایل و اجرا
go build -o v2ray-collector
./v2ray-collector

# یا اجرا بدون کامپایل
go run main.go
```

### استفاده با Docker (اختیاری)

```dockerfile
# Dockerfile
FROM golang:1.22-alpine AS builder
WORKDIR /app
COPY . .
RUN go mod download
RUN go build -o v2ray-collector

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/v2ray-collector .
COPY channels.csv .
CMD ["./v2ray-collector", "-sort"]
```

```bash
# ساخت و اجرا
docker build -t v2ray-collector .
docker run -v $(pwd)/output:/app v2ray-collector
```

### استفاده در GitHub Actions (Fork)

1. این مخزن را Fork کنید
2. به تنظیمات مخزن خود بروید
3. Actions را فعال کنید
4. Workflow به صورت خودکار هر 5 ساعت اجرا می‌شود
5. از لینک‌های مخزن خود استفاده کنید:
   ```
   https://raw.githubusercontent.com/YOUR_USERNAME/V2rayCollector/main/vmess_iran.txt
   ```

## ⚙️ نحوه کار پروژه

### معماری سیستم

```
┌─────────────────┐
│  channels.csv   │  ← لیست کانال‌های تلگرام
└────────┬────────┘
         │
         ▼
┌─────────────────────────────────┐
│   خواندن لیست کانال‌ها          │
└────────┬────────────────────────┘
         │
         ▼
┌─────────────────────────────────┐
│  تبدیل URL به فرمت وب تلگرام   │
└────────┬────────────────────────┘
         │
         ▼
┌─────────────────────────────────┐
│   ارسال درخواست HTTP           │
└────────┬────────────────────────┘
         │
         ▼
┌─────────────────────────────────┐
│   پارس کردن HTML با goquery    │
└────────┬────────────────────────┘
         │
         ▼
┌─────────────────────────────────┐
│   استخراج کانفیگ‌ها با Regex   │
└────────┬────────────────────────┘
         │
         ▼
┌─────────────────────────────────┐
│   پردازش ویژه VMess            │
│   (Decode → Edit → Encode)      │
└────────┬────────────────────────┘
         │
         ▼
┌─────────────────────────────────┐
│   حذف تکراری‌ها                │
└────────┬────────────────────────┘
         │
         ▼
┌─────────────────────────────────┐
│   افزودن نام و ID               │
└────────┬────────────────────────┘
         │
         ▼
┌─────────────────────────────────┐
│   مرتب‌سازی (جدید → قدیم)      │
└────────┬────────────────────────┘
         │
         ▼
┌─────────────────────────────────┐
│   ذخیره در فایل‌های خروجی       │
│   (vmess_iran.txt, ...)         │
└─────────────────────────────────┘
```

### فرآیند استخراج

1. **خواندن کانال‌ها**: از فایل `channels.csv` لیست کانال‌ها خوانده می‌شود
2. **تبدیل URL**: URLها به فرمت `web.telegram.org/k/` تبدیل می‌شوند
3. **دریافت پیام‌ها**: پیام‌های کانال با HTTP GET دریافت می‌شود
4. **پارس HTML**: با استفاده از کتابخانه goquery
5. **استخراج با Regex**: کانفیگ‌ها با الگوهای regex شناسایی می‌شوند:
   - `ss://...` → ShadowSocks
   - `vmess://...` → VMess
   - `trojan://...` → Trojan
   - `vless://...` → VLESS
6. **پردازش VMess**: رمزگشایی Base64، ویرایش JSON، رمزگذاری مجدد
7. **حذف تکراری**: با استفاده از `slices.Compact()`
8. **افزودن شناسه**: هر کانفیگ یک نام و ID منحصر به فرد می‌گیرد
9. **ذخیره‌سازی**: در فایل‌های جداگانه ذخیره می‌شوند

### پروتکل‌های پشتیبانی شده

#### 🔷 VMess
- پردازش ویژه: رمزگشایی و ویرایش فیلد `ps`
- فرمت: Base64(JSON)
- مناسب برای: اکثر کاربران

#### 🔶 ShadowSocks (SS)
- سبک‌ترین پروتکل
- فرمت: `ss://BASE64#REMARK`
- مناسب برای: سرعت بالا

#### 🔺 Trojan
- امنیت بالا با TLS
- فرمت: `trojan://PASSWORD@HOST:PORT#REMARK`
- مناسب برای: امنیت بیشتر

#### 🔻 VLESS
- نسخه سبک VMess
- فرمت: `vless://UUID@HOST:PORT?params#REMARK`
- مناسب برای: کارایی بهتر

## 🔧 تنظیمات پیشرفته

### فایل channels.csv

فرمت فایل:
```csv
URL,AllMessagesFlag
https://t.me/FreeV2rays,false
https://t.me/v2ray_configs_pool,true
```

**ستون‌ها**:
- **URL**: آدرس کانال تلگرام (با فرمت `https://t.me/CHANNEL_NAME`)
- **AllMessagesFlag**:
  - `true`: تمام متن پیام‌ها را بررسی کن (برای کانال‌های با توضیحات زیاد)
  - `false`: فقط تگ‌های `<code>` و `<pre>` را بررسی کن (سریع‌تر)

### پارامترهای خط فرمان

```bash
# مرتب‌سازی از جدید به قدیم
./v2ray-collector -sort

# حالت پیش‌فرض (قدیم به جدید)
./v2ray-collector
```

### تنظیمات داخلی (main.go)

```go
// تعداد حداکثر پیام‌ها برای هر کانال
maxMessages = 100

// نام کانفیگ‌ها
ConfigsNames = "@Vip_Security join us"
```

برای تغییر این تنظیمات، فایل `main.go` را ویرایش کنید.

### تنظیمات GitHub Actions

برای تغییر زمان‌بندی به‌روزرسانی، فایل `.github/workflows/Collector.yml` را ویرایش کنید:

```yaml
# هر 5 ساعت (پیش‌فرض)
schedule:
  - cron: "0 */5 * * *"

# هر 3 ساعت
schedule:
  - cron: "0 */3 * * *"

# هر روز ساعت 12 ظهر
schedule:
  - cron: "0 12 * * *"
```

## 📁 ساختار پروژه

```
V2rayCollector/
├── main.go                      # فایل اصلی برنامه
├── collector/
│   └── helpers.go              # توابع کمکی
├── go.mod                       # تعریف ماژول Go
├── go.sum                       # چک‌سام وابستگی‌ها
├── channels.csv                 # لیست کانال‌های تلگرام
├── channels_to_check.txt        # لیست پشتیبان
├── README.md                    # مستندات انگلیسی
├── README_FA.md                 # مستندات فارسی
├── .github/
│   └── workflows/
│       └── Collector.yml       # GitHub Actions
└── [فایل‌های خروجی]
    ├── vmess_iran.txt
    ├── ss_iran.txt
    ├── trojan_iran.txt
    ├── vless_iran.txt
    └── mixed_iran.txt
```

## 🤝 مشارکت در پروژه

ما از مشارکت شما استقبال می‌کنیم! راه‌های مختلفی برای کمک به پروژه وجود دارد:

### 1️⃣ اضافه کردن کانال جدید

اگر کانال تلگرامی می‌شناسید که کانفیگ V2Ray منتشر می‌کند:

1. فایل `channels.csv` را ویرایش کنید
2. یک خط جدید اضافه کنید:
   ```csv
   https://t.me/CHANNEL_NAME,false
   ```
3. Pull Request ارسال کنید

### 2️⃣ گزارش باگ

اگر باگی پیدا کردید:

1. به صفحه [Issues](https://github.com/mrvcoder/V2rayCollector/issues) بروید
2. یک Issue جدید ایجاد کنید
3. باگ را به تفصیل شرح دهید

### 3️⃣ پیشنهاد ویژگی جدید

ایده‌ای برای بهبود پروژه دارید؟

1. یک Issue با برچسب "enhancement" ایجاد کنید
2. ویژگی پیشنهادی را توضیح دهید
3. منتظر بازخورد باشید

### 4️⃣ بهبود کد

می‌خواهید کد را بهبود دهید؟

1. مخزن را Fork کنید
2. یک branch جدید بسازید: `git checkout -b feature/amazing-feature`
3. تغییرات را commit کنید: `git commit -m 'Add amazing feature'`
4. به branch خود push کنید: `git push origin feature/amazing-feature`
5. یک Pull Request باز کنید

### 📜 راهنمای کدنویسی

- از Go conventions استفاده کنید
- کد را document کنید
- تست‌های واحد اضافه کنید
- از gofmt برای فرمت کد استفاده کنید

## ❓ سوالات متداول

### چرا برخی کانفیگ‌ها کار نمی‌کنند؟

کانال‌های تلگرام ممکن است کانفیگ‌های منقضی شده یا نادرست منتشر کنند. پیشنهادات:
- چندین کانفیگ را امتحان کنید
- از کانفیگ‌های جدیدتر استفاده کنید (ابتدای لیست)
- پروتکل دیگری را امتحان کنید

### هر چند وقت یکبار کانفیگ‌ها به‌روز می‌شوند؟

در حالت پیش‌فرض، GitHub Actions هر **5 ساعت** یکبار اجرا می‌شود. می‌توانید این بازه را در فایل workflow تغییر دهید.

### آیا می‌توانم کانال خودم را اضافه کنم؟

بله! کافیست فایل `channels.csv` را ویرایش کرده و Pull Request ارسال کنید.

### تفاوت AllMessagesFlag چیست؟

- **true**: تمام متن پیام را بررسی می‌کند (برای کانال‌های با توضیحات و تبلیغات زیاد)
- **false**: فقط کدهای داخل تگ‌های `<code>` و `<pre>` را بررسی می‌کند (سریع‌تر و دقیق‌تر)

### چگونه یک کانفیگ خاص را پیدا کنم؟

هر کانفیگ یک ID منحصر به فرد دارد:
```
vmess://... @Vip_Security join us - 123
```
عدد آخر (123) شناسه کانفیگ است که به ترتیب افزایش می‌یابد.

### آیا این سرویس رایگان است؟

بله! این پروژه کاملاً **رایگان** و **متن‌باز** است. کانفیگ‌های جمع‌آوری شده نیز از کانال‌های عمومی تلگرام هستند.

### چرا GitHub Actions من کار نمی‌کند؟

مطمئن شوید:
1. در تنظیمات مخزن، Actions فعال است
2. مجوزهای لازم برای workflow داده شده است
3. از branch اصلی (main) استفاده می‌کنید

### می‌توانم از این پروژه برای کشورهای دیگر استفاده کنم؟

بله! فقط کافیست نام فایل‌های خروجی را تغییر دهید (مثلاً `vmess_iran.txt` به `vmess_global.txt`)

## 🔧 مشکلات رایج

### خطای "Permission Denied"

```bash
chmod +x v2ray-collector
./v2ray-collector
```

### خطای "Cannot find channels.csv"

مطمئن شوید که فایل `channels.csv` در همان پوشه اجرایی وجود دارد.

### خطای "Too Many Requests"

تلگرام ممکن است درخواست‌های زیاد را محدود کند. راه‌حل‌ها:
- تعداد کانال‌ها را کاهش دهید
- تاخیر بین درخواست‌ها اضافه کنید
- از proxy استفاده کنید

### کانفیگ‌های تکراری

اگر کانفیگ‌های تکراری مشاهده می‌کنید:
- مطمئن شوید که آخرین نسخه را دارید
- فایل‌های خروجی را پاک کرده و دوباره اجرا کنید

### خطای Build

```bash
# پاک کردن cache
go clean -modcache

# نصب مجدد وابستگی‌ها
go mod download
go mod tidy

# Build مجدد
go build -o v2ray-collector
```

## 📊 آمار پروژه

- **+170 کانال تلگرام**: پوشش گسترده
- **4 پروتکل**: VMess, SS, Trojan, VLESS
- **هر 5 ساعت**: به‌روزرسانی خودکار
- **+1000 کانفیگ**: در هر به‌روزرسانی
- **100% رایگان**: هیچ هزینه‌ای ندارد

## 🗺️ نقشه راه (Roadmap)

- [x] افزودن کامنت به توابع
- [x] ماژولار کردن دریافت پیام‌ها
- [x] فیلتر زمانی (فقط x روز اخیر)
- [x] مرتب‌سازی (از جدید به قدیم)
- [x] بهینه‌سازی استخراج کانفیگ
- [x] خواندن کانال‌ها از CSV
- [x] حذف تکراری‌ها
- [ ] پشتیبانی از کانفیگ‌های JSON
- [ ] محدود کردن تعداد کانفیگ‌ها در هر فایل
- [ ] تست خودکار کانفیگ‌ها (پینگ و سرعت)
- [ ] فیلتر بر اساس کشور سرور
- [ ] API برای دسترسی برنامه‌نویسی
- [ ] داشبورد وب برای مدیریت
- [ ] پشتیبانی از Hysteria و Tuic

## 📄 لیست کانال‌های منبع

برای مشاهده لیست کامل +170 کانال تلگرامی که به عنوان منبع استفاده می‌شوند، [اینجا کلیک کنید](https://github.com/mrvcoder/V2rayCollector/blob/main/channels.csv).

اگر کانال تلگرامی می‌شناسید که کانفیگ V2Ray منتشر می‌کند، لطفاً Pull Request ارسال کنید!

## 📞 ارتباط با ما

- **GitHub Issues**: [گزارش مشکل](https://github.com/mrvcoder/V2rayCollector/issues)
- **Pull Requests**: [مشارکت کنید](https://github.com/mrvcoder/V2rayCollector/pulls)
- **Telegram**: [@Vip_Security](https://t.me/Vip_Security)

## 📜 مجوز (License)

این پروژه تحت مجوز MIT منتشر شده است - فایل [LICENSE](LICENSE) را برای جزئیات بیشتر مشاهده کنید.

## 🌟 حمایت از پروژه

اگر این پروژه برای شما مفید بود، لطفاً:

- ⭐ به مخزن Star بدهید
- 🔄 آن را Fork کنید
- 📢 با دوستان خود به اشتراک بگذارید
- 🐛 باگ‌ها را گزارش دهید
- 💡 ایده‌های جدید پیشنهاد دهید

## ⚠️ سلب مسئولیت

این ابزار صرفاً برای اهداف آموزشی و تحقیقاتی ارائه شده است. استفاده از VPN ممکن است در برخی کشورها محدودیت داشته باشد. لطفاً قوانین محلی خود را بررسی کنید. سازندگان این پروژه هیچ مسئولیتی در قبال نحوه استفاده از این ابزار ندارند.

---

<div align="center">

**ساخته شده با ❤️ توسط جامعه متن‌باز**

اگر این پروژه را دوست دارید، یک ⭐ بدهید!

[گزارش مشکل](https://github.com/mrvcoder/V2rayCollector/issues) · [درخواست ویژگی](https://github.com/mrvcoder/V2rayCollector/issues) · [مشاهده کد](https://github.com/mrvcoder/V2rayCollector)

</div>

</div>
