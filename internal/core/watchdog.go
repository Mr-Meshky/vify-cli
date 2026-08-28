package core

import (
	"context"
	"fmt"
	"time"

	"github.com/Mr-Meshky/vify-cli/internal/ping"
)

// Watchdog monitors the active proxy connection and triggers failover when disconnected
type Watchdog struct {
	Interval       time.Duration
	TestURL        string
	ProxyURL       string
	FailThreshold  int
	OnFailover     func(reason string)
	OnHealthUpdate func(latency time.Duration, ok bool)
	stopChan       chan struct{}
}

// NewWatchdog creates a new connection Watchdog
func NewWatchdog(intervalSec int, socksPort int, testURL string) *Watchdog {
	if intervalSec <= 0 {
		intervalSec = 7
	}
	if testURL == "" {
		testURL = "http://cp.cloudflare.com/generate_204"
	}

	return &Watchdog{
		Interval:      time.Duration(intervalSec) * time.Second,
		TestURL:       testURL,
		ProxyURL:      fmt.Sprintf("socks5://127.0.0.1:%d", socksPort),
		FailThreshold: 2,
		stopChan:      make(chan struct{}),
	}
}

// Start begins background health monitoring
func (w *Watchdog) Start(ctx context.Context) {
	ticker := time.NewTicker(w.Interval)
	defer ticker.Stop()

	consecutiveFails := 0

	for {
		select {
		case <-ctx.Done():
			return
		case <-w.stopChan:
			return
		case <-ticker.C:
			testCtx, cancel := context.WithTimeout(ctx, 3*time.Second)
			latency, code, err := ping.HTTP204ViaLocalProxy(testCtx, w.ProxyURL, w.TestURL, 3*time.Second)
			cancel()

			isHealthy := (err == nil && (code == 204 || code == 200))

			if isHealthy {
				consecutiveFails = 0
				if w.OnHealthUpdate != nil {
					w.OnHealthUpdate(latency, true)
				}
			} else {
				consecutiveFails++
				if w.OnHealthUpdate != nil {
					w.OnHealthUpdate(0, false)
				}

				if consecutiveFails >= w.FailThreshold {
					consecutiveFails = 0
					if w.OnFailover != nil {
						reason := "HTTP 204 check failed consecutively"
						if err != nil {
							reason = err.Error()
						}
						w.OnFailover(reason)
					}
				}
			}
		}
	}
}

// Stop terminates the watchdog loop
func (w *Watchdog) Stop() {
	select {
	case <-w.stopChan:
	default:
		close(w.stopChan)
	}
}
