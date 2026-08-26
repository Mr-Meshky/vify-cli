package core

import (
	"fmt"
	"sync"
	"time"
)

// TrafficStats tracks bandwidth and upload/download rates
type TrafficStats struct {
	mu             sync.RWMutex
	BytesUpload    int64
	BytesDownload  int64
	UploadSpeed    int64 // bytes per second
	DownloadSpeed  int64 // bytes per second
	lastUpload     int64
	lastDownload   int64
	lastCalculated time.Time
}

// NewTrafficStats initializes a new traffic monitor
func NewTrafficStats() *TrafficStats {
	return &TrafficStats{
		lastCalculated: time.Now(),
	}
}

// AddUpload records uploaded bytes
func (s *TrafficStats) AddUpload(n int64) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.BytesUpload += n
}

// AddDownload records downloaded bytes
func (s *TrafficStats) AddDownload(n int64) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.BytesDownload += n
}

// SetTotals overwrites cumulative byte counters with absolute totals
// reported by an external source (e.g. sing-box's Clash API).
func (s *TrafficStats) SetTotals(download, upload int64) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.BytesDownload = download
	s.BytesUpload = upload
}

// Tick calculates current speed rates (should be called once per second)
func (s *TrafficStats) Tick() {
	s.mu.Lock()
	defer s.mu.Unlock()

	now := time.Now()
	duration := now.Sub(s.lastCalculated).Seconds()
	if duration <= 0 {
		duration = 1
	}

	diffUp := s.BytesUpload - s.lastUpload
	diffDown := s.BytesDownload - s.lastDownload

	s.UploadSpeed = int64(float64(diffUp) / duration)
	s.DownloadSpeed = int64(float64(diffDown) / duration)

	s.lastUpload = s.BytesUpload
	s.lastDownload = s.BytesDownload
	s.lastCalculated = now
}

// Snapshot returns a copy of current metrics
func (s *TrafficStats) Snapshot() (uploadTotal, downloadTotal, upSpeed, downSpeed int64) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.BytesUpload, s.BytesDownload, s.UploadSpeed, s.DownloadSpeed
}

// FormatBytes formats byte counts into human-readable strings (KB, MB, GB)
func FormatBytes(b int64) string {
	const unit = 1024
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := int64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.2f %cB", float64(b)/float64(div), "KMGTPE"[exp])
}

// FormatSpeed formats speed bytes/sec into human-readable string
func FormatSpeed(bytesPerSec int64) string {
	return fmt.Sprintf("%s/s", FormatBytes(bytesPerSec))
}
