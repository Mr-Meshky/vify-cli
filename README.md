# ⚡ Vify

<p align="center">
  <img src="https://raw.githubusercontent.com/Mr-Meshky/vify-cli/main/ui/macos/Runner/Assets.xcassets/AppIcon.appiconset/app_icon_128.png" alt="Vify Logo" width="96" height="96" />
</p>

<p align="center">
  <b>فیلترشکن فوق‌سریع، هوشمند و کراس‌پلتفرم Vify برای دور زدن فیلترینگ شدید در ایران</b><br>
  <b>Blazing-Fast, Intelligent Anti-Censorship VPN & Proxy Client (Desktop GUI, Mobile & CLI)</b>
</p>

<p align="center">
  <a href="https://github.com/Mr-Meshky/vify-cli/releases"><img src="https://img.shields.io/github/v/release/Mr-Meshky/vify-cli?color=00F5D4&label=Release&logo=github&style=flat-square" alt="Release"></a>
  <img src="https://img.shields.io/badge/Platform-macOS%20%7C%20Linux%20%7C%20Windows%20%7C%20Android-05FFA1?style=flat-square" alt="Platforms">
  <img src="https://img.shields.io/badge/Language-Go%20%2B%20Flutter-00BBF9?style=flat-square" alt="Stack">
  <a href="LICENSE"><img src="https://img.shields.io/badge/License-MIT-9B5DE5?style=flat-square" alt="License"></a>
</p>

<p align="center">
  <a href="#-راهنمای-فارسی">🇮🇷 راهنمای فارسی</a> •
  <a href="#-دانلود-و-نصب--downloads">📥 دانلود / Downloads</a> •
  <a href="#-راهنمای-کامل-نسخه-اندروید-android-guide">📱 راهنمای اندروید</a> •
  <a href="#-پرسش‌های-متداول-و-عیب‌یابی-faq">❓ پرسش‌های متداول</a> •
  <a href="#english">🇬🇧 English Guide</a>
</p>

---

<a id="-دانلود-و-نصب--downloads"></a>
## 📥 دانلود و نصب / Downloads

تمام نسخه‌ها از صفحه **[GitHub Releases](https://github.com/Mr-Meshky/vify-cli/releases)** قابل دریافت هستند:

| سیستم‌عامل / پلتفرم | نوع برنامه | فرمت فایل | راهنمای اجرا |
| :--- | :--- | :--- | :--- |
| 🤖 **Android (Universal)** | 📱 اپلیکیشن موبایل | `.apk` (Universal) | نصب فایل APK روی تمام گوشی‌های اندرویدی (Android 7+) |
| 🤖 **Android (ARM64)** | 📱 اپلیکیشن بهینه | `.apk` (arm64-v8a) | نسخه کم‌حجم‌تر مخصوص پردازنده‌های ۶۴ بیتی مدرن |
| 🍏 **macOS** (Universal) | 🖥️ Desktop GUI | `.dmg` / `.zip` | نصب فایل DMG یا انتقال `Vify.app` به Applications (در صورت خطای damaged: `xattr -cr /Applications/Vify.app`) |
| 🪟 **Windows** (x64) | 🖥️ Desktop GUI | Setup `.exe` / `.zip` | نصب خودکار با اینستالر یا اجرای پرتابل `Vify.exe` |
| 🐧 **Linux** (x86_64) | 🖥️ Desktop GUI | `.deb` / `.tar.gz` | نصب با `sudo dpkg -i` در اوبونتو/دبیان یا اجرای باینری `vify-ui` |
| 💻 **Linux CLI** (x86_64) | ⚡ ترمینال CLI | باینری مستقل | `chmod +x vify-cli-linux-amd64 && ./vify-cli-linux-amd64` |
| 💻 **macOS CLI** | ⚡ ترمینال CLI | باینری مستقل | `chmod +x vify-cli-darwin-universal && ./vify-cli-darwin-universal` |

---

<a id="-راهنمای-فارسی"></a>
## 🇮🇷 راهنمای فارسی

> **کلاینت ضدفیلترینگ مدرن، فوق‌سریع و کراس‌پلتفرم — با هسته پرقدرت Go/Xray و رابط کاربری نئونی Flutter**  
> دسترسی آسان و پایدار به پراکسی‌های پرسرعت، سالم و ضدفیلتر از [اکوسیستم Vify](https://github.com/Mr-Meshky/vify) برای دور زدن فیلترینگ شدید همراه اول، ایرانسل، رایتل، مخابرات، شاتل و سایر اپراتورها.

---

### 🖥️ ویژگی‌های برجسته اپلیکیشن گرافیکی (Desktop & Mobile)

- ⚡ **اتصال تک‌کلیکه هوشمند (One-Tap Fast-Pass):** تست موازی ده‌ها کانفیگ و اتصال خودکار به سریع‌ترین سرور سالم زیر ۱ ثانیه.
- 🎯 **تست ترافیک واقعی اینترنت (Real HTTP 204 Verification):** حذف خودکار کانفیگ‌های فیک کلادفلر که پینگ TCP دارند اما دیتایی عبور نمی‌دهند؛ تضمین برقراری تبادل دانلود و آپلود واقعی.
- 🎨 **طراحی نئونی شیشه‌ای (Modern Glassmorphism):** رابط کاربری تیره، چشم‌نواز با رنگ‌های نئونی، رینگ وضعیت پویا و جلوه‌های بصری خیره‌کننده.
- 📊 **مانیتورینگ زنده پهنای‌باند:** نمایش سرعت زنده دانلود، سرعت آپلود و پینگ واقعی به میلی‌ثانیه.
- 🛡️ **سوییچ فوری حالت تونل (در دسکتاپ):** سوییچ آسان بین **TUN Mode (وی‌پی‌ان کامل سیستمی)** و **System Proxy (پراکسی سیستمی بدون نیاز به رمز عبور روت)**.
- 🌐 **دراور هوشمند مدیریت سرورها:** همراه با پرچم کشورها، نوع پروتکل (VLESS, VMess, Trojan, Shadowsocks) و برچسب‌های رنگی پینگ (سبز، زرد، قرمز و تایم‌اوت).
- 🔄 **سیستم بررسی آپدیت خودکار (In-App Auto Updater):** بررسی خودکار انتشار نسخه جدید و دانلود مستقیم فایل‌ها بدون نیاز به فیلترشکن جانبی.
- 📥 **پشتیبانی از System Tray و Menu Bar:** مینی‌مایز شدن بی‌صدا در پس‌زمینه بدون قطع شدن اتصال هنگام بستن پنجره.

---

<a id="-راهنمای-کامل-نسخه-اندروید-android-guide"></a>
## 📱 راهنمای کامل نسخه اندروید

### ۱. نصب و راه‌اندازی اولیه
1. آخرین فایل APK را از بخش [GitHub Releases](https://github.com/Mr-Meshky/vify-cli/releases) دانلود کنید.
2. فایل را نصب و اجرا کنید.
3. در اولین اتصال، سیستم‌عامل اندروید یک پنجره باز کرده و درخواست مجوز **VPN Service** می‌کند؛ روی **OK / Allow** کلیک کنید.
4. در اندروید ۱۳ به بالا، مجوز **Notifications** برای نمایش وضعیت زنده اتصال در نوار اعلان فعال می‌شود.

### ۲. نحوه استفاده و اتصال به سریع‌ترین سرور
- **اتصال با یک لمس:** کافی است روی دکمه بزرگ مرکزی (Power Button) بزنید. برنامه به صورت خودکار سرورها را بررسی کرده و پس از تست عبور ترافیک، به بهترین سرور وصل می‌شود.
- **انتخاب دستی سرور:**
  1. روی کارت سرور در پایین صفحه کلیک کنید تا دراور سرورها باز شود.
  2. روی آیکون **صاعقه (`⚡`)** در بالای دراور بزنید تا تمام سرورها تست شوند و سریع‌ترین سرورهای سالم در بالای لیست چیده شوند.
  3. سرور مورد نظر خود را با توجه به پرچم کشور و پینگ لمس کنید. برنامه ابتدا پکت اینترنت آن را چک کرده و سپس تونل را فعال می‌کند.
- **وارد کردن کانفیگ اختصاصی:**
  - در دراور سرورها روی آیکون **`+` (Import Link)** بزنید.
  - لینک کانفیگ خود (`vless://`، `vmess://`، `trojan://` یا `ss://`) را جای‌گذاری کرده و دکمه **Connect Now** را لمس کنید.

---

### ⚡ ویژگی‌های نسخه خط فرمان (CLI)

- 🚀 **فوق‌العاده سبک و مستقل:** باینری مستقل و کم‌مصرف نوشته‌شده به زبان Go.
- 🎯 **اعتبارسنجی ترافیک در سطح کرنل:** بررسی عبور پکت‌ها با پروسه موقت sing-box.
- 🇮🇷 **بایپس خودکار سایت‌های ایرانی:** عبور مستقیم و نیم‌بهای سایت‌های داخلی و بانکی (`.ir`) بدون عبور از پروکسی.
- 🔄 **واچ‌داگ و بازیابی خودکار (Watchdog Failover):** مانیتورینگ اتصال در پس‌زمینه و سوییچ خودکار به سرور سالم در صورت افت کیفیت.
- 🔍 **اسکنر Clean-IP کلادفلر:** یافتن بهترین آی‌پی‌های تمیز متناسب با اپراتور اینترنت شما.

---

### 🚀 شروع سریع با ترمینال (CLI)

**نصب خودکار با اسکریپت یک‌خطی (macOS / Linux / Termux):**
```bash
curl -fsSL https://raw.githubusercontent.com/Mr-Meshky/vify-cli/main/scripts/install.sh | bash
```

**دستورات پرکاربرد:**
```bash
# اتصال خودکار به سریع‌ترین سرور با حالت TUN
vify connect

# اتصال در حالت System Proxy (بدون دسترسی sudo/root)
vify connect --system-proxy

# اتصال به کشور یا پروتکل دلخواه
vify connect --country DE --protocol vless

# اتصال به یک کانفیگ دلخواه با لینک مستقیم
vify connect --uri "vless://..."

# داشبورد گرافیکی ترمینال و انتخاب تعاملی سرورها
vify list

# اجرای بنچمارک پینگ و سلامت سرورها
vify test --batch 40

# بررسی وضعیت زنده اتصال، پینگ و ترافیک
vify status

# بررسی و آپدیت خودکار به آخرین نسخه برنامه
vify update

# اسکن و استخراج آی‌پی‌های تمیز کلادفلر برای اینترنت شما
vify clean-ip --count 10

# قطع اتصال و بازگردانی تنظیمات شبکه
vify disconnect
```

---

<a id="-پرسش‌های-متداول-و-عیب‌یابی-faq"></a>
## ❓ پرسش‌های متداول و عیب‌یابی (FAQ)

<details>
<summary><b>۱. چرا متصل می‌شدم اما دانلود و آپلود صفر بود؟ (رفع شده در v1.2.3 به بعد)</b></summary>
بسیاری از کانفیگ‌های عمومی از لبه شبکه کلادفلر استفاده می‌کنند که به پینگ اولیه لایه TCP سریع پاسخ می‌دهند (زیر ۱۰۰ میلی‌ثانیه)، اما اینترنت پشت آن‌ها مسدود است. در نسخه جدید، الگوریتم اعتبارسنجی دو مرحله‌ای Vify تا زمانی که تبادل واقعی دیتای اینترنت (HTTP 204) از سرور تایید نشود به آن متصل نمی‌شود.
</details>

<details>
<summary><b>۲. تفاوت حالت TUN و System Proxy چیست؟</b></summary>
<ul>
  <li><b>TUN Mode:</b> یک کارت شبکه مجازی می‌سازد و ۱۰۰٪ ترافیک کل سیستم و تمام برنامه‌ها، بازی‌ها و ترمینال را از فیلترشکن رد می‌کند (نیازمند دسترسی ادمین/رووت در دسکتاپ).</li>
  <li><b>System Proxy:</b> پراکسی را در تنظیمات شبکه سیستم ثبت می‌کند و برنامه‌های استاندارد مانند مرورگرها و تلگرام از آن استفاده می‌کنند؛ این حالت به دسترسی ادمین/رووت نیازی ندارد.</li>
</ul>
</details>

<details>
<summary><b>۳. آیا ترافیک سایت‌های ایرانی مصرف بین‌الملل حساب می‌شود؟</b></summary>
خیر، Vify مجهز به سیستم مسیریابی هوشمند ایران است. ترافیک تمام وب‌سایت‌های ایرانی (`.ir`، بانک‌ها، اسنپ، دیجی‌کالا و ...) به صورت خودکار بایپس شده و با سرعت کامل و تعرفه نیم‌بها محاسبه می‌شود.
</details>

<details>
<summary><b>۴. چگونه آدرس سابسکریپشن‌های اختصاصی خود را وارد کنم؟</b></summary>
در نسخه CLI می‌توانید فایل <code>~/.vify/config.yaml</code> را ویرایش کنید. در نسخه موبایل و دسکتاپ نیز می‌توانید لینک‌های خام کانفیگ‌ها را مستقیماً از طریق دکمه Import در دراور سرورها اضافه فرمایید.
</details>

<details>
<summary><b>۵. در مک با ارور “Vify is damaged and can’t be opened” مواجه می‌شوم، چه کنم؟</b></summary>
این پیام امنیتی Gatekeeper مک برای برنامه‌های اوپن‌سورس و بدون امضای تجاری اپل است. پس از کپی کردن برنامه به پوشه Applications، کافی است ترمینال مک (Terminal) را باز کرده و این دستور تک‌خطی را اجرا کنید:
<pre><code>xattr -cr /Applications/Vify.app</code></pre>
سپس برنامه به راحتی و بدون هیچ خطایی اجرا می‌شود.
</details>

---

<a id="english"></a>
## 🇬🇧 English Documentation

> **Modern, Blazing-Fast, Multi-Platform Anti-Censorship Client — Powered by Go, Xray & Flutter**  
> Provides resilient, unrestricted internet connectivity for users facing heavy censorship by automatically discovering, benchmarking, and routing through high-performance nodes from the [Vify Ecosystem](https://github.com/Mr-Meshky/vify).

---

### 🖥️ Desktop & Mobile Features

- ⚡ **One-Tap Smart Fast-Pass:** Concurrent node testing that locks onto the lowest latency healthy server in under a second.
- 🎯 **Real HTTP 204 Verification:** Filters out fake CDN edge TCP pings and verifies end-to-end data delivery before connecting.
- 🎨 **Modern Glassmorphism UI:** Striking dark interface with neon accents, dynamic pulsing status rings, and smooth animations.
- 📊 **Real-Time Speedometer:** Live gauges measuring download throughput, upload rates, and real round-trip latency.
- 🛡️ **Instant Tunnel Switch:** Effortlessly toggle between **TUN Mode (System Virtual Adapter)** and **System Proxy (Zero-admin mode)**.
- 🌐 **Interactive Server Drawer:** Country flags, protocol indicators, interactive ping benchmarks, and custom config imports.
- 🔄 **Built-in Auto Updater:** Integrated update checker across Desktop and Mobile with direct binary downloads.
- 📥 **System Tray & Menu Bar Resident:** Operates quietly in the background without dropping connections when windows are closed.

---

### 🚀 CLI Quick Start

**One-Line Auto-Install (macOS / Linux / Termux):**
```bash
curl -fsSL https://raw.githubusercontent.com/Mr-Meshky/vify-cli/main/scripts/install.sh | bash
```

**Essential Commands:**
```bash
# Connect to fastest node in TUN mode
vify connect

# Connect without root permissions (System Proxy)
vify connect --system-proxy

# Filter by country or protocol
vify connect --country DE --protocol vless

# Direct connect via share link
vify connect --uri "vless://..."

# Interactive terminal dashboard
vify list

# Benchmark latency leaderboard
vify test --batch 40

# Check live connection telemetry
vify status

# Check and perform automated updates
vify update

# Scan Clean Cloudflare IPs for your ISP
vify clean-ip --count 10

# Disconnect session and restore routing
vify disconnect
```

---

### ⚙️ CLI Configuration (`~/.vify/config.yaml`)

```yaml
subscriptions:
  - "https://raw.githubusercontent.com/Mr-Meshky/vify/main/configs/vless.txt"
  - "https://raw.githubusercontent.com/Mr-Meshky/vify/main/configs/vmess.txt"
  - "https://raw.githubusercontent.com/Mr-Meshky/vify/main/configs/trojan.txt"
  - "https://raw.githubusercontent.com/Mr-Meshky/vify/main/configs/ss.txt"

default_mode: "tun" # "tun" or "system_proxy"
local_socks_port: 2080
local_http_port: 2081
test_timeout_ms: 2500
concurrency_limit: 35
fastpass_threshold_ms: 800
test_url: "http://cp.cloudflare.com/generate_204"
auto_failover: true
watchdog_interval_sec: 7
direct_iran_bypass: true
log_level: "warn"
```

---

### 📄 License

Distributed under the [MIT License](LICENSE).  
Part of the **Vify Anti-Censorship Ecosystem** by [@Mr-Meshky](https://github.com/Mr-Meshky).
