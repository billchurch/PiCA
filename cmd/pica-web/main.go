package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/billchurch/PiCA/internal/ca"
	"github.com/billchurch/PiCA/internal/config"
	"github.com/billchurch/PiCA/internal/crypto"
	"github.com/billchurch/PiCA/internal/yubikey"
	"github.com/billchurch/PiCA/web/api"
)

func main() {
	// Load configuration
	cfg, err := config.Load(os.Args[1:], "")
	if err != nil {
		log.Fatalf("Error loading configuration: %v", err)
	}

	// Validate required settings
	if cfg.CAConfigFile == "" {
		log.Fatal("Config file is required. Use --ca-config flag or CA_CONFIG environment variable.")
	}
	if cfg.CACertFile == "" {
		log.Fatal("Certificate file is required. Use --ca-cert flag or CA_CERT environment variable.")
	}

	// Parse YubiKey slot
	slot := yubikey.PIVSlot(0)
	slotVal, err := strconv.ParseInt(cfg.KeySlot, 16, 64)
	if err != nil {
		log.Fatalf("Invalid YubiKey PIV slot: %s", cfg.KeySlot)
	}
	slot = yubikey.PIVSlot(slotVal)

	// Determine CA type
	caType := ca.SubCA
	if cfg.CAType == "root" {
		caType = ca.RootCA
	}

	// Create CA instance
	caInstance := ca.NewCA(
		caType,
		cfg.CAConfigFile,
		"", // Key file not needed when using YubiKey
		cfg.CACertFile,
	)

	// Set up crypto provider if specified
	if cfg.ProviderType != "" {
		// Force specific provider type
		os.Setenv("PICA_PROVIDER", cfg.ProviderType)
	}

	// Initialize provider
	provider, err := crypto.CreateDefaultProvider()
	if err != nil {
		log.Fatalf("Error creating crypto provider: %v", err)
	}
	defer provider.Close()

	log.Printf("Using crypto provider: %s (Hardware: %t)", provider.Name(), provider.IsHardware())

	// Create API server
	server := api.NewServer(caInstance, slot, cfg.CertDir, cfg.CSRDir)

	// Create required directories
	for _, dir := range []string{cfg.CertDir, cfg.CSRDir, cfg.LogDir} {
		if err := os.MkdirAll(dir, 0755); err != nil {
			log.Fatalf("Error creating directory %s: %v", dir, err)
		}
	}

	// Set up static file serving
	webDir, err := filepath.Abs(cfg.WebRoot)
	if err != nil {
		log.Fatalf("Error resolving web root path: %v", err)
	}
	if _, err := os.Stat(webDir); os.IsNotExist(err) {
		log.Fatalf("Web root directory does not exist: %s", webDir)
	}
	http.Handle("/", http.FileServer(http.Dir(webDir)))

	// Start the server
	addr := fmt.Sprintf(":%d", cfg.WebPort)
	log.Printf("Starting web server on %s", addr)
	log.Printf("Web root: %s", webDir)
	
	// Start with appropriate protocol
	if cfg.EnableHTTPS && cfg.WebTLSCert != "" && cfg.WebTLSKey != "" {
		log.Printf("HTTPS enabled with certificate: %s", cfg.WebTLSCert)
		if err := http.ListenAndServeTLS(addr, cfg.WebTLSCert, cfg.WebTLSKey, nil); err != nil {
			log.Fatalf("Error starting HTTPS server: %v", err)
		}
	} else {
		if err := server.StartServer(addr); err != nil {
			log.Fatalf("Error starting HTTP server: %v", err)
		}
	}
}
