package core

import (
	"fmt"
	"os/exec"
	"runtime"
	"strings"
)

// SystemProxyController manages enabling and disabling OS system proxy
type SystemProxyController struct {
	httpPort  int
	socksPort int
}

// NewSystemProxyController creates a new controller
func NewSystemProxyController(httpPort, socksPort int) *SystemProxyController {
	return &SystemProxyController{
		httpPort:  httpPort,
		socksPort: socksPort,
	}
}

// Enable sets the system proxy for the current OS
func (s *SystemProxyController) Enable() error {
	switch runtime.GOOS {
	case "darwin":
		return s.enableDarwin()
	case "linux":
		return s.enableLinux()
	case "windows":
		return s.enableWindows()
	default:
		return nil
	}
}

// Disable turns off system proxy settings and restores original state
func (s *SystemProxyController) Disable() error {
	switch runtime.GOOS {
	case "darwin":
		return s.disableDarwin()
	case "linux":
		return s.disableLinux()
	case "windows":
		return s.disableWindows()
	default:
		return nil
	}
}

func (s *SystemProxyController) enableDarwin() error {
	services := getDarwinActiveServices()
	for _, svc := range services {
		// Set HTTP proxy
		_ = exec.Command("networksetup", "-setwebproxy", svc, "127.0.0.1", fmt.Sprintf("%d", s.httpPort)).Run()
		_ = exec.Command("networksetup", "-setwebproxystate", svc, "on").Run()

		// Set HTTPS proxy
		_ = exec.Command("networksetup", "-setsecurewebproxy", svc, "127.0.0.1", fmt.Sprintf("%d", s.httpPort)).Run()
		_ = exec.Command("networksetup", "-setsecurewebproxystate", svc, "on").Run()

		// Set SOCKS proxy
		_ = exec.Command("networksetup", "-setsocksfirewallproxy", svc, "127.0.0.1", fmt.Sprintf("%d", s.socksPort)).Run()
		_ = exec.Command("networksetup", "-setsocksfirewallproxystate", svc, "on").Run()

		// Bypass local and .ir domains
		_ = exec.Command("networksetup", "-setproxybypassdomains", svc, "127.0.0.1", "localhost", "192.168.0.0/16", "10.0.0.0/8", "*.ir").Run()
	}
	return nil
}

func (s *SystemProxyController) disableDarwin() error {
	services := getDarwinActiveServices()
	for _, svc := range services {
		_ = exec.Command("networksetup", "-setwebproxystate", svc, "off").Run()
		_ = exec.Command("networksetup", "-setsecurewebproxystate", svc, "off").Run()
		_ = exec.Command("networksetup", "-setsocksfirewallproxystate", svc, "off").Run()
	}
	return nil
}

func getDarwinActiveServices() []string {
	out, err := exec.Command("networksetup", "-listallnetworkservices").Output()
	if err != nil {
		return []string{"Wi-Fi", "Ethernet"}
	}

	lines := strings.Split(string(out), "\n")
	var active []string
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.Contains(line, "*") {
			continue
		}
		active = append(active, line)
	}

	if len(active) == 0 {
		active = []string{"Wi-Fi", "Ethernet"}
	}
	return active
}

func (s *SystemProxyController) enableLinux() error {
	// gsettings only exists on GNOME. On KDE/XFCE/i3/headless servers this
	// silently does nothing, which used to leave the user thinking they were
	// proxied system-wide when literally no application outside of one
	// manually pointed at these ports was. Since the README explicitly
	// advertises "Linux (Servers & Desktops)" and non-root system-proxy mode
	// as the go-to for servers (which never have a GNOME session at all),
	// print the manual fallback instead of failing silently.
	if _, err := exec.LookPath("gsettings"); err != nil {
		fmt.Printf("\n⚠ No GNOME desktop detected (gsettings not found) — system-wide proxy settings were not applied automatically.\n")
		fmt.Printf("  The local proxy is still up; point apps at it manually, e.g.:\n")
		fmt.Printf("    export http_proxy=http://127.0.0.1:%d https_proxy=http://127.0.0.1:%d all_proxy=socks5://127.0.0.1:%d\n\n", s.httpPort, s.httpPort, s.socksPort)
		return nil
	}

	_ = exec.Command("gsettings", "set", "org.gnome.system.proxy", "mode", "'manual'").Run()
	_ = exec.Command("gsettings", "set", "org.gnome.system.proxy.http", "host", "'127.0.0.1'").Run()
	_ = exec.Command("gsettings", "set", "org.gnome.system.proxy.http", "port", fmt.Sprintf("%d", s.httpPort)).Run()
	_ = exec.Command("gsettings", "set", "org.gnome.system.proxy.https", "host", "'127.0.0.1'").Run()
	_ = exec.Command("gsettings", "set", "org.gnome.system.proxy.https", "port", fmt.Sprintf("%d", s.httpPort)).Run()
	_ = exec.Command("gsettings", "set", "org.gnome.system.proxy.socks", "host", "'127.0.0.1'").Run()
	_ = exec.Command("gsettings", "set", "org.gnome.system.proxy.socks", "port", fmt.Sprintf("%d", s.socksPort)).Run()
	_ = exec.Command("gsettings", "set", "org.gnome.system.proxy", "ignore-hosts", "['localhost', '127.0.0.0/8', '::1', '*.ir']").Run()
	return nil
}

func (s *SystemProxyController) disableLinux() error {
	_ = exec.Command("gsettings", "set", "org.gnome.system.proxy", "mode", "'none'").Run()
	return nil
}

func (s *SystemProxyController) enableWindows() error {
	proxyServer := fmt.Sprintf("http=127.0.0.1:%d;https=127.0.0.1:%d;socks=127.0.0.1:%d", s.httpPort, s.httpPort, s.socksPort)
	_ = exec.Command("reg", "add", `HKCU\Software\Microsoft\Windows\CurrentVersion\Internet Settings`, "/v", "ProxyEnable", "/t", "REG_DWORD", "/d", "1", "/f").Run()
	_ = exec.Command("reg", "add", `HKCU\Software\Microsoft\Windows\CurrentVersion\Internet Settings`, "/v", "ProxyServer", "/t", "REG_SZ", "/d", proxyServer, "/f").Run()
	_ = exec.Command("reg", "add", `HKCU\Software\Microsoft\Windows\CurrentVersion\Internet Settings`, "/v", "ProxyOverride", "/t", "REG_SZ", "/d", "<local>;*.ir;127.0.0.1", "/f").Run()
	return nil
}

func (s *SystemProxyController) disableWindows() error {
	_ = exec.Command("reg", "add", `HKCU\Software\Microsoft\Windows\CurrentVersion\Internet Settings`, "/v", "ProxyEnable", "/t", "REG_DWORD", "/d", "0", "/f").Run()
	return nil
}
