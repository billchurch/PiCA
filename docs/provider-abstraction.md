# PiCA Crypto Provider Abstraction Layer

## Overview

We've implemented a provider abstraction layer for cryptographic operations in PiCA that allows the system to work with or without a YubiKey. This enhancement is particularly useful for development, testing, and deployment flexibility.

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

## Usage

The system now supports two modes of operation:

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

## Security Considerations

- Software provider should be used for development and testing only
- Keys stored by the software provider are protected with file permissions
- For production use, YubiKey provider is still recommended for Root CA operations

## Implementation Notes

1. The software provider mimics YubiKey slot behavior
2. Both providers implement the same interface so code doesn't need to be aware of which is in use
3. The provider abstraction is transparent to higher-level code

## Next Steps

1. Implement more robust key protection for the software provider (e.g., encryption)
2. Add more thorough testing, especially for failover scenarios
3. Consider supporting additional HSM types in the future
4. Complete integration with CRL generation and OCSP functionality
