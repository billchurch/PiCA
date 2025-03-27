#!/bin/bash
# Script to set up a YubiKey for use with PiCA

set -e

# Default values
KEY_TYPE="ECCP384"
PIN="123456"
PUK="12345678"
MANAGEMENT_KEY="010203040506070801020304050607080102030405060708"
SLOT="82"  # 9A - Default for Root CA
NAME="PiCA Root CA"

print_help() {
    echo "Usage: $0 [options]"
    echo
    echo "Options:"
    echo "  -h, --help                   Show this help message"
    echo "  -t, --type <key_type>        Key type (RSA2048, ECCP256, ECCP384)"
    echo "  -p, --pin <pin>              New PIN (6-8 digits)"
    echo "  -u, --puk <puk>              New PUK (8 digits)"
    echo "  -m, --mgmt-key <key>         New management key (24 byte hex)"
    echo "  -s, --slot <slot>            PIV slot to use (82=9A for Root CA, 83=9B for Sub CA)"
    echo "  -n, --name <n>            Name for the CA"
    echo
}

# Parse command line arguments
while [[ $# -gt 0 ]]; do
    case $1 in
        -h|--help)
            print_help
            exit 0
            ;;
        -t|--type)
            KEY_TYPE="$2"
            shift 2
            ;;
        -p|--pin)
            PIN="$2"
            shift 2
            ;;
        -u|--puk)
            PUK="$2"
            shift 2
            ;;
        -m|--mgmt-key)
            MANAGEMENT_KEY="$2"
            shift 2
            ;;
        -s|--slot)
            SLOT="$2"
            shift 2
            ;;
        -n|--name)
            NAME="$2"
            shift 2
            ;;
        *)
            echo "Unknown option: $1"
            print_help
            exit 1
            ;;
    esac
done

# Validate input
if [[ ! "$KEY_TYPE" =~ ^(RSA2048|ECCP256|ECCP384)$ ]]; then
    echo "Error: Key type must be RSA2048, ECCP256, or ECCP384"
    exit 1
fi

if [[ ! "$PIN" =~ ^[0-9]{6,8}$ ]]; then
    echo "Error: PIN must be 6-8 digits"
    exit 1
fi

if [[ ! "$PUK" =~ ^[0-9]{8}$ ]]; then
    echo "Error: PUK must be 8 digits"
    exit 1
fi

if [[ ! "$MANAGEMENT_KEY" =~ ^[0-9A-Fa-f]{48}$ ]]; then
    echo "Error: Management key must be 24 bytes (48 hex characters)"
    exit 1
fi

if [[ ! "$SLOT" =~ ^(82|83)$ ]]; then
    echo "Error: Slot must be 82 (9A) for Root CA or 83 (9B) for Sub CA"
    exit 1
fi

# Convert slot to hex
SLOT_HEX="0x$SLOT"

echo "========================================================"
echo "Setting up YubiKey for PiCA"
echo "========================================================"
echo "Key Type:       $KEY_TYPE"
echo "PIN:            $PIN"
echo "PUK:            $PUK"
echo "Management Key: $MANAGEMENT_KEY"
echo "Slot:           $SLOT ($SLOT_HEX)"
echo "Name:           $NAME"
echo "========================================================"
echo

# Confirm before proceeding
read -p "This will reset the PIV application on your YubiKey. Continue? (y/n) " -n 1 -r
echo
if [[ ! $REPLY =~ ^[Yy]$ ]]; then
    echo "Operation cancelled."
    exit 0
fi

echo "Please insert your YubiKey now..."
sleep 2

# Check if YubiKey is connected
if ! ykman piv info >/dev/null 2>&1; then
    echo "Error: No YubiKey detected."
    exit 1
fi

echo "YubiKey detected."
echo

# Reset PIV application
echo "Resetting PIV application..."
ykman piv reset -f

# Set PIN, PUK, and management key
echo "Setting PIN, PUK, and management key..."
ykman piv access change-pin -P 123456 -n "$PIN"
ykman piv access change-puk -P 12345678 -n "$PUK"
ykman piv access change-management-key -P 010203040506070801020304050607080102030405060708 -n "$MANAGEMENT_KEY" --protect

# Generate key and self-signed certificate
echo "Generating key and self-signed certificate..."
ykman piv generate-key --algorithm="$KEY_TYPE" "$SLOT_HEX" - | \
ykman piv generate-certificate "$SLOT_HEX" - \
    --subject "/CN=$NAME/O=PiCA Certificate Authority" \
    --valid-days 3650 \
    -P "$PIN"

echo "YubiKey setup complete."
echo
echo "Important: Keep the PIN, PUK, and management key in a safe place."
echo "PIN:            $PIN"
echo "PUK:            $PUK"
echo "Management Key: $MANAGEMENT_KEY"
echo

# Export the certificate
echo "Exporting certificate..."
CERT_FILENAME="${NAME// /_}.crt"
ykman piv export-certificate "$SLOT_HEX" "$CERT_FILENAME"
echo "Certificate exported to $CERT_FILENAME"

echo "YubiKey is now ready for use with PiCA."
