package core

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/Mr-Meshky/vify-cli/internal/config"
	"github.com/Mr-Meshky/vify-cli/internal/model"
	"github.com/Mr-Meshky/vify-cli/internal/ping"
)

// VerifyPortBase is the start of the scratch port range used only by
// short-lived VerifyNodeReal processes. Deliberately distinct from the
// user's configured ports so verification never collides with an
// already-running session. Callers running N verifications concurrently
// should pass ports spaced apart within this range (see app.verifyRealConnectivity).
const VerifyPortBase = 29000

// VerifyNodeReal confirms a candidate node's proxy protocol actually works,
// not just that its TCP port is open. The initial benchmark (ping.DialTest)
// only dials a raw socket and, for TLS nodes, completes a generic TLS
// handshake with certificate verification disabled — that passes for any
// listening port, or any server willing to complete a TLS handshake,
// regardless of whether a real VLESS/VMess/Trojan/Shadowsocks proxy is
// actually behind it. This instead starts a short-lived, real sing-box
// process for the node and performs one genuine HTTP request through it.
// socksPort/httpPort must be unique per concurrently-running call.
func VerifyNodeReal(ctx context.Context, node *model.ProxyNode, testURL string, socksPort, httpPort int) (time.Duration, error) {
	sbConfig, err := GenerateConfig(node, model.ModeSystemProxy, socksPort, httpPort, 0, false)
	if err != nil {
		return 0, fmt.Errorf("config generation failed: %w", err)
	}

	configJSON, err := sbConfig.ToJSON()
	if err != nil {
		return 0, err
	}

	cfgFilePath := filepath.Join(config.GetVifyDir(), fmt.Sprintf("verify-%d.json", httpPort))
	if err := os.WriteFile(cfgFilePath, configJSON, 0644); err != nil {
		return 0, err
	}
	defer os.Remove(cfgFilePath)

	singBoxPath, err := EnsureSingBoxBinary()
	if err != nil {
		return 0, err
	}

	// 12s budget: 400ms for sing-box to bind its listeners, then room for the
	// HTTP204ViaLocalProxy call below, which itself needs two full round-trips
	// through the remote node (a proxied DNS lookup, then the HTTP fetch) —
	// tight on higher-latency links if squeezed into a single-hop timeout.
	procCtx, cancel := context.WithTimeout(ctx, 12*time.Second)
	defer cancel()

	cmd := exec.CommandContext(procCtx, singBoxPath, "run", "-c", cfgFilePath)
	if err := cmd.Start(); err != nil {
		return 0, fmt.Errorf("failed to start verification process: %w", err)
	}
	defer func() {
		if cmd.Process != nil {
			_ = cmd.Process.Kill()
			_ = cmd.Wait()
		}
	}()

	select {
	case <-time.After(400 * time.Millisecond):
		// give sing-box a moment to bind its local listeners
	case <-procCtx.Done():
		return 0, fmt.Errorf("verification timed out during startup")
	}

	proxyURL := fmt.Sprintf("http://127.0.0.1:%d", httpPort)
	lat, code, err := ping.HTTP204ViaLocalProxy(procCtx, proxyURL, testURL, 8*time.Second)
	if err != nil {
		return 0, fmt.Errorf("real proxy request failed: %w", err)
	}
	if code != 204 && code != 200 {
		return 0, fmt.Errorf("unexpected status %d from real proxy request", code)
	}

	return lat, nil
}
