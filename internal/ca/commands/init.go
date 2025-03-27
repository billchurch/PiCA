package commands

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/billchurch/pica/internal/ca"
	"github.com/billchurch/pica/internal/crypto"
	"github.com/cloudflare/cfssl/csr"
)

// InitCommand represents the command to initialize a new CA
type InitCommand struct {
	CAType          ca.CAType
	ConfigFile      string
	CSRFile         string
	CertificateFile string
	Slot            crypto.Slot
	Provider        crypto.Provider
	RootCACertFile  string
	RootCAConfigFile string
	Profile         string
}

// NewInitCommand creates a new InitCommand
func NewInitCommand(caType ca.CAType, configFile, csrFile, certFile string, slot crypto.Slot) *InitCommand {
	return &InitCommand{
		CAType:          caType,
		ConfigFile:      configFile,
		CSRFile:         csrFile,
		CertificateFile: certFile,
		Slot:            slot,
	}
}

// Execute initializes a new CA
func (cmd *InitCommand) Execute() error {
	// Read the CSR file
	csrBytes, err := os.ReadFile(cmd.CSRFile)
	if err != nil {
		return fmt.Errorf("error reading CSR file: %w", err)
	}

	// Parse the CSR
	req := csr.CertificateRequest{}
	err = json.Unmarshal(csrBytes, &req)
	if err != nil {
		return fmt.Errorf("error parsing CSR: %w", err)
	}

	// Create or use the crypto provider
	if cmd.Provider == nil {
		provider, err := crypto.CreateDefaultProvider()
		if err != nil {
			return fmt.Errorf("error creating crypto provider: %w", err)
		}
		cmd.Provider = provider
	}

	// Create the CA directory if it doesn't exist
	caDir := filepath.Dir(cmd.CertificateFile)
	if err := os.MkdirAll(caDir, 0700); err != nil {
		return fmt.Errorf("error creating CA directory: %w", err)
	}

	// Initialize the CA based on type
	switch cmd.CAType {
	case ca.RootCA:
		fmt.Println("Initializing Root CA")
		fmt.Println("This operation will generate a self-signed certificate")

		if cmd.Provider.IsHardware() {
			fmt.Println("Using hardware security module:", cmd.Provider.Name())
			fmt.Println("Please ensure your security device is inserted and press Enter to continue...")
			fmt.Scanln()
		} else {
			fmt.Println("Using software-based key storage:", cmd.Provider.Name())
		}

		// Set expiry from the CSR (default to 10 years if not specified)
		expiry := 10 * 365 * 24 * time.Hour
		if req.CA != nil && req.CA.Expiry != "" {
			parsedExpiry, err := time.ParseDuration(req.CA.Expiry)
			if err == nil {
				expiry = parsedExpiry
			}
		}

		// Generate the Root CA certificate
		err := ca.GenerateRootCA(&req, cmd.Provider, cmd.Slot, cmd.CertificateFile, expiry)
		if err != nil {
			return fmt.Errorf("error generating Root CA: %w", err)
		}

		fmt.Println("Root CA initialized successfully!")
		fmt.Println("Certificate saved to:", cmd.CertificateFile)
		return nil

	case ca.SubCA:
		// For Sub CA, we need a Root CA certificate to sign this one
		fmt.Println("Initializing Sub CA")
		fmt.Println("This operation requires:")
		fmt.Println("1. A security module for the Sub CA")
		fmt.Println("2. The Root CA certificate")
		fmt.Println("3. Access to the Root CA security module for signing")

		if cmd.Provider.IsHardware() {
			fmt.Println("Using hardware security module:", cmd.Provider.Name())
			fmt.Println("Please ensure your security device is inserted and press Enter to continue...")
			fmt.Scanln()
		} else {
			fmt.Println("Using software-based key storage:", cmd.Provider.Name())
		}

		// Use the provided Root CA certificate path or fall back to default
		rootCACertFile := "./certs/root-ca.pem"
		if cmd.RootCACertFile != "" {
			rootCACertFile = cmd.RootCACertFile
		}

		// Create a root CA provider (can be the same as the sub CA provider)
		rootProvider := cmd.Provider

		// Set slots for root and sub CA
		rootSlot := crypto.SlotCA1
		subSlot := crypto.SlotCA2

		// Set expiry from the CSR (default to 5 years if not specified)
		expiry := 5 * 365 * 24 * time.Hour
		if req.CA != nil && req.CA.Expiry != "" {
			parsedExpiry, err := time.ParseDuration(req.CA.Expiry)
			if err == nil {
				expiry = parsedExpiry
			}
		}

		// Generate the Sub CA certificate
		err := ca.GenerateSubCA(&req, rootProvider, rootSlot, cmd.Provider, subSlot,
			rootCACertFile, cmd.CertificateFile, expiry)
		if err != nil {
			return fmt.Errorf("error generating Sub CA: %w", err)
		}

		fmt.Println("Sub CA initialized successfully!")
		fmt.Println("Certificate saved to:", cmd.CertificateFile)
		return nil

	default:
		return fmt.Errorf("unknown CA type: %d", cmd.CAType)
	}
}
