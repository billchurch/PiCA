# PiCA Usage Guide

This document provides comprehensive instructions for using the PiCA Certificate Authority system.

## Table of Contents

1. [Root CA Operations](#root-ca-operations)
2. [Sub CA Operations](#sub-ca-operations)
3. [Certificate Management](#certificate-management)
4. [YubiKey Operations](#yubikey-operations)
5. [Provider Selection](#provider-selection)
6. [Web Interface](#web-interface)
7. [Command Line Interface](#command-line-interface)
8. [Maintenance Tasks](#maintenance-tasks)
9. [Troubleshooting](#troubleshooting)

## Root CA Operations

### Initializing the Root CA

1. Start the PiCA CLI application:

   ```bash
   ./bin/pica
   ```

2. In the Terminal UI, navigate to the Root CA Management screen using the Tab key.

3. Fill in the required fields:
   - Path to Root CA config file (e.g., `./configs/cfssl/root-ca-config.json`)
   - Path to Root CA CSR file (e.g., `./configs/cfssl/root-ca-csr.json`)
   - Path to save certificate (e.g., `./certs/root-ca.pem`)

4. Press Enter to submit the form.

5. When prompted, insert your YubiKey for the Root CA.

6. Follow the on-screen instructions to complete the initialization.

### Signing Sub CA Certificates

1. Start the PiCA CLI application:

   ```bash
   ./bin/pica
   ```

2. In the Terminal UI, navigate to the Certificate Management screen using the Tab key.

3. Press 's' to select the Sign Certificate option.

4. Fill in the required fields:
   - Path to CSR file (the Sub CA CSR)
   - Path to save certificate
   - Profile (use "subca")
   - Path to Root CA config file

5. Press Enter to submit the form.

6. When prompted, insert your Root CA YubiKey.

7. Follow the on-screen instructions to complete the signing process.

### Creating Certificate Revocation Lists (CRLs)

1. Start the PiCA CLI application.

2. Navigate to the Certificate Management screen.

3. Press 'c' to create or update the CRL.

4. Follow the on-screen instructions.

## Sub CA Operations

### Initializing the Sub CA

1. Start the PiCA CLI application:

   ```bash
   ./bin/pica
   ```

2. In the Terminal UI, navigate to the Sub CA Management screen using the Tab key.

3. Fill in the required fields:
   - Path to Sub CA config file (e.g., `./configs/cfssl/sub-ca-config.json`)
   - Path to Sub CA CSR file (e.g., `./configs/cfssl/sub-ca-csr.json`)
   - Path to save certificate (e.g., `./certs/sub-ca.pem`)
   - Path to Root CA certificate (e.g., `./certs/root-ca.pem`)
   - Path to Root CA config file
   - Profile (e.g., "subca")

4. Press Enter to submit the form.

5. When prompted, insert your YubiKey for the Sub CA.

6. Follow the on-screen instructions to complete the initialization.

### Starting the Web Server

1. Start the PiCA Web server:

   ```bash
   ./bin/pica-web \
     --config ./configs/cfssl/sub-ca-config.json \
     --cert ./certs/sub-ca.pem \
     --port 8080 \
     --webroot ./web/html \
     --certdir ./certs \
     --csrdir ./csrs
   ```

2. Access the web interface at <http://localhost:8080> (or the appropriate IP address/hostname).

### Creating or Updating CRLs

1. Start the PiCA CLI application.

2. Navigate to the Certificate Management screen.

3. Press 'c' to create or update the CRL.

4. Follow the on-screen instructions.

## Certificate Management

### Issuing Certificates via the CLI

1. Start the PiCA CLI application.

2. Navigate to the Certificate Management screen.

3. Press 's' to select the Sign Certificate option.

4. Fill in the required fields:
   - Path to CSR file
   - Path to save certificate
   - Profile (server, client, etc.)
   - Path to CA config file

5. Press Enter to submit the form.

6. When prompted, insert your CA YubiKey.

7. Follow the on-screen instructions to complete the signing process.

### Issuing Certificates via the Web Interface

1. Access the PiCA web interface at <http://localhost:8080> (or the appropriate IP address/hostname).

2. Click on the "Submit CSR" tab.

3. Paste your CSR in PEM format into the text area.

4. Select the appropriate certificate profile from the dropdown.

5. Click the "Submit" button.

6. Download the signed certificate when it appears.

### Revoking Certificates via the CLI

1. Start the PiCA CLI application.

2. Navigate to the Certificate Management screen.

3. Press 'r' to select the Revoke Certificate option.

4. Fill in the required fields:
   - Serial number of certificate to revoke
   - Revocation reason
   - Path to CA config file

5. Press Enter to submit the form.

6. When prompted, insert your CA YubiKey.

7. Follow the on-screen instructions to complete the revocation process.

### Revoking Certificates via the Web Interface

1. Access the PiCA web interface.

2. Click on the "List Certificates" tab.

3. Find the certificate you want to revoke and click the "Revoke" button.

4. Select a revocation reason from the dropdown menu.

5. Click "Revoke" to confirm.

### Listing Certificates via the CLI

1. Start the PiCA CLI application.

2. Navigate to the Certificate Management screen.

3. Press 'l' to list all certificates.

### Listing Certificates via the Web Interface

1. Access the PiCA web interface.

2. Click on the "List Certificates" tab to view all certificates.

## YubiKey Operations

### Setting Up a New YubiKey

Use the provided setup script:

```bash
./scripts/setup-yubikey.sh --slot 82 --name "PiCA Root CA" --pin 123456 --puk 12345678
```

Or follow the manual steps in the [YubiKey Setup Guide](yubikey-setup.md).

### Changing YubiKey PIN

```bash
ykman piv access change-pin -P CURRENT_PIN -n NEW_PIN
```

### Changing YubiKey PUK

```bash
ykman piv access change-puk -P CURRENT_PUK -n NEW_PUK
```

### Changing YubiKey Management Key

```bash
ykman piv access change-management-key -P CURRENT_KEY -n NEW_KEY --protect
```

### Resetting YubiKey PIV Application

Warning: This will delete all certificates and keys!

```bash
ykman piv reset
```

## Provider Selection

PiCA supports multiple cryptographic providers. By default, it will try to use a YubiKey if one is available, and fall back to software-based keys if not.

### Using Hardware-based Keys (YubiKey)

For production use, ensure a YubiKey is connected and use:

```bash
# Force YubiKey provider
export PICA_PROVIDER=yubikey
```

### Using Software-based Keys

For development or testing without a YubiKey:

```bash
# Force software provider
export PICA_PROVIDER=software
```

### Provider Auto-detection

By default, PiCA will automatically detect the best available provider:

1. If a YubiKey is available, it will be used
2. If no YubiKey is available, software-based keys will be used

This behavior can be relied upon for most development workflows.

## Web Interface

### Navigating the Web Interface

The web interface has three main tabs:

- **Submit CSR**: Submit Certificate Signing Requests
- **List Certificates**: View and manage issued certificates
- **Revoke Certificate**: Revoke specific certificates

### Submitting a CSR

1. Click on the "Submit CSR" tab.
2. Paste your CSR in PEM format into the text area.
3. Select the certificate profile from the dropdown.
4. Click "Submit".
5. Download the certificate when it appears.

### Managing Certificates

1. Click on the "List Certificates" tab.
2. View certificate details including subject, serial number, validity dates, and status for all certificates stored in the certificate directory.
3. Use the "Download" button to download certificates.
4. Use the "Revoke" button to revoke valid certificates.

Note: The certificate list displays all valid certificates found in the directory specified by the `--certdir` flag when starting the web server.

### Revoking a Certificate

1. Click on the "Revoke Certificate" tab.
2. Enter the serial number of the certificate to revoke.
3. Select a revocation reason from the dropdown.
4. Click "Revoke" to confirm.

## Command Line Interface

### Basic Navigation

- Use the Tab key to switch between pages
- Press ? for help
- Press q to quit

### Root CA Management Page

- Fill in the form fields for Root CA initialization
- Press Enter to submit

### Sub CA Management Page

- Fill in the form fields for Sub CA initialization
- Press Enter to submit

### Certificate Management Page

- Press s to sign a certificate
- Press r to revoke a certificate
- Press l to list certificates
- Press c to create/update CRL
- Press Esc to cancel current action

## Maintenance Tasks

### Backing Up CA Certificates

```bash
# Backup Root CA certificate
cp ./certs/root-ca.pem /secure/backup/root-ca-$(date +%Y%m%d).pem

# Backup Sub CA certificate
cp ./certs/sub-ca.pem /secure/backup/sub-ca-$(date +%Y%m%d).pem
```

### Backing Up YubiKey Information

Document and securely store the following for each YubiKey:

- YubiKey serial number
- PIN
- PUK
- Management key
- Slot information

### Updating CRLs

CRLs should be updated regularly and after any certificate revocations:

1. Start the PiCA CLI application.
2. Navigate to the Certificate Management screen.
3. Press 'c' to create/update the CRL.
4. Follow the on-screen instructions.

### Checking Certificate Database

Periodically check the certificate database for:

- Expiring certificates
- Unexpected revocations
- Database integrity

### Rotating YubiKey PINs

For security, periodically rotate YubiKey PINs and management keys:

```bash
# Change PIN
ykman piv access change-pin -P CURRENT_PIN -n NEW_PIN

# Change management key
ykman piv access change-management-key -P CURRENT_KEY -n NEW_KEY --protect
```

### System Updates

1. Backup your configuration files:

   ```bash
   mkdir -p backup/configs
   cp -r configs/* backup/configs/
   ```

2. Update the PiCA software:

   ```bash
   git pull
   make clean
   make build
   ```

3. Restart any running services.

## Troubleshooting

### Provider Issues

#### Provider Auto-detection Problems

1. Check the environment to see if a provider is forced:

   ```bash
   echo $PICA_PROVIDER
   ```

2. If using YubiKey provider, ensure the YubiKey is properly inserted.

3. Force a specific provider for testing:

   ```bash
   # Force software provider
   export PICA_PROVIDER=software
   
   # Or force YubiKey provider
   export PICA_PROVIDER=yubikey
   ```

4. Check that the software provider directories exist:

   ```bash
   # Keys directory
   ls -la ~/.pica/keys/
   
   # Certificates directory
   ls -la ~/.pica/certs/
   ```

5. Create the directories if they don't exist:

   ```bash
   mkdir -p ~/.pica/keys ~/.pica/certs
   chmod 700 ~/.pica/keys
   ```

### YubiKey Issues

#### YubiKey Not Detected

1. Check that the YubiKey is properly inserted
2. Verify the PCSC daemon is running:

   ```bash
   systemctl status pcscd
   ```

3. Try re-inserting the YubiKey or using a different USB port
4. Restart the PCSC daemon:

   ```bash
   systemctl restart pcscd
   ```

#### PIN Locked

If the PIN is locked due to too many incorrect attempts:

1. Unlock the PIN using the PUK:

   ```bash
   ykman piv access unblock-pin -P YOUR_PUK -n NEW_PIN
   ```

2. If the PUK is also locked, you'll need to reset the PIV application:

   ```bash
   ykman piv reset
   ```

   Note: This will delete all keys and certificates on the YubiKey!

### Web Server Issues

#### Web Server Won't Start

1. Check for port conflicts:

   ```bash
   sudo lsof -i :8080
   ```

2. Verify the certificate and configuration paths are correct

3. Check the logs for detailed error messages

#### CSR Submission Errors

1. Verify the CSR is in valid PEM format
2. Check that the CSR meets the requirements in the profile
3. Ensure the YubiKey is inserted and functioning

### Certificate Issues

#### Certificate Signing Fails

1. Verify the CSR format is correct
2. Check that the YubiKey is properly inserted
3. Ensure the PIN is correct
4. Verify the CA configuration file exists and is valid

#### Certificate Revocation Fails

1. Verify the serial number is correct
2. Check that the YubiKey is properly inserted
3. Ensure the PIN is correct
4. Verify the CA has permission to update the CRL

### Common Error Messages

#### "Error connecting to YubiKey"

1. Ensure the YubiKey is properly inserted
2. Verify the PCSC daemon is running
3. Try using the YubiKey in another application to confirm functionality

#### "Invalid PIN"

1. Ensure you're using the correct PIN
2. Be aware of PIN attempt limits (typically 3)
3. If locked, use the PUK to reset the PIN

#### "Certificate signing failed"

1. Check the CSR format
2. Verify the YubiKey is functioning
3. Ensure the CA configuration is correct

#### "CRL update failed"

1. Check file permissions
2. Verify the YubiKey is functioning
3. Ensure the CA has the correct configuration for CRL generation

#### "Failed to create provider"

1. Verify that the requested provider type is available
2. If using YubiKey provider, ensure the YubiKey is properly inserted
3. If using software provider, ensure ~/.pica/keys and ~/.pica/certs directories exist
4. Check if PICA_PROVIDER environment variable is set to a valid value ("yubikey" or "software")

#### "No certificates shown in web interface"

1. Verify the `--certdir` flag points to the correct directory containing your certificates
2. Check that certificates have `.pem` or `.crt` file extensions
3. Ensure certificates are in valid PEM format with the correct headers
4. Check file permissions on the certificate directory
5. Refresh the browser page

### Getting Help

If you encounter issues not covered in this troubleshooting guide:

1. Check the logs for detailed error messages
2. Consult the [GitHub repository](https://github.com/billchurch/PiCA) for known issues
3. Open a new issue on GitHub with detailed information about the problem
