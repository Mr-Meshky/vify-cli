# ⚡ vify-cli

<p align="center">
  <a href="#فارسی">🇮🇷 فارسی</a> •
  <a href="#english">🇬🇧 English</a>
</p>

---

<a id="فارسی"></a>
## 🇮🇷 فارسی

> **کلاینت VPN ترمینالی، فوق‌سریع، سبک و کراس‌پلتفرم — نوشته‌شده با Go**
> کاربران ایرانی رو به پراکسی‌های رایگان، پرسرعت و ضدفیلترینگ از [اکوسیستم Vify](https://github.com/Mr-Meshky/vify) وصل می‌کنه.

روی **macOS**، **لینوکس (سرور و دسکتاپ)**، **ویندوز** و **اندروید (Termux)** به‌طور کامل کار می‌کنه.

### 🌟 ویژگی‌های کلیدی

- 🚀 **سریع و سبک:** باینری خالص Go با مصرف حافظه‌ی کم و بدون وابستگی اضافه.
- 🎯 **تست دسته‌ای واقعی سمت کلاینت:** سرورهای کاندید رو به‌صورت هم‌زمان (۲۰ تا ۵۰ گوروتین) با endpointهای واقعی HTTP/204 (`cp.cloudflare.com` / `google.com`) و هندشیک TLS بنچمارک می‌کنه تا از عبور واقعی از DPI ایران (همراه اول، ایرانسل، رایتل، مخابرات، شاتل) مطمئن بشه.
- ✅ **اعتبارسنجی واقعی پروتکل:** برخلاف یه تست TCP/TLS ساده که هر پورت بازی رو قبول می‌کنه، قبل از اتصال نهایی یه پروسه‌ی sing-box واقعی و کوتاه‌مدت راه می‌ندازه و یه درخواست HTTP واقعی از تونل رد می‌کنه — نودهایی که واقعاً کار نمی‌کنن (حتی اگه پورتشون باز باشه) رد می‌شن و اتصال خودکار سراغ کاندید بعدی می‌ره.
- ⚡ **اتصال آنی Fast-Pass:** به محض تأیید اولین نود سالم با تأخیر زیر ۸۰۰ میلی‌ثانیه، فوراً وصل می‌شه و کشف بقیه‌ی سرورها در پس‌زمینه ادامه پیدا می‌کنه.
- 🇮🇷 **بایپس هوشمند ترافیک داخلی ایران:** قوانین مسیریابی مستقیم برای GeoIP ایران (`geoip:ir`)، GeoSite ایران (`geosite:ir`) و دامنه‌های داخلی (`.ir`، بانک‌ها، شاپرک، دیجی‌کالا، دیوار، اسنپ) برای سرعت کامل و مصرف نصف‌قیمت پهنای‌باند.
- 🛡️ **دو حالت اتصال (TUN و System Proxy):**
  - **حالت TUN (VPN کامل):** یه اینترفیس شبکه‌ی مجازی می‌سازه و کل ترافیک سیستم رو مسیریابی می‌کنه (نیاز به دسترسی root/sudo).
  - **حالت System Proxy:** پراکسی HTTP و SOCKS5 سیستمی رو بدون نیاز به دسترسی root تنظیم می‌کنه.
- 🎨 **رابط ترمینالی مدرن:** ساخته‌شده با [Charmbracelet Bubble Tea و Lipgloss](https://charm.sh) با اسپینر زنده، پرچم کشور سرورها و نمایش لحظه‌ای سرعت آپلود/دانلود.
- 🔄 **واچ‌داگ Auto-Failover:** به‌صورت پس‌زمینه throttle یا قطعی ناشی از DPI رو تشخیص می‌ده و فوراً به بهترین سرور سالم بعدی سوییچ می‌کنه — چه به‌خاطر throttle تدریجی (پایش دوره‌ای) و چه کرش ناگهانی پروسه (تشخیص آنی).
- 🔍 **اسکنر Clean-IP کلادفلر:** به‌صورت خودکار از رنج‌های واقعی IP کلادفلر نمونه‌برداری تصادفی می‌کنه (نه فقط یه لیست ثابت) تا آی‌پی‌های تمیز و کم‌تأخیر مخصوص ISP خودت رو کشف کنه.
- 💾 **کش آفلاین:** کانفیگ‌های سالم رو در `~/.vify/cache.json` ذخیره می‌کنه تا موقع قطعی ناگهانی هم بشه سریع وصل شد.

### 📦 نصب

**اسکریپت یک‌خطی (macOS / لینوکس / Termux):**
```bash
curl -fsSL https://raw.githubusercontent.com/Mr-Meshky/vify-cli/main/scripts/install.sh | bash
```

**نصب با Go:**
```bash
go install github.com/Mr-Meshky/vify-cli@latest
```

**بیلد از سورس:**
```bash
git clone https://github.com/Mr-Meshky/vify-cli.git
cd vify-cli
make build
# باینری در ./bin/vify ساخته می‌شه
```

### 🚀 شروع سریع و دستورات

**۱. اتصال آنی (انتخاب خودکار بهترین سرور):**
```bash
vify connect
```

**۲. اتصال با فیلتر و گزینه‌های دلخواه:**
```bash
# اتصال به یه سرور VLESS آلمانی در حالت TUN (VPN کامل)
vify connect --country DE --protocol vless --tun

# اتصال در حالت System Proxy با دسته‌ی تست بزرگ‌تر
vify connect --country NL --batch 80
```

**۳. انتخاب تعاملی سرور (TUI):**
```bash
vify list
```

**۴. بنچمارک و رتبه‌بندی تأخیر:**
```bash
vify test --batch 50
```

**۵. بررسی وضعیت اتصال:**
```bash
vify status
```

**۶. قطع اتصال و بازگردانی تنظیمات شبکه:**
```bash
vify disconnect
```

**۷. پیدا کردن آی‌پی‌های تمیز کلادفلر:**
```bash
vify clean-ip --count 10
```

> **نکته:** حالت TUN (پیش‌فرض `vify connect`) چون یه اینترفیس شبکه‌ی مجازی می‌سازه، نیاز به `sudo` داره:
> ```bash
> sudo vify connect
> ```
> اگه نمی‌خوای دسترسی root بدی، از حالت System Proxy استفاده کن که بدون sudo کار می‌کنه:
> ```bash
> vify connect --system-proxy
> ```

### ⚙️ پیکربندی (`~/.vify/config.yaml`)

می‌تونی endpointهای subscription، پورت‌ها و timeoutها رو در `~/.vify/config.yaml` شخصی‌سازی کنی:

```yaml
subscriptions:
  - "https://raw.githubusercontent.com/Mr-Meshky/vify/main/sub/all_normal.txt"
  - "https://raw.githubusercontent.com/Mr-Meshky/vify/main/sub/reality_sub.txt"
  - "https://raw.githubusercontent.com/Mr-Meshky/vify/main/sub/vless_sub.txt"
  - "https://raw.githubusercontent.com/Mr-Meshky/vify/main/sub/vmess_sub.txt"
  - "https://raw.githubusercontent.com/Mr-Meshky/vify/main/sub/trojan_sub.txt"
  - "https://raw.githubusercontent.com/Mr-Meshky/vify/main/sub/ss_sub.txt"

default_mode: "system_proxy" # "system_proxy" یا "tun"
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

### 🏗️ معماری

```
                        ┌──────────────────────────────────────────────┐
                        │              vify-cli (Cobra)                │
                        │  connect | list | test | status | clean-ip   │
                        └──────────────────────┬───────────────────────┘
                                               │
               ┌───────────────────────────────┼───────────────────────────────┐
               │                               │                               │
               ▼                               ▼                               ▼
    ┌────────────────────┐          ┌────────────────────┐          ┌────────────────────┐
    │  Subscription &    │          │ Concurrency Engine │          │ Interactive TUI    │
    │  Config Parser     │          │ & Fast-Pass Tester │          │ (Bubble Tea /      │
    │  (VLESS, VMess,    │          │ (Real HTTP/204 +   │          │  Lipgloss)         │
    │   Trojan, SS)      │          │  Real-Proxy Verify)│          │                    │
    └──────────┬─────────┘          └──────────┬─────────┘          └────────────────────┘
               │                               │
               └────────────────┬──────────────┘
                                │
                                ▼
               ┌─────────────────────────────────┐
               │         VPN Core Engine         │
               │  - TUN Mode (sing-box / gVisor) │
               │  - System Proxy Mode            │
               │  - Smart Iran Bypass (.ir, etc) │
               │  - Live Bandwidth / Stats       │
               │  - Watchdog & Auto-Failover     │
               └─────────────────────────────────┘
```

### 🛠️ کراس-کامپایل

با `Makefile` می‌تونی برای همه‌ی پلتفرم‌های پشتیبانی‌شده بیلد بگیری:

```bash
make cross-build
```

باینری‌های تولیدشده در `./bin/`:
- `vify-darwin-arm64` (مک اپل سیلیکون)
- `vify-darwin-amd64` (مک اینتل)
- `vify-linux-amd64` (لینوکس x86_64)
- `vify-linux-arm64` (لینوکس ARM64 / رزبری‌پای / سرور)
- `vify-windows-amd64.exe` (ویندوز x64)

### 📄 لایسنس
تحت [لایسنس MIT](LICENSE) منتشر شده.
بخشی از **اکوسیستم ضدفیلترینگ Vify** ساخته‌ی [@Mr-Meshky](https://github.com/Mr-Meshky).

---

<a id="english"></a>
## 🇬🇧 English

> **Blazing-Fast, Lightweight, Cross-Platform Terminal VPN Client in Go**
> Connects Iranian users to free, high-speed, anti-censorship proxies fetched directly from the [Vify Ecosystem](https://github.com/Mr-Meshky/vify).

Works seamlessly on **macOS**, **Linux (Servers & Desktops)**, **Windows**, and **Android (Termux)**.

### 🌟 Highlights & Key Features

- 🚀 **Blazing-Fast & Lightweight:** Pure Go binary with minimal memory footprint and zero external bloat.
- 🎯 **Client-Side Real-World Batch Testing:** Concurrently benchmarks candidate servers (20-50 goroutines) using real HTTP/204 endpoints (`cp.cloudflare.com` / `google.com`) and TLS handshakes to guarantee bypass through Iranian DPI (MCI, Irancell, Rightel, TCI, Shatel).
- ✅ **Real Protocol Verification:** Unlike a plain TCP/TLS dial (which any open port passes), Vify spins up a short-lived, real sing-box process and performs an actual HTTP request through it before committing to a connection — nodes that don't genuinely relay traffic are rejected and the next best candidate is tried automatically.
- ⚡ **Fast-Pass Instant Connection:** Automatically connects as soon as the first sub-800ms healthy node is verified, continuing background discovery without blocking.
- 🇮🇷 **Smart Iran Domestic Direct Bypass:** Built-in direct routing rules for Iranian GeoIP (`geoip:ir`), GeoSite (`geosite:ir`), and domestic domains (`.ir`, banking, Shaparak, Digikala, Divar, Snapp) for full speed and half-price bandwidth.
- 🛡️ **Dual Modes (TUN & System Proxy):**
  - **TUN Mode (Full VPN):** Creates a virtual network interface and routes all device/system traffic (requires root/sudo).
  - **System Proxy Mode:** Configures system-wide HTTP & SOCKS5 proxies without requiring root/sudo privileges.
- 🎨 **Modern Terminal UI:** Built with [Charmbracelet Bubble Tea & Lipgloss](https://charm.sh) featuring live spinners, server country flags, and real-time upload/download speed gauges.
- 🔄 **Auto-Failover Watchdog:** Detects DPI throttling or drops — both gradual (periodic health probes) and sudden (instant process-crash detection) — and fails over to the next healthy cached node automatically.
- 🔍 **Cloudflare Clean-IP Scanner:** Randomly samples real Cloudflare CIDR ranges (not just a static list) to discover low-latency, zero-loss edge IPs tailored to your ISP.
- 💾 **Offline Cache:** Saves healthy configs to `~/.vify/cache.json` for immediate connection during sudden network blackouts.

### 📦 Installation

**Quick One-Line Script (macOS / Linux / Termux):**
```bash
curl -fsSL https://raw.githubusercontent.com/Mr-Meshky/vify-cli/main/scripts/install.sh | bash
```

**Go Install:**
```bash
go install github.com/Mr-Meshky/vify-cli@latest
```

**Build from Source:**
```bash
git clone https://github.com/Mr-Meshky/vify-cli.git
cd vify-cli
make build
# Binary will be placed in ./bin/vify
```

### 🚀 Quick Start & Commands

**1. Connect Immediately (Auto-Select Best Server):**
```bash
vify connect
```

**2. Connect with Filters & Options:**
```bash
# Connect to a German VLESS server in TUN (Full VPN) mode
vify connect --country DE --protocol vless --tun

# Connect in System Proxy mode with a larger test batch
vify connect --country NL --batch 80
```

**3. Interactive Server Selection (TUI):**
```bash
vify list
```

**4. Benchmark & Latency Leaderboard:**
```bash
vify test --batch 50
```

**5. Check Connection Status:**
```bash
vify status
```

**6. Disconnect & Restore Network Settings:**
```bash
vify disconnect
```

**7. Find Cloudflare Clean IPs:**
```bash
vify clean-ip --count 10
```

> **Note:** TUN mode (the default for `vify connect`) creates a virtual network interface, so it requires `sudo`:
> ```bash
> sudo vify connect
> ```
> If you'd rather not grant root access, use System Proxy mode instead — no sudo required:
> ```bash
> vify connect --system-proxy
> ```

### ⚙️ Configuration (`~/.vify/config.yaml`)

You can customize subscription endpoints, ports, and timeouts in `~/.vify/config.yaml`:

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

### 🏗️ Architecture

```
                        ┌──────────────────────────────────────────────┐
                        │              vify-cli (Cobra)                │
                        │  connect | list | test | status | clean-ip   │
                        └──────────────────────┬───────────────────────┘
                                               │
               ┌───────────────────────────────┼───────────────────────────────┐
               │                               │                               │
               ▼                               ▼                               ▼
    ┌────────────────────┐          ┌────────────────────┐          ┌────────────────────┐
    │  Subscription &    │          │ Concurrency Engine │          │ Interactive TUI    │
    │  Config Parser     │          │ & Fast-Pass Tester │          │ (Bubble Tea /      │
    │  (VLESS, VMess,    │          │ (Real HTTP/204 +   │          │  Lipgloss)         │
    │   Trojan, SS)      │          │  Real-Proxy Verify)│          │                    │
    └──────────┬─────────┘          └──────────┬─────────┘          └────────────────────┘
               │                               │
               └────────────────┬──────────────┘
                                │
                                ▼
               ┌─────────────────────────────────┐
               │         VPN Core Engine         │
               │  - TUN Mode (sing-box / gVisor) │
               │  - System Proxy Mode            │
               │  - Smart Iran Bypass (.ir, etc) │
               │  - Live Bandwidth / Stats       │
               │  - Watchdog & Auto-Failover     │
               └─────────────────────────────────┘
```

### 🛠️ Cross-Compilation

Use the included `Makefile` to compile binaries for all supported platforms:

```bash
make cross-build
```

Generated binaries in `./bin/`:
- `vify-darwin-arm64` (macOS Apple Silicon)
- `vify-darwin-amd64` (macOS Intel)
- `vify-linux-amd64` (Linux x86_64)
- `vify-linux-arm64` (Linux ARM64 / Raspberry Pi / Servers)
- `vify-windows-amd64.exe` (Windows x64)

### 📄 License
Released under the [MIT License](LICENSE).
Part of the **Vify Anti-Censorship Ecosystem** by [@Mr-Meshky](https://github.com/Mr-Meshky).
