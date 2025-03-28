# PiCA Project Goals

This document outlines the goals and milestones for the PiCA (Raspberry Pi Certificate Authority) project, tracking both accomplished items and future objectives.

## Core Infrastructure

- [x] Set up basic project structure
- [x] Implement Go module and dependency management
- [x] Create Makefile for build automation
- [x] Set up GitHub repository
- [x] Configure development environment with devcontainer
- [x] Create documentation framework
- [x] Basic command-line scaffolding

## Cryptographic Backend

- [x] YubiKey integration for hardware key storage
- [x] Software-based key provider for development
- [x] Provider abstraction layer
- [x] Auto-detection of available providers
- [x] Environment variable override for provider selection
- [x] Support for PIV slot selection
- [ ] Encrypted storage for software-based keys
- [ ] Support for additional HSM types
- [ ] Cloud KMS provider option
- [ ] Key migration between providers
- [ ] Multiple key algorithm support (RSA, ECDSA, Ed25519)
- [ ] Custom certificate extensions

## Certificate Authority Features

- [x] Certificate signing request (CSR) handling
- [x] Certificate generation and signing
- [x] Root CA initialization
- [x] Sub CA initialization and delegation
- [x] CFSSL integration
- [ ] Certificate revocation list (CRL) generation
- [ ] OCSP responder implementation
- [ ] Certificate transparency logging
- [ ] Certificate lifecycle management
- [ ] Automated certificate renewal
- [ ] Scheduled CRL updates
- [ ] Multi-tier CA hierarchies (beyond Root/Sub)
- [ ] Advanced certificate policies
- [ ] Certificate template management

## User Interfaces

- [x] Basic terminal UI framework using Charm
- [x] YubiKey management screens
- [x] Certificate management screens
- [ ] Complete terminal UI experience
- [x] Web interface for CSR submission
- [x] Web interface for certificate issuance
- [ ] Web interface for certificate management
- [ ] Web interface for revocation
- [ ] REST API for programmatic access
- [ ] Responsive design for mobile compatibility
- [ ] User authentication and role-based access control
- [ ] Localization support

## Raspberry Pi Integration

- [ ] Custom Root CA image with minimal attack surface
- [ ] Custom Sub CA image with secure networking
- [ ] Hardware security enhancements
- [ ] Boot security (signed boot, secure boot)
- [ ] Disk encryption
- [ ] Automatic YubiKey detection
- [ ] Optimized performance for Raspberry Pi hardware
- [ ] Power failure safety mechanisms
- [ ] Air-gap management tools for Root CA
- [ ] Pi-specific installation scripts

## Security Features

- [x] PIV slot-based key isolation
- [x] PIN protection for private key operations
- [ ] Secure audit logging
- [ ] Key ceremony documentation and tooling
- [ ] Tamper-evident seals for physical security
- [ ] Intrusion detection
- [ ] Network security hardening for Sub CA
- [ ] Threat modeling documentation
- [ ] Security policy enforcement
- [ ] Regular security scanning integration

## Documentation

- [x] Architecture overview
- [x] YubiKey setup and operations guides
- [x] Provider abstraction documentation
- [x] Basic usage documentation
- [ ] Full user manual
- [ ] Administrator guide
- [ ] Security best practices guide
- [ ] Developer documentation
- [ ] API reference
- [ ] Deployment scenarios
- [ ] Disaster recovery procedures
- [ ] Video tutorials
- [ ] Comprehensive examples

## Testing and Quality Assurance

- [ ] Basic unit tests
- [ ] Integration tests
- [ ] End-to-end tests
- [ ] Performance benchmarks
- [ ] Security testing
- [ ] Continuous integration pipeline
- [ ] Test coverage metrics
- [ ] Fuzz testing for cryptographic operations
- [ ] Cross-platform testing
- [ ] Usability testing

## Deployment and Operations

- [ ] Docker compose deployment option
- [ ] Kubernetes deployment option
- [ ] Backup and restore procedures
- [ ] Monitoring and alerting
- [ ] Performance tuning
- [ ] High availability options
- [ ] Disaster recovery procedures
- [ ] Upgrade and migration paths
- [ ] Integration with external services (LDAP, etc.)
- [ ] Configuration management

## Community and Ecosystem

- [x] Comprehensive contributing guidelines
- [ ] Code of conduct
- [ ] Community forum or discussion platform
- [ ] Example integrations with common services
- [ ] Plugin architecture for extensions
- [ ] Public demo environment
- [ ] Regular release schedule
- [ ] Roadmap publication
- [ ] Package distribution (apt, brew, etc.)
- [ ] Container distribution

## Future Directions

- [ ] Multi-site distributed CA
- [ ] Integration with blockchain for transparency
- [ ] ACME support for Let's Encrypt-like functionality
- [ ] Quantum-resistant cryptography options
- [ ] Zero-trust architecture integration
- [ ] WebAuthn/FIDO support
- [ ] Regulatory compliance tooling (e.g., eIDAS)

## Release Milestones

- [x] Project initialization
- [x] Provider abstraction implementation
- [ ] Alpha release with basic CA functionality
- [ ] Beta release with complete terminal UI
- [ ] 1.0 release with production-ready features
- [ ] 1.x releases with additional features and improvements
- [ ] 2.0 release with advanced ecosystem integration

---

Note: This document will be updated regularly as the project progresses. Checked items indicate completed work, while unchecked items represent planned or in-progress work.
