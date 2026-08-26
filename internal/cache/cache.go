package cache

import (
	"encoding/json"
	"os"
	"sync"
	"time"

	"github.com/Mr-Meshky/vify-cli/internal/config"
	"github.com/Mr-Meshky/vify-cli/internal/model"
)

// CacheData holds cached healthy proxies and last update timestamp
type CacheData struct {
	UpdatedAt time.Time          `json:"updated_at"`
	Nodes     []*model.ProxyNode `json:"nodes"`
}

var (
	cacheMu sync.Mutex
)

// LoadCache loads previously tested healthy nodes from ~/.vify/cache.json
func LoadCache() (*CacheData, error) {
	cacheMu.Lock()
	defer cacheMu.Unlock()

	filePath := config.GetCachePath()
	data, err := os.ReadFile(filePath)
	if err != nil {
		return &CacheData{
			UpdatedAt: time.Time{},
			Nodes:     []*model.ProxyNode{},
		}, nil
	}

	var c CacheData
	if err := json.Unmarshal(data, &c); err != nil {
		return &CacheData{
			UpdatedAt: time.Time{},
			Nodes:     []*model.ProxyNode{},
		}, nil
	}

	return &c, nil
}

// SaveCache writes healthy nodes to ~/.vify/cache.json
func SaveCache(nodes []*model.ProxyNode) error {
	cacheMu.Lock()
	defer cacheMu.Unlock()

	c := CacheData{
		UpdatedAt: time.Now(),
		Nodes:     nodes,
	}

	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(config.GetCachePath(), data, 0644)
}

// SaveSession saves the currently running session to ~/.vify/session.json
func SaveSession(session *model.ActiveSession) error {
	data, err := json.MarshalIndent(session, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(config.GetSessionPath(), data, 0644)
}

// LoadSession loads the active session from ~/.vify/session.json
func LoadSession() (*model.ActiveSession, error) {
	data, err := os.ReadFile(config.GetSessionPath())
	if err != nil {
		return nil, err
	}

	var s model.ActiveSession
	if err := json.Unmarshal(data, &s); err != nil {
		return nil, err
	}
	return &s, nil
}

// ClearSession removes ~/.vify/session.json
func ClearSession() {
	_ = os.Remove(config.GetSessionPath())
}
