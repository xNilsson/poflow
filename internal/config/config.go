package config

import (
	"fmt"
	"os"
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

// GetAllPOFiles returns paths to all .po files in the gettext directory
func (c *Config) GetAllPOFiles() ([]string, error) {
	if c.GettextPath == "" {
		return nil, fmt.Errorf("gettext_path not set in config file")
	}

	var poFiles []string

	// Find all .po files recursively
	err := filepath.Walk(c.GettextPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && filepath.Ext(path) == ".po" {
			poFiles = append(poFiles, path)
		}
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to walk gettext directory: %w", err)
	}

	return poFiles, nil
}

// GetPOTFile returns the path to the .pot template file if it exists
func (c *Config) GetPOTFile() (string, error) {
	potPath, err := c.ResolvePOTPath()
	if err != nil {
		return "", err
	}

	// Check if file exists
	if _, err := os.Stat(potPath); err != nil {
		if os.IsNotExist(err) {
			return "", fmt.Errorf("template file not found: %s", potPath)
		}
		return "", fmt.Errorf("failed to stat template file: %w", err)
	}

	return potPath, nil
}
