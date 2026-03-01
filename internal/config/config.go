package config

import (
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

type Config struct {
	GeminiAPIKey string `mapstructure:"GEMINI_API_KEY"`
	OPGGURL      string `mapstructure:"OPGG_MCP_URL"`
}

var AppConfig *Config

func Load() error {
	viper.SetConfigName(".env")
	viper.SetConfigType("env")

	// Look for .env in the current working directory (dev mode)
	viper.AddConfigPath(".")

	// Also look in $HOME/.config/lol/ (installed binary)
	if home, err := os.UserHomeDir(); err == nil {
		viper.AddConfigPath(filepath.Join(home, ".config", "lol"))
	}

	// AutomaticEnv ensures real environment variables always take precedence
	viper.AutomaticEnv()

	// Ignore "config not found" errors — env vars alone are sufficient
	_ = viper.ReadInConfig()

	AppConfig = &Config{}
	return viper.Unmarshal(AppConfig)
}