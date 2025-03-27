package crypto

import (
	"crypto"
	"crypto/ecdsa"
	"encoding/hex"
	"fmt"
	"io"
)

// ProviderSigner implements crypto.Signer interface using a Provider
type ProviderSigner struct {
	Provider  Provider
	Slot      Slot
	PublicKey crypto.PublicKey
}

// Public returns the public key associated with the signer
func (s *ProviderSigner) Public() crypto.PublicKey {
	return s.PublicKey
}

// Sign signs the digest with the private key in the provider
func (s *ProviderSigner) Sign(rand io.Reader, digest []byte, opts crypto.SignerOpts) ([]byte, error) {
	if s.Provider == nil {
		return nil, fmt.Errorf("provider not initialized")
	}
	
	// Debug information about the signing request
	fmt.Printf("Signing request: data length = %d bytes\n", len(digest))
	if len(digest) > 0 {
		fmt.Printf("Data prefix: %s\n", hex.EncodeToString(digest[:min(len(digest), 16)]))
	}
	
	// Determine key type for debugging
	switch s.PublicKey.(type) {
	case *ecdsa.PublicKey:
		fmt.Println("Using ECDSA public key")
	}
	
	// Use the provider to sign the digest
	signature, err := s.Provider.Sign(s.Slot, digest, opts)
	if err != nil {
		return nil, fmt.Errorf("provider signing failed: %w", err)
	}
	
	// Debug information about the result
	fmt.Printf("Signature generated successfully (length: %d bytes)\n", len(signature))
	
	return signature, nil
}

// CreateProviderSigner creates a ProviderSigner for the given provider and slot
func CreateProviderSigner(provider Provider, slot Slot) (*ProviderSigner, error) {
	if provider == nil {
		return nil, fmt.Errorf("provider is required")
	}
	
	// Get the public key from the provider
	pubKey, err := provider.GetPublicKey(slot)
	if err != nil {
		return nil, fmt.Errorf("failed to get public key: %w", err)
	}
	
	return &ProviderSigner{
		Provider:  provider,
		Slot:      slot,
		PublicKey: pubKey,
	}, nil
}

// min returns the smaller of a and b
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
