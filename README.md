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
  <a href="#فارسی">🇮🇷 فارسی</a> •
  <a href="#english">🇬🇧 English</a> •
  <a href="#-دانلود-و-نصب--downloads">📥 دانلود / Downloads</a>
</p>

---

<a id="-دانلود-و-نصب--downloads"></a>
## 📥 دانلود و نصب / Downloads

تمام نسخه‌ها از صفحه **[GitHub Releases](https://github.com/Mr-Meshky/vify-cli/releases)** قابل دریافت هستند:

| سیستم‌عامل / پلتفرم | نوع برنامه | فرمت فایل | راهنمای اجرا |
| :--- | :--- | :--- | :--- |
| 🍏 **macOS** (Universal) | 🖥️ Desktop GUI | `.dmg` / `.zip` (شامل `Vify.app` + CLI) | نصب فایل DMG یا انتقال `Vify.app` به `/Applications` |
| 🪟 **Windows** (x64) | 🖥️ Desktop GUI | Setup `.exe` / `.zip` | نصب خودکار با اینستالر یا اجرای مستقیم `Vify.exe` |
| 🐧 **Linux** (x86_64) | 🖥️ Desktop GUI | `.deb` / `.tar.gz` | نصب با `sudo dpkg -i` یا اجرای `vify-ui` |
| 🤖 **Android** | 📱 Mobile App | `.apk` (Universal / ARM64) | نصب فایل APK روی گوشی اندرویدی |
| 💻 **Linux CLI** (x86_64) | ⚡ Terminal CLI | باینری مستقل | `chmod +x vify-cli-linux-amd64 && ./vify-cli-linux-amd64` |
| 💻 **macOS CLI** | ⚡ Terminal CLI | باینری مستقل Universal | `chmod +x vify-cli-darwin-universal && ./vify-cli-darwin-universal` |

---

<a id="فارسی"></a>
## 🇮🇷 راهنمای فارسی

> **کلاینت ضدفیلترینگ مدرن، فوق‌سریع و کراس‌پلتفرم — با هسته پرقدرت Go و رابط کاربری شیشه‌ای Flutter**  
> دسترسی آسان و پایدار به پراکسی‌های پرسرعت، سالم و ضدفیلتر از [اکوسیستم Vify](https://github.com/Mr-Meshky/vify) برای دور زدن فیلترینگ همراه اول، ایرانسل، رایتل، مخابرات و شاتل.

---

### 🖥️ ویژگی‌های اپلیکیشن گرافیکی (Desktop & Mobile)

- ⚡ **اتصال تک‌کلیکه (One-Tap Connect):** انتخاب آنی و بدون وقفه سریع‌ترین سرور سالم در کمتر از ۱ ثانیه.
- 🎨 **طراحی نئونی و مدرن (Modern Glassmorphism):** رابط کاربری تاریک، انیمیشن‌های پویا و حلقه وضعیت ضربان‌دار.
- 📊 **مانیتورینگ زنده پهنای‌باند:** نمایش لحظه‌ای سرعت دانلود، آپلود و پینگ (Latency) به میلی‌ثانیه.
- 🛡️ **سوییچ سریع حالت تونل:** جابه‌جایی آنی با یک کلیک بین **TUN Mode (وی‌پی‌ان کامل سیستمی)** و **System Proxy (پراکسی سیستمی بدون نیاز به پسورد ادمین)**.
- 🌐 **دراور هوشمند انتخاب سرور:** مشاهده لیست سرورها همراه با پرچم کشور، پروتکل و پینگ زنده.
- 📥 **پشتیبانی از System Tray و Menu Bar:** کارکرد بی‌صدا در پس‌زمینه سیستم؛ بستن پنجره اتصال را قطع نمی‌کند و دسترسی سریع در منوبار مک یا تسک‌بار ویندوز و لینوکس باقی می‌ماند.

---

### ⚡ ویژگی‌های هسته و خط فرمان (CLI)

- 🚀 **فوق‌العاده سبک و کم‌مصرف:** نوشته‌شده با زبان Go بدون تحمیل بار پردازشی به سیستم.
- 🎯 **تست دسته‌ای واقعی (Real Verification):** اجرای تست پروتکل واقعی از طریق پروسه موقت sing-box و اعتبارسنجی عبور بسته از فایروال DPI.
- ⚡ **الگوریتم Fast-Pass:** اتصال سریع به محض کشف اولین سرور زیر ۸۰۰ میلی‌ثانیه، و تکمیل پایش در پس‌زمینه.
- 🇮🇷 **بایپس خودکار سایت‌های ایران:** عبور مستقیم و بدون فیلتر ترافیک سایت‌های داخلی (`.ir`، بانک‌ها، اسنپ، دیوار، دیجی‌کالا و ...) با ترافیک نیم‌بها.
- 🔄 **واچ‌داگ و سوییچ خودکار (Auto-Failover):** شناسایی فوری افت سرعت یا مسدودی نود و جابه‌جایی بلادرنگ به سرور سالم بعدی.
- 🔍 **اسکنر Clean-IP کلادفلر:** شناسایی خودکار بهترین IPهای تمیز کلادفلر متناسب با ارائه‌دهنده اینترنت (ISP) شما.

---

### 🚀 شروع سریع با ترمینال (CLI)

**نصب سریع با اسکریپت یک‌خطی (macOS / Linux / Termux):**
```bash
curl -fsSL https://raw.githubusercontent.com/Mr-Meshky/vify-cli/main/scripts/install.sh | bash
```

**دستورات پرکاربرد:**
```bash
# اتصال خودکار به سریع‌ترین نود (حالت پیش‌فرض TUN)
vify connect

# اتصال بدون دسترسی روت (حالت System Proxy)
vify connect --system-proxy

# اتصال با فیلتر کشور یا پروتکل دلخواه
vify connect --country DE --protocol vless

# انتخاب تعاملی سرور از لیست (TUI)
vify list

# اجرای بنچمارک و رتبه‌بندی پینگ سرورها
vify test --batch 50

# بررسی وضعیت زنده اتصال
vify status

# قطع اتصال و بازگردانی تنظیمات شبکه
vify disconnect

# اسکن آی‌پی‌های تمیز کلادفلر
vify clean-ip --count 10
```

---

<a id="english"></a>
## 🇬🇧 English Documentation

> **Modern, Blazing-Fast, Multi-Platform Anti-Censorship Client — Powered by Go & Modern Flutter UI**  
> Provides resilient, unrestricted internet connectivity for Iranian users by relaying through high-performance nodes from the [Vify Ecosystem](https://github.com/Mr-Meshky/vify).

---

### 🖥️ Desktop & Mobile GUI Features

- ⚡ **One-Tap Instant Connection:** Benchmarks and locks onto the fastest healthy node in sub-seconds.
- 🎨 **State-of-the-Art Glassmorphism:** Sleek dark mode design, glowing status rings, and fluid micro-animations.
- 📊 **Real-Time Speedometer:** Live gauges measuring download bandwidth, upload throughput, and ping latency.
- 🛡️ **Instant Mode Switch:** Effortlessly switch between **TUN Mode (Full Virtual Adapter)** and **System Proxy (No root required)**.
- 🌐 **Interactive Server Drawer:** Browse nodes categorized by country, protocol, and verified round-trip ping.
- 📥 **System Tray & Menu Bar Resident:** Runs quietly in the background without dropping VPN sessions when the window is closed.

---

### ⚡ CLI & Core Capabilities

- 🚀 **High Performance, Zero Dependencies:** Pure Go core engineered for minimal CPU and memory footprints.
- 🎯 **Real-World Protocol Verification:** Spins up an ephemeral sing-box instance and transmits verified HTTP/204 payloads to defeat deep packet inspection (DPI).
- ⚡ **Fast-Pass Instant Connection:** Triggers connection upon verifying the first node under 800ms, continuing discovery asynchronously.
- 🇮🇷 **Smart Iran Domestic Bypass:** Direct routing for Iranian GeoIP (`geoip:ir`), GeoSite (`geosite:ir`), and domestic domains (`.ir`, banks, local apps) ensuring half-price bandwidth and full ISP speed.
- 🔄 **Continuous Watchdog & Auto-Failover:** Automatically detects ISP throttling or server dropouts and reroutes through the next best node without manual intervention.
- 🔍 **Cloudflare Clean-IP Discovery:** Randomly samples active Cloudflare edge subnets to locate the lowest-latency IP for your current provider.

---

### 🚀 CLI Quick Start

**One-Line Installer (macOS / Linux / Termux):**
```bash
curl -fsSL https://raw.githubusercontent.com/Mr-Meshky/vify-cli/main/scripts/install.sh | bash
```

**Common Commands:**
```bash
# Instant auto-connect to lowest latency server
vify connect

# Connect in System Proxy mode (no root/sudo needed)
vify connect --system-proxy

# Filter by country or protocol
vify connect --country DE --protocol vless

# Interactive terminal dashboard
vify list

# Benchmark latency leaderboard
vify test --batch 50

# Check connection status & live throughput
vify status

# Disconnect and reset network routes
vify disconnect

# Discover optimized Cloudflare Clean IPs
vify clean-ip --count 10
```

---

### ⚙️ Configuration (`~/.vify/config.yaml`)

You can customize subscription endpoints, ports, and watchdog thresholds in `~/.vify/config.yaml`:

```yaml
subscriptions:
  - "https://raw.githubusercontent.com/Mr-Meshky/vify/main/sub/all_normal.txt"
  - "https://raw.githubusercontent.com/Mr-Meshky/vify/main/sub/reality_sub.txt"
  - "https://raw.githubusercontent.com/Mr-Meshky/vify/main/sub/vless_sub.txt"
  - "https://raw.githubusercontent.com/Mr-Meshky/vify/main/sub/vmess_sub.txt"
  - "https://raw.githubusercontent.com/Mr-Meshky/vify/main/sub/trojan_sub.txt"
  - "https://raw.githubusercontent.com/Mr-Meshky/vify/main/sub/ss_sub.txt"

default_mode: "system_proxy" # "system_proxy" or "tun"
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
