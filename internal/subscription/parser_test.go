package subscription

import (
	"testing"

	"github.com/Mr-Meshky/vify-cli/internal/model"
)

func TestParseVLESSReality(t *testing.T) {
	raw := "vless://432a5789-f529-4d64-9694-81492b4520b7@185.190.140.20:443?security=reality&sni=speedtest.net&fp=chrome&pbk=J8_3z0P4QJ9x4Wv_qT11U8K72b4Wv_qT11U8K72b4Wv&sid=123456&type=grpc&serviceName=vify-grpc#DE-Frankfurt-VLESS"
	node, err := ParseURI(raw)
	if err != nil {
		t.Fatalf("unexpected error parsing vless: %v", err)
	}

	if node.Protocol != model.ProtocolVLESS {
		t.Errorf("expected protocol vless, got %s", node.Protocol)
	}
	if node.Server != "185.190.140.20" {
		t.Errorf("expected server 185.190.140.20, got %s", node.Server)
	}
	if node.Port != 443 {
		t.Errorf("expected port 443, got %d", node.Port)
	}
	if node.Security != "reality" {
		t.Errorf("expected security reality, got %s", node.Security)
	}
	if node.CountryCode != "DE" {
		t.Errorf("expected country DE, got %s", node.CountryCode)
	}
}

func TestParseTrojan(t *testing.T) {
	raw := "trojan://password123@nl.vify.org:443?security=tls&sni=nl.vify.org&type=ws&path=%2Fws#NL-Amsterdam"
	node, err := ParseURI(raw)
	if err != nil {
		t.Fatalf("unexpected error parsing trojan: %v", err)
	}

	if node.Protocol != model.ProtocolTrojan {
		t.Errorf("expected protocol trojan, got %s", node.Protocol)
	}
	if node.Password != "password123" {
		t.Errorf("expected password password123, got %s", node.Password)
	}
	if node.CountryCode != "NL" {
		t.Errorf("expected country NL, got %s", node.CountryCode)
	}
}

func TestParseTrojanComplex(t *testing.T) {
	raw := "trojan://ND91608427@44.246.163.102:443?headerType=none&sni=fleet-bonefish.rooster465.autos&fp=chrome&type=tcp&insecure=0&security=tls&allowInsecure=0#%F0%9F%87%BA%F0%9F%87%B8%20%40MrMeshkyChannel%20385450"
	node, err := ParseURI(raw)
	if err != nil {
		t.Fatalf("unexpected error parsing complex trojan: %v", err)
	}

	if node.Protocol != model.ProtocolTrojan {
		t.Errorf("expected protocol trojan, got %s", node.Protocol)
	}
	if node.Password != "ND91608427" {
		t.Errorf("expected password ND91608427, got %s", node.Password)
	}
	if node.Server != "44.246.163.102" {
		t.Errorf("expected server 44.246.163.102, got %s", node.Server)
	}
	if node.SNI != "fleet-bonefish.rooster465.autos" {
		t.Errorf("expected sni fleet-bonefish.rooster465.autos, got %s", node.SNI)
	}
	if node.Fingerprint != "chrome" {
		t.Errorf("expected fp chrome, got %s", node.Fingerprint)
	}
	if node.CountryCode != "US" {
		t.Errorf("expected country US, got %s", node.CountryCode)
	}
}
