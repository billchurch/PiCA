# PiCA Integration Guide

This document provides guidelines for integrating PiCA with existing infrastructure and systems.

## Overview

PiCA can be integrated with various systems to provide certificate services. This guide covers common integration scenarios and best practices.

## Certificate Trust Store Integration

### Adding Root CA to Trust Stores

To get your PiCA Root CA certificate trusted by client systems:

#### Linux Systems

1. Copy the Root CA certificate to the trust store:

   ```bash
   sudo cp root-ca.crt /usr/local/share/ca-certificates/pica-root-ca.crt
   sudo update-ca-certificates
   ```

#### Windows Systems

1. Import the Root CA certificate to the Trusted Root Certification Authorities store:

   ```cmd
   certutil -addstore -f "ROOT" root-ca.crt
   ```

   Or through GUI:
   - Double-click the certificate file
   - Select "Install Certificate"
   - Select "Local Machine"
   - Select "Place all certificates in the following store"
   - Choose "Trusted Root Certification Authorities"
   - Finish the wizard

#### macOS Systems

1. Import the Root CA certificate to the System keychain:

   ```bash
   sudo security add-trusted-cert -d -r trustRoot -k /Library/Keychains/System.keychain root-ca.crt
   ```

#### Web Browsers

- **Firefox**: Import through Settings → Privacy & Security → Certificates → View Certificates → Authorities → Import
- **Chrome/Edge**: These browsers use the system trust store on most platforms

### Distribution Mechanisms

Consider these methods for distributing your Root CA certificate:

- Configuration management tools (Ansible, Puppet, Chef)
- MDM (Mobile Device Management) solutions
- Group Policy for Windows domains
- Internal package repositories

## Web Server Integration

### Apache

1. Configure SSL in Apache to use PiCA-issued certificates:

   ```apache
   <VirtualHost *:443>
     ServerName example.com
     SSLEngine on
     SSLCertificateFile /path/to/server.crt
     SSLCertificateKeyFile /path/to/server.key
     SSLCertificateChainFile /path/to/sub-ca.crt
     # Other directives...
   </VirtualHost>
   ```

### Nginx

1. Configure SSL in Nginx to use PiCA-issued certificates:

   ```nginx
   server {
     listen 443 ssl;
     server_name example.com;
     ssl_certificate /path/to/server.crt;
     ssl_certificate_key /path/to/server.key;
     ssl_trusted_certificate /path/to/sub-ca.crt;
     # Other directives...
   }
   ```

### Certificate Automation

For automated certificate renewal and deployment:

1. Set up scripts to generate CSRs and submit them to PiCA
2. Create hooks to reload services after certificate renewal
3. Consider implementing with cron or systemd timers

## OpenVPN Integration

### Server Configuration

1. Add the following to your OpenVPN server configuration:

   ```bash
   ca /path/to/sub-ca.crt
   cert /path/to/server.crt
   key /path/to/server.key
   ```

### Client Configuration

1. Add the following to your OpenVPN client configuration:

   ```bash
   ca /path/to/sub-ca.crt
   cert /path/to/client.crt
   key /path/to/client.key
   ```

## LDAP/Active Directory Integration

### Using Certificates for LDAPS

1. Issue a server certificate for your LDAP server
2. Configure LDAP to use the certificate:

   For OpenLDAP:

   ```bash
   TLSCACertificateFile /path/to/sub-ca.crt
   TLSCertificateFile /path/to/ldap-server.crt
   TLSCertificateKeyFile /path/to/ldap-server.key
   ```

   For Active Directory, use the Certificates MMC snap-in to import the certificate

### Client Authentication with Certificates

1. Issue client certificates with appropriate subject names
2. Configure LDAP to verify client certificates
3. Map certificates to user accounts in the directory

## Email and S/MIME Integration

### Email Client Configuration

1. Import the Root CA and Sub CA certificates into email clients
2. Import personal certificates and private keys for signing/encryption
3. Configure S/MIME settings in the email client

### Exchange/Email Server Configuration

1. Configure your email server to validate S/MIME signatures using your CA
2. Set up directory lookup for user certificates

## Code Signing Integration

### Windows Code Signing

1. Issue a code signing certificate from your PiCA Sub CA
2. Use signtool to sign executables:

   ```bash
   signtool sign /f certificate.pfx /p password /tr http://timestamp.server.com /td sha256 file.exe
   ```

### JAR Signing

1. Issue a code signing certificate from your PiCA Sub CA
2. Sign JAR files with jarsigner:

   ```bash
   jarsigner -keystore keystore.jks -signedjar signed.jar unsigned.jar alias
   ```

## Docker Content Trust

1. Issue certificates for Docker Content Trust
2. Configure Docker to use your PiCA Sub CA certificate
3. Set up signing processes for container images

## Kubernetes Integration

### TLS for Kubernetes API Server

1. Issue certificates for the Kubernetes API server
2. Configure the API server to use these certificates
3. Distribute the CA certificate to all nodes

### Service Mesh Integration

For integrating with service meshes like Istio:

1. Configure the mesh to use PiCA for certificate issuance
2. Set up automatic certificate rotation
3. Integrate with Kubernetes secrets for certificate storage

## Certificate Issuance API Integration

### API Client Implementation

Create scripts or applications that interact with the PiCA API:

1. Generate a CSR
2. Submit the CSR to the API
3. Retrieve and install the signed certificate

Example Python client:

```python
import requests
import json
from cryptography.hazmat.primitives import serialization
from cryptography.hazmat.primitives.asymmetric import rsa
from cryptography import x509
from cryptography.x509.oid import NameOID
from cryptography.hazmat.primitives import hashes

# Generate key and CSR
private_key = rsa.generate_private_key(
    public_exponent=65537,
    key_size=2048,
)

csr = x509.CertificateSigningRequestBuilder().subject_name(x509.Name([
    x509.NameAttribute(NameOID.COMMON_NAME, u'example.com'),
    x509.NameAttribute(NameOID.ORGANIZATION_NAME, u'Example Org'),
])).sign(private_key, hashes.SHA256())

# Get CSR in PEM format
csr_pem = csr.public_bytes(serialization.Encoding.PEM).decode('utf-8')

# Submit to PiCA API
response = requests.post(
    'https://pica-sub-ca.example.com/api/submit-csr',
    json={
        'csr': csr_pem,
        'profile': 'server'
    },
    headers={'Content-Type': 'application/json'}
)

# Process response
if response.status_code == 200:
    cert_data = response.json()['certificate']
    # Save the certificate
    with open('certificate.pem', 'w') as f:
        f.write(cert_data)
else:
    print(f"Error: {response.text}")
```

## CRL and OCSP Integration

### CRL Distribution

1. Configure web servers to serve CRLs at defined distribution points
2. Set up regular CRL updates and publication
3. Ensure CRL URLs are included in issued certificates

### OCSP Responder

1. Set up an OCSP responder using the PiCA SubCA credentials
2. Configure clients to validate certificates using OCSP
3. Ensure OCSP URLs are included in issued certificates

## Monitoring and Alerting Integration

1. Configure monitoring for certificate expiration
2. Set up alerts for CRL updates and OCSP availability
3. Integrate with existing monitoring solutions like Prometheus, Nagios, or Zabbix

## Best Practices for Integration

1. **Automation**: Automate certificate issuance, renewal, and deployment processes
2. **Inventory**: Maintain an inventory of all issued certificates and their locations
3. **Monitoring**: Monitor certificate expiration and CA service health
4. **Documentation**: Document CA policy and certificate lifecycle management procedures
5. **Testing**: Test certificate issuance and revocation before deploying to production
6. **Disaster Recovery**: Plan for CA compromise or failure scenarios
7. **Security**: Protect private keys and credentials used in automated systems
8. **Validation**: Regularly validate certificate chains and trust relationships
