package util

import "testing"

func TestCountryCodeToFlag(t *testing.T) {
	tests := []struct {
		code string
		flag string
	}{
		{"DE", "🇩🇪"},
		{"NL", "🇳🇱"},
		{"US", "🇺🇸"},
		{"GB", "🇬🇧"},
		{"FR", "🇫🇷"},
		{"TR", "🇹🇷"},
		{"IR", "🇮🇷"},
	}

	for _, tt := range tests {
		got := CountryCodeToFlag(tt.code)
		if got != tt.flag {
			t.Errorf("CountryCodeToFlag(%s) = %s; want %s", tt.code, got, tt.flag)
		}
	}
}

func TestDetectCountry(t *testing.T) {
	tests := []struct {
		remark string
		want   string
	}{
		{"🇩🇪 Germany High Speed", "DE"},
		{"[NL] Netherlands Fast", "NL"},
		{"US-01-Reality", "US"},
		{"Server in France", "FR"},
		{"Random Server", "UN"},
		{"🇸🇨 @MrMeshkyChannel 888771", "SC"},
		{"🇷🇺 @MrMeshkyChannel 837211", "RU"},
		{"🇭🇰 @MrMeshkyChannel 812345", "HK"},
	}

	for _, tt := range tests {
		got := DetectCountry(tt.remark)
		if got != tt.want {
			t.Errorf("DetectCountry(%s) = %s; want %s", tt.remark, got, tt.want)
		}
	}
}
