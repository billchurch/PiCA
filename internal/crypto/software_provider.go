package crypto

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"
	"path/filepath"
	"sync"
)

// SoftwareProvider implements the Provider interface using software-based keys
type SoftwareProvider struct {
	name         string
	keys         map[Slot]crypto.PrivateKey
	certificates map[Slot]*x509.Certificate
	keyDir       string
	certDir      string
	connected    bool
	mutex        sync.RWMutex
}

// NewSoftwareProvider creates a new software-based key provider
func NewSoftwareProvider(opts map[string]interface{}) (Provider, error) {
	// Get the directory for storing keys and certificates
	var keyDir, certDir string
	
	if dir, ok := opts["directory"].(string); ok && dir != "" {
		keyDir = filepath.Join(dir, "keys")
		certDir = filepath.Join(dir, "certs")
	} else {
		// Default to storing keys in the user's home directory
		homeDir, err := os.UserHomeDir()
		if err != nil {
			homeDir = "."
		}
		
		keyDir = filepath.Join(homeDir, ".pica", "keys")
		certDir = filepath.Join(homeDir, ".pica", "certs")
	}
	
	name := "Software Provider"
	if n, ok := opts["name"].(string); ok && n != "" {
		name = n
	}
	
	// Create directories if they don't exist
	if err := os.MkdirAll(keyDir, 0700); err != nil {
		return nil, fmt.Errorf("failed to create key directory: %w", err)
	}
	
	if err := os.MkdirAll(certDir, 0700); err != nil {
		return nil, fmt.Errorf("failed to create certificate directory: %w", err)
	}
	
	return &SoftwareProvider{
		name:         name,
		keys:         make(map[Slot]crypto.PrivateKey),
		certificates: make(map[Slot]*x509.Certificate),
		keyDir:       keyDir,
		certDir:      certDir,
		connected:    false,
	}, nil
}

// Type returns the type of the provider
func (p *SoftwareProvider) Type() ProviderType {
	return SoftwareProviderType
}

// Name returns a human-readable name for the provider
func (p *SoftwareProvider) Name() string {
	return p.name
}

// Connect establishes a connection to the provider
func (p *SoftwareProvider) Connect() error {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	
	// Load any existing keys and certificates from the filesystem
	if err := p.loadKeysAndCertificates(); err != nil {
		return fmt.Errorf("failed to load keys and certificates: %w", err)
	}
	
	p.connected = true
	return nil
}

// Close terminates the connection to the provider
func (p *SoftwareProvider) Close() error {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	
	p.connected = false
	return nil
}

// GenerateKey generates a new key pair in the specified slot
func (p *SoftwareProvider) GenerateKey(slot Slot, algorithm string, bits int) error {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	
	if !p.connected {
		return ErrNotConnected
	}
	
	var privateKey crypto.PrivateKey
	var err error
	
	switch algorithm {
	case "RSA":
		privateKey, err = rsa.GenerateKey(rand.Reader, bits)
	case "ECDSA":
		var curve elliptic.Curve
		switch bits {
		case 256:
			curve = elliptic.P256()
		case 384:
			curve = elliptic.P384()
		case 521:
			curve = elliptic.P521()
		default:
			return fmt.Errorf("unsupported ECDSA key size: %d", bits)
		}
		privateKey, err = ecdsa.GenerateKey(curve, rand.Reader)
	default:
		return fmt.Errorf("unsupported algorithm: %s", algorithm)
	}
	
	if err != nil {
		return fmt.Errorf("failed to generate key: %w", err)
	}
	
	// Store the key in memory
	p.keys[slot] = privateKey
	
	// Save the key to disk
	return p.saveKey(slot, privateKey)
}

// GetPublicKey retrieves the public key from a slot
func (p *SoftwareProvider) GetPublicKey(slot Slot) (crypto.PublicKey, error) {
	p.mutex.RLock()
	defer p.mutex.RUnlock()
	
	if !p.connected {
		return nil, ErrNotConnected
	}
	
	privateKey, ok := p.keys[slot]
	if !ok {
		return nil, ErrKeyNotFound
	}
	
	switch key := privateKey.(type) {
	case *rsa.PrivateKey:
		return &key.PublicKey, nil
	case *ecdsa.PrivateKey:
		return &key.PublicKey, nil
	default:
		return nil, fmt.Errorf("unsupported key type: %T", privateKey)
	}
}

// Sign signs data using the private key in the specified slot
func (p *SoftwareProvider) Sign(slot Slot, digest []byte, opts crypto.SignerOpts) ([]byte, error) {
	p.mutex.RLock()
	defer p.mutex.RUnlock()
	
	if !p.connected {
		return nil, ErrNotConnected
	}
	
	privateKey, ok := p.keys[slot]
	if !ok {
		return nil, ErrKeyNotFound
	}
	
	switch key := privateKey.(type) {
	case *rsa.PrivateKey:
		return key.Sign(rand.Reader, digest, opts)
	case *ecdsa.PrivateKey:
		// For X.509 certificate operations, input is a tbs (to-be-signed) block that needs
		// to be signed directly with ecdsa.SignASN1
		fmt.Printf("ECDSA signing input length: %d\n", len(digest))
		
		// For X.509 certificate signing operations, we need to use ASN.1 DER encoding
		// and we need to pass the raw digest directly to SignASN1
		signature, err := ecdsa.SignASN1(rand.Reader, key, digest)
		if err != nil {
			fmt.Printf("ECDSA signing failed: %v\n", err)
			return nil, fmt.Errorf("failed to sign with ECDSA: %w", err)
		}
		
		fmt.Printf("ECDSA signature produced (length: %d)\n", len(signature))
		return signature, nil
	default:
		return nil, fmt.Errorf("unsupported key type: %T", privateKey)
	}
}

// ImportCertificate imports a certificate into a slot
func (p *SoftwareProvider) ImportCertificate(slot Slot, cert *x509.Certificate) error {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	
	if !p.connected {
		return ErrNotConnected
	}
	
	if cert == nil {
		return ErrInvalidCertificate
	}
	
	// Store the certificate in memory
	p.certificates[slot] = cert
	
	// Save the certificate to disk
	return p.saveCertificate(slot, cert)
}

// GetCertificate retrieves a certificate from a slot
func (p *SoftwareProvider) GetCertificate(slot Slot) (*x509.Certificate, error) {
	p.mutex.RLock()
	defer p.mutex.RUnlock()
	
	if !p.connected {
		return nil, ErrNotConnected
	}
	
	cert, ok := p.certificates[slot]
	if !ok {
		return nil, ErrCertNotFound
	}
	
	return cert, nil
}

// IsHardware returns false for software-based provider
func (p *SoftwareProvider) IsHardware() bool {
	return false
}

// Helper methods

// loadKeysAndCertificates loads keys and certificates from disk
func (p *SoftwareProvider) loadKeysAndCertificates() error {
	// Load keys
	keyFiles, err := filepath.Glob(filepath.Join(p.keyDir, "slot_*.key"))
	if err != nil {
		return fmt.Errorf("failed to list key files: %w", err)
	}
	
	for _, keyFile := range keyFiles {
		// Extract slot number from filename
		var slot Slot
		_, err := fmt.Sscanf(filepath.Base(keyFile), "slot_%x.key", &slot)
		if err != nil {
			continue
		}
		
		// Read and parse the key
		keyData, err := os.ReadFile(keyFile)
		if err != nil {
			continue
		}
		
		block, _ := pem.Decode(keyData)
		if block == nil {
			continue
		}
		
		var privateKey crypto.PrivateKey
		
		switch block.Type {
		case "RSA PRIVATE KEY":
			privateKey, err = x509.ParsePKCS1PrivateKey(block.Bytes)
		case "EC PRIVATE KEY":
			privateKey, err = x509.ParseECPrivateKey(block.Bytes)
		case "PRIVATE KEY":
			privateKey, err = x509.ParsePKCS8PrivateKey(block.Bytes)
		default:
			continue
		}
		
		if err != nil {
			continue
		}
		
		p.keys[slot] = privateKey
	}
	
	// Load certificates
	certFiles, err := filepath.Glob(filepath.Join(p.certDir, "slot_*.crt"))
	if err != nil {
		return fmt.Errorf("failed to list certificate files: %w", err)
	}
	
	for _, certFile := range certFiles {
		// Extract slot number from filename
		var slot Slot
		_, err := fmt.Sscanf(filepath.Base(certFile), "slot_%x.crt", &slot)
		if err != nil {
			continue
		}
		
		// Read and parse the certificate
		certData, err := os.ReadFile(certFile)
		if err != nil {
			continue
		}
		
		block, _ := pem.Decode(certData)
		if block == nil || block.Type != "CERTIFICATE" {
			continue
		}
		
		cert, err := x509.ParseCertificate(block.Bytes)
		if err != nil {
			continue
		}
		
		p.certificates[slot] = cert
	}
	
	return nil
}

// saveKey saves a private key to disk
func (p *SoftwareProvider) saveKey(slot Slot, privateKey crypto.PrivateKey) error {
	var keyPEM *pem.Block
	
	switch key := privateKey.(type) {
	case *rsa.PrivateKey:
		keyPEM = &pem.Block{
			Type:  "RSA PRIVATE KEY",
			Bytes: x509.MarshalPKCS1PrivateKey(key),
		}
	case *ecdsa.PrivateKey:
		keyBytes, err := x509.MarshalECPrivateKey(key)
		if err != nil {
			return fmt.Errorf("failed to marshal EC private key: %w", err)
		}
		keyPEM = &pem.Block{
			Type:  "EC PRIVATE KEY",
			Bytes: keyBytes,
		}
	default:
		return fmt.Errorf("unsupported key type: %T", privateKey)
	}
	
	// Write the key to disk
	filename := filepath.Join(p.keyDir, fmt.Sprintf("slot_%x.key", slot))
	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return fmt.Errorf("failed to open key file: %w", err)
	}
	defer file.Close()
	
	if err := pem.Encode(file, keyPEM); err != nil {
		return fmt.Errorf("failed to write key to file: %w", err)
	}
	
	return nil
}

// saveCertificate saves a certificate to disk
func (p *SoftwareProvider) saveCertificate(slot Slot, cert *x509.Certificate) error {
	// Write the certificate to disk
	filename := filepath.Join(p.certDir, fmt.Sprintf("slot_%x.crt", slot))
	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return fmt.Errorf("failed to open certificate file: %w", err)
	}
	defer file.Close()
	
	certPEM := &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: cert.Raw,
	}
	
	if err := pem.Encode(file, certPEM); err != nil {
		return fmt.Errorf("failed to write certificate to file: %w", err)
	}
	
	return nil
}

// Register the provider
func init() {
	RegisterProvider(SoftwareProviderType, NewSoftwareProvider)
}