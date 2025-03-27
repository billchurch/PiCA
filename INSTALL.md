# PiCA Installation Guide

This guide provides instructions for setting up both the Root CA and Sub CA components of the PiCA system.

## Prerequisites

For running the software directly:
- Go 1.21 or newer
- YubiKey with PIV support
- PCSC daemon (pcscd)
- YubiKey tools (yubico-piv-tool, yubikey-manager)
- CFSSL toolkit

For building custom Raspberry Pi images:
- Raspberry Pi OS (64-bit)
- rpi-image-gen tool
- Raspberry Pi Imager

## Installation Steps

### Option 1: Install on an existing system

1. Clone the repository:
   ```
   git clone https://github.com/billchurch/pica.git
   cd pica
   ```

2. Install dependencies:
   ```
   # Debian/Ubuntu
   sudo apt-get update
   sudo apt-get install -y ca-certificates openssl pcscd pcsc-tools yubikey-manager yubico-piv-tool golang-cfssl
   ```

3. Initialize directories:
   ```
   make init
   ```

4. Build the applications:
   ```
   make build
   ```

5. Initialize your YubiKey for CA use:
   - For Root CA, prepare a YubiKey with a PIV slot (typically 82) for the Root CA key
   - For Sub CA, prepare a YubiKey with a PIV slot for the Sub CA key

### Option 2: Build and use custom Raspberry Pi images

#### For Root CA (Offline)

1. Clone the rpi-image-gen repository:
   ```
   git clone https://github.com/raspberrypi/rpi-image-gen.git
   cd rpi-image-gen
   sudo ./install_deps.sh
   ```

2. Build the Root CA image:
   ```
   cd /path/to/pica
   make build-root-ca-image
   ```

3. Flash the image to an SD card:
   ```
   sudo rpi-imager --cli rpi-image-gen/work/root-ca/artefacts/root-ca.img /dev/mmcblk0
   ```

4. Boot the Raspberry Pi from the SD card and follow the on-screen instructions to set up the Root CA.

#### For Sub CA (Online)

1. Build the Sub CA image:
   ```
   cd /path/to/pica
   make build-sub-ca-image
   ```

2. Flash the image to an SD card:
   ```
   sudo rpi-imager --cli rpi-image-gen/work/sub-ca/artefacts/sub-ca.img /dev/mmcblk0
   ```

3. Boot the Raspberry Pi, configure networking, and follow the on-screen instructions to set up the Sub CA.

### Option 3: Using Docker

1. Run using Docker Compose:
   ```
   docker-compose up -d
   ```

2. Access the web interface at https://localhost:8443

## Post-Installation

1. Create CSR for the Root CA:
   ```
   cfssl genkey -initca configs/cfssl/root-ca-csr.json | cfssljson -bare root-ca
   ```

2. Initialize the Root CA:
   ```
   bin/pica
   ```
   Follow the on-screen instructions to set up the Root CA.

3. Create CSR for the Sub CA:
   ```
   cfssl genkey -initca configs/cfssl/sub-ca-csr.json | cfssljson -bare sub-ca
   ```

4. Sign the Sub CA CSR using the Root CA:
   ```
   bin/pica
   ```
   Use the Root CA management interface to sign the Sub CA certificate.

5. Start the Sub CA web server:
   ```
   bin/pica-web \
     --config ./configs/cfssl/sub-ca-config.json \
     --cert ./certs/sub-ca.pem \
     --port 8080 \
     --webroot ./web/html \
     --certdir ./certs \
     --csrdir ./csrs
   ```
   
   The key parameters are:
   - `--certdir`: Directory where issued certificates are stored and listed in the web interface
   - `--csrdir`: Directory where CSRs are stored when submitted through the web interface

6. Access the Sub CA web interface at http://localhost:8080
