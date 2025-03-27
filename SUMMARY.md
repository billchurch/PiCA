# PiCA Project Summary

This document provides an overview of the PiCA project structure and components.

## Project Structure

```
PiCA/
├── cmd/                       # Command-line applications
│   ├── pica/                  # Main TUI application
│   └── pica-web/              # Web server application
├── configs/                   # Configuration files
│   └── cfssl/                 # CFSSL configuration files
├── internal/                  # Internal packages
│   ├── ca/                    # Certificate authority implementation
│   │   └── commands/          # CA command implementations
│   ├── ui/                    # Terminal UI components
│   │   └── pages/             # UI page components
│   └── yubikey/               # YubiKey integration
├── pkg/                       # Public packages (empty for now)
├── rpi-images/                # Custom Raspberry Pi image configurations
│   ├── root-ca/               # Root CA image configuration
│   └── sub-ca/                # Sub CA image configuration
├── scripts/                   # Helper scripts
│   └── setup-yubikey.sh       # YubiKey setup script
├── web/                       # Web interface
│   ├── api/                   # API server
│   └── html/                  # Web UI files
├── .gitignore                 # Git ignore file
├── CONTRIBUTING.md            # Contribution guidelines
├── Dockerfile                 # Docker build file
├── docker-compose.yml         # Docker Compose configuration
├── go.mod                     # Go module definition
├── INSTALL.md                 # Installation instructions
├── LICENSE                    # License file
├── Makefile                   # Build automation
└── README.md                  # Project overview
```

## Components

### Certificate Authority (CA)

- **Root CA**: Offline, high-security CA used only to sign Sub CA certificates
- **Sub CA**: Online CA for issuing end-entity certificates

### YubiKey Integration

- Secure storage of private keys in YubiKey PIV slots
- Hardware-backed signing operations

### User Interfaces

- **Terminal UI**: Charm-based TUI for direct management of CAs
- **Web Interface**: Web UI for certificate management and CSR submission

### Raspberry Pi Image Generation

- Custom OS images for Root CA and Sub CA using rpi-image-gen
- Preloaded with all necessary software and configurations

### Deployment Options

- Native installation on Raspberry Pi devices
- Docker containerization for the Sub CA

## Implementation Status

### Completed Components

- Project structure and organization
- Basic UI components for TUI
- Web interface for certificate management
- YubiKey integration framework
- CFSSL configuration files
- Raspberry Pi image configurations
- Build system and Makefile
- Docker containerization
- Installation and setup scripts

### Components Requiring Further Development

- Complete YubiKey PIV integration for key operations
- Full implementation of certificate lifecycle management
- CRL and OCSP responder functionality
- Audit logging and reporting
- Comprehensive testing

## Next Steps

1. Implement YubiKey signing operations
2. Complete the certificate management workflows
3. Add CRL and OCSP support
4. Improve error handling and recovery
5. Add comprehensive logging
6. Write tests for all components
7. Create detailed documentation
8. Build and test the Raspberry Pi images

## Usage

The PiCA system is designed to be used in a two-tier architecture:

1. The Root CA is set up on an air-gapped Raspberry Pi with a YubiKey for key storage
2. The Sub CA is set up on a network-connected Raspberry Pi with a YubiKey for key storage
3. The Root CA signs the Sub CA certificate
4. The Sub CA issues and manages end-entity certificates for services and users

See the `INSTALL.md` file for detailed setup instructions.
