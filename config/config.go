package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// Context represents a named Elasticsearch connection configuration.
type Context struct {
	CloudID          string `yaml:"cloud_id,omitempty"`
	ElasticsearchURL string `yaml:"elasticsearch_url,omitempty"`
	KibanaURL        string `yaml:"kibana_url,omitempty"`
	APIKey           string `yaml:"api_key,omitempty"`
}

// Config represents the full configuration file.
type Config struct {
	CurrentContext string             `yaml:"current-context,omitempty"`
	Contexts       map[string]Context `yaml:"contexts,omitempty"`
}

// ConfigFilePath returns the OS-appropriate config file path:
//   - Linux:   $XDG_CONFIG_HOME/elastic/config.yaml  (fallback ~/.config/elastic/config.yaml)
//   - macOS:   ~/Library/Application Support/elastic/config.yaml
//   - Windows: %AppData%\elastic\config.yaml
func ConfigFilePath() (string, error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return "", fmt.Errorf("getting user config dir: %w", err)
	}
	return filepath.Join(configDir, "es", "config.yaml"), nil
}

// Load reads the config file and returns the parsed Config.
// If the file does not exist, an empty Config with an initialised Contexts map
// is returned (no error).
func Load() (*Config, error) {
	path, err := ConfigFilePath()
	if err != nil {
		return nil, err
	}

	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return &Config{Contexts: make(map[string]Context)}, nil
	}
	if err != nil {
		return nil, fmt.Errorf("reading config file %s: %w", path, err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("parsing config file %s: %w", path, err)
	}

	if cfg.Contexts == nil {
		cfg.Contexts = make(map[string]Context)
	}

	return &cfg, nil
}

// Save writes cfg to the OS-appropriate config file, creating parent
// directories as needed.  The file is written with mode 0600.
func Save(cfg *Config) error {
	path, err := ConfigFilePath()
	if err != nil {
		return err
	}

	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return fmt.Errorf("creating config directory: %w", err)
	}

	data, err := yaml.Marshal(cfg)
	if err != nil {
		return fmt.Errorf("serializing config: %w", err)
	}

	if err := os.WriteFile(path, data, 0o600); err != nil {
		return fmt.Errorf("writing config file %s: %w", path, err)
	}

	return nil
}
