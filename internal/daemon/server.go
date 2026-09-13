package daemon

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/Mr-Meshky/vify-cli/internal/app"
	"github.com/Mr-Meshky/vify-cli/internal/cache"
	"github.com/Mr-Meshky/vify-cli/internal/model"
)

// DaemonStatus represents the current state of the Vify connection
type DaemonStatus struct {
	Connected     bool            `json:"connected"`
	Connecting    bool            `json:"connecting"`
	StatusMessage string          `json:"status_message"`
	Mode          string          `json:"mode"`
	CurrentNode   *model.ProxyNode `json:"current_node,omitempty"`
	UploadSpeed   int64           `json:"upload_speed"`
	DownloadSpeed int64           `json:"download_speed"`
	TotalUpload   int64           `json:"total_upload"`
	TotalDownload int64           `json:"total_download"`
	LastError     string          `json:"last_error,omitempty"`
}

// ConnectRequest is the JSON payload for /api/connect
type ConnectRequest struct {
	Target       string `json:"target"`
	Country      string `json:"country"`
	Protocol     string `json:"protocol"`
	Mode         string `json:"mode"`
	FastPass     bool   `json:"fast_pass"`
	BatchSize    int    `json:"batch_size"`
	UseCacheOnly bool   `json:"cache_only"`
	RawURL       string `json:"raw_url"`
}

// Server handles the HTTP daemon API for Vify Desktop / GUI clients
type Server struct {
	mu         sync.RWMutex
	app        *app.App
	status     DaemonStatus
	cancelConn context.CancelFunc
	clients    map[chan string]bool
	clientsMu  sync.Mutex
	httpServer *http.Server
}

// NewServer creates a new Daemon server
func NewServer(application *app.App) *Server {
	return &Server{
		app: application,
		status: DaemonStatus{
			Connected:     false,
			Connecting:    false,
			StatusMessage: "Disconnected",
			Mode:          "tun",
		},
		clients: make(map[chan string]bool),
	}
}

// Start runs the HTTP listener on the given port
func (s *Server) Start(port int) error {
	mux := http.NewServeMux()

	mux.HandleFunc("/api/status", s.handleStatus)
	mux.HandleFunc("/api/connect", s.handleConnect)
	mux.HandleFunc("/api/disconnect", s.handleDisconnect)
	mux.HandleFunc("/api/nodes", s.handleNodes)
	mux.HandleFunc("/api/events", s.handleEvents)
	mux.HandleFunc("/api/ping", func(w http.ResponseWriter, r *http.Request) {
		enableCORS(w)
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]string{"status": "ok", "version": "1.0.0"})
	})

	s.httpServer = &http.Server{
		Addr:    fmt.Sprintf("127.0.0.1:%d", port),
		Handler: mux,
	}

	// Start background ticker for traffic stats & events
	go s.statsLoop()

	return s.httpServer.ListenAndServe()
}

// Stop gracefully terminates the daemon and active VPN session
func (s *Server) Stop(ctx context.Context) error {
	s.mu.Lock()
	if s.cancelConn != nil {
		s.cancelConn()
	}
	if s.app.Engine != nil {
		_ = s.app.Engine.Stop()
	}
	s.mu.Unlock()

	if s.httpServer != nil {
		return s.httpServer.Shutdown(ctx)
	}
	return nil
}

func enableCORS(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
}

func (s *Server) handleStatus(w http.ResponseWriter, r *http.Request) {
	enableCORS(w)
	if r.Method == http.MethodOptions {
		return
	}

	s.mu.RLock()
	st := s.status
	s.mu.RUnlock()

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(st)
}

func (s *Server) handleConnect(w http.ResponseWriter, r *http.Request) {
	enableCORS(w)
	if r.Method == http.MethodOptions {
		return
	}
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req ConnectRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		req = ConnectRequest{FastPass: true, Mode: "tun", BatchSize: 60}
	}

	if req.Mode == "" {
		req.Mode = "tun"
	}
	if req.BatchSize <= 0 {
		req.BatchSize = 60
	}

	s.mu.Lock()
	if s.status.Connecting || s.status.Connected {
		if s.cancelConn != nil {
			s.cancelConn()
		}
		if s.app.Engine != nil {
			_ = s.app.Engine.Stop()
		}
	}

	ctx, cancel := context.WithCancel(context.Background())
	s.cancelConn = cancel
	s.status.Connecting = true
	s.status.Connected = false
	s.status.StatusMessage = "Finding fastest server..."
	s.status.Mode = req.Mode
	s.status.LastError = ""
	s.mu.Unlock()

	s.broadcastEvent("status_change", s.status)

	go func() {
		testURL := s.app.Config.TestURL
		if req.Target == "gemini" || req.Target == "ai" {
			testURL = "https://gemini.google.com/"
		}

		s.app.Config.TestURL = testURL

		opts := app.ConnectOptions{
			Country:      req.Country,
			Protocol:     req.Protocol,
			Mode:         req.Mode,
			FastPass:     req.FastPass,
			BatchSize:    req.BatchSize,
			UseCacheOnly: req.UseCacheOnly,
			ManualURI:    req.RawURL,
			Headless:     true,
		}

		err := s.app.RunConnect(ctx, opts)

		s.mu.Lock()
		s.status.Connecting = false
		if err != nil {
			s.status.Connected = false
			s.status.StatusMessage = "Disconnected"
			s.status.LastError = err.Error()
		} else {
			s.status.Connected = true
			s.status.StatusMessage = "Connected"
			if s.app.Engine != nil {
				s.status.CurrentNode = s.app.Engine.GetCurrentNode()
			}
		}
		s.mu.Unlock()

		s.broadcastEvent("status_change", s.status)
	}()

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"status":  "connecting",
		"message": "Connection sequence initiated",
	})
}

func (s *Server) handleDisconnect(w http.ResponseWriter, r *http.Request) {
	enableCORS(w)
	if r.Method == http.MethodOptions {
		return
	}

	s.mu.Lock()
	if s.cancelConn != nil {
		s.cancelConn()
	}
	if s.app.Engine != nil {
		_ = s.app.Engine.Stop()
	}
	s.status.Connected = false
	s.status.Connecting = false
	s.status.StatusMessage = "Disconnected"
	s.status.CurrentNode = nil
	s.status.UploadSpeed = 0
	s.status.DownloadSpeed = 0
	s.mu.Unlock()

	s.broadcastEvent("status_change", s.status)

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]string{"status": "disconnected"})
}

func (s *Server) handleNodes(w http.ResponseWriter, r *http.Request) {
	enableCORS(w)
	if r.Method == http.MethodOptions {
		return
	}

	cached, err := cache.LoadCache()
	var nodes []*model.ProxyNode
	if err == nil && len(cached.Nodes) > 0 {
		nodes = cached.Nodes
	} else {
		fetched, _ := s.app.Fetcher.FetchAll(r.Context(), s.app.Config.Subscriptions)
		nodes = fetched
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"count": len(nodes),
		"nodes": nodes,
	})
}

func (s *Server) handleEvents(w http.ResponseWriter, r *http.Request) {
	enableCORS(w)
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming unsupported", http.StatusInternalServerError)
		return
	}

	msgChan := make(chan string, 16)
	s.clientsMu.Lock()
	s.clients[msgChan] = true
	s.clientsMu.Unlock()

	defer func() {
		s.clientsMu.Lock()
		delete(s.clients, msgChan)
		s.clientsMu.Unlock()
		close(msgChan)
	}()

	// Send initial status
	s.mu.RLock()
	initJSON, _ := json.Marshal(s.status)
	s.mu.RUnlock()
	fmt.Fprintf(w, "event: status\ndata: %s\n\n", string(initJSON))
	flusher.Flush()

	for {
		select {
		case <-r.Context().Done():
			return
		case msg := <-msgChan:
			fmt.Fprintf(w, "%s\n\n", msg)
			flusher.Flush()
		}
	}
}

func (s *Server) broadcastEvent(eventType string, data interface{}) {
	payload, err := json.Marshal(data)
	if err != nil {
		return
	}
	msg := fmt.Sprintf("event: %s\ndata: %s", eventType, string(payload))

	s.clientsMu.Lock()
	defer s.clientsMu.Unlock()
	for ch := range s.clients {
		select {
		case ch <- msg:
		default:
		}
	}
}

func (s *Server) statsLoop() {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		s.mu.RLock()
		connected := s.status.Connected
		s.mu.RUnlock()

		if connected && s.app.Engine != nil {
			stats := s.app.Engine.GetStats()
			if stats != nil {
				upTotal, downTotal, upSpeed, downSpeed := stats.Snapshot()
				s.mu.Lock()
				s.status.UploadSpeed = upSpeed
				s.status.DownloadSpeed = downSpeed
				s.status.TotalUpload = upTotal
				s.status.TotalDownload = downTotal
				if s.status.CurrentNode == nil {
					s.status.CurrentNode = s.app.Engine.GetCurrentNode()
				}
				currStatus := s.status
				s.mu.Unlock()

				s.broadcastEvent("stats_tick", currStatus)
			}
		}
	}
}
