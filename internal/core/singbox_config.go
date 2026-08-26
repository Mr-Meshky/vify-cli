package core

import (
	"encoding/json"
	"fmt"
	"runtime"

	"github.com/Mr-Meshky/vify-cli/internal/model"
)

// SingBoxConfig generates valid sing-box JSON configuration for TUN or System Proxy modes
type SingBoxConfig struct {
	Log          map[string]interface{}   `json:"log"`
	DNS          map[string]interface{}   `json:"dns"`
	Inbounds     []map[string]interface{} `json:"inbounds"`
	Outbounds    []map[string]interface{} `json:"outbounds"`
	Route        map[string]interface{}   `json:"route"`
	Experimental map[string]interface{}   `json:"experimental,omitempty"`
}

// GenerateConfig generates a complete sing-box configuration struct.
// clashAPIPort enables sing-box's Clash-compatible API (used to read live
// traffic totals) on 127.0.0.1:clashAPIPort; pass 0 to omit it entirely —
// callers spinning up multiple concurrent short-lived instances (e.g. real
// connectivity verification) that don't need traffic stats MUST pass 0,
// since every instance would otherwise fight over the same hardcoded port
// and fail to start.
func GenerateConfig(node *model.ProxyNode, mode model.ConnectionMode, socksPort, httpPort, clashAPIPort int, bypassIran bool) (*SingBoxConfig, error) {
	if node == nil {
		return nil, fmt.Errorf("node cannot be nil")
	}

	cfg := &SingBoxConfig{
		Log: map[string]interface{}{
			"level":     "warn",
			"timestamp": true,
		},
		DNS: map[string]interface{}{
			"servers": []map[string]interface{}{
				{
					"tag":              "dns-remote",
					"address":          "tcp://8.8.8.8",
					"address_resolver": "dns-direct",
					"detour":           "proxy",
				},
				{
					"tag":     "dns-direct",
					"address": "1.1.1.1",
					"detour":  "direct",
				},
				{
					"tag":     "dns-block",
					"address": "rcode://success",
				},
			},
			"rules": []map[string]interface{}{
				{
					"domain_suffix": []string{".ir", ".xn--mgba3a4f16a"},
					"server":        "dns-direct",
				},
				{
					"outbound": "any",
					"server":   "dns-remote",
				},
			},
			"strategy": "prefer_ipv4",
		},
		Inbounds: []map[string]interface{}{
			{
				"type":        "mixed",
				"tag":         "mixed-in",
				"listen":      "127.0.0.1",
				"listen_port": socksPort,
			},
			{
				"type":        "http",
				"tag":         "http-in",
				"listen":      "127.0.0.1",
				"listen_port": httpPort,
			},
		},
	}

	if clashAPIPort > 0 {
		cfg.Experimental = map[string]interface{}{
			"clash_api": map[string]interface{}{
				"external_controller": fmt.Sprintf("127.0.0.1:%d", clashAPIPort),
				"secret":              "",
			},
		}
	}

	// Add TUN inbound if requested
	if mode == model.ModeTUN {
		ifaceName := "vify-tun0"
		if runtime.GOOS == "darwin" {
			ifaceName = "utun99"
		}
		tunInbound := map[string]interface{}{
			"type":           "tun",
			"tag":            "tun-in",
			"interface_name": ifaceName,
			"address":        []string{"172.19.0.1/30"},
			"auto_route":     true,
			"strict_route":   true,
			"stack":          "gvisor",
			"sniff":          true,
		}
		cfg.Inbounds = append(cfg.Inbounds, tunInbound)
	}

	// Generate Outbound for Proxy Node
	proxyOutbound, err := generateProxyOutbound(node)
	if err != nil {
		return nil, err
	}

	cfg.Outbounds = []map[string]interface{}{
		proxyOutbound,
		{
			"type": "direct",
			"tag":  "direct",
		},
		{
			"type": "block",
			"tag":  "block",
		},
		{
			"type": "dns",
			"tag":  "dns-out",
		},
	}

	// Setup Routing Rules
	var routeRules []map[string]interface{}

	// DNS traffic routing
	routeRules = append(routeRules, map[string]interface{}{
		"protocol": "dns",
		"outbound": "dns-out",
	})

	// Direct Iran & LAN bypass
	if bypassIran {
		routeRules = append(routeRules, map[string]interface{}{
			"domain_suffix": []string{".ir", ".xn--mgba3a4f16a"},
			"outbound":      "direct",
		})
		routeRules = append(routeRules, map[string]interface{}{
			"ip_cidr":  IranDirectIPRanges[2:], // private ranges
			"outbound": "direct",
		})
	}

	cfg.Route = map[string]interface{}{
		"rules":                 routeRules,
		"final":                 "proxy",
		"auto_detect_interface": true,
	}

	return cfg, nil
}

func generateProxyOutbound(node *model.ProxyNode) (map[string]interface{}, error) {
	out := map[string]interface{}{
		"tag":         "proxy",
		"server":      node.Server,
		"server_port": node.Port,
	}

	switch node.Protocol {
	case model.ProtocolVLESS:
		out["type"] = "vless"
		out["uuid"] = node.UUID
		if node.Flow != "" {
			out["flow"] = node.Flow
		}

		tlsObj := map[string]interface{}{
			"enabled": node.Security == "tls" || node.Security == "reality",
		}
		if node.SNI != "" {
			tlsObj["server_name"] = node.SNI
		}
		fp := node.Fingerprint
		if fp == "" && node.Security == "reality" {
			fp = "chrome"
		}
		if fp != "" {
			tlsObj["utls"] = map[string]interface{}{
				"enabled":     true,
				"fingerprint": fp,
			}
		}
		if node.Security == "reality" {
			tlsObj["reality"] = map[string]interface{}{
				"enabled":    true,
				"public_key": node.PublicKey,
				"short_id":   node.ShortID,
			}
		}
		if node.Security == "tls" {
			tlsObj["insecure"] = true
		}
		if tlsObj["enabled"].(bool) {
			out["tls"] = tlsObj
		}

		if node.Type == "ws" {
			out["transport"] = map[string]interface{}{
				"type": "ws",
				"path": node.Path,
				"headers": map[string]string{
					"Host": node.Host,
				},
			}
		} else if node.Type == "grpc" {
			out["transport"] = map[string]interface{}{
				"type":         "grpc",
				"service_name": node.ServiceName,
			}
		}

	case model.ProtocolVMess:
		out["type"] = "vmess"
		out["uuid"] = node.UUID
		out["security"] = "auto"
		out["alter_id"] = 0

		if node.Security == "tls" {
			tlsObj := map[string]interface{}{
				"enabled":  true,
				"insecure": true,
			}
			if node.SNI != "" {
				tlsObj["server_name"] = node.SNI
			}
			if node.Fingerprint != "" {
				tlsObj["utls"] = map[string]interface{}{
					"enabled":     true,
					"fingerprint": node.Fingerprint,
				}
			}
			out["tls"] = tlsObj
		}

		if node.Type == "ws" {
			out["transport"] = map[string]interface{}{
				"type": "ws",
				"path": node.Path,
				"headers": map[string]string{
					"Host": node.Host,
				},
			}
		}

	case model.ProtocolTrojan:
		out["type"] = "trojan"
		out["password"] = node.Password

		tlsObj := map[string]interface{}{
			"enabled":  true,
			"insecure": true,
		}
		if node.SNI != "" {
			tlsObj["server_name"] = node.SNI
		}
		if node.Fingerprint != "" {
			tlsObj["utls"] = map[string]interface{}{
				"enabled":     true,
				"fingerprint": node.Fingerprint,
			}
		}
		out["tls"] = tlsObj

		if node.Type == "ws" {
			out["transport"] = map[string]interface{}{
				"type": "ws",
				"path": node.Path,
				"headers": map[string]string{
					"Host": node.Host,
				},
			}
		} else if node.Type == "grpc" {
			out["transport"] = map[string]interface{}{
				"type":         "grpc",
				"service_name": node.ServiceName,
			}
		}

	case model.ProtocolShadowsocks:
		out["type"] = "shadowsocks"
		out["method"] = node.Method
		out["password"] = node.Password

	default:
		return nil, fmt.Errorf("unsupported protocol: %s", node.Protocol)
	}

	return out, nil
}

// ToJSON serializes SingBoxConfig to formatted JSON
func (c *SingBoxConfig) ToJSON() ([]byte, error) {
	return json.MarshalIndent(c, "", "  ")
}
