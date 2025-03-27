// Package crypto provides an abstraction layer for cryptographic operations
// with support for both hardware (YubiKey) and software-based implementations.
package crypto

import (
	"crypto"
	"crypto/x509"
)

// ProviderType represents the type of crypto provider
type ProviderType int

const (
	// YubiKeyProviderType represents a YubiKey hardware-based provider
	YubiKeyProviderType ProviderType = iota
	// SoftwareProviderType represents a software-based provider
	SoftwareProviderType
)

// Slot represents a key slot in a provider
type Slot int

// Define common slots that match YubiKey PIV slots for compatibility
const (
	SlotAuthentication Slot = 0x9A
	SlotSignature      Slot = 0x9C
	SlotCardAuth       Slot = 0x9E
	SlotKeyManagement  Slot = 0x9D
	// 0x82-0x95 are retirement slots that can be used for CA keys
	SlotCA1 Slot = 0x82
	SlotCA2 Slot = 0x83
)

// Provider defines the interface for cryptographic operations
// that can be implemented by different providers (hardware or software)
type Provider interface {
	// Type returns the type of the provider
	Type() ProviderType
	
	// Name returns a human-readable name for the provider
	Name() string
	
	// Connect establishes a connection to the provider
	Connect() error
	
	// Close terminates the connection to the provider
	Close() error
	
	// GenerateKey generates a new key pair in the specified slot
	GenerateKey(slot Slot, algorithm string, bits int) error
	
	// GetPublicKey retrieves the public key from a slot
	GetPublicKey(slot Slot) (crypto.PublicKey, error)
	
	// Sign signs data using the private key in the specified slot
	Sign(slot Slot, digest []byte, opts crypto.SignerOpts) ([]byte, error)
	
	// ImportCertificate imports a certificate into a slot
	ImportCertificate(slot Slot, cert *x509.Certificate) error
	
	// GetCertificate retrieves a certificate from a slot
	GetCertificate(slot Slot) (*x509.Certificate, error)
	
	// IsHardware returns true if this is a hardware-based provider
	IsHardware() bool
}

// ProviderFactory is a function that creates a new Provider
type ProviderFactory func(opts map[string]interface{}) (Provider, error)

var (
	// registeredProviders holds the registered provider factories
	registeredProviders = make(map[ProviderType]ProviderFactory)
)

// RegisterProvider registers a provider factory for the given type
func RegisterProvider(providerType ProviderType, factory ProviderFactory) {
	registeredProviders[providerType] = factory
}

// NewProvider creates a new provider of the given type
func NewProvider(providerType ProviderType, opts map[string]interface{}) (Provider, error) {
	if factory, ok := registeredProviders[providerType]; ok {
		return factory(opts)
	}
	return nil, ErrProviderNotFound
}
