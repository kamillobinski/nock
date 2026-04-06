package config

import (
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

func Init() error {
	dir, err := os.UserConfigDir()
	if err != nil {
		return err
	}

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(filepath.Join(dir, "nock"))

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return err
		}
	}

	return nil
}