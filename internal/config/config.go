package config

import (
	"fmt"
	"path/filepath"

	"github.com/spf13/viper"
)

// Config holds the application configuration
type Config struct {
	GettextPath string `mapstructure:"gettext_path"`
}

// Load returns the loaded configuration
func Load() (*Config, error) {
	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("unable to decode config: %w", err)
	}
	return &cfg, nil
}

// ResolvePOPath resolves the .po file path for a given language
// Returns: {gettext_path}/{lang}/LC_MESSAGES/default.po
func (c *Config) ResolvePOPath(lang string) (string, error) {
	if c.GettextPath == "" {
		return "", fmt.Errorf("gettext_path not set in config file")
	}
	if lang == "" {
		return "", fmt.Errorf("language code is required")
	}

	path := filepath.Join(c.GettextPath, lang, "LC_MESSAGES", "default.po")
	return path, nil
}

// ResolvePOTPath resolves the .pot template file path
// Returns: {gettext_path}/default.pot
func (c *Config) ResolvePOTPath() (string, error) {
	if c.GettextPath == "" {
		return "", fmt.Errorf("gettext_path not set in config file")
	}

	path := filepath.Join(c.GettextPath, "default.pot")
	return path, nil
}
