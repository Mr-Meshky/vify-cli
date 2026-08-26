package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

// AppConfig represents global persistent application settings
type AppConfig struct {
	Subscriptions       []string `mapstructure:"subscriptions" yaml:"subscriptions"`
	DefaultMode         string   `mapstructure:"default_mode" yaml:"default_mode"` // "tun" or "system_proxy"
	LocalSocksPort      int      `mapstructure:"local_socks_port" yaml:"local_socks_port"`
	LocalHTTPPort       int      `mapstructure:"local_http_port" yaml:"local_http_port"`
	TestTimeoutMS       int      `mapstructure:"test_timeout_ms" yaml:"test_timeout_ms"`
	ConcurrencyLimit    int      `mapstructure:"concurrency_limit" yaml:"concurrency_limit"`
	FastPassThresholdMS int      `mapstructure:"fastpass_threshold_ms" yaml:"fastpass_threshold_ms"`
	TestURL             string   `mapstructure:"test_url" yaml:"test_url"`
	AutoFailover        bool     `mapstructure:"auto_failover" yaml:"auto_failover"`
	WatchdogInterval    int      `mapstructure:"watchdog_interval_sec" yaml:"watchdog_interval_sec"`
	DirectIranBypass    bool     `mapstructure:"direct_iran_bypass" yaml:"direct_iran_bypass"`
	LogLevel            string   `mapstructure:"log_level" yaml:"log_level"`
}

// DefaultConfig returns safe out-of-the-box defaults tuned for Iranian network conditions
func DefaultConfig() *AppConfig {
	return &AppConfig{
		Subscriptions: []string{
			"https://raw.githubusercontent.com/Mr-Meshky/vify/main/configs/vless.txt",
			"https://raw.githubusercontent.com/Mr-Meshky/vify/main/configs/vmess.txt",
			"https://raw.githubusercontent.com/Mr-Meshky/vify/main/configs/trojan.txt",
			"https://raw.githubusercontent.com/Mr-Meshky/vify/main/configs/ss.txt",
		},
		DefaultMode:         "system_proxy",
		LocalSocksPort:      2080,
		LocalHTTPPort:       2081,
		TestTimeoutMS:       2500,
		ConcurrencyLimit:    35,
		FastPassThresholdMS: 800,
		TestURL:             "http://cp.cloudflare.com/generate_204",
		AutoFailover:        true,
		WatchdogInterval:    7,
		DirectIranBypass:    true,
		LogLevel:            "warn",
	}
}

// GetVifyDir returns the ~/.vify path
func GetVifyDir() string {
	home, err := os.UserHomeDir()
	if err != nil {
		home = "."
	}
	dir := filepath.Join(home, ".vify")
	_ = os.MkdirAll(dir, 0755)
	return dir
}

// GetConfigPath returns ~/.vify/config.yaml
func GetConfigPath() string {
	return filepath.Join(GetVifyDir(), "config.yaml")
}

// GetCachePath returns ~/.vify/cache.json
func GetCachePath() string {
	return filepath.Join(GetVifyDir(), "cache.json")
}

// GetSessionPath returns ~/.vify/session.json
func GetSessionPath() string {
	return filepath.Join(GetVifyDir(), "session.json")
}

// LoadConfig reads the configuration file or generates defaults if absent
func LoadConfig() (*AppConfig, error) {
	cfgPath := GetConfigPath()
	v := viper.New()
	v.SetConfigFile(cfgPath)
	v.SetConfigType("yaml")

	defaults := DefaultConfig()
	v.SetDefault("subscriptions", defaults.Subscriptions)
	v.SetDefault("default_mode", defaults.DefaultMode)
	v.SetDefault("local_socks_port", defaults.LocalSocksPort)
	v.SetDefault("local_http_port", defaults.LocalHTTPPort)
	v.SetDefault("test_timeout_ms", defaults.TestTimeoutMS)
	v.SetDefault("concurrency_limit", defaults.ConcurrencyLimit)
	v.SetDefault("fastpass_threshold_ms", defaults.FastPassThresholdMS)
	v.SetDefault("test_url", defaults.TestURL)
	v.SetDefault("auto_failover", defaults.AutoFailover)
	v.SetDefault("watchdog_interval_sec", defaults.WatchdogInterval)
	v.SetDefault("direct_iran_bypass", defaults.DirectIranBypass)
	v.SetDefault("log_level", defaults.LogLevel)

	if _, err := os.Stat(cfgPath); os.IsNotExist(err) {
		// Save default configuration
		cfg := defaults
		_ = SaveConfig(cfg)
		return cfg, nil
	}

	if err := v.ReadInConfig(); err != nil {
		return defaults, fmt.Errorf("failed to read config: %w", err)
	}

	var cfg AppConfig
	if err := v.Unmarshal(&cfg); err != nil {
		return defaults, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return &cfg, nil
}

// SaveConfig persists the current configuration to ~/.vify/config.yaml
func SaveConfig(cfg *AppConfig) error {
	v := viper.New()
	v.SetConfigFile(GetConfigPath())
	v.SetConfigType("yaml")

	v.Set("subscriptions", cfg.Subscriptions)
	v.Set("default_mode", cfg.DefaultMode)
	v.Set("local_socks_port", cfg.LocalSocksPort)
	v.Set("local_http_port", cfg.LocalHTTPPort)
	v.Set("test_timeout_ms", cfg.TestTimeoutMS)
	v.Set("concurrency_limit", cfg.ConcurrencyLimit)
	v.Set("fastpass_threshold_ms", cfg.FastPassThresholdMS)
	v.Set("test_url", cfg.TestURL)
	v.Set("auto_failover", cfg.AutoFailover)
	v.Set("watchdog_interval_sec", cfg.WatchdogInterval)
	v.Set("direct_iran_bypass", cfg.DirectIranBypass)
	v.Set("log_level", cfg.LogLevel)

	return v.WriteConfigAs(GetConfigPath())
}
