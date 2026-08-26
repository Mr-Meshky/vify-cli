package core

import (
	"testing"

	"github.com/Mr-Meshky/vify-cli/internal/model"
)

func TestGenerateConfig(t *testing.T) {
	node := &model.ProxyNode{
		ID:          "test-1",
		Protocol:    model.ProtocolVLESS,
		Name:        "Test-DE",
		Server:      "1.2.3.4",
		Port:        443,
		UUID:        "e990c012-32a7-47b8-b7ae-24ad2657e2d9",
		Security:    "reality",
		SNI:         "speedtest.net",
		PublicKey:   "abcdef",
		CountryCode: "DE",
	}

	// Test TUN Mode
	cfgTUN, err := GenerateConfig(node, model.ModeTUN, 2080, 2081, 9090, true)
	if err != nil {
		t.Fatalf("GenerateConfig TUN error: %v", err)
	}

	if len(cfgTUN.Inbounds) < 3 {
		t.Errorf("expected at least 3 inbounds in TUN mode, got %d", len(cfgTUN.Inbounds))
	}

	jsonBytes, err := cfgTUN.ToJSON()
	if err != nil {
		t.Fatalf("ToJSON error: %v", err)
	}

	if len(jsonBytes) == 0 {
		t.Errorf("empty json output")
	}

	// Test System Proxy Mode
	cfgSys, err := GenerateConfig(node, model.ModeSystemProxy, 2080, 2081, 0, true)
	if err != nil {
		t.Fatalf("GenerateConfig SystemProxy error: %v", err)
	}

	if len(cfgSys.Inbounds) != 2 {
		t.Errorf("expected 2 inbounds in SystemProxy mode, got %d", len(cfgSys.Inbounds))
	}
}
