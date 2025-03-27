# YubiKey Setup for PiCA

This document provides detailed instructions for setting up YubiKey devices for use with the PiCA Certificate Authority system.

## Prerequisites

- YubiKey 5 Series (or newer) with PIV support
- PC/Workstation with YubiKey Manager installed
- PiCA software installed

## YubiKey Overview

The YubiKey provides secure hardware-based key storage through its PIV (Personal Identity Verification) application. PiCA uses the following PIV slots for CA keys:

- **Slot 9A (82)**: Used for the Root CA
- **Slot 9B (83)**: Used for the Sub CA

## Security Considerations

Before setting up your YubiKeys, consider the following security practices:

1. Use dedicated YubiKeys for the Root CA and Sub CA
2. Store the Root CA YubiKey in a secure location when not in use
3. Use strong, unique PINs and PUKs for each YubiKey
4. Document and securely store the PINs, PUKs, and management keys
5. Consider using PIN protection for the management key

## Setup Process

### Option 1: Using the Automated Script

PiCA includes a script to automate the YubiKey setup process:

```bash
# Set up YubiKey for Root CA
./scripts/setup-yubikey.sh --slot 82 --name "PiCA Root CA" --pin 123456 --puk 12345678

# Set up YubiKey for Sub CA
./scripts/setup-yubikey.sh --slot 83 --name "PiCA Sub CA" --pin 123456 --puk 12345678
```

### Option 2: Manual Setup

#### Reset the PIV Application

1. Insert your YubiKey
2. Open a terminal
3. Reset the PIV application:

   ```bash
   ykman piv reset
   ```

#### Set PIN, PUK, and Management Key

1. Change the PIN from default (123456):

   ```bash
   ykman piv access change-pin -P 123456 -n YOUR_NEW_PIN
   ```

2. Change the PUK from default (12345678):

   ```bash
   ykman piv access change-puk -P 12345678 -n YOUR_NEW_PUK
   ```

3. Change the management key from default (010203040506070801020304050607080102030405060708):

   ```bash
   ykman piv access change-management-key -P 010203040506070801020304050607080102030405060708 -n YOUR_NEW_MANAGEMENT_KEY --protect
   ```

   Note: The `--protect` flag requires the PIN to be entered when using the management key.

#### Generate Keys and Certificates

1. For Root CA (Slot 9A/82):

   ```bash
   # Generate key
   ykman piv generate-key --algorithm ECCP384 0x9a root-ca-public-key.pem

   # Generate self-signed certificate
   ykman piv generate-certificate 0x9a root-ca-public-key.pem \
     --subject "/CN=PiCA Root CA/O=PiCA Certificate Authority" \
     --valid-days 3650
   ```

2. For Sub CA (Slot 9B/83):

   ```bash
   # Generate key
   ykman piv generate-key --algorithm ECCP384 0x9b sub-ca-public-key.pem

   # Generate CSR for Sub CA
   ykman piv generate-csr 0x9b sub-ca-public-key.pem sub-ca.csr \
     --subject "/CN=PiCA Sub CA/O=PiCA Certificate Authority"
   ```

3. Sign the Sub CA CSR with the Root CA:
   This step will be done using the PiCA software after initial setup.

#### Export Certificates

1. Export the Root CA certificate:

   ```bash
   ykman piv export-certificate 0x9a root-ca.crt
   ```

2. After the Sub CA certificate is signed and imported, export it:

   ```bash
   ykman piv export-certificate 0x9b sub-ca.crt
   ```

## Verifying Setup

After setup, verify your YubiKey configuration:

```bash
ykman piv info
```

You should see information about the PIV application, including:

- PIN and PUK retry counters
- Certificates in slots 9A and/or 9B

## Troubleshooting

### PIN/PUK Locked

If you enter the wrong PIN too many times (default 3), it will be locked. Use the PUK to unlock it:

```bash
ykman piv access unblock-pin -P YOUR_PUK -n NEW_PIN
```

If you enter the wrong PUK too many times, you'll need to reset the PIV application, which will delete all keys and certificates.

### Management Key Issues

If you forget the management key but have PIN access:

1. If management key is PIN protected, you can use PIN instead
2. Otherwise, you'll need to reset the PIV application

### Key Generation Errors

If you encounter errors during key generation:

1. Ensure the YubiKey is properly inserted
2. Check that the PCSC daemon is running
3. Verify you have proper permissions to access the YubiKey
4. Try another USB port or cable

## Backup Considerations

### What to Back Up

For each YubiKey, securely document:

1. PIN
2. PUK
3. Management key
4. Exported public certificate

### What Cannot Be Backed Up

The private keys generated on the YubiKey **cannot be exported**. This is a security feature that ensures the keys remain secure in the hardware.

### Recovery Planning

In case of YubiKey loss or failure:

1. For Root CA: Generate a new Root CA and distribute to all relying parties
2. For Sub CA: Generate a new Sub CA and get it signed by the Root CA

## Best Practices

1. **Physical Security**: Store YubiKeys securely when not in use
2. **PIN Management**: Use strong PINs and change them periodically
3. **Documentation**: Document the setup process and store credentials securely
4. **Testing**: Test certificate signing and revocation processes after setup
5. **Inventory**: Maintain an inventory of YubiKeys used for CA operations

## Advanced Configuration

### Hardware-based Attestation

YubiKey supports hardware-based attestation to verify that keys were generated on the device:

```bash
# Generate key with attestation
ykman piv generate-key --algorithm ECCP384 0x9a key.pem --attestation

# Export attestation certificate
ykman piv export-certificate f9 attestation.pem
```

### Touch Policy

You can configure the YubiKey to require touch for private key operations:

```bash
ykman piv generate-key --algorithm ECCP384 0x9a key.pem --touch-policy always
```

Options for touch policy:

- `always`: Touch required for every operation
- `cached`: Touch required once, then cached for 15 seconds
- `never`: No touch required (default)
