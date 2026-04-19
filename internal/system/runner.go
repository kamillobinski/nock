package system

import (
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

func FetchValue(item Item, systemDir string, cache *Cache, force bool) string {
	cacheKey := item.Script

	if item.Cache && !force {
		if item.CacheTTL == "" {
			if val, ok := cache.Get(cacheKey, 0); ok {
				return val
			}
		} else if ttl, err := time.ParseDuration(item.CacheTTL); err == nil {
			if val, ok := cache.Get(cacheKey, ttl); ok {
				return val
			}
		}
	}

	scriptPath := item.Script
	if !filepath.IsAbs(scriptPath) {
		scriptPath = filepath.Join(systemDir, scriptPath)
	}

	// try direct execution first (respects shebang)
	out, err := exec.Command(scriptPath).Output()
	if err != nil {
		// fallback to sh
		out, err = exec.Command("sh", scriptPath).Output()
		if err != nil {
			return "N/A"
		}
	}

	value := strings.TrimSpace(string(out))

	if item.Cache {
		cache.Set(cacheKey, value)
	}

	return value
}
