// This is an example of how to properly create a self-signed ECDSA certificate
// which can help diagnose issues with the main implementation
package main

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"io"
	"math/big"
	"os"
	"time"
)

// Simple implementation of crypto.Signer
type ECDSASigner struct {
	PrivateKey *ecdsa.PrivateKey
}

func (s *ECDSASigner) Public() crypto.PublicKey {
	return &s.PrivateKey.PublicKey
}

func (s *ECDSASigner) Sign(rand io.Reader, digest []byte, opts crypto.SignerOpts) ([]byte, error) {
	// Directly use ECDSA to sign the digest
	// CreateCertificate expects DER-encoded ASN.1 signature
	return ecdsa.SignASN1(rand, s.PrivateKey, digest)
}

func main() {
	// 1. Generate a new ECDSA key pair
	fmt.Println("Generating ECDSA key pair...")
	privateKey, err := ecdsa.GenerateKey(elliptic.P384(), rand.Reader)
	if err != nil {
		fmt.Printf("Failed to generate key: %v\n", err)
		os.Exit(1)
	}

	publicKey := &privateKey.PublicKey
	fmt.Printf("Key generated. Curve: %s\n", publicKey.Curve.Params().Name)

	// 2. Create a self-signed certificate template
	serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128)
	serialNumber, err := rand.Int(rand.Reader, serialNumberLimit)
	if err != nil {
		fmt.Printf("Failed to generate serial number: %v\n", err)
		os.Exit(1)
	}

	template := x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			Organization: []string{"Example Organization"},
			CommonName:   "Example Root CA",
		},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().Add(365 * 24 * time.Hour),
		KeyUsage:              x509.KeyUsageCertSign | x509.KeyUsageCRLSign,
		BasicConstraintsValid: true,
		IsCA:                  true,
	}

	// 3. Create the signer
	signer := &ECDSASigner{
		PrivateKey: privateKey,
	}

	// 4. Create the certificate
	fmt.Println("Creating self-signed certificate...")
	derBytes, err := x509.CreateCertificate(rand.Reader, &template, &template, publicKey, signer)
	if err != nil {
		fmt.Printf("Failed to create certificate: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Certificate created successfully. Size: %d bytes\n", len(derBytes))

	// 5. Save certificate and key to files
	// Certificate
	certOut, err := os.Create("example_cert.pem")
	if err != nil {
		fmt.Printf("Failed to open cert.pem for writing: %v\n", err)
		os.Exit(1)
	}
	pem.Encode(certOut, &pem.Block{Type: "CERTIFICATE", Bytes: derBytes})
	certOut.Close()
	fmt.Println("Certificate saved to example_cert.pem")

	// Private key
	keyBytes, err := x509.MarshalPKCS8PrivateKey(privateKey)
	if err != nil {
		fmt.Printf("Failed to marshal private key: %v\n", err)
		os.Exit(1)
	}
	keyOut, err := os.Create("example_key.pem")
	if err != nil {
		fmt.Printf("Failed to open key.pem for writing: %v\n", err)
		os.Exit(1)
	}
	pem.Encode(keyOut, &pem.Block{Type: "PRIVATE KEY", Bytes: keyBytes})
	keyOut.Close()
	fmt.Println("Private key saved to example_key.pem")

	fmt.Println("Done!")
}
