package ca

import (
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"errors"
	"fmt"
	"math/big"
	"os"
	"path/filepath"
	"time"

	"github.com/cloudflare/cfssl/config"
	"github.com/cloudflare/cfssl/csr"

	"github.com/billchurch/PiCA/internal/crypto"
)

// CAType represents the type of CA we're dealing with
type CAType int

const (
	RootCA CAType = iota
	SubCA
)

// CA represents a certificate authority
type CA struct {
	Type       CAType
	ConfigFile string
	KeyFile    string
	CertFile   string
	Provider   crypto.Provider
	Slot       crypto.Slot
}

// NewCA creates a new CA instance
func NewCA(caType CAType, configFile, keyFile, certFile string) *CA {
	return &CA{
		Type:       caType,
		ConfigFile: configFile,
		KeyFile:    keyFile,
		CertFile:   certFile,
	}
}

// NewCAWithProvider creates a new CA instance with a specific crypto provider
func NewCAWithProvider(caType CAType, configFile, keyFile, certFile string, provider crypto.Provider, slot crypto.Slot) *CA {
	return &CA{
		Type:       caType,
		ConfigFile: configFile,
		KeyFile:    keyFile,
		CertFile:   certFile,
		Provider:   provider,
		Slot:       slot,
	}
}

// InitializeProvider initializes the crypto provider if not already done
func (ca *CA) InitializeProvider() error {
	if ca.Provider != nil {
		// Provider is already initialized
		return nil
	}

	// Create a default provider
	provider, err := crypto.CreateDefaultProvider()
	if err != nil {
		return fmt.Errorf("failed to create default crypto provider: %w", err)
	}

	ca.Provider = provider

	// Set default slot based on CA type
	if ca.Type == RootCA {
		ca.Slot = crypto.SlotCA1
	} else {
		ca.Slot = crypto.SlotCA2
	}

	return nil
}

// GenerateRootCA generates a new root CA certificate
func GenerateRootCA(req *csr.CertificateRequest, provider crypto.Provider, slot crypto.Slot, certFile string, expiry time.Duration) error {
	if provider == nil {
		return errors.New("crypto provider is required")
	}

	// Generate key in the specified slot
	algorithm := "ECDSA"
	bits := 384
	if req.KeyRequest != nil {
		if req.KeyRequest.A == "rsa" {
			algorithm = "RSA"
			bits = req.KeyRequest.Size()
		} else if req.KeyRequest.A == "ecdsa" {
			algorithm = "ECDSA"
			bits = req.KeyRequest.Size()
		}
	}

	fmt.Printf("Generating %s key with size/curve %d\n", algorithm, bits)

	if err := provider.GenerateKey(slot, algorithm, bits); err != nil {
		return fmt.Errorf("failed to generate key: %w", err)
	}

	// Get the public key
	pubKey, err := provider.GetPublicKey(slot)
	if err != nil {
		return fmt.Errorf("failed to get public key: %w", err)
	}

	fmt.Printf("Generated public key of type: %T\n", pubKey)

	// Create a self-signed certificate
	template := &x509.Certificate{
		SerialNumber: big.NewInt(time.Now().UnixNano() / 1000000),
		Subject: pkix.Name{
			CommonName:   req.CN,
			Organization: []string{req.Names[0].O},
			Country:      []string{req.Names[0].C},
			Province:     []string{req.Names[0].ST},
			Locality:     []string{req.Names[0].L},
		},
		NotBefore:             time.Now().Add(-5 * time.Minute),
		NotAfter:              time.Now().Add(expiry),
		KeyUsage:              x509.KeyUsageCertSign | x509.KeyUsageCRLSign,
		ExtKeyUsage:           []x509.ExtKeyUsage{},
		BasicConstraintsValid: true,
		IsCA:                  true,
		MaxPathLenZero:        false,
	}

	// Create directory if it doesn't exist
	if certFile != "" {
		certDir := filepath.Dir(certFile)
		if err := os.MkdirAll(certDir, 0755); err != nil {
			return fmt.Errorf("failed to create certificate directory: %w", err)
		}
	}

	// Create a signer that uses our provider
	signer := &crypto.ProviderSigner{
		Provider:  provider,
		Slot:      slot,
		PublicKey: pubKey,
	}

	fmt.Printf("Created signer with provider: %s\n", provider.Name())

	// Let's add debug info about the template and key
	fmt.Printf("Creating self-signed certificate with key type: %T\n", pubKey)

	// Sign the certificate
	certDER, err := x509.CreateCertificate(rand.Reader, template, template, pubKey, signer)
	if err != nil {
		return fmt.Errorf("failed to create certificate: %w", err)
	}

	// Parse the certificate
	cert, err := x509.ParseCertificate(certDER)
	if err != nil {
		return fmt.Errorf("failed to parse certificate: %w", err)
	}

	// Store the certificate in the provider
	if err := provider.ImportCertificate(slot, cert); err != nil {
		return fmt.Errorf("failed to import certificate: %w", err)
	}

	// Save the certificate to disk if requested
	if certFile != "" {
		certPEM := &pem.Block{
			Type:  "CERTIFICATE",
			Bytes: certDER,
		}

		certBytes := pem.EncodeToMemory(certPEM)
		if err := os.WriteFile(certFile, certBytes, 0644); err != nil {
			return fmt.Errorf("failed to write certificate file: %w", err)
		}
	}

	return nil
}

// GenerateSubCA generates a new sub CA certificate
func GenerateSubCA(req *csr.CertificateRequest, parentProvider crypto.Provider, parentSlot crypto.Slot,
	subProvider crypto.Provider, subSlot crypto.Slot,
	parentCACertFile, certFile string, expiry time.Duration) error {
	if parentProvider == nil || subProvider == nil {
		return errors.New("crypto providers are required")
	}

	// Read the parent CA certificate
	parentCACertBytes, err := os.ReadFile(parentCACertFile)
	if err != nil {
		return fmt.Errorf("failed to read parent CA certificate: %w", err)
	}

	parentCACertBlock, _ := pem.Decode(parentCACertBytes)
	if parentCACertBlock == nil || parentCACertBlock.Type != "CERTIFICATE" {
		return errors.New("failed to decode parent CA certificate")
	}

	parentCACert, err := x509.ParseCertificate(parentCACertBlock.Bytes)
	if err != nil {
		return fmt.Errorf("failed to parse parent CA certificate: %w", err)
	}

	// Generate key for sub CA
	algorithm := "ECDSA"
	bits := 384
	if req.KeyRequest != nil {
		if req.KeyRequest.A == "rsa" {
			algorithm = "RSA"
			bits = req.KeyRequest.Size()
		} else if req.KeyRequest.A == "ecdsa" {
			algorithm = "ECDSA"
			bits = req.KeyRequest.Size()
		}
	}

	if err := subProvider.GenerateKey(subSlot, algorithm, bits); err != nil {
		return fmt.Errorf("failed to generate key: %w", err)
	}

	// Get the public key
	pubKey, err := subProvider.GetPublicKey(subSlot)
	if err != nil {
		return fmt.Errorf("failed to get public key: %w", err)
	}

	// Create a certificate for the sub CA
	template := &x509.Certificate{
		SerialNumber: big.NewInt(time.Now().UnixNano() / 1000000),
		Subject: pkix.Name{
			CommonName:   req.CN,
			Organization: []string{req.Names[0].O},
			Country:      []string{req.Names[0].C},
			Province:     []string{req.Names[0].ST},
			Locality:     []string{req.Names[0].L},
		},
		NotBefore:             time.Now().Add(-5 * time.Minute),
		NotAfter:              time.Now().Add(expiry),
		KeyUsage:              x509.KeyUsageCertSign | x509.KeyUsageCRLSign,
		ExtKeyUsage:           []x509.ExtKeyUsage{},
		BasicConstraintsValid: true,
		IsCA:                  true,
		MaxPathLenZero:        true,
	}

	// Create a signer that uses the parent provider
	signer := &crypto.ProviderSigner{
		Provider:  parentProvider,
		Slot:      parentSlot,
		PublicKey: parentCACert.PublicKey,
	}

	// Let's add debug info about the keys
	fmt.Printf("Creating sub CA certificate with key type: %T, signed by key type: %T\n", pubKey, parentCACert.PublicKey)

	// Sign the certificate
	certDER, err := x509.CreateCertificate(rand.Reader, template, parentCACert, pubKey, signer)
	if err != nil {
		return fmt.Errorf("failed to create certificate: %w", err)
	}

	// Parse the certificate
	cert, err := x509.ParseCertificate(certDER)
	if err != nil {
		return fmt.Errorf("failed to parse certificate: %w", err)
	}

	// Store the certificate in the provider
	if err := subProvider.ImportCertificate(subSlot, cert); err != nil {
		return fmt.Errorf("failed to import certificate: %w", err)
	}

	// Save the certificate to disk if requested
	if certFile != "" {
		certPEM := &pem.Block{
			Type:  "CERTIFICATE",
			Bytes: certDER,
		}

		certBytes := pem.EncodeToMemory(certPEM)
		if err := os.WriteFile(certFile, certBytes, 0644); err != nil {
			return fmt.Errorf("failed to write certificate file: %w", err)
		}
	}

	return nil
}

// SignCertificate signs a CSR using the CA
func (ca *CA) SignCertificate(csrBytes []byte, profile string) ([]byte, error) {
	// Ensure the provider is initialized
	if err := ca.InitializeProvider(); err != nil {
		return nil, err
	}

	// Load CA certificate
	caCertBytes, err := os.ReadFile(ca.CertFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read CA certificate: %w", err)
	}

	caCertBlock, _ := pem.Decode(caCertBytes)
	if caCertBlock == nil || caCertBlock.Type != "CERTIFICATE" {
		return nil, errors.New("failed to decode CA certificate")
	}

	caCert, err := x509.ParseCertificate(caCertBlock.Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse CA certificate: %w", err)
	}

	// Parse CSR
	csrBlock, _ := pem.Decode(csrBytes)
	if csrBlock == nil || csrBlock.Type != "CERTIFICATE REQUEST" {
		return nil, errors.New("failed to decode CSR")
	}

	csr, err := x509.ParseCertificateRequest(csrBlock.Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse CSR: %w", err)
	}

	// Verify CSR signature
	if err := csr.CheckSignature(); err != nil {
		return nil, fmt.Errorf("invalid CSR signature: %w", err)
	}

	// Load config
	configData, err := ca.LoadConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	// Get signing profile
	var signingProfile *config.SigningProfile
	if profile == "" {
		signingProfile = configData.Signing.Default
	} else {
		if p, ok := configData.Signing.Profiles[profile]; ok {
			signingProfile = p
		} else {
			return nil, fmt.Errorf("profile '%s' not found", profile)
		}
	}

	// Create a signer that uses our provider
	signer := &crypto.ProviderSigner{
		Provider:  ca.Provider,
		Slot:      ca.Slot,
		PublicKey: caCert.PublicKey,
	}

	// Create certificate template
	template := &x509.Certificate{
		SerialNumber: big.NewInt(time.Now().UnixNano() / 1000000),
		Subject:      csr.Subject,
		NotBefore:    time.Now().Add(-5 * time.Minute),
		NotAfter:     time.Now().Add(signingProfile.Expiry),
		SubjectKeyId: nil, // Will be calculated
		ExtKeyUsage:  []x509.ExtKeyUsage{},
	}

	// Set key usage based on profile
	ku, eku, _ := signingProfile.Usages()
	template.KeyUsage = ku
	template.ExtKeyUsage = eku

	// Set CA constraints if present
	if signingProfile.CAConstraint.IsCA {
		template.BasicConstraintsValid = true
		template.IsCA = true
		if signingProfile.CAConstraint.MaxPathLen > 0 {
			template.MaxPathLen = signingProfile.CAConstraint.MaxPathLen
		} else if signingProfile.CAConstraint.MaxPathLenZero {
			template.MaxPathLenZero = true
		}
	}

	// Add debug info
	fmt.Printf("Signing certificate with CSR key type: %T, signed by CA key type: %T\n", csr.PublicKey, caCert.PublicKey)

	// Sign the certificate
	certDER, err := x509.CreateCertificate(rand.Reader, template, caCert, csr.PublicKey, signer)
	if err != nil {
		return nil, fmt.Errorf("failed to create certificate: %w", err)
	}

	// Return the PEM-encoded certificate
	certPEM := &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: certDER,
	}

	return pem.EncodeToMemory(certPEM), nil
}

// RevokeCertificate revokes a certificate
func (ca *CA) RevokeCertificate(serialNumber string) error {
	// Ensure the provider is initialized
	if err := ca.InitializeProvider(); err != nil {
		return err
	}

	// In a real implementation, this would create/update a CRL
	// Need to implement CRL management
	return errors.New("CRL management not yet implemented")
}

// LoadConfig loads the CFSSL configuration file
func (ca *CA) LoadConfig() (*config.Config, error) {
	configBytes, err := os.ReadFile(ca.ConfigFile)
	if err != nil {
		return nil, err
	}
	return config.LoadConfig(configBytes)
}
