package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/pelletier/go-toml"
)

// Load creates a Config with the following priority:
// 1. Command-line flags
// 2. Environment variables
// 3. Config file
// 4. Default values
func Load(args []string, configFile string) (*Config, error) {
	// Start with default configuration
	cfg := DefaultConfig()

	// Load from config file if specified
	if configFile != "" {
		if err := cfg.LoadConfigFromFile(configFile); err != nil {
			return nil, fmt.Errorf("error loading config file: %w", err)
		}
	} else {
		// Check default config locations
		configLocations := []string{
			"./pica.json",
			"./pica.toml",
			"./configs/pica.json",
			"./configs/pica.toml",
			filepath.Join(os.Getenv("HOME"), ".pica/config.json"),
			filepath.Join(os.Getenv("HOME"), ".pica/config.toml"),
		}

		for _, location := range configLocations {
			if _, err := os.Stat(location); err == nil {
				if err := cfg.LoadConfigFromFile(location); err == nil {
					break
				}
			}
		}
	}

	// Load from environment variables (overrides config file)
	cfg.LoadFromEnvironment()

	// Parse command-line flags (overrides environment and config file)
	if err := cfg.ParseFlags(args); err != nil {
		return nil, fmt.Errorf("error parsing command-line flags: %w", err)
	}

	// Ensure directories exist
	if err := cfg.LoadDefaults(); err != nil {
		return nil, fmt.Errorf("error loading defaults: %w", err)
	}

	// Validate configuration
	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}

	return cfg, nil
}

// SaveConfig saves the current configuration to a file
func (cfg *Config) SaveConfig(filename string) error {
	dir := filepath.Dir(filename)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	var data []byte
	var err error

	// Determine file type based on extension
	ext := filepath.Ext(filename)
	switch ext {
	case ".json":
		data, err = json.MarshalIndent(cfg, "", "  ")
		if err != nil {
			return fmt.Errorf("error encoding JSON config: %w", err)
		}
	case ".toml":
		data, err = toml.Marshal(cfg)
		if err != nil {
			return fmt.Errorf("error encoding TOML config: %w", err)
		}
	default:
		return fmt.Errorf("unsupported config file format: %s", ext)
	}

	if err := os.WriteFile(filename, data, 0644); err != nil {
		return fmt.Errorf("error writing config file: %w", err)
	}

	return nil
}
