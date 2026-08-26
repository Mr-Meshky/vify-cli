package core

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/Mr-Meshky/vify-cli/internal/cache"
	"github.com/Mr-Meshky/vify-cli/internal/config"
	"github.com/Mr-Meshky/vify-cli/internal/model"
)

// clashAPIPort is the local port sing-box's Clash-compatible API listens on
// for the main session, used to pull live traffic totals for the status
// dashboard. Verification runs (internal/core/verify.go) pass 0 to
// GenerateConfig instead, since they run multiple short-lived instances
// concurrently and must not all fight over this same port.
const clashAPIPort = 9090

var clashAPIAddr = fmt.Sprintf("127.0.0.1:%d", clashAPIPort)

// Engine orchestrates the VPN / Proxy connection, sing-box runner, TUN and System Proxy modes
type Engine struct {
	mu           sync.Mutex
	cfg          *config.AppConfig
	currentNode  *model.ProxyNode
	mode         model.ConnectionMode
	sysProxy     *SystemProxyController
	stats        *TrafficStats
	watchdog     *Watchdog
	cmd          *exec.Cmd
	exitCh       chan struct{} // closed exactly once, when cmd.Wait() returns for the current process
	ctx          context.Context
	cancel       context.CancelFunc
	running      bool
	healthyNodes []*model.ProxyNode
	nodeIndex    int
}

// NewEngine creates a new VPN/Proxy Engine instance
func NewEngine(cfg *config.AppConfig) *Engine {
	return &Engine{
		cfg:      cfg,
		sysProxy: NewSystemProxyController(cfg.LocalHTTPPort, cfg.LocalSocksPort),
		stats:    NewTrafficStats(),
	}
}

// StartNode initiates connection to a specific proxy node
func (e *Engine) StartNode(ctx context.Context, node *model.ProxyNode, mode model.ConnectionMode, backupNodes []*model.ProxyNode) error {
	e.mu.Lock()
	defer e.mu.Unlock()

	if e.running {
		_ = e.stopInternal()
	}

	e.currentNode = node
	e.mode = mode
	e.healthyNodes = backupNodes
	e.nodeIndex = 0
	e.ctx, e.cancel = context.WithCancel(ctx)

	// 1. Generate sing-box configuration
	sbConfig, err := GenerateConfig(node, mode, e.cfg.LocalSocksPort, e.cfg.LocalHTTPPort, clashAPIPort, e.cfg.DirectIranBypass)
	if err != nil {
		return fmt.Errorf("failed to generate sing-box configuration: %w", err)
	}

	configJSON, err := sbConfig.ToJSON()
	if err != nil {
		return fmt.Errorf("failed to marshal sing-box configuration: %w", err)
	}

	cfgFilePath := filepath.Join(config.GetVifyDir(), "singbox.json")
	if err := os.WriteFile(cfgFilePath, configJSON, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	// 2. Launch sing-box process
	singBoxPath, err := EnsureSingBoxBinary()
	if err != nil {
		return fmt.Errorf("failed to locate or download sing-box core: %w", err)
	}

	logFile, _ := os.OpenFile(filepath.Join(config.GetVifyDir(), "singbox.log"), os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	cmd := exec.CommandContext(e.ctx, singBoxPath, "run", "-c", cfgFilePath)
	if logFile != nil {
		cmd.Stdout = logFile
		cmd.Stderr = logFile
	} else {
		cmd.Stdout = io.Discard
		cmd.Stderr = io.Discard
	}
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to start sing-box: %w", err)
	}
	e.cmd = cmd
	exitCh := make(chan struct{})
	e.exitCh = exitCh
	var waitErr error
	go func() {
		// waitErr is written here, then exitCh is closed in the same
		// goroutine: per Go's memory model that ordering is visible to any
		// goroutine that wakes up because the channel closed, so reading
		// waitErr below (after <-exitCh) needs no extra locking.
		waitErr = cmd.Wait()
		close(exitCh)
	}()

	// Guard against a bad config/binary causing sing-box to exit immediately;
	// without this check the engine would report a healthy connection even
	// though no proxy process is actually running.
	select {
	case <-exitCh:
		e.cmd = nil
		logBytes, _ := os.ReadFile(filepath.Join(config.GetVifyDir(), "singbox.log"))
		return fmt.Errorf("sing-box exited immediately (%v); log:\n%s", waitErr, strings.TrimSpace(string(logBytes)))
	case <-time.After(700 * time.Millisecond):
		// Still running after the startup window. A node can still pass this
		// check and then die later mid-session (e.g. a broken/CDN-fronted
		// node returning garbage on the first real request instead of a
		// proxy handshake). Watch for that in the background so a crash
		// triggers failover immediately instead of waiting on the slower
		// periodic watchdog probe (up to ~2x its interval before it reacts).
		go e.watchForCrash(exitCh, cmd)
	}

	// 3. Enable System Proxy if in system proxy mode
	if mode == model.ModeSystemProxy || mode == model.ModeMixed {
		_ = e.sysProxy.Enable()
	}

	// 4. Save session metadata
	session := &model.ActiveSession{
		PID:        os.Getpid(),
		Node:       *node,
		Mode:       mode,
		LocalSocks: e.cfg.LocalSocksPort,
		LocalHTTP:  e.cfg.LocalHTTPPort,
		StartedAt:  time.Now(),
	}
	_ = cache.SaveSession(session)

	// 5. Start Watchdog if enabled
	if e.cfg.AutoFailover {
		e.watchdog = NewWatchdog(e.cfg.WatchdogInterval, e.cfg.LocalSocksPort, e.cfg.TestURL)
		e.watchdog.OnFailover = func(reason string) {
			go e.FailoverToNext()
		}
		go e.watchdog.Start(e.ctx)
	}

	// 6. Start speed ticker routine
	sessionCtx := e.ctx // capture once: e.ctx gets reassigned by later StartNode
	// calls, and re-reading the struct field on every loop iteration below
	// (instead of this local copy) raced with those writes.
	go func() {
		ticker := time.NewTicker(1 * time.Second)
		defer ticker.Stop()
		httpClient := &http.Client{Timeout: 1500 * time.Millisecond}
		for {
			select {
			case <-sessionCtx.Done():
				return
			case <-ticker.C:
				if down, up, ok := fetchClashTraffic(httpClient); ok {
					e.stats.SetTotals(down, up)
				}
				e.stats.Tick()
			}
		}
	}()

	e.running = true
	return nil
}

// FailoverToNext seamlessly switches outbound node to the next available healthy node
func (e *Engine) FailoverToNext() {
	e.mu.Lock()

	if len(e.healthyNodes) == 0 {
		e.mu.Unlock()
		return
	}

	e.nodeIndex = (e.nodeIndex + 1) % len(e.healthyNodes)
	nextNode := e.healthyNodes[e.nodeIndex]
	nextMode := e.mode
	backups := e.healthyNodes

	_ = e.stopInternal()
	e.mu.Unlock()

	// StartNode acquires e.mu itself; must not still be held here, or this
	// deadlocks the engine (mu is not reentrant) every time a failover fires.
	_ = e.StartNode(context.Background(), nextNode, nextMode, backups)
}

// Stop cleanly terminates proxy process, restores system settings, and cleans session
func (e *Engine) Stop() error {
	e.mu.Lock()
	defer e.mu.Unlock()
	return e.stopInternal()
}

func (e *Engine) stopInternal() error {
	if e.cancel != nil {
		e.cancel()
	}

	if e.watchdog != nil {
		e.watchdog.Stop()
	}

	if e.sysProxy != nil {
		_ = e.sysProxy.Disable()
	}

	if e.cmd != nil && e.cmd.Process != nil {
		dying := e.cmd
		exitCh := e.exitCh
		// Clear before killing so watchForCrash (which checks e.cmd under
		// e.mu once exitCh closes) recognizes this as an intentional stop
		// rather than a crash to react to.
		e.cmd = nil
		e.exitCh = nil
		_ = dying.Process.Kill()
		if exitCh != nil {
			<-exitCh // wait for the reap to finish so a following StartNode doesn't race for the same ports
		}
	}

	cache.ClearSession()
	e.running = false
	return nil
}

// watchForCrash waits for a running sing-box process to exit and, if that
// exit was unexpected (nobody called Stop/FailoverToNext for it already),
// triggers an immediate failover instead of leaving the user disconnected
// until the next periodic watchdog probe.
func (e *Engine) watchForCrash(exitCh chan struct{}, myCmd *exec.Cmd) {
	<-exitCh

	e.mu.Lock()
	isStale := e.cmd != myCmd // already stopped/replaced intentionally
	if !isStale {
		e.cmd = nil
		e.exitCh = nil
		e.running = false
	}
	e.mu.Unlock()

	if isStale {
		return
	}

	go e.FailoverToNext()
}

// GetCurrentNode returns currently connected node
func (e *Engine) GetCurrentNode() *model.ProxyNode {
	e.mu.Lock()
	defer e.mu.Unlock()
	return e.currentNode
}

// GetStats returns current traffic statistics
func (e *Engine) GetStats() *TrafficStats {
	return e.stats
}

// fetchClashTraffic pulls cumulative upload/download totals from sing-box's
// Clash-compatible API. Returns ok=false while sing-box is still starting up
// or the API is unreachable, so callers should simply skip that tick.
func fetchClashTraffic(client *http.Client) (download, upload int64, ok bool) {
	resp, err := client.Get("http://" + clashAPIAddr + "/connections")
	if err != nil {
		return 0, 0, false
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, 0, false
	}

	var payload struct {
		DownloadTotal int64 `json:"downloadTotal"`
		UploadTotal   int64 `json:"uploadTotal"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return 0, 0, false
	}

	return payload.DownloadTotal, payload.UploadTotal, true
}
