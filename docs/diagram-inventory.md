# PiCA Diagram Inventory

This document provides an inventory of all Mermaid diagrams used throughout the PiCA documentation.

## System Architecture Diagrams

| Diagram Name | Location | Description |
|--------------|----------|-------------|
| System Overview | [architecture-diagrams.md](architecture-diagrams.md) | High-level view of the PiCA system components |
| Certificate Hierarchy | [architecture-diagrams.md](architecture-diagrams.md) | Hierarchy of certificates in the system |
| Software Architecture | [architecture-diagrams.md](architecture-diagrams.md) | Software stack and components |
| Deployment Architecture | [architecture-diagrams.md](architecture-diagrams.md) | Physical deployment topology |
| Project Structure | [architecture-diagrams.md](architecture-diagrams.md) | Code organization and project structure |
| YubiKey Integration | [architecture-diagrams.md](architecture-diagrams.md) | YubiKey integration with the application |

## Process Flow Diagrams

| Diagram Name | Location | Description |
|--------------|----------|-------------|
| Certificate Issuance Process | [architecture-diagrams.md](architecture-diagrams.md) | End-to-end certificate issuance flow |
| Certificate Revocation Process | [architecture-diagrams.md](architecture-diagrams.md) | Process for revoking certificates |
| Root CA Certificate Creation | [certificate-lifecycle.md](certificate-lifecycle.md) | Process of creating the Root CA |
| Sub CA Certificate Creation | [certificate-lifecycle.md](certificate-lifecycle.md) | Process of creating the Sub CA |
| End-Entity Certificate Issuance | [certificate-lifecycle.md](certificate-lifecycle.md) | Detailed flow for issuing certificates |
| Certificate Renewal Process | [certificate-lifecycle.md](certificate-lifecycle.md) | Process for renewing certificates |
| Certificate Status Checking | [certificate-lifecycle.md](certificate-lifecycle.md) | OCSP and CRL checking processes |
| Root CA Renewal Process | [certificate-lifecycle.md](certificate-lifecycle.md) | Process for renewing the Root CA |

## State Diagrams

| Diagram Name | Location | Description |
|--------------|----------|-------------|
| Certificate Lifecycle Overview | [certificate-lifecycle.md](certificate-lifecycle.md) | States in a certificate's lifecycle |
| YubiKey PIN Protection | [yubikey-operations.md](yubikey-operations.md) | States and transitions for PIN protection |
| YubiKey States and Transitions | [yubikey-operations.md](yubikey-operations.md) | Lifecycle states of a YubiKey |

## Component Diagrams

| Diagram Name | Location | Description |
|--------------|----------|-------------|
| YubiKey PIV Slot Allocation | [yubikey-operations.md](yubikey-operations.md) | PIV slot usage in PiCA |
| YubiKey Touch Policy Options | [yubikey-operations.md](yubikey-operations.md) | Touch policies for YubiKey operations |
| YubiKey Backup and Recovery | [yubikey-operations.md](yubikey-operations.md) | Backup strategy for YubiKeys |
| YubiKey in Software Architecture | [yubikey-operations.md](yubikey-operations.md) | YubiKey integration in software |

## Sequence Diagrams

| Diagram Name | Location | Description |
|--------------|----------|-------------|
| Certificate Signing with YubiKey | [yubikey-operations.md](yubikey-operations.md) | Detailed signing flow with YubiKey |
| PIN Protection for Management Key | [yubikey-operations.md](yubikey-operations.md) | PIN protection flow |
| Hardware-Based Attestation | [yubikey-operations.md](yubikey-operations.md) | Attestation process with YubiKey |

## Flow Charts

| Diagram Name | Location | Description |
|--------------|----------|-------------|
| YubiKey Setup Process | [yubikey-operations.md](yubikey-operations.md) | Steps to set up a YubiKey |
| YubiKey Management Workflow | [yubikey-operations.md](yubikey-operations.md) | Lifecycle management of YubiKeys |

## How to Use These Diagrams

These diagrams are created using Mermaid, a JavaScript-based diagramming tool. You can:

1. View them directly in the Markdown files when rendered on GitHub or other Markdown viewers that support Mermaid
2. Copy the Mermaid code to the [Mermaid Live Editor](https://mermaid.live/) for editing
3. Generate images from the diagrams for presentations or documentation

## Updating Diagrams

To update a diagram:

1. Locate the diagram in the appropriate Markdown file
2. Edit the Mermaid code between the triple backticks and "mermaid" tag
3. Test your changes in the [Mermaid Live Editor](https://mermaid.live/)
4. Submit your changes through a pull request

## Diagram Styling

We use a consistent styling convention for our diagrams:

- Root CA components: #f9f fill color
- Sub CA components: #bbf fill color
- YubiKey components: #fdd or #dfd fill colors
- Error states: #ffdddd fill color
- Success states: #ddffdd fill color
- Neutral states: #ddddff fill color

When creating new diagrams, please adhere to these styling conventions for consistency.
