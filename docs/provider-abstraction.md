# Crypto Providers in PiCA

## Overview

PiCA implements a provider abstraction layer for cryptographic operations that allows the system to work with or without a YubiKey. This flexible architecture enhances the development experience and provides deployment options for different security requirements.

## Components Implemented

1. **Provider Interface**: A common interface for all crypto providers
   - `provider.go`: Defines the Provider interface and related types

2. **Provider Implementations**:
   - `yubikey_provider.go`: Implementation using YubiKey hardware
   - `software_provider.go`: Implementation using software-based keys

3. **Factory and Detection**:
   - `factory.go`: Functions to create providers
   - `detect.go`: Auto-detection and environment variable handling

4. **Support Classes**:
   - `signer.go`: Implementation of crypto.Signer using our providers
   - `errors.go`: Common errors for the package

5. **CA Integration**:
   - Updates to `ca.go` and command implementations to use the provider abstraction

6. **Documentation and Examples**:
   - `README.md`: Documentation for the crypto package
   - `examples/provider_example.go`: Example usage

## Key Features

1. **Auto-detection**: Automatically chooses YubiKey if available, falls back to software
2. **Environment Control**: Can force provider type via `PICA_PROVIDER` environment variable
3. **Graceful Fallback**: Attempts YubiKey first, gracefully falls back to software provider
4. **Common Interface**: All providers use the same interface, simplifying code
5. **Seamless Integration**: CA module functions with either provider type

## Provider Types

PiCA currently supports two provider types:

1. **YubiKeyProvider**: Uses YubiKey hardware for key storage and operations
2. **SoftwareProvider**: Uses software-based keys stored on disk

## Provider Interface

All providers implement the common `Provider` interface:

```go
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
```

## Usage

The system supports two modes of operation:

### Hardware Mode (YubiKey)

When a YubiKey is available or specifically requested, the system will:

- Use the YubiKey for all cryptographic operations
- Store keys securely in YubiKey PIV slots
- Prompt for YubiKey insertion when needed

### Software Mode

When no YubiKey is available or software mode is requested, the system will:

- Generate and use software-based keys
- Store keys on disk (in `~/.pica/keys/`)
- Store certificates on disk (in `~/.pica/certs/`)

### Controlling Provider Mode

Set the `PICA_PROVIDER` environment variable:

- `export PICA_PROVIDER=yubikey` - Force YubiKey provider
- `export PICA_PROVIDER=software` - Force software provider
- Not set - Auto-detect (YubiKey if available, otherwise software)

### Slots

The YubiKey provider uses the following PIV slots:

- `SlotCA1` (0x82/9A): Recommended for Root CA keys
- `SlotCA2` (0x83/9B): Recommended for Sub CA keys

## Security Considerations

- Software provider should be used for development and testing only
- Keys stored by the software provider are protected with file permissions
- For production use, YubiKey provider is still recommended for Root CA operations

## Implementation Notes

1. The software provider mimics YubiKey slot behavior
2. Both providers implement the same interface so code doesn't need to be aware of which is in use
3. The provider abstraction is transparent to higher-level code

## Usage Examples

### YubiKey Provider Example

```go
// Force YubiKey provider
os.Setenv("PICA_PROVIDER", "yubikey")

// Create provider
provider, err := crypto.CreateDefaultProvider()
if err != nil {
    log.Fatalf("Failed to create provider: %v", err)
}
defer provider.Close()

// Generate a key
err = provider.GenerateKey(crypto.SlotCA1, "ECDSA", 384)
if err != nil {
    log.Fatalf("Failed to generate key: %v", err)
}

// Get the public key
pubKey, err := provider.GetPublicKey(crypto.SlotCA1)
if err != nil {
    log.Fatalf("Failed to get public key: %v", err)
}
```

### Software Provider Example

```go
// Force software provider
os.Setenv("PICA_PROVIDER", "software")

// Create provider
provider, err := crypto.CreateDefaultProvider()
if err != nil {
    log.Fatalf("Failed to create provider: %v", err)
}
defer provider.Close()

// Generate a key
err = provider.GenerateKey(crypto.SlotCA1, "ECDSA", 384)
if err != nil {
    log.Fatalf("Failed to generate key: %v", err)
}

// Get the public key
pubKey, err := provider.GetPublicKey(crypto.SlotCA1)
if err != nil {
    log.Fatalf("Failed to get public key: %v", err)
}
```

Note that the code is identical regardless of which provider is being used.

## Future Enhancements

1. Encryption for software-stored keys
2. Support for additional HSM types
3. Enhanced key protection mechanisms
4. Cloud KMS integration
5. More thorough testing, especially for failover scenarios
6. Complete integration with CRL generation and OCSP functionality
