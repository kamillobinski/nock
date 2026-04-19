package system

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sync"
	"time"
)

type cacheEntry struct {
	Value    string    `json:"value"`
	CachedAt time.Time `json:"cached_at"`
}

type Cache struct {
	mu      sync.Mutex
	path    string
	entries map[string]cacheEntry
	dirty   bool
}

func LoadCache(systemDir string) *Cache {
	path := filepath.Join(systemDir, "cache.json")
	c := &Cache{path: path, entries: make(map[string]cacheEntry)}
	data, err := os.ReadFile(path)
	if err != nil {
		return c
	}
	_ = json.Unmarshal(data, &c.entries)
	return c
}

func (c *Cache) Get(key string, ttl time.Duration) (string, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	entry, ok := c.entries[key]
	if !ok {
		return "", false
	}
	if ttl > 0 && time.Since(entry.CachedAt) > ttl {
		return "", false
	}
	return entry.Value, true
}

func (c *Cache) Set(key, value string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.entries[key] = cacheEntry{Value: value, CachedAt: time.Now()}
	c.dirty = true
}

func (c *Cache) Save() error {
	if !c.dirty {
		return nil
	}
	data, err := json.MarshalIndent(c.entries, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(c.path, data, 0644)
}
