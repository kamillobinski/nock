package system

import (
	"os"
	"path/filepath"

	"go.yaml.in/yaml/v3"
)

type Config struct {
	ASCII  string  `yaml:"ascii"`
	Groups []Group `yaml:"groups"`
}

type Group struct {
	Name    string `yaml:"name"`
	ShowName bool   `yaml:"show_name"`
	Items   []Item `yaml:"items"`
}

type Item struct {
	Label    string `yaml:"label"`
	Script   string `yaml:"script"`
	Cache    bool   `yaml:"cache"`
	CacheTTL string `yaml:"cache_ttl"`
}

func Dir() (string, error) {
	dir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "nock", "system"), nil
}

func LoadConfig(systemDir string) (*Config, error) {
	data, err := os.ReadFile(filepath.Join(systemDir, "config.yaml"))
	if err != nil {
		return nil, err
	}
	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}
