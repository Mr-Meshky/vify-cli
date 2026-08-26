package model

import "time"

// Protocol represents the supported proxy protocol types
type Protocol string

const (
	ProtocolVLESS       Protocol = "vless"
	ProtocolVMess       Protocol = "vmess"
	ProtocolTrojan      Protocol = "trojan"
	ProtocolShadowsocks Protocol = "shadowsocks"
	ProtocolUnknown     Protocol = "unknown"
)

// ProxyNode represents a parsed proxy server configuration
type ProxyNode struct {
	ID          string            `json:"id"`
	RawURI      string            `json:"raw_uri"`
	Protocol    Protocol          `json:"protocol"`
	Name        string            `json:"name"`
	Server      string            `json:"server"`
	Port        int               `json:"port"`
	UUID        string            `json:"uuid,omitempty"`
	Password    string            `json:"password,omitempty"`
	Method      string            `json:"method,omitempty"`   // For Shadowsocks
	Security    string            `json:"security,omitempty"` // tls, reality, none
	SNI         string            `json:"sni,omitempty"`
	Host        string            `json:"host,omitempty"`
	Path        string            `json:"path,omitempty"`
	Type        string            `json:"type,omitempty"`         // ws, grpc, tcp, http
	ServiceName string            `json:"service_name,omitempty"` // for gRPC
	PublicKey   string            `json:"public_key,omitempty"`   // for Reality
	ShortID     string            `json:"short_id,omitempty"`     // for Reality
	SpiderX     string            `json:"spider_x,omitempty"`     // for Reality
	Fingerprint string            `json:"fingerprint,omitempty"`  // chrome, safari, etc.
	Flow        string            `json:"flow,omitempty"`         // xtls-rprx-vision
	CountryCode string            `json:"country_code"`
	CountryFlag string            `json:"country_flag"`
	Latency     time.Duration     `json:"latency"`
	IsHealthy   bool              `json:"is_healthy"`
	LastTested  time.Time         `json:"last_tested,omitempty"`
	Extra       map[string]string `json:"extra,omitempty"`
}

// BenchmarkResult represents the result of a latency test on a node
type BenchmarkResult struct {
	Node       *ProxyNode    `json:"node"`
	Latency    time.Duration `json:"latency"`
	Success    bool          `json:"success"`
	StatusCode int           `json:"status_code"`
	Error      string        `json:"error,omitempty"`
}

// ConnectionMode specifies whether to use TUN (Full VPN) or System Proxy
type ConnectionMode string

const (
	ModeTUN         ConnectionMode = "tun"
	ModeSystemProxy ConnectionMode = "system_proxy"
	ModeMixed       ConnectionMode = "mixed"
)

// ActiveSession represents the currently running proxy session
type ActiveSession struct {
	PID           int            `json:"pid"`
	Node          ProxyNode      `json:"node"`
	Mode          ConnectionMode `json:"mode"`
	LocalSocks    int            `json:"local_socks"`
	LocalHTTP     int            `json:"local_http"`
	StartedAt     time.Time      `json:"started_at"`
	OutboundIP    string         `json:"outbound_ip,omitempty"`
	BytesUpload   int64          `json:"bytes_upload"`
	BytesDownload int64          `json:"bytes_download"`
	UploadSpeed   int64          `json:"upload_speed"`   // bytes/sec
	DownloadSpeed int64          `json:"download_speed"` // bytes/sec
}
