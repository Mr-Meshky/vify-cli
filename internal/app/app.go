package app

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/Mr-Meshky/vify-cli/internal/cache"
	"github.com/Mr-Meshky/vify-cli/internal/cleanip"
	"github.com/Mr-Meshky/vify-cli/internal/config"
	"github.com/Mr-Meshky/vify-cli/internal/core"
	"github.com/Mr-Meshky/vify-cli/internal/model"
	"github.com/Mr-Meshky/vify-cli/internal/ping"
	"github.com/Mr-Meshky/vify-cli/internal/subscription"
	"github.com/Mr-Meshky/vify-cli/internal/tui"
	tea "github.com/charmbracelet/bubbletea"
)

// App manages command workflows
type App struct {
	Config  *config.AppConfig
	Fetcher *subscription.Fetcher
	Engine  *core.Engine
}

// New creates a new initialized App
func New() (*App, error) {
	cfg, err := config.LoadConfig()
	if err != nil {
		return nil, err
	}

	return &App{
		Config:  cfg,
		Fetcher: subscription.NewFetcher(12 * time.Second),
		Engine:  core.NewEngine(cfg),
	}, nil
}

// ConnectOptions contains flags for the connect command
type ConnectOptions struct {
	Country         string
	Protocol        string
	Mode            string
	FastPass        bool
	BatchSize       int
	UseCacheOnly    bool
	ManualURI       string
	PreselectedNode *model.ProxyNode
}

// RunConnect handles automatic or filtered fast-pass connection or manual config link
func (a *App) RunConnect(ctx context.Context, opts ConnectOptions) error {
	fmt.Println(tui.TitleStyle.Render("⚡ Vify CLI — High-Speed Anti-Censorship VPN (TUN Mode)"))
	fmt.Println(tui.SubtitleStyle.Render("Initializing VPN connection..."))
	fmt.Println()

	if err := validateProtocol(opts.Protocol); err != nil {
		return err
	}

	// Determine Connection Mode (Default to TUN Mode)
	connMode := model.ModeTUN
	if opts.Mode == "system_proxy" {
		connMode = model.ModeSystemProxy
	}

	// 1. Manual Config Link
	if opts.ManualURI != "" {
		fmt.Printf(" %s Parsing manual configuration link...\n", tui.BadgeStyle.Render("CONFIG"))
		node, err := subscription.ParseURI(opts.ManualURI)
		if err != nil {
			return fmt.Errorf("failed to parse config link: %w", err)
		}
		return a.connectToNode(ctx, node, connMode)
	}

	// 2. A specific node the user already picked (e.g. via `vify list`) —
	// connect straight to it instead of re-running discovery/benchmarking.
	if opts.PreselectedNode != nil {
		return a.connectToNode(ctx, opts.PreselectedNode, connMode)
	}

	var candidateNodes []*model.ProxyNode

	if !opts.UseCacheOnly {
		fmt.Printf(" %s Fetching latest subscriptions...\n", tui.BadgeStyle.Render("FETCH"))
		nodes, err := a.Fetcher.FetchAll(ctx, a.Config.Subscriptions)
		if err == nil && len(nodes) > 0 {
			candidateNodes = nodes
		}
	}

	if len(candidateNodes) == 0 {
		cached, _ := cache.LoadCache()
		if len(cached.Nodes) > 0 {
			fmt.Printf(" %s Using %d cached healthy servers.\n", tui.WarningBadge.Render("CACHE"), len(cached.Nodes))
			candidateNodes = cached.Nodes
		}
	}

	if len(candidateNodes) == 0 {
		return fmt.Errorf("no server configurations available")
	}

	// Filter by Country and Protocol if requested
	candidateNodes = filterNodes(candidateNodes, opts.Country, opts.Protocol)
	if len(candidateNodes) == 0 {
		return fmt.Errorf("no servers matched filters (Country: %s, Protocol: %s)", opts.Country, opts.Protocol)
	}

	// Limit batch size
	if opts.BatchSize > 0 && len(candidateNodes) > opts.BatchSize {
		candidateNodes = candidateNodes[:opts.BatchSize]
	}

	fmt.Printf(" %s Benchmarking %d candidates concurrently (Fast-Pass < %dms)...\n",
		tui.BadgeStyle.Render("TEST"),
		len(candidateNodes),
		a.Config.FastPassThresholdMS,
	)

	tester := ping.NewTester(
		a.Config.TestTimeoutMS,
		a.Config.ConcurrencyLimit,
		a.Config.FastPassThresholdMS,
		a.Config.TestURL,
	)

	fastPassChan := make(chan *model.ProxyNode, 1)
	if opts.FastPass {
		tester.OnFastPassCandidate = func(node *model.ProxyNode) {
			select {
			case fastPassChan <- node:
			default:
			}
		}
	}

	tester.OnProgress = func(result model.BenchmarkResult, completed, total int) {
		if completed%10 == 0 || completed == total {
			fmt.Printf("\r  ⚡ Progress: %d/%d verified...", completed, total)
		}
	}

	var healthyNodes []*model.ProxyNode
	testDone := make(chan struct{})

	go func() {
		healthyNodes = tester.TestBatch(ctx, candidateNodes)
		close(testDone)
	}()

	var targetNode *model.ProxyNode

	select {
	case fastNode := <-fastPassChan:
		targetNode = fastNode
		fmt.Printf("\n %s Fast-Pass discovered: %s %s (%dms)\n",
			tui.SuccessBadge.Render("FAST-PASS"),
			targetNode.CountryFlag,
			targetNode.Name,
			targetNode.Latency.Milliseconds(),
		)
	case <-testDone:
		if len(healthyNodes) > 0 {
			targetNode = healthyNodes[0]
		}
	}

	if targetNode == nil {
		<-testDone
		if len(healthyNodes) > 0 {
			targetNode = healthyNodes[0]
		} else {
			return fmt.Errorf("no healthy servers found. Try 'vify clean-ip' or check your connection")
		}
	}

	// Persist tested nodes to cache in background
	go func() {
		<-testDone
		if len(healthyNodes) > 0 {
			_ = cache.SaveCache(healthyNodes)
		}
	}()

	// The benchmark above (and Fast-Pass in particular) only confirms a raw
	// TCP/TLS dial succeeds — it does NOT confirm the proxy protocol itself
	// relays real traffic. Verify the chosen node with an actual request,
	// falling back through the next-best benchmarked candidates if it turns
	// out to be broken, instead of silently connecting a session that can't
	// actually reach anything.
	fmt.Printf("\n %s Verifying real proxy connectivity for %s %s...\n",
		tui.BadgeStyle.Render("VERIFY"), targetNode.CountryFlag, targetNode.Name)

	verifiedNode, verifiedLat, verr := verifyRealConnectivity(ctx, targetNode, func() []*model.ProxyNode {
		<-testDone
		return healthyNodes
	}, a.Config.TestURL, 16, 4, func(attempt int, n *model.ProxyNode, err error) {
		if err != nil {
			fmt.Printf("  [%d] %s %s failed: %v\n", attempt, n.CountryFlag, n.Name, err)
		}
	})
	if verr != nil {
		return fmt.Errorf("no server passed real connectivity verification: %w", verr)
	}
	if verifiedNode != targetNode {
		fmt.Printf(" %s Original pick failed the real check — switched to %s %s\n",
			tui.WarningBadge.Render("SWITCHED"), verifiedNode.CountryFlag, verifiedNode.Name)
	}
	targetNode = verifiedNode
	fmt.Printf(" %s Verified working (real latency: %s)\n",
		tui.SuccessBadge.Render("VERIFIED"), tui.FormatLatency(verifiedLat.Milliseconds()))

	// Determine Connection Mode
	if opts.Mode == "system_proxy" {
		connMode = model.ModeSystemProxy
	} else {
		connMode = model.ModeTUN
	}

	fmt.Printf("\n %s Starting connection via %s mode...\n",
		tui.BadgeStyle.Render("CONNECT"),
		strings.ToUpper(string(connMode)),
	)

	if err := a.Engine.StartNode(ctx, targetNode, connMode, healthyNodes); err != nil {
		return fmt.Errorf("failed to start proxy session: %w", err)
	}

	return a.startSessionUI(targetNode, connMode)
}

// verifyRealConnectivity confirms a candidate node actually relays traffic by
// trying `first` immediately (the fast path — no extra latency paid), then
// only if that fails, calling `more` (typically waiting for the full
// benchmark batch to finish) and trying its next-best candidates, up to
// maxAttempts total, `concurrency` at a time. Free public subscriptions can
// have a very low rate of nodes that pass real verification (raw TCP/TLS
// reachability is a weak signal — see VerifyNodeReal), so attempts run
// concurrently to keep total wait time reasonable even when trying a dozen
// or more candidates. Returns as soon as one passes; the rest are cancelled.
// onAttempt (optional) is called for every finished attempt (in completion
// order, not candidate order) so callers can surface why each one failed.
func verifyRealConnectivity(ctx context.Context, first *model.ProxyNode, more func() []*model.ProxyNode, testURL string, maxAttempts, concurrency int, onAttempt func(attempt int, n *model.ProxyNode, err error)) (*model.ProxyNode, time.Duration, error) {
	nodeKey := func(n *model.ProxyNode) string {
		return fmt.Sprintf("%s:%d:%s", n.Server, n.Port, n.Protocol)
	}

	lat, err := core.VerifyNodeReal(ctx, first, testURL, core.VerifyPortBase, core.VerifyPortBase+1)
	if onAttempt != nil {
		onAttempt(1, first, err)
	}
	if err == nil {
		return first, lat, nil
	}

	// Build the fallback queue: deduped, capped at maxAttempts-1 (first already used one slot).
	tried := map[string]bool{nodeKey(first): true}
	var queue []*model.ProxyNode
	for _, n := range more() {
		if len(queue) >= maxAttempts-1 {
			break
		}
		if tried[nodeKey(n)] {
			continue
		}
		tried[nodeKey(n)] = true
		queue = append(queue, n)
	}

	type result struct {
		node *model.ProxyNode
		lat  time.Duration
		err  error
	}

	verifyCtx, cancel := context.WithCancel(ctx)
	defer cancel()

	resultCh := make(chan result, len(queue))
	sem := make(chan struct{}, concurrency)
	var wg sync.WaitGroup

	for i, n := range queue {
		wg.Add(1)
		go func(slot int, node *model.ProxyNode) {
			defer wg.Done()
			sem <- struct{}{}
			defer func() { <-sem }()

			socksPort := core.VerifyPortBase + 2 + slot*2
			httpPort := socksPort + 1
			lat, err := core.VerifyNodeReal(verifyCtx, node, testURL, socksPort, httpPort)
			resultCh <- result{node, lat, err}
		}(i, n)
	}
	go func() {
		wg.Wait()
		close(resultCh)
	}()

	lastErr := err
	attempts := 1
	for r := range resultCh {
		attempts++
		if onAttempt != nil {
			onAttempt(attempts, r.node, r.err)
		}
		if r.err == nil {
			cancel() // stop any still-running verification attempts
			return r.node, r.lat, nil
		}
		lastErr = r.err
	}

	return nil, 0, fmt.Errorf("tried %d candidates, none passed real verification (last error: %v)", attempts, lastErr)
}

// connectToNode dials, verifies, and activates a session for a single,
// already-known node — shared by the manual-URI path and by connecting
// directly to a server the user picked in `vify list`.
func (a *App) connectToNode(ctx context.Context, node *model.ProxyNode, connMode model.ConnectionMode) error {
	fmt.Printf(" %s Testing connection to %s %s (%s:%d)...\n",
		tui.BadgeStyle.Render("TEST"),
		node.CountryFlag,
		node.Name,
		node.Server,
		node.Port,
	)

	testCtx, cancel := context.WithTimeout(ctx, 3*time.Second)
	lat, err := ping.DialTest(testCtx, node, a.Config.TestURL)
	cancel()

	if err == nil && lat > 0 {
		node.Latency = lat
		node.IsHealthy = true
		fmt.Printf(" %s Node reachable (latency: %s)\n", tui.SuccessBadge.Render("ONLINE"), tui.FormatLatency(lat.Milliseconds()))
	} else {
		fmt.Printf(" %s Node test warning: %v (attempting connection anyway)\n", tui.WarningBadge.Render("WARN"), err)
	}

	// The check above only confirms a raw TCP/TLS dial succeeds — it does
	// NOT confirm the proxy protocol itself works. Actually exercise it
	// with a real request before committing to a full session.
	fmt.Printf(" %s Verifying real proxy handshake...\n", tui.BadgeStyle.Render("VERIFY"))
	if realLat, verr := core.VerifyNodeReal(ctx, node, a.Config.TestURL, core.VerifyPortBase, core.VerifyPortBase+1); verr != nil {
		fmt.Printf(" %s Real proxy check failed: %v (attempting connection anyway — this node may not work)\n", tui.WarningBadge.Render("WARN"), verr)
	} else {
		fmt.Printf(" %s Verified working (real latency: %s)\n", tui.SuccessBadge.Render("VERIFIED"), tui.FormatLatency(realLat.Milliseconds()))
	}

	fmt.Printf("\n %s Activating %s mode...\n",
		tui.BadgeStyle.Render("CONNECT"),
		strings.ToUpper(string(connMode)),
	)

	if err := a.Engine.StartNode(ctx, node, connMode, []*model.ProxyNode{node}); err != nil {
		return fmt.Errorf("failed to start VPN: %w", err)
	}

	return a.startSessionUI(node, connMode)
}

func (a *App) startSessionUI(targetNode *model.ProxyNode, connMode model.ConnectionMode) error {
	// Set up OS signal interceptor
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	// Render interactive status dashboard
	session := &model.ActiveSession{
		PID:        os.Getpid(),
		Node:       *targetNode,
		Mode:       connMode,
		LocalSocks: a.Config.LocalSocksPort,
		LocalHTTP:  a.Config.LocalHTTPPort,
		StartedAt:  time.Now(),
	}

	statusModel := tui.NewStatusModel(session, a.Engine.GetStats())
	p := tea.NewProgram(statusModel)

	go func() {
		<-sigChan
		_ = a.Engine.Stop()
		p.Quit()
	}()

	if _, err := p.Run(); err != nil {
		_ = a.Engine.Stop()
		return err
	}

	_ = a.Engine.Stop()
	fmt.Println(tui.SuccessBadge.Render("✓") + " Disconnected cleanly. System proxy and DNS restored.")
	return nil
}

// RunList displays the interactive server selection TUI
func (a *App) RunList(ctx context.Context, country, protocol string, systemProxy bool) error {
	if err := validateProtocol(protocol); err != nil {
		return err
	}

	fmt.Println(tui.TitleStyle.Render("⚡ Fetching available servers for selection..."))

	nodes, err := a.Fetcher.FetchAll(ctx, a.Config.Subscriptions)
	if err != nil || len(nodes) == 0 {
		cached, _ := cache.LoadCache()
		nodes = cached.Nodes
	}

	nodes = filterNodes(nodes, country, protocol)
	if len(nodes) == 0 {
		return fmt.Errorf("no servers found matching criteria")
	}

	listModel := tui.NewListModel(nodes)
	p := tea.NewProgram(listModel)
	m, err := p.Run()
	if err != nil {
		return err
	}

	finalModel := m.(tui.ListModel)
	if finalModel.Selected == nil {
		return nil
	}

	mode := "tun"
	if systemProxy {
		mode = "system_proxy"
	}

	return a.RunConnect(ctx, ConnectOptions{
		Mode:            mode,
		PreselectedNode: finalModel.Selected,
	})
}

// RunTest benchmarks servers and prints the leaderboard
func (a *App) RunTest(ctx context.Context, batchSize int, country, protocol string) error {
	if err := validateProtocol(protocol); err != nil {
		return err
	}

	fmt.Println(tui.TitleStyle.Render("⚡ Vify Benchmark & Latency Tester"))
	fmt.Printf(" %s Fetching candidate servers...\n", tui.BadgeStyle.Render("FETCH"))

	nodes, err := a.Fetcher.FetchAll(ctx, a.Config.Subscriptions)
	if err != nil || len(nodes) == 0 {
		cached, _ := cache.LoadCache()
		nodes = cached.Nodes
	}

	nodes = filterNodes(nodes, country, protocol)
	if batchSize > 0 && len(nodes) > batchSize {
		nodes = nodes[:batchSize]
	}

	fmt.Printf(" %s Testing %d servers with real HTTP/204 get (concurrency: %d)...\n",
		tui.BadgeStyle.Render("TEST"),
		len(nodes),
		a.Config.ConcurrencyLimit,
	)

	tester := ping.NewTester(a.Config.TestTimeoutMS, a.Config.ConcurrencyLimit, a.Config.FastPassThresholdMS, a.Config.TestURL)
	tester.OnProgress = func(result model.BenchmarkResult, completed, total int) {
		if completed%5 == 0 || completed == total {
			fmt.Printf("\r  ⚡ Benchmarking: %d/%d completed...", completed, total)
		}
	}

	results := tester.TestBatch(ctx, nodes)
	fmt.Println()

	if len(results) > 0 {
		_ = cache.SaveCache(results)
	}

	fmt.Print(tui.RenderLeaderboard(results))
	return nil
}

// RunCleanIP scans Cloudflare clean IPs
func (a *App) RunCleanIP(ctx context.Context, count int) error {
	fmt.Println(tui.TitleStyle.Render("⚡ Cloudflare Clean IP Scanner for Iranian ISPs"))
	fmt.Printf(" %s Scanning candidate IP ranges for lowest latency and zero packet loss...\n\n", tui.BadgeStyle.Render("SCAN"))

	candidates := append([]string{}, cleanip.PopularCandidateIPs...)
	candidates = append(candidates, cleanip.RandomIPsFromCIDRs(cleanip.CloudflareCIDRs, 40)...)

	scanner := cleanip.NewScanner(30, 1500*time.Millisecond)
	results := scanner.ScanIPs(ctx, candidates, func(res cleanip.CleanIPResult, done, total int) {
		if done%5 == 0 || done == total {
			fmt.Printf("\r  ⚡ Scanned: %d/%d IPs...", done, total)
		}
	})

	fmt.Println()
	if len(results) == 0 {
		fmt.Println(tui.DangerBadge.Render(" No clean IPs reachable on current network interface. "))
		return nil
	}

	fmt.Println(tui.HeaderStyle.Render("🎯 Top Clean Cloudflare IPs (Optimized for your ISP)"))
	fmt.Printf(" %-4s %-20s %-12s %-12s\n", "#", "IP ADDRESS", "PORT", "LATENCY")
	fmt.Println(strings.Repeat("─", 50))

	limit := count
	if limit <= 0 || limit > len(results) {
		limit = len(results)
	}

	for i := 0; i < limit; i++ {
		res := results[i]
		fmt.Printf(" #%-3d %-20s %-12d %s\n",
			i+1,
			res.IP,
			res.Port,
			tui.FormatLatency(res.Latency.Milliseconds()),
		)
	}
	fmt.Println(strings.Repeat("─", 50))
	return nil
}

// RunStatus prints current session status
func (a *App) RunStatus() error {
	session, err := cache.LoadSession()
	if err != nil || session == nil {
		fmt.Println(tui.WarningBadge.Render(" No active Vify connection found. "))
		return nil
	}

	fmt.Println(tui.HeaderStyle.Render("⚡ Vify Active Connection Status"))
	fmt.Printf(" Status:     %s\n", tui.SuccessBadge.Render("ONLINE"))
	fmt.Printf(" Mode:       %s\n", tui.BadgeStyle.Render(strings.ToUpper(string(session.Mode))))
	fmt.Printf(" Server:     %s %s\n", session.Node.CountryFlag, session.Node.Name)
	fmt.Printf(" Endpoint:   %s:%d\n", session.Node.Server, session.Node.Port)
	fmt.Printf(" Protocol:   %s\n", session.Node.Protocol)
	fmt.Printf(" Local SOCKS:127.0.0.1:%d\n", session.LocalSocks)
	fmt.Printf(" Local HTTP: 127.0.0.1:%d\n", session.LocalHTTP)
	fmt.Printf(" Started At: %s (%s ago)\n",
		session.StartedAt.Format("15:04:05"),
		time.Since(session.StartedAt).Truncate(time.Second),
	)
	return nil
}

// RunDisconnect stops any active session and restores proxy/DNS
func (a *App) RunDisconnect() error {
	session, err := cache.LoadSession()
	if err == nil && session != nil {
		// Stop system proxy
		sys := core.NewSystemProxyController(session.LocalHTTP, session.LocalSocks)
		_ = sys.Disable()

		// Kill process if still alive
		if proc, err := os.FindProcess(session.PID); err == nil {
			_ = proc.Kill()
		}
		cache.ClearSession()
	}

	// Also ensure local system proxy is disabled
	sys := core.NewSystemProxyController(a.Config.LocalHTTPPort, a.Config.LocalSocksPort)
	_ = sys.Disable()

	fmt.Println(tui.SuccessBadge.Render("✓") + " Vify disconnected and all system proxy/DNS settings successfully restored.")
	return nil
}

// validProtocolFlags lists every value the --protocol flag accepts, including
// the "ss" shorthand alias for shadowsocks handled in filterNodes.
var validProtocolFlags = []string{"vless", "vmess", "trojan", "shadowsocks", "ss"}

// validateProtocol rejects an unrecognized --protocol value up front with a
// helpful message, instead of letting it silently filter out every server
// and surface as a generic "no servers matched" error further down.
func validateProtocol(protocol string) error {
	protocol = strings.ToLower(strings.TrimSpace(protocol))
	if protocol == "" {
		return nil
	}
	for _, p := range validProtocolFlags {
		if protocol == p {
			return nil
		}
	}
	return fmt.Errorf("unknown protocol %q — valid options are: %s", protocol, strings.Join(validProtocolFlags, ", "))
}

func filterNodes(nodes []*model.ProxyNode, country, protocol string) []*model.ProxyNode {
	if country == "" && protocol == "" {
		return nodes
	}

	country = strings.ToUpper(strings.TrimSpace(country))
	protocol = strings.ToLower(strings.TrimSpace(protocol))
	if protocol == "ss" {
		protocol = string(model.ProtocolShadowsocks)
	}

	var filtered []*model.ProxyNode
	for _, n := range nodes {
		if country != "" && n.CountryCode != country {
			continue
		}
		if protocol != "" && strings.ToLower(string(n.Protocol)) != protocol {
			continue
		}
		filtered = append(filtered, n)
	}
	return filtered
}
