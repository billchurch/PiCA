package config

import (
	"os"
	"testing"
)

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()
	
	// Test a few default values
	if cfg.LogLevel != "info" {
		t.Errorf("Expected LogLevel to be 'info', got '%s'", cfg.LogLevel)
	}
	
	if cfg.WebPort != 8080 {
		t.Errorf("Expected WebPort to be 8080, got %d", cfg.WebPort)
	}
	
	if cfg.CAType != "sub" {
		t.Errorf("Expected CAType to be 'sub', got '%s'", cfg.CAType)
	}
}

func TestLoadFromEnvironment(t *testing.T) {
	// Set environment variables
	os.Setenv("LOG_LEVEL", "debug")
	os.Setenv("WEB_PORT", "9090")
	os.Setenv("CA_TYPE", "root")
	defer func() {
		os.Unsetenv("LOG_LEVEL")
		os.Unsetenv("WEB_PORT")
		os.Unsetenv("CA_TYPE")
	}()
	
	// Create config and load from environment
	cfg := DefaultConfig()
	cfg.LoadFromEnvironment()
	
	// Test if variables were properly loaded
	if cfg.LogLevel != "debug" {
		t.Errorf("Expected LogLevel to be 'debug', got '%s'", cfg.LogLevel)
	}
	
	if cfg.WebPort != 9090 {
		t.Errorf("Expected WebPort to be 9090, got %d", cfg.WebPort)
	}
	
	if cfg.CAType != "root" {
		t.Errorf("Expected CAType to be 'root', got '%s'", cfg.CAType)
	}
}

func TestLoadConfigFromFile(t *testing.T) {
	// Create a temporary JSON config file
	tempFile, err := os.CreateTemp("", "pica-config-*.json")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tempFile.Name())
	
	// Write test configuration to the file
	testConfig := `{
		"log_level": "trace",
		"web_port": 7070,
		"ca_type": "root",
		"provider": "software"
	}`
	
	if _, err := tempFile.Write([]byte(testConfig)); err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}
	tempFile.Close()
	
	// Load from the test file
	cfg := DefaultConfig()
	err = cfg.LoadConfigFromFile(tempFile.Name())
	if err != nil {
		t.Fatalf("Failed to load config file: %v", err)
	}
	
	// Test if variables were properly loaded
	if cfg.LogLevel != "trace" {
		t.Errorf("Expected LogLevel to be 'trace', got '%s'", cfg.LogLevel)
	}
	
	if cfg.WebPort != 7070 {
		t.Errorf("Expected WebPort to be 7070, got %d", cfg.WebPort)
	}
	
	if cfg.CAType != "root" {
		t.Errorf("Expected CAType to be 'root', got '%s'", cfg.CAType)
	}
	
	if cfg.ProviderType != "software" {
		t.Errorf("Expected ProviderType to be 'software', got '%s'", cfg.ProviderType)
	}
}

func TestFullConfigLoad(t *testing.T) {
	// Create a temporary config file
	tempFile, err := os.CreateTemp("", "pica-config-*.json")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tempFile.Name())
	
	// Write test configuration to the file
	testConfig := `{
		"log_level": "debug",
		"web_port": 6060,
		"ca_type": "sub",
		"provider": "software"
	}`
	
	if _, err := tempFile.Write([]byte(testConfig)); err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}
	tempFile.Close()
	
	// Set environment variables (should override file)
	os.Setenv("LOG_LEVEL", "info")
	os.Setenv("CA_TYPE", "root")
	defer func() {
		os.Unsetenv("LOG_LEVEL")
		os.Unsetenv("CA_TYPE")
	}()
	
	// Create mock command-line args (should override env and file)
	args := []string{"--port", "7777", "--provider", "yubikey"}
	
	// Load with all sources
	cfg, err := Load(args, tempFile.Name())
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}
	
	// Test the final priority - command-line > env > file > defaults
	if cfg.LogLevel != "info" { // from env, overrides file
		t.Errorf("Expected LogLevel to be 'info', got '%s'", cfg.LogLevel)
	}
	
	if cfg.WebPort != 7777 { // from args, overrides file
		t.Errorf("Expected WebPort to be 7777, got %d", cfg.WebPort)
	}
	
	if cfg.CAType != "root" { // from env, overrides file
		t.Errorf("Expected CAType to be 'root', got '%s'", cfg.CAType)
	}
	
	if cfg.ProviderType != "yubikey" { // from args, overrides file
		t.Errorf("Expected ProviderType to be 'yubikey', got '%s'", cfg.ProviderType)
	}
}
