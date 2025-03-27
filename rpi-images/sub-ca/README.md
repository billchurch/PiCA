# Sub CA Raspberry Pi Image

This directory contains configuration files and assets needed to build a custom Raspberry Pi image for the online Sub CA using rpi-image-gen.

## Requirements

- Raspberry Pi OS 64-bit (Bookworm)
- Minimal installation with essential network services
- Firewall configured to allow only necessary connections
- Web server for certificate issuance interface

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
   ./build.sh -c /path/to/pica/rpi-images/sub-ca/config.ini
   ```

4. Use Raspberry Pi Imager to flash the resulting image to an SD card.

## Usage

1. Boot the Raspberry Pi using the created image
2. Configure network settings
3. The PiCA software will be accessible via the console
4. The web interface will be available at https://<ip-address>:8443
5. Insert your YubiKey when prompted to initialize the Sub CA

## Security Considerations

- Keep the Sub CA on a secured network
- Regularly update the system
- Monitor access logs
- Consider implementing network security measures like IDS/IPS
- Consider using full disk encryption
