package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

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
		if err := viper.ReadInConfig(); err != nil {
			return err
		}
	}

	if err := initKnowledgeDirs(nockDir); err != nil {
		return err
	}
	if err := initSystemDir(nockDir); err != nil {
		return err
	}
	return initLogDir()
}

func initKnowledgeDirs(nockDir string) error {
	knowledgeDir := viper.GetString("knowledge.dir")
	if knowledgeDir == "" {
		knowledgeDir = filepath.Join(nockDir, "knowledge")
	} else if strings.HasPrefix(knowledgeDir, "~/") {
		home, _ := os.UserHomeDir()
		knowledgeDir = filepath.Join(home, knowledgeDir[2:])
	}

	return os.MkdirAll(filepath.Join(knowledgeDir, "pending"), 0755)
}

func initLogDir() error {
	logDir := viper.GetString("log.dir")
	if logDir == "" {
		home, _ := os.UserHomeDir()
		logDir = filepath.Join(home, "Documents", "logs")
	} else if strings.HasPrefix(logDir, "~/") {
		home, _ := os.UserHomeDir()
		logDir = filepath.Join(home, logDir[2:])
	}
	return os.MkdirAll(logDir, 0755)
}

func seed(nockDir string) error {
	configPath := filepath.Join(nockDir, "config.yaml")
	content := fmt.Sprintf("nav:\n  config: %s\n", nockDir)
	return os.WriteFile(configPath, []byte(content), 0644)
}

func initSystemDir(nockDir string) error {
	systemDir := filepath.Join(nockDir, "system")

	if err := os.MkdirAll(filepath.Join(systemDir, "scripts"), 0755); err != nil {
		return err
	}

	configPath := filepath.Join(systemDir, "config.yaml")
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return os.WriteFile(configPath, []byte("groups: []\n"), 0644)
	}

	return nil
}
