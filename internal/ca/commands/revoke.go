package commands

import (
	"fmt"

	"github.com/billchurch/PiCA/internal/ca"
	"github.com/billchurch/PiCA/internal/crypto"
)

// RevokeCommand represents the command to revoke a certificate
type RevokeCommand struct {
	CA           *ca.CA
	SerialNumber string
	Reason       string
	Slot         crypto.Slot
	Provider     crypto.Provider
}

// NewRevokeCommand creates a new RevokeCommand with default provider
func NewRevokeCommand(ca *ca.CA, serialNumber, reason string, slot crypto.Slot) *RevokeCommand {
	return &RevokeCommand{
		CA:           ca,
		SerialNumber: serialNumber,
		Reason:       reason,
		Slot:         slot,
	}
}

// NewRevokeCommandWithProvider creates a new RevokeCommand with a specific provider
func NewRevokeCommandWithProvider(ca *ca.CA, serialNumber, reason string, provider crypto.Provider, slot crypto.Slot) *RevokeCommand {
	return &RevokeCommand{
		CA:           ca,
		SerialNumber: serialNumber,
		Reason:       reason,
		Provider:     provider,
		Slot:         slot,
	}
}

// Execute revokes a certificate and updates the CRL
func (cmd *RevokeCommand) Execute() error {
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
		fmt.Println("Ready to revoke certificate with serial number:", cmd.SerialNumber)
		fmt.Println("Using hardware security module:", cmd.CA.Provider.Name())
		fmt.Println("Please ensure your security device is inserted and press Enter to continue...")
		fmt.Scanln()
	} else {
		fmt.Println("Ready to revoke certificate with serial number:", cmd.SerialNumber)
		fmt.Println("Using software-based key storage:", cmd.CA.Provider.Name())
	}

	// Revoke the certificate
	err := cmd.CA.RevokeCertificate(cmd.SerialNumber)
	if err != nil {
		return fmt.Errorf("error revoking certificate: %w", err)
	}

	fmt.Println("Certificate revoked successfully.")
	return nil
}
