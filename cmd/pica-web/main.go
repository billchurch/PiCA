package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/billchurch/PiCA/internal/ca"
	"github.com/billchurch/PiCA/internal/yubikey"
	"github.com/billchurch/PiCA/web/api"
)

func main() {
	// Define command-line flags
	configFile := flag.String("config", "", "Path to CA config file")
	certFile := flag.String("cert", "", "Path to CA certificate file")
	slotStr := flag.String("slot", "82", "YubiKey PIV slot to use (82-95 for CA slots)")
	port := flag.Int("port", 8080, "Port to listen on")
	webRoot := flag.String("webroot", "./web/html", "Directory containing web files")
	certDir := flag.String("certdir", "./certs", "Directory to store certificates")
	csrDir := flag.String("csrdir", "./csrs", "Directory to store CSRs")
	flag.Parse()

	// Validate inputs
	if *configFile == "" {
		log.Fatal("Config file is required")
	}
	if *certFile == "" {
		log.Fatal("Certificate file is required")
	}

	// Parse YubiKey slot
	slot := yubikey.PIVSlot(0)
	switch *slotStr {
	case "82":
		slot = yubikey.SlotCA1
	case "83":
		slot = yubikey.SlotCA2
	default:
		log.Fatalf("Invalid YubiKey PIV slot: %s", *slotStr)
	}

	// Create CA instance
	ca := ca.NewCA(
		ca.SubCA,
		*configFile,
		"", // Key file not needed when using YubiKey
		*certFile,
	)

	// Create API server
	server := api.NewServer(ca, slot, *certDir, *csrDir)

	// Set up static file serving
	webDir, err := filepath.Abs(*webRoot)
	if err != nil {
		log.Fatalf("Error resolving web root path: %v", err)
	}
	if _, err := os.Stat(webDir); os.IsNotExist(err) {
		log.Fatalf("Web root directory does not exist: %s", webDir)
	}
	http.Handle("/", http.FileServer(http.Dir(webDir)))

	// Start the server
	addr := fmt.Sprintf(":%d", *port)
	log.Printf("Starting web server on %s", addr)
	log.Printf("Web root: %s", webDir)
	if err := server.StartServer(addr); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
