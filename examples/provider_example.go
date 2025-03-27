package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/billchurch/pica/internal/crypto"
)

func main() {
	// Set environment variable to force provider type if needed
	// Comment out the following line to use auto-detection
	// os.Setenv("PICA_PROVIDER", "software")

	// Print detected provider type
	providerType := crypto.GetPreferredProviderType()
	fmt.Printf("Preferred provider type: %s\n", crypto.GetProviderNameByType(providerType))
	
	// Create default provider (will use auto-detection)
	provider, err := crypto.CreateDefaultProvider()
	if err != nil {
		fmt.Printf("Error creating provider: %v\n", err)
		os.Exit(1)
	}
	defer provider.Close()
	
	fmt.Printf("Using provider: %s (Hardware: %t)\n", provider.Name(), provider.IsHardware())
	
	// Generate a key in slot CA1
	slot := crypto.SlotCA1
	fmt.Printf("Generating key in slot %X...\n", slot)
	
	err = provider.GenerateKey(slot, "ECDSA", 256)
	if err != nil {
		fmt.Printf("Error generating key: %v\n", err)
		os.Exit(1)
	}
	
	// Get the public key
	pubKey, err := provider.GetPublicKey(slot)
	if err != nil {
		fmt.Printf("Error getting public key: %v\n", err)
		os.Exit(1)
	}
	
	fmt.Printf("Generated key successfully: %T\n", pubKey)
	
	// Create a signer from the provider and slot
	signer, err := crypto.CreateProviderSigner(provider, slot)
	if err != nil {
		fmt.Printf("Error creating signer: %v\n", err)
		os.Exit(1)
	}
	
	fmt.Printf("Created signer: %T\n", signer)
	
	// Now demonstrate signing something
	dataToSign := []byte("Hello, PiCA!")
	fmt.Printf("Signing data: %s\n", string(dataToSign))
	
	// In a real application, you'd use a proper hash function
	signature, err := signer.Sign(nil, dataToSign, nil)
	if err != nil {
		fmt.Printf("Error signing data: %v\n", err)
		os.Exit(1)
	}
	
	fmt.Printf("Signed data successfully. Signature length: %d bytes\n", len(signature))
	
	// Generate a self-signed certificate for testing
	fmt.Println("\nGenerating a self-signed test certificate...")
	
	// Create output directory
	outDir := "test-certs"
	if err := os.MkdirAll(outDir, 0755); err != nil {
		fmt.Printf("Error creating output directory: %v\n", err)
		os.Exit(1)
	}
	
	certFile := filepath.Join(outDir, "test-cert.pem")
	
	// Example values for a test certificate
	subject := map[string]string{
		"CN": "Test Certificate",
		"O":  "PiCA Example",
		"C":  "US",
	}
	
	// Create a basic certificate (simplified for example)
	// In a real application, you'd use something like ca.GenerateRootCA()
	fmt.Printf("Certificate will be saved to: %s\n", certFile)
	fmt.Println("This is just a placeholder - in a real application, you would use the CA implementation")
	
	fmt.Println("\nProvider abstraction layer test completed successfully!")
}
