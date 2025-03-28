package commands

import (
	"encoding/pem"
	"fmt"
	"os"
	"strings"

	"github.com/billchurch/PiCA/internal/ca"
	"github.com/billchurch/PiCA/internal/crypto"
)

// SignCommand represents the command to sign a certificate
type SignCommand struct {
	CA       *ca.CA
	CSRFile  string
	CertFile string
	Profile  string
	Slot     crypto.Slot
	Provider crypto.Provider
}

// NewSignCommand creates a new SignCommand with default provider
func NewSignCommand(ca *ca.CA, csrFile, certFile, profile string, slot crypto.Slot) *SignCommand {
	return &SignCommand{
		CA:       ca,
		CSRFile:  csrFile,
		CertFile: certFile,
		Profile:  profile,
		Slot:     slot,
	}
}

// NewSignCommandWithProvider creates a new SignCommand with a specific provider
func NewSignCommandWithProvider(ca *ca.CA, csrFile, certFile, profile string, provider crypto.Provider, slot crypto.Slot) *SignCommand {
	return &SignCommand{
		CA:       ca,
		CSRFile:  csrFile,
		CertFile: certFile,
		Profile:  profile,
		Provider: provider,
		Slot:     slot,
	}
}

// Execute signs a certificate
func (cmd *SignCommand) Execute() error {
	// Read CSR
	csrBytes, err := os.ReadFile(cmd.CSRFile)
	if err != nil {
		return fmt.Errorf("error reading CSR file: %w", err)
	}

	// Parse PEM
	csrBlock, _ := pem.Decode(csrBytes)
	if csrBlock == nil || csrBlock.Type != "CERTIFICATE REQUEST" {
		return fmt.Errorf("failed to parse CSR PEM block")
	}

	// Load config
	configData, err := cmd.CA.LoadConfig()
	if err != nil {
		return fmt.Errorf("error loading config: %w", err)
	}

	// Validate profile
	if _, ok := configData.Signing.Profiles[cmd.Profile]; !ok && cmd.Profile != "" && cmd.Profile != "default" {
		// List available profiles
		profiles := make([]string, 0, len(configData.Signing.Profiles))
		for name := range configData.Signing.Profiles {
			profiles = append(profiles, name)
		}

		return fmt.Errorf("profile '%s' not found. Available profiles: %s",
			cmd.Profile, strings.Join(profiles, ", "))
	}

	// Use either the provided provider or initialize the CA's provider
	if cmd.Provider != nil {
		cmd.CA.Provider = cmd.Provider
		cmd.CA.Slot = cmd.Slot
	} else {
		if err := cmd.CA.InitializeProvider(); err != nil {
			return fmt.Errorf("error initializing provider: %w", err)
		}
	}

	if cmd.CA.Provider.IsHardware() {
		fmt.Println("Ready to sign certificate.")
		fmt.Println("Using hardware security module:", cmd.CA.Provider.Name())
		fmt.Println("Please ensure your security device is inserted and press Enter to continue...")
		fmt.Scanln()
	} else {
		fmt.Println("Ready to sign certificate.")
		fmt.Println("Using software-based key storage:", cmd.CA.Provider.Name())
	}

	// Sign the certificate
	certPEM, err := cmd.CA.SignCertificate(csrBytes, cmd.Profile)
	if err != nil {
		return fmt.Errorf("error signing certificate: %w", err)
	}

	// Write the certificate to file if requested
	if cmd.CertFile != "" {
		if err := os.WriteFile(cmd.CertFile, certPEM, 0644); err != nil {
			return fmt.Errorf("error writing certificate file: %w", err)
		}
		fmt.Println("Certificate signed and saved to:", cmd.CertFile)
	} else {
		fmt.Println("Certificate signed successfully.")
	}

	return nil
}
