package crypto

import (
	"crypto"
	"crypto/x509"
	"fmt"

	"github.com/billchurch/PiCA/internal/yubikey"
)

// YubiKeyProvider implements the Provider interface using a YubiKey
type YubiKeyProvider struct {
	name      string
	yubikey   *yubikey.YubiKey
	connected bool
}

// NewYubiKeyProvider creates a new YubiKey-based provider
func NewYubiKeyProvider(opts map[string]interface{}) (Provider, error) {
	name := "YubiKey Provider"
	if n, ok := opts["name"].(string); ok && n != "" {
		name = n
	}

	return &YubiKeyProvider{
		name:      name,
		yubikey:   nil,
		connected: false,
	}, nil
}

// Type returns the type of the provider
func (p *YubiKeyProvider) Type() ProviderType {
	return YubiKeyProviderType
}

// Name returns a human-readable name for the provider
func (p *YubiKeyProvider) Name() string {
	return p.name
}

// Connect establishes a connection to the provider
func (p *YubiKeyProvider) Connect() error {
	// Connect to the YubiKey
	yk, err := yubikey.Connect()
	if err != nil {
		return fmt.Errorf("failed to connect to YubiKey: %w", err)
	}

	p.yubikey = yk
	p.connected = true
	return nil
}

// Close terminates the connection to the provider
func (p *YubiKeyProvider) Close() error {
	if p.yubikey != nil {
		err := p.yubikey.Close()
		p.yubikey = nil
		p.connected = false
		return err
	}

	p.connected = false
	return nil
}

// GenerateKey generates a new key pair in the specified slot
func (p *YubiKeyProvider) GenerateKey(slot Slot, algorithm string, bits int) error {
	if !p.connected || p.yubikey == nil {
		return ErrNotConnected
	}

	// Convert Provider Slot to YubiKey PIVSlot
	pivSlot := yubikey.PIVSlot(slot)

	// For now, just pass through to the YubiKey
	// In a real implementation, we'd handle the algorithm and bits parameters
	return p.yubikey.GenerateKey(pivSlot)
}

// GetPublicKey retrieves the public key from a slot
func (p *YubiKeyProvider) GetPublicKey(slot Slot) (crypto.PublicKey, error) {
	if !p.connected || p.yubikey == nil {
		return nil, ErrNotConnected
	}

	// Convert Provider Slot to YubiKey PIVSlot
	pivSlot := yubikey.PIVSlot(slot)

	return p.yubikey.GetPublicKey(pivSlot)
}

// Sign signs data using the private key in the specified slot
func (p *YubiKeyProvider) Sign(slot Slot, digest []byte, opts crypto.SignerOpts) ([]byte, error) {
	if !p.connected || p.yubikey == nil {
		return nil, ErrNotConnected
	}

	// Convert Provider Slot to YubiKey PIVSlot
	pivSlot := yubikey.PIVSlot(slot)

	// The YubiKey.Sign method doesn't take SignerOpts, so ignore it for now
	// In a real implementation, we'd handle the SignerOpts appropriately
	return p.yubikey.Sign(pivSlot, digest)
}

// ImportCertificate imports a certificate into a slot
func (p *YubiKeyProvider) ImportCertificate(slot Slot, cert *x509.Certificate) error {
	if !p.connected || p.yubikey == nil {
		return ErrNotConnected
	}

	// Convert Provider Slot to YubiKey PIVSlot
	pivSlot := yubikey.PIVSlot(slot)

	return p.yubikey.ImportCertificate(pivSlot, cert)
}

// GetCertificate retrieves a certificate from a slot
func (p *YubiKeyProvider) GetCertificate(slot Slot) (*x509.Certificate, error) {
	if !p.connected || p.yubikey == nil {
		return nil, ErrNotConnected
	}

	// Convert Provider Slot to YubiKey PIVSlot
	pivSlot := yubikey.PIVSlot(slot)

	return p.yubikey.GetCertificate(pivSlot)
}

// IsHardware returns true for YubiKey-based provider
func (p *YubiKeyProvider) IsHardware() bool {
	return true
}

// Register the provider
func init() {
	RegisterProvider(YubiKeyProviderType, NewYubiKeyProvider)
}
