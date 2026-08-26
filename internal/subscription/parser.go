package subscription

import (
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"github.com/Mr-Meshky/vify-cli/internal/model"
	"github.com/Mr-Meshky/vify-cli/internal/util"
)

// VMessConfigSchema represents standard VMess base64 JSON payload
type VMessConfigSchema struct {
	V    interface{} `json:"v"`
	Ps   string      `json:"ps"`
	Add  string      `json:"add"`
	Port interface{} `json:"port"`
	ID   string      `json:"id"`
	Aid  interface{} `json:"aid"`
	Scy  string      `json:"scy"`
	Net  string      `json:"net"`
	Type string      `json:"type"`
	Host string      `json:"host"`
	Path string      `json:"path"`
	TLS  string      `json:"tls"`
	Sni  string      `json:"sni"`
	Alpn string      `json:"alpn"`
	Fp   string      `json:"fp"`
}

// ParseURI parses a single proxy URI (vless, vmess, trojan, ss) into a ProxyNode
func ParseURI(rawURI string) (*model.ProxyNode, error) {
	rawURI = strings.TrimSpace(rawURI)
	if rawURI == "" {
		return nil, fmt.Errorf("empty uri")
	}

	if strings.HasPrefix(rawURI, "vless://") {
		return parseVLESS(rawURI)
	} else if strings.HasPrefix(rawURI, "vmess://") {
		return parseVMess(rawURI)
	} else if strings.HasPrefix(rawURI, "trojan://") {
		return parseTrojan(rawURI)
	} else if strings.HasPrefix(rawURI, "ss://") {
		return parseShadowsocks(rawURI)
	}

	return nil, fmt.Errorf("unsupported protocol scheme: %s", rawURI)
}

func parseVLESS(rawURI string) (*model.ProxyNode, error) {
	u, err := url.Parse(rawURI)
	if err != nil {
		return nil, fmt.Errorf("invalid vless uri: %w", err)
	}

	port, _ := strconv.Atoi(u.Port())
	if port == 0 {
		port = 443
	}

	q := u.Query()
	name := u.Fragment
	if name == "" {
		name = fmt.Sprintf("VLESS-%s:%d", u.Hostname(), port)
	} else {
		name, _ = url.QueryUnescape(name)
	}

	countryCode := util.DetectCountry(name)

	node := &model.ProxyNode{
		ID:          generateID("vless", u.Hostname(), port, u.User.Username()),
		RawURI:      rawURI,
		Protocol:    model.ProtocolVLESS,
		Name:        name,
		Server:      u.Hostname(),
		Port:        port,
		UUID:        u.User.Username(),
		Security:    q.Get("security"),
		SNI:         q.Get("sni"),
		Host:        q.Get("host"),
		Path:        q.Get("path"),
		Type:        q.Get("type"),
		ServiceName: q.Get("serviceName"),
		PublicKey:   q.Get("pbk"),
		ShortID:     q.Get("sid"),
		SpiderX:     q.Get("spx"),
		Fingerprint: q.Get("fp"),
		Flow:        q.Get("flow"),
		CountryCode: countryCode,
		CountryFlag: util.CountryCodeToFlag(countryCode),
		Extra:       make(map[string]string),
	}

	if node.SNI == "" && node.Host != "" {
		node.SNI = node.Host
	}

	return node, nil
}

func parseVMess(rawURI string) (*model.ProxyNode, error) {
	b64Data := strings.TrimPrefix(rawURI, "vmess://")
	// Fix base64 padding if needed
	if rem := len(b64Data) % 4; rem != 0 {
		b64Data += strings.Repeat("=", 4-rem)
	}

	decoded, err := base64.StdEncoding.DecodeString(b64Data)
	if err != nil {
		decoded, err = base64.URLEncoding.DecodeString(b64Data)
		if err != nil {
			return nil, fmt.Errorf("failed to decode vmess base64: %w", err)
		}
	}

	var schema VMessConfigSchema
	if err := json.Unmarshal(decoded, &schema); err != nil {
		return nil, fmt.Errorf("failed to parse vmess json: %w", err)
	}

	var port int
	switch p := schema.Port.(type) {
	case float64:
		port = int(p)
	case string:
		port, _ = strconv.Atoi(p)
	default:
		port = 443
	}

	name := schema.Ps
	if name == "" {
		name = fmt.Sprintf("VMess-%s:%d", schema.Add, port)
	}
	countryCode := util.DetectCountry(name)

	sec := ""
	if schema.TLS == "tls" || schema.TLS == "1" {
		sec = "tls"
	}

	node := &model.ProxyNode{
		ID:          generateID("vmess", schema.Add, port, schema.ID),
		RawURI:      rawURI,
		Protocol:    model.ProtocolVMess,
		Name:        name,
		Server:      schema.Add,
		Port:        port,
		UUID:        schema.ID,
		Security:    sec,
		SNI:         schema.Sni,
		Host:        schema.Host,
		Path:        schema.Path,
		Type:        schema.Net,
		Fingerprint: schema.Fp,
		CountryCode: countryCode,
		CountryFlag: util.CountryCodeToFlag(countryCode),
		Extra:       make(map[string]string),
	}

	return node, nil
}

func parseTrojan(rawURI string) (*model.ProxyNode, error) {
	u, err := url.Parse(rawURI)
	if err != nil {
		return nil, fmt.Errorf("invalid trojan uri: %w", err)
	}

	port, _ := strconv.Atoi(u.Port())
	if port == 0 {
		port = 443
	}

	q := u.Query()
	name := u.Fragment
	if name == "" {
		name = fmt.Sprintf("Trojan-%s:%d", u.Hostname(), port)
	} else {
		name, _ = url.QueryUnescape(name)
	}

	countryCode := util.DetectCountry(name)

	node := &model.ProxyNode{
		ID:          generateID("trojan", u.Hostname(), port, u.User.Username()),
		RawURI:      rawURI,
		Protocol:    model.ProtocolTrojan,
		Name:        name,
		Server:      u.Hostname(),
		Port:        port,
		Password:    u.User.Username(),
		Security:    "tls",
		SNI:         q.Get("sni"),
		Host:        q.Get("host"),
		Path:        q.Get("path"),
		Type:        q.Get("type"),
		ServiceName: q.Get("serviceName"),
		Fingerprint: q.Get("fp"),
		CountryCode: countryCode,
		CountryFlag: util.CountryCodeToFlag(countryCode),
		Extra:       make(map[string]string),
	}

	if node.SNI == "" && node.Host != "" {
		node.SNI = node.Host
	}

	return node, nil
}

func parseShadowsocks(rawURI string) (*model.ProxyNode, error) {
	// ss://base64(method:password@server:port)#name
	// or ss://base64(method:password)@server:port#name
	u, err := url.Parse(rawURI)
	if err != nil {
		return nil, fmt.Errorf("invalid ss uri: %w", err)
	}

	name := u.Fragment
	if name != "" {
		name, _ = url.QueryUnescape(name)
	}

	var method, password, server string
	var port int

	if u.User != nil {
		userPart := u.User.Username()
		// Try decode userPart
		if decoded, err := decodeB64(userPart); err == nil && strings.Contains(string(decoded), ":") {
			parts := strings.SplitN(string(decoded), ":", 2)
			method = parts[0]
			password = parts[1]
			server = u.Hostname()
			port, _ = strconv.Atoi(u.Port())
		} else {
			method = u.User.Username()
			password, _ = u.User.Password()
			server = u.Hostname()
			port, _ = strconv.Atoi(u.Port())
		}
	} else {
		// Entire host might be base64
		host := u.Host
		if decoded, err := decodeB64(host); err == nil {
			// format: method:password@server:port
			s := string(decoded)
			if atIdx := strings.LastIndex(s, "@"); atIdx != -1 {
				cred := s[:atIdx]
				srv := s[atIdx+1:]
				if credParts := strings.SplitN(cred, ":", 2); len(credParts) == 2 {
					method = credParts[0]
					password = credParts[1]
				}
				if hostPort := strings.Split(srv, ":"); len(hostPort) == 2 {
					server = hostPort[0]
					port, _ = strconv.Atoi(hostPort[1])
				}
			}
		}
	}

	if server == "" || port == 0 {
		return nil, fmt.Errorf("failed to parse shadowsocks server/port")
	}

	if name == "" {
		name = fmt.Sprintf("SS-%s:%d", server, port)
	}
	countryCode := util.DetectCountry(name)

	return &model.ProxyNode{
		ID:          generateID("ss", server, port, method),
		RawURI:      rawURI,
		Protocol:    model.ProtocolShadowsocks,
		Name:        name,
		Server:      server,
		Port:        port,
		Password:    password,
		Method:      method,
		CountryCode: countryCode,
		CountryFlag: util.CountryCodeToFlag(countryCode),
		Extra:       make(map[string]string),
	}, nil
}

func decodeB64(s string) ([]byte, error) {
	if rem := len(s) % 4; rem != 0 {
		s += strings.Repeat("=", 4-rem)
	}
	b, err := base64.StdEncoding.DecodeString(s)
	if err == nil {
		return b, nil
	}
	return base64.URLEncoding.DecodeString(s)
}

func generateID(proto, server string, port int, cred string) string {
	h := md5.New()
	h.Write([]byte(fmt.Sprintf("%s:%s:%d:%s", proto, server, port, cred)))
	return hex.EncodeToString(h.Sum(nil))[:12]
}
