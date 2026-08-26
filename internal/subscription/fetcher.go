package subscription

import (
	"bufio"
	"context"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/Mr-Meshky/vify-cli/internal/model"
)

// Fetcher handles downloading and decoding subscription links
type Fetcher struct {
	client *http.Client
}

// NewFetcher creates a new subscription Fetcher
func NewFetcher(timeout time.Duration) *Fetcher {
	if timeout <= 0 {
		timeout = 10 * time.Second
	}
	return &Fetcher{
		client: &http.Client{
			Timeout: timeout,
		},
	}
}

// FetchAll downloads configs from multiple URLs concurrently
func (f *Fetcher) FetchAll(ctx context.Context, urls []string) ([]*model.ProxyNode, error) {
	var (
		wg      sync.WaitGroup
		mu      sync.Mutex
		nodes   []*model.ProxyNode
		seenIDs = make(map[string]bool)
	)

	for _, u := range urls {
		u := strings.TrimSpace(u)
		if u == "" {
			continue
		}

		wg.Add(1)
		go func(urlStr string) {
			defer wg.Done()

			subNodes, err := f.FetchSingle(ctx, urlStr)
			if err != nil {
				return
			}

			mu.Lock()
			for _, n := range subNodes {
				if !seenIDs[n.ID] {
					seenIDs[n.ID] = true
					nodes = append(nodes, n)
				}
			}
			mu.Unlock()
		}(u)
	}

	wg.Wait()

	if len(nodes) == 0 {
		return nil, fmt.Errorf("no valid proxy configurations found from subscriptions")
	}

	return nodes, nil
}

// FetchSingle downloads or reads from a local file and parses all proxy links
func (f *Fetcher) FetchSingle(ctx context.Context, target string) ([]*model.ProxyNode, error) {
	var content string

	if strings.HasPrefix(target, "http://") || strings.HasPrefix(target, "https://") {
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, target, nil)
		if err != nil {
			return nil, err
		}
		req.Header.Set("User-Agent", "v2rayNG/1.8.5 (vify-cli)")

		resp, err := f.client.Do(req)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("http error %d fetching %s", resp.StatusCode, target)
		}

		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		content = string(bodyBytes)
	} else {
		// Read local file
		data, err := os.ReadFile(target)
		if err != nil {
			return nil, err
		}
		content = string(data)
	}

	return f.ParseContent(content)
}

// ParseContent splits raw content (plain lines or base64 encoded) and extracts nodes
func (f *Fetcher) ParseContent(content string) ([]*model.ProxyNode, error) {
	trimmed := strings.TrimSpace(content)
	if trimmed == "" {
		return nil, fmt.Errorf("empty content")
	}

	// If the content is base64 encoded as a whole
	if !strings.Contains(trimmed, "://") {
		if decoded, err := decodeSubscriptionBase64(trimmed); err == nil {
			trimmed = string(decoded)
		}
	}

	var nodes []*model.ProxyNode
	scanner := bufio.NewScanner(strings.NewReader(trimmed))
	// Support large lines
	buf := make([]byte, 1024*1024)
	scanner.Buffer(buf, 1024*1024)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") || strings.HasPrefix(line, "//") {
			continue
		}

		node, err := ParseURI(line)
		if err == nil && node != nil {
			nodes = append(nodes, node)
		}
	}

	return nodes, nil
}

func decodeSubscriptionBase64(s string) ([]byte, error) {
	s = strings.TrimSpace(s)
	s = strings.ReplaceAll(s, "\n", "")
	s = strings.ReplaceAll(s, "\r", "")
	if rem := len(s) % 4; rem != 0 {
		s += strings.Repeat("=", 4-rem)
	}

	b, err := base64.StdEncoding.DecodeString(s)
	if err == nil {
		return b, nil
	}
	return base64.URLEncoding.DecodeString(s)
}
