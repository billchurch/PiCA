// Package config provides a unified configuration system for PiCA
// that supports loading settings from environment variables,
// command-line flags, and configuration files.
package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/pelletier/go-toml"
)

// Config holds all configuration settings for PiCA
type Config struct {
	// General settings
	LogLevel  string `env:"LOG_LEVEL" flag:"log-level" config:"log_level" default:"info"`
	ConfigDir string `env:"CONFIG_DIR" flag:"config-dir" config:"config_dir" default:"./configs"`

	// CA settings
	CAType           string `env:"CA_TYPE" flag:"ca-type" config:"ca_type" default:"sub"`
	CAConfigFile     string `env:"CA_CONFIG" flag:"ca-config" config:"ca_config" default:""`
	CACertFile       string `env:"CA_CERT" flag:"ca-cert" config:"ca_cert" default:""`
	CRLFile          string `env:"CRL_FILE" flag:"crl-file" config:"crl_file" default:""`
	RootCACertFile   string `env:"ROOT_CA_CERT" flag:"root-ca-cert" config:"root_ca_cert" default:""`
	RootCAConfigFile string `env:"ROOT_CA_CONFIG" flag:"root-ca-config" config:"root_ca_config" default:""`
	CAProfile        string `env:"CA_PROFILE" flag:"ca-profile" config:"ca_profile" default:""`

	// Provider settings
	ProviderType string `env:"PICA_PROVIDER" flag:"provider" config:"provider" default:""`
	KeySlot      string `env:"KEY_SLOT" flag:"key-slot" config:"key_slot" default:"82"`

	// Web server settings
	WebPort      int    `env:"WEB_PORT" flag:"port" config:"web_port" default:"8080"`
	WebRoot      string `env:"WEB_ROOT" flag:"webroot" config:"web_root" default:"./web/html"`
	WebTLSCert   string `env:"WEB_TLS_CERT" flag:"tls-cert" config:"web_tls_cert" default:""`
	WebTLSKey    string `env:"WEB_TLS_KEY" flag:"tls-key" config:"web_tls_key" default:""`
	EnableHTTPS  bool   `env:"ENABLE_HTTPS" flag:"https" config:"enable_https" default:"false"`
	RedirectHTTP bool   `env:"REDIRECT_HTTP" flag:"redirect-http" config:"redirect_http" default:"true"`

	// Storage settings
	CertDir     string `env:"CERT_DIR" flag:"certdir" config:"cert_dir" default:"./certs"`
	CSRDir      string `env:"CSR_DIR" flag:"csrdir" config:"csr_dir" default:"./csrs"`
	LogDir      string `env:"LOG_DIR" flag:"logdir" config:"log_dir" default:"./logs"`
	DatabaseDir string `env:"DB_DIR" flag:"dbdir" config:"db_dir" default:"./db"`
}

// DefaultConfig returns a new Config with default values
func DefaultConfig() *Config {
	cfg := &Config{}

	// Set default values from struct tags
	t := reflect.TypeOf(*cfg)
	v := reflect.ValueOf(cfg).Elem()

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		defaultVal := field.Tag.Get("default")

		if defaultVal != "" {
			fieldValue := v.Field(i)
			switch fieldValue.Kind() {
			case reflect.String:
				fieldValue.SetString(defaultVal)
			case reflect.Int:
				if intVal, err := strconv.Atoi(defaultVal); err == nil {
					fieldValue.SetInt(int64(intVal))
				}
			case reflect.Bool:
				if boolVal, err := strconv.ParseBool(defaultVal); err == nil {
					fieldValue.SetBool(boolVal)
				}
			case reflect.Float64:
				if floatVal, err := strconv.ParseFloat(defaultVal, 64); err == nil {
					fieldValue.SetFloat(floatVal)
				}
			}
		}
	}

	return cfg
}

// LoadConfigFromFile loads configuration from a file
func (cfg *Config) LoadConfigFromFile(configFile string) error {
	if configFile == "" {
		return nil
	}

	data, err := os.ReadFile(configFile)
	if err != nil {
		return fmt.Errorf("error reading config file: %w", err)
	}

	// Determine file type based on extension
	ext := strings.ToLower(filepath.Ext(configFile))
	switch ext {
	case ".json":
		if err := json.Unmarshal(data, cfg); err != nil {
			return fmt.Errorf("error parsing JSON config: %w", err)
		}
	case ".toml":
		if err := toml.Unmarshal(data, cfg); err != nil {
			return fmt.Errorf("error parsing TOML config: %w", err)
		}
	default:
		return fmt.Errorf("unsupported config file format: %s", ext)
	}

	return nil
}

// LoadFromEnvironment loads configuration from environment variables
func (cfg *Config) LoadFromEnvironment() {
	t := reflect.TypeOf(*cfg)
	v := reflect.ValueOf(cfg).Elem()

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		envName := field.Tag.Get("env")

		if envName != "" {
			envValue := os.Getenv(envName)
			if envValue != "" {
				fieldValue := v.Field(i)
				switch fieldValue.Kind() {
				case reflect.String:
					fieldValue.SetString(envValue)
				case reflect.Int:
					if intVal, err := strconv.Atoi(envValue); err == nil {
						fieldValue.SetInt(int64(intVal))
					}
				case reflect.Bool:
					if boolVal, err := strconv.ParseBool(envValue); err == nil {
						fieldValue.SetBool(boolVal)
					}
				case reflect.Float64:
					if floatVal, err := strconv.ParseFloat(envValue, 64); err == nil {
						fieldValue.SetFloat(floatVal)
					}
				}
			}
		}
	}
}

// LoadDefaults sets up all required directories and ensures proper configurations
func (cfg *Config) LoadDefaults() error {
	// Create required directories if they don't exist
	directories := []string{
		cfg.CertDir,
		cfg.CSRDir,
		cfg.LogDir,
		cfg.DatabaseDir,
	}

	for _, dir := range directories {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", dir, err)
		}
	}

	return nil
}

// Validate checks if the configuration is valid
func (cfg *Config) Validate() error {
	// TODO: Add validation for required fields
	return nil
}

// Duration is a helper for parsing time.Duration from strings
func Duration(val string) (time.Duration, error) {
	return time.ParseDuration(val)
}
