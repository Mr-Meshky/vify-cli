package ping

import (
	"context"
	"crypto/tls"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"time"

	"github.com/Mr-Meshky/vify-cli/internal/model"
)

// DialTest executes a real-world latency test for a proxy node
// It tests socket connection + TLS handshake (with SNI/ALPN) and measures true round-trip time.
func DialTest(ctx context.Context, node *model.ProxyNode, testURL string) (time.Duration, error) {
	if node.Server == "" || node.Port == 0 {
		return 0, fmt.Errorf("invalid server/port")
	}

	targetAddr := fmt.Sprintf("%s:%d", node.Server, node.Port)
	start := time.Now()

	// 1. Direct TCP dial with timeout
	var d net.Dialer
	conn, err := d.DialContext(ctx, "tcp", targetAddr)
	if err != nil {
		return 0, fmt.Errorf("tcp dial failed: %w", err)
	}
	defer conn.Close()

	// 2. If TLS or Reality is specified, perform TLS handshake to verify DPI bypass
	if node.Security == "tls" || node.Security == "reality" {
		serverName := node.SNI
		if serverName == "" {
			serverName = node.Host
		}
		if serverName == "" {
			serverName = node.Server
		}

		tlsConfig := &tls.Config{
			ServerName:         serverName,
			InsecureSkipVerify: true,
		}

		tlsConn := tls.Client(conn, tlsConfig)
		tlsConn.SetDeadline(time.Now().Add(2 * time.Second))
		if err := tlsConn.HandshakeContext(ctx); err != nil {
			return 0, fmt.Errorf("tls handshake failed: %w", err)
		}
	}

	latency := time.Since(start)
	return latency, nil
}

// HTTP204ViaLocalProxy performs an actual HTTP GET request through a running
// local SOCKS5 or HTTP proxy. Each call typically costs the caller two full
// round-trips through the remote node (a DNS lookup, then the HTTP fetch),
// so timeout should be sized generously on high-latency links rather than
// assuming a single-hop request budget.
func HTTP204ViaLocalProxy(ctx context.Context, proxyURLStr, targetURL string, timeout time.Duration) (time.Duration, int, error) {
	proxyURL, err := url.Parse(proxyURLStr)
	if err != nil {
		return 0, 0, err
	}

	transport := &http.Transport{
		Proxy: http.ProxyURL(proxyURL),
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
		DisableKeepAlives: true,
	}

	client := &http.Client{
		Transport: transport,
		Timeout:   timeout,
	}

	start := time.Now()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, targetURL, nil)
	if err != nil {
		return 0, 0, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return 0, 0, err
	}
	defer resp.Body.Close()

	latency := time.Since(start)
	return latency, resp.StatusCode, nil
}
