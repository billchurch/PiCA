# Crypto Provider Abstraction Layer

This package provides an abstraction layer for cryptographic operations in PiCA, supporting both hardware (YubiKey) and software-based implementations.

## Overview

The crypto provider abstraction layer allows PiCA to work with or without a YubiKey. This is especially useful for:

- Development and testing environments where YubiKeys may not be available
- Deployment scenarios with different security requirements
- Gracefully degrading to software-based keys when hardware is unavailable

## Provider Types

The package currently supports two provider types:

1. **YubiKeyProvider**: Uses YubiKey hardware for key storage and operations
2. **SoftwareProvider**: Uses software-based keys stored on disk (with appropriate permissions)

## Usage

### Auto-Detection

By default, the system will try to detect the best provider to use:

```go
// Create a provider using auto-detection
provider, err := crypto.CreateDefaultProvider()
if err != nil {
    // Handle error
}
defer provider.Close()

// Use the provider
err = provider.GenerateKey(crypto.SlotCA1, "ECDSA", 384)
// ...
```

### Enforcing a Provider Type

You can force a specific provider type via environment variable:

```bash
# Force software provider
export PICA_PROVIDER=software

# Force YubiKey provider
export PICA_PROVIDER=yubikey
```

Or programmatically:

```go
// Force software provider
provider, err := crypto.CreateProviderFromConfig(map[string]interface{}{
    "type": "software",
    "name": "My Software Provider",
})
```

### Working with Providers

All providers implement the common `Provider` interface:

```go
// Generate a key
err := provider.GenerateKey(crypto.SlotCA1, "ECDSA", 384)

// Get the public key
pubKey, err := provider.GetPublicKey(crypto.SlotCA1)

// Sign data
signature, err := provider.Sign(crypto.SlotCA1, digest, opts)

// Import/export certificates
cert, err := provider.GetCertificate(crypto.SlotCA1)
err = provider.ImportCertificate(crypto.SlotCA1, cert)
```

### Using with the CA Module

The CA module has been updated to work with both provider types:

```go
// Create a CA with the default provider
ca := ca.NewCA(ca.RootCA, configFile, keyFile, certFile)

// Or with a specific provider
provider, _ := crypto.CreateProviderFromConfig(map[string]interface{}{
    "type": "software",
})
ca := ca.NewCAWithProvider(ca.RootCA, configFile, keyFile, certFile, provider, crypto.SlotCA1)
```

## Slots

The crypto provider uses the concept of slots for key storage, which maps directly to YubiKey PIV slots:

- `SlotCA1` (0x82): Recommended for Root CA keys
- `SlotCA2` (0x83): Recommended for Sub CA keys

## Development and Testing

For development and testing purposes, you can use the software provider which stores keys and certificates on disk:

1. Set environment variable: `export PICA_PROVIDER=software`
2. Run your application as normal

The software provider will create and use keys in `~/.pica/keys/` and certificates in `~/.pica/certs/`.

## Security Considerations

- The software provider should be used for development and testing only
- For production use, prefer the YubiKey provider for Root CA operations
- The software provider attempts to use secure permissions but cannot match the security of hardware tokens
