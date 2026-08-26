package cleanip

import (
	"context"
	"crypto/tls"
	"fmt"
	"net"
	"sort"
	"sync"
	"time"
)

// CleanIPResult stores the scan result for a Cloudflare IP
type CleanIPResult struct {
	IP         string        `json:"ip"`
	Port       int           `json:"port"`
	Latency    time.Duration `json:"latency"`
	Success    bool          `json:"success"`
	PacketLoss float64       `json:"packet_loss"`
}

// Scanner scans candidate Cloudflare IPs concurrently
type Scanner struct {
	Timeout     time.Duration
	Concurrency int
	SNI         string
	Port        int
}

// NewScanner creates a new Cloudflare Clean IP scanner
func NewScanner(concurrency int, timeout time.Duration) *Scanner {
	if concurrency <= 0 {
		concurrency = 30
	}
	if timeout <= 0 {
		timeout = 1500 * time.Millisecond
	}
	return &Scanner{
		Timeout:     timeout,
		Concurrency: concurrency,
		SNI:         "cloudflare.com",
		Port:        443,
	}
}

// ScanIPs tests candidate IPs concurrently and returns them sorted by lowest latency
func (s *Scanner) ScanIPs(ctx context.Context, ipList []string, onProgress func(res CleanIPResult, done, total int)) []CleanIPResult {
	if len(ipList) == 0 {
		ipList = PopularCandidateIPs
	}

	total := len(ipList)
	sem := make(chan struct{}, s.Concurrency)
	var (
		wg      sync.WaitGroup
		mu      sync.Mutex
		results []CleanIPResult
		done    int
	)

	for _, ipStr := range ipList {
		ip := ipStr
		wg.Add(1)

		go func() {
			defer wg.Done()
			sem <- struct{}{}
			defer func() { <-sem }()

			res := s.testSingleIP(ctx, ip)

			mu.Lock()
			done++
			if res.Success {
				results = append(results, res)
			}
			currentDone := done
			mu.Unlock()

			if onProgress != nil {
				onProgress(res, currentDone, total)
			}
		}()
	}

	wg.Wait()

	sort.Slice(results, func(i, j int) bool {
		return results[i].Latency < results[j].Latency
	})

	return results
}

func (s *Scanner) testSingleIP(ctx context.Context, ip string) CleanIPResult {
	target := fmt.Sprintf("%s:%d", ip, s.Port)
	start := time.Now()

	d := &net.Dialer{Timeout: s.Timeout}
	conn, err := d.DialContext(ctx, "tcp", target)
	if err != nil {
		return CleanIPResult{
			IP:      ip,
			Port:    s.Port,
			Success: false,
		}
	}
	defer conn.Close()

	tlsConn := tls.Client(conn, &tls.Config{
		ServerName:         s.SNI,
		InsecureSkipVerify: true,
	})
	tlsConn.SetDeadline(time.Now().Add(s.Timeout))

	if err := tlsConn.HandshakeContext(ctx); err != nil {
		return CleanIPResult{
			IP:      ip,
			Port:    s.Port,
			Success: false,
		}
	}

	latency := time.Since(start)
	return CleanIPResult{
		IP:         ip,
		Port:       s.Port,
		Latency:    latency,
		Success:    true,
		PacketLoss: 0,
	}
}
