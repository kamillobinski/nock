package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

func Init() error {
	dir, err := os.UserConfigDir()
	if err != nil {
		return err
	}

	nockDir := filepath.Join(dir, "nock")

	if err := os.MkdirAll(nockDir, 0755); err != nil {
		return err
	}

	if err := os.MkdirAll(filepath.Join(nockDir, "scripts"), 0755); err != nil {
		return err
	}

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(nockDir)

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return err
		}
		if err := seed(nockDir); err != nil {
			return err
		}
		return viper.ReadInConfig()
	}

	return nil
}

func seed(nockDir string) error {
	configPath := filepath.Join(nockDir, "config.yaml")
	content := fmt.Sprintf("nav:\n  config: %s\n", nockDir)
	return os.WriteFile(configPath, []byte(content), 0644)
}
