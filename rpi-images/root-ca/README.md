# Root CA Raspberry Pi Image

This directory contains configuration files and assets needed to build a custom Raspberry Pi image for the offline Root CA using rpi-image-gen.

## Requirements

- Raspberry Pi OS 64-bit (Bookworm)
- Minimal installation with only required packages
- No networking services activated by default
- Firewall configured to block all incoming connections

## Build Instructions

1. Clone the rpi-image-gen repository:
   ```
   git clone https://github.com/raspberrypi/rpi-image-gen.git
   cd rpi-image-gen
   sudo ./install_deps.sh
   ```

2. Copy the configuration files from this directory to the appropriate locations in rpi-image-gen.

3. Run the build script:
   ```
   ./build.sh -c /path/to/pica/rpi-images/root-ca/config.ini
   ```

4. Use Raspberry Pi Imager to flash the resulting image to an SD card.

## Usage

1. Boot the Raspberry Pi using the created image
2. The PiCA software will be accessible via the console
3. Insert your YubiKey when prompted to initialize the Root CA

## Security Considerations

- Always keep the Root CA offline
- Store backup keys securely
- Consider using full disk encryption
- Only power on the Root CA when absolutely necessary to sign Sub CA certificates
