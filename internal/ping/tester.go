package ping

import (
	"context"
	"sort"
	"sync"
	"time"

	"github.com/Mr-Meshky/vify-cli/internal/model"
)

// ProgressCallback is invoked as each node finishes testing
type ProgressCallback func(result model.BenchmarkResult, completed, total int)

// Tester manages high-concurrency node benchmarking
type Tester struct {
	Timeout             time.Duration
	Concurrency         int
	FastPassThreshold   time.Duration
	TestURL             string
	OnProgress          ProgressCallback
	OnFastPassCandidate func(node *model.ProxyNode)
}

// NewTester creates a new Tester with given settings
func NewTester(timeoutMS, concurrency, fastPassMS int, testURL string) *Tester {
	if timeoutMS <= 0 {
		timeoutMS = 2500
	}
	if concurrency <= 0 {
		concurrency = 30
	}
	if fastPassMS <= 0 {
		fastPassMS = 800
	}
	if testURL == "" {
		testURL = "http://cp.cloudflare.com/generate_204"
	}

	return &Tester{
		Timeout:           time.Duration(timeoutMS) * time.Millisecond,
		Concurrency:       concurrency,
		FastPassThreshold: time.Duration(fastPassMS) * time.Millisecond,
		TestURL:           testURL,
	}
}

// TestBatch benchmarks a list of nodes concurrently and returns healthy nodes sorted by lowest latency
func (t *Tester) TestBatch(ctx context.Context, nodes []*model.ProxyNode) []*model.ProxyNode {
	total := len(nodes)
	if total == 0 {
		return nil
	}

	semaphore := make(chan struct{}, t.Concurrency)
	var (
		wg           sync.WaitGroup
		mu           sync.Mutex
		completed    int
		results      []*model.ProxyNode
		fastPassOnce sync.Once
	)

	for _, n := range nodes {
		node := n
		wg.Add(1)

		go func() {
			defer wg.Done()

			semaphore <- struct{}{}
			defer func() { <-semaphore }()

			testCtx, cancel := context.WithTimeout(ctx, t.Timeout)
			defer cancel()

			lat, err := DialTest(testCtx, node, t.TestURL)

			mu.Lock()
			completed++
			res := model.BenchmarkResult{
				Node:    node,
				Latency: lat,
			}

			if err == nil && lat > 0 {
				node.Latency = lat
				node.IsHealthy = true
				node.LastTested = time.Now()
				res.Success = true
				res.StatusCode = 204
				results = append(results, node)

				// Fast-Pass Trigger
				if lat <= t.FastPassThreshold && t.OnFastPassCandidate != nil {
					fastPassOnce.Do(func() {
						t.OnFastPassCandidate(node)
					})
				}
			} else {
				node.IsHealthy = false
				node.Latency = 0
				res.Success = false
				if err != nil {
					res.Error = err.Error()
				}
			}

			currentCompleted := completed
			progressCb := t.OnProgress
			mu.Unlock()

			if progressCb != nil {
				progressCb(res, currentCompleted, total)
			}
		}()
	}

	wg.Wait()

	// Sort healthy nodes by latency ascending
	sort.Slice(results, func(i, j int) bool {
		return results[i].Latency < results[j].Latency
	})

	return results
}
